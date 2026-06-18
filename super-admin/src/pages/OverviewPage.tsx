import { Card, Spinner, type PipelineCounts } from "@team1/shared";
import { useOverview } from "../lib/admin";

function MetricCard({ label, value, sub }: { label: string; value: string; sub?: string }) {
  return (
    <Card>
      <p className="text-xs font-medium uppercase tracking-wide text-zinc-500">{label}</p>
      <p className="mt-2 text-2xl font-semibold text-zinc-100">{value}</p>
      {sub && <p className="mt-1 text-xs text-zinc-500">{sub}</p>}
    </Card>
  );
}

const PIPELINE_LABELS: { key: keyof PipelineCounts; label: string; color: string }[] = [
  { key: "draft", label: "Draft", color: "bg-zinc-600" },
  { key: "submitted", label: "Submitted", color: "bg-amber-500" },
  { key: "changesRequested", label: "Changes Requested", color: "bg-orange-500" },
  { key: "resubmitted", label: "Resubmitted", color: "bg-yellow-500" },
  { key: "editorialApproved", label: "Approved", color: "bg-sky-500" },
  { key: "bannerUploaded", label: "Banner Ready", color: "bg-violet-500" },
  { key: "published", label: "Published", color: "bg-emerald-500" },
  { key: "paymentInitiated", label: "Payment Sent", color: "bg-teal-500" },
  { key: "paymentConfirmed", label: "Paid", color: "bg-emerald-600" },
];

export function OverviewPage() {
  const { data: overview, isLoading } = useOverview();

  if (isLoading || !overview) {
    return (
      <div className="flex justify-center py-16">
        <Spinner />
      </div>
    );
  }

  const totalInPipeline = Object.values(overview.pipeline).reduce((a, b) => a + b, 0);

  return (
    <div className="mx-auto max-w-5xl space-y-8">
      <div>
        <h1 className="mb-1 text-xl font-semibold text-zinc-100">Overview</h1>
        <p className="text-sm text-zinc-500">Platform-wide health, at a glance.</p>
      </div>

      <div className="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-3">
        <MetricCard
          label="Published (all time)"
          value={String(overview.totalPublishedAllTime)}
          sub={`${overview.totalPublished30d} in the last 30 days`}
        />
        <MetricCard
          label="Paid out (all time)"
          value={`$${overview.totalPaidUsdAllTime.toFixed(2)}`}
          sub={`$${overview.totalPaidUsd30d.toFixed(2)} in the last 30 days`}
        />
        <MetricCard
          label="Active contributors"
          value={String(overview.activeContributors60d)}
          sub="submitted in the last 60 days"
        />
        <MetricCard
          label="Pending payments"
          value={String(overview.pendingPaymentCount)}
          sub={`$${overview.pendingPaymentUsd.toFixed(2)} outstanding`}
        />
      </div>

      <Card>
        <p className="mb-4 text-xs font-medium uppercase tracking-wide text-zinc-500">Pipeline</p>
        <div className="mb-4 flex h-3 overflow-hidden rounded-full bg-surface-base">
          {PIPELINE_LABELS.map(({ key, color }) => {
            const count = overview.pipeline[key];
            if (!count || totalInPipeline === 0) return null;
            return <div key={key} className={color} style={{ width: `${(count / totalInPipeline) * 100}%` }} />;
          })}
        </div>
        <div className="grid grid-cols-2 gap-3 sm:grid-cols-3">
          {PIPELINE_LABELS.map(({ key, label, color }) => (
            <div key={key} className="flex items-center gap-2 text-sm">
              <span className={`h-2.5 w-2.5 rounded-full ${color}`} />
              <span className="text-zinc-400">{label}</span>
              <span className="ml-auto font-medium text-zinc-200">{overview.pipeline[key]}</span>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}
