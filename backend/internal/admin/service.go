package admin

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/audit"
	"team1blog/backend/internal/users"
)

var ErrReasonRequired = errors.New("a reason is required to override an article's status")

type Service struct {
	repo      *Repository
	usersRepo *users.Repository
	audit     *audit.Logger
}

func NewService(repo *Repository, usersRepo *users.Repository, auditLogger *audit.Logger) *Service {
	return &Service{repo: repo, usersRepo: usersRepo, audit: auditLogger}
}

func (s *Service) ListContributors(ctx context.Context) ([]*ContributorSummary, error) {
	return s.repo.ListContributors(ctx)
}

func (s *Service) ListPendingInvitations(ctx context.Context) ([]*PendingInvitation, error) {
	return s.repo.ListPendingInvitations(ctx)
}

func (s *Service) GetOverview(ctx context.Context) (*Overview, error) {
	return s.repo.GetOverview(ctx)
}

func (s *Service) GetAnalytics(ctx context.Context) (*PlatformMetrics, error) {
	contributorMetrics, err := s.repo.GetContributorMetrics(ctx)
	if err != nil {
		return nil, err
	}
	pubVolume, err := s.repo.GetPublicationVolume(ctx)
	if err != nil {
		return nil, err
	}
	payVolume, err := s.repo.GetPaymentVolume(ctx)
	if err != nil {
		return nil, err
	}
	avgDays, err := s.repo.GetAvgPipelineDays(ctx)
	if err != nil {
		return nil, err
	}
	return &PlatformMetrics{
		ContributorMetrics: contributorMetrics,
		PublicationVolume:  pubVolume,
		PaymentVolume:      payVolume,
		AvgPipelineDays:    avgDays,
	}, nil
}

func (s *Service) SetUserStatus(ctx context.Context, actorID, userID uuid.UUID, status users.Status) error {
	if err := s.usersRepo.SetStatus(ctx, userID, status); err != nil {
		return err
	}
	return s.audit.Log(ctx, &actorID, "user_status_changed", "user", &userID, map[string]any{"status": status})
}

func (s *Service) UpdateUserRole(ctx context.Context, actorID, userID uuid.UUID, role users.Role) error {
	if !role.Valid() {
		return errors.New("invalid role")
	}
	if err := s.usersRepo.UpdateRole(ctx, userID, role); err != nil {
		return err
	}
	return s.audit.Log(ctx, &actorID, "user_role_changed", "user", &userID, map[string]any{"role": role})
}

// OverrideArticleStatus lets a Super Admin force an article into any state
// in exceptional circumstances, bypassing the normal transition rules -
// per the PRD this requires a reason, which is the only thing distinguishing
// it from a routine transition in the audit log.
func (s *Service) OverrideArticleStatus(ctx context.Context, actorID, articleID uuid.UUID, status articles.Status, reason string) error {
	if reason == "" {
		return ErrReasonRequired
	}
	if err := s.repo.OverrideArticleStatus(ctx, articleID, string(status)); err != nil {
		return err
	}
	return s.audit.Log(ctx, &actorID, "article_status_overridden", "article", &articleID, map[string]any{
		"to":     status,
		"reason": reason,
	})
}
