// Package webhooks handles inbound webhook calls from third-party
// providers - currently just Resend's delivery-event webhook, used to
// flag bounced emails for Super Admins.
package webhooks

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strings"
)

// verifyResendSignature checks a Resend (Svix-based) webhook signature.
// secret is the whsec_... value from the Resend dashboard. Algorithm:
// https://docs.svix.com/receiving/verifying-payloads/how-manual-implementation
func verifyResendSignature(secret, svixID, svixTimestamp, svixSignatureHeader string, body []byte) bool {
	secret = strings.TrimPrefix(secret, "whsec_")
	key, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return false
	}

	signedContent := svixID + "." + svixTimestamp + "." + string(body)
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(signedContent))
	expected := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	// The header is a space-separated list of "v1,<base64sig>" entries -
	// match against any of them.
	for _, part := range strings.Fields(svixSignatureHeader) {
		sig := strings.TrimPrefix(part, "v1,")
		if hmac.Equal([]byte(sig), []byte(expected)) {
			return true
		}
	}
	return false
}

type resendEvent struct {
	Type string `json:"type"`
	Data struct {
		To      []string `json:"to"`
		Subject string   `json:"subject"`
	} `json:"data"`
}
