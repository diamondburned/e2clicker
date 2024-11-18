import debounce_ from "debounce-promise";
import { untrack } from "svelte";
import { get, derived } from "svelte/store";
import { goto } from "$app/navigation";
import { DateTime } from "luxon";
import { persisted } from "svelte-persisted-store";

import * as api from "./api";
export * from "./api";

// token is the current session token.
// It is persisted to the local storage automatically.
const token = persisted<string | null>("e2clicker-token", null);

// isLoggedIn is true if the user is logged in.
export const isLoggedIn = derived(token, (token) => !!token);

token.subscribe((token) => {
  api.defaults.headers = { Authorization: token ? `Bearer ${token}` : undefined };
  api.defaults.baseUrl = "/api";
});

export const user = persisted<(api.User & { secret: api.UserSecret }) | null>("e2clicker-me", null);
let meLastUpdated: DateTime | undefined;

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
    meLastUpdated = DateTime.now();
  } catch (err) {
    // If unauthorized, then clear token and redirect to /login.
    if (api.isStatus(err, 401)) {
      token.set(null);
      goto("/login");
    }
  }
}

export async function auth(secret: string, { redirect = true }: { redirect?: boolean } = {}) {
  const resp = await api.auth({ secret });
  token.set(resp.token);
  await updateUser({ force: true });
  if (redirect) {
    goto("/");
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
