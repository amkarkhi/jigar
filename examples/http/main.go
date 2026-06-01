// Example: minimal net/http server instrumented with jigar.
//
//	go run ./examples/http
//	curl -H 'X-Should-Trace: true' http://localhost:8080/hello
package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/amkarkhi/jigar"
	jigarhttp "github.com/amkarkhi/jigar/nethttp"
)

func main() {
	ctx := context.Background()
	t, err := jigar.Init(ctx,
		jigar.WithService("http-example"),
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

	mux := http.NewServeMux()
	mux.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		_, span := jigar.Start(r.Context(), "say-hello")
		defer span.End()
		span.SetString("hello.who", "world")
		_, _ = w.Write([]byte("hello, world\n"))
	})

	log.Fatal(http.ListenAndServe(":8080", jigarhttp.Middleware(mux, "http-example")))
}
