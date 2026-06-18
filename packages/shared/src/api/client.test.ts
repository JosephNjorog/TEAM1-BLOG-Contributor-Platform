import { describe, expect, it } from "vitest";
import { ApiClientError, configureApiClient, wsURL } from "./client";

describe("wsURL", () => {
  it("builds against the current page's origin when baseUrl is relative", () => {
    configureApiClient({ baseUrl: "/api/v1" });
    expect(wsURL("/notifications/ws")).toBe("wss://app.example.com/api/v1/notifications/ws");
  });

  it("builds against the configured host when baseUrl is absolute (http)", () => {
    configureApiClient({ baseUrl: "http://api.example.com/api/v1" });
    expect(wsURL("/notifications/ws")).toBe("ws://api.example.com/api/v1/notifications/ws");
  });

  it("builds against the configured host when baseUrl is absolute (https)", () => {
    configureApiClient({ baseUrl: "https://api.example.com/api/v1" });
    expect(wsURL("/notifications/ws")).toBe("wss://api.example.com/api/v1/notifications/ws");
  });
});

describe("ApiClientError", () => {
  it("uses the server message when present", () => {
    const err = new ApiClientError(403, { error: "forbidden", message: "you do not have access" });
    expect(err.message).toBe("you do not have access");
    expect(err.status).toBe(403);
  });

  it("falls back to a generic message when the body is null", () => {
    const err = new ApiClientError(500, null);
    expect(err.message).toBe("Request failed with status 500");
  });
});
