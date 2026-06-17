package users

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("user not found")

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

const selectColumns = `id, name, email, password_hash, role, wallet_address, bio, status, invited_by, created_at, updated_at`

func scanUser(row pgx.Row) (*User, error) {
	var u User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.Role, &u.WalletAddress, &u.Bio, &u.Status, &u.InvitedBy, &u.CreatedAt, &u.UpdatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) Create(ctx context.Context, name, email, passwordHash string, role Role, invitedBy *uuid.UUID) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO users (name, email, password_hash, role, invited_by)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING `+selectColumns, name, email, passwordHash, role, invitedBy)
	return scanUser(row)
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+selectColumns+` FROM users WHERE email = $1`, email)
	return scanUser(row)
}

func (r *Repository) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	row := r.pool.QueryRow(ctx, `SELECT `+selectColumns+` FROM users WHERE id = $1`, id)
	return scanUser(row)
}

func (r *Repository) UpdateWallet(ctx context.Context, id uuid.UUID, wallet string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET wallet_address = $1, updated_at = now() WHERE id = $2`, wallet, id)
	return err
}

func (r *Repository) UpdateRole(ctx context.Context, id uuid.UUID, role Role) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET role = $1, updated_at = now() WHERE id = $2`, role, id)
	return err
}

func (r *Repository) SetStatus(ctx context.Context, id uuid.UUID, status Status) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET status = $1, updated_at = now() WHERE id = $2`, status, id)
	return err
}

func (r *Repository) UpdateProfile(ctx context.Context, id uuid.UUID, name, bio string) error {
	_, err := r.pool.Exec(ctx, `UPDATE users SET name = $1, bio = $2, updated_at = now() WHERE id = $3`, name, bio, id)
	return err
}

func (r *Repository) List(ctx context.Context) ([]*User, error) {
	rows, err := r.pool.Query(ctx, `SELECT `+selectColumns+` FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*User
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, u)
	}
	return out, rows.Err()
}
