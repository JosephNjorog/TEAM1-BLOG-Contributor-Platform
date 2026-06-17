import { useAuth } from "../lib/auth";
import { ArticlesPage } from "./contributor/ArticlesPage";
import { ReviewQueuePage } from "./moderator/ReviewQueuePage";
import { BannerQueuePage } from "./designer/BannerQueuePage";
import { ReadyToPublishPage } from "./publisher/ReadyToPublishPage";
import { ComingSoonPage } from "./ComingSoonPage";

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
    default:
      return <ComingSoonPage title="Overview" />;
  }
}
