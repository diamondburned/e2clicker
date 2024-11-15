import { persisted } from "svelte-persisted-store";

import * as api from "./openapi.gen.js";
import { derived } from "svelte/store";
import { goto } from "$app/navigation";

export * from "./openapi.gen.js";

// token is the current session token.
// It is persisted to the local storage automatically.
export const token = persisted<string | null>("e2clicker-token", null);

// isLoggedIn is true if the user is logged in.
export const isLoggedIn = derived(token, (token) => !!token);

token.subscribe((token) => {
  api.defaults.headers = { Authorization: `Bearer ${token}` };
  api.defaults.baseUrl = "/api";
});

export function setToken(newToken: string | null) {
  token.set(newToken);
  goto("/");
}

// secretQRRegex is a regular expression that matches the secret QR code format.
export const secretQRRegex = /^e2clicker:secret-v1:(.*)$/;

// secretQRData returns the secret QR code data for a given secret.
export function secretQRData(secret: string): string {
  return `e2clicker:secret-v1:${secret}`;
}
