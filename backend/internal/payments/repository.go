package payments

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("payment not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const baseSelect = `
	SELECT p.id, p.article_id, a.title, p.contributor_id, c.name, p.wallet_address, p.amount_usd,
	       p.tx_hash, p.status, p.initiated_by, p.initiated_at, p.confirmed_at, p.created_at
	FROM payments p
	JOIN articles a ON a.id = p.article_id
	JOIN users c ON c.id = p.contributor_id
`

func scan(row pgx.Row) (*Payment, error) {
	var p Payment
	err := row.Scan(
		&p.ID, &p.ArticleID, &p.ArticleTitle, &p.ContributorID, &p.ContributorName, &p.WalletAddress, &p.AmountUSD,
		&p.TxHash, &p.Status, &p.InitiatedBy, &p.InitiatedAt, &p.ConfirmedAt, &p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &p, nil
}

// GetOrCreate returns the existing payment row for an article, or creates
// a fresh pending one - articles only get a payment row once they're ready
// for release, not eagerly at publish time.
func (r *Repository) GetOrCreate(ctx context.Context, articleID, contributorID uuid.UUID, walletAddress string) (*Payment, error) {
	existing, err := r.GetByArticle(ctx, articleID)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	row := r.pool.QueryRow(ctx, `
		INSERT INTO payments (article_id, contributor_id, wallet_address)
		VALUES ($1, $2, $3)
		RETURNING id
	`, articleID, contributorID, walletAddress)
	var id uuid.UUID
	if err := row.Scan(&id); err != nil {
		return nil, err
	}
	return r.GetByID(ctx, id)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*Payment, error) {
	return scan(r.pool.QueryRow(ctx, baseSelect+` WHERE p.id = $1`, id))
}

func (r *Repository) GetByArticle(ctx context.Context, articleID uuid.UUID) (*Payment, error) {
	return scan(r.pool.QueryRow(ctx, baseSelect+` WHERE p.article_id = $1`, articleID))
}

func (r *Repository) MarkInitiated(ctx context.Context, id uuid.UUID, txHash string, initiatedBy uuid.UUID, status Status) error {
	_, err := r.pool.Exec(ctx, `
		UPDATE payments SET tx_hash = $1, status = $2, initiated_by = $3, initiated_at = now() WHERE id = $4
	`, txHash, status, initiatedBy, id)
	return err
}

func (r *Repository) MarkConfirmed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE payments SET status = 'confirmed', confirmed_at = now() WHERE id = $1`, id)
	return err
}

func (r *Repository) MarkFailed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE payments SET status = 'failed' WHERE id = $1`, id)
	return err
}

func (r *Repository) ListLedger(ctx context.Context) ([]*Payment, error) {
	rows, err := r.pool.Query(ctx, baseSelect+` ORDER BY p.created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Payment
	for rows.Next() {
		p, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *Repository) ListForContributor(ctx context.Context, contributorID uuid.UUID) ([]*Payment, error) {
	rows, err := r.pool.Query(ctx, baseSelect+` WHERE p.contributor_id = $1 ORDER BY p.created_at DESC`, contributorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*Payment
	for rows.Next() {
		p, err := scan(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
