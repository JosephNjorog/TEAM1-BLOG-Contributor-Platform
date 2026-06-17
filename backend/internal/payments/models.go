package payments

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusPending   Status = "pending"
	StatusInitiated Status = "initiated"
	StatusSimulated Status = "simulated"
	StatusConfirmed Status = "confirmed"
	StatusFailed    Status = "failed"
)

type Payment struct {
	ID              uuid.UUID
	ArticleID       uuid.UUID
	ArticleTitle    string
	ContributorID   uuid.UUID
	ContributorName string
	WalletAddress   string
	AmountUSD       float64
	TxHash          *string
	Status          Status
	InitiatedBy     *uuid.UUID
	InitiatedAt     *time.Time
	ConfirmedAt     *time.Time
	CreatedAt       time.Time
}
