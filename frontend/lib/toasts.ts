import { DateTime, Duration } from "luxon";
import type { Snippet } from "svelte";
import { writable } from "svelte/store";

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

export const toasts = writable<AssignedToast[]>([
  {
    message: "Hello, world!",
    description: "This is a test toast.",
    urgency: "info",
    timeout: Duration.fromObject({ year: 5 }),
    id: 0,
    deleteAt: DateTime.now().plus(Duration.fromObject({ year: 5 })),
  },
  {
    message: "Hello, world!",
    urgency: "error",
    timeout: Duration.fromObject({ year: 5 }),
    id: 1,
    deleteAt: DateTime.now().plus(Duration.fromObject({ year: 5 })),
  },
]);

let nextToastID = 0;

export function addToast(toast: Toast) {
  toasts.update((toasts) => {
    toast.timeout ??= Duration.fromObject({ seconds: 5 });
    toast.urgency ??= "info";

    toasts.push({
      id: ++nextToastID,
      deleteAt: DateTime.now().plus(toast.timeout),
      ...toast,
    });

    setTimeout(clearToasts, toast.timeout.as("milliseconds"));
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
