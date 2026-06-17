package reviews

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
	service     *Service
	articlesSvc *articles.Service
}

func NewHandler(service *Service, articlesSvc *articles.Service) *Handler {
	return &Handler{service: service, articlesSvc: articlesSvc}
}

type suggestionDTO struct {
	ID             uuid.UUID `json:"id"`
	ReviewCycleID  uuid.UUID `json:"reviewCycleId"`
	ReviewerID     uuid.UUID `json:"reviewerId"`
	ReviewerName   string    `json:"reviewerName"`
	RangeStart     int       `json:"rangeStart"`
	RangeEnd       int       `json:"rangeEnd"`
	SuggestionText string    `json:"suggestionText"`
	Status         string    `json:"status"`
	CreatedAt      string    `json:"createdAt"`
}

func toSuggestionDTO(s *Suggestion) suggestionDTO {
	return suggestionDTO{
		ID:             s.ID,
		ReviewCycleID:  s.ReviewCycleID,
		ReviewerID:     s.ReviewerID,
		ReviewerName:   s.ReviewerName,
		RangeStart:     s.RangeStart,
		RangeEnd:       s.RangeEnd,
		SuggestionText: s.SuggestionText,
		Status:         string(s.Status),
		CreatedAt:      s.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type cycleDTO struct {
	ID           uuid.UUID `json:"id"`
	ArticleTitle string    `json:"articleTitle,omitempty"`
	ReviewerName string    `json:"reviewerName"`
	Decision     string    `json:"decision"`
	Summary      string    `json:"summary"`
	CreatedAt    string    `json:"createdAt"`
}

func toCycleDTO(c *ReviewCycle) cycleDTO {
	return cycleDTO{
		ID:           c.ID,
		ArticleTitle: c.ArticleTitle,
		ReviewerName: c.ReviewerName,
		Decision:     string(c.Decision),
		Summary:      c.Summary,
		CreatedAt:    c.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

type submitReviewRequest struct {
	ArticleID   uuid.UUID         `json:"articleId"`
	Decision    Decision          `json:"decision"`
	Summary     string            `json:"summary"`
	Suggestions []SuggestionInput `json:"suggestions"`
}

func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := auth.UserIDFromContext(r.Context())
	var req submitReviewRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}
	if req.Decision != DecisionApproved && req.Decision != DecisionChangesRequested {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "decision must be 'approved' or 'changes_requested'")
		return
	}

	a, err := h.service.SubmitReview(r.Context(), SubmitReviewInput{
		ArticleID:   req.ArticleID,
		ReviewerID:  reviewerID,
		Decision:    req.Decision,
		Summary:     req.Summary,
		Suggestions: req.Suggestions,
	})
	if err != nil {
		status, code := mapError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, articles.ToDTO(a))
}

func (h *Handler) ListForArticle(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	claims, _ := auth.ClaimsFromContext(r.Context())
	articleID, err := uuid.Parse(chi.URLParam(r, "articleId"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}

	// Reuse the article visibility rules as the authorization gate for its
	// review history - whoever can see the article can see its feedback.
	if _, err := h.articlesSvc.GetVisible(r.Context(), articleID, userID, claims.Role); err != nil {
		status, code := mapError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}

	cycles, suggestions, err := h.service.ListForArticle(r.Context(), articleID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	cycleDTOs := make([]cycleDTO, 0, len(cycles))
	for _, c := range cycles {
		cycleDTOs = append(cycleDTOs, toCycleDTO(c))
	}
	suggestionDTOs := make([]suggestionDTO, 0, len(suggestions))
	for _, s := range suggestions {
		suggestionDTOs = append(suggestionDTOs, toSuggestionDTO(s))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"reviewCycles": cycleDTOs, "suggestions": suggestionDTOs})
}

func (h *Handler) Activity(w http.ResponseWriter, r *http.Request) {
	reviewerID, _ := auth.UserIDFromContext(r.Context())
	cycles, err := h.service.ListActivity(r.Context(), reviewerID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	dtos := make([]cycleDTO, 0, len(cycles))
	for _, c := range cycles {
		dtos = append(dtos, toCycleDTO(c))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"activity": dtos})
}

func (h *Handler) suggestionStatusHandler(status SuggestionStatus) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		contributorID, _ := auth.UserIDFromContext(r.Context())
		id, err := uuid.Parse(chi.URLParam(r, "id"))
		if err != nil {
			httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid suggestion id")
			return
		}
		if err := h.service.SetSuggestionStatus(r.Context(), id, contributorID, status); err != nil {
			s, code := mapError(err)
			httpx.Error(w, s, code, err.Error())
			return
		}
		httpx.JSON(w, http.StatusNoContent, nil)
	}
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, articles.ErrNotFound):
		return http.StatusNotFound, "article_not_found"
	case errors.Is(err, articles.ErrForbidden), errors.Is(err, ErrForbidden):
		return http.StatusForbidden, "forbidden"
	case errors.Is(err, ErrNotReviewable):
		return http.StatusConflict, "article_not_reviewable"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
