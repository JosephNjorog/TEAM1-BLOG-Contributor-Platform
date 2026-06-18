package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"team1blog/backend/internal/users"
)

func TestIssueAndParseRoundTrip(t *testing.T) {
	issuer := NewTokenIssuer("test-secret", time.Hour)
	userID := uuid.New()

	token, expiresAt, err := issuer.Issue(userID, users.RoleModerator)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}
	if token == "" {
		t.Fatal("Issue returned an empty token")
	}
	if !expiresAt.After(time.Now()) {
		t.Fatalf("expected expiresAt in the future, got %v", expiresAt)
	}

	claims, err := issuer.Parse(token)
	if err != nil {
		t.Fatalf("Parse returned error for a freshly issued token: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %v, want %v", claims.UserID, userID)
	}
	if claims.Role != users.RoleModerator {
		t.Errorf("Role = %v, want %v", claims.Role, users.RoleModerator)
	}
}

func TestParseRejectsTamperedToken(t *testing.T) {
	issuer := NewTokenIssuer("test-secret", time.Hour)
	token, _, err := issuer.Issue(uuid.New(), users.RoleContributor)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	tampered := token[:len(token)-2] + "xx"
	if _, err := issuer.Parse(tampered); err == nil {
		t.Fatal("expected Parse to reject a tampered token, got nil error")
	}
}

func TestParseRejectsWrongSecret(t *testing.T) {
	issued := NewTokenIssuer("secret-a", time.Hour)
	token, _, err := issued.Issue(uuid.New(), users.RoleSuperAdmin)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	verifier := NewTokenIssuer("secret-b", time.Hour)
	if _, err := verifier.Parse(token); err == nil {
		t.Fatal("expected Parse to reject a token signed with a different secret, got nil error")
	}
}

func TestParseRejectsExpiredToken(t *testing.T) {
	issuer := NewTokenIssuer("test-secret", -time.Minute) // already expired
	token, _, err := issuer.Issue(uuid.New(), users.RoleContributor)
	if err != nil {
		t.Fatalf("Issue returned error: %v", err)
	}

	if _, err := issuer.Parse(token); err == nil {
		t.Fatal("expected Parse to reject an expired token, got nil error")
	}
}
