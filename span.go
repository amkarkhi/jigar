package jigar

import (
	"context"
	"encoding/json"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// Span is a thin wrapper around an OpenTelemetry span that:
//
//   - is safe to use when tracing is disabled (all methods become no-ops),
//   - exposes convenience setters for the common attribute types,
//   - lets you marshal arbitrary values to JSON for inspection in the UI.
//
// Always pair Start with span.End, typically with defer.
type Span struct {
	span    trace.Span
	enabled bool
}

// Start begins a new span on the package's default tracer (set by Init).
// If the context isn't marked for tracing or Init wasn't called, the
// returned Span is a no-op and the context is returned unchanged.
func Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, *Span) {
	return Default().Start(ctx, name, opts...)
}

// Start begins a new span. If the context isn't marked for tracing (see
// middleware or WithShouldTrace), the returned Span is a no-op. A nil
// receiver is safe and also yields a no-op span.
func (t *Tracer) Start(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, *Span) {
	if t == nil || !ShouldTrace(ctx) {
		return ctx, &Span{}
	}
	ctx, s := t.tracer.Start(ctx, name, opts...)
	return ctx, &Span{span: s, enabled: true}
}

// End completes the span. No-op if tracing is disabled or the span is nil.
func (s *Span) End(opts ...trace.SpanEndOption) {
	if s != nil && s.enabled {
		s.span.End(opts...)
	}
}

// AddError records err on the span and marks the span as errored. No-op
// when err is nil or tracing is disabled.
func (s *Span) AddError(err error) {
	if s == nil || !s.enabled || err == nil {
		return
	}
	s.span.SetStatus(codes.Error, err.Error())
	s.span.RecordError(err)
}

// EndWithError records err and ends the span. Equivalent to AddError + End.
func (s *Span) EndWithError(err error) {
	s.AddError(err)
	s.End()
}

// SetString sets a string attribute on the span.
func (s *Span) SetString(key, value string) {
	if s != nil && s.enabled {
		s.span.SetAttributes(attribute.String(key, value))
	}
}

// SetInt sets an int64 attribute on the span.
func (s *Span) SetInt(key string, value int64) {
	if s != nil && s.enabled {
		s.span.SetAttributes(attribute.Int64(key, value))
	}
}

// SetFloat sets a float64 attribute on the span.
func (s *Span) SetFloat(key string, value float64) {
	if s != nil && s.enabled {
		s.span.SetAttributes(attribute.Float64(key, value))
	}
}

// SetBool sets a bool attribute on the span.
func (s *Span) SetBool(key string, value bool) {
	if s != nil && s.enabled {
		s.span.SetAttributes(attribute.Bool(key, value))
	}
}

// SetJSON JSON-marshals value and stores it as a string attribute. Useful
// for dumping a struct into a span for debugging. If marshaling fails the
// call is a no-op.
func (s *Span) SetJSON(key string, value any) {
	if s == nil || !s.enabled {
		return
	}
	b, err := json.Marshal(value)
	if err != nil {
		return
	}
	s.span.SetAttributes(attribute.String(key, string(b)))
}

// SetAttributes sets one or more raw OpenTelemetry attributes on the span.
// Use this when none of the typed helpers fit.
func (s *Span) SetAttributes(kv ...attribute.KeyValue) {
	if s != nil && s.enabled {
		s.span.SetAttributes(kv...)
	}
}

// Underlying returns the wrapped OpenTelemetry span. May be nil when
// tracing is disabled.
func (s *Span) Underlying() trace.Span {
	if s == nil {
		return nil
	}
	return s.span
}

// SpanID returns the span ID as a hex string. Returns an empty string if
// the span is nil, tracing is disabled, or the span context is invalid.
func (s *Span) SpanID() string {
	if s == nil || !s.enabled {
		return ""
	}
	spanContext := s.span.SpanContext()
	if !spanContext.IsValid() {
		return ""
	}
	return spanContext.SpanID().String()
}

// TraceID returns the trace ID as a hex string. Returns an empty string if
// the span is nil, tracing is disabled, or the span context is invalid.
func (s *Span) TraceID() string {
	if s == nil || !s.enabled {
		return ""
	}
	spanContext := s.span.SpanContext()
	if !spanContext.IsValid() {
		return ""
	}
	return spanContext.TraceID().String()
}

// TraceInfo returns both the trace ID and span ID as hex strings for easy
// linking to trace backends like Jaeger. Returns empty strings if the span
// is nil, tracing is disabled, or the span context is invalid.
func (s *Span) TraceInfo() (traceID, spanID string) {
	if s == nil || !s.enabled {
		return "", ""
	}
	spanContext := s.span.SpanContext()
	if !spanContext.IsValid() {
		return "", ""
	}
	return spanContext.TraceID().String(), spanContext.SpanID().String()
}
