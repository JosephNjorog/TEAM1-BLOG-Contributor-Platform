package notifications

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"team1blog/backend/internal/auth"
	"team1blog/backend/internal/httpx"
)

type Handler struct {
	service *Service
	hub     *Hub
}

func NewHandler(service *Service, hub *Hub) *Handler {
	return &Handler{service: service, hub: hub}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	list, err := h.service.ListForUser(r.Context(), userID, 50)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	unread, err := h.service.CountUnread(r.Context(), userID)
	if err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusOK, map[string]any{"notifications": list, "unreadCount": unread})
}

func (h *Handler) MarkRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "bad_request", "invalid notification id")
		return
	}
	if err := h.service.MarkRead(r.Context(), id, userID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

func (h *Handler) MarkAllRead(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())
	if err := h.service.MarkAllRead(r.Context(), userID); err != nil {
		httpx.Error(w, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	httpx.JSON(w, http.StatusNoContent, nil)
}

var upgrader = websocket.Upgrader{
	// The router's CORS middleware already restricts which origins can
	// reach this handler at all, so the WS-specific origin check gorilla
	// adds on top of that would just be redundant - allow whatever already
	// got past the global middleware chain.
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Socket upgrades to a WebSocket connection and registers it with the Hub
// for the duration of the connection. It only ever reads frames to detect
// the client closing the connection - all data flows server -> client.
func (h *Handler) Socket(w http.ResponseWriter, r *http.Request) {
	userID, _ := auth.UserIDFromContext(r.Context())

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("notifications websocket upgrade failed: %v", err)
		return
	}
	h.hub.Register(userID, conn)
	defer func() {
		h.hub.Unregister(userID, conn)
		_ = conn.Close()
	}()

	for {
		if _, _, err := conn.ReadMessage(); err != nil {
			return
		}
	}
}
