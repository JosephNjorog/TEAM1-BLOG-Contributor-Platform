package webhooks

import "github.com/go-chi/chi/v5"

// Routes are intentionally not behind RequireAuth - Resend calls this
// endpoint directly, with no user session, authenticating itself via the
// webhook signature instead.
func Routes(handler *Handler) chi.Router {
	r := chi.NewRouter()
	r.Post("/resend", handler.Resend)
	return r
}
