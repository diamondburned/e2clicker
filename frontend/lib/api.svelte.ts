import debounce_ from "debounce-promise";

export * from "./api";

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
  return ignoreFirst //
    ? debounce_(ignoreFirstRun(fn), wait, { leading })
    : debounce_(fn, wait, { leading });
}
