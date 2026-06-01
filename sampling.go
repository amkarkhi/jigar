package jigar

import "math/rand"

// DecideShouldTrace returns true when the X-Should-Trace header forces
// sampling on ("true") or, failing that, when a draw against the configured
// ratio succeeds. ratio<=0 disables probabilistic sampling; ratio>=1 always
// samples.
//
// Exposed so framework integrations (and tests) can reuse the same
// decision logic without duplicating it.
func DecideShouldTrace(headerValue string, ratio float64) bool {
	if headerValue == "true" {
		return true
	}
	if ratio <= 0 {
		return false
	}
	if ratio >= 1 {
		return true
	}
	return rand.Float64() < ratio
}
