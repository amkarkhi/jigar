package jigar

import (
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// Protocol selects the OTLP transport used by the default exporter.
type Protocol string

const (
	// ProtocolGRPC sends spans over OTLP/gRPC. Default.
	ProtocolGRPC Protocol = "grpc"
	// ProtocolHTTP sends spans over OTLP/HTTP (protobuf).
	ProtocolHTTP Protocol = "http"
)

// Config holds tracer configuration. Build with Option functions rather than
// constructing this directly; the zero value is not usable.
type Config struct {
	Service        string
	ServiceVersion string
	Endpoint       string
	Protocol       Protocol
	Insecure       bool
	TraceRatio     float64
	TracerName     string
	ResourceAttrs  []attribute.KeyValue
	// Exporter, if set, is used instead of building one from Endpoint/Protocol.
	Exporter sdktrace.SpanExporter
}

// Option configures a Tracer.
type Option func(*Config)

// WithService sets the service name reported in spans. Required.
func WithService(name string) Option {
	return func(c *Config) { c.Service = name }
}

// WithServiceVersion sets the service version reported in spans.
func WithServiceVersion(v string) Option {
	return func(c *Config) { c.ServiceVersion = v }
}

// WithEndpoint sets the OTLP collector endpoint, e.g. "localhost:4317".
// Required unless WithExporter is used.
func WithEndpoint(endpoint string) Option {
	return func(c *Config) { c.Endpoint = endpoint }
}

// WithProtocol selects gRPC (default) or HTTP transport for the default exporter.
func WithProtocol(p Protocol) Option {
	return func(c *Config) { c.Protocol = p }
}

// WithInsecure controls whether the OTLP connection uses TLS.
// Default is true, set false for collectors that require TLS.
func WithInsecure(insecure bool) Option {
	return func(c *Config) { c.Insecure = insecure }
}

// WithTraceRatio sets the probabilistic sampling rate for requests that do
// not carry the X-Should-Trace header. 0 disables probabilistic sampling
// (only header-forced requests are traced); 1 samples every request.
func WithTraceRatio(ratio float64) Option {
	return func(c *Config) { c.TraceRatio = ratio }
}

// WithTracerName overrides the OpenTelemetry tracer name. Defaults to the
// service name.
func WithTracerName(name string) Option {
	return func(c *Config) { c.TracerName = name }
}

// WithResourceAttributes adds extra OpenTelemetry resource attributes.
func WithResourceAttributes(attrs ...attribute.KeyValue) Option {
	return func(c *Config) { c.ResourceAttrs = append(c.ResourceAttrs, attrs...) }
}

// WithExporter installs a custom span exporter, bypassing the built-in OTLP
// exporter. Useful for tests or non-OTLP backends.
func WithExporter(e sdktrace.SpanExporter) Option {
	return func(c *Config) { c.Exporter = e }
}
