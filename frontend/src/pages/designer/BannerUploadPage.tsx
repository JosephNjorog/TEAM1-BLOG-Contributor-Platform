import { useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Button, Spinner } from "@team1/shared";
import { useArticle } from "../../lib/articles";
import { useMarkBannerReady, useUploadBanner } from "../../lib/banners";
import { validateBannerFile } from "../../lib/validateImage";

export function BannerUploadPage() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { data: article, isLoading } = useArticle(id);
  const uploadBanner = useUploadBanner(id);
  const markReady = useMarkBannerReady(id);

  const fileInputRef = useRef<HTMLInputElement>(null);
  const [error, setError] = useState<string | null>(null);

  if (isLoading || !article) {
    return (
      <div className="flex justify-center py-16">
        <Spinner />
      </div>
    );
  }

  const onFileSelected = async (file: File) => {
    setError(null);
    const validationError = await validateBannerFile(file);
    if (validationError) {
      setError(validationError);
      return;
    }
    try {
      await uploadBanner.mutateAsync(file);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Upload failed.");
    }
  };

  const canMarkReady = !!article.cloudinaryBannerUrl;

  return (
    <div className="mx-auto max-w-3xl">
      <button onClick={() => navigate("/app")} className="mb-4 text-sm text-zinc-500 hover:text-zinc-300">
        ← Back to banner queue
      </button>

      <h1 className="mb-1 text-xl font-semibold text-zinc-100">{article.title}</h1>
      <p className="mb-6 text-sm text-zinc-500">by {article.contributorName} &middot; {article.wordCount} words</p>

      <div className="mb-6 rounded-xl2 border border-surface-border bg-surface-card p-5">
        <p className="mb-3 text-xs font-medium uppercase tracking-wide text-zinc-500">Article preview</p>
        <div className="article-editor max-h-72 overflow-y-auto text-sm text-zinc-300" dangerouslySetInnerHTML={{ __html: article.content }} />
      </div>

      <div className="rounded-xl2 border border-surface-border bg-surface-card p-5">
        <p className="mb-3 text-xs font-medium uppercase tracking-wide text-zinc-500">Cover banner</p>

        {article.cloudinaryBannerUrl ? (
          <img src={article.cloudinaryBannerUrl} alt="Current banner" className="mb-4 w-full rounded-xl2 object-cover" />
        ) : (
          <p className="mb-4 text-sm text-zinc-500">No banner uploaded yet.</p>
        )}

        <input
          ref={fileInputRef}
          type="file"
          accept="image/jpeg,image/png"
          className="hidden"
          onChange={(e) => {
            const file = e.target.files?.[0];
            if (file) onFileSelected(file);
            e.target.value = "";
          }}
        />
        <div className="flex flex-wrap items-center gap-3">
          <Button variant="secondary" onClick={() => fileInputRef.current?.click()} loading={uploadBanner.isPending}>
            {article.cloudinaryBannerUrl ? "Re-upload banner" : "Upload banner"}
          </Button>
          <Button onClick={() => markReady.mutateAsync().then(() => navigate("/app"))} disabled={!canMarkReady} loading={markReady.isPending}>
            Mark Banner as Ready
          </Button>
        </div>
        <p className="mt-3 text-xs text-zinc-600">JPG or PNG &middot; minimum 1360&times;1360px &middot; maximum 5MB</p>

        {error && <p className="mt-3 text-sm text-red-400">{error}</p>}
      </div>
    </div>
  );
}
