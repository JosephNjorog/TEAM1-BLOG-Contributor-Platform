package notifications

import (
	"time"

	"github.com/google/uuid"
)

type Type string

const (
	TypeArticleSubmitted   Type = "article_submitted"
	TypeChangesRequested   Type = "changes_requested"
	TypeArticleResubmitted Type = "article_resubmitted"
	TypeArticleApproved    Type = "article_approved"
	TypeBannerReady        Type = "banner_ready"
	TypeReadyToPublish     Type = "ready_to_publish"
	TypeArticlePublished   Type = "article_published"
	TypePaymentInitiated   Type = "payment_initiated"
	TypePaymentConfirmed   Type = "payment_confirmed"
	TypePaymentFailed      Type = "payment_failed"
	TypeUserRegistered     Type = "user_registered"
	TypeEmailBounced       Type = "email_bounced"
)

type Notification struct {
	ID        uuid.UUID  `json:"id"`
	UserID    uuid.UUID  `json:"userId"`
	Type      Type       `json:"type"`
	ArticleID *uuid.UUID `json:"articleId"`
	Message   string     `json:"message"`
	Read      bool       `json:"read"`
	CreatedAt time.Time  `json:"createdAt"`
}
