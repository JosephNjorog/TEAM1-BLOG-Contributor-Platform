package notifications

import (
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// Hub tracks live WebSocket connections per user and pushes new
// notifications to them in real time. The REST endpoints remain the
// source of truth for history and unread counts; this just delivers an
// immediate nudge (with the notification payload itself) so the bell
// updates without waiting on the next poll.
type Hub struct {
	mu    sync.Mutex
	conns map[uuid.UUID]map[*websocket.Conn]struct{}
}

func NewHub() *Hub {
	return &Hub{conns: make(map[uuid.UUID]map[*websocket.Conn]struct{})}
}

func (h *Hub) Register(userID uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.conns[userID] == nil {
		h.conns[userID] = make(map[*websocket.Conn]struct{})
	}
	h.conns[userID][conn] = struct{}{}
}

func (h *Hub) Unregister(userID uuid.UUID, conn *websocket.Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.conns[userID], conn)
	if len(h.conns[userID]) == 0 {
		delete(h.conns, userID)
	}
}

// Push sends a notification to every connection the user currently has
// open (e.g. several browser tabs). Best-effort: a write failure just
// drops that one connection - the client's own reconnect logic in
// useNotificationSocket handles re-establishing it.
func (h *Hub) Push(userID uuid.UUID, n *Notification) {
	h.mu.Lock()
	conns := make([]*websocket.Conn, 0, len(h.conns[userID]))
	for c := range h.conns[userID] {
		conns = append(conns, c)
	}
	h.mu.Unlock()

	for _, c := range conns {
		if err := c.WriteJSON(map[string]any{"type": "notification", "payload": n}); err != nil {
			h.Unregister(userID, c)
			_ = c.Close()
		}
	}
}
