package jigar

import "context"

// ShouldTraceHeader is the HTTP header that forces a request to be traced
// when set to "true", regardless of the configured sampling ratio.
const ShouldTraceHeader = "X-Should-Trace"

// ctxKey is a private, typed context key so callers can't collide with us.
type ctxKey struct{}

var shouldTraceKey = ctxKey{}

// WithShouldTrace returns a context that carries the should-trace decision.
// Middleware sets this; Start consults it.
func WithShouldTrace(ctx context.Context, on bool) context.Context {
	return context.WithValue(ctx, shouldTraceKey, on)
}

// ShouldTrace reports whether the context was marked for tracing by a
// jigar middleware. Returns false when the value is absent.
func ShouldTrace(ctx context.Context) bool {
	v, _ := ctx.Value(shouldTraceKey).(bool)
	return v
}
