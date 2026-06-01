package jigar

import "errors"

var (
	// ErrEndpointRequired is returned by Init/New when WithEndpoint was not set.
	ErrEndpointRequired = errors.New("jigar: endpoint is required")
	// ErrServiceRequired is returned by Init/New when WithService was not set.
	ErrServiceRequired = errors.New("jigar: service is required")
	// ErrNotInitialized is returned by package-level helpers when called
	// before Init.
	ErrNotInitialized = errors.New("jigar: tracer not initialized; call jigar.Init first")
)
