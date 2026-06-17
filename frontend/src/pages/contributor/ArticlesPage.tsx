import { useNavigate } from "react-router-dom";
import { Button, Card, EmptyState, Spinner, StatusPill } from "@team1/shared";
import { useArticles, useCreateArticle } from "../../lib/articles";

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, { month: "short", day: "numeric", year: "numeric" });
}

export function ArticlesPage() {
  const { data: articles, isLoading } = useArticles();
  const createArticle = useCreateArticle();
  const navigate = useNavigate();

  const onNewDraft = async () => {
    const article = await createArticle.mutateAsync({ title: "Untitled draft", content: "" });
    navigate(`/app/articles/${article.id}`);
  };

  return (
    <div className="mx-auto max-w-4xl">
      <div className="mb-6 flex items-center justify-between">
        <div>
          <h1 className="text-xl font-semibold text-zinc-100">My Articles</h1>
          <p className="text-sm text-zinc-500">Drafts, submissions, and everything in between.</p>
        </div>
        <Button onClick={onNewDraft} loading={createArticle.isPending}>
          + New Draft
        </Button>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Spinner />
        </div>
      ) : !articles || articles.length === 0 ? (
        <EmptyState title="No articles yet" hint="Start your first draft to see it tracked here." />
      ) : (
        <div className="space-y-3">
          {articles.map((a) => (
            <Card
              key={a.id}
              className="cursor-pointer transition-colors hover:border-zinc-700"
              onClick={() => navigate(`/app/articles/${a.id}`)}
            >
              <div className="flex items-start justify-between gap-4">
                <div className="min-w-0">
                  <p className="truncate font-medium text-zinc-100">{a.title || "Untitled draft"}</p>
                  <p className="mt-1 text-xs text-zinc-500">
                    {a.wordCount} words · updated {formatDate(a.updatedAt)}
                    {a.reviewerName ? ` · reviewer: ${a.reviewerName}` : ""}
                  </p>
                </div>
                <StatusPill status={a.status} />
              </div>
            </Card>
          ))}
        </div>
      )}
    </div>
  );
}
