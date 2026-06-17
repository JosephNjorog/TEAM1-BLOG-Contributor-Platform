package banners

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/audit"
	"team1blog/backend/internal/cloudinary"
	"team1blog/backend/internal/email"
	"team1blog/backend/internal/notifications"
	"team1blog/backend/internal/users"
)

var (
	ErrNotAwaitingBanner = errors.New("article is not awaiting a banner")
	ErrNoBannerUploaded  = errors.New("upload a banner before marking it ready")
)

type Service struct {
	repo          *Repository
	articlesRepo  *articles.Repository
	usersRepo     *users.Repository
	notifications *notifications.Repository
	uploader      cloudinary.Uploader
	audit         *audit.Logger
	mailer        email.Sender
	appURL        string
}

func NewService(
	repo *Repository,
	articlesRepo *articles.Repository,
	usersRepo *users.Repository,
	notificationsRepo *notifications.Repository,
	uploader cloudinary.Uploader,
	auditLogger *audit.Logger,
	mailer email.Sender,
	appURL string,
) *Service {
	return &Service{
		repo:          repo,
		articlesRepo:  articlesRepo,
		usersRepo:     usersRepo,
		notifications: notificationsRepo,
		uploader:      uploader,
		audit:         auditLogger,
		mailer:        mailer,
		appURL:        appURL,
	}
}

func (s *Service) Upload(ctx context.Context, articleID, designerID uuid.UUID, file []byte, filename string) (*Banner, *articles.Article, error) {
	a, err := s.articlesRepo.GetByID(ctx, articleID)
	if err != nil {
		return nil, nil, err
	}
	if a.Status != articles.StatusEditorialApproved {
		return nil, nil, ErrNotAwaitingBanner
	}
	if err := cloudinary.ValidateBanner(file); err != nil {
		return nil, nil, err
	}

	folder := fmt.Sprintf("team1/articles/%s/banner", articleID)
	url, err := s.uploader.Upload(ctx, file, filename, folder)
	if err != nil {
		return nil, nil, err
	}

	if err := s.articlesRepo.AssignDesigner(ctx, articleID, designerID); err != nil {
		return nil, nil, err
	}
	if err := s.articlesRepo.SetBannerURL(ctx, articleID, url); err != nil {
		return nil, nil, err
	}
	banner, err := s.repo.Create(ctx, articleID, designerID, url)
	if err != nil {
		return nil, nil, err
	}

	_ = s.audit.Log(ctx, &designerID, "banner_uploaded", "article", &articleID, map[string]any{"url": url})

	updated, err := s.articlesRepo.GetByID(ctx, articleID)
	return banner, updated, err
}

func (s *Service) MarkReady(ctx context.Context, articleID, designerID uuid.UUID) (*articles.Article, error) {
	a, err := s.articlesRepo.GetByID(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if a.Status != articles.StatusEditorialApproved {
		return nil, ErrNotAwaitingBanner
	}
	if a.CloudinaryBannerURL == nil {
		return nil, ErrNoBannerUploaded
	}

	banner, err := s.repo.GetLatestForArticle(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if err := s.repo.MarkReady(ctx, banner.ID); err != nil {
		return nil, err
	}
	if err := s.articlesRepo.MarkBannerReady(ctx, articleID); err != nil {
		return nil, err
	}

	dashboardURL := fmt.Sprintf("%s/articles/%s", s.appURL, articleID)
	_ = s.notifications.CreateForRole(ctx, string(users.RolePublisher), notifications.TypeBannerReady, &articleID,
		fmt.Sprintf("%q has a banner ready to publish", a.Title))
	if publishers, err := s.usersRepo.ListByRole(ctx, users.RolePublisher); err == nil {
		subject, html := email.BannerReadyEmail(a.Title, dashboardURL)
		for _, p := range publishers {
			_ = s.mailer.Send(ctx, p.Email, subject, html)
		}
	}

	_ = s.audit.Log(ctx, &designerID, "banner_marked_ready", "article", &articleID, nil)

	return s.articlesRepo.GetByID(ctx, articleID)
}
