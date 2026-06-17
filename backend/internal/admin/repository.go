package admin

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

func (r *Repository) ListContributors(ctx context.Context) ([]*ContributorSummary, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			u.id, u.name, u.email, u.wallet_address, u.status, u.created_at,
			(SELECT COUNT(*) FROM articles a WHERE a.contributor_id = u.id AND a.submitted_at IS NOT NULL),
			(SELECT COUNT(*) FROM articles a WHERE a.contributor_id = u.id AND a.status IN ('published','payment_initiated','payment_confirmed')),
			(SELECT COALESCE(SUM(p.amount_usd), 0) FROM payments p WHERE p.contributor_id = u.id AND p.status = 'confirmed'),
			(SELECT MAX(a.submitted_at) FROM articles a WHERE a.contributor_id = u.id)
		FROM users u
		WHERE u.role = 'contributor'
		ORDER BY u.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*ContributorSummary
	for rows.Next() {
		var c ContributorSummary
		if err := rows.Scan(&c.ID, &c.Name, &c.Email, &c.WalletAddress, &c.Status, &c.RegisteredAt,
			&c.ArticlesSubmitted, &c.ArticlesPublished, &c.TotalPaidUSD, &c.LastSubmissionAt); err != nil {
			return nil, err
		}
		out = append(out, &c)
	}
	return out, rows.Err()
}

func (r *Repository) ListPendingInvitations(ctx context.Context) ([]*PendingInvitation, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, email, role, expires_at, used_at, created_at
		FROM invitations
		ORDER BY created_at DESC
		LIMIT 200
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*PendingInvitation
	for rows.Next() {
		var inv PendingInvitation
		if err := rows.Scan(&inv.ID, &inv.Email, &inv.Role, &inv.ExpiresAt, &inv.UsedAt, &inv.CreatedAt); err != nil {
			return nil, err
		}
		out = append(out, &inv)
	}
	return out, rows.Err()
}

func (r *Repository) GetOverview(ctx context.Context) (*Overview, error) {
	var o Overview

	if err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM articles WHERE status IN ('published','payment_initiated','payment_confirmed')
	`).Scan(&o.TotalPublishedAllTime); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM articles
		WHERE status IN ('published','payment_initiated','payment_confirmed') AND published_at >= now() - interval '30 days'
	`).Scan(&o.TotalPublished30d); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount_usd), 0) FROM payments WHERE status = 'confirmed'
	`).Scan(&o.TotalPaidUSDAllTime); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(SUM(amount_usd), 0) FROM payments WHERE status = 'confirmed' AND confirmed_at >= now() - interval '30 days'
	`).Scan(&o.TotalPaidUSD30d); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `
		SELECT COUNT(DISTINCT contributor_id) FROM articles WHERE submitted_at >= now() - interval '60 days'
	`).Scan(&o.ActiveContributors60d); err != nil {
		return nil, err
	}
	if err := r.pool.QueryRow(ctx, `
		SELECT COUNT(*), COALESCE(COUNT(*) * 100.00, 0) FROM articles WHERE status = 'published'
	`).Scan(&o.PendingPaymentCount, &o.PendingPaymentUSD); err != nil {
		return nil, err
	}

	rows, err := r.pool.Query(ctx, `SELECT status, COUNT(*) FROM articles GROUP BY status`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		switch status {
		case "draft":
			o.Pipeline.Draft = count
		case "submitted":
			o.Pipeline.Submitted = count
		case "changes_requested":
			o.Pipeline.ChangesRequested = count
		case "resubmitted":
			o.Pipeline.Resubmitted = count
		case "editorial_approved":
			o.Pipeline.EditorialApproved = count
		case "banner_uploaded":
			o.Pipeline.BannerUploaded = count
		case "published":
			o.Pipeline.Published = count
		case "payment_initiated":
			o.Pipeline.PaymentInitiated = count
		case "payment_confirmed":
			o.Pipeline.PaymentConfirmed = count
		}
	}
	return &o, rows.Err()
}

func (r *Repository) GetContributorMetrics(ctx context.Context) ([]ContributorMetric, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT
			u.id, u.name,
			COUNT(*) FILTER (WHERE a.submitted_at IS NOT NULL) AS submitted,
			COUNT(*) FILTER (WHERE a.status IN ('published','payment_initiated','payment_confirmed')) AS published,
			COALESCE(AVG(a.review_cycle_count) FILTER (WHERE a.submitted_at IS NOT NULL), 0) AS avg_cycles,
			COALESCE(AVG(EXTRACT(EPOCH FROM (a.published_at - a.submitted_at)) / 86400) FILTER (WHERE a.published_at IS NOT NULL), 0) AS avg_days
		FROM users u
		LEFT JOIN articles a ON a.contributor_id = u.id
		WHERE u.role = 'contributor'
		GROUP BY u.id, u.name
		ORDER BY u.name
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []ContributorMetric
	for rows.Next() {
		var m ContributorMetric
		var submitted, published int
		if err := rows.Scan(&m.ContributorID, &m.ContributorName, &submitted, &published, &m.AvgReviewCycles, &m.AvgDaysToPublish); err != nil {
			return nil, err
		}
		m.ArticlesSubmitted = submitted
		m.ArticlesPublished = published
		if submitted > 0 {
			m.AcceptanceRate = float64(published) / float64(submitted)
		}
		out = append(out, m)
	}
	return out, rows.Err()
}

func (r *Repository) GetPublicationVolume(ctx context.Context) ([]VolumePoint, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT to_char(date_trunc('week', published_at), 'YYYY-MM-DD'), COUNT(*)
		FROM articles
		WHERE published_at IS NOT NULL
		GROUP BY 1
		ORDER BY 1 DESC
		LIMIT 12
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []VolumePoint
	for rows.Next() {
		var v VolumePoint
		if err := rows.Scan(&v.Period, &v.Count); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *Repository) GetPaymentVolume(ctx context.Context) ([]VolumePoint, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT to_char(date_trunc('month', confirmed_at), 'YYYY-MM'), COUNT(*), COALESCE(SUM(amount_usd), 0)
		FROM payments
		WHERE status = 'confirmed' AND confirmed_at IS NOT NULL
		GROUP BY 1
		ORDER BY 1 DESC
		LIMIT 12
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []VolumePoint
	for rows.Next() {
		var v VolumePoint
		if err := rows.Scan(&v.Period, &v.Count, &v.Amount); err != nil {
			return nil, err
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

func (r *Repository) GetAvgPipelineDays(ctx context.Context) (float64, error) {
	var avg float64
	err := r.pool.QueryRow(ctx, `
		SELECT COALESCE(AVG(EXTRACT(EPOCH FROM (published_at - submitted_at)) / 86400), 0)
		FROM articles WHERE published_at IS NOT NULL AND submitted_at IS NOT NULL
	`).Scan(&avg)
	return avg, err
}

func (r *Repository) OverrideArticleStatus(ctx context.Context, articleID uuid.UUID, status string) error {
	_, err := r.pool.Exec(ctx, `UPDATE articles SET status = $1, updated_at = now() WHERE id = $2`, status, articleID)
	return err
}
