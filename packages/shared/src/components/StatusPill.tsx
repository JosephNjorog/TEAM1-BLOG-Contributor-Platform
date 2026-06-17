import type { ArticleStatus } from "../types";

const labels: Record<ArticleStatus, string> = {
  draft: "Draft",
  submitted: "In Review",
  changes_requested: "Changes Requested",
  resubmitted: "Resubmitted",
  editorial_approved: "Approved",
  banner_uploaded: "Banner Ready",
  published: "Published",
  payment_initiated: "Payment Sent",
  payment_confirmed: "Paid",
};

const classes: Record<ArticleStatus, string> = {
  draft: "bg-status-draft/15 text-zinc-300 border-status-draft/40",
  submitted: "bg-status-submitted/15 text-amber-300 border-status-submitted/40",
  changes_requested: "bg-status-changes/15 text-orange-300 border-status-changes/40",
  resubmitted: "bg-status-resubmitted/15 text-yellow-300 border-status-resubmitted/40",
  editorial_approved: "bg-status-approved/15 text-sky-300 border-status-approved/40",
  banner_uploaded: "bg-status-banner/15 text-violet-300 border-status-banner/40",
  published: "bg-status-published/15 text-emerald-300 border-status-published/40",
  payment_initiated: "bg-status-payment-initiated/15 text-teal-300 border-status-payment-initiated/40",
  payment_confirmed: "bg-status-payment-confirmed/15 text-emerald-300 border-status-payment-confirmed/40",
};

export function StatusPill({ status }: { status: ArticleStatus }) {
  return (
    <span
      className={`inline-flex items-center rounded-full border px-2.5 py-1 text-xs font-medium ${classes[status]}`}
    >
      {labels[status]}
    </span>
  );
}
