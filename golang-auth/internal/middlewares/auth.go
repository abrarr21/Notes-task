package middlewares

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/abrarr21/test/internal/utils"
)

type contextKey string

const UserIDKey contextKey = "userID"

func RequireAuth(secret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := ""

			cookie, err := r.Cookie("accessToken")
			if err == nil {
				tokenString = cookie.Value
			} else {
				_, token, found := strings.Cut(r.Header.Get("Authorization"), "Bearer ")
				if found {
					tokenString = token
				}
			}

			if tokenString == "" {
				utils.JSON(w, http.StatusUnauthorized, "missing token", nil)
				return
			}

			claims, err := utils.ParseToken(tokenString, secret)
			if err != nil {
				if errors.Is(err, utils.ErrTokenExpired) {
					utils.JSON(w, http.StatusUnauthorized, "token has expired", nil)
					return
				}
				utils.JSON(w, http.StatusUnauthorized, "invalid token", nil)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(r *http.Request) (string, bool) {
	v, ok := r.Context().Value(UserIDKey).(string)
	return v, ok
}
