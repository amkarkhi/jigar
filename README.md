# jigar

[![Go Reference](https://pkg.go.dev/badge/github.com/amkarkhi/jigar.svg)](https://pkg.go.dev/github.com/amkarkhi/jigar)
[![CI](https://github.com/amkarkhi/jigar/actions/workflows/ci.yml/badge.svg)](https://github.com/amkarkhi/jigar/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/amkarkhi/jigar)](https://goreportcard.com/report/github.com/amkarkhi/jigar)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A small, batteries-included OpenTelemetry tracing helper for Go services.

`jigar` wraps OTLP exporter setup, per-request sampling, and span ergonomics
so application code can add tracing in a few lines without learning the full
OpenTelemetry API. It ships with middlewares for **Gin**, **Fiber**, and
**net/http**.

> "jigar" (جیگر) is Persian for *liver* — also a term of endearment for
> someone or something you adore. This package is the friendly little
> instrument that keeps your traces healthy.

---

## Features

- One-call setup for an OTLP/gRPC or OTLP/HTTP exporter.
- Per-request opt-in tracing via the `X-Should-Trace` header.
- Probabilistic sampling via `TraceRatio` for the rest of your traffic.
- Drop-in middlewares for Gin, Fiber, and `net/http`.
- A `*Span` handle with typed setters (`SetString`, `SetInt`, `SetJSON`, ...)
  that are safe no-ops when a request isn't sampled.
- No internal globals leaking out, no `log.Fatal` in library code, no panics.

## Install

```bash
go get github.com/amkarkhi/jigar
```

For framework support, also install the matching sub-package:

```bash
go get github.com/amkarkhi/jigar/gin       # Gin
go get github.com/amkarkhi/jigar/fiber     # Fiber
go get github.com/amkarkhi/jigar/nethttp   # net/http
```

## Try it locally with Jaeger

The repo ships a `docker-compose.yml` that runs the Jaeger all-in-one image
with the OTLP receivers enabled:

```bash
docker compose up -d
# then run one of the examples, e.g.
go run ./examples/http
# in another terminal
curl -H 'X-Should-Trace: true' http://localhost:8080/hello
```

Open <http://localhost:16686>, pick your service from the dropdown, and the
trace will be there. Stop it with `docker compose down`.

## Quick start

```go
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
        jigar.WithService("my-service"),
        jigar.WithEndpoint("localhost:4317"),
        jigar.WithTraceRatio(0.01), // sample 1% of requests
    )
    if err != nil {
        log.Fatal(err)
    }
    defer func() {
        ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
        defer cancel()
        _ = t.Shutdown(ctx)
    }()

    mux := http.NewServeMux()
    mux.HandleFunc("/work", func(w http.ResponseWriter, r *http.Request) {
        _, span := jigar.Start(r.Context(), "do-work")
        defer span.End()
        span.SetString("user.id", r.Header.Get("X-User-Id"))
        // ... do work ...
        if err := doWork(); err != nil {
            span.EndWithError(err)
            http.Error(w, err.Error(), 500)
            return
        }
        _, _ = w.Write([]byte("ok"))
    })

    log.Fatal(http.ListenAndServe(":8080", jigarhttp.Middleware(mux, "my-service")))
}
```

## Per-request sampling

Every middleware checks two things:

1. If the incoming request has `X-Should-Trace: true`, the request is **always**
   traced — useful for forcing a trace on a specific call from production.
2. Otherwise, the request is sampled with probability `TraceRatio`.

Set `TraceRatio` to `0` to disable probabilistic sampling and only trace
requests that carry the header. Set it to `1` to trace every request.

Inside handlers, `jigar.Start` is a no-op for unsampled requests, so calls
like `span.SetJSON(...)` cost nothing when tracing is off.

## Framework middlewares

### Gin

```go
import jigargin "github.com/amkarkhi/jigar/gin"

r := gin.New()
r.Use(jigargin.Middleware())
```

### Fiber

```go
import jigarfiber "github.com/amkarkhi/jigar/fiber"

app := fiber.New()
app.Use(jigarfiber.Middleware())
```

### net/http

```go
import jigarhttp "github.com/amkarkhi/jigar/nethttp"

handler := jigarhttp.Middleware(mux, "service-name")
http.ListenAndServe(":8080", handler)
```

Full runnable programs are under [`examples/`](examples/).

## Span API

```go
ctx, span := jigar.Start(ctx, "fetch-user")
defer span.End()

span.SetString("user.id", id)
span.SetInt("rows.scanned", n)
span.SetBool("cache.hit", false)
span.SetFloat("duration.ms", 12.3)
span.SetJSON("query", queryStruct)         // marshals to JSON

if err := db.QueryRow(...).Scan(...); err != nil {
    span.EndWithError(err)                 // records + ends
    return err
}
```

All setters are safe to call on an unsampled span — they simply do nothing.

## Configuration

| Option                       | Default     | Notes                                         |
|------------------------------|-------------|-----------------------------------------------|
| `WithService(name)`          | _required_  | Service name reported on every span.          |
| `WithEndpoint(addr)`         | _required_  | OTLP collector address, e.g. `localhost:4317`.|
| `WithServiceVersion(v)`      | unset       | Optional semantic version of the service.     |
| `WithProtocol(p)`            | `ProtocolGRPC` | `ProtocolGRPC` or `ProtocolHTTP`.          |
| `WithInsecure(b)`            | `true`      | Disable TLS to the collector.                 |
| `WithTraceRatio(r)`          | `0`         | Probabilistic sampling rate (0–1).            |
| `WithTracerName(name)`       | service     | OpenTelemetry tracer name.                    |
| `WithResourceAttributes(kv...)` | none     | Extra OTel resource attributes.               |
| `WithExporter(e)`            | OTLP        | Override the exporter (tests, custom sinks).  |

## Versioning

`jigar` follows [Semantic Versioning](https://semver.org/). The API is
unstable until `v1.0.0`; breaking changes will be called out in
[CHANGELOG.md](CHANGELOG.md).

## Contributing

Contributions are welcome — see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

[MIT](LICENSE)
