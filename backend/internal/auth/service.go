package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"team1blog/backend/internal/audit"
	"team1blog/backend/internal/email"
	"team1blog/backend/internal/users"
)

var (
	ErrEmailInUse        = errors.New("email already registered")
	ErrInvalidCreds      = errors.New("invalid email or password")
	ErrAccountInactive   = errors.New("account is inactive")
	ErrInvalidWalletAddr = errors.New("invalid Avalanche C-Chain wallet address")
)

var avalancheAddrRE = regexp.MustCompile(`^0x[a-fA-F0-9]{40}$`)

func IsValidAvalancheAddress(addr string) bool {
	return avalancheAddrRE.MatchString(addr)
}

type Service struct {
	users       *users.Repository
	repo        *Repository
	issuer      *TokenIssuer
	audit       *audit.Logger
	mailer      email.Sender
	inviteTTL   time.Duration
	refreshTTL  time.Duration
	appURL      string
	adminAppURL string
	// onRegistered fires after a successful registration. It's a plain
	// callback rather than a direct dependency on internal/notifications,
	// since that package's routes already import internal/auth (for its
	// own RequireAuth middleware) - importing it back here would be a
	// cycle. main.go wires this up to a real notification call.
	onRegistered func(ctx context.Context, u *users.User)
}

func NewService(
	usersRepo *users.Repository,
	repo *Repository,
	issuer *TokenIssuer,
	auditLogger *audit.Logger,
	mailer email.Sender,
	inviteTTL, refreshTTL time.Duration,
	appURL, adminAppURL string,
	onRegistered func(ctx context.Context, u *users.User),
) *Service {
	return &Service{
		users:        usersRepo,
		repo:         repo,
		issuer:       issuer,
		audit:        auditLogger,
		mailer:       mailer,
		inviteTTL:    inviteTTL,
		refreshTTL:   refreshTTL,
		appURL:       appURL,
		adminAppURL:  adminAppURL,
		onRegistered: onRegistered,
	}
}

// registerURLFor picks which app's /register page an invitation should
// point to: Super Admins use the dedicated admin app (a separate
// deployment from the contributor/moderator/designer/publisher frontend),
// everyone else uses the main frontend.
func (s *Service) registerURLFor(role users.Role) string {
	if role == users.RoleSuperAdmin {
		return s.adminAppURL
	}
	return s.appURL
}

type TokenPair struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    time.Time
}

func (s *Service) issueTokenPair(ctx context.Context, u *users.User) (*TokenPair, error) {
	access, expiresAt, err := s.issuer.Issue(u.ID, u.Role)
	if err != nil {
		return nil, err
	}
	rawRefresh, refreshHash, err := generateOpaqueToken()
	if err != nil {
		return nil, err
	}
	if err := s.repo.CreateRefreshToken(ctx, u.ID, refreshHash, time.Now().Add(s.refreshTTL)); err != nil {
		return nil, err
	}
	return &TokenPair{AccessToken: access, RefreshToken: rawRefresh, ExpiresAt: expiresAt}, nil
}

func (s *Service) Invite(ctx context.Context, actorID uuid.UUID, emailAddr string, role users.Role) error {
	if !role.Valid() {
		return fmt.Errorf("invalid role: %s", role)
	}
	rawToken, tokenHash, err := generateOpaqueToken()
	if err != nil {
		return err
	}
	inv, err := s.repo.CreateInvitation(ctx, emailAddr, role, tokenHash, time.Now().Add(s.inviteTTL), actorID)
	if err != nil {
		return err
	}

	registerURL := fmt.Sprintf("%s/register?token=%s", s.registerURLFor(role), rawToken)
	subject, html := email.InvitationEmail(string(role), registerURL)
	if err := s.mailer.Send(ctx, emailAddr, subject, html); err != nil {
		return fmt.Errorf("send invitation email: %w", err)
	}

	return s.audit.Log(ctx, &actorID, "invite_sent", "invitation", &inv.ID, map[string]any{"email": emailAddr, "role": role})
}

type RegisterInput struct {
	Token         string
	Name          string
	Password      string
	Bio           string
	WalletAddress string
}

func (s *Service) RegisterFromInvite(ctx context.Context, in RegisterInput) (*users.User, *TokenPair, error) {
	tokenHash := hashToken(in.Token)
	inv, err := s.repo.GetInvitationByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, nil, err
	}
	if inv.UsedAt != nil {
		return nil, nil, ErrInvitationUsed
	}
	if time.Now().After(inv.ExpiresAt) {
		return nil, nil, ErrInvitationExpired
	}

	if _, err := s.users.GetByEmail(ctx, inv.Email); err == nil {
		return nil, nil, ErrEmailInUse
	} else if !errors.Is(err, users.ErrNotFound) {
		return nil, nil, err
	}

	if inv.Role == users.RoleContributor && in.WalletAddress != "" && !IsValidAvalancheAddress(in.WalletAddress) {
		return nil, nil, ErrInvalidWalletAddr
	}

	passwordHash, err := HashPassword(in.Password)
	if err != nil {
		return nil, nil, err
	}

	invitedBy := inv.InvitedBy
	u, err := s.users.Create(ctx, in.Name, inv.Email, passwordHash, inv.Role, &invitedBy)
	if err != nil {
		return nil, nil, err
	}

	if inv.Role == users.RoleContributor && in.WalletAddress != "" {
		if err := s.users.UpdateWallet(ctx, u.ID, in.WalletAddress); err != nil {
			return nil, nil, err
		}
		u.WalletAddress = &in.WalletAddress
	}
	if in.Bio != "" {
		if err := s.users.UpdateProfile(ctx, u.ID, u.Name, in.Bio); err != nil {
			return nil, nil, err
		}
		u.Bio = &in.Bio
	}

	if err := s.repo.MarkInvitationUsed(ctx, inv.ID); err != nil {
		return nil, nil, err
	}

	tokens, err := s.issueTokenPair(ctx, u)
	if err != nil {
		return nil, nil, err
	}

	_ = s.audit.Log(ctx, &u.ID, "user_registered", "user", &u.ID, map[string]any{"role": u.Role})
	if s.onRegistered != nil {
		s.onRegistered(ctx, u)
	}

	return u, tokens, nil
}

func (s *Service) Login(ctx context.Context, emailAddr, password string) (*users.User, *TokenPair, error) {
	u, err := s.users.GetByEmail(ctx, emailAddr)
	if err != nil {
		if errors.Is(err, users.ErrNotFound) {
			return nil, nil, ErrInvalidCreds
		}
		return nil, nil, err
	}
	if !CheckPassword(u.PasswordHash, password) {
		return nil, nil, ErrInvalidCreds
	}
	if u.Status != users.StatusActive {
		return nil, nil, ErrAccountInactive
	}

	tokens, err := s.issueTokenPair(ctx, u)
	if err != nil {
		return nil, nil, err
	}

	_ = s.audit.Log(ctx, &u.ID, "login", "user", &u.ID, nil)
	return u, tokens, nil
}

func (s *Service) Refresh(ctx context.Context, rawRefreshToken string) (*TokenPair, error) {
	tokenHash := hashToken(rawRefreshToken)
	row, err := s.repo.GetRefreshToken(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if row.RevokedAt != nil || time.Now().After(row.ExpiresAt) {
		return nil, ErrRefreshRevoked
	}

	u, err := s.users.GetByID(ctx, row.UserID)
	if err != nil {
		return nil, err
	}
	if u.Status != users.StatusActive {
		return nil, ErrAccountInactive
	}

	if err := s.repo.RevokeRefreshToken(ctx, tokenHash); err != nil {
		return nil, err
	}

	return s.issueTokenPair(ctx, u)
}

func (s *Service) Logout(ctx context.Context, rawRefreshToken string) error {
	return s.repo.RevokeRefreshToken(ctx, hashToken(rawRefreshToken))
}
