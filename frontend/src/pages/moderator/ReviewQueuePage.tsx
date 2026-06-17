import { useNavigate } from "react-router-dom";
import { Card, EmptyState, Spinner } from "@team1/shared";
import { useArticles } from "../../lib/articles";
import { useReviewActivity } from "../../lib/reviews";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

const decisionLabels: Record<string, string> = {
  approved: "Approved",
  changes_requested: "Changes requested",
};

export function ReviewQueuePage() {
  const { data: articles, isLoading } = useArticles();
  const { data: activity } = useReviewActivity();
  const navigate = useNavigate();

  const queue = (articles ?? [])
    .filter((a) => a.status === "submitted" || a.status === "resubmitted")
    .sort((a, b) => new Date(a.submittedAt ?? a.createdAt).getTime() - new Date(b.submittedAt ?? b.createdAt).getTime());

  return (
    <div className="mx-auto max-w-4xl space-y-10">
      <div>
        <h1 className="mb-1 text-xl font-semibold text-zinc-100">Review Queue</h1>
        <p className="mb-6 text-sm text-zinc-500">Oldest submissions first. Click an article to review it.</p>

        {isLoading ? (
          <div className="flex justify-center py-12">
            <Spinner />
          </div>
        ) : queue.length === 0 ? (
          <EmptyState title="Nothing to review" hint="New submissions will show up here." />
        ) : (
          <div className="space-y-3">
            {queue.map((a) => (
              <Card
                key={a.id}
                className="cursor-pointer transition-colors hover:border-zinc-700"
                onClick={() => navigate(`/app/review/${a.id}`)}
              >
                <div className="flex items-start justify-between gap-4">
                  <div className="min-w-0">
                    <p className="truncate font-medium text-zinc-100">{a.title || "Untitled draft"}</p>
                    <p className="mt-1 text-xs text-zinc-500">
                      by {a.contributorName} &middot; {a.wordCount} words &middot; submitted {formatDate(a.submittedAt ?? a.createdAt)}
                      {a.reviewCycleCount > 0 ? ` · cycle ${a.reviewCycleCount + 1}` : ""}
                    </p>
                  </div>
                  <span
                    className={`shrink-0 rounded-full border px-2.5 py-0.5 text-xs font-medium ${
                      a.status === "resubmitted"
                        ? "border-amber-900 bg-amber-950/40 text-amber-400"
                        : "border-sky-900 bg-sky-950/40 text-sky-400"
                    }`}
                  >
                    {a.status === "resubmitted" ? "Resubmitted" : "New"}
                  </span>
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>

      <div>
        <h2 className="mb-1 text-lg font-semibold text-zinc-100">Activity Log</h2>
        <p className="mb-4 text-sm text-zinc-500">Articles you've reviewed.</p>

        {!activity || activity.length === 0 ? (
          <p className="text-sm text-zinc-600">No reviews yet.</p>
        ) : (
          <div className="space-y-2">
            {activity.map((c) => (
              <div key={c.id} className="flex items-center justify-between gap-4 rounded-lg border border-surface-border bg-surface-card px-4 py-3">
                <div className="min-w-0">
                  <p className="truncate text-sm text-zinc-200">{c.articleTitle}</p>
                  <p className="text-xs text-zinc-500">by {c.contributorName} &middot; {formatDate(c.createdAt)}</p>
                </div>
                <span
                  className={`shrink-0 text-xs font-medium ${c.decision === "approved" ? "text-emerald-400" : "text-orange-400"}`}
                >
                  {decisionLabels[c.decision]}
                </span>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
