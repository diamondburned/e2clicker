/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

import { cleanupOutdatedCaches } from "workbox-precaching";

import type * as api from "./lib/api";
// import { updatePushSubscription } from "./lib/notification";
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

self.addEventListener("notificationclick", (ev) => {
  handleEvent(ev, async () => {
    const notification = ev.notification.data as api.Notification | undefined;
    if (!notification) {
      postWorkerMessage("error", new WorkerError("notification data is missing"));
      return;
    }

    const route = {
      welcome: null,
      reminder: "/dashboard",
      account_notice: "/settings",
      web_push_expiring: "/settings",
    }[notification.type];
    if (route) {
      await gotoWindow(route);
    }

    ev.notification.close();
  });
});

/*
self.addEventListener("pushsubscriptionchange", (ev) => {
  // Supposedly, Chrome has no support for this event. However, we're handling
  // this event just in case. It is fully supported in Firefox and Safari.
  //
  // For more information, see:
  // - https://stackoverflow.com/questions/61487043/why-are-my-pushsubscriptions-expiring-so-quickly
  // - https://issues.chromium.org/issues/41275327
  // - https://developer.mozilla.org/en-US/docs/Web/API/ServiceWorkerGlobalScope/pushsubscriptionchange_event#browser_compatibility
  if ("waitUntil" in ev) {
    try {
      handleEvent(ev as ExtendableEvent, async () => {
        await updatePushSubscription({ toggle: true });
      });
      return;
    } catch (err) {
      console.error("pushsubscriptionchange error", err);
      // defer to showing a notification as below.
    }
  }

  self.registration.showNotification("Notifications broke!", {
    body:
      "A web quirk happened... " +
      "Please re-open the app or switch to more reliable notifying methods.",
    tag: "_system",
    requireInteraction: true,
  });
});
*/

function handleEvent<Event extends ExtendableEvent>(ev: Event, fn: () => Promise<any>) {
  ev.waitUntil(
    fn().catch((err) => {
      console.error("Service Worker event failed", err);
      postWorkerMessage("error", new WorkerError("failed to handle event", { cause: err }));
    }),
  );
}

async function gotoWindow(url: string) {
  const windows = await self.clients.matchAll({ type: "window" });

  const window = windows.find((w) => w.url === url);
  if (window) {
    await window.focus();
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
