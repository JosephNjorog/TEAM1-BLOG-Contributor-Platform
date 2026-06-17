import React from "react";
import ReactDOM from "react-dom/client";
import { BrowserRouter } from "react-router-dom";
import { QueryClientProvider } from "@tanstack/react-query";
import { configureApiClient } from "@team1/shared";
import { queryClient } from "./lib/queryClient";
import { AuthProvider } from "./lib/auth";
import { App } from "./App";
import "./index.css";

configureApiClient({ baseUrl: "/api/v1" });

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
