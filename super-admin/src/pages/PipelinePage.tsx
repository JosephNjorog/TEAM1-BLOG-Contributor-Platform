import { useState } from "react";
import { Button, Spinner, StatusPill, type Article, type ArticleStatus } from "@team1/shared";
import { useArticles } from "../lib/articles";
import { useOverrideArticleStatus } from "../lib/admin";
import { useArticleFeedback } from "../lib/reviews";
import { useLatestBanner } from "../lib/banners";
import { FeedbackPanel } from "../components/FeedbackPanel";

const COLUMNS: { status: ArticleStatus; label: string }[] = [
  { status: "draft", label: "Draft" },
  { status: "submitted", label: "Submitted" },
  { status: "changes_requested", label: "Changes Requested" },
  { status: "resubmitted", label: "Resubmitted" },
  { status: "editorial_approved", label: "Approved" },
  { status: "banner_uploaded", label: "Banner Ready" },
  { status: "published", label: "Published" },
  { status: "payment_initiated", label: "Payment Sent" },
  { status: "payment_confirmed", label: "Paid" },
];

const ALL_STATUSES = COLUMNS.map((c) => c.status);

export function PipelinePage() {
  const { data: articles, isLoading } = useArticles();
  const [selected, setSelected] = useState<Article | null>(null);

  if (isLoading || !articles) {
    return (
      <div className="flex justify-center py-16">
        <Spinner />
      </div>
    );
  }

  return (
    <div>
      <div className="mb-6">
        <h1 className="text-xl font-semibold text-zinc-100">Content Pipeline</h1>
        <p className="text-sm text-zinc-500">Every article on the platform, grouped by stage.</p>
      </div>

      <div className="flex gap-4 overflow-x-auto pb-4">
        {COLUMNS.map((col) => {
          const items = articles.filter((a) => a.status === col.status);
          return (
            <div key={col.status} className="w-64 shrink-0">
              <p className="mb-2 flex items-center justify-between text-xs font-medium uppercase tracking-wide text-zinc-500">
                {col.label}
                <span className="rounded-full bg-surface-raised px-1.5 text-zinc-400">{items.length}</span>
              </p>
              <div className="space-y-2">
                {items.map((a) => (
                  <button
                    key={a.id}
                    onClick={() => setSelected(a)}
                    className="block w-full rounded-xl2 border border-surface-border bg-surface-card p-3 text-left transition-colors hover:border-zinc-700"
                  >
                    <p className="truncate text-sm font-medium text-zinc-100">{a.title}</p>
                    <p className="mt-1 truncate text-xs text-zinc-500">{a.contributorName}</p>
                  </button>
                ))}
              </div>
            </div>
          );
        })}
      </div>

      {selected && <OverridePanel article={selected} onClose={() => setSelected(null)} />}
    </div>
  );
}

function OverridePanel({ article, onClose }: { article: Article; onClose: () => void }) {
  const [status, setStatus] = useState<ArticleStatus>(article.status);
  const [reason, setReason] = useState("");
  const override = useOverrideArticleStatus();
  const { data: feedback } = useArticleFeedback(article.id);
  const { data: banner } = useLatestBanner(article.id);

  const onSubmit = async () => {
    if (!reason.trim()) return;
    await override.mutateAsync({ articleId: article.id, status, reason: reason.trim() });
    onClose();
  };

  return (
    <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/60 px-4 py-8">
      <div className="max-h-full w-full max-w-2xl overflow-y-auto rounded-xl2 border border-surface-border bg-surface-card p-6 shadow-2xl">
        <div className="mb-4 flex items-start justify-between">
          <div>
            <h2 className="text-lg font-semibold text-zinc-100">{article.title}</h2>
            <p className="text-sm text-zinc-500">
              by {article.contributorName} &middot; {article.wordCount} words
            </p>
            <div className="mt-2">
              <StatusPill status={article.status} />
            </div>
          </div>
          <button onClick={onClose} className="text-zinc-500 hover:text-zinc-300">
            ✕
          </button>
        </div>

        <div className="article-editor mb-6 max-h-60 overflow-y-auto rounded-xl2 border border-surface-border bg-surface-base px-4 py-3 text-sm text-zinc-300">
          <div dangerouslySetInnerHTML={{ __html: article.content || "<p><em>No content.</em></p>" }} />
        </div>

        {banner && (
          <div className="mb-6">
            <p className="mb-2 text-xs font-medium uppercase tracking-wide text-zinc-500">Banner</p>
            <img src={banner.cloudinaryUrl} alt="" className="w-full rounded-xl2 object-cover" />
            <p className="mt-1 text-xs text-zinc-500">
              uploaded by {banner.designerName}
              {banner.markedReadyAt ? " · marked ready" : " · not yet marked ready"}
            </p>
          </div>
        )}

        {feedback && (feedback.reviewCycles.length > 0 || feedback.suggestions.length > 0) && (
          <div className="mb-6">
            <p className="mb-2 text-xs font-medium uppercase tracking-wide text-zinc-500">Review history</p>
            <FeedbackPanel reviewCycles={feedback.reviewCycles} suggestions={feedback.suggestions} />
          </div>
        )}

        <div className="rounded-xl2 border border-amber-900/60 bg-amber-950/20 p-4">
          <p className="mb-3 text-xs font-medium uppercase tracking-wide text-amber-400">Manual override</p>
          <label className="mb-1 block text-xs font-medium text-zinc-400">New status</label>
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value as ArticleStatus)}
            className="mb-3 w-full rounded-lg border border-surface-border bg-surface-card px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
          >
            {ALL_STATUSES.map((s) => (
              <option key={s} value={s}>
                {s}
              </option>
            ))}
          </select>
          <label className="mb-1 block text-xs font-medium text-zinc-400">Reason (required)</label>
          <textarea
            value={reason}
            onChange={(e) => setReason(e.target.value)}
            rows={3}
            placeholder="Why this override is necessary..."
            className="mb-3 w-full rounded-lg border border-surface-border bg-surface-card px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
          />
          <div className="flex justify-end">
            <Button onClick={onSubmit} disabled={!reason.trim()} loading={override.isPending}>
              Override status
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}
