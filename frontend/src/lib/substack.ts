import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiRequest, type SubstackPost } from "@team1/shared";

export function useMySubstackPosts() {
  return useQuery({
    queryKey: ["substack", "mine"],
    queryFn: () => apiRequest<{ posts: SubstackPost[] }>("/sync/substack/mine").then((r) => r.posts ?? []),
  });
}

export function useTriggerSubstackSync() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => apiRequest<{ synced: number }>("/sync/substack/sync", { method: "POST" }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["substack"] }),
  });
}
