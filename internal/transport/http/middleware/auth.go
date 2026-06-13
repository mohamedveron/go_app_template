package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"

	api_gen "github.com/mohamedveron/go_app_template/internal/transport/http/api_server_gen/v1"
)

type Permission int

const (
	PermissionAdmin    Permission = iota
	PermissionReadonly Permission = iota
)

type contextKey string

const permissionKey contextKey = "permission"

type AuthConfig struct {
	Workspace string
}

type AuthMiddleware struct {
	workspace string
	mu        sync.RWMutex
}

func NewAuthMiddleware(config AuthConfig) (*AuthMiddleware, error) {
	return &AuthMiddleware{
		workspace: config.Workspace,
	}, nil
}

func (am *AuthMiddleware) RequireAuth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := am.extractToken(r)
			if token == "" {
				http.Error(w, `{"error":"Authentication token required"}`, http.StatusUnauthorized)
				return
			}

			permission := am.getPermission(token)
			if permission == -1 {
				http.Error(w, `{"error":"Invalid authentication token"}`, http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), permissionKey, permission)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func (am *AuthMiddleware) RequireAdmin(r *http.Request) *api_gen.Error {
	permission, ok := r.Context().Value(permissionKey).(Permission)
	if !ok {
		return &api_gen.Error{Message: "Authentication required", Code: http.StatusUnauthorized}
	}

	if permission != PermissionAdmin {
		return &api_gen.Error{Message: "Admin permission required", Code: http.StatusForbidden}
	}
	return nil
}

func (am *AuthMiddleware) extractToken(r *http.Request) string {
	authHeader := r.Header.Get("Authorization")
	if authHeader != "" {
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) == 2 && parts[0] == "Bearer" {
			return parts[1]
		}
		return authHeader
	}
	return r.URL.Query().Get("token")
}

func (am *AuthMiddleware) getPermission(_ string) Permission {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return -1
}
