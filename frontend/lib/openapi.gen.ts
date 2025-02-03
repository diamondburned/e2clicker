/**
 * e2clicker service
 * 0
 * DO NOT MODIFY - This file has been generated using oazapfts.
 * See https://www.npmjs.com/package/oazapfts
 */
import * as Oazapfts from "@oazapfts/runtime";
import * as QS from "@oazapfts/runtime/query";
export const defaults: Oazapfts.Defaults<Oazapfts.CustomHeaders> = {
    headers: {},
    baseUrl: "https://e2clicker.app/api",
};
const oazapfts = Oazapfts.runtime(defaults);
export const servers = {
    server1: "https://e2clicker.app/api",
    server2: "/api"
};
export type DeliveryMethod = {
    /** A short string representing the delivery method. This is what goes into the DeliveryMethod fields. */
    id: string;
    /** The units of the delivery method. */
    units: string;
    /** The full name of the delivery method. */
    name: string;
    /** A description of the delivery method. */
    description?: string;
};
export type Error = {
    /** A message describing the error */
    message: string;
    /** An array of errors that caused this error. If this is populated, then [details] is omitted. */
    errors?: Error[];
    /** Additional details about the error. Ignored if [errors] is used. */
    details?: any;
    /** Whether the error is internal */
    internal?: boolean;
    /** An internal code for the error (useless for clients) */
    internalCode?: string;
};
export type Dosage = {
    /** The delivery method to use. */
    deliveryMethod: string;
    /** The dosage amount. */
    dose: number;
    /** The interval between doses in days. */
    interval: number;
    /** The number of estrogen patches on the body at once. Only relevant if delivery method is patch. */
    concurrence?: number;
};
export type Dose = {
    /** The delivery method used. */
    deliveryMethod: string;
    /** The dosage amount. */
    dose: number;
    /** The time the dosage was taken. */
    takenAt: string;
    /** The time the dosage was taken off. This is only relevant for patch delivery methods. */
    takenOffAt?: string;
    /** A comment about the dosage, if any. */
    comment?: string;
};
export type DosageHistory = Dose[];
export type DosageHistoryCsv = string;
export type PushInfo = {
    /** A Base64-encoded string or ArrayBuffer containing an ECDSA P-256 public key that the push server will use to authenticate your application server. If specified, all messages from your application server must use the VAPID authentication scheme, and include a JWT signed with the corresponding private key. This key IS NOT the same ECDH key that you use to encrypt the data. For more information, see "Using VAPID with WebPush". */
    applicationServerKey: string;
};
export type PushDeviceId = string;
export type ReturnedPushSubscription = {
    deviceID: PushDeviceId;
    /** The time at which the subscription expires. This is the time when the subscription will be automatically deleted by the browser. */
    expirationTime?: string;
    keys: {
        /** An Elliptic curve Diffie–Hellman public key on the P-256 curve (that is, the NIST secp256r1 elliptic curve). The resulting key is an uncompressed point in ANSI X9.62 format. */
        p256dh: string;
    };
};
export type ReturnedNotificationMethods = {
    webPush?: ReturnedPushSubscription[];
};
export type PushSubscription = {
    deviceID: PushDeviceId;
    /** The endpoint to send the notification to. */
    endpoint: string;
    /** The time at which the subscription expires. This is the time when the subscription will be automatically deleted by the browser. */
    expirationTime?: string;
    /** The VAPID keys to encrypt the push notification. */
    keys: {
        /** An Elliptic curve Diffie–Hellman public key on the P-256 curve (that is, the NIST secp256r1 elliptic curve). The resulting key is an uncompressed point in ANSI X9.62 format. */
        p256dh: string;
        /** An authentication secret, as described in Message Encryption for Web Push. */
        auth: string;
    };
};
export type NotificationType = "welcome_message" | "reminder_message" | "account_notice_message" | "web_push_expiring_message";
export type NotificationMessage = {
    /** The title of the notification. */
    title: string;
    /** The message of the notification. */
    message: string;
};
export type Notification = {
    "type": NotificationType;
    /** The message of the notification. */
    message: NotificationMessage;
    /** The username of the user to send the notification to. */
    username: string;
};
export type Locale = string;
export type User = {
    /** The user's name */
    name: string;
    locale: Locale;
};
export type UserSecret = string;
export type Session = {
    /** The session identifier */
    id: number;
    /** The time the session was created */
    createdAt: string;
    /** The last time the session was used */
    lastUsed: string;
    /** The time the session expires, or null if it never expires */
    expiresAt?: string;
};
/**
 * List all available delivery methods
 */
export function deliveryMethods(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: DeliveryMethod[];
    } | {
        status: number;
        data: Error;
    }>("/delivery-methods", {
        ...opts
    }));
}
/**
 * Get the user's dosage and optionally their history
 */
export function dosage({ start, end }: {
    start?: string;
    end?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: {
            /** The user's current dosage schedule. This is null if the user has no dosage set. */
            dosage?: Dosage;
            /** The user's dosage history within the requested time range. If either historyStart or historyEnd are not provided, this will be null. */
            history?: DosageHistory;
        };
    }>(`/dosage${QS.query(QS.explode({
        start,
        end
    }))}`, {
        ...opts
    }));
}
/**
 * Set the user's dosage
 */
export function setDosage(dosage: Dosage, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>("/dosage", oazapfts.json({
        ...opts,
        method: "PUT",
        body: dosage
    })));
}
/**
 * Clear the user's dosage schedule
 */
export function clearDosage(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>("/dosage", {
        ...opts,
        method: "DELETE"
    }));
}
/**
 * Record a new dosage to the user's history
 */
export function recordDose(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: Dose;
    }>("/dosage/dose", {
        ...opts,
        method: "POST"
    }));
}
/**
 * Delete multiple dosages from the user's history
 */
export function forgetDoses(doseTimes: string[], opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>(`/dosage/dose${QS.query(QS.explode({
        doseTimes
    }))}`, {
        ...opts,
        method: "DELETE"
    }));
}
/**
 * Update a dosage in the user's history
 */
export function editDose(doseTime: string, body: Dose, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>(`/dosage/dose/${encodeURIComponent(doseTime)}`, oazapfts.json({
        ...opts,
        method: "PUT",
        body
    })));
}
/**
 * Delete a dosage from the user's history
 */
export function forgetDose(doseTime: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>(`/dosage/dose/${encodeURIComponent(doseTime)}`, {
        ...opts,
        method: "DELETE"
    }));
}
/**
 * Export the user's dosage history
 */
export function exportDoses(accept: "text/csv" | "application/json", { start, end }: {
    start?: string;
    end?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: DosageHistoryCsv;
    } | {
        status: 429;
        data: Error;
    } | {
        status: number;
        data: Error;
    }>(`/dosage/export-doses${QS.query(QS.explode({
        start,
        end
    }))}`, {
        ...opts,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            Accept: accept
        })
    }));
}
/**
 * Import a CSV file of dosage history
 */
export function importDoses(contentType: "text/csv" | "application/json", dosageHistoryCsv: DosageHistoryCsv, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: {
            /** The number of records in the file. */
            records: number;
            /** The number of records actually imported successfully. This is not equal to #records if there were errors or duplicate entries. */
            succeeded: number;
            error?: Error;
        };
    }>("/dosage/import-doses", oazapfts.json({
        ...opts,
        method: "POST",
        body: dosageHistoryCsv,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            "Content-Type": contentType
        })
    })));
}
/**
 * Get the server's push notification information
 */
export function webPushInfo(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: PushInfo;
    } | {
        status: number;
        data: Error;
    }>("/push-info", {
        ...opts
    }));
}
/**
 * Get the user's notification methods
 */
export function userNotificationMethods(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: ReturnedNotificationMethods;
    } | {
        status: number;
        data: Error;
    }>("/notifications/methods", {
        ...opts
    }));
}
/**
 * Get the user's push notification subscription
 */
export function userPushSubscription(deviceId: PushDeviceId, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: PushSubscription;
    } | {
        status: 404;
        data: Error;
    } | {
        status: number;
        data: Error;
    }>(`/notifications/methods/push/${encodeURIComponent(deviceId)}`, {
        ...opts
    }));
}
/**
 * Unsubscribe from push notifications
 */
export function userUnsubscribePush(deviceId: PushDeviceId, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>(`/notifications/methods/push/${encodeURIComponent(deviceId)}`, {
        ...opts,
        method: "DELETE"
    }));
}
/**
 * Create or update a push subscription
 */
export function userSubscribePush(pushSubscription: PushSubscription, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>("/notifications/methods/push", oazapfts.json({
        ...opts,
        method: "PUT",
        body: pushSubscription
    })));
}
export function getIgnoreNotificationHahaAnythingCanGoHereLol(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 500;
        data: Notification;
    }>("/_ignore/notification/_haha_anything_can_go_here_lol", {
        ...opts
    }));
}
/**
 * Register a new account
 */
export function register(body: {
    /** The name to register with */
    name: string;
}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: User & {
            secret: UserSecret;
        };
    } | {
        status: number;
        data: Error;
    }>("/register", oazapfts.json({
        ...opts,
        method: "POST",
        body
    })));
}
/**
 * Authenticate a user and obtain a session
 */
export function auth(body: {
    secret: UserSecret;
}, { userAgent }: {
    userAgent?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: {
            /** The session token */
            token: string;
        };
    } | {
        status: number;
        data: Error;
    }>("/auth", oazapfts.json({
        ...opts,
        method: "POST",
        body,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            "User-Agent": userAgent
        })
    })));
}
/**
 * Get the current user
 */
export function currentUser(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: User & {
            secret: UserSecret;
        };
    } | {
        status: number;
        data: Error;
    }>("/me", {
        ...opts
    }));
}
/**
 * List the current user's sessions
 */
export function currentUserSessions(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: Session[];
    } | {
        status: number;
        data: Error;
    }>("/me/sessions", {
        ...opts
    }));
}
/**
 * Delete one of the current user's sessions
 */
export function deleteUserSession(id: number, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>(`/me/sessions${QS.query(QS.explode({
        id
    }))}`, {
        ...opts,
        method: "DELETE"
    }));
}
