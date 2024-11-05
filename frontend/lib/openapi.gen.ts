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
export type Locale = string;
export type User = {
    /** The user's name */
    name: string;
    locale: Locale;
    /** Whether the user has an avatar. */
    has_avatar: boolean;
};
export type UserSecret = string;
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
export type Session = {
    /** The session identifier */
    id: number;
    /** The time the session was created */
    created_at: string;
    /** The last time the session was used */
    last_used: string;
    /** The time the session expires, or null if it never expires */
    expires_at?: string;
};
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
        method: "POST",
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
