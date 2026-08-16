package middleware

import (
	"car-valuation-agent/internal/infrastructure/external/keycloak"
	"context"
	"net/http"
	"strings"
)

type ContextKey string

const (
	UserKey ContextKey = "user"
)

type User struct {
	ID       string
	Username string
	Roles    []string
}

func Auth(client *keycloak.Client) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, ok := extractBearerToken(r)
			if !ok {
				http.Error(w, "Authorization Bearer token is required", http.StatusUnauthorized)
				return
			}

			claims, err := client.VerifyToken(r.Context(), token)
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			user := User{
				ID:       claims.Subject,
				Username: claims.PreferredUsername,
				Roles:    claims.Roles(),
			}

			ctx := context.WithValue(r.Context(), UserKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (User, bool) {
	user, ok := ctx.Value(UserKey).(User)
	return user, ok
}

func extractBearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "
	authorization := r.Header.Get("Authorization")
	if !strings.HasPrefix(authorization, prefix) {
		return "", false
	}
	token := strings.TrimSpace(strings.TrimPrefix(authorization, prefix))
	if token == "" {
		return "", false
	}
	return token, true
}
