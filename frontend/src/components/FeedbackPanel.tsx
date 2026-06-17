import type { ReviewCycle, Suggestion } from "@team1/shared";

function formatDateTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, { month: "short", day: "numeric", hour: "numeric", minute: "2-digit" });
}

const decisionStyles: Record<string, string> = {
  approved: "text-emerald-400 border-emerald-900 bg-emerald-950/40",
  changes_requested: "text-orange-400 border-orange-900 bg-orange-950/40",
};

const decisionLabels: Record<string, string> = {
  approved: "Approved",
  changes_requested: "Changes requested",
};

interface FeedbackPanelProps {
  reviewCycles: ReviewCycle[];
  suggestions: Suggestion[];
  onAccept?: (suggestionId: string) => void;
  onReject?: (suggestionId: string) => void;
  busySuggestionId?: string;
}

export function FeedbackPanel({ reviewCycles, suggestions, onAccept, onReject, busySuggestionId }: FeedbackPanelProps) {
  if (reviewCycles.length === 0 && suggestions.length === 0) {
    return <p className="text-sm text-zinc-500">No reviewer feedback yet.</p>;
  }

  return (
    <div className="space-y-4">
      {reviewCycles.map((cycle) => (
        <div key={cycle.id} className="rounded-xl2 border border-surface-border bg-surface-card p-4">
          <div className="mb-2 flex items-center justify-between gap-3">
            <span className={`rounded-full border px-2.5 py-0.5 text-xs font-medium ${decisionStyles[cycle.decision]}`}>
              {decisionLabels[cycle.decision]}
            </span>
            <span className="text-xs text-zinc-500">
              {cycle.reviewerName} &middot; {formatDateTime(cycle.createdAt)}
            </span>
          </div>
          {cycle.summary && <p className="text-sm text-zinc-300">{cycle.summary}</p>}
        </div>
      ))}

      {suggestions.length > 0 && (
        <div className="space-y-2.5">
          <p className="text-xs font-medium uppercase tracking-wide text-zinc-500">Inline suggestions</p>
          {suggestions.map((s) => (
            <div key={s.id} className="rounded-xl2 border border-surface-border bg-surface-card p-4">
              <div className="mb-2 flex items-center justify-between gap-3">
                <span className="text-xs text-zinc-500">
                  {s.reviewerName} &middot; {formatDateTime(s.createdAt)}
                </span>
                <span
                  className={`rounded-full px-2 py-0.5 text-xs font-medium ${
                    s.status === "accepted"
                      ? "text-emerald-400"
                      : s.status === "rejected"
                        ? "text-zinc-500"
                        : "text-amber-400"
                  }`}
                >
                  {s.status}
                </span>
              </div>
              <p className="text-sm text-zinc-200">{s.suggestionText}</p>
              {s.status === "pending" && onAccept && onReject && (
                <div className="mt-3 flex gap-2">
                  <button
                    onClick={() => onAccept(s.id)}
                    disabled={busySuggestionId === s.id}
                    className="rounded-lg border border-surface-border px-3 py-1 text-xs font-medium text-emerald-400 hover:bg-surface-raised disabled:opacity-50"
                  >
                    Accept
                  </button>
                  <button
                    onClick={() => onReject(s.id)}
                    disabled={busySuggestionId === s.id}
                    className="rounded-lg border border-surface-border px-3 py-1 text-xs font-medium text-zinc-400 hover:bg-surface-raised disabled:opacity-50"
                  >
                    Reject
                  </button>
                </div>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}
