/**
 * Avalanche-derived design tokens shared by the contributor app and the
 * super-admin app. Single source of truth so both products read as one.
 */
export const brand = {
  red: "#E84142",
  redDark: "#C5302F",
  redGlow: "rgba(232, 65, 66, 0.35)",
} as const;

export const surface = {
  app: "#0A0A0B",
  base: "#0F0F11",
  card: "#141416",
  raised: "#1F1F22",
  border: "#27272A",
} as const;

/** Article / payment lifecycle colors — each state must be visually distinct at a glance. */
export const status = {
  draft: "#71717A",
  submitted: "#F59E0B",
  changesRequested: "#FB923C",
  resubmitted: "#FBBF24",
  editorialApproved: "#38BDF8",
  bannerUploaded: "#A78BFA",
  published: "#34D399",
  paymentInitiated: "#2DD4BF",
  paymentConfirmed: "#10B981",
} as const;

export type StatusKey = keyof typeof status;
