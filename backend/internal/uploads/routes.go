package uploads

import (
	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/users"
)

func Routes(handler *Handler, issuer *auth.TokenIssuer) chi.Router {
	r := chi.NewRouter()
	r.Use(auth.RequireAuth(issuer))
	r.Use(auth.RequireRole(users.RoleContributor))

	r.Post("/inline/{articleId}", handler.UploadInline)

	return r
}
