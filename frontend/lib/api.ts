import { persisted } from "svelte-persisted-store";

import * as api from "./openapi.gen.js";
import { derived } from "svelte/store";
export * from "./openapi.gen.js";

// token is the current session token.
// It is persisted to the local storage automatically.
export const token = persisted<string | null>("e2clicker-token", null);

// isLoggedIn is true if the user is logged in.
export const isLoggedIn = derived(token, (token) => !!token);

token.subscribe((token) => {
  api.defaults.headers = { Authorization: `Bearer ${token}` };
});
