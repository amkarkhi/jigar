package jigar

import (
	"context"
	"fmt"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

// Tracer is a configured OpenTelemetry tracer with an OTLP exporter and
// the helpers in this package. Create one with New or Init.
type Tracer struct {
	cfg      Config
	tracer   trace.Tracer
	provider *sdktrace.TracerProvider
}

var (
	defaultMu     sync.RWMutex
	defaultTracer *Tracer
)

// New builds a Tracer from the given options without touching package-level
// state. Use this when you want explicit control of the tracer's lifetime
// or when wiring multiple tracers in tests.
//
// The returned Tracer's TracerProvider is also installed as the OpenTelemetry
// global TracerProvider so that third-party instrumentations pick it up.
// Call Shutdown to flush and close.
func New(ctx context.Context, opts ...Option) (*Tracer, error) {
	cfg := Config{
		Protocol: ProtocolGRPC,
		Insecure: true,
	}
	for _, o := range opts {
		o(&cfg)
	}
	if cfg.Service == "" {
		return nil, ErrServiceRequired
	}
	if cfg.Exporter == nil && cfg.Endpoint == "" {
		return nil, ErrEndpointRequired
	}
	if cfg.TracerName == "" {
		cfg.TracerName = cfg.Service
	}

	exporter := cfg.Exporter
	if exporter == nil {
		var err error
		exporter, err = buildExporter(ctx, cfg)
		if err != nil {
			return nil, fmt.Errorf("jigar: create OTLP exporter: %w", err)
		}
	}

	res, err := buildResource(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("jigar: build resource: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Tracer{
		cfg:      cfg,
		tracer:   tp.Tracer(cfg.TracerName),
		provider: tp,
	}, nil
}

// Init builds a Tracer and stores it as the package default used by the
// package-level Start, Default, Shutdown and the middleware sub-packages.
// Call this once at process startup.
//
// Init is safe to call concurrently with reads of Default but should not be
// called more than once; subsequent calls return the existing tracer and
// ignore the new options.
func Init(ctx context.Context, opts ...Option) (*Tracer, error) {
	defaultMu.Lock()
	if defaultTracer != nil {
		t := defaultTracer
		defaultMu.Unlock()
		return t, nil
	}
	defaultMu.Unlock()

	t, err := New(ctx, opts...)
	if err != nil {
		return nil, err
	}
	defaultMu.Lock()
	defaultTracer = t
	defaultMu.Unlock()
	return t, nil
}

// Default returns the tracer installed by Init, or nil if Init has not been
// called. The package-level Start uses this.
func Default() *Tracer {
	defaultMu.RLock()
	defer defaultMu.RUnlock()
	return defaultTracer
}

// Shutdown flushes pending spans and shuts down the underlying provider.
// Safe to call on a nil receiver.
func (t *Tracer) Shutdown(ctx context.Context) error {
	if t == nil || t.provider == nil {
		return nil
	}
	return t.provider.Shutdown(ctx)
}

// Provider returns the underlying OpenTelemetry TracerProvider. Useful when
// wiring third-party instrumentations that take a provider explicitly.
func (t *Tracer) Provider() trace.TracerProvider {
	if t == nil {
		return nil
	}
	return t.provider
}

// Config returns a copy of the configuration the tracer was built with.
func (t *Tracer) Config() Config {
	if t == nil {
		return Config{}
	}
	return t.cfg
}

func buildExporter(ctx context.Context, cfg Config) (sdktrace.SpanExporter, error) {
	switch cfg.Protocol {
	case ProtocolHTTP:
		opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.Endpoint)}
		if cfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		return otlptrace.New(ctx, otlptracehttp.NewClient(opts...))
	case ProtocolGRPC, "":
		opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.Endpoint)}
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		return otlptrace.New(ctx, otlptracegrpc.NewClient(opts...))
	default:
		return nil, fmt.Errorf("jigar: unknown protocol %q", cfg.Protocol)
	}
}

func buildResource(_ context.Context, cfg Config) (*resource.Resource, error) {
	attrs := []attribute.KeyValue{semconv.ServiceName(cfg.Service)}
	if cfg.ServiceVersion != "" {
		attrs = append(attrs, semconv.ServiceVersion(cfg.ServiceVersion))
	}
	attrs = append(attrs, cfg.ResourceAttrs...)
	return resource.NewWithAttributes(semconv.SchemaURL, attrs...), nil
}
