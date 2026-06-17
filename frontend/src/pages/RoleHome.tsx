import { useAuth } from "../lib/auth";
import { ArticlesPage } from "./contributor/ArticlesPage";
import { ComingSoonPage } from "./ComingSoonPage";

const titles: Record<string, string> = {
  moderator: "Review Queue",
  graphic_designer: "Banner Queue",
  publisher: "Ready to Publish",
  super_admin: "Overview",
};

export function RoleHome() {
  const { user } = useAuth();
  if (user?.role === "contributor") return <ArticlesPage />;
  return <ComingSoonPage title={titles[user?.role ?? ""] ?? "Dashboard"} />;
}
