import { useRef, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { Button, Spinner } from "@team1/shared";
import { useArticle } from "../../lib/articles";
import { useArticleFeedback, useSubmitReview, type SuggestionDraft } from "../../lib/reviews";
import { FeedbackPanel } from "../../components/FeedbackPanel";
import { captureSelection, type SelectionRange } from "../../lib/textSelection";

interface PendingSuggestion extends SuggestionDraft {
  quotedText: string;
}

export function ReviewArticlePage() {
  const { id = "" } = useParams();
  const navigate = useNavigate();
  const { data: article, isLoading } = useArticle(id);
  const { data: feedback } = useArticleFeedback(id);
  const submitReview = useSubmitReview();

  const contentRef = useRef<HTMLDivElement>(null);
  const [pending, setPending] = useState<PendingSuggestion[]>([]);
  const [selection, setSelection] = useState<SelectionRange | null>(null);
  const [commentDraft, setCommentDraft] = useState("");
  const [summary, setSummary] = useState("");

  if (isLoading || !article) {
    return (
      <div className="flex justify-center py-16">
        <Spinner />
      </div>
    );
  }

  const onMouseUp = () => {
    if (!contentRef.current) return;
    const range = captureSelection(contentRef.current);
    setSelection(range);
    setCommentDraft("");
  };

  const addSuggestion = () => {
    if (!selection || commentDraft.trim() === "") return;
    setPending((prev) => [
      ...prev,
      { rangeStart: selection.start, rangeEnd: selection.end, suggestionText: commentDraft.trim(), quotedText: selection.text },
    ]);
    setSelection(null);
    setCommentDraft("");
    window.getSelection()?.removeAllRanges();
  };

  const removePending = (index: number) => setPending((prev) => prev.filter((_, i) => i !== index));

  const decide = async (decision: "approved" | "changes_requested") => {
    await submitReview.mutateAsync({
      articleId: id,
      decision,
      summary,
      suggestions: pending.map(({ rangeStart, rangeEnd, suggestionText }) => ({ rangeStart, rangeEnd, suggestionText })),
    });
    navigate("/app");
  };

  return (
    <div className="mx-auto max-w-3xl">
      <button onClick={() => navigate("/app")} className="mb-4 text-sm text-zinc-500 hover:text-zinc-300">
        ← Back to queue
      </button>

      <div className="mb-4">
        <h1 className="text-xl font-semibold text-zinc-100">{article.title}</h1>
        <p className="mt-1 text-sm text-zinc-500">
          by {article.contributorName} &middot; {article.wordCount} words
          {article.sourceCitation && (
            <>
              {" "}
              &middot; source: <span className="text-zinc-400">{article.sourceCitation}</span>
            </>
          )}
        </p>
      </div>

      <div className="relative">
        <div
          ref={contentRef}
          onMouseUp={onMouseUp}
          className="article-editor select-text rounded-xl2 border border-surface-border bg-surface-card px-5 py-4 text-zinc-100"
          dangerouslySetInnerHTML={{ __html: article.content }}
        />

        {selection && (
          <div
            className="fixed z-30 w-72 rounded-xl border border-surface-border bg-surface-raised p-3 shadow-2xl"
            style={{ top: selection.rect.bottom + window.scrollY + 8, left: Math.min(selection.rect.left + window.scrollX, window.innerWidth - 300) }}
          >
            <p className="mb-2 line-clamp-2 text-xs italic text-zinc-400">&ldquo;{selection.text}&rdquo;</p>
            <textarea
              autoFocus
              value={commentDraft}
              onChange={(e) => setCommentDraft(e.target.value)}
              placeholder="Suggest a change or leave a comment..."
              rows={3}
              className="mb-2 w-full rounded-lg border border-surface-border bg-surface-base px-2.5 py-1.5 text-sm text-zinc-100 outline-none focus:border-brand-red"
            />
            <div className="flex justify-end gap-2">
              <button onClick={() => setSelection(null)} className="rounded-lg px-3 py-1 text-xs text-zinc-500 hover:text-zinc-300">
                Cancel
              </button>
              <Button onClick={addSuggestion} className="px-3 py-1 text-xs">
                Add suggestion
              </Button>
            </div>
          </div>
        )}
      </div>

      {pending.length > 0 && (
        <div className="mt-4 space-y-2">
          <p className="text-xs font-medium uppercase tracking-wide text-zinc-500">New suggestions ({pending.length})</p>
          {pending.map((p, i) => (
            <div key={i} className="flex items-start justify-between gap-3 rounded-lg border border-surface-border bg-surface-card px-4 py-3">
              <div className="min-w-0">
                <p className="truncate text-xs italic text-zinc-500">&ldquo;{p.quotedText}&rdquo;</p>
                <p className="text-sm text-zinc-200">{p.suggestionText}</p>
              </div>
              <button onClick={() => removePending(i)} className="shrink-0 text-xs text-zinc-500 hover:text-red-400">
                Remove
              </button>
            </div>
          ))}
        </div>
      )}

      {feedback && (feedback.reviewCycles.length > 0 || feedback.suggestions.length > 0) && (
        <div className="mt-8">
          <p className="mb-3 text-xs font-medium uppercase tracking-wide text-zinc-500">Previous review history</p>
          <FeedbackPanel reviewCycles={feedback.reviewCycles} suggestions={feedback.suggestions} />
        </div>
      )}

      <div className="mt-8">
        <label className="mb-1 block text-xs font-medium text-zinc-400">Review summary</label>
        <textarea
          value={summary}
          onChange={(e) => setSummary(e.target.value)}
          rows={4}
          placeholder="Overall feedback before submitting your decision..."
          className="w-full rounded-lg border border-surface-border bg-surface-card px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
        />
      </div>

      <div className="mt-5 flex justify-end gap-3">
        <Button variant="secondary" onClick={() => decide("changes_requested")} loading={submitReview.isPending}>
          Request Changes
        </Button>
        <Button onClick={() => decide("approved")} loading={submitReview.isPending}>
          Approve
        </Button>
      </div>
    </div>
  );
}
