import * as api from "$lib/api";
import { DateTime } from "luxon";

const deviceIDStorageKey = "e2clicker-device-id";

function deviceID() {
  let id = localStorage.getItem(deviceIDStorageKey);
  if (!id) {
    id = crypto.randomUUID().slice(0, 8);
    localStorage.setItem(deviceIDStorageKey, id);
  }
  return id;
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
}: {
  // Whether to enable or disable the subscription.
  // If not provided, the subscription is kept the same.
  toggle?: boolean;
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

    if (toggle == true) {
      const result = await Notification.requestPermission();
      if (result != "granted") {
        return {
          enabled: false,
          available: true,
          reason: NotificationError.Blocked,
        };
      }
    }

    const registration = await navigator.serviceWorker.ready;
    let subscription = await registration.pushManager.getSubscription();

    if (toggle == false && subscription) {
      await subscription.unsubscribe();
      await api.userUnsubscribePush(deviceID());
      return { enabled: false, available: true };
    }

    if (toggle == true && !subscription) {
      subscription = await registration.pushManager.subscribe({
        userVisibleOnly: true,
        applicationServerKey: webPushInfo.applicationServerKey,
      });
    }

    if (subscription) {
      await syncSubscription(subscription);
      return { enabled: true, available: true };
    }

    return { enabled: false, available: true };
  } catch (err) {
    throw new Error(`Error updating push subscription`, { cause: err });
  }
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
  await api.userSubscribePush({
    deviceID: deviceID(),
    endpoint: subscription.endpoint,
    expirationTime,
    keys: { p256dh, auth },
  });
}

async function isSupported(): Promise<true | NotificationError> {
  if (!("showNotification" in ServiceWorkerRegistration.prototype)) {
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
