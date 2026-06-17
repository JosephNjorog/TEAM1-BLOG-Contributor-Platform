import { useEffect, useRef, useState } from "react";
import { getAccessToken } from "./client";

export interface WsMessage<T = unknown> {
  type: string;
  payload: T;
}

/**
 * Connects to the notifications WebSocket and reconnects with backoff on drop.
 * `wsUrl` should be the absolute ws(s):// endpoint, e.g. ws://localhost:8080/api/v1/notifications/ws
 */
export function useNotificationSocket(wsUrl: string, onMessage: (msg: WsMessage) => void) {
  const [connected, setConnected] = useState(false);
  const onMessageRef = useRef(onMessage);
  onMessageRef.current = onMessage;

  useEffect(() => {
    let socket: WebSocket | null = null;
    let retryDelay = 1000;
    let closedByEffect = false;
    let retryTimer: ReturnType<typeof setTimeout> | null = null;

    const connect = () => {
      const token = getAccessToken();
      if (!token) return;
      socket = new WebSocket(`${wsUrl}?token=${encodeURIComponent(token)}`);

      socket.onopen = () => {
        retryDelay = 1000;
        setConnected(true);
      };
      socket.onmessage = (event) => {
        try {
          const msg = JSON.parse(event.data) as WsMessage;
          onMessageRef.current(msg);
        } catch {
          /* ignore malformed frame */
        }
      };
      socket.onclose = () => {
        setConnected(false);
        if (closedByEffect) return;
        retryTimer = setTimeout(connect, retryDelay);
        retryDelay = Math.min(retryDelay * 2, 15000);
      };
      socket.onerror = () => {
        socket?.close();
      };
    };

    connect();

    return () => {
      closedByEffect = true;
      if (retryTimer) clearTimeout(retryTimer);
      socket?.close();
    };
  }, [wsUrl]);

  return { connected };
}
