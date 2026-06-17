import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiRequest, type Article } from "@team1/shared";

export function useArticles() {
  return useQuery({
    queryKey: ["articles"],
    queryFn: () => apiRequest<{ articles: Article[] }>("/articles").then((r) => r.articles),
  });
}

export function useArticle(id: string | undefined) {
  return useQuery({
    queryKey: ["articles", id],
    queryFn: () => apiRequest<Article>(`/articles/${id}`),
    enabled: !!id,
  });
}

export interface ArticleInput {
  title: string;
  content: string;
  sourceCitation: string;
}

export function useCreateArticle() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: Partial<ArticleInput>) => apiRequest<Article>("/articles", { method: "POST", body: input }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["articles"] }),
  });
}

export function useUpdateArticle(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: ArticleInput) => apiRequest<Article>(`/articles/${id}`, { method: "PUT", body: input }),
    onSuccess: (data) => {
      queryClient.setQueryData(["articles", id], data);
      queryClient.invalidateQueries({ queryKey: ["articles"] });
    },
  });
}

export function useSubmitArticle(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: () => apiRequest<Article>(`/articles/${id}/submit`, { method: "POST" }),
    onSuccess: (data) => {
      queryClient.setQueryData(["articles", id], data);
      queryClient.invalidateQueries({ queryKey: ["articles"] });
    },
  });
}

export function useDeleteArticle() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => apiRequest<void>(`/articles/${id}`, { method: "DELETE" }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["articles"] }),
  });
}

export function usePublishArticle(id: string) {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (substackUrl: string) => apiRequest<Article>(`/articles/${id}/publish`, { method: "POST", body: { substackUrl } }),
    onSuccess: (data) => {
      queryClient.setQueryData(["articles", id], data);
      queryClient.invalidateQueries({ queryKey: ["articles"] });
    },
  });
}
