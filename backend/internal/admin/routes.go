package admin

import (
	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/users"
)

func Routes(handler *Handler, issuer *auth.TokenIssuer) chi.Router {
	r := chi.NewRouter()
	r.Use(auth.RequireAuth(issuer))
	r.Use(auth.RequireRole(users.RoleSuperAdmin))

	r.Get("/overview", handler.Overview)
	r.Get("/analytics", handler.Analytics)
	r.Get("/contributors", handler.ListContributors)
	r.Get("/staff", handler.ListStaff)
	r.Get("/invitations", handler.ListInvitations)
	r.Patch("/users/{id}/status", handler.SetUserStatus)
	r.Patch("/users/{id}/role", handler.UpdateUserRole)
	r.Post("/articles/{id}/override", handler.OverrideArticleStatus)

	return r
}
