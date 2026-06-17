import { EmptyState, Spinner } from "@team1/shared";
import { useMyPayments } from "../../lib/payments";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

const statusLabels: Record<string, string> = {
  pending: "Pending",
  initiated: "Sent — awaiting confirmation",
  simulated: "Sent — awaiting confirmation",
  confirmed: "Confirmed",
  failed: "Failed",
};

const statusColors: Record<string, string> = {
  pending: "text-zinc-500",
  initiated: "text-teal-400",
  simulated: "text-teal-400",
  confirmed: "text-emerald-400",
  failed: "text-red-400",
};

export function PaymentHistoryPage() {
  const { data: payments, isLoading } = useMyPayments();

  return (
    <div className="mx-auto max-w-3xl">
      <h1 className="mb-1 text-xl font-semibold text-zinc-100">Payments</h1>
      <p className="mb-6 text-sm text-zinc-500">USDC payments for your published articles, $100 each.</p>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Spinner />
        </div>
      ) : !payments || payments.length === 0 ? (
        <EmptyState title="No payments yet" hint="Once an article you wrote is published, payment release will show up here." />
      ) : (
        <div className="space-y-3">
          {payments.map((p) => (
            <div key={p.id} className="rounded-xl2 border border-surface-border bg-surface-card p-4">
              <div className="mb-2 flex items-start justify-between gap-3">
                <p className="font-medium text-zinc-100">{p.articleTitle}</p>
                <span className={`shrink-0 text-xs font-medium ${statusColors[p.status] ?? "text-zinc-400"}`}>
                  {statusLabels[p.status] ?? p.status}
                </span>
              </div>
              <p className="text-sm text-zinc-400">${p.amountUsd.toFixed(2)} to {p.walletAddress}</p>
              <p className="mt-1 text-xs text-zinc-500">
                {p.confirmedAt
                  ? `Confirmed ${formatDate(p.confirmedAt)}`
                  : p.initiatedAt
                    ? `Initiated ${formatDate(p.initiatedAt)}`
                    : `Created ${formatDate(p.createdAt)}`}
              </p>
              {p.txHash && (
                <p className="mt-2 truncate font-mono text-xs text-zinc-600">{p.txHash}</p>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
