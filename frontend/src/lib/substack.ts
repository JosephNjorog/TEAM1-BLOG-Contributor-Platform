import { useQuery } from "@tanstack/react-query";
import { apiRequest, type SubstackPost } from "@team1/shared";

export function useMySubstackPosts() {
  return useQuery({
    queryKey: ["substack", "mine"],
    queryFn: () => apiRequest<{ posts: SubstackPost[] }>("/sync/substack/mine").then((r) => r.posts ?? []),
  });
}
