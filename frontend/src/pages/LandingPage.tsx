import { useState } from "react";
import { Link } from "react-router-dom";

interface RoleGuide {
  id: string;
  label: string;
  tagline: string;
  bullets: string[];
}

const ROLE_GUIDES: RoleGuide[] = [
  {
    id: "contributor",
    label: "Contributor",
    tagline: "Write articles, track every step, get paid on publish.",
    bullets: [
      "Write in a rich text editor with autosave every 60 seconds, plus a manual save whenever you want.",
      "Click \"Submit for Review\" when it's ready — the editor locks until a moderator responds.",
      "See reviewer feedback as inline suggestions, right where the text is, and resubmit once you've addressed them.",
      "Track every article's status — Draft, In Review, Changes Requested, Approved, Published, Paid — from one list.",
      "Connect your Core wallet once during onboarding. You're paid 100 USDC on Avalanche the same day your article goes live.",
    ],
  },
  {
    id: "moderator",
    label: "Moderator",
    tagline: "The editorial gate. Review, suggest, decide.",
    bullets: [
      "A review queue shows every submitted article, oldest first, with word count and how many revision cycles it's had.",
      "Highlight any text to leave an inline suggestion or comment, the same way you would in a shared doc.",
      "Leave an overall summary, then either request changes or approve — approval hands the article straight to a graphic designer.",
      "Your past decisions are kept in an activity log, so there's a record of what you approved and when.",
    ],
  },
  {
    id: "graphic_designer",
    label: "Graphic Designer",
    tagline: "Turn approved articles into a published-ready banner.",
    bullets: [
      "A banner queue shows every article that's cleared editorial review and is waiting on a cover image.",
      "Read the article first for context, then upload a banner — JPG or PNG, at least 1360×1360px, under 5MB.",
      "Re-upload as many times as you like before you're happy with it.",
      "Click \"Mark Banner as Ready\" to hand the article off to a publisher.",
    ],
  },
  {
    id: "publisher",
    label: "Publisher",
    tagline: "The last step before an article goes live.",
    bullets: [
      "A ready-to-publish queue shows every article with both editorial approval and a finished banner.",
      "Review the full article and banner once more, then confirm publication.",
      "Paste in the live Substack URL — it's saved against the article permanently.",
      "The contributor, the reviewing moderator, and the Super Admin are all notified the moment you confirm.",
    ],
  },
  {
    id: "super_admin",
    label: "Super Admin",
    tagline: "Full visibility, payment authority, and platform control.",
    bullets: [
      "An overview dashboard tracks publication volume, payments released, and where every article sits in the pipeline.",
      "Invite new contributors, moderators, designers, and publishers by email — invitations expire after 72 hours.",
      "Release the 100 USDC payment for any published article with a confirmation step before the onchain transfer fires.",
      "Full read access to every article at every stage, plus the ability to override a stuck article's state with a reason note.",
    ],
  },
];

const LIFECYCLE = [
  "Draft",
  "Submitted",
  "In Review",
  "Approved",
  "Banner Ready",
  "Published",
  "Paid",
];

export function LandingPage() {
  const [openRole, setOpenRole] = useState<string | null>("contributor");

  return (
    <div className="min-h-screen bg-surface-app text-zinc-100">
      <header className="flex items-center justify-between px-6 py-5 sm:px-10">
        <div className="flex items-center gap-2">
          <span className="h-2.5 w-2.5 rounded-full bg-brand-red shadow-glow-red" />
          <span className="text-sm font-bold uppercase tracking-wide">Team1 Blog</span>
        </div>
        <Link
          to="/login"
          className="rounded-xl bg-brand-red px-4 py-2 text-sm font-medium text-white shadow-glow-red transition-colors hover:bg-brand-red-dark"
        >
          Sign in
        </Link>
      </header>

      {/* Hero */}
      <section className="mx-auto max-w-4xl px-6 pt-16 pb-20 text-center sm:px-10">
        <p className="mb-4 text-xs font-semibold uppercase tracking-widest text-brand-red">
          Contributor Platform
        </p>
        <h1 className="mb-6 text-4xl font-bold leading-tight sm:text-5xl">
          One place to write, review, design,
          <br />
          publish, and get paid.
        </h1>
        <p className="mx-auto mb-10 max-w-2xl text-base text-zinc-400 sm:text-lg">
          Team1 Blog used to run across Telegram chats, Google Docs, and manual transfers — seven
          to eight working days from submission to a contributor seeing payment. This platform
          replaces every handoff that didn't need to exist with one auditable workflow, settled
          onchain via Core wallet the same day an article goes live.
        </p>
        <Link
          to="/login"
          className="inline-flex items-center justify-center rounded-xl bg-brand-red px-6 py-3 text-sm font-semibold text-white shadow-glow-red transition-colors hover:bg-brand-red-dark"
        >
          Sign in to your dashboard
        </Link>
        <p className="mt-4 text-xs text-zinc-500">
          Access is invite-only. Don't have an invitation yet? Ask your Super Admin.
        </p>
      </section>

      {/* Lifecycle */}
      <section className="border-t border-surface-border bg-surface-base px-6 py-16 sm:px-10">
        <div className="mx-auto max-w-5xl">
          <h2 className="mb-2 text-center text-2xl font-semibold">How an article moves through the platform</h2>
          <p className="mb-10 text-center text-sm text-zinc-500">
            Every article follows the same path. Every step notifies exactly who needs to act next.
          </p>
          <div className="flex flex-wrap items-center justify-center gap-2">
            {LIFECYCLE.map((step, i) => (
              <div key={step} className="flex items-center gap-2">
                <span className="rounded-full border border-surface-border bg-surface-card px-4 py-2 text-sm font-medium text-zinc-200">
                  {step}
                </span>
                {i < LIFECYCLE.length - 1 && <span className="text-zinc-600">&rarr;</span>}
              </div>
            ))}
          </div>
        </div>
      </section>

      {/* Role guides */}
      <section className="px-6 py-16 sm:px-10">
        <div className="mx-auto max-w-4xl">
          <h2 className="mb-2 text-center text-2xl font-semibold">Find your role</h2>
          <p className="mb-10 text-center text-sm text-zinc-500">
            Click a role to see exactly what you'll do on the platform.
          </p>

          <div className="space-y-3">
            {ROLE_GUIDES.map((role) => {
              const isOpen = openRole === role.id;
              return (
                <div
                  key={role.id}
                  className="overflow-hidden rounded-xl2 border border-surface-border bg-surface-card"
                >
                  <button
                    onClick={() => setOpenRole(isOpen ? null : role.id)}
                    className="flex w-full items-center justify-between gap-4 px-5 py-4 text-left"
                    aria-expanded={isOpen}
                  >
                    <span>
                      <span className="block font-semibold text-zinc-100">{role.label}</span>
                      <span className="block text-sm text-zinc-500">{role.tagline}</span>
                    </span>
                    <span
                      className={`shrink-0 text-zinc-500 transition-transform ${isOpen ? "rotate-180" : ""}`}
                      aria-hidden
                    >
                      ▾
                    </span>
                  </button>

                  {isOpen && (
                    <div className="border-t border-surface-border px-5 py-4">
                      <ul className="space-y-2.5">
                        {role.bullets.map((b, i) => (
                          <li key={i} className="flex gap-3 text-sm text-zinc-300">
                            <span className="mt-1.5 h-1.5 w-1.5 shrink-0 rounded-full bg-brand-red" />
                            <span>{b}</span>
                          </li>
                        ))}
                      </ul>
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        </div>
      </section>

      {/* New here */}
      <section className="border-t border-surface-border bg-surface-base px-6 py-16 sm:px-10">
        <div className="mx-auto max-w-3xl text-center">
          <h2 className="mb-4 text-2xl font-semibold">New here?</h2>
          <p className="text-sm leading-relaxed text-zinc-400">
            Access is invitation-only. Your Super Admin sends you an email with a registration
            link valid for 72 hours — open it, set a password, fill in a short bio, and (if
            you're joining as a contributor) connect the Core wallet you want paid to. From there
            you land straight on your dashboard, ready to go.
          </p>
        </div>
      </section>

      <footer className="px-6 py-8 text-center text-xs text-zinc-600 sm:px-10">
        Team1 Blog Contributor Platform — internal use only.
      </footer>
    </div>
  );
}
