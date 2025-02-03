import { DateTime, Duration } from "luxon";
import type { Snippet } from "svelte";
import { writable } from "svelte/store";

export type Toast = {
  message: string | Snippet;
  description?: string | Snippet;
  error?: any;
  urgency?: "info" | "error";
  timeout?: Duration;
};

export type AssignedToast = Toast & {
  readonly id: number;
  readonly deleteAt: DateTime;
};

export const toasts = writable<AssignedToast[]>([]);

let nextToastID = 0;

export function logErrorToast(what: string, error: any) {
  console.error(what, error);
  addToast({
    urgency: "error",
    message: what,
    description: error.message ?? error,
  });
}

export function addToast(toast: Toast) {
  toast.timeout ??= Duration.fromObject({ seconds: 5 });
  toast.urgency ??= "info";

  toasts.update((toasts) => {
    toasts.push({
      id: ++nextToastID,
      deleteAt: DateTime.now().plus(toast.timeout!),
      ...toast,
    });

    setTimeout(clearToasts, toast.timeout!.as("milliseconds"));
    return toasts;
  });
}

export function clearToasts() {
  const now = DateTime.now();
  toasts.update((toasts) => toasts.filter((toast) => now.diff(toast.deleteAt).toMillis() <= 0));
}

export function closeToast(toast: Pick<AssignedToast, "id">) {
  toasts.update((toasts) => toasts.filter((t) => t.id != toast.id));
}
