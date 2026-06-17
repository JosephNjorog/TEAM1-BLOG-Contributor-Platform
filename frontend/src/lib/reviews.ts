import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiRequest, type Article, type ReviewCycle, type Suggestion } from "@team1/shared";

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

export function useReviewActivity() {
  return useQuery({
    queryKey: ["reviews", "activity"],
    queryFn: () => apiRequest<{ activity: ReviewCycle[] }>("/reviews/activity").then((r) => r.activity),
  });
}

export interface SuggestionDraft {
  rangeStart: number;
  rangeEnd: number;
  suggestionText: string;
}

export interface SubmitReviewInput {
  articleId: string;
  decision: "approved" | "changes_requested";
  summary: string;
  suggestions: SuggestionDraft[];
}

export function useSubmitReview() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: SubmitReviewInput) => apiRequest<Article>("/reviews", { method: "POST", body: input }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["articles"] });
      queryClient.invalidateQueries({ queryKey: ["reviews"] });
    },
  });
}

function useSetSuggestionStatus(action: "accept" | "reject") {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (suggestionId: string) => apiRequest<void>(`/reviews/suggestions/${suggestionId}/${action}`, { method: "POST" }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["reviews"] }),
  });
}

export const useAcceptSuggestion = () => useSetSuggestionStatus("accept");
export const useRejectSuggestion = () => useSetSuggestionStatus("reject");
