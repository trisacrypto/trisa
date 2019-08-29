package handler

import "context"

type contextKey string

var (
	contextKeyClientSide = contextKey("client-side")
)

func WithClientSide(ctx context.Context) context.Context {
	return context.WithValue(ctx, contextKeyClientSide, true)
}

func HasClientSideFromContext(ctx context.Context) bool {
	if _, ok := ctx.Value(contextKeyClientSide).(bool); ok {
		return true
	}
	return false
}
