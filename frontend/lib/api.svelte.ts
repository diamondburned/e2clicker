import debounce_ from "debounce-promise";
import { untrack } from "svelte";
import { get, derived } from "svelte/store";
import { goto } from "$app/navigation";
import { DateTime } from "luxon";
import { persisted } from "svelte-persisted-store";

import * as api from "./api";
import { clearPushOptIn } from "./notification/webpush";
export * from "./api";

// token is the current session token.
// It is persisted to the local storage automatically.
const token = persisted<string | null>("e2clicker-token", null);

// isLoggedIn is true if the user is logged in.
export const isLoggedIn = derived(token, (token) => !!token);

export function currentToken(): string | null {
  return get(token);
}

const fetchAPI: Exclude<typeof api.defaults.fetch, undefined> = async (input, init = {}) => {
  const resp = await window.fetch("/api" + input, init);
  switch (resp.status) {
    case 401:
      // Unauthorized. Clear token, push opt-in, and redirect to /login.
      clearSession();
      if (typeof window !== "undefined" && window.location.pathname !== "/login") {
        goto("/login");
      }
      break;
    case 500:
      console.error("500 Internal Server Error", resp);
      break;
  }
  return resp;
};

api.defaults.fetch = fetchAPI;
api.defaults.baseUrl = "";

token.subscribe((tok) => {
  api.defaults.headers = { Authorization: tok ? `Bearer ${tok}` : undefined };
});

function clearSession() {
  token.set(null);
  clearPushOptIn();
}

// fetch performs an API request with the given path and parameters.
// The path is relative to the API base URL.
export async function fetch(
  path: string,
  params: Omit<RequestInit, "headers"> & {
    headers?: Record<string, string>;
  },
) {
  const resp = await fetchAPI(path, {
    ...params,
    headers: {
      ...(params.headers ?? {}),
      ...(api.defaults.headers as Record<string, string>),
    },
  });
  if (!resp.ok) {
    let err: api.Error;
    try {
      err = (await resp.json()) as api.Error;
    } catch (e) {
      throw new Error(`HTTP ${resp.status} ${resp.statusText}`);
    }
    throw {
      status: resp.status,
      data: err,
    };
  }
  return resp;
}

export const user = persisted<(api.User & { secret: api.UserSecret }) | null>("e2clicker-me", null);
// savedSecret persists the account secret independently of session state so
// the login field stays populated if a session token is invalidated.
export const savedSecret = persisted<api.UserSecret | null>("e2clicker-secret", null);
let meLastUpdated: DateTime | undefined;

if (!get(savedSecret)) {
  const existingSecret = get(user)?.secret;
  if (existingSecret) {
    savedSecret.set(existingSecret);
  }
}

export async function updateUser({ force = false }: { force?: boolean } = {}) {
  if (!get(isLoggedIn)) {
    return;
  }

  if (!force && meLastUpdated && meLastUpdated.diffNow().as("minutes") < 1) {
    // Too soon to update again.
    return;
  }

  try {
    const v = await api.currentUser();
    user.set(v);
    savedSecret.set(v.secret);
    meLastUpdated = DateTime.now();
  } catch (err) {
    // If unauthorized, then clear token and redirect to /login.
    if (api.isStatus(err, 401)) {
      clearSession();
      goto("/login");
    }
  }
}

export async function auth(secret: string, { redirect = true }: { redirect?: boolean } = {}) {
  const previousSecret = get(savedSecret) ?? get(user)?.secret;
  if (previousSecret && previousSecret !== secret) {
    clearPushOptIn();
  }
  const resp = await api.auth({ secret });
  savedSecret.set(secret);
  token.set(resp.token);
  await updateUser({ force: true });
  if (redirect) {
    goto("/dashboard");
  }
}

export function ignoreFirstRun<T extends (...args: any[]) => any>(fn: T) {
  let first = true;
  return (...args: Parameters<T>) => {
    if (first) {
      console.debug(
        "Ignoring first run of function",
        args.map((a) => $state.snapshot(a)),
      );
      first = false;
      return Promise.resolve();
    }
    console.debug(
      "Running function",
      args.map((a) => $state.snapshot(a)),
    );
    return fn(...args);
  };
}

export function debounced<T extends (...args: any[]) => any>(
  fn: T,
  {
    wait = 500,
    leading = true,
    ignoreFirstRun: ignoreFirst = false,
  }: {
    wait?: number;
    ignoreFirstRun?: boolean;
  } & debounce_.DebounceOptions = {},
) {
  const fn2 = debounce_(fn, wait, { leading });
  return ignoreFirst ? ignoreFirstRun(fn2) : fn2;
}

export type AnyAsyncToOK = AsyncToOK<(...args: any[]) => Promise<any>, any>;

export class AsyncToOK<
  T extends (...args: any[]) => Promise<any>,
  ValueT = Awaited<ReturnType<T>>,
> {
  value = $state<ValueT | undefined>(undefined);
  error = $state<Error | undefined>(undefined);
  loading = $state(true);
  promise = $state(new Promise<ValueT>(() => {}));

  constructor(
    private fn: T,
    private opts: {
      firstPromiseOnly?: boolean;
      debounce?: number | boolean;
      initial?: ValueT;
    } = {},
  ) {
    if (this.opts.initial !== undefined) {
      this.loading = false;
      this.promise = Promise.resolve(this.opts.initial);
      this.value = this.opts.initial;
    }

    if (this.opts.debounce) {
      this.fn = debounced(
        this.fn,
        typeof this.opts.debounce == "number" ? { wait: this.opts.debounce } : {},
      ) as T;
    }
  }

  public do = this.load; // alias
  public load(...args: Parameters<T>) {
    this.load_(...args);
  }

  private async load_(...args: Parameters<T>) {
    try {
      this.error = undefined;
      this.loading = true;

      const p = this.fn(...args);
      if (!(this.opts.firstPromiseOnly ?? false) || untrack(() => this.value) === undefined) {
        this.promise = p;
      }
      this.value = await p;
    } catch (err) {
      this.error = err instanceof Error ? err : new Error(`${err}`);
    } finally {
      this.loading = false;
    }
  }
}

// discard wraps fn for a function that returns void.
export function discard<
  T extends (...args: any[]) => any,
  ReturnT = ReturnType<T> extends Promise<any> ? Promise<void> : void,
>(fn: T): (...args: Parameters<T>) => ReturnT {
  return fn;
}
