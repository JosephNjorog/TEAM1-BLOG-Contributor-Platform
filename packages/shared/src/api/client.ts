import type { ApiError, AuthTokens } from "../types";

const ACCESS_TOKEN_KEY = "team1_access_token";
const REFRESH_TOKEN_KEY = "team1_refresh_token";

export function getAccessToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

export function setTokens(tokens: AuthTokens): void {
  localStorage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken);
}

export function clearTokens(): void {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
}

export class ApiClientError extends Error {
  status: number;
  body: ApiError | null;
  constructor(status: number, body: ApiError | null) {
    super(body?.message ?? `Request failed with status ${status}`);
    this.status = status;
    this.body = body;
  }
}

let baseUrl = "/api/v1";
let refreshing: Promise<boolean> | null = null;

export function configureApiClient(opts: { baseUrl: string }): void {
  baseUrl = opts.baseUrl;
}

// Builds an absolute ws(s):// URL against the current page's origin, using
// the same baseUrl apiRequest uses for plain HTTP - in dev that's proxied
// to the backend by Vite (vite.config.ts sets ws: true on the /api proxy);
// in production whatever serves the frontend needs to proxy WS upgrades
// under this same path too.
export function wsURL(path: string): string {
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  return `${protocol}//${window.location.host}${baseUrl}${path}`;
}

async function tryRefresh(): Promise<boolean> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) return false;
  if (!refreshing) {
    refreshing = fetch(`${baseUrl}/auth/refresh`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ refreshToken }),
    })
      .then(async (res) => {
        if (!res.ok) {
          clearTokens();
          return false;
        }
        const data: AuthTokens = await res.json();
        setTokens(data);
        return true;
      })
      .catch(() => {
        clearTokens();
        return false;
      })
      .finally(() => {
        refreshing = null;
      });
  }
  return refreshing;
}

export interface RequestOptions {
  method?: "GET" | "POST" | "PUT" | "PATCH" | "DELETE";
  body?: unknown;
  signal?: AbortSignal;
  skipAuth?: boolean;
}

export async function apiRequest<T>(path: string, opts: RequestOptions = {}): Promise<T> {
  const isFormData = opts.body instanceof FormData;

  const doFetch = async (): Promise<Response> => {
    const headers: Record<string, string> = {};
    // Let the browser set the multipart boundary itself for FormData bodies.
    if (!isFormData) headers["Content-Type"] = "application/json";
    if (!opts.skipAuth) {
      const token = getAccessToken();
      if (token) headers.Authorization = `Bearer ${token}`;
    }
    return fetch(`${baseUrl}${path}`, {
      method: opts.method ?? "GET",
      headers,
      body: isFormData ? (opts.body as FormData) : opts.body !== undefined ? JSON.stringify(opts.body) : undefined,
      signal: opts.signal,
    });
  };

  let res = await doFetch();

  if (res.status === 401 && !opts.skipAuth) {
    const refreshed = await tryRefresh();
    if (refreshed) {
      res = await doFetch();
    }
  }

  if (!res.ok) {
    let body: ApiError | null = null;
    try {
      body = await res.json();
    } catch {
      body = null;
    }
    throw new ApiClientError(res.status, body);
  }

  if (res.status === 204) return undefined as T;
  return res.json() as Promise<T>;
}
