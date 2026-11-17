import * as api from "$lib/api";
import { clearSWCredentials, saveSWCredentials } from "$lib/sw-creds";
import { DateTime } from "luxon";

const deviceIDStorageKey = "e2clicker-device-id";
const webPushOptInStorageKey = "e2clicker-webpush-account";

export function deviceID() {
  let id = localStorage.getItem(deviceIDStorageKey);
  if (!id) {
    id = crypto.randomUUID().slice(0, 8);
    localStorage.setItem(deviceIDStorageKey, id);
  }
  return id;
}

export function clearPushOptIn() {
  if (typeof localStorage !== "undefined") {
    localStorage.removeItem(webPushOptInStorageKey);
  }
  clearSWCredentials().catch(() => {});
  if (typeof navigator !== "undefined" && "serviceWorker" in navigator) {
    navigator.serviceWorker
      .getRegistration()
      .then((reg) => reg?.pushManager.getSubscription())
      .then((sub) => sub?.unsubscribe())
      .catch(() => {});
  }
}

export async function pushIsEnabled() {
  const registration = await navigator.serviceWorker.ready;
  const subscription = await registration.pushManager.getSubscription();
  return !!subscription;
}

export enum NotificationError {
  Unsupported = "Notifications aren't supported in your browser.",
  PushMessaging = "Push messaging isn't supported in your browser.",
  Blocked = "You have blocked notifications for this app.",
  ServerUnsupported = "Server does not support this feature.",
}

export type UpdatePushSubscriptionResult = {
  enabled: boolean;
  available: boolean;
  reason?: NotificationError;
};

export async function updatePushSubscription({
  toggle,
  accountSecret,
  token,
}: {
  // Whether to enable or disable the subscription.
  // If not provided, the subscription is synced only when previously opted in for this account.
  toggle?: boolean;
  accountSecret?: string;
  token?: string;
} = {}): Promise<UpdatePushSubscriptionResult> {
  try {
    const supported = await isSupported();
    if (supported != true) {
      return {
        enabled: false,
        available: false,
        reason: supported,
      };
    }

    const webPushInfo = await api //
      .webPushInfo()
      .catch((err) => {
        if (api.isStatus(err, 400)) return null;
        throw err;
      });
    if (!webPushInfo) {
      // Server does not support handling Web Push.
      return {
        enabled: false,
        available: false,
        reason: NotificationError.ServerUnsupported,
      };
    }

    const registration = await navigator.serviceWorker.ready;
    let subscription = await registration.pushManager.getSubscription();

    if (toggle === false) {
      localStorage.removeItem(webPushOptInStorageKey);
      await clearSWCredentials().catch(() => {});
      if (subscription) {
        await subscription.unsubscribe().catch(() => {});
      }
      await api.userUnsubscribePush(deviceID()).catch((err) => {
        if (!api.isStatus(err, 404)) throw err;
      });
      return { enabled: false, available: true };
    }

    if (toggle === undefined) {
      const optedIn = localStorage.getItem(webPushOptInStorageKey);
      const sameAccount =
        !!optedIn && (!accountSecret || optedIn === accountSecret || optedIn === "true");
      if (!sameAccount) {
        return { enabled: false, available: true };
      }
      if (Notification.permission !== "granted") {
        localStorage.removeItem(webPushOptInStorageKey);
        await clearSWCredentials().catch(() => {});
        return { enabled: false, available: true, reason: NotificationError.Blocked };
      }
    }

    if (toggle === true) {
      const result = await Notification.requestPermission();
      if (result !== "granted") {
        return {
          enabled: false,
          available: true,
          reason: NotificationError.Blocked,
        };
      }
    }

    if (subscription && shouldResetSubscription(subscription, webPushInfo.applicationServerKey)) {
      await subscription.unsubscribe().catch(() => {});
      subscription = null;
    }

    if (!subscription) {
      subscription = await subscribeWithVAPIDKey(registration, webPushInfo.applicationServerKey);
    }

    await syncSubscription(subscription);

    localStorage.setItem(webPushOptInStorageKey, accountSecret ?? "true");
    if (token) {
      await saveSWCredentials({ token, deviceID: deviceID() }).catch(() => {});
    }

    return { enabled: true, available: true };
  } catch (err) {
    throw new Error(`Error updating push subscription`, { cause: err });
  }
}

function shouldResetSubscription(subscription: PushSubscription, expectedKeyBase64: string): boolean {
  if (subscription.expirationTime && subscription.expirationTime <= Date.now()) {
    return true;
  }

  const existingKey = subscription.options?.applicationServerKey;
  if (!existingKey) {
    return false;
  }

  const expected = base64UrlToUint8Array(expectedKeyBase64);
  const actual = new Uint8Array(existingKey);
  if (actual.length !== expected.length) {
    return true;
  }
  for (let i = 0; i < actual.length; i++) {
    if (actual[i] !== expected[i]) {
      return true;
    }
  }
  return false;
}

async function subscribeWithVAPIDKey(
  registration: ServiceWorkerRegistration,
  applicationServerKey: string,
): Promise<PushSubscription> {
  try {
    return await registration.pushManager.subscribe({
      userVisibleOnly: true,
      applicationServerKey,
    });
  } catch (err) {
    if (err instanceof DOMException && err.name === "InvalidStateError") {
      const existing = await registration.pushManager.getSubscription();
      if (existing) {
        await existing.unsubscribe().catch(() => {});
      }
      return await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey,
      });
    }
    throw err;
  }
}

function base64UrlToUint8Array(base64Url: string): Uint8Array {
  const padding = "=".repeat((4 - (base64Url.length % 4)) % 4);
  const base64 = (base64Url + padding).replaceAll("-", "+").replaceAll("_", "/");
  const raw = atob(base64);
  const output = new Uint8Array(raw.length);
  for (let i = 0; i < raw.length; i++) {
    output[i] = raw.charCodeAt(i);
  }
  return output;
}

async function syncSubscription(subscription: PushSubscription) {
  let expirationTime: string | undefined;
  if (subscription.expirationTime) {
    const t = DateTime.fromMillis(subscription.expirationTime);
    if (!t.isValid) {
      throw new Error("invalid expirationTime", { cause: t.invalidExplanation });
    }
    expirationTime = t.toISO();
  }

  const { p256dh, auth } = subscription.toJSON().keys as Record<string, string>;
  const id = deviceID();
  await api.userSubscribePush(id, {
    deviceID: id,
    endpoint: subscription.endpoint,
    expirationTime,
    keys: { p256dh, auth },
  });
}

async function isSupported(): Promise<true | NotificationError> {
  if (
    !("ServiceWorkerRegistration" in window) ||
    !("showNotification" in ServiceWorkerRegistration.prototype)
  ) {
    return NotificationError.Unsupported;
  }

  // Check if push messaging is supported
  if (!("PushManager" in window)) {
    return NotificationError.PushMessaging;
  }

  // Check the current Notification permission.
  // If its denied, it's a permanent block until the
  // user changes the permission
  if (Notification.permission === "denied") {
    return NotificationError.Blocked;
  }

  return true;
}
