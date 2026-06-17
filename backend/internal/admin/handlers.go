package admin

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/httpx"
	"team1blog/backend/internal/users"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

const timeFmt = "2006-01-02T15:04:05Z07:00"

func contributorDTO(c *ContributorSummary) map[string]any {
	var lastSubmission *string
	if c.LastSubmissionAt != nil {
		v := c.LastSubmissionAt.Format(timeFmt)
		lastSubmission = &v
	}
	return map[string]any{
		"id":                c.ID,
		"name":              c.Name,
		"email":             c.Email,
		"walletAddress":     c.WalletAddress,
		"status":            c.Status,
		"registeredAt":      c.RegisteredAt.Format(timeFmt),
		"articlesSubmitted": c.ArticlesSubmitted,
		"articlesPublished": c.ArticlesPublished,
		"totalPaidUsd":      c.TotalPaidUSD,
		"lastSubmissionAt":  lastSubmission,
	}
}

func (h *Handler) ListContributors(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.ListContributors(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	dtos := make([]map[string]any, 0, len(list))
	for _, c := range list {
		dtos = append(dtos, contributorDTO(c))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"contributors": dtos})
}

func (h *Handler) ListInvitations(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.ListPendingInvitations(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	dtos := make([]map[string]any, 0, len(list))
	for _, inv := range list {
		var usedAt *string
		if inv.UsedAt != nil {
			v := inv.UsedAt.Format(timeFmt)
			usedAt = &v
		}
		dtos = append(dtos, map[string]any{
			"id":        inv.ID,
			"email":     inv.Email,
			"role":      inv.Role,
			"expiresAt": inv.ExpiresAt.Format(timeFmt),
			"usedAt":    usedAt,
			"createdAt": inv.CreatedAt.Format(timeFmt),
		})
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"invitations": dtos})
}

func (h *Handler) Overview(w http.ResponseWriter, r *http.Request) {
	o, err := h.service.GetOverview(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, o)
}

func (h *Handler) Analytics(w http.ResponseWriter, r *http.Request) {
	m, err := h.service.GetAnalytics(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, m)
}

type setStatusRequest struct {
	Status users.Status `json:"status"`
}

func (h *Handler) SetUserStatus(w http.ResponseWriter, r *http.Request) {
	actorID, _ := auth.UserIDFromContext(r.Context())
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid user id")
		return
	}
	var req setStatusRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || (req.Status != users.StatusActive && req.Status != users.StatusInactive) {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "status must be 'active' or 'inactive'")
		return
	}
	if err := h.service.SetUserStatus(r.Context(), actorID, userID, req.Status); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

type updateRoleRequest struct {
	Role users.Role `json:"role"`
}

func (h *Handler) UpdateUserRole(w http.ResponseWriter, r *http.Request) {
	actorID, _ := auth.UserIDFromContext(r.Context())
	userID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid user id")
		return
	}
	var req updateRoleRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}
	if err := h.service.UpdateUserRole(r.Context(), actorID, userID, req.Role); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", err.Error())
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

type overrideRequest struct {
	Status articles.Status `json:"status"`
	Reason string          `json:"reason"`
}

func (h *Handler) OverrideArticleStatus(w http.ResponseWriter, r *http.Request) {
	actorID, _ := auth.UserIDFromContext(r.Context())
	articleID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	var req overrideRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}
	if err := h.service.OverrideArticleStatus(r.Context(), actorID, articleID, req.Status, req.Reason); err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrReasonRequired) {
			status = http.StatusBadRequest
		}
		httpx.Error(w, status, "override_failed", err.Error())
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}
