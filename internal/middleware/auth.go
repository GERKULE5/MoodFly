package middleware

import (
	apperror "MoodFly/pkg/error"
	auth_jwt "MoodFly/pkg/jwt"
	"context"
	"net/http"
	"strings"
)

type contextKey string

const userIDKey contextKey = "userID"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			handleError(w, apperror.Unauthorized("authorization header required"))
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)

		if len(parts) != 2 || parts[0] != "Bearer" {
			handleError(w, apperror.Unauthorized("invalid authorization header format"))
			return
		}

		claims, err := auth_jwt.ParseAccessToken(parts[1])
		if err != nil {
			handleError(w, apperror.Unauthorized("invalid token"))
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, claims.UserID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserIDFromContext(ctx context.Context) (int, bool) {
	userID, ok := ctx.Value(userIDKey).(int)
	return userID, ok
}

func handleError(w http.ResponseWriter, err error) {
	code, msg := apperror.ToHTTP(err)
	http.Error(w, msg, code)
}
