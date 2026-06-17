import { Navigate } from "react-router-dom";
import type { Role } from "@team1/shared";
import { Spinner } from "@team1/shared";
import { useAuth } from "../lib/auth";

export function ProtectedRoute({ allow, children }: { allow: Role[]; children: React.ReactNode }) {
  const { user, isLoading } = useAuth();

  if (isLoading) {
    return (
      <div className="flex h-screen items-center justify-center bg-surface-app">
        <Spinner />
      </div>
    );
  }

  if (!user) return <Navigate to="/login" replace />;
  if (!allow.includes(user.role)) return <Navigate to="/app" replace />;

  return <>{children}</>;
}
