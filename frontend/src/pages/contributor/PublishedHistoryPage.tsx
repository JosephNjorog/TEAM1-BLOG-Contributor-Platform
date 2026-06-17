import { EmptyState, Spinner, StatusPill } from "@team1/shared";
import { useArticles } from "../../lib/articles";
import { useMySubstackPosts } from "../../lib/substack";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

const PUBLISHED_STATUSES = new Set(["published", "payment_initiated", "payment_confirmed"]);

export function PublishedHistoryPage() {
  const { data: articles, isLoading: articlesLoading } = useArticles();
  const { data: substackPosts, isLoading: substackLoading } = useMySubstackPosts();

  const published = (articles ?? []).filter((a) => PUBLISHED_STATUSES.has(a.status));
  const isLoading = articlesLoading || substackLoading;

  return (
    <div className="mx-auto max-w-4xl">
      <h1 className="mb-1 text-xl font-semibold text-zinc-100">Published</h1>
      <p className="mb-6 text-sm text-zinc-500">Everything you've published, on the platform and synced from Substack.</p>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Spinner />
        </div>
      ) : published.length === 0 && (substackPosts ?? []).length === 0 ? (
        <EmptyState title="Nothing published yet" hint="Articles you publish through the platform will show up here." />
      ) : (
        <div className="space-y-3">
          {published.map((a) => (
            <div key={a.id} className="flex items-center justify-between gap-4 rounded-xl2 border border-surface-border bg-surface-card p-4">
              <div className="min-w-0">
                <div className="mb-1 flex items-center gap-2">
                  <span className="rounded-full bg-brand-red/15 px-2 py-0.5 text-xs font-medium text-brand-red">Platform</span>
                  <StatusPill status={a.status} />
                </div>
                <p className="truncate font-medium text-zinc-100">{a.title}</p>
                <p className="mt-1 text-xs text-zinc-500">{a.publishedAt && formatDate(a.publishedAt)}</p>
              </div>
              {a.substackUrl && (
                <a href={a.substackUrl} target="_blank" rel="noreferrer" className="shrink-0 text-xs font-medium text-brand-red hover:underline">
                  View ↗
                </a>
              )}
            </div>
          ))}

          {(substackPosts ?? []).map((p) => (
            <div key={p.id} className="flex items-center justify-between gap-4 rounded-xl2 border border-surface-border bg-surface-card p-4">
              <div className="min-w-0">
                <div className="mb-1">
                  <span className="rounded-full bg-zinc-700/40 px-2 py-0.5 text-xs font-medium text-zinc-400">Substack history</span>
                </div>
                <p className="truncate font-medium text-zinc-100">{p.title}</p>
                <p className="mt-1 text-xs text-zinc-500">{formatDate(p.publishedAt)}</p>
              </div>
              <a href={p.url} target="_blank" rel="noreferrer" className="shrink-0 text-xs font-medium text-brand-red hover:underline">
                View ↗
              </a>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
