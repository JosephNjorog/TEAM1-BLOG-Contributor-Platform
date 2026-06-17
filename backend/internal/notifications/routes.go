package notifications

import (
	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/auth"
)

func Routes(handler *Handler, issuer *auth.TokenIssuer) chi.Router {
	r := chi.NewRouter()
	r.Use(auth.RequireAuth(issuer))

	r.Get("/", handler.List)
	r.Post("/{id}/read", handler.MarkRead)
	r.Post("/read-all", handler.MarkAllRead)

	return r
}
