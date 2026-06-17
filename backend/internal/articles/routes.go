package articles

import (
	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/users"
)

func Routes(handler *Handler, issuer *auth.TokenIssuer) chi.Router {
	r := chi.NewRouter()
	r.Use(auth.RequireAuth(issuer))

	// Visible to any authenticated role; results are scoped per-user inside
	// the service/repository per PRD section 11.
	r.Get("/", handler.List)
	r.Get("/{id}", handler.Get)

	r.Group(func(contributorOnly chi.Router) {
		contributorOnly.Use(auth.RequireRole(users.RoleContributor))
		contributorOnly.Post("/", handler.Create)
		contributorOnly.Put("/{id}", handler.Update)
		contributorOnly.Post("/{id}/submit", handler.Submit)
		contributorOnly.Delete("/{id}", handler.Delete)
	})

	r.Group(func(publisherOnly chi.Router) {
		publisherOnly.Use(auth.RequireRole(users.RolePublisher))
		publisherOnly.Post("/{id}/publish", handler.Publish)
	})

	return r
}
