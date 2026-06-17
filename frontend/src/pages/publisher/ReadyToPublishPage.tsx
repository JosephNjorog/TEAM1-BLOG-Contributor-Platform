import { useNavigate } from "react-router-dom";
import { Card, EmptyState, Spinner } from "@team1/shared";
import { useArticles } from "../../lib/articles";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

export function ReadyToPublishPage() {
  const { data: articles, isLoading } = useArticles();
  const navigate = useNavigate();

  const queue = (articles ?? []).filter((a) => a.status === "banner_uploaded");
  const log = (articles ?? [])
    .filter((a) => a.publishedAt)
    .sort((a, b) => new Date(b.publishedAt!).getTime() - new Date(a.publishedAt!).getTime());

  return (
    <div className="mx-auto max-w-4xl space-y-10">
      <div>
        <h1 className="mb-1 text-xl font-semibold text-zinc-100">Ready to Publish</h1>
        <p className="mb-6 text-sm text-zinc-500">Editorially approved, bannered, and waiting to go live.</p>

        {isLoading ? (
          <div className="flex justify-center py-12">
            <Spinner />
          </div>
        ) : queue.length === 0 ? (
          <EmptyState title="Nothing ready to publish" hint="Articles with a finished banner will show up here." />
        ) : (
          <div className="space-y-3">
            {queue.map((a) => (
              <Card
                key={a.id}
                className="cursor-pointer transition-colors hover:border-zinc-700"
                onClick={() => navigate(`/app/publish/${a.id}`)}
              >
                <div className="flex items-center justify-between gap-4">
                  <div className="min-w-0">
                    <p className="truncate font-medium text-zinc-100">{a.title}</p>
                    <p className="mt-1 text-xs text-zinc-500">
                      by {a.contributorName} &middot; reviewed by {a.reviewerName ?? "—"} &middot; banner added {formatDate(a.updatedAt)}
                    </p>
                  </div>
                  {a.cloudinaryBannerUrl && (
                    <img src={a.cloudinaryBannerUrl} alt="" className="h-12 w-20 shrink-0 rounded-lg object-cover" />
                  )}
                </div>
              </Card>
            ))}
          </div>
        )}
      </div>

      <div>
        <h2 className="mb-1 text-lg font-semibold text-zinc-100">Published Log</h2>
        <p className="mb-4 text-sm text-zinc-500">Articles you've posted.</p>

        {log.length === 0 ? (
          <p className="text-sm text-zinc-600">Nothing published yet.</p>
        ) : (
          <div className="space-y-2">
            {log.map((a) => (
              <div key={a.id} className="flex items-center justify-between gap-4 rounded-lg border border-surface-border bg-surface-card px-4 py-3">
                <div className="min-w-0">
                  <p className="truncate text-sm text-zinc-200">{a.title}</p>
                  <p className="text-xs text-zinc-500">{formatDate(a.publishedAt!)}</p>
                </div>
                {a.substackUrl && (
                  <a
                    href={a.substackUrl}
                    target="_blank"
                    rel="noreferrer"
                    onClick={(e) => e.stopPropagation()}
                    className="shrink-0 text-xs font-medium text-brand-red hover:underline"
                  >
                    View on Substack ↗
                  </a>
                )}
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  );
}
