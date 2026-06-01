// Package gin provides a Gin middleware that integrates with jigar.
//
// Typical use:
//
//	import jigargin "github.com/amkarkhi/jigar/gin"
//	router.Use(jigargin.Middleware())
//
// The middleware decides per request whether to record a trace (using the
// X-Should-Trace header or the configured ratio), stores the decision on
// the request context, and only then delegates to otelgin so that no spans
// are produced for untraced requests.
package gin

import (
	"net/http"

	"github.com/amkarkhi/jigar"
	gingonic "github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

// ContextKey is the gin.Context key under which the should-trace decision
// is stored. Handlers can read it with c.GetBool(ContextKey).
const ContextKey = "shouldTrace"

// Middleware returns a Gin middleware that wires jigar into the request
// pipeline. It must be installed after jigar.Init.
func Middleware() gingonic.HandlerFunc {
	return func(c *gingonic.Context) {
		t := jigar.Default()
		if t == nil {
			c.Next()
			return
		}
		cfg := t.Config()
		should := jigar.DecideShouldTrace(c.GetHeader(jigar.ShouldTraceHeader), cfg.TraceRatio)

		ctx := jigar.WithShouldTrace(c.Request.Context(), should)
		c.Set(ContextKey, should)
		c.Request = c.Request.WithContext(ctx)

		otelgin.Middleware(cfg.Service,
			otelgin.WithTracerProvider(t.Provider()),
			otelgin.WithFilter(func(r *http.Request) bool {
				return jigar.ShouldTrace(r.Context())
			}),
			otelgin.WithSpanNameFormatter(func(r *http.Request) string {
				return r.URL.Path
			}),
		)(c)
		c.Next()
	}
}
