package middleware

import (
	"context"
	"net/http"
)

const (
	requestContextDataCtxKey ctxKey = "requestContextData"
	requestContextHeaderName        = "X-REQUEST-CONTEXT"
)

func ExtractRequestContextData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		requestContextData := request.Header.Get(requestContextHeaderName)
		ctx := context.WithValue(request.Context(), requestContextDataCtxKey, requestContextData)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetRequestContextData(ctx context.Context) string {
	requestContextData, ok := ctx.Value(requestContextDataCtxKey).(string)
	if ok {
		return requestContextData
	}
	return ""
}

func SetRequestContextData(ctx context.Context, requestContextData string) context.Context {
	return context.WithValue(ctx, requestContextDataCtxKey, requestContextData)
}
