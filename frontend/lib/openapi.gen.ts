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
export type UserId = string;
export type SessionToken = string;
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
export type Locale = string;
export type User = {
    id: UserId;
    /** The user's email address */
    email: string;
    /** The user's name */
    name: string;
    locale: Locale;
};
/**
 * Log into an existing account
 */
export function login(body?: {
    /** The username to log in with */
    email: string;
    /** The password to log in with */
    password: string;
}, { userAgent }: {
    userAgent?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: {
            userID: UserId;
            token: SessionToken;
        };
    } | {
        status: number;
        data: Error;
    }>("/login", oazapfts.json({
        ...opts,
        method: "POST",
        body,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            "User-Agent": userAgent
        })
    })));
}
/**
 * Register a new account
 */
export function register(body?: {
    /** The name to register with */
    name: string;
    /** The username to register with */
    email: string;
    /** The password to register with */
    password: string;
}, { userAgent }: {
    userAgent?: string;
} = {}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: {
            user: User;
            token: SessionToken;
        };
    } | {
        status: number;
        data: Error;
    }>("/register", oazapfts.json({
        ...opts,
        method: "POST",
        body,
        headers: oazapfts.mergeHeaders(opts?.headers, {
            "User-Agent": userAgent
        })
    })));
}
/**
 * Get a user by ID
 */
export function user(userId: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: User;
    } | {
        status: number;
        data: Error;
    }>(`/user/${encodeURIComponent(userId)}`, {
        ...opts
    }));
}
/**
 * Get a user's avatar by ID
 */
export function userAvatar(userId: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: Blob;
    } | {
        status: number;
        data: Error;
    }>(`/user/${encodeURIComponent(userId)}/avatar`, {
        ...opts
    }));
}
