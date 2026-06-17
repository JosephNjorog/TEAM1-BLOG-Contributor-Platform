import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiRequest, type Article, type Banner } from "@team1/shared";

export function useLatestBanner(articleId: string | undefined) {
  return useQuery({
    queryKey: ["banners", articleId],
    queryFn: () => apiRequest<Banner>(`/banners/${articleId}`),
    enabled: !!articleId,
    retry: false,
  });
}

export function useUploadBanner(articleId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (file: File) => {
      const form = new FormData();
      form.append("file", file);
      return apiRequest<{ id: string; cloudinaryUrl: string; uploadedAt: string; article: Article }>(
        `/banners/${articleId}/upload`,
        { method: "POST", body: form },
      );
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["banners", articleId] });
      queryClient.invalidateQueries({ queryKey: ["articles"] });
    },
  });
}

export function useMarkBannerReady(articleId: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => apiRequest<Article>(`/banners/${articleId}/mark-ready`, { method: "POST" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["banners", articleId] });
      queryClient.invalidateQueries({ queryKey: ["articles"] });
    },
  });
}
