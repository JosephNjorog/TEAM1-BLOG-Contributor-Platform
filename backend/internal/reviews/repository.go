package reviews

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateCycle(ctx context.Context, articleID, reviewerID uuid.UUID, decision Decision, summary string) (*ReviewCycle, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO review_cycles (article_id, reviewer_id, decision, summary)
		VALUES ($1, $2, $3, $4)
		RETURNING id, article_id, reviewer_id, decision, summary, created_at
	`, articleID, reviewerID, decision, summary)

	var c ReviewCycle
	if err := row.Scan(&c.ID, &c.ArticleID, &c.ReviewerID, &c.Decision, &c.Summary, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) CreateSuggestion(ctx context.Context, articleID, cycleID, reviewerID uuid.UUID, in SuggestionInput) (*Suggestion, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO suggestions (article_id, review_cycle_id, reviewer_id, range_start, range_end, suggestion_text)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, article_id, review_cycle_id, reviewer_id, range_start, range_end, suggestion_text, status, created_at
	`, articleID, cycleID, reviewerID, in.RangeStart, in.RangeEnd, in.SuggestionText)

	var s Suggestion
	if err := row.Scan(&s.ID, &s.ArticleID, &s.ReviewCycleID, &s.ReviewerID, &s.RangeStart, &s.RangeEnd, &s.SuggestionText, &s.Status, &s.CreatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListSuggestionsForArticle(ctx context.Context, articleID uuid.UUID) ([]*Suggestion, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT s.id, s.article_id, s.review_cycle_id, s.reviewer_id, u.name, s.range_start, s.range_end, s.suggestion_text, s.status, s.created_at
		FROM suggestions s
		JOIN users u ON u.id = s.reviewer_id
		WHERE s.article_id = $1
		ORDER BY s.range_start ASC, s.created_at ASC
	`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Suggestion
	for rows.Next() {
		s, err := scanSuggestion(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, s)
	}
	return out, rows.Err()
}

func scanSuggestion(row pgx.Row) (*Suggestion, error) {
	var s Suggestion
	if err := row.Scan(&s.ID, &s.ArticleID, &s.ReviewCycleID, &s.ReviewerID, &s.ReviewerName, &s.RangeStart, &s.RangeEnd, &s.SuggestionText, &s.Status, &s.CreatedAt); err != nil {
		return nil, err
	}
	return &s, nil
}

// GetSuggestionArticleID is used to authorize accept/reject actions: the
// caller must own the article the suggestion belongs to.
func (r *Repository) GetSuggestionArticleID(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	var articleID uuid.UUID
	err := r.pool.QueryRow(ctx, `SELECT article_id FROM suggestions WHERE id = $1`, id).Scan(&articleID)
	return articleID, err
}

func (r *Repository) SetSuggestionStatus(ctx context.Context, id uuid.UUID, status SuggestionStatus) error {
	_, err := r.pool.Exec(ctx, `UPDATE suggestions SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *Repository) ListCyclesForArticle(ctx context.Context, articleID uuid.UUID) ([]*ReviewCycle, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT rc.id, rc.article_id, rc.reviewer_id, u.name, rc.decision, rc.summary, rc.created_at
		FROM review_cycles rc
		JOIN users u ON u.id = rc.reviewer_id
		WHERE rc.article_id = $1
		ORDER BY rc.created_at ASC
	`, articleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*ReviewCycle
	for rows.Next() {
		c, err := scanCycle(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// ListActivityForReviewer backs the moderator's activity log: every review
// decision they've made, newest first, with the article title for context.
func (r *Repository) ListActivityForReviewer(ctx context.Context, reviewerID uuid.UUID) ([]*ReviewCycle, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT rc.id, rc.article_id, a.title, rc.reviewer_id, u.name, rc.decision, rc.summary, rc.created_at
		FROM review_cycles rc
		JOIN users u ON u.id = rc.reviewer_id
		JOIN articles a ON a.id = rc.article_id
		WHERE rc.reviewer_id = $1
		ORDER BY rc.created_at DESC
	`, reviewerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*ReviewCycle
	for rows.Next() {
		var c ReviewCycle
		if err := rows.Scan(&c.ID, &c.ArticleID, &c.ArticleTitle, &c.ReviewerID, &c.ReviewerName, &c.Decision, &c.Summary, &c.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}

func scanCycle(row pgx.Row) (*ReviewCycle, error) {
	var c ReviewCycle
	if err := row.Scan(&c.ID, &c.ArticleID, &c.ReviewerID, &c.ReviewerName, &c.Decision, &c.Summary, &c.CreatedAt); err != nil {
		return nil, err
	}
	return &c, nil
}
