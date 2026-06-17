package articles

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"team1blog/backend/internal/audit"
	"team1blog/backend/internal/email"
	"team1blog/backend/internal/notifications"
	"team1blog/backend/internal/users"
)

var ErrNotEditable = errors.New("article cannot be edited in its current status")

type Service struct {
	repo          *Repository
	usersRepo     *users.Repository
	notifications *notifications.Repository
	audit         *audit.Logger
	mailer        email.Sender
	appURL        string
}

func NewService(repo *Repository, usersRepo *users.Repository, notificationsRepo *notifications.Repository, auditLogger *audit.Logger, mailer email.Sender, appURL string) *Service {
	return &Service{
		repo:          repo,
		usersRepo:     usersRepo,
		notifications: notificationsRepo,
		audit:         auditLogger,
		mailer:        mailer,
		appURL:        appURL,
	}
}

func (s *Service) CreateDraft(ctx context.Context, contributorID uuid.UUID, title, content, sourceCitation string) (*Article, error) {
	return s.repo.Create(ctx, contributorID, title, content, sourceCitation, CountWords(content))
}

type UpdateInput struct {
	Title          string
	Content        string
	SourceCitation string
}

// UpdateDraft handles both the 60s autosave and manual save. It's only
// valid while the article is in an editable state (draft, or changes
// requested - mid-revision before resubmission).
func (s *Service) UpdateDraft(ctx context.Context, articleID, contributorID uuid.UUID, in UpdateInput) (*Article, error) {
	a, err := s.repo.GetByID(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if a.ContributorID != contributorID {
		return nil, ErrForbidden
	}
	if !a.Status.Editable() {
		return nil, ErrNotEditable
	}

	if err := s.repo.UpdateDraft(ctx, articleID, UpdateDraftInput{
		Title:          in.Title,
		Content:        in.Content,
		SourceCitation: in.SourceCitation,
		WordCount:      CountWords(in.Content),
	}); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, articleID)
}

func (s *Service) GetVisible(ctx context.Context, id, userID uuid.UUID, role users.Role) (*Article, error) {
	return s.repo.GetVisible(ctx, id, userID, role)
}

func (s *Service) ListVisible(ctx context.Context, userID uuid.UUID, role users.Role) ([]*Article, error) {
	return s.repo.ListVisible(ctx, userID, role)
}

func (s *Service) Delete(ctx context.Context, articleID, contributorID uuid.UUID) error {
	a, err := s.repo.GetByID(ctx, articleID)
	if err != nil {
		return err
	}
	if a.ContributorID != contributorID {
		return ErrForbidden
	}
	if a.Status != StatusDraft {
		return ErrNotEditable
	}
	return s.repo.Delete(ctx, articleID)
}

// Submit moves a draft (or a changes-requested article being resubmitted)
// into the review queue, notifying every active moderator by email and
// in-app notification, plus acknowledging the submission to the contributor.
func (s *Service) Submit(ctx context.Context, articleID, contributorID uuid.UUID) (*Article, error) {
	a, err := s.repo.GetByID(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if a.ContributorID != contributorID {
		return nil, ErrForbidden
	}

	target := StatusSubmitted
	if a.Status == StatusChangesRequested {
		target = StatusResubmitted
	}
	if err := ValidateTransition(a.Status, target); err != nil {
		return nil, err
	}

	if err := s.repo.MarkSubmitted(ctx, articleID, target); err != nil {
		return nil, err
	}

	dashboardURL := fmt.Sprintf("%s/articles/%s", s.appURL, articleID)
	moderatorMsg := fmt.Sprintf("%q was submitted for review by %s", a.Title, a.ContributorName)
	notifType := notifications.TypeArticleSubmitted
	if target == StatusResubmitted {
		notifType = notifications.TypeArticleResubmitted
		moderatorMsg = fmt.Sprintf("%q was resubmitted after changes by %s", a.Title, a.ContributorName)
	}
	_ = s.notifications.CreateForRole(ctx, string(users.RoleModerator), notifType, &articleID, moderatorMsg)

	if moderators, err := s.usersRepo.ListByRole(ctx, users.RoleModerator); err == nil {
		subject, html := email.SubmissionAcknowledgedEmail(a.Title, dashboardURL)
		for _, m := range moderators {
			_ = s.mailer.Send(ctx, m.Email, subject, html)
		}
	}

	if contributor, err := s.usersRepo.GetByID(ctx, contributorID); err == nil {
		subject, html := email.SubmissionAcknowledgedEmail(a.Title, dashboardURL)
		_ = s.mailer.Send(ctx, contributor.Email, subject, html)
		_, _ = s.notifications.Create(ctx, contributorID, notifications.TypeArticleSubmitted, &articleID, fmt.Sprintf("%q was submitted for review", a.Title))
	}

	_ = s.audit.Log(ctx, &contributorID, "article_submitted", "article", &articleID, map[string]any{"from": a.Status, "to": target})

	return s.repo.GetByID(ctx, articleID)
}

var ErrNotReadyToPublish = errors.New("article is not ready to publish")

// Publish confirms a banner_uploaded article is live, recording the
// Substack URL and notifying the contributor, the moderator who approved
// it, and every Super Admin (who needs to act on releasing payment).
func (s *Service) Publish(ctx context.Context, articleID, publisherID uuid.UUID, substackURL string) (*Article, error) {
	a, err := s.repo.GetByID(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if a.Status != StatusBannerUploaded {
		return nil, ErrNotReadyToPublish
	}
	if err := ValidateTransition(a.Status, StatusPublished); err != nil {
		return nil, err
	}

	if err := s.repo.MarkPublished(ctx, articleID, publisherID, substackURL); err != nil {
		return nil, err
	}

	dashboardURL := fmt.Sprintf("%s/articles/%s", s.appURL, articleID)
	subject, html := email.ArticlePublishedEmail(a.Title, dashboardURL)

	_, _ = s.notifications.Create(ctx, a.ContributorID, notifications.TypeArticlePublished, &articleID,
		fmt.Sprintf("%q is now live", a.Title))
	if contributor, err := s.usersRepo.GetByID(ctx, a.ContributorID); err == nil {
		_ = s.mailer.Send(ctx, contributor.Email, subject, html)
	}

	if a.ReviewerID != nil {
		_, _ = s.notifications.Create(ctx, *a.ReviewerID, notifications.TypeArticlePublished, &articleID,
			fmt.Sprintf("%q was published", a.Title))
		if reviewer, err := s.usersRepo.GetByID(ctx, *a.ReviewerID); err == nil {
			_ = s.mailer.Send(ctx, reviewer.Email, subject, html)
		}
	}

	_ = s.notifications.CreateForRole(ctx, string(users.RoleSuperAdmin), notifications.TypeArticlePublished, &articleID,
		fmt.Sprintf("%q was published - payment release needed", a.Title))

	_ = s.audit.Log(ctx, &publisherID, "article_published", "article", &articleID, map[string]any{"substackUrl": substackURL})

	return s.repo.GetByID(ctx, articleID)
}
