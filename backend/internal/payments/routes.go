package payments

import (
	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/users"
)

func Routes(handler *Handler, issuer *auth.TokenIssuer) chi.Router {
	r := chi.NewRouter()
	r.Use(auth.RequireAuth(issuer))

	r.Get("/mine", handler.ListMine)

	r.Group(func(adminOnly chi.Router) {
		adminOnly.Use(auth.RequireRole(users.RoleSuperAdmin))
		adminOnly.Get("/", handler.ListLedger)
		adminOnly.Post("/{articleId}/release", handler.Release)
	})

	return r
}
