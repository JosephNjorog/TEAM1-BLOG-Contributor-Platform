package auth

import (
	"net/http"
	"strings"

	"team1blog/backend/internal/httpx"
	"team1blog/backend/internal/users"
)

// RequireAuth validates the JWT bearer token and stashes the parsed claims
// in the request context for downstream handlers / RBAC checks.
func RequireAuth(issuer *TokenIssuer) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractToken(r)
			if tokenString == "" {
				httpx.Error(w, http.StatusUnauthorized, "unauthorized", "missing or invalid token")
				return
			}
			claims, err := issuer.Parse(tokenString)
			if err != nil {
				httpx.Error(w, http.StatusUnauthorized, "unauthorized", "missing or invalid token")
				return
			}
			next.ServeHTTP(w, r.WithContext(SetClaims(r.Context(), claims)))
		})
	}
}

func extractToken(r *http.Request) string {
	if h := r.Header.Get("Authorization"); strings.HasPrefix(h, "Bearer ") {
		return strings.TrimPrefix(h, "Bearer ")
	}
	// WebSocket connections can't set custom headers from the browser API, so
	// allow the token as a query param on that one route.
	return r.URL.Query().Get("token")
}

func RequireRole(roles ...users.Role) func(http.Handler) http.Handler {
	allowed := make(map[users.Role]bool, len(roles))
	for _, r := range roles {
		allowed[r] = true
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := ClaimsFromContext(r.Context())
			if !ok || !allowed[claims.Role] {
				httpx.Error(w, http.StatusForbidden, "forbidden", "you do not have access to this resource")
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
