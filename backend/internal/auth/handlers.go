package auth

import (
	"errors"
	"net/http"

	"team1blog/backend/internal/httpx"
	"team1blog/backend/internal/users"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type inviteRequest struct {
	Email string     `json:"email"`
	Role  users.Role `json:"role"`
}

func (h *Handler) Invite(w http.ResponseWriter, r *http.Request) {
	var req inviteRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}
	actorID, _ := UserIDFromContext(r.Context())

	if err := h.service.Invite(r.Context(), actorID, req.Email, req.Role); err != nil {
		httpx.Error(w, http.StatusBadRequest, "invite_failed", err.Error())
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]string{"status": "invited"})
}

type registerRequest struct {
	Token         string `json:"token"`
	Name          string `json:"name"`
	Password      string `json:"password"`
	Bio           string `json:"bio"`
	WalletAddress string `json:"walletAddress"`
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req registerRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}
	if req.Token == "" || req.Name == "" || len(req.Password) < 8 {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "name, a token, and an 8+ character password are required")
		return
	}

	u, tokens, err := h.service.RegisterFromInvite(r.Context(), RegisterInput{
		Token:         req.Token,
		Name:          req.Name,
		Password:      req.Password,
		Bio:           req.Bio,
		WalletAddress: req.WalletAddress,
	})
	if err != nil {
		status, code := mapAuthError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}

	httpx.JSON(w, http.StatusCreated, authResponse(u, tokens))
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}

	u, tokens, err := h.service.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		status, code := mapAuthError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}

	httpx.JSON(w, http.StatusOK, authResponse(u, tokens))
}

type refreshRequest struct {
	RefreshToken string `json:"refreshToken"`
}

func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.RefreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "refreshToken is required")
		return
	}

	tokens, err := h.service.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		status, code := mapAuthError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}

	httpx.JSON(w, http.StatusOK, map[string]any{
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"expiresAt":    tokens.ExpiresAt,
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req refreshRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.RefreshToken == "" {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "refreshToken is required")
		return
	}
	_ = h.service.Logout(r.Context(), req.RefreshToken)
	httpx.JSON(w, http.StatusNoContent, nil)
}

func authResponse(u *users.User, tokens *TokenPair) map[string]any {
	return map[string]any{
		"user": map[string]any{
			"id":            u.ID,
			"name":          u.Name,
			"email":         u.Email,
			"role":          u.Role,
			"walletAddress": u.WalletAddress,
			"status":        u.Status,
		},
		"accessToken":  tokens.AccessToken,
		"refreshToken": tokens.RefreshToken,
		"expiresAt":    tokens.ExpiresAt,
	}
}

func mapAuthError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrInvalidCreds):
		return http.StatusUnauthorized, "invalid_credentials"
	case errors.Is(err, ErrAccountInactive):
		return http.StatusForbidden, "account_inactive"
	case errors.Is(err, ErrEmailInUse):
		return http.StatusConflict, "email_in_use"
	case errors.Is(err, ErrInvitationNotFound):
		return http.StatusNotFound, "invitation_not_found"
	case errors.Is(err, ErrInvitationExpired):
		return http.StatusGone, "invitation_expired"
	case errors.Is(err, ErrInvitationUsed):
		return http.StatusConflict, "invitation_used"
	case errors.Is(err, ErrInvalidWalletAddr):
		return http.StatusBadRequest, "invalid_wallet_address"
	case errors.Is(err, ErrRefreshNotFound), errors.Is(err, ErrRefreshRevoked):
		return http.StatusUnauthorized, "invalid_refresh_token"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
