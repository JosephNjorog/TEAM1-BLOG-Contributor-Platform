package notifications

import (
	"context"

	"github.com/google/uuid"
)

// Service wraps the repository with realtime delivery: every notification
// written to the database is also pushed over WebSocket to any connection
// its recipient currently has open. Exposes the same Create/CreateForRole
// signatures as the repository so callers elsewhere in the codebase don't
// need to change anything beyond their dependency's type.
type Service struct {
	repo *Repository
	hub  *Hub
}

func NewService(repo *Repository, hub *Hub) *Service {
	return &Service{repo: repo, hub: hub}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, t Type, articleID *uuid.UUID, message string) (*Notification, error) {
	n, err := s.repo.Create(ctx, userID, t, articleID, message)
	if err != nil {
		return nil, err
	}
	s.hub.Push(userID, n)
	return n, nil
}

func (s *Service) CreateForRole(ctx context.Context, role string, t Type, articleID *uuid.UUID, message string) ([]*Notification, error) {
	created, err := s.repo.CreateForRole(ctx, role, t, articleID, message)
	if err != nil {
		return nil, err
	}
	for _, n := range created {
		s.hub.Push(n.UserID, n)
	}
	return created, nil
}

func (s *Service) ListForUser(ctx context.Context, userID uuid.UUID, limit int) ([]*Notification, error) {
	return s.repo.ListForUser(ctx, userID, limit)
}

func (s *Service) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	return s.repo.CountUnread(ctx, userID)
}

func (s *Service) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	return s.repo.MarkRead(ctx, id, userID)
}

func (s *Service) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	return s.repo.MarkAllRead(ctx, userID)
}
