import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { apiRequest, type Payment } from "@team1/shared";

export function usePaymentLedger() {
  return useQuery({
    queryKey: ["payments"],
    queryFn: () => apiRequest<{ payments: Payment[] }>("/payments").then((r) => r.payments),
  });
}

export function useMyPayments() {
  return useQuery({
    queryKey: ["payments", "mine"],
    queryFn: () => apiRequest<{ payments: Payment[] }>("/payments/mine").then((r) => r.payments),
  });
}

export function useReleasePayment() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (articleId: string) => apiRequest<Payment>(`/payments/${articleId}/release`, { method: "POST" }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["payments"] });
      queryClient.invalidateQueries({ queryKey: ["articles"] });
      queryClient.invalidateQueries({ queryKey: ["admin"] });
    },
  });
}
