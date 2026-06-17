import { FormEvent, useState } from "react";
import { Navigate, useNavigate, useSearchParams } from "react-router-dom";
import { Button, ApiClientError } from "@team1/shared";
import { useAuth } from "../lib/auth";

const AVAX_ADDR_RE = /^0x[a-fA-F0-9]{40}$/;

export function RegisterPage() {
  const { user, registerFromInvite } = useAuth();
  const navigate = useNavigate();
  const [params] = useSearchParams();
  const token = params.get("token") ?? "";

  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [bio, setBio] = useState("");
  const [walletAddress, setWalletAddress] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);

  if (user) return <Navigate to="/" replace />;

  if (!token) {
    return (
      <div className="flex min-h-screen items-center justify-center bg-surface-app px-4 text-center">
        <p className="text-sm text-zinc-400">
          This page needs an invitation link. Check the email you were sent, or ask your Super Admin to resend it.
        </p>
      </div>
    );
  }

  const walletInvalid = walletAddress !== "" && !AVAX_ADDR_RE.test(walletAddress);

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    if (password.length < 8) {
      setError("Password must be at least 8 characters.");
      return;
    }
    if (walletInvalid) {
      setError("Wallet address must be a valid Avalanche C-Chain address (0x + 40 hex characters).");
      return;
    }
    setLoading(true);
    try {
      await registerFromInvite({ token, name, password, bio, walletAddress });
      navigate("/", { replace: true });
    } catch (err) {
      setError(err instanceof ApiClientError ? err.message : "Something went wrong. Try again.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="flex min-h-screen items-center justify-center bg-surface-app px-4 py-10">
      <div className="w-full max-w-md">
        <div className="mb-8 flex items-center justify-center gap-2">
          <span className="h-2.5 w-2.5 rounded-full bg-brand-red shadow-glow-red" />
          <span className="text-sm font-bold uppercase tracking-wide text-zinc-100">Team1 Blog</span>
        </div>

        <form onSubmit={onSubmit} className="rounded-xl2 border border-surface-border bg-surface-card p-6 shadow-xl">
          <h1 className="mb-1 text-lg font-semibold text-zinc-100">Accept your invitation</h1>
          <p className="mb-6 text-sm text-zinc-500">Set up your account to join the platform.</p>

          <label className="mb-1 block text-xs font-medium text-zinc-400">Display name</label>
          <input
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="mb-4 w-full rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
            placeholder="Ada Lovelace"
          />

          <label className="mb-1 block text-xs font-medium text-zinc-400">Password</label>
          <input
            type="password"
            required
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="mb-4 w-full rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
            placeholder="At least 8 characters"
          />

          <label className="mb-1 block text-xs font-medium text-zinc-400">Bio (optional)</label>
          <textarea
            value={bio}
            onChange={(e) => setBio(e.target.value)}
            rows={3}
            className="mb-4 w-full rounded-lg border border-surface-border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red"
            placeholder="A short line about you"
          />

          <label className="mb-1 block text-xs font-medium text-zinc-400">
            Core wallet address <span className="text-zinc-600">(contributors only)</span>
          </label>
          <input
            value={walletAddress}
            onChange={(e) => setWalletAddress(e.target.value)}
            className={`mb-1 w-full rounded-lg border bg-surface-base px-3 py-2 text-sm text-zinc-100 outline-none focus:border-brand-red ${
              walletInvalid ? "border-red-700" : "border-surface-border"
            }`}
            placeholder="0x..."
          />
          <p className="mb-5 text-xs text-zinc-600">Payments for published articles are sent here. Skip if this isn't a contributor account.</p>

          {error && <p className="mb-4 text-sm text-red-400">{error}</p>}

          <Button type="submit" loading={loading} className="w-full">
            Create account
          </Button>
        </form>
      </div>
    </div>
  );
}
