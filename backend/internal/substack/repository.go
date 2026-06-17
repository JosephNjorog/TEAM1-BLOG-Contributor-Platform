package substack

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Upsert records a fetched post, matched to a contributor if one was
// found. Re-syncing the same post (by substack_post_id) just refreshes
// synced_at instead of duplicating it.
func (r *Repository) Upsert(ctx context.Context, contributorID *uuid.UUID, p Post) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO substack_articles (contributor_id, substack_post_id, title, url, published_at, synced_at)
		VALUES ($1, $2, $3, $4, $5, now())
		ON CONFLICT (substack_post_id) DO UPDATE SET
			contributor_id = EXCLUDED.contributor_id,
			title = EXCLUDED.title,
			url = EXCLUDED.url,
			synced_at = now()
	`, contributorID, p.SubstackPostID, p.Title, p.URL, p.PublishedAt)
	return err
}

type ContributorPost struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"publishedAt"`
}

func (r *Repository) ListForContributor(ctx context.Context, contributorID uuid.UUID) ([]ContributorPost, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, url, published_at
		FROM substack_articles
		WHERE contributor_id = $1
		ORDER BY published_at DESC
	`, contributorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ContributorPost
	for rows.Next() {
		var p ContributorPost
		if err := rows.Scan(&p.ID, &p.Title, &p.URL, &p.PublishedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
