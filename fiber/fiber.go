// Package fiber provides a Fiber middleware that integrates with jigar.
//
// Typical use:
//
//	import jigarfiber "github.com/amkarkhi/jigar/fiber"
//	app.Use(jigarfiber.Middleware())
//
// See the github.com/amkarkhi/jigar/gin sub-package for the equivalent Gin
// integration.
package fiber

import (
	"github.com/amkarkhi/jigar"
	"github.com/gofiber/contrib/otelfiber"
	gofiber "github.com/gofiber/fiber/v2"
)

// Middleware returns a Fiber middleware that wires jigar into the request
// pipeline. It must be installed after jigar.Init.
func Middleware() gofiber.Handler {
	return func(c *gofiber.Ctx) error {
		t := jigar.Default()
		if t == nil {
			return c.Next()
		}
		cfg := t.Config()
		should := jigar.DecideShouldTrace(c.Get(jigar.ShouldTraceHeader), cfg.TraceRatio)
		c.SetUserContext(jigar.WithShouldTrace(c.UserContext(), should))

		return otelfiber.Middleware(
			otelfiber.WithTracerProvider(t.Provider()),
			otelfiber.WithNext(func(c *gofiber.Ctx) bool {
				return !jigar.ShouldTrace(c.UserContext())
			}),
		)(c)
	}
}
