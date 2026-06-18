import { useQuery } from "@tanstack/react-query";
import { apiRequest, type Payment } from "@team1/shared";

export function useMyPayments() {
  return useQuery({
    queryKey: ["payments", "mine"],
    queryFn: () => apiRequest<{ payments: Payment[] }>("/payments/mine").then((r) => r.payments),
  });
}
