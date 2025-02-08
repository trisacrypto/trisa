package client

import "context"

// API-specific context keys for passing values to requests via the context. These keys
// are unexported to reduce the size of the public interface an prevent incorrect handling.
type contextKey uint8

// Allocate context keys to simplify context key usage in helper functions.
const (
	contextKeyUnknown contextKey = iota
	contextKeyRequestID
	contextKeyTracing
)

// Adds a request identifer to the context which is sent with the request
// in the request-identifier header. (OpenVASP specific - should specify the request
// is about a specific transaction).
func ContextWithRequestID(parent context.Context, requestID string) context.Context {
	return context.WithValue(parent, contextKeyRequestID, requestID)
}

// Extracts a request identifier from the context.
func RequestIDFromContext(ctx context.Context) (string, bool) {
	requestID, ok := ctx.Value(contextKeyRequestID).(string)
	return requestID, ok
}

// Adds a tracing ID to the context which is sent with the request in the
// X-Request-ID header. This is an per-request unique identifier that is used
// for tracing and observability, not for the protocol.
func ContextWithTracing(parent context.Context, tracingID string) context.Context {
	return context.WithValue(parent, contextKeyTracing, tracingID)
}

// Extracts a tracing ID from the context.
func TracingFromContext(ctx context.Context) (string, bool) {
	tracingID, ok := ctx.Value(contextKeyTracing).(string)
	return tracingID, ok
}
