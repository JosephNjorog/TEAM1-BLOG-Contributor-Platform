package reviews

import (
	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/users"
)

func Routes(handler *Handler, issuer *auth.TokenIssuer) chi.Router {
	r := chi.NewRouter()
	r.Use(auth.RequireAuth(issuer))

	r.Get("/article/{articleId}", handler.ListForArticle)

	r.Group(func(moderatorOnly chi.Router) {
		moderatorOnly.Use(auth.RequireRole(users.RoleModerator))
		moderatorOnly.Post("/", handler.Submit)
		moderatorOnly.Get("/activity", handler.Activity)
	})

	r.Group(func(contributorOnly chi.Router) {
		contributorOnly.Use(auth.RequireRole(users.RoleContributor))
		contributorOnly.Post("/suggestions/{id}/accept", handler.suggestionStatusHandler(SuggestionAccepted))
		contributorOnly.Post("/suggestions/{id}/reject", handler.suggestionStatusHandler(SuggestionRejected))
	})

	return r
}
