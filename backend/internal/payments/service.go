package payments

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"team1blog/backend/internal/articles"
	"team1blog/backend/internal/audit"
	"team1blog/backend/internal/avalanche"
	"team1blog/backend/internal/email"
	"team1blog/backend/internal/notifications"
	"team1blog/backend/internal/users"
)

const amountUSD = 100.00

var (
	ErrNotPublished    = errors.New("article must be published before payment can be released")
	ErrAlreadyReleased = errors.New("payment has already been released for this article")
	ErrNoWallet        = errors.New("contributor has not connected a wallet")
)

type Service struct {
	repo          *Repository
	articlesRepo  *articles.Repository
	usersRepo     *users.Repository
	notifications *notifications.Repository
	sender        avalanche.Sender
	audit         *audit.Logger
	mailer        email.Sender
	appURL        string
	mock          bool
}

func NewService(
	repo *Repository,
	articlesRepo *articles.Repository,
	usersRepo *users.Repository,
	notificationsRepo *notifications.Repository,
	sender avalanche.Sender,
	auditLogger *audit.Logger,
	mailer email.Sender,
	appURL string,
	mock bool,
) *Service {
	return &Service{
		repo:          repo,
		articlesRepo:  articlesRepo,
		usersRepo:     usersRepo,
		notifications: notificationsRepo,
		sender:        sender,
		audit:         auditLogger,
		mailer:        mailer,
		appURL:        appURL,
		mock:          mock,
	}
}

// Release initiates the USDC transfer for a published article's payment.
// It returns once the transfer has been submitted; confirmation happens
// asynchronously and updates the payment/article state when it lands.
func (s *Service) Release(ctx context.Context, articleID, adminID uuid.UUID) (*Payment, error) {
	a, err := s.articlesRepo.GetByID(ctx, articleID)
	if err != nil {
		return nil, err
	}
	if a.Status != articles.StatusPublished {
		return nil, ErrNotPublished
	}

	contributor, err := s.usersRepo.GetByID(ctx, a.ContributorID)
	if err != nil {
		return nil, err
	}
	if contributor.WalletAddress == nil || *contributor.WalletAddress == "" {
		return nil, ErrNoWallet
	}

	payment, err := s.repo.GetOrCreate(ctx, articleID, a.ContributorID, *contributor.WalletAddress)
	if err != nil {
		return nil, err
	}
	if payment.Status != StatusPending {
		return nil, ErrAlreadyReleased
	}

	txHash, err := s.sender.Send(ctx, payment.WalletAddress, amountUSD)
	if err != nil {
		_ = s.repo.MarkFailed(ctx, payment.ID)
		return nil, fmt.Errorf("onchain transfer failed: %w", err)
	}

	initiatedStatus := StatusInitiated
	if s.mock {
		initiatedStatus = StatusSimulated
	}
	if err := s.repo.MarkInitiated(ctx, payment.ID, txHash, adminID, initiatedStatus); err != nil {
		return nil, err
	}
	if err := s.articlesRepo.TransitionStatus(ctx, articleID, articles.StatusPaymentInitiated); err != nil {
		return nil, err
	}

	dashboardURL := fmt.Sprintf("%s/articles/%s", s.appURL, articleID)
	amountStr := fmt.Sprintf("$%.2f", amountUSD)
	_, _ = s.notifications.Create(ctx, a.ContributorID, notifications.TypePaymentInitiated, &articleID,
		fmt.Sprintf("Payment for %q has been initiated", a.Title))
	subject, html := email.PaymentInitiatedEmail(a.Title, amountStr, dashboardURL)
	_ = s.mailer.Send(ctx, contributor.Email, subject, html)

	_ = s.audit.Log(ctx, &adminID, "payment_initiated", "payment", &payment.ID, map[string]any{
		"articleId": articleID, "txHash": txHash, "amountUsd": amountUSD,
	})

	go s.awaitConfirmation(payment.ID, articleID, txHash)

	return s.repo.GetByID(ctx, payment.ID)
}

// awaitConfirmation runs detached from the originating request - it must
// use its own context since the request's is cancelled once the HTTP
// response is written.
func (s *Service) awaitConfirmation(paymentID, articleID uuid.UUID, txHash string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	confirmed, err := s.sender.Confirm(ctx, txHash)
	if err != nil || !confirmed {
		log.Printf("payment %s (article %s, tx %s) did not confirm: %v", paymentID, articleID, txHash, err)
		_ = s.repo.MarkFailed(ctx, paymentID)
		return
	}

	if err := s.repo.MarkConfirmed(ctx, paymentID); err != nil {
		log.Printf("mark payment %s confirmed: %v", paymentID, err)
		return
	}
	if err := s.articlesRepo.TransitionStatus(ctx, articleID, articles.StatusPaymentConfirmed); err != nil {
		log.Printf("transition article %s to payment_confirmed: %v", articleID, err)
		return
	}

	a, err := s.articlesRepo.GetByID(ctx, articleID)
	if err != nil {
		return
	}
	dashboardURL := fmt.Sprintf("%s/articles/%s", s.appURL, articleID)
	amountStr := fmt.Sprintf("$%.2f", amountUSD)
	_, _ = s.notifications.Create(ctx, a.ContributorID, notifications.TypePaymentConfirmed, &articleID,
		fmt.Sprintf("Payment for %q is confirmed onchain", a.Title))
	if contributor, err := s.usersRepo.GetByID(ctx, a.ContributorID); err == nil {
		subject, html := email.PaymentConfirmedEmail(a.Title, amountStr, txHash, dashboardURL)
		_ = s.mailer.Send(ctx, contributor.Email, subject, html)
	}
}

func (s *Service) ListLedger(ctx context.Context) ([]*Payment, error) {
	return s.repo.ListLedger(ctx)
}

func (s *Service) ListForContributor(ctx context.Context, contributorID uuid.UUID) ([]*Payment, error) {
	return s.repo.ListForContributor(ctx, contributorID)
}

func (s *Service) GetByArticle(ctx context.Context, articleID uuid.UUID) (*Payment, error) {
	return s.repo.GetByArticle(ctx, articleID)
}
