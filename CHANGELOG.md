# Changelog

All notable changes to this project are documented here. The format is based
on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/) and the project
follows [Semantic Versioning](https://semver.org/).

## [Unreleased]

## [0.1.0] - 2026-05-30

First public, open-source release. The API was reshaped from the original
Digikala-internal version; if you used the pre-OSS package, treat this as a
breaking re-release.

### Added
- Functional options (`WithService`, `WithEndpoint`, `WithTraceRatio`, ...)
  on `jigar.Init` / `jigar.New`.
- `*Tracer` type with `Shutdown`, `Provider`, and `Config` accessors.
- `*Span` handle with typed setters: `SetString`, `SetInt`, `SetFloat`,
  `SetBool`, `SetJSON`, and `SetAttributes` for raw KVs.
- OTLP/HTTP exporter support alongside OTLP/gRPC (`WithProtocol`).
- net/http middleware in `github.com/amkarkhi/jigar/nethttp`.
- Typed context key (`jigar.WithShouldTrace`, `jigar.ShouldTrace`).
- Examples for Gin, Fiber, and net/http under `examples/`.

### Changed
- Module path is now `github.com/amkarkhi/jigar`.
- Gin and Fiber middlewares moved to sub-packages
  (`github.com/amkarkhi/jigar/gin`, `.../fiber`).
- Error variables renamed to `ErrEndpointRequired`, `ErrServiceRequired`,
  `ErrNotInitialized`.
- `StartSpan` is now `jigar.Start`; the per-span handle was renamed from
  `CustomTracer` to `Span`.

### Removed
- The false-alarm error-skip list and the internal `demeter` dependency.
  Filter errors in your own code before calling `span.AddError`.
- The library no longer panics when `Start` is called before `Init`; it
  returns a no-op span instead.
- The zerolog dependency.

[Unreleased]: https://github.com/amkarkhi/jigar/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/amkarkhi/jigar/releases/tag/v0.1.0
