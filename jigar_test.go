package jigar_test

import (
	"context"
	"errors"
	"testing"

	"github.com/amkarkhi/jigar"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"gotest.tools/v3/assert"
)

func newTestTracer(t *testing.T) (*jigar.Tracer, *tracetest.InMemoryExporter) {
	t.Helper()
	exp := tracetest.NewInMemoryExporter()
	tr, err := jigar.New(context.Background(),
		jigar.WithService("test"),
		jigar.WithExporter(exp),
	)
	assert.NilError(t, err)
	t.Cleanup(func() { _ = tr.Shutdown(context.Background()) })
	return tr, exp
}

func TestNew_RequiresService(t *testing.T) {
	_, err := jigar.New(context.Background(), jigar.WithEndpoint("localhost:4317"))
	assert.ErrorIs(t, err, jigar.ErrServiceRequired)
}

func TestNew_RequiresEndpointOrExporter(t *testing.T) {
	_, err := jigar.New(context.Background(), jigar.WithService("test"))
	assert.ErrorIs(t, err, jigar.ErrEndpointRequired)
}

func TestStart_NoOpWithoutShouldTrace(t *testing.T) {
	tr, exp := newTestTracer(t)
	_, span := tr.Start(context.Background(), "noop")
	span.SetString("k", "v")
	span.End()
	assert.NilError(t, forceFlush(tr))
	assert.Equal(t, 0, len(exp.GetSpans()))
}

func TestStart_RecordsAttributesWhenSampled(t *testing.T) {
	tr, exp := newTestTracer(t)
	ctx := jigar.WithShouldTrace(context.Background(), true)
	_, span := tr.Start(ctx, "work")
	span.SetString("user.id", "u1")
	span.SetInt("rows", 3)
	span.SetBool("cache.hit", false)
	span.SetFloat("dur.ms", 1.5)
	span.SetJSON("body", map[string]int{"n": 1})
	span.End()

	assert.NilError(t, forceFlush(tr))
	spans := exp.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, "work", spans[0].Name)

	got := map[string]any{}
	for _, kv := range spans[0].Attributes {
		got[string(kv.Key)] = kv.Value.AsInterface()
	}
	assert.Equal(t, "u1", got["user.id"])
	assert.Equal(t, int64(3), got["rows"])
	assert.Equal(t, false, got["cache.hit"])
	assert.Equal(t, 1.5, got["dur.ms"])
	assert.Equal(t, `{"n":1}`, got["body"])
}

func TestEndWithError(t *testing.T) {
	tr, exp := newTestTracer(t)
	ctx := jigar.WithShouldTrace(context.Background(), true)
	_, span := tr.Start(ctx, "fail")
	span.EndWithError(errors.New("boom"))

	assert.NilError(t, forceFlush(tr))
	spans := exp.GetSpans()
	assert.Equal(t, 1, len(spans))
	assert.Equal(t, 1, len(spans[0].Events))
	assert.Equal(t, "exception", spans[0].Events[0].Name)
}

func TestDefaultBeforeInitIsNil(t *testing.T) {
	// Default may be set by another test; just assert the type behaves as
	// a no-op when nil.
	var nilTracer *jigar.Tracer
	_, span := nilTracer.Start(context.Background(), "x")
	span.SetString("k", "v")
	span.End() // must not panic
}

func forceFlush(tr *jigar.Tracer) error {
	if p, ok := tr.Provider().(*sdktrace.TracerProvider); ok {
		return p.ForceFlush(context.Background())
	}
	return nil
}
