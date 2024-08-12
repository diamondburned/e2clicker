import * as api from "./api/api.js";
export * from "./api/api.js";

// token is the current session token.
// It is persisted to the local storage automatically.
export let token = $state(localStorage.getItem("token"));

$effect(() => {
  token ? localStorage.setItem("token", token) : localStorage.removeItem("token");
});
$effect(() => {
  api.defaults.headers = { Authorization: `Bearer ${token}` };
});

// userID returns the user ID of the current session if any.
export const userID = () => (token ? token.split(".")[0] : null);

// isLoggedIn returns true if the user is logged in.
export const isLoggedIn = () => !!userID;
