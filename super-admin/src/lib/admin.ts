import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  apiRequest,
  type ContributorSummary,
  type Overview,
  type PendingInvitation,
  type PlatformMetrics,
  type Role,
} from "@team1/shared";

export function useOverview() {
  return useQuery({
    queryKey: ["admin", "overview"],
    queryFn: () => apiRequest<Overview>("/admin/overview"),
    refetchInterval: 30_000,
  });
}

export function useAnalytics() {
  return useQuery({
    queryKey: ["admin", "analytics"],
    queryFn: () => apiRequest<PlatformMetrics>("/admin/analytics"),
  });
}

export function useContributors() {
  return useQuery({
    queryKey: ["admin", "contributors"],
    queryFn: () => apiRequest<{ contributors: ContributorSummary[] }>("/admin/contributors").then((r) => r.contributors),
  });
}

export interface StaffMember {
  id: string;
  name: string;
  email: string;
  role: Role;
  status: "active" | "inactive";
  createdAt: string;
}

export function useStaff() {
  return useQuery({
    queryKey: ["admin", "staff"],
    queryFn: () => apiRequest<{ staff: StaffMember[] }>("/admin/staff").then((r) => r.staff),
  });
}

export function usePendingInvitations() {
  return useQuery({
    queryKey: ["admin", "invitations"],
    queryFn: () => apiRequest<{ invitations: PendingInvitation[] }>("/admin/invitations").then((r) => r.invitations),
  });
}

export function useSendInvite() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (input: { email: string; role: Role }) =>
      apiRequest<void>("/auth/invite", { method: "POST", body: input }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["admin", "invitations"] }),
  });
}

export function useSetUserStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, status }: { userId: string; status: "active" | "inactive" }) =>
      apiRequest<void>(`/admin/users/${userId}/status`, { method: "PATCH", body: { status } }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "contributors"] });
      queryClient.invalidateQueries({ queryKey: ["admin", "staff"] });
    },
  });
}

export function useUpdateUserRole() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ userId, role }: { userId: string; role: Role }) =>
      apiRequest<void>(`/admin/users/${userId}/role`, { method: "PATCH", body: { role } }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["admin", "contributors"] });
      queryClient.invalidateQueries({ queryKey: ["admin", "staff"] });
    },
  });
}

export function useOverrideArticleStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ articleId, status, reason }: { articleId: string; status: string; reason: string }) =>
      apiRequest<void>(`/admin/articles/${articleId}/override`, { method: "POST", body: { status, reason } }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["articles"] });
      queryClient.invalidateQueries({ queryKey: ["admin"] });
    },
  });
}
