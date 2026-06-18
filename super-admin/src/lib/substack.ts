import { useMutation, useQueryClient } from "@tanstack/react-query";
import { apiRequest } from "@team1/shared";

export function useTriggerSubstackSync() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => apiRequest<{ synced: number }>("/sync/substack/sync", { method: "POST" }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["substack"] }),
  });
}
