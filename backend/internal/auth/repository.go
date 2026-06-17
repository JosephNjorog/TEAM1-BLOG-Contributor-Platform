package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"team1blog/backend/internal/users"
)

var (
	ErrInvitationNotFound = errors.New("invitation not found")
	ErrInvitationExpired  = errors.New("invitation expired")
	ErrInvitationUsed     = errors.New("invitation already used")
	ErrRefreshNotFound    = errors.New("refresh token not found")
	ErrRefreshRevoked     = errors.New("refresh token revoked or expired")
)

type Invitation struct {
	ID        uuid.UUID
	Email     string
	Role      users.Role
	TokenHash string
	ExpiresAt time.Time
	UsedAt    *time.Time
	InvitedBy uuid.UUID
	CreatedAt time.Time
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) CreateInvitation(ctx context.Context, email string, role users.Role, tokenHash string, expiresAt time.Time, invitedBy uuid.UUID) (*Invitation, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO invitations (email, role, token_hash, expires_at, invited_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, email, role, token_hash, expires_at, used_at, invited_by, created_at
	`, email, role, tokenHash, expiresAt, invitedBy)
	return scanInvitation(row)
}

func (r *Repository) GetInvitationByTokenHash(ctx context.Context, tokenHash string) (*Invitation, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, email, role, token_hash, expires_at, used_at, invited_by, created_at
		FROM invitations WHERE token_hash = $1
	`, tokenHash)
	return scanInvitation(row)
}

func (r *Repository) MarkInvitationUsed(ctx context.Context, id uuid.UUID) error {
	_, err := r.pool.Exec(ctx, `UPDATE invitations SET used_at = now() WHERE id = $1`, id)
	return err
}

func scanInvitation(row pgx.Row) (*Invitation, error) {
	var inv Invitation
	if err := row.Scan(&inv.ID, &inv.Email, &inv.Role, &inv.TokenHash, &inv.ExpiresAt, &inv.UsedAt, &inv.InvitedBy, &inv.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrInvitationNotFound
		}
		return nil, err
	}
	return &inv, nil
}

// Refresh tokens

func (r *Repository) CreateRefreshToken(ctx context.Context, userID uuid.UUID, tokenHash string, expiresAt time.Time) error {
	_, err := r.pool.Exec(ctx, `
		INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)
	`, userID, tokenHash, expiresAt)
	return err
}

type refreshTokenRow struct {
	UserID    uuid.UUID
	ExpiresAt time.Time
	RevokedAt *time.Time
}

func (r *Repository) GetRefreshToken(ctx context.Context, tokenHash string) (*refreshTokenRow, error) {
	var row refreshTokenRow
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, expires_at, revoked_at FROM refresh_tokens WHERE token_hash = $1
	`, tokenHash).Scan(&row.UserID, &row.ExpiresAt, &row.RevokedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrRefreshNotFound
		}
		return nil, err
	}
	return &row, nil
}

func (r *Repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	_, err := r.pool.Exec(ctx, `UPDATE refresh_tokens SET revoked_at = now() WHERE token_hash = $1`, tokenHash)
	return err
}
