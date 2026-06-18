import { useQuery } from "@tanstack/react-query";
import { apiRequest, type Banner } from "@team1/shared";

export function useLatestBanner(articleId: string | undefined) {
  return useQuery({
    queryKey: ["banners", articleId],
    queryFn: () => apiRequest<Banner>(`/banners/${articleId}`),
    enabled: !!articleId,
    retry: false,
  });
}
