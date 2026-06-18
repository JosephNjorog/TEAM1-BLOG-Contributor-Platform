import { Navigate } from "react-router-dom";
import { Spinner } from "@team1/shared";
import { useAuth } from "../lib/auth";

const FRONTEND_URL = import.meta.env.VITE_FRONTEND_URL ?? "http://localhost:5173";

export function ProtectedRoute({ children }: { children: React.ReactNode }) {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-surface-app">
        <Spinner />
      </div>
    );
  }

  if (!user) return <Navigate to="/login" replace />;

  if (user.role !== "super_admin") {
    return (
      <div className="flex h-screen flex-col items-center justify-center gap-3 bg-surface-app px-4 text-center">
        <p className="text-zinc-300">This dashboard is for Super Admins only.</p>
        <a href={FRONTEND_URL} className="text-sm font-medium text-brand-red hover:underline">
          Go to your dashboard →
        </a>
      </div>
    );
  }

  return <>{children}</>;
}
