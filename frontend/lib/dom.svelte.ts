import { writable } from "svelte/store";

const prefersReducedMotionQuery = window.matchMedia("(prefers-reduced-motion: reduce)");
export let prefersReducedMotion = writable(prefersReducedMotionQuery.matches);
prefersReducedMotionQuery.addEventListener("change", () => {
  prefersReducedMotion.set(prefersReducedMotionQuery.matches);
});
