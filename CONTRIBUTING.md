# Contributing to jigar

Thanks for your interest! This project is small, so the rules are short.

## Development setup

```bash
git clone https://github.com/amkarkhi/jigar.git
cd jigar
go mod download
go test ./...
```

You'll need Go 1.23+.

## Workflow

1. Open an issue describing the change before sending a large PR — it's the
   fastest way to make sure the change has a home.
2. Fork, branch from `main`, and keep changes focused. One PR per change.
3. Run `go test ./...`, `go vet ./...`, and `gofmt -s -w .` before pushing.
4. Add or update tests for any behavior change.
5. Update [CHANGELOG.md](CHANGELOG.md) under the `## Unreleased` heading.

## Style

- Standard `gofmt`. No extra linters beyond `go vet` and `staticcheck`.
- Public symbols need a doc comment that starts with the symbol name.
- Keep the public API minimal — favor configuration over new functions.

## Reporting bugs

Please include:
- jigar version (`go list -m github.com/amkarkhi/jigar`)
- Go version (`go version`)
- A minimal reproduction (a few lines of code is ideal).

## Code of conduct

Be kind. Disagree about code, not people. The maintainers reserve the right
to remove comments or contributors that make the project hostile.
