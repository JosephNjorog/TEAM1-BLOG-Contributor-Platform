import { FormEvent, useState } from "react";
import { Button, ApiClientError, type Role } from "@team1/shared";
import { useStaff, usePendingInvitations, useSendInvite, useSetUserStatus, useUpdateUserRole } from "../../lib/admin";

function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, { month: "short", day: "numeric", hour: "numeric", minute: "2-digit" });
}

const INVITABLE_ROLES: Role[] = ["moderator", "graphic_designer", "publisher", "contributor", "super_admin"];

export function UsersPage() {
  const { data: staff } = useStaff();
  const { data: invitations } = usePendingInvitations();
  const sendInvite = useSendInvite();
  const setStatus = useSetUserStatus();
  const updateRole = useUpdateUserRole();

  const [email, setEmail] = useState("");
  const [role, setRole] = useState<Role>("moderator");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const onInvite = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(null);
    try {
      await sendInvite.mutateAsync({ email, role });
      setSuccess(`Invitation sent to ${email}.`);
      setEmail("");
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : "Could not send invitation.");
    }
  };

  return (
    <div className="mx-auto max-w-4xl space-y-10">
      <div>
        <h1 className="mb-1 text-xl font-semibold text-zinc-100">Users &amp; Invites</h1>
        <p className="mb-6 text-sm text-zinc-500">Invite new team members and manage existing staff accounts.</p>

        <form onSubmit={onInvite} className="mb-8 flex flex-wrap items-end gap-3 rounded-xl2 border border-surface-border bg-surface-card p-5">
          <div className="min-w-[240px] flex-1">
            <label className="mb-1 block text-xs font-medium text-zinc-400">Email</label>
            <input
              type="email"
              required
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              placeholder="name@example.com"
              className="w-full rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
            />
          </div>
          <div>
            <label className="mb-1 block text-xs font-medium text-zinc-400">Role</label>
            <select
              value={role}
              onChange={(e) => setRole(e.target.value as Role)}
              className="rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
            >
              {INVITABLE_ROLES.map((r) => (
                <option key={r} value={r}>
                  {r.replace("_", " ")}
                </option>
              ))}
            </select>
          </div>
          <Button type="submit" loading={sendInvite.isPending}>
            Send Invite
          </Button>
        </form>
        {error && <p className="mb-4 text-sm text-red-400">{error}</p>}
        {success && <p className="mb-4 text-sm text-emerald-400">{success}</p>}
      </div>

      <div>
        <h2 className="mb-3 text-lg font-semibold text-zinc-100">Pending Invitations</h2>
        {!invitations || invitations.length === 0 ? (
          <p className="text-sm text-zinc-600">No invitations sent yet.</p>
        ) : (
          <div className="space-y-2">
            {invitations.map((inv) => {
              const expired = new Date(inv.expiresAt) < new Date();
              return (
                <div key={inv.id} className="flex items-center justify-between gap-3 rounded-lg border border-surface-border bg-surface-card px-4 py-3 text-sm">
                  <div className="min-w-0">
                    <p className="truncate text-zinc-200">{inv.email}</p>
                    <p className="text-xs text-zinc-500">{inv.role.replace("_", " ")}</p>
                  </div>
                  <span className={`text-xs font-medium ${inv.usedAt ? "text-emerald-400" : expired ? "text-red-400" : "text-amber-400"}`}>
                    {inv.usedAt ? `Accepted ${formatDateTime(inv.usedAt)}` : expired ? "Expired" : `Expires ${formatDateTime(inv.expiresAt)}`}
                  </span>
                </div>
              );
            })}
          </div>
        )}
      </div>

      <div>
        <h2 className="mb-3 text-lg font-semibold text-zinc-100">Staff</h2>
        {!staff || staff.length === 0 ? (
          <p className="text-sm text-zinc-600">No staff accounts yet.</p>
        ) : (
          <div className="space-y-2">
            {staff.map((u) => (
              <div key={u.id} className="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-surface-border bg-surface-card px-4 py-3 text-sm">
                <div className="min-w-0">
                  <p className="text-zinc-200">{u.name}</p>
                  <p className="text-xs text-zinc-500">{u.email}</p>
                </div>
                <div className="flex items-center gap-3">
                  <select
                    value={u.role}
                    onChange={(e) => updateRole.mutate({ userId: u.id, role: e.target.value as Role })}
                    className="rounded-lg border border-surface-border bg-surface-base px-2 py-1 text-xs text-zinc-200 outline-none focus:border-brand-red"
                  >
                    {INVITABLE_ROLES.map((r) => (
                      <option key={r} value={r}>
                        {r.replace("_", " ")}
                      </option>
                    ))}
                  </select>
                  <button
                    onClick={() => setStatus.mutate({ userId: u.id, status: u.status === "active" ? "inactive" : "active" })}
                    className={`text-xs font-medium ${u.status === "active" ? "text-zinc-400 hover:text-red-400" : "text-emerald-400 hover:underline"}`}
                  >
                    {u.status === "active" ? "Deactivate" : "Reactivate"}
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
