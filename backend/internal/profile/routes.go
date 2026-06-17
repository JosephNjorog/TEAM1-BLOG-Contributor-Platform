package profile

import (
	"github.com/go-chi/chi/v5"
	"team1blog/backend/internal/auth"
)

func Routes(handler *Handler, issuer *auth.TokenIssuer) chi.Router {
	r := chi.NewRouter()
	r.Use(auth.RequireAuth(issuer))

	r.Get("/", handler.Me)
	r.Patch("/", handler.UpdateProfile)
	r.Patch("/wallet", handler.UpdateWallet)

	return r
}
