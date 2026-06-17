import { Button } from "@team1/shared";

interface SubmitDialogProps {
  title: string;
  wordCount: number;
  loading: boolean;
  onConfirm: () => void;
  onCancel: () => void;
}

export function SubmitDialog({ title, wordCount, loading, onConfirm, onCancel }: SubmitDialogProps) {
  return (
    <div className="fixed inset-0 z-30 flex items-center justify-center bg-black/60 px-4">
      <div className="w-full max-w-md rounded-xl2 border border-surface-border bg-surface-card p-6 shadow-2xl">
        <h2 className="mb-2 text-lg font-semibold text-zinc-100">Submit for review?</h2>
        <p className="mb-5 text-sm text-zinc-400">
          Once submitted, the editor locks until a moderator responds. Double-check the title and length below.
        </p>

        <div className="mb-6 rounded-lg border border-surface-border bg-surface-base p-4">
          <p className="text-sm font-medium text-zinc-100">{title || "Untitled draft"}</p>
          <p className="mt-1 text-xs text-zinc-500">{wordCount} words</p>
        </div>

        <div className="flex justify-end gap-3">
          <Button variant="secondary" onClick={onCancel} disabled={loading}>
            Keep editing
          </Button>
          <Button onClick={onConfirm} loading={loading}>
            Submit for review
          </Button>
        </div>
      </div>
    </div>
  );
}
