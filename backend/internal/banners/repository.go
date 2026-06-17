package banners

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, articleID, designerID uuid.UUID, cloudinaryURL string) (*Banner, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO banners (article_id, designer_id, cloudinary_url)
		VALUES ($1, $2, $3)
		RETURNING id, article_id, designer_id, cloudinary_url, uploaded_at, marked_ready_at
	`, articleID, designerID, cloudinaryURL)
	return scan(row)
}

func (r *Repository) GetLatestForArticle(ctx context.Context, articleID uuid.UUID) (*Banner, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT b.id, b.article_id, b.designer_id, u.name, b.cloudinary_url, b.uploaded_at, b.marked_ready_at
		FROM banners b
		JOIN users u ON u.id = b.designer_id
		WHERE b.article_id = $1
		ORDER BY b.uploaded_at DESC
		LIMIT 1
	`, articleID)

	var b Banner
	if err := row.Scan(&b.ID, &b.ArticleID, &b.DesignerID, &b.DesignerName, &b.CloudinaryURL, &b.UploadedAt, &b.MarkedReadyAt); err != nil {
		return nil, err
	}
	return &b, nil
}

func (r *Repository) MarkReady(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE banners SET marked_ready_at = now() WHERE id = $1`, id)
	return err
}

func scan(row interface {
	Scan(dest ...any) error
}) (*Banner, error) {
	var b Banner
	if err := row.Scan(&b.ID, &b.ArticleID, &b.DesignerID, &b.CloudinaryURL, &b.UploadedAt, &b.MarkedReadyAt); err != nil {
		return nil, err
	}
	return &b, nil
}
