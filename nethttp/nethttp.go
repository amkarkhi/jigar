// Package nethttp provides a net/http middleware that integrates with jigar.
//
//	mux := http.NewServeMux()
//	handler := jigarhttp.Middleware(mux)
//	http.ListenAndServe(":8080", handler)
package nethttp

import (
	"net/http"

	"github.com/amkarkhi/jigar"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// Middleware wraps next so that incoming requests are evaluated against the
// jigar sampling rules and, when sampled, instrumented with otelhttp.
// Pass the operation name otelhttp should report (typically the route or
// service name).
func Middleware(next http.Handler, operation string) http.Handler {
	instrumented := otelhttp.NewHandler(next, operation)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t := jigar.Default()
		if t == nil {
			next.ServeHTTP(w, r)
			return
		}
		cfg := t.Config()
		should := jigar.DecideShouldTrace(r.Header.Get(jigar.ShouldTraceHeader), cfg.TraceRatio)
		r = r.WithContext(jigar.WithShouldTrace(r.Context(), should))
		if !should {
			next.ServeHTTP(w, r)
			return
		}
		instrumented.ServeHTTP(w, r)
	})
}
