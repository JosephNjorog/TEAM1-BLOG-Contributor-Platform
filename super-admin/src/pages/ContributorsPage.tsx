import { useMemo, useState } from "react";
import { Button, Spinner, StatusPill, type ContributorSummary, type Role } from "@team1/shared";
import { useContributors, useSetUserStatus, useUpdateUserRole } from "../lib/admin";
import { useArticles } from "../lib/articles";
import { usePaymentLedger } from "../lib/payments";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

const ROLES: Role[] = ["contributor", "moderator", "graphic_designer", "publisher", "super_admin"];

export function ContributorsPage() {
  const { data: contributors, isLoading } = useContributors();
  const { data: articles } = useArticles();
  const { data: payments } = usePaymentLedger();
  const setStatus = useSetUserStatus();
  const updateRole = useUpdateUserRole();

  const [statusFilter, setStatusFilter] = useState<"all" | "active" | "inactive">("all");
  const [selected, setSelected] = useState<ContributorSummary | null>(null);

  const filtered = (contributors ?? []).filter((c) => statusFilter === "all" || c.status === statusFilter);

  const selectedArticles = useMemo(
    () => (selected ? (articles ?? []).filter((a) => a.contributorId === selected.id) : []),
    [selected, articles],
  );
  const selectedPayments = useMemo(
    () => (selected ? (payments ?? []).filter((p) => p.contributorId === selected.id) : []),
    [selected, payments],
  );

  return (
    <div className="mx-auto max-w-5xl">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-zinc-100">Contributors</h1>
          <p className="text-sm text-zinc-500">{contributors?.length ?? 0} registered.</p>
        </div>
        <div className="flex gap-1 rounded-lg border border-surface-border bg-surface-card p-1 text-xs">
          {(["all", "active", "inactive"] as const).map((f) => (
            <button
              key={f}
              onClick={() => setStatusFilter(f)}
              className={`rounded-md px-3 py-1.5 font-medium capitalize ${
                statusFilter === f ? "bg-surface-raised text-zinc-100" : "text-zinc-500 hover:text-zinc-300"
              }`}
            >
              {f}
            </button>
          ))}
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Spinner />
        </div>
      ) : (
        <div className="overflow-hidden rounded-xl2 border border-surface-border bg-surface-card">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-surface-border text-left text-xs uppercase tracking-wide text-zinc-500">
                <th className="px-4 py-3">Name</th>
                <th className="px-4 py-3">Submitted</th>
                <th className="px-4 py-3">Published</th>
                <th className="px-4 py-3">Paid</th>
                <th className="px-4 py-3">Status</th>
                <th className="px-4 py-3">Joined</th>
                <th className="px-4 py-3" />
              </tr>
            </thead>
            <tbody>
              {filtered.map((c) => (
                <tr key={c.id} className="border-b border-surface-border/60 last:border-b-0 hover:bg-surface-raised/40">
                  <td className="px-4 py-3">
                    <p className="font-medium text-zinc-100">{c.name}</p>
                    <p className="text-xs text-zinc-500">{c.email}</p>
                  </td>
                  <td className="px-4 py-3 text-zinc-300">{c.articlesSubmitted}</td>
                  <td className="px-4 py-3 text-zinc-300">{c.articlesPublished}</td>
                  <td className="px-4 py-3 text-zinc-300">${c.totalPaidUsd.toFixed(2)}</td>
                  <td className="px-4 py-3">
                    <span className={`rounded-full px-2 py-0.5 text-xs font-medium ${c.status === "active" ? "text-emerald-400" : "text-zinc-500"}`}>
                      {c.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-zinc-500">{formatDate(c.registeredAt)}</td>
                  <td className="px-4 py-3 text-right">
                    <button onClick={() => setSelected(c)} className="text-xs font-medium text-brand-red hover:underline">
                      View
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}

      {selected && (
        <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/60 px-4 py-8">
          <div className="max-h-full w-full max-w-2xl overflow-y-auto rounded-xl2 border border-surface-border bg-surface-card p-6 shadow-2xl">
            <div className="mb-4 flex items-start justify-between">
              <div>
                <h2 className="text-lg font-semibold text-zinc-100">{selected.name}</h2>
                <p className="text-sm text-zinc-500">{selected.email}</p>
                {selected.walletAddress && <p className="mt-1 font-mono text-xs text-zinc-600">{selected.walletAddress}</p>}
              </div>
              <button onClick={() => setSelected(null)} className="text-zinc-500 hover:text-zinc-300">
                ✕
              </button>
            </div>

            <div className="mb-6 flex flex-wrap gap-3">
              <Button
                variant="secondary"
                onClick={() => setStatus.mutate({ userId: selected.id, status: selected.status === "active" ? "inactive" : "active" })}
                loading={setStatus.isPending}
              >
                {selected.status === "active" ? "Deactivate" : "Reactivate"}
              </Button>
              <select
                defaultValue="contributor"
                onChange={(e) => updateRole.mutate({ userId: selected.id, role: e.target.value as Role })}
                className="rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
              >
                {ROLES.map((r) => (
                  <option key={r} value={r}>
                    Change role to {r.replace("_", " ")}
                  </option>
                ))}
              </select>
            </div>

            <p className="mb-2 text-xs font-medium uppercase tracking-wide text-zinc-500">Articles ({selectedArticles.length})</p>
            <div className="mb-6 space-y-2">
              {selectedArticles.length === 0 ? (
                <p className="text-sm text-zinc-600">No articles yet.</p>
              ) : (
                selectedArticles.map((a) => (
                  <div key={a.id} className="flex items-center justify-between gap-3 rounded-lg border border-surface-border bg-surface-base px-3 py-2">
                    <span className="truncate text-sm text-zinc-200">{a.title}</span>
                    <StatusPill status={a.status} />
                  </div>
                ))
              )}
            </div>

            <p className="mb-2 text-xs font-medium uppercase tracking-wide text-zinc-500">Payments ({selectedPayments.length})</p>
            <div className="space-y-2">
              {selectedPayments.length === 0 ? (
                <p className="text-sm text-zinc-600">No payments yet.</p>
              ) : (
                selectedPayments.map((p) => (
                  <div key={p.id} className="flex items-center justify-between gap-3 rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm">
                    <span className="truncate text-zinc-200">{p.articleTitle}</span>
                    <span className="text-zinc-400">${p.amountUsd.toFixed(2)}</span>
                    <span className="text-xs text-zinc-500">{p.status}</span>
                  </div>
                ))
              )}
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
