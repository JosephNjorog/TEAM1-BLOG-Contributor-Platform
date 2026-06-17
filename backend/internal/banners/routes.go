package banners

import (
	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/users"
)

func Routes(handler *Handler, issuer *auth.TokenIssuer) chi.Router {
	r := chi.NewRouter()
	r.Use(auth.RequireAuth(issuer))

	r.Get("/{articleId}", handler.GetLatest)

	r.Group(func(designerOnly chi.Router) {
		designerOnly.Use(auth.RequireRole(users.RoleGraphicDesigner))
		designerOnly.Post("/{articleId}/upload", handler.Upload)
		designerOnly.Post("/{articleId}/mark-ready", handler.MarkReady)
	})

	return r
}
