import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Button, Spinner, StatusPill } from "@team1/shared";
import { useArticle, useSubmitArticle, useUpdateArticle } from "../../lib/articles";
import { RichTextEditor } from "../../components/RichTextEditor";
import { SubmitDialog } from "../../components/SubmitDialog";

const AUTOSAVE_INTERVAL_MS = 60_000;

function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, { month: "short", day: "numeric", hour: "numeric", minute: "2-digit" });
}

export function ArticleEditorPage() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { data: article, isLoading } = useArticle(id);
  const updateArticle = useUpdateArticle(id);
  const submitArticle = useSubmitArticle(id);

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [sourceCitation, setSourceCitation] = useState("");
  const [wordCount, setWordCount] = useState(0);
  const [dirty, setDirty] = useState(false);
  const [lastSavedAt, setLastSavedAt] = useState<Date | null>(null);
  const [showSubmitDialog, setShowSubmitDialog] = useState(false);

  const hydrated = useRef(false);

  useEffect(() => {
    if (article && !hydrated.current) {
      setTitle(article.title);
      setContent(article.content);
      setSourceCitation(article.sourceCitation ?? "");
      setWordCount(article.wordCount);
      hydrated.current = true;
    }
  }, [article]);

  const editable = article?.status === "draft" || article?.status === "changes_requested";

  const save = async () => {
    if (!dirty) return;
    await updateArticle.mutateAsync({ title, content, sourceCitation });
    setDirty(false);
    setLastSavedAt(new Date());
  };

  useEffect(() => {
    if (!editable) return;
    const interval = setInterval(() => {
      save();
    }, AUTOSAVE_INTERVAL_MS);
    return () => clearInterval(interval);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [editable, dirty, title, content, sourceCitation]);

  if (isLoading || !article) {
    return (
      <div className="flex justify-center py-16">
        <Spinner />
      </div>
    );
  }

  const onConfirmSubmit = async () => {
    await save();
    await submitArticle.mutateAsync();
    setShowSubmitDialog(false);
  };

  return (
    <div className="mx-auto max-w-3xl">
      <button onClick={() => navigate("/app")} className="mb-4 text-sm text-zinc-500 hover:text-zinc-300">
        ← Back to my articles
      </button>

      <div className="mb-4 flex flex-wrap items-center gap-3">
        <StatusPill status={article.status} />
        <span className="text-xs text-zinc-500">{wordCount} words</span>
        {article.reviewerName && <span className="text-xs text-zinc-500">Reviewer: {article.reviewerName}</span>}
        <span className="text-xs text-zinc-500">Updated {formatDateTime(article.updatedAt)}</span>
        {editable && (
          <span className="text-xs text-zinc-600">
            {updateArticle.isPending ? "Saving…" : lastSavedAt ? `Saved ${formatDateTime(lastSavedAt.toISOString())}` : dirty ? "Unsaved changes" : ""}
          </span>
        )}
      </div>

      <input
        value={title}
        disabled={!editable}
        onChange={(e) => {
          setTitle(e.target.value);
          setDirty(true);
        }}
        placeholder="Article title"
        className="mb-4 w-full rounded-xl2 border border-surface-border bg-surface-card px-4 py-3 text-lg font-semibold text-zinc-100 outline-none focus:border-brand-red disabled:opacity-70"
      />

      <RichTextEditor
        content={content}
        editable={!!editable}
        onChange={(html, wc) => {
          setContent(html);
          setWordCount(wc);
          setDirty(true);
        }}
      />

      <label className="mb-1 mt-4 block text-xs font-medium text-zinc-400">External source citation</label>
      <input
        value={sourceCitation}
        disabled={!editable}
        onChange={(e) => {
          setSourceCitation(e.target.value);
          setDirty(true);
        }}
        placeholder="https://..."
        className="mb-6 w-full rounded-lg border border-surface-border bg-surface-card px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red disabled:opacity-70"
      />

      {editable && (
        <div className="flex justify-end gap-3">
          <Button variant="secondary" onClick={save} loading={updateArticle.isPending} disabled={!dirty}>
            Save
          </Button>
          <Button onClick={() => setShowSubmitDialog(true)}>Submit for Review</Button>
        </div>
      )}

      {showSubmitDialog && (
        <SubmitDialog
          title={title}
          wordCount={wordCount}
          loading={submitArticle.isPending || updateArticle.isPending}
          onConfirm={onConfirmSubmit}
          onCancel={() => setShowSubmitDialog(false)}
        />
      )}
    </div>
  );
}
