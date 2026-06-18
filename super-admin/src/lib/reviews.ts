import { useQuery } from "@tanstack/react-query";
import { apiRequest, type ReviewCycle, type Suggestion } from "@team1/shared";

interface ArticleFeedback {
  reviewCycles: ReviewCycle[];
  suggestions: Suggestion[];
}

export function useArticleFeedback(articleId: string | undefined) {
  return useQuery({
    queryKey: ["reviews", "article", articleId],
    queryFn: () => apiRequest<ArticleFeedback>(`/reviews/article/${articleId}`),
    enabled: !!articleId,
  });
}
