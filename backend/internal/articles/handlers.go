package articles

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/httpx"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

type articleDTO struct {
	ID                  uuid.UUID `json:"id"`
	ContributorID       uuid.UUID `json:"contributorId"`
	ContributorName     string    `json:"contributorName"`
	Title               string    `json:"title"`
	Content             string    `json:"content"`
	Status              Status    `json:"status"`
	WordCount           int       `json:"wordCount"`
	SourceCitation      *string   `json:"sourceCitation"`
	SubstackURL         *string   `json:"substackUrl"`
	CloudinaryBannerURL *string   `json:"cloudinaryBannerUrl"`
	ReviewCycleCount    int       `json:"reviewCycleCount"`
	ReviewerName        *string   `json:"reviewerName"`
	CreatedAt           string    `json:"createdAt"`
	UpdatedAt           string    `json:"updatedAt"`
	SubmittedAt         *string   `json:"submittedAt"`
	PublishedAt         *string   `json:"publishedAt"`
}

func ToDTO(a *Article) articleDTO {
	dto := articleDTO{
		ID:                  a.ID,
		ContributorID:       a.ContributorID,
		ContributorName:     a.ContributorName,
		Title:               a.Title,
		Content:             a.Content,
		Status:              a.Status,
		WordCount:           a.WordCount,
		SourceCitation:      a.SourceCitation,
		SubstackURL:         a.SubstackURL,
		CloudinaryBannerURL: a.CloudinaryBannerURL,
		ReviewCycleCount:    a.ReviewCycleCount,
		ReviewerName:        a.ReviewerName,
		CreatedAt:           a.CreatedAt.Format(timeFmt),
		UpdatedAt:           a.UpdatedAt.Format(timeFmt),
	}
	if a.SubmittedAt != nil {
		v := a.SubmittedAt.Format(timeFmt)
		dto.SubmittedAt = &v
	}
	if a.PublishedAt != nil {
		v := a.PublishedAt.Format(timeFmt)
		dto.PublishedAt = &v
	}
	return dto
}

const timeFmt = "2006-01-02T15:04:05Z07:00"

type createRequest struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	SourceCitation string `json:"sourceCitation"`
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	var req createRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}
	if req.Title == "" {
		req.Title = "Untitled draft"
	}
	a, err := h.service.CreateDraft(r.Context(), userID, req.Title, req.Content, req.SourceCitation)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusCreated, ToDTO(a))
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	claims, _ := auth.ClaimsFromContext(r.Context())
	list, err := h.service.ListVisible(r.Context(), userID, claims.Role)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	dtos := make([]articleDTO, 0, len(list))
	for _, a := range list {
		dtos = append(dtos, ToDTO(a))
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"articles": dtos})
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	claims, _ := auth.ClaimsFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	a, err := h.service.GetVisible(r.Context(), id, userID, claims.Role)
	if err != nil {
		status, code := mapArticleError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, ToDTO(a))
}

type updateRequest struct {
	Title          string `json:"title"`
	Content        string `json:"content"`
	SourceCitation string `json:"sourceCitation"`
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	var req updateRequest
	if err := httpx.DecodeJSON(r, &req); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid request body")
		return
	}
	a, err := h.service.UpdateDraft(r.Context(), id, userID, UpdateInput{
		Title:          req.Title,
		Content:        req.Content,
		SourceCitation: req.SourceCitation,
	})
	if err != nil {
		status, code := mapArticleError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, ToDTO(a))
}

func (h *Handler) Submit(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	a, err := h.service.Submit(r.Context(), id, userID)
	if err != nil {
		status, code := mapArticleError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, ToDTO(a))
}

type publishRequest struct {
	SubstackURL string `json:"substackUrl"`
}

func (h *Handler) Publish(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	var req publishRequest
	if err := httpx.DecodeJSON(r, &req); err != nil || req.SubstackURL == "" {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "substackUrl is required")
		return
	}
	a, err := h.service.Publish(r.Context(), id, userID, req.SubstackURL)
	if err != nil {
		status, code := mapArticleError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, ToDTO(a))
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		status, code := mapArticleError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

func mapArticleError(err error) (int, string) {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound, "article_not_found"
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden, "forbidden"
	case errors.Is(err, ErrNotEditable):
		return http.StatusConflict, "article_not_editable"
	case errors.Is(err, ErrNotReadyToPublish):
		return http.StatusConflict, "article_not_ready_to_publish"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
