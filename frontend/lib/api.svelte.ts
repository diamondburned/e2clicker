import * as api from "./openapi.gen.js";
export * from "./openapi.gen.js";

// token is the current session token.
// It is persisted to the local storage automatically.
export let token = $state(localStorage.getItem("token"));

// isLoggedIn returns true if the user is logged in.
export const isLoggedIn = () => !!token;

export function init() {
  $effect(() => {
    token ? localStorage.setItem("token", token) : localStorage.removeItem("token");
  });
  $effect(() => {
    api.defaults.headers = { Authorization: `Bearer ${token}` };
  });
}
