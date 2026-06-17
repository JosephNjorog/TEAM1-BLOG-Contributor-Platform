package notifications

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

func (r *Repository) Create(ctx context.Context, userID uuid.UUID, t Type, articleID *uuid.UUID, message string) (*Notification, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO notifications (user_id, type, article_id, message)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, type, article_id, message, read, created_at
	`, userID, t, articleID, message)
	return scan(row)
}

// CreateForRole fans a notification out to every active user with the given
// role - used for moderator/designer/publisher queue events where there's
// no single assigned recipient yet.
func (r *Repository) CreateForRole(ctx context.Context, role string, t Type, articleID *uuid.UUID, message string) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO notifications (user_id, type, article_id, message)
		SELECT id, $1, $2, $3 FROM users WHERE role = $4 AND status = 'active'
	`, t, articleID, message, role)
	return err
}

func (r *Repository) ListForUser(ctx context.Context, userID uuid.UUID, limit int) ([]*Notification, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, user_id, type, article_id, message, read, created_at
		FROM notifications WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2
	`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Notification
	for rows.Next() {
		n, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, n)
	}
	return out, rows.Err()
}

func (r *Repository) CountUnread(ctx context.Context, userID uuid.UUID) (int, error) {
	var count int
	err := r.pool.QueryRow(ctx, `SELECT count(*) FROM notifications WHERE user_id = $1 AND read = false`, userID).Scan(&count)
	return count, err
}

func (r *Repository) MarkRead(ctx context.Context, id, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET read = true WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}

func (r *Repository) MarkAllRead(ctx context.Context, userID uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE notifications SET read = true WHERE user_id = $1 AND read = false`, userID)
	return err
}

func scan(row pgx.Row) (*Notification, error) {
	var n Notification
	if err := row.Scan(&n.ID, &n.UserID, &n.Type, &n.ArticleID, &n.Message, &n.Read, &n.CreatedAt); err != nil {
		return nil, err
	}
	return &n, nil
}
