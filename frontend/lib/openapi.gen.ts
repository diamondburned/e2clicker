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
    baseUrl: "/api",
};
const oazapfts = Oazapfts.runtime(defaults);
export const servers = {
    server1: "/api"
};
export type DeliveryMethod = {
    /** A short string representing the delivery method. This is what goes into the DeliveryMethod fields. */
    id: string;
    /** The units of the delivery method. */
    units: string;
    /** The full name of the delivery method. */
    name: string;
};
export type Error = {
    /** A message describing the error */
    message: string;
    /** Additional details about the error */
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
export type DosageObservation = {
    /** The unique identifier for the observation. */
    id: number;
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
export type DosageHistory = DosageObservation[];
export type Locale = string;
export type User = {
    /** The user's name */
    name: string;
    locale: Locale;
    /** Whether the user has an avatar. */
    hasAvatar: boolean;
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
export function dosage({ historyStart, historyEnd }: {
    historyStart?: string;
    historyEnd?: string;
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
        historyStart,
        historyEnd
    }))}`, {
        ...opts
    }));
}
/**
 * Set the user's dosage
 */
export function setDosage(dosage?: Dosage, opts?: Oazapfts.RequestOpts) {
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
export function recordDose(body?: {
    /** The time the dosage was taken. */
    takenAt: string;
}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: DosageObservation;
    }>("/dosage/dose", oazapfts.json({
        ...opts,
        method: "POST",
        body
    })));
}
/**
 * Update a dosage in the user's history
 */
export function editDose(body?: DosageObservation, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>("/dosage/dose", oazapfts.json({
        ...opts,
        method: "PUT",
        body
    })));
}
/**
 * Delete multiple dosages from the user's history
 */
export function forgetDoses(doseIds: number[], opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>(`/dosage/dose${QS.query(QS.explode({
        dose_ids: doseIds
    }))}`, {
        ...opts,
        method: "DELETE"
    }));
}
/**
 * Register a new account
 */
export function register(body?: {
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
export function auth(body?: {
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
        data: User;
    } | {
        status: number;
        data: Error;
    }>("/me", {
        ...opts
    }));
}
/**
 * Get the current user's avatar
 */
export function currentUserAvatar(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: Blob;
    } | {
        status: number;
        data: Error;
    }>("/me/avatar", {
        ...opts
    }));
}
/**
 * Set the current user's avatar
 */
export function setCurrentUserAvatar(body?: Blob, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>("/me/avatar", {
        ...opts,
        method: "PUT",
        body
    }));
}
/**
 * Get the current user's secret
 */
export function currentUserSecret(opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: {
            secret: UserSecret;
        };
    } | {
        status: number;
        data: Error;
    }>("/me/secret", {
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
export function deleteUserSession(body?: {
    /** The session identifier to delete */
    id: number;
}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 204;
    } | {
        status: number;
        data: Error;
    }>("/me/sessions", oazapfts.json({
        ...opts,
        method: "DELETE",
        body
    })));
}
