package substack

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"team1blog/backend/internal/users"
)

type Service struct {
	repo      *Repository
	usersRepo *users.Repository
	fetcher   Fetcher
}

func NewService(repo *Repository, usersRepo *users.Repository, fetcher Fetcher) *Service {
	return &Service{repo: repo, usersRepo: usersRepo, fetcher: fetcher}
}

// Sync fetches the publication's post history and matches each post to a
// contributor by author/byline mention, per the PRD - there's no stronger
// signal available from a public RSS feed than a name match.
func (s *Service) Sync(ctx context.Context) (int, error) {
	posts, err := s.fetcher.FetchPosts(ctx)
	if err != nil {
		return 0, err
	}

	contributors, err := s.usersRepo.ListByRole(ctx, users.RoleContributor)
	if err != nil {
		return 0, err
	}

	synced := 0
	for _, post := range posts {
		var matchedID *uuid.UUID
		if id := matchContributor(post.Author, contributors); id != nil {
			matchedID = id
		}
		if err := s.repo.Upsert(ctx, matchedID, post); err != nil {
			return synced, err
		}
		synced++
	}
	return synced, nil
}

func matchContributor(author string, contributors []*users.User) *uuid.UUID {
	author = strings.TrimSpace(strings.ToLower(author))
	if author == "" {
		return nil
	}
	for _, c := range contributors {
		if strings.ToLower(c.Name) == author || strings.Contains(author, strings.ToLower(c.Name)) {
			id := c.ID
			return &id
		}
	}
	return nil
}

func (s *Service) ListForContributor(ctx context.Context, contributorID uuid.UUID) ([]ContributorPost, error) {
	return s.repo.ListForContributor(ctx, contributorID)
}
