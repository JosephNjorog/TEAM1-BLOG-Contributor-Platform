package articles

import (
	"time"

	"github.com/google/uuid"
)

type Status string

const (
	StatusDraft             Status = "draft"
	StatusSubmitted         Status = "submitted"
	StatusChangesRequested  Status = "changes_requested"
	StatusResubmitted       Status = "resubmitted"
	StatusEditorialApproved Status = "editorial_approved"
	StatusBannerUploaded    Status = "banner_uploaded"
	StatusPublished         Status = "published"
	StatusPaymentInitiated  Status = "payment_initiated"
	StatusPaymentConfirmed  Status = "payment_confirmed"
)

// Editable reports whether a contributor may still change the article body.
func (s Status) Editable() bool {
	return s == StatusDraft || s == StatusChangesRequested
}

type Article struct {
	ID                  uuid.UUID
	ContributorID       uuid.UUID
	ContributorName     string
	ReviewerID          *uuid.UUID
	ReviewerName        *string
	DesignerID          *uuid.UUID
	PublisherID         *uuid.UUID
	Title               string
	Content             string
	SourceCitation      *string
	Status              Status
	WordCount           int
	ReviewCycleCount    int
	SubstackURL         *string
	CloudinaryBannerURL *string
	CreatedAt           time.Time
	UpdatedAt           time.Time
	SubmittedAt         *time.Time
	PublishedAt         *time.Time
}
