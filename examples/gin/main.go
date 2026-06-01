// Example: minimal Gin server instrumented with jigar.
//
//	go run ./examples/gin
//	curl -H 'X-Should-Trace: true' http://localhost:8080/hello/world
package main

import (
	"context"
	"log"
	"time"

	"github.com/amkarkhi/jigar"
	jigargin "github.com/amkarkhi/jigar/gin"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()
	t, err := jigar.Init(ctx,
		jigar.WithService("gin-example"),
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

	r := gin.Default()
	r.Use(jigargin.Middleware())

	r.GET("/hello/:name", func(c *gin.Context) {
		_, span := jigar.Start(c.Request.Context(), "say-hello")
		defer span.End()
		name := c.Param("name")
		span.SetString("hello.name", name)
		c.JSON(200, gin.H{"hello": name})
	})

	_ = r.Run(":8080")
}
