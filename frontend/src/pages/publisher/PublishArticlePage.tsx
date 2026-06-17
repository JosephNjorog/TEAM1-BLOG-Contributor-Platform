import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Button, Spinner } from "@team1/shared";
import { useArticle, usePublishArticle } from "../../lib/articles";

export function PublishArticlePage() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { data: article, isLoading } = useArticle(id);
  const publish = usePublishArticle(id);

  const [showDialog, setShowDialog] = useState(false);
  const [substackUrl, setSubstackUrl] = useState("");
  const [error, setError] = useState<string | null>(null);

  if (isLoading || !article) {
    return (
      <div className="flex justify-center py-16">
        <Spinner />
      </div>
    );
  }

  const onConfirm = async () => {
    if (!substackUrl.trim()) {
      setError("Enter the live Substack URL.");
      return;
    }
    try {
      await publish.mutateAsync(substackUrl.trim());
      navigate("/app");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Could not confirm publication.");
    }
  };

  return (
    <div className="mx-auto max-w-3xl">
      <button onClick={() => navigate("/app")} className="mb-4 text-sm text-zinc-500 hover:text-zinc-300">
        ← Back to ready-to-publish
      </button>

      <h1 className="mb-1 text-xl font-semibold text-zinc-100">{article.title}</h1>
      <p className="mb-6 text-sm text-zinc-500">
        by {article.contributorName} &middot; reviewed by {article.reviewerName ?? "—"} &middot; {article.wordCount} words
      </p>

      {article.cloudinaryBannerUrl && (
        <img src={article.cloudinaryBannerUrl} alt="Banner" className="mb-6 w-full rounded-xl2 object-cover" />
      )}

      <div className="article-editor rounded-xl2 border border-surface-border bg-surface-card px-5 py-4 text-zinc-100">
        <div dangerouslySetInnerHTML={{ __html: article.content }} />
      </div>

      <div className="mt-6 flex justify-end">
        <Button onClick={() => setShowDialog(true)}>Confirm Publication</Button>
      </div>

      {showDialog && (
        <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/60 px-4">
          <div className="w-full max-w-md rounded-xl2 border border-surface-border bg-surface-card p-6 shadow-2xl">
            <h2 className="mb-2 text-lg font-semibold text-zinc-100">Confirm publication</h2>
            <p className="mb-4 text-sm text-zinc-400">Paste the live Substack URL. This is saved against the article permanently.</p>
            <input
              autoFocus
              value={substackUrl}
              onChange={(e) => setSubstackUrl(e.target.value)}
              placeholder="https://team1blog.substack.com/p/..."
              className="mb-3 w-full rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
            />
            {error && <p className="mb-3 text-sm text-red-400">{error}</p>}
            <div className="flex justify-end gap-3">
              <Button variant="secondary" onClick={() => setShowDialog(false)} disabled={publish.isPending}>
                Cancel
              </Button>
              <Button onClick={onConfirm} loading={publish.isPending}>
                Confirm &amp; Publish
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
