import { Link, NavLink, Outlet } from "react-router-dom";
import { useAuth } from "../lib/auth";
import { NotificationBell } from "./NotificationBell";

const roleLabels: Record<string, string> = {
  super_admin: "Super Admin",
  moderator: "Moderator",
  graphic_designer: "Graphic Designer",
  publisher: "Publisher",
  contributor: "Contributor",
};

interface NavItem {
  to: string;
  label: string;
}

const navByRole: Record<string, NavItem[]> = {
  contributor: [
    { to: "/app", label: "My Articles" },
    { to: "/app/published", label: "Published" },
    { to: "/app/my-payments", label: "Payments" },
  ],
  moderator: [{ to: "/app", label: "Review Queue" }],
  graphic_designer: [{ to: "/app", label: "Banner Queue" }],
  publisher: [{ to: "/app", label: "Ready to Publish" }],
  super_admin: [], // Super Admins use the separate super-admin app instead
};

export function AppShell() {
  const { user, logout } = useAuth();
  const navItems = navByRole[user?.role ?? ""] ?? [];

  return (
    <div className="flex min-h-screen bg-surface-app">
      <aside className="flex w-64 flex-col border-r border-surface-border bg-surface-base">
        <Link to="/" className="flex items-center gap-2 px-6 py-5">
          <span className="h-2.5 w-2.5 rounded-full bg-brand-red shadow-glow-red" />
          <span className="text-sm font-bold uppercase tracking-wide text-zinc-100">Team1 Blog</span>
        </Link>
        <nav className="flex-1 space-y-1 px-3 py-2">
          {navItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              end
              className={({ isActive }) =>
                `block rounded-lg px-3 py-2 text-sm font-medium transition-colors ${
                  isActive ? "bg-surface-raised text-zinc-100" : "text-zinc-400 hover:bg-surface-raised hover:text-zinc-100"
                }`
              }
            >
              {item.label}
            </NavLink>
          ))}
        </nav>
        <div className="border-t border-surface-border px-4 py-4">
          <p className="truncate text-sm font-medium text-zinc-100">{user?.name}</p>
          <p className="text-xs text-zinc-500">{user ? roleLabels[user.role] : ""}</p>
          <button
            onClick={() => logout()}
            className="mt-3 text-xs font-medium text-zinc-500 hover:text-brand-red"
          >
            Sign out
          </button>
        </div>
      </aside>

      <div className="flex flex-1 flex-col">
        <header className="flex h-16 items-center justify-end border-b border-surface-border px-6">
          <NotificationBell />
        </header>
        <main className="flex-1 overflow-y-auto p-8">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
