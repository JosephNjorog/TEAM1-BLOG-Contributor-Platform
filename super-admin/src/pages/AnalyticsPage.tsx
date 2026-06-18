import { Button, Card, Spinner } from "@team1/shared";
import { useAnalytics } from "../lib/admin";
import { downloadCSV } from "../lib/csv";

export function AnalyticsPage() {
  const { data, isLoading } = useAnalytics();

  if (isLoading || !data) {
    return (
      <div className="flex justify-center py-16">
        <Spinner />
      </div>
    );
  }

  const maxPubCount = Math.max(1, ...data.publicationVolume.map((p) => p.count));
  const maxPayAmount = Math.max(1, ...data.paymentVolume.map((p) => p.amount));

  return (
    <div className="mx-auto max-w-5xl space-y-10">
      <div>
        <h1 className="mb-1 text-xl font-semibold text-zinc-100">Analytics &amp; Reporting</h1>
        <p className="text-sm text-zinc-500">Average {data.avgPipelineDays.toFixed(1)} days from submission to publication, platform-wide.</p>
      </div>

      <div>
        <div className="mb-3 flex items-center justify-between">
          <h2 className="text-lg font-semibold text-zinc-100">Per-Contributor Metrics</h2>
          <Button
            variant="secondary"
            className="px-3 py-1.5 text-xs"
            onClick={() =>
              downloadCSV(
                "contributor-metrics.csv",
                data.contributorMetrics.map((m) => ({
                  contributor: m.contributorName,
                  articlesSubmitted: m.articlesSubmitted,
                  articlesPublished: m.articlesPublished,
                  acceptanceRate: (m.acceptanceRate * 100).toFixed(1) + "%",
                  avgReviewCycles: m.avgReviewCycles.toFixed(2),
                  avgDaysToPublish: m.avgDaysToPublish.toFixed(2),
                })),
              )
            }
          >
            Export CSV
          </Button>
        </div>
        <div className="overflow-hidden rounded-xl2 border border-surface-border bg-surface-card">
          <table className="w-full text-sm">
            <thead>
              <tr className="border-b border-surface-border text-left text-xs uppercase tracking-wide text-zinc-500">
                <th className="px-4 py-3">Contributor</th>
                <th className="px-4 py-3">Submitted</th>
                <th className="px-4 py-3">Published</th>
                <th className="px-4 py-3">Acceptance</th>
                <th className="px-4 py-3">Avg Cycles</th>
                <th className="px-4 py-3">Avg Days to Publish</th>
              </tr>
            </thead>
            <tbody>
              {data.contributorMetrics.map((m) => (
                <tr key={m.contributorId} className="border-b border-surface-border/60 last:border-b-0">
                  <td className="px-4 py-3 text-zinc-200">{m.contributorName}</td>
                  <td className="px-4 py-3 text-zinc-300">{m.articlesSubmitted}</td>
                  <td className="px-4 py-3 text-zinc-300">{m.articlesPublished}</td>
                  <td className="px-4 py-3 text-zinc-300">{(m.acceptanceRate * 100).toFixed(0)}%</td>
                  <td className="px-4 py-3 text-zinc-300">{m.avgReviewCycles.toFixed(1)}</td>
                  <td className="px-4 py-3 text-zinc-300">{m.avgDaysToPublish.toFixed(1)}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="grid grid-cols-1 gap-6 lg:grid-cols-2">
        <Card>
          <div className="mb-3 flex items-center justify-between">
            <p className="text-xs font-medium uppercase tracking-wide text-zinc-500">Publication volume (weekly)</p>
            <Button
              variant="secondary"
              className="px-2 py-1 text-xs"
              onClick={() => downloadCSV("publication-volume.csv", data.publicationVolume.map((p) => ({ week: p.period, count: p.count })))}
            >
              CSV
            </Button>
          </div>
          <div className="space-y-2">
            {data.publicationVolume.map((p) => (
              <div key={p.period} className="flex items-center gap-3 text-xs">
                <span className="w-20 shrink-0 text-zinc-500">{p.period}</span>
                <div className="h-2 flex-1 overflow-hidden rounded-full bg-surface-base">
                  <div className="h-full rounded-full bg-brand-red" style={{ width: `${(p.count / maxPubCount) * 100}%` }} />
                </div>
                <span className="w-6 text-right text-zinc-300">{p.count}</span>
              </div>
            ))}
          </div>
        </Card>

        <Card>
          <div className="mb-3 flex items-center justify-between">
            <p className="text-xs font-medium uppercase tracking-wide text-zinc-500">Payment volume (monthly)</p>
            <Button
              variant="secondary"
              className="px-2 py-1 text-xs"
              onClick={() =>
                downloadCSV(
                  "payment-volume.csv",
                  data.paymentVolume.map((p) => ({ month: p.period, count: p.count, amountUsd: p.amount.toFixed(2) })),
                )
              }
            >
              CSV
            </Button>
          </div>
          <div className="space-y-2">
            {data.paymentVolume.map((p) => (
              <div key={p.period} className="flex items-center gap-3 text-xs">
                <span className="w-16 shrink-0 text-zinc-500">{p.period}</span>
                <div className="h-2 flex-1 overflow-hidden rounded-full bg-surface-base">
                  <div className="h-full rounded-full bg-emerald-500" style={{ width: `${(p.amount / maxPayAmount) * 100}%` }} />
                </div>
                <span className="w-14 text-right text-zinc-300">${p.amount.toFixed(0)}</span>
              </div>
            ))}
          </div>
        </Card>
      </div>
    </div>
  );
}
