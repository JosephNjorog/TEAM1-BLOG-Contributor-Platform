import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { configureApiClient } from "@team1/shared";
import { queryClient } from "./lib/queryClient";
import { AuthProvider } from "./lib/auth";
import { App } from "./App";
import "./index.css";

// Relative by default - works when this app and the API share an origin
// (e.g. a reverse proxy routes /api to the backend). Set VITE_API_BASE_URL
// at build time for deployments where they're on separate hosts/domains.
configureApiClient({ baseUrl: import.meta.env.VITE_API_BASE_URL ?? "/api/v1" });

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <BrowserRouter>
          <App />
        </BrowserRouter>
      </AuthProvider>
    </QueryClientProvider>
  </React.StrictMode>,
);
