package admin

import (
	"time"

	"github.com/google/uuid"
	"team1blog/backend/internal/users"
)

// ContributorSummary backs the Super Admin's contributor management table.
type ContributorSummary struct {
	ID                uuid.UUID
	Name              string
	Email             string
	WalletAddress     *string
	Status            users.Status
	RegisteredAt      time.Time
	ArticlesSubmitted int
	ArticlesPublished int
	TotalPaidUSD      float64
	LastSubmissionAt  *time.Time
}

// PendingInvitation backs the invitations list in user management.
type PendingInvitation struct {
	ID        uuid.UUID
	Email     string
	Role      users.Role
	ExpiresAt time.Time
	UsedAt    *time.Time
	CreatedAt time.Time
}

type PipelineCounts struct {
	Draft             int `json:"draft"`
	Submitted         int `json:"submitted"`
	ChangesRequested  int `json:"changesRequested"`
	Resubmitted       int `json:"resubmitted"`
	EditorialApproved int `json:"editorialApproved"`
	BannerUploaded    int `json:"bannerUploaded"`
	Published         int `json:"published"`
	PaymentInitiated  int `json:"paymentInitiated"`
	PaymentConfirmed  int `json:"paymentConfirmed"`
}

type Overview struct {
	TotalPublishedAllTime int            `json:"totalPublishedAllTime"`
	TotalPublished30d     int            `json:"totalPublished30d"`
	TotalPaidUSDAllTime   float64        `json:"totalPaidUsdAllTime"`
	TotalPaidUSD30d       float64        `json:"totalPaidUsd30d"`
	ActiveContributors60d int            `json:"activeContributors60d"`
	PendingPaymentCount   int            `json:"pendingPaymentCount"`
	PendingPaymentUSD     float64        `json:"pendingPaymentUsd"`
	Pipeline              PipelineCounts `json:"pipeline"`
}

// ContributorMetric backs the per-contributor analytics table.
type ContributorMetric struct {
	ContributorID     uuid.UUID `json:"contributorId"`
	ContributorName   string    `json:"contributorName"`
	ArticlesSubmitted int       `json:"articlesSubmitted"`
	ArticlesPublished int       `json:"articlesPublished"`
	AcceptanceRate    float64   `json:"acceptanceRate"` // published / submitted
	AvgReviewCycles   float64   `json:"avgReviewCycles"`
	AvgDaysToPublish  float64   `json:"avgDaysToPublish"`
}

// VolumePoint is one bucket of a time-series metric (e.g. one week/month).
type VolumePoint struct {
	Period string  `json:"period"`
	Count  int     `json:"count"`
	Amount float64 `json:"amount"`
}

type PlatformMetrics struct {
	ContributorMetrics []ContributorMetric `json:"contributorMetrics"`
	PublicationVolume  []VolumePoint       `json:"publicationVolume"`
	PaymentVolume      []VolumePoint       `json:"paymentVolume"`
	AvgPipelineDays    float64             `json:"avgPipelineDays"`
}
