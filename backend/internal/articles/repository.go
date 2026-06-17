package articles

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"team1blog/backend/internal/users"
)

var (
	ErrNotFound  = errors.New("article not found")
	ErrForbidden = errors.New("you do not have access to this article")
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const baseSelect = `
	SELECT a.id, a.contributor_id, c.name, a.reviewer_id, r.name, a.designer_id, a.publisher_id,
	       a.title, a.content, a.source_citation, a.status, a.word_count, a.review_cycle_count,
	       a.substack_url, a.cloudinary_banner_url, a.created_at, a.updated_at, a.submitted_at, a.published_at
	FROM articles a
	JOIN users c ON c.id = a.contributor_id
	LEFT JOIN users r ON r.id = a.reviewer_id
`

func scanArticle(row pgx.Row) (*Article, error) {
	var a Article
	err := row.Scan(
		&a.ID, &a.ContributorID, &a.ContributorName, &a.ReviewerID, &a.ReviewerName, &a.DesignerID, &a.PublisherID,
		&a.Title, &a.Content, &a.SourceCitation, &a.Status, &a.WordCount, &a.ReviewCycleCount,
		&a.SubstackURL, &a.CloudinaryBannerURL, &a.CreatedAt, &a.UpdatedAt, &a.SubmittedAt, &a.PublishedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &a, nil
}

func (r *Repository) Create(ctx context.Context, contributorID uuid.UUID, title, content, sourceCitation string, wordCount int) (*Article, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO articles (contributor_id, title, content, source_citation, word_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, contributorID, title, content, nullableString(sourceCitation), wordCount)

	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Article, error) {
	row := r.pool.QueryRow(ctx, baseSelect+` WHERE a.id = $1`, id)
	return scanArticle(row)
}

// GetVisible fetches an article and enforces the per-role access rules from
// PRD section 11 at the query layer, not just in middleware.
func (r *Repository) GetVisible(ctx context.Context, id uuid.UUID, userID uuid.UUID, role users.Role) (*Article, error) {
	query := baseSelect + ` WHERE a.id = $1 AND ` + visibilityClause(role, 2)
	row := r.pool.QueryRow(ctx, query, id, userID)
	a, err := scanArticle(row)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			// Distinguish "doesn't exist" from "exists but not visible to you"
			// without leaking existence to unauthorized roles either way.
			if _, rawErr := r.GetByID(ctx, id); rawErr == nil {
				return nil, ErrForbidden
			}
			return nil, ErrNotFound
		}
		return nil, err
	}
	return a, nil
}

func (r *Repository) ListVisible(ctx context.Context, userID uuid.UUID, role users.Role) ([]*Article, error) {
	query := baseSelect + ` WHERE ` + visibilityClause(role, 1) + ` ORDER BY a.updated_at DESC`
	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Article
	for rows.Next() {
		a, err := scanArticle(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, a)
	}
	return out, rows.Err()
}

// visibilityClause returns a SQL predicate implementing PRD section 11's
// access rules per role, referencing the user id as positional arg $<argIndex>
// since callers bind it at different positions (GetVisible has the article
// id at $1; ListVisible only binds the user id, at $1).
func visibilityClause(role users.Role, argIndex int) string {
	// Every branch references userIDParam at least once, even when the role
	// doesn't actually need it for scoping (super_admin/unknown role), so
	// Postgres can always infer the bound parameter's type.
	userIDParam := fmt.Sprintf("$%d", argIndex)
	switch role {
	case users.RoleSuperAdmin:
		return userIDParam + `::uuid IS NOT NULL`
	case users.RoleContributor:
		return `a.contributor_id = ` + userIDParam
	case users.RoleModerator:
		return `(a.status IN ('submitted', 'resubmitted') OR a.reviewer_id = ` + userIDParam + `)`
	case users.RoleGraphicDesigner:
		return `(a.status = 'editorial_approved' OR a.designer_id = ` + userIDParam + `)`
	case users.RolePublisher:
		return `(a.status = 'banner_uploaded' OR a.publisher_id = ` + userIDParam + `)`
	default:
		return userIDParam + `::uuid IS NULL`
	}
}

type UpdateDraftInput struct {
	Title          string
	Content        string
	SourceCitation string
	WordCount      int
}

func (r *Repository) UpdateDraft(ctx context.Context, id uuid.UUID, in UpdateDraftInput) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE articles
		SET title = $1, content = $2, source_citation = $3, word_count = $4, updated_at = now()
		WHERE id = $5
	`, in.Title, in.Content, nullableString(in.SourceCitation), in.WordCount, id)
	return err
}

// TransitionStatus moves an article to a new status with no other field
// changes.
func (r *Repository) TransitionStatus(ctx context.Context, id uuid.UUID, to Status) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE articles SET status = $1, updated_at = now() WHERE id = $2
	`, to, id)
	return err
}

// SetReviewDecision assigns the reviewing moderator and transitions the
// article in one statement. The reviewer assignment matters beyond this
// one transition: once the article leaves the shared submitted/resubmitted
// queue, reviewer_id is what keeps it visible to the moderator who handled
// it (see visibilityClause), e.g. for their activity log.
func (r *Repository) SetReviewDecision(ctx context.Context, id, reviewerID uuid.UUID, to Status) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE articles SET status = $1, reviewer_id = $2, updated_at = now() WHERE id = $3
	`, to, reviewerID, id)
	return err
}

// AssignDesigner records which designer picked up an editorial_approved
// article, without changing its status - banner work happens before the
// banner_uploaded transition fires.
func (r *Repository) AssignDesigner(ctx context.Context, id, designerID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE articles SET designer_id = $1, updated_at = now() WHERE id = $2
	`, designerID, id)
	return err
}

// SetBannerURL records an uploaded banner without changing article status -
// a designer can re-upload a better banner before marking it ready.
func (r *Repository) SetBannerURL(ctx context.Context, id uuid.UUID, bannerURL string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE articles SET cloudinary_banner_url = $1, updated_at = now() WHERE id = $2
	`, bannerURL, id)
	return err
}

// MarkBannerReady transitions an article into banner_uploaded once the
// designer confirms the banner is final.
func (r *Repository) MarkBannerReady(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE articles SET status = 'banner_uploaded', updated_at = now() WHERE id = $1
	`, id)
	return err
}

// MarkPublished records the publisher, the live Substack URL, and the
// publication timestamp in one statement.
func (r *Repository) MarkPublished(ctx context.Context, id, publisherID uuid.UUID, substackURL string) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE articles
		SET status = 'published', publisher_id = $1, substack_url = $2, published_at = now(), updated_at = now()
		WHERE id = $3
	`, publisherID, substackURL, id)
	return err
}

func (r *Repository) MarkSubmitted(ctx context.Context, id uuid.UUID, to Status) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE articles
		SET status = $1, submitted_at = now(), updated_at = now()
		WHERE id = $2
	`, to, id)
	return err
}

func (r *Repository) IncrementReviewCycle(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE articles SET review_cycle_count = review_cycle_count + 1 WHERE id = $1`, id)
	return err
}

func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM articles WHERE id = $1`, id)
	return err
}

func nullableString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
