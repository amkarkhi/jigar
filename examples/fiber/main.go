// Example: minimal Fiber server instrumented with jigar.
//
//	go run ./examples/fiber
//	curl -H 'X-Should-Trace: true' http://localhost:8080/hello/world
package main

import (
	"context"
	"log"
	"time"

	"github.com/amkarkhi/jigar"
	jigarfiber "github.com/amkarkhi/jigar/fiber"
	"github.com/gofiber/fiber/v2"
)

func main() {
	ctx := context.Background()
	t, err := jigar.Init(ctx,
		jigar.WithService("fiber-example"),
		jigar.WithEndpoint("localhost:4317"),
		jigar.WithTraceRatio(0),
	)
	if err != nil {
		log.Fatalf("init tracer: %v", err)
	}
	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = t.Shutdown(ctx)
	}()

	app := fiber.New()
	app.Use(jigarfiber.Middleware())

	app.Get("/hello/:name", func(c *fiber.Ctx) error {
		_, span := jigar.Start(c.UserContext(), "say-hello")
		defer span.End()
		name := c.Params("name")
		span.SetString("hello.name", name)
		return c.JSON(fiber.Map{"hello": name})
	})

	_ = app.Listen(":8080")
}
