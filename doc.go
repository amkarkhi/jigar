// Package jigar is a small, batteries-included OpenTelemetry tracing helper
// for Go services.
//
// It wraps OTLP exporter setup, request-level sampling via a header, and
// ergonomic span helpers so application code can add tracing in a few lines
// without learning the full OpenTelemetry API.
//
// Quick start:
//
//	t, err := jigar.Init(ctx,
//	    jigar.WithService("my-service"),
//	    jigar.WithEndpoint("localhost:4317"),
//	    jigar.WithTraceRatio(0.01),
//	)
//	if err != nil { return err }
//	defer t.Shutdown(context.Background())
//
//	ctx, span := jigar.Start(ctx, "do-work")
//	defer span.End()
//	span.SetString("user.id", userID)
//
// Per-request sampling is controlled by the X-Should-Trace header (forced on
// when "true") or the configured TraceRatio (probabilistic). Set TraceRatio
// to 1 to sample every request, or 0 to disable probabilistic sampling and
// only trace requests that carry the header.
package jigar
