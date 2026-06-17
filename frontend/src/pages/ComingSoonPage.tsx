import { EmptyState } from "@team1/shared";

export function ComingSoonPage({ title }: { title: string }) {
  return (
    <div className="mx-auto max-w-2xl pt-16">
      <h1 className="mb-6 text-xl font-semibold text-zinc-100">{title}</h1>
      <EmptyState
        title="This dashboard is coming in the next build phase"
        hint="The platform is being delivered in stages — this role's screens are next up."
      />
    </div>
  );
}
