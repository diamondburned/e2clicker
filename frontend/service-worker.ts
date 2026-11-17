/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

import { cleanupOutdatedCaches, precache } from "workbox-precaching";

import type * as api from "./lib/api";
import { loadSWCredentials } from "./lib/sw-creds";
import { WorkerError, type WorkerMessage } from "./lib/shared-worker";

// TypeScript support:
// https://svelte.dev/docs/kit/service-workers#Type-safety
// https://vite-pwa-org.netlify.app/guide/inject-manifest.html#service-worker-code-3
declare let self: ServiceWorkerGlobalScope;

// Clean up old workbox caches.
cleanupOutdatedCaches();

// Do not route the service worker. We're not trying to make this work offline
// yet, and there are some issues with doing so.
precache(self.__WB_MANIFEST);

self.addEventListener("message", (ev) => {
  console.debug("Service Worker received message", ev.data);
  if (ev.data && ev.data.type === "SKIP_WAITING") {
    self.skipWaiting();
  }
});

self.addEventListener("push", (ev) => {
  handleEvent(ev, async () => {
    console.debug("Service Worker received push", ev.data);

    if (!ev.data) {
      throw new Error("no data in push event");
    }

    let notification: api.Notification | undefined;
    let message: api.NotificationMessage;
    try {
      notification = ev.data?.json() as api.Notification;
      message = notification.message;
    } catch (err) {
      const text = ev.data.text();
      message = { title: "Notification", message: text };
    }

    // TODO: handle notification actions like Snooze.
    await self.registration.showNotification(message.title, {
      body: message.message,
      data: notification,
      tag: notification?.type,
      requireInteraction: true,
      silent: false,
    });
  });
});

const notificationRoutes: Record<api.NotificationType, string | null> = {
  welcome_message: "/dashboard",
  reminder_message: "/dashboard",
  account_notice_message: "/settings",
  web_push_expiring_message: "/settings",
  test_message: "/dashboard",
};

self.addEventListener("notificationclick", (ev) => {
  handleEvent(ev, async () => {
    const notification = ev.notification.data as api.Notification | undefined;
    if (!notification) {
      postWorkerMessage("error", new WorkerError("notification data is missing"));
      return;
    }

    const route = notificationRoutes[notification.type];
    if (route) {
      await gotoWindow(route);
    }

    ev.notification.close();
  });
});

self.addEventListener("pushsubscriptionchange", (ev) => {
  const pushEv = ev as ExtendableEvent & {
    newSubscription?: PushSubscription;
  };
  pushEv.waitUntil(
    (async () => {
      try {
        const creds = await loadSWCredentials();
        if (!creds) {
          throw new Error("missing stored credentials for push resubscription");
        }

        let subscription = pushEv.newSubscription;
        if (!subscription) {
          const infoResp = await fetch("/api/push-info");
          if (!infoResp.ok) {
            throw new Error(`failed to fetch push-info: ${infoResp.status}`);
          }
          const info = (await infoResp.json()) as api.PushInfo;
          subscription = await self.registration.pushManager.subscribe({
            userVisibleOnly: true,
            applicationServerKey: info.applicationServerKey,
          });
        }

        const keys = subscription.toJSON().keys as Record<string, string> | undefined;
        if (!keys?.p256dh || !keys?.auth) {
          throw new Error("push subscription keys missing");
        }

        const expirationTime = subscription.expirationTime
          ? new Date(subscription.expirationTime).toISOString()
          : undefined;

        const resp = await fetch(`/api/notifications/web-push/${encodeURIComponent(creds.deviceID)}`, {
          method: "PUT",
          headers: {
            "Content-Type": "application/json",
            Authorization: `Bearer ${creds.token}`,
          },
          body: JSON.stringify({
            deviceID: creds.deviceID,
            endpoint: subscription.endpoint,
            expirationTime,
            keys: { p256dh: keys.p256dh, auth: keys.auth },
          }),
        });
        if (!resp.ok) {
          throw new Error(`failed to update push subscription: ${resp.status}`);
        }
      } catch (err) {
        console.error("pushsubscriptionchange error", err);
        await self.registration.showNotification("Notifications broke!", {
          body:
            "A web quirk happened... " +
            "Please re-open the app or switch to more reliable notifying methods.",
          tag: "_system",
          requireInteraction: true,
        });
      }
    })(),
  );
});

function handleEvent<Event extends ExtendableEvent>(ev: Event, fn: () => Promise<any>) {
  ev.waitUntil(
    fn().catch((err) => {
      console.error("Service Worker event failed", err);
      postWorkerMessage("error", new WorkerError("failed to handle event", { cause: err }));
    }),
  );
}

async function gotoWindow(url: string) {
  const windows = await self.clients.matchAll({ type: "window", includeUncontrolled: true });

  const existing = windows.find((w) => {
    try {
      return new URL(w.url).pathname === url;
    } catch {
      return false;
    }
  });
  if (existing) {
    await existing.focus();
    return;
  }

  if (self.clients.openWindow) {
    await self.clients.openWindow(url);
  }
}

async function postWorkerMessage<
  Type extends string,
  Message extends Extract<WorkerMessage, { type: Type }>,
>(type: Type, data: Message["data"]) {
  const clients = await self.clients.matchAll();
  for (const client of clients) {
    client.postMessage({ type, data });
  }
}
