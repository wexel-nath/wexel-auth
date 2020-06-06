package jwt

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/wexel-nath/auth"
	"wexel-auth/pkg/logger"
)

var (
	userKey contextKey = "user"
)

type contextKey string

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := Authenticate(r)
		if err != nil {
			if err != ErrExpiredToken {
				logger.Error(err)
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = fmt.Fprint(w, `{"message":"Unauthorized"}`)
			return
		}

		ctx := context.WithValue(r.Context(), userKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthorizationMiddleware(capability string) func(handler http.Handler) http.Handler {
	return func (next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, err := Authorize(r, capability)
			if err != nil {
				if err != ErrExpiredToken {
					logger.Error(err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				_, _ = fmt.Fprint(w, `{"message":"Unauthorized"}`)
				return
			}

			ctx := context.WithValue(r.Context(), userKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (auth.User, error) {
	if user, ok := ctx.Value(userKey).(auth.User); ok {
		return user, nil
	}

	return auth.User{}, errors.New("no user in context")
}
