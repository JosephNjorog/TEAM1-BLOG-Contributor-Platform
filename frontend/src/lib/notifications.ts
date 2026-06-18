import { useCallback } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiRequest, useNotificationSocket, wsURL, type Notification, type WsMessage } from "@team1/shared";

interface NotificationsResponse {
  notifications: Notification[];
  unreadCount: number;
}

export function useNotifications() {
  const queryClient = useQueryClient();

  const onMessage = useCallback(
    (msg: WsMessage) => {
      if (msg.type === "notification") {
        queryClient.invalidateQueries({ queryKey: ["notifications"] });
      }
    },
    [queryClient],
  );
  useNotificationSocket(wsURL("/notifications/ws"), onMessage);

  // Polling stays on as a fallback (e.g. a proxy that doesn't support WS
  // upgrades) - the socket just makes updates feel instant instead of
  // waiting up to this interval.
  return useQuery({
    queryKey: ["notifications"],
    queryFn: () => apiRequest<NotificationsResponse>("/notifications"),
    refetchInterval: 60_000,
  });
}

export function useMarkNotificationRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => apiRequest<void>(`/notifications/${id}/read`, { method: "POST" }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications"] }),
  });
}

export function useMarkAllNotificationsRead() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => apiRequest<void>("/notifications/read-all", { method: "POST" }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["notifications"] }),
  });
}
