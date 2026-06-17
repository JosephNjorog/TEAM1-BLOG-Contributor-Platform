import { createContext, ReactNode, useContext, useEffect, useState } from "react";
import {
  apiRequest,
  ApiClientError,
  clearTokens,
  getAccessToken,
  getRefreshToken,
  setTokens,
  type AuthTokens,
  type User,
} from "@team1/shared";

interface AuthResponse {
  user: User;
  accessToken: string;
  refreshToken: string;
  expiresAt: string;
}

interface AuthContextValue {
  user: User | null;
  isLoading: boolean;
  login: (email: string, password: string) => Promise<User>;
  registerFromInvite: (input: {
    token: string;
    name: string;
    password: string;
    bio?: string;
    walletAddress?: string;
  }) => Promise<User>;
  logout: () => Promise<void>;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

function storeAuthResponse(res: AuthResponse) {
  const tokens: AuthTokens = {
    accessToken: res.accessToken,
    refreshToken: res.refreshToken,
    expiresAt: res.expiresAt,
  };
  setTokens(tokens);
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  const refreshUser = async () => {
    if (!getAccessToken()) {
      setUser(null);
      return;
    }
    try {
      const me = await apiRequest<User>("/me");
      setUser(me);
    } catch (err) {
      if (err instanceof ApiClientError && err.status === 401) {
        clearTokens();
      }
      setUser(null);
    }
  };

  useEffect(() => {
    refreshUser().finally(() => setIsLoading(false));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const login = async (email: string, password: string) => {
    const res = await apiRequest<AuthResponse>("/auth/login", {
      method: "POST",
      body: { email, password },
      skipAuth: true,
    });
    storeAuthResponse(res);
    setUser(res.user);
    return res.user;
  };

  const registerFromInvite: AuthContextValue["registerFromInvite"] = async (input) => {
    const res = await apiRequest<AuthResponse>("/auth/register", {
      method: "POST",
      body: input,
      skipAuth: true,
    });
    storeAuthResponse(res);
    setUser(res.user);
    return res.user;
  };

  const logout = async () => {
    const refreshToken = getRefreshToken();
    if (refreshToken) {
      try {
        await apiRequest("/auth/logout", { method: "POST", body: { refreshToken }, skipAuth: true });
      } catch {
        /* best-effort revoke */
      }
    }
    clearTokens();
    setUser(null);
  };

  return (
    <AuthContext.Provider value={{ user, isLoading, login, registerFromInvite, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used within an AuthProvider");
  return ctx;
}
