package substack

import (
	"net/http"

	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Sync(w http.ResponseWriter, r *http.Request) {
	count, err := h.service.Sync(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "sync_failed", err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"synced": count})
}

func (h *Handler) ListMine(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	posts, err := h.service.ListForContributor(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"posts": posts})
}
