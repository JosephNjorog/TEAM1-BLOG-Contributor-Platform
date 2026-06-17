export type Role =
  | "super_admin"
  | "moderator"
  | "graphic_designer"
  | "publisher"
  | "contributor";

export type ArticleStatus =
  | "draft"
  | "submitted"
  | "changes_requested"
  | "resubmitted"
  | "editorial_approved"
  | "banner_uploaded"
  | "published"
  | "payment_initiated"
  | "payment_confirmed";

export type PaymentStatus = "pending" | "initiated" | "simulated" | "confirmed" | "failed";

export interface User {
  id: string;
  name: string;
  email: string;
  role: Role;
  walletAddress: string | null;
  status: "active" | "inactive";
  bio: string | null;
  createdAt: string;
}

export interface Article {
  id: string;
  contributorId: string;
  contributorName: string;
  title: string;
  content: string;
  status: ArticleStatus;
  wordCount: number;
  sourceCitation: string | null;
  substackUrl: string | null;
  cloudinaryBannerUrl: string | null;
  reviewCycleCount: number;
  reviewerName: string | null;
  createdAt: string;
  updatedAt: string;
  submittedAt: string | null;
  publishedAt: string | null;
}

export interface Suggestion {
  id: string;
  reviewCycleId: string;
  reviewerId: string;
  reviewerName: string;
  rangeStart: number;
  rangeEnd: number;
  suggestionText: string;
  status: "pending" | "accepted" | "rejected";
  createdAt: string;
}

export interface ReviewCycle {
  id: string;
  articleTitle?: string;
  contributorName?: string;
  reviewerName: string;
  decision: "approved" | "changes_requested";
  summary: string;
  createdAt: string;
}

export interface Banner {
  id: string;
  articleId: string;
  designerId: string;
  designerName: string;
  cloudinaryUrl: string;
  uploadedAt: string;
  markedReadyAt: string | null;
}

export interface Payment {
  id: string;
  articleId: string;
  articleTitle: string;
  contributorId: string;
  contributorName: string;
  walletAddress: string;
  amountUsd: number;
  txHash: string | null;
  status: PaymentStatus;
  initiatedAt: string | null;
  confirmedAt: string | null;
  createdAt: string;
}

export interface Notification {
  id: string;
  userId: string;
  type: string;
  articleId: string | null;
  message: string;
  read: boolean;
  createdAt: string;
}

export interface SubstackArticle {
  id: string;
  contributorId: string | null;
  substackPostId: string;
  title: string;
  url: string;
  publishedAt: string;
  syncedAt: string;
}

export interface SubstackPost {
  id: string;
  title: string;
  url: string;
  publishedAt: string;
}

export interface ContributorSummary {
  id: string;
  name: string;
  email: string;
  walletAddress: string | null;
  status: "active" | "inactive";
  registeredAt: string;
  articlesSubmitted: number;
  articlesPublished: number;
  totalPaidUsd: number;
  lastSubmissionAt: string | null;
}

export interface PendingInvitation {
  id: string;
  email: string;
  role: Role;
  expiresAt: string;
  usedAt: string | null;
  createdAt: string;
}

export interface PipelineCounts {
  draft: number;
  submitted: number;
  changesRequested: number;
  resubmitted: number;
  editorialApproved: number;
  bannerUploaded: number;
  published: number;
  paymentInitiated: number;
  paymentConfirmed: number;
}

export interface Overview {
  totalPublishedAllTime: number;
  totalPublished30d: number;
  totalPaidUsdAllTime: number;
  totalPaidUsd30d: number;
  activeContributors60d: number;
  pendingPaymentCount: number;
  pendingPaymentUsd: number;
  pipeline: PipelineCounts;
}

export interface ContributorMetric {
  contributorId: string;
  contributorName: string;
  articlesSubmitted: number;
  articlesPublished: number;
  acceptanceRate: number;
  avgReviewCycles: number;
  avgDaysToPublish: number;
}

export interface VolumePoint {
  period: string;
  count: number;
  amount: number;
}

export interface PlatformMetrics {
  contributorMetrics: ContributorMetric[];
  publicationVolume: VolumePoint[];
  paymentVolume: VolumePoint[];
  avgPipelineDays: number;
}

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  expiresAt: string;
}

export interface ApiError {
  error: string;
  message: string;
}
