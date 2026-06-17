import { useState } from "react";
import { useMarkAllNotificationsRead, useMarkNotificationRead, useNotifications } from "../lib/notifications";

function timeAgo(iso: string): string {
  const diffMs = Date.now() - new Date(iso).getTime();
  const mins = Math.floor(diffMs / 60_000);
  if (mins < 1) return "just now";
  if (mins < 60) return `${mins}m ago`;
  const hours = Math.floor(mins / 60);
  if (hours < 24) return `${hours}h ago`;
  return `${Math.floor(hours / 24)}d ago`;
}

export function NotificationBell() {
  const [open, setOpen] = useState(false);
  const { data } = useNotifications();
  const markRead = useMarkNotificationRead();
  const markAllRead = useMarkAllNotificationsRead();

  const unread = data?.unreadCount ?? 0;

  return (
    <div className="relative">
      <button
        onClick={() => setOpen((o) => !o)}
        className="relative rounded-full p-2 text-zinc-400 hover:bg-surface-raised hover:text-zinc-100"
        aria-label="Notifications"
      >
        <BellIcon />
        {unread > 0 && (
          <span className="absolute -right-0.5 -top-0.5 flex h-4 min-w-4 items-center justify-center rounded-full bg-brand-red px-1 text-[10px] font-semibold text-white">
            {unread > 9 ? "9+" : unread}
          </span>
        )}
      </button>

      {open && (
        <div className="absolute right-0 z-20 mt-2 w-80 rounded-xl2 border border-surface-border bg-surface-card shadow-xl">
          <div className="flex items-center justify-between border-b border-surface-border px-4 py-3">
            <span className="text-sm font-semibold text-zinc-100">Notifications</span>
            {unread > 0 && (
              <button
                onClick={() => markAllRead.mutate()}
                className="text-xs font-medium text-brand-red hover:underline"
              >
                Mark all read
              </button>
            )}
          </div>
          <div className="max-h-96 overflow-y-auto">
            {!data || data.notifications.length === 0 ? (
              <p className="px-4 py-6 text-center text-sm text-zinc-500">You're all caught up.</p>
            ) : (
              data.notifications.map((n) => (
                <button
                  key={n.id}
                  onClick={() => !n.read && markRead.mutate(n.id)}
                  className={`block w-full border-b border-surface-border/60 px-4 py-3 text-left text-sm last:border-b-0 hover:bg-surface-raised ${
                    n.read ? "text-zinc-500" : "text-zinc-100"
                  }`}
                >
                  <span className="block">{n.message}</span>
                  <span className="mt-1 block text-xs text-zinc-500">{timeAgo(n.createdAt)}</span>
                </button>
              ))
            )}
          </div>
        </div>
      )}
    </div>
  );
}

function BellIcon() {
  return (
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
      <path d="M18 8a6 6 0 0 0-12 0c0 7-3 9-3 9h18s-3-2-3-9" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M13.73 21a2 2 0 0 1-3.46 0" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}
