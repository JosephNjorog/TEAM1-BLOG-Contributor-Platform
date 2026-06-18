import { useAuth } from "../lib/auth";
import { ArticlesPage } from "./contributor/ArticlesPage";
import { ReviewQueuePage } from "./moderator/ReviewQueuePage";
import { BannerQueuePage } from "./designer/BannerQueuePage";
import { ReadyToPublishPage } from "./publisher/ReadyToPublishPage";

const ADMIN_APP_URL = import.meta.env.VITE_ADMIN_APP_URL ?? "http://localhost:5174";

export function RoleHome() {
  const { user } = useAuth();
  switch (user?.role) {
    case "contributor":
      return <ArticlesPage />;
    case "moderator":
      return <ReviewQueuePage />;
    case "graphic_designer":
      return <BannerQueuePage />;
    case "publisher":
      return <ReadyToPublishPage />;
    case "super_admin":
      return (
        <div className="flex flex-col items-center justify-center gap-3 py-24 text-center">
          <p className="text-zinc-300">Super Admins use the dedicated admin dashboard.</p>
          <a href={ADMIN_APP_URL} className="text-sm font-medium text-brand-red hover:underline">
            Go to the admin dashboard →
          </a>
        </div>
      );
    default:
      return null;
  }
}
