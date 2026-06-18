package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"testing"
)

// signFor computes a valid svix-signature header for the given inputs,
// mirroring exactly what Resend/Svix does on their end - used so the test
// doesn't just duplicate verifyResendSignature's own logic to check itself.
func signFor(t *testing.T, secret, svixID, svixTimestamp string, body []byte) string {
	t.Helper()
	key, err := base64.StdEncoding.DecodeString(secret[len("whsec_"):])
	if err != nil {
		t.Fatalf("failed to decode test secret: %v", err)
	}
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(svixID + "." + svixTimestamp + "." + string(body)))
	return "v1," + base64.StdEncoding.EncodeToString(mac.Sum(nil))
}

func TestVerifyResendSignature_Valid(t *testing.T) {
	secret := "whsec_dGVzdHNlY3JldA=="
	svixID, svixTS := "msg_123", "1700000000"
	body := []byte(`{"type":"email.bounced"}`)
	sig := signFor(t, secret, svixID, svixTS, body)

	if !verifyResendSignature(secret, svixID, svixTS, sig, body) {
		t.Error("expected a correctly signed payload to verify")
	}
}

func TestVerifyResendSignature_WrongSecret(t *testing.T) {
	body := []byte(`{"type":"email.bounced"}`)
	sig := signFor(t, "whsec_dGVzdHNlY3JldA==", "msg_123", "1700000000", body)

	if verifyResendSignature("whsec_d3JvbmdzZWNyZXQ=", "msg_123", "1700000000", sig, body) {
		t.Error("expected verification to fail with the wrong secret")
	}
}

func TestVerifyResendSignature_TamperedBody(t *testing.T) {
	secret := "whsec_dGVzdHNlY3JldA=="
	sig := signFor(t, secret, "msg_123", "1700000000", []byte(`{"type":"email.bounced"}`))

	tampered := []byte(`{"type":"email.delivered"}`)
	if verifyResendSignature(secret, "msg_123", "1700000000", sig, tampered) {
		t.Error("expected verification to fail when the body doesn't match what was signed")
	}
}

func TestVerifyResendSignature_MultipleSignaturesInHeader(t *testing.T) {
	secret := "whsec_dGVzdHNlY3JldA=="
	body := []byte(`{"type":"email.complained"}`)
	valid := signFor(t, secret, "msg_123", "1700000000", body)
	header := "v1,bm90dGhlcmlnaHRvbmU= " + valid // a decoy entry plus the real one

	if !verifyResendSignature(secret, "msg_123", "1700000000", header, body) {
		t.Error("expected verification to succeed when the valid signature is one of several in the header")
	}
}
