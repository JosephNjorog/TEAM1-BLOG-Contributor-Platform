import { useEffect, useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Button, Spinner, StatusPill } from "@team1/shared";
import { useArticle, useSubmitArticle, useUpdateArticle } from "../../lib/articles";
import { useAcceptSuggestion, useArticleFeedback, useRejectSuggestion } from "../../lib/reviews";
import { RichTextEditor } from "../../components/RichTextEditor";
import { SubmitDialog } from "../../components/SubmitDialog";
import { FeedbackPanel } from "../../components/FeedbackPanel";

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
  const { data: feedback } = useArticleFeedback(id);
  const acceptSuggestion = useAcceptSuggestion();
  const rejectSuggestion = useRejectSuggestion();

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [sourceCitation, setSourceCitation] = useState("");
  const [wordCount, setWordCount] = useState(0);
  const [dirty, setDirty] = useState(false);
  const [lastSavedAt, setLastSavedAt] = useState<Date | null>(null);
  const [showSubmitDialog, setShowSubmitDialog] = useState(false);
  const [ready, setReady] = useState(false);

  const hydrated = useRef(false);

  useEffect(() => {
    if (article && !hydrated.current) {
      setTitle(article.title);
      setContent(article.content);
      setSourceCitation(article.sourceCitation ?? "");
      setWordCount(article.wordCount);
      hydrated.current = true;
      // Only mount the rich text editor once its initial content is the
      // real article body - Tiptap's `content` option is read once at
      // mount and never re-applied, so mounting it before hydration would
      // leave it permanently empty even after this state updates.
      setReady(true);
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

  if (isLoading || !article || !ready) {
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
        <div className="mb-8 flex justify-end gap-3">
          <Button variant="secondary" onClick={save} loading={updateArticle.isPending} disabled={!dirty}>
            Save
          </Button>
          <Button onClick={() => setShowSubmitDialog(true)}>Submit for Review</Button>
        </div>
      )}

      {feedback && (feedback.reviewCycles.length > 0 || feedback.suggestions.length > 0) && (
        <div>
          <p className="mb-3 text-xs font-medium uppercase tracking-wide text-zinc-500">Reviewer feedback</p>
          <FeedbackPanel
            reviewCycles={feedback.reviewCycles}
            suggestions={feedback.suggestions}
            onAccept={(suggestionId) => acceptSuggestion.mutate(suggestionId)}
            onReject={(suggestionId) => rejectSuggestion.mutate(suggestionId)}
            busySuggestionId={acceptSuggestion.isPending || rejectSuggestion.isPending ? (acceptSuggestion.variables ?? rejectSuggestion.variables) : undefined}
          />
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
