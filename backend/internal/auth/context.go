package auth

import (
	"context"

	"github.com/google/uuid"
)

type ctxKey int

const claimsCtxKey ctxKey = iota

func SetClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsCtxKey, claims)
}

func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsCtxKey).(*Claims)
	return claims, ok
}

func UserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	claims, ok := ClaimsFromContext(ctx)
	if !ok {
		return uuid.UUID{}, false
	}
	return claims.UserID, true
}
