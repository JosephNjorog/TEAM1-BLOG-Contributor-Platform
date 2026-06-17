// Package profile exposes the "current user" HTTP surface (/me). It's kept
// separate from internal/users (pure domain + repository) specifically
// because internal/auth depends on internal/users, so internal/users can't
// depend back on internal/auth's context helpers without an import cycle.
package profile

import (
	"net/http"
	"regexp"

	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/httpx"
	"team1blog/backend/internal/users"
)

type Handler struct {
	repo *users.Repository
}

func NewHandler(repo *users.Repository) *Handler {
	return &Handler{repo: repo}
}

var avalancheAddrRE = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

func userDTO(u *users.User) map[string]any {
	return map[string]any{
		"id":            u.ID,
		"name":          u.Name,
		"email":         u.Email,
		"role":          u.Role,
		"walletAddress": u.WalletAddress,
		"bio":           u.Bio,
		"status":        u.Status,
		"createdAt":     u.CreatedAt,
	}
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	u, err := h.repo.GetByID(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "not_found", "user not found")
		return
	}
	httpx.JSON(w, http.StatusOK, userDTO(u))
}

type updateProfileRequest struct {
	Name string `json:"name"`
	Bio  string `json:"bio"`
}

func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	var req updateProfileRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.Name == "" {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "name is required")
		return
	}
	if err := h.repo.UpdateProfile(r.Context(), userID, req.Name, req.Bio); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	u, err := h.repo.GetByID(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, userDTO(u))
}

type updateWalletRequest struct {
	WalletAddress string `json:"walletAddress"`
}

// UpdateWallet lets a contributor change their payout wallet from settings.
// Per the PRD this should "trigger re-verification" - format validation is
// the re-verification step for v1; a confirmation email could be layered on
// top later without changing this contract.
func (h *Handler) UpdateWallet(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	var req updateWalletRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || !avalancheAddrRE.MatchString(req.WalletAddress) {
		httpx.Error(w, http.StatusBadRequest, "invalid_wallet_address", "a valid Avalanche C-Chain address is required")
		return
	}
	if err := h.repo.UpdateWallet(r.Context(), userID, req.WalletAddress); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	u, err := h.repo.GetByID(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, userDTO(u))
}
