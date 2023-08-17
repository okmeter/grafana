package middleware

import (
	"context"
)

const (
	querierContextKey ctxKey = "querierContext"
)

func NewQuerierContext(ctx context.Context, querier string) context.Context {
	return context.WithValue(ctx, querierContextKey, querier)
}

func GetQuerier(ctx context.Context) string {
	querier, ok := ctx.Value(querierContextKey).(string)
	if ok {
		return querier
	}
	return ""
}
