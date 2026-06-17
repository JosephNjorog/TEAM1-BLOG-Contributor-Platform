package banners

import (
	"errors"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/cloudinary"
	"team1blog/backend/internal/httpx"
)

const maxUploadBytes = 6 * 1024 * 1024 // a little above BannerMaxBytes so we can return a clean validation error instead of a truncated read

type Handler struct {
	service     *Service
	articlesSvc *articles.Service
}

func NewHandler(service *Service, articlesSvc *articles.Service) *Handler {
	return &Handler{service: service, articlesSvc: articlesSvc}
}

func (h *Handler) Upload(w http.ResponseWriter, r *http.Request) {
	designerID, _ := auth.UserIDFromContext(r.Context())
	articleID, err := uuid.Parse(chi.URLParam(r, "articleId"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
	if err := r.ParseMultipartForm(maxUploadBytes); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "file too large or malformed upload")
		return
	}
	file, header, err := r.FormFile("file")
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "missing file field")
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "could not read uploaded file")
		return
	}

	banner, article, err := h.service.Upload(r.Context(), articleID, designerID, data, header.Filename)
	if err != nil {
		status, code := mapError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]any{
		"id":            banner.ID,
		"cloudinaryUrl": banner.CloudinaryURL,
		"uploadedAt":    banner.UploadedAt,
		"article":       articles.ToDTO(article),
	})
}

func (h *Handler) MarkReady(w http.ResponseWriter, r *http.Request) {
	designerID, _ := auth.UserIDFromContext(r.Context())
	articleID, err := uuid.Parse(chi.URLParam(r, "articleId"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	a, err := h.service.MarkReady(r.Context(), articleID, designerID)
	if err != nil {
		status, code := mapError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, articles.ToDTO(a))
}

func (h *Handler) GetLatest(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	claims, _ := auth.ClaimsFromContext(r.Context())
	articleID, err := uuid.Parse(chi.URLParam(r, "articleId"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}
	if _, err := h.articlesSvc.GetVisible(r.Context(), articleID, userID, claims.Role); err != nil {
		status, code := mapError(err)
		httpx.Error(w, status, code, err.Error())
		return
	}
	banner, err := h.service.repo.GetLatestForArticle(r.Context(), articleID)
	if err != nil {
		httpx.Error(w, http.StatusNotFound, "banner_not_found", "no banner uploaded yet")
		return
	}
	httpx.JSON(w, http.StatusOK, banner)
}

func mapError(err error) (int, string) {
	switch {
	case errors.Is(err, articles.ErrNotFound):
		return http.StatusNotFound, "article_not_found"
	case errors.Is(err, articles.ErrForbidden):
		return http.StatusForbidden, "forbidden"
	case errors.Is(err, ErrNotAwaitingBanner):
		return http.StatusConflict, "article_not_awaiting_banner"
	case errors.Is(err, ErrNoBannerUploaded):
		return http.StatusConflict, "no_banner_uploaded"
	case errors.Is(err, cloudinary.ErrImageTooLarge), errors.Is(err, cloudinary.ErrImageTooSmall), errors.Is(err, cloudinary.ErrUnsupportedFormat):
		return http.StatusBadRequest, "invalid_banner_image"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}
