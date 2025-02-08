package client_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	. "github.com/trisacrypto/trisa/pkg/openvasp/client"
)

func TestContextWithRequestID(t *testing.T) {
	parent := context.Background()
	requestID := "12345"
	ctx := ContextWithRequestID(parent, requestID)

	val, ok := RequestIDFromContext(ctx)
	require.True(t, ok, "expected to find requestID in context")
	require.Equal(t, requestID, val, "expected requestID %s, got %s", requestID, val)
}

func TestContextWithTracing(t *testing.T) {
	parent := context.Background()
	tracingID := "trace-12345"
	ctx := ContextWithTracing(parent, tracingID)

	val, ok := TracingFromContext(ctx)
	require.True(t, ok, "expected to find tracingID in context")
	require.Equal(t, tracingID, val, "expected tracingID %s, got %s", tracingID, val)
}
