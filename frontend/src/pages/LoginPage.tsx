import { FormEvent, useState } from "react";
import { Navigate, useLocation, useNavigate } from "react-router-dom";
import { Button, ApiClientError } from "@team1/shared";
import { useAuth } from "../lib/auth";

export function LoginPage() {
  const { user, login } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  if (user) return <Navigate to="/" replace />;

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);
    try {
      await login(email, password);
      const redirectTo = (location.state as { from?: string } | null)?.from ?? "/";
      navigate(redirectTo, { replace: true });
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : "Something went wrong. Try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-surface-app px-4">
      <div className="w-full max-w-sm">
        <div className="mb-8 flex items-center justify-center gap-2">
          <span className="h-2.5 w-2.5 rounded-full bg-brand-red shadow-glow-red" />
          <span className="text-sm font-bold uppercase tracking-wide text-zinc-100">Team1 Blog</span>
        </div>

        <form
          onSubmit={onSubmit}
          className="rounded-xl2 border border-surface-border bg-surface-card p-6 shadow-xl"
        >
          <h1 className="mb-1 text-lg font-semibold text-zinc-100">Sign in</h1>
          <p className="mb-6 text-sm text-zinc-500">Contributor Platform — invite only.</p>

          <label className="mb-1 block text-xs font-medium text-zinc-400">Email</label>
          <input
            type="email"
            required
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="mb-4 w-full rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
            placeholder="you@team1.blog"
          />

          <label className="mb-1 block text-xs font-medium text-zinc-400">Password</label>
          <input
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="mb-5 w-full rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
            placeholder="••••••••"
          />

          {error && <p className="mb-4 text-sm text-red-400">{error}</p>}

          <Button type="submit" loading={loading} className="w-full">
            Sign in
          </Button>
        </form>
      </div>
    </div>
  );
}
