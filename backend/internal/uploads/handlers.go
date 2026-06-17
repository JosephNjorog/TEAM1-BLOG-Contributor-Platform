// Package uploads handles inline images dropped into the article body
// while drafting - distinct from internal/banners, which has its own
// stricter format/dimension rules for cover banners.
package uploads

import (
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/cloudinary"
	"team1blog/backend/internal/httpx"
)

const maxInlineImageBytes = 8 * 1024 * 1024

type Handler struct {
	articlesRepo *articles.Repository
	uploader     cloudinary.Uploader
}

func NewHandler(articlesRepo *articles.Repository, uploader cloudinary.Uploader) *Handler {
	return &Handler{articlesRepo: articlesRepo, uploader: uploader}
}

func (h *Handler) UploadInline(w http.ResponseWriter, r *http.Request) {
	contributorID, _ := auth.UserIDFromContext(r.Context())
	articleID, err := uuid.Parse(chi.URLParam(r, "articleId"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid article id")
		return
	}

	a, err := h.articlesRepo.GetByID(r.Context(), articleID)
	if err != nil {
		if errors.Is(err, articles.ErrNotFound) {
			httpx.Error(w, http.StatusNotFound, "article_not_found", err.Error())
			return
		}
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	if a.ContributorID != contributorID {
		httpx.Error(w, http.StatusForbidden, "forbidden", "you do not have access to this article")
		return
	}
	if !a.Status.Editable() {
		httpx.Error(w, http.StatusConflict, "article_not_editable", "article cannot be edited in its current status")
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, maxInlineImageBytes)
	if err := r.ParseMultipartForm(maxInlineImageBytes); err != nil {
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

	folder := fmt.Sprintf("team1/articles/%s/inline", articleID)
	url, err := h.uploader.Upload(r.Context(), data, header.Filename, folder)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "upload_failed", err.Error())
		return
	}
	httpx.JSON(w, http.StatusCreated, map[string]string{"url": url})
}
