import { DateTime, Duration } from "luxon";
import type { Snippet } from "svelte";

export type Toast = {
  message: string | Snippet;
  description?: string | Snippet;
  urgency?: "info" | "error";
  timeout?: Duration;
};

export type AssignedToast = Toast & {
  readonly id: number;
  readonly deleteAt: DateTime;
};

export let toasts = $state<AssignedToast[]>([]);

let nextToastID = 0;

export function addToast(toast: Toast) {
  toast.timeout ??= Duration.fromObject({ seconds: 5 });
  toast.urgency ??= "info";

  toasts.push({
    id: ++nextToastID,
    deleteAt: DateTime.now().plus(toast.timeout),
    ...toast,
  });

  setTimeout(clearToasts, toast.timeout.as("milliseconds"));
}

export function clearToasts() {
  const now = DateTime.now();
  toasts = toasts.filter((toast) => now.diff(toast.deleteAt).toMillis() <= 0);
}

export function closeToast(toast: Pick<AssignedToast, "id">) {
  toasts = toasts.filter((t) => t.id != toast.id);
}
