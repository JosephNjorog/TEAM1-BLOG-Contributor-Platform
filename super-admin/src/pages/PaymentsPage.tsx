import { useState } from "react";
import { Button, EmptyState, Spinner, type Article } from "@team1/shared";
import { useContributors } from "../lib/admin";
import { useArticles } from "../lib/articles";
import { usePaymentLedger, useReleasePayment } from "../lib/payments";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

const statusColors: Record<string, string> = {
  pending: "text-zinc-500",
  initiated: "text-teal-400",
  simulated: "text-teal-400",
  confirmed: "text-emerald-400",
  failed: "text-red-400",
};

export function PaymentsPage() {
  const { data: articles, isLoading: articlesLoading } = useArticles();
  const { data: ledger, isLoading: ledgerLoading } = usePaymentLedger();
  const { data: contributors } = useContributors();
  const releasePayment = useReleasePayment();
  const [confirmArticle, setConfirmArticle] = useState<Article | null>(null);

  const confirmWallet = confirmArticle
    ? contributors?.find((c) => c.id === confirmArticle.contributorId)?.walletAddress
    : null;

  const queue = (articles ?? []).filter((a) => a.status === "published");

  const onConfirmRelease = async () => {
    if (!confirmArticle) return;
    await releasePayment.mutateAsync(confirmArticle.id);
    setConfirmArticle(null);
  };

  return (
    <div className="mx-auto max-w-5xl space-y-10">
      <div>
        <h1 className="mb-1 text-xl font-semibold text-zinc-100">Payments Awaiting Release</h1>
        <p className="mb-6 text-sm text-zinc-500">Published articles ready for their $100 USDC payout.</p>

        {articlesLoading ? (
          <div className="flex justify-center py-12">
            <Spinner />
          </div>
        ) : queue.length === 0 ? (
          <EmptyState title="Nothing awaiting payment" hint="Published articles will show up here." />
        ) : (
          <div className="overflow-hidden rounded-xl2 border border-surface-border bg-surface-card">
            <table className="w-full text-sm">
              <thead>
                <tr className="border-b border-surface-border text-left text-xs uppercase tracking-wide text-zinc-500">
                  <th className="px-4 py-3">Contributor</th>
                  <th className="px-4 py-3">Article</th>
                  <th className="px-4 py-3">Published</th>
                  <th className="px-4 py-3" />
                </tr>
              </thead>
              <tbody>
                {queue.map((a) => (
                  <tr key={a.id} className="border-b border-surface-border/60 last:border-b-0">
                    <td className="px-4 py-3 text-zinc-200">{a.contributorName}</td>
                    <td className="px-4 py-3 text-zinc-300">{a.title}</td>
                    <td className="px-4 py-3 text-zinc-500">{a.publishedAt ? formatDate(a.publishedAt) : "—"}</td>
                    <td className="px-4 py-3 text-right">
                      <Button onClick={() => setConfirmArticle(a)} className="px-3 py-1.5 text-xs">
                        Release Payment
                      </Button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>

      <div>
        <h2 className="mb-1 text-lg font-semibold text-zinc-100">Payment History</h2>
        <p className="mb-4 text-sm text-zinc-500">Full ledger of every payment released.</p>

        {ledgerLoading ? (
          <div className="flex justify-center py-12">
            <Spinner />
          </div>
        ) : !ledger || ledger.length === 0 ? (
          <p className="text-sm text-zinc-600">No payments released yet.</p>
        ) : (
          <div className="space-y-2">
            {ledger.map((p) => (
              <div key={p.id} className="flex flex-wrap items-center justify-between gap-3 rounded-lg border border-surface-border bg-surface-card px-4 py-3">
                <div className="min-w-0">
                  <p className="truncate text-sm text-zinc-200">{p.articleTitle}</p>
                  <p className="text-xs text-zinc-500">{p.contributorName} &middot; ${p.amountUsd.toFixed(2)}</p>
                </div>
                <div className="text-right">
                  <p className={`text-xs font-medium ${statusColors[p.status] ?? "text-zinc-400"}`}>{p.status}</p>
                  {p.txHash && <p className="font-mono text-xs text-zinc-600">{p.txHash.slice(0, 10)}...{p.txHash.slice(-6)}</p>}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>

      {confirmArticle && (
        <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/60 px-4">
          <div className="w-full max-w-md rounded-xl2 border border-surface-border bg-surface-card p-6 shadow-2xl">
            <h2 className="mb-2 text-lg font-semibold text-zinc-100">Release $100.00 USDC?</h2>
            <p className="mb-4 text-sm text-zinc-400">
              This sends an onchain payment to <span className="text-zinc-200">{confirmArticle.contributorName}</span> for &ldquo;{confirmArticle.title}&rdquo;.
            </p>
            <div className="mb-6 rounded-lg border border-surface-border bg-surface-base p-3">
              <p className="text-xs text-zinc-500">Wallet address</p>
              <p className="break-all font-mono text-sm text-zinc-200">{confirmWallet ?? "unknown"}</p>
            </div>
            <div className="flex justify-end gap-3">
              <Button variant="secondary" onClick={() => setConfirmArticle(null)} disabled={releasePayment.isPending}>
                Cancel
              </Button>
              <Button onClick={onConfirmRelease} loading={releasePayment.isPending}>
                Confirm Release
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
