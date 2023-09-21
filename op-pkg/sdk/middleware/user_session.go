package middleware

import (
	"context"
	"net/http"
)

const (
	userSessionDataCtxKey ctxKey = "userSessionData"
	userSessionCookieName        = "user_session"
)

func ExtractUserSessionData(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		var userSessionData string
		cookie, _ := request.Cookie(userSessionCookieName)
		if cookie != nil {
			userSessionData = cookie.Value
		}
		ctx := context.WithValue(request.Context(), userSessionDataCtxKey, userSessionData)
		next.ServeHTTP(writer, request.WithContext(ctx))
	})
}

func GetUserSessionData(ctx context.Context) string {
	userSessionData, ok := ctx.Value(userSessionDataCtxKey).(string)
	if ok {
		return userSessionData
	}
	return ""
}

func SetUserSessionData(ctx context.Context, userSessionData string) context.Context {
	return context.WithValue(ctx, userSessionDataCtxKey, userSessionData)
}
