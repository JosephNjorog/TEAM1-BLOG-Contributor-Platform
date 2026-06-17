import { useNavigate } from "react-router-dom";
import { Card, EmptyState, Spinner } from "@team1/shared";
import { useArticles } from "../../lib/articles";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

export function BannerQueuePage() {
  const { data: articles, isLoading } = useArticles();
  const navigate = useNavigate();

  const queue = (articles ?? []).filter((a) => a.status === "editorial_approved");

  return (
    <div className="mx-auto max-w-4xl">
      <h1 className="mb-1 text-xl font-semibold text-zinc-100">Banner Queue</h1>
      <p className="mb-6 text-sm text-zinc-500">Editorially approved articles waiting on a cover banner.</p>

      {isLoading ? (
        <div className="flex justify-center py-12">
          <Spinner />
        </div>
      ) : queue.length === 0 ? (
        <EmptyState title="Nothing waiting on a banner" hint="Approved articles will show up here." />
      ) : (
        <div className="space-y-3">
          {queue.map((a) => (
            <Card
              key={a.id}
              className="cursor-pointer transition-colors hover:border-zinc-700"
              onClick={() => navigate(`/app/banner/${a.id}`)}
            >
              <div className="flex items-center justify-between gap-4">
                <div className="min-w-0">
                  <p className="truncate font-medium text-zinc-100">{a.title}</p>
                  <p className="mt-1 text-xs text-zinc-500">
                    by {a.contributorName} &middot; approved {formatDate(a.updatedAt)}
                  </p>
                </div>
                {a.cloudinaryBannerUrl && (
                  <img src={a.cloudinaryBannerUrl} alt="" className="h-12 w-20 shrink-0 rounded-lg object-cover" />
                )}
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
