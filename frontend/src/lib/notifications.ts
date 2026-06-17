import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiRequest, type Notification } from "@team1/shared";

interface NotificationsResponse {
  notifications: Notification[];
  unreadCount: number;
}

export function useNotifications() {
  // Polled for now; swapped for the live WebSocket push once the realtime
  // notification hub lands later in the build.
  return useQuery({
    queryKey: ["notifications"],
    queryFn: () => apiRequest<NotificationsResponse>("/notifications"),
    refetchInterval: 15_000,
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
