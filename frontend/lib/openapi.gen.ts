/**
 * REST API definition for the e2clicker service.
 * 0
 * DO NOT MODIFY - This file has been generated using oazapfts.
 * See https://www.npmjs.com/package/oazapfts
 */
import * as Oazapfts from "@oazapfts/runtime";
import * as QS from "@oazapfts/runtime/query";
export const defaults: Oazapfts.Defaults<Oazapfts.CustomHeaders> = {
    headers: {},
    baseUrl: "/api/user/v1",
};
const oazapfts = Oazapfts.runtime(defaults);
export const servers = {
    server1: "/api/user/v1"
};
export type SessionToken = string;
export type UserId = string;
export type User = {
    id: UserId;
    /** The user's email address */
    email: string;
    /** The user's name */
    name: string;
    /** The user's preferred locale */
    locale: string;
};
export function login(body?: {
    /** The username to log in with */
    email: string;
    /** The password to log in with */
    password: string;
}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: SessionToken;
    } | {
        status: 401;
    }>("/login", oazapfts.json({
        ...opts,
        method: "POST",
        body
    })));
}
export function register(body?: {
    /** The name to register with */
    name: string;
    /** The username to register with */
    email: string;
    /** The password to register with */
    password: string;
}, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 201;
        data: SessionToken;
    } | {
        status: 409;
    }>("/register", oazapfts.json({
        ...opts,
        method: "POST",
        body
    })));
}
export function user(userId: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchJson<{
        status: 200;
        data: User;
    }>(`/user/${encodeURIComponent(userId)}`, {
        ...opts
    }));
}
export function userAvatar(userId: string, opts?: Oazapfts.RequestOpts) {
    return oazapfts.ok(oazapfts.fetchBlob<{
        status: 200;
        data: Blob;
    }>(`/user/${encodeURIComponent(userId)}/avatar`, {
        ...opts
    }));
}
