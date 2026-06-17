package reviews

import (
	"time"

	"github.com/google/uuid"
)

type Decision string

const (
	DecisionApproved         Decision = "approved"
	DecisionChangesRequested Decision = "changes_requested"
)

type SuggestionStatus string

const (
	SuggestionPending  SuggestionStatus = "pending"
	SuggestionAccepted SuggestionStatus = "accepted"
	SuggestionRejected SuggestionStatus = "rejected"
)

type ReviewCycle struct {
	ID           uuid.UUID
	ArticleID    uuid.UUID
	ArticleTitle string
	ReviewerID   uuid.UUID
	ReviewerName string
	Decision     Decision
	Summary      string
	CreatedAt    time.Time
}

type Suggestion struct {
	ID             uuid.UUID
	ArticleID      uuid.UUID
	ReviewCycleID  uuid.UUID
	ReviewerID     uuid.UUID
	ReviewerName   string
	RangeStart     int
	RangeEnd       int
	SuggestionText string
	Status         SuggestionStatus
	CreatedAt      time.Time
}

type SuggestionInput struct {
	RangeStart     int    `json:"rangeStart"`
	RangeEnd       int    `json:"rangeEnd"`
	SuggestionText string `json:"suggestionText"`
}
