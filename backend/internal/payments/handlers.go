package payments

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type paymentDTO struct {
	ID              uuid.UUID `json:"id"`
	ArticleID       uuid.UUID `json:"articleId"`
	ArticleTitle    string    `json:"articleTitle"`
	ContributorID   uuid.UUID `json:"contributorId"`
	ContributorName string    `json:"contributorName"`
	WalletAddress   string    `json:"walletAddress"`
	AmountUSD       float64   `json:"amountUsd"`
	TxHash          *string   `json:"txHash"`
	Status          string    `json:"status"`
	InitiatedAt     *string   `json:"initiatedAt"`
	ConfirmedAt     *string   `json:"confirmedAt"`
	CreatedAt       string    `json:"createdAt"`
}

const timeFmt = "2006-01-02T15:04:05Z07:00"

func toDTO(p *Payment) paymentDTO {
	dto := paymentDTO{
		ID:              p.ID,
		ArticleID:       p.ArticleID,
		ArticleTitle:    p.ArticleTitle,
		ContributorID:   p.ContributorID,
		ContributorName: p.ContributorName,
		WalletAddress:   p.WalletAddress,
		AmountUSD:       p.AmountUSD,
		TxHash:          p.TxHash,
		Status:          string(p.Status),
		CreatedAt:       p.CreatedAt.Format(timeFmt),
	}
	if p.InitiatedAt != nil {
		v := p.InitiatedAt.Format(timeFmt)
		dto.InitiatedAt = &v
	}
	if p.ConfirmedAt != nil {
		v := p.ConfirmedAt.Format(timeFmt)
		dto.ConfirmedAt = &v
	}
	return dto
}

func (h *Handler) ListLedger(w http.ResponseWriter, r *http.Request) {
	list, err := h.service.ListLedger(r.Context())
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	dtos := make([]paymentDTO, 0, len(list))
	for _, p := range list {
		dtos = append(dtos, toDTO(p))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"payments": dtos})
}

func (h *Handler) Release(w http.ResponseWriter, r *http.Request) {
	adminID, _ := auth.UserIDFromContext(r.Context())
	articleID, err := uuid.Parse(chi.URLParam(r, "articleId"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	payment, err := h.service.Release(r.Context(), articleID, adminID)
	if err != nil {
		status, code := mapError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, toDTO(payment))
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, articles.ErrNotFound):
		return http.StatusNotFound, "article_not_found"
	case errors.Is(err, ErrNotPublished):
		return http.StatusConflict, "article_not_published"
	case errors.Is(err, ErrAlreadyReleased):
		return http.StatusConflict, "payment_already_released"
	case errors.Is(err, ErrNoWallet):
		return http.StatusConflict, "contributor_no_wallet"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
