package auth

import (
	"time"

	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/middleware"
	"team1blog/backend/internal/users"
)

func Routes(handler *Handler, issuer *TokenIssuer) chi.Router {
	r := chi.NewRouter()

	loginLimiter := middleware.RateLimit(10, 15*time.Minute)

	r.With(loginLimiter).Post("/login", handler.Login)
	r.With(loginLimiter).Post("/refresh", handler.Refresh)
	r.Post("/register", handler.Register)
	r.Post("/logout", handler.Logout)

	r.Group(func(protected chi.Router) {
		protected.Use(RequireAuth(issuer))
		protected.With(RequireRole(users.RoleSuperAdmin)).Post("/invite", handler.Invite)
	})

	return r
}
