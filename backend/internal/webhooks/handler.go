package webhooks

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"team1blog/backend/internal/httpx"
	"team1blog/backend/internal/notifications"
	"team1blog/backend/internal/users"
)

type Handler struct {
	notifications *notifications.Service
	webhookSecret string
}

func NewHandler(notificationsService *notifications.Service, webhookSecret string) *Handler {
	return &Handler{notifications: notificationsService, webhookSecret: webhookSecret}
}

// Resend handles Resend's delivery-event webhook. We only care about
// bounces/complaints here - flagging them is the only thing the PRD asks
// for ("Resend webhook will be used to track delivery status and flag
// bounced emails in the admin panel").
func (h *Handler) Resend(w http.ResponseWriter, r *http.Request) {
	if h.webhookSecret == "" {
		httpx.Error(w, http.StatusServiceUnavailable, "webhook_not_configured", "RESEND_WEBHOOK_SECRET is not set")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "could not read request body")
		return
	}

	ok := verifyResendSignature(
		h.webhookSecret,
		r.Header.Get("svix-id"),
		r.Header.Get("svix-timestamp"),
		r.Header.Get("svix-signature"),
		body,
	)
	if !ok {
		httpx.Error(w, http.StatusUnauthorized, "invalid_signature", "webhook signature verification failed")
		return
	}

	var event resendEvent
	if err := json.Unmarshal(body, &event); err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid webhook payload")
		return
	}

	switch event.Type {
	case "email.bounced", "email.complained":
		recipients := "unknown recipient"
		if len(event.Data.To) > 0 {
			recipients = event.Data.To[0]
		}
		message := fmt.Sprintf("Email %s for %s to %s", strippedEventName(event.Type), event.Data.Subject, recipients)
		_, _ = h.notifications.CreateForRole(r.Context(), string(users.RoleSuperAdmin), notifications.TypeEmailBounced, nil, message)
	}

	httpx.JSON(w, http.StatusOK, map[string]any{"received": true})
}

func strippedEventName(eventType string) string {
	switch eventType {
	case "email.bounced":
		return "bounced"
	case "email.complained":
		return "was marked as spam"
	default:
		return eventType
	}
}
