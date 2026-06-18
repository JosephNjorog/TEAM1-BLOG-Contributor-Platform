import { useQuery } from "@tanstack/react-query";
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
