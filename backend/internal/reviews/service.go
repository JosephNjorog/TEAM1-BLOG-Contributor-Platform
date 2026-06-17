package reviews

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/audit"
	"team1blog/backend/internal/email"
	"team1blog/backend/internal/notifications"
	"team1blog/backend/internal/users"
)

var ErrNotReviewable = errors.New("article is not awaiting review")

type Service struct {
	repo          *Repository
	articlesRepo  *articles.Repository
	usersRepo     *users.Repository
	notifications *notifications.Repository
	audit         *audit.Logger
	mailer        email.Sender
	appURL        string
}

func NewService(
	repo *Repository,
	articlesRepo *articles.Repository,
	usersRepo *users.Repository,
	notificationsRepo *notifications.Repository,
	auditLogger *audit.Logger,
	mailer email.Sender,
	appURL string,
) *Service {
	return &Service{
		repo:          repo,
		articlesRepo:  articlesRepo,
		usersRepo:     usersRepo,
		notifications: notificationsRepo,
		audit:         auditLogger,
		mailer:        mailer,
		appURL:        appURL,
	}
}

type SubmitReviewInput struct {
	ArticleID   uuid.UUID
	ReviewerID  uuid.UUID
	Decision    Decision
	Summary     string
	Suggestions []SuggestionInput
}

func (s *Service) SubmitReview(ctx context.Context, in SubmitReviewInput) (*articles.Article, error) {
	a, err := s.articlesRepo.GetByID(ctx, in.ArticleID)
	if err != nil {
		return nil, err
	}
	if a.Status != articles.StatusSubmitted && a.Status != articles.StatusResubmitted {
		return nil, ErrNotReviewable
	}

	target := articles.StatusEditorialApproved
	if in.Decision == DecisionChangesRequested {
		target = articles.StatusChangesRequested
	}
	if err := articles.ValidateTransition(a.Status, target); err != nil {
		return nil, err
	}

	if err := s.articlesRepo.SetReviewDecision(ctx, in.ArticleID, in.ReviewerID, target); err != nil {
		return nil, err
	}
	if err := s.articlesRepo.IncrementReviewCycle(ctx, in.ArticleID); err != nil {
		return nil, err
	}

	cycle, err := s.repo.CreateCycle(ctx, in.ArticleID, in.ReviewerID, in.Decision, in.Summary)
	if err != nil {
		return nil, err
	}
	for _, sug := range in.Suggestions {
		if _, err := s.repo.CreateSuggestion(ctx, in.ArticleID, cycle.ID, in.ReviewerID, sug); err != nil {
			return nil, err
		}
	}

	s.notifyOutcome(ctx, a, in.Decision)

	_ = s.audit.Log(ctx, &in.ReviewerID, "article_reviewed", "article", &in.ArticleID, map[string]any{
		"decision": in.Decision,
		"from":     a.Status,
		"to":       target,
	})

	return s.articlesRepo.GetByID(ctx, in.ArticleID)
}

func (s *Service) notifyOutcome(ctx context.Context, a *articles.Article, decision Decision) {
	dashboardURL := fmt.Sprintf("%s/articles/%s", s.appURL, a.ID)

	if decision == DecisionApproved {
		_ = s.notifications.CreateForRole(ctx, string(users.RoleGraphicDesigner), notifications.TypeArticleApproved, &a.ID,
			fmt.Sprintf("%q is ready for a banner", a.Title))
		_, _ = s.notifications.Create(ctx, a.ContributorID, notifications.TypeArticleApproved, &a.ID,
			fmt.Sprintf("%q was approved by editorial", a.Title))

		subject, html := email.ArticleApprovedEmail(a.Title, dashboardURL)
		if designers, err := s.usersRepo.ListByRole(ctx, users.RoleGraphicDesigner); err == nil {
			for _, d := range designers {
				_ = s.mailer.Send(ctx, d.Email, subject, html)
			}
		}
		if contributor, err := s.usersRepo.GetByID(ctx, a.ContributorID); err == nil {
			_ = s.mailer.Send(ctx, contributor.Email, subject, html)
		}
		return
	}

	_, _ = s.notifications.Create(ctx, a.ContributorID, notifications.TypeChangesRequested, &a.ID,
		fmt.Sprintf("Changes were requested on %q", a.Title))
	if contributor, err := s.usersRepo.GetByID(ctx, a.ContributorID); err == nil {
		subject, html := email.ReviewFeedbackEmail(a.Title, dashboardURL)
		_ = s.mailer.Send(ctx, contributor.Email, subject, html)
	}
}

func (s *Service) ListForArticle(ctx context.Context, articleID uuid.UUID) ([]*ReviewCycle, []*Suggestion, error) {
	cycles, err := s.repo.ListCyclesForArticle(ctx, articleID)
	if err != nil {
		return nil, nil, err
	}
	suggestions, err := s.repo.ListSuggestionsForArticle(ctx, articleID)
	if err != nil {
		return nil, nil, err
	}
	return cycles, suggestions, nil
}

func (s *Service) ListActivity(ctx context.Context, reviewerID uuid.UUID) ([]*ReviewCycle, error) {
	return s.repo.ListActivityForReviewer(ctx, reviewerID)
}

var ErrForbidden = errors.New("you do not have access to this suggestion")

// SetSuggestionStatus lets a contributor accept or reject a suggestion on
// their own article.
func (s *Service) SetSuggestionStatus(ctx context.Context, suggestionID, contributorID uuid.UUID, status SuggestionStatus) error {
	articleID, err := s.repo.GetSuggestionArticleID(ctx, suggestionID)
	if err != nil {
		return err
	}
	a, err := s.articlesRepo.GetByID(ctx, articleID)
	if err != nil {
		return err
	}
	if a.ContributorID != contributorID {
		return ErrForbidden
	}
	return s.repo.SetSuggestionStatus(ctx, suggestionID, status)
}
