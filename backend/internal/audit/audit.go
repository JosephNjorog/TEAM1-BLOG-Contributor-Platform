package audit

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Logger struct {
	pool *pgxpool.Pool
}

func NewLogger(pool *pgxpool.Pool) *Logger {
	return &Logger{pool: pool}
}

// Log records a state-changing action. metadata is typically a map with
// "before"/"after" keys; pass nil if there's nothing to capture.
func (l *Logger) Log(ctx context.Context, actorID *uuid.UUID, action, entityType string, entityID *uuid.UUID, metadata map[string]any) error {
	if metadata == nil {
		metadata = map[string]any{}
	}
	raw, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	_, err = l.pool.Exec(ctx, `
		INSERT INTO audit_log (actor_id, action, entity_type, entity_id, metadata)
		VALUES ($1, $2, $3, $4, $5)
	`, actorID, action, entityType, entityID, raw)
	return err
}
