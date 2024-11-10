import type { Action } from "svelte/action";

import { Duration } from "luxon";
import parseDuration from "parse-duration";

export const daysInput: Action<
  HTMLInputElement,
  { initial?: number },
  {
    ondayschange: (e: CustomEvent<{ days: number }>) => void;
  }
> = (input, { initial }) => {
  function durationstr(duration: Duration): string {
    return duration.normalize().rescale().toHuman({ listStyle: "narrow" });
  }

  function onchange() {
    const parsed = parseDuration(input.value, "ms");
    if (!parsed) {
      input.value = "";

      input.setCustomValidity("Invalid duration.");
      input.reportValidity();
      return;
    }

    const duration = Duration.fromMillis(parsed);
    if (duration.as("hours") < 1 || duration.as("months") > 1) {
      input.value = "";

      input.setCustomValidity("Interval must be between 1 hour and 1 month.");
      input.reportValidity();
      return;
    }

    input.value = durationstr(duration);
    input.dispatchEvent(
      new CustomEvent("dayschange", {
        detail: { days: duration.as("days") },
      }),
    );
  }

  $effect(() => {
    if (initial && !input.value) {
      input.value = durationstr(Duration.fromObject({ days: initial }));
    }

    input.addEventListener("change", onchange);
    return () => {
      input.removeEventListener("change", onchange);
    };
  });
};

export const quantityInput: Action<
  HTMLInputElement,
  { unit: string; initial?: number },
  {
    onquantitychange: (e: CustomEvent<{ quantity: number }>) => void;
  }
> = (input, { unit, initial }) => {
  const re = `^([0-9]+(?:\\.[0-9]+)?) *(${unit})$`;

  input.pattern = re;
  input.placeholder = `0 ${unit}`;
  input.inputMode = "decimal";

  function onchange() {
    let match = input.value.match(re);
    if (!match) {
      // Check if the user didn't include the unit.
      const num = parseFloat(input.value);
      if (isNaN(num)) {
        input.value = "";
        return;
      }

      // Yes! We can manually add the unit.
      input.value = `${num} ${unit}`;
      match = input.value.match(re);
    }

    const quantity = parseFloat(match?.[1] ?? "");

    input.dispatchEvent(
      new CustomEvent("quantitychange", {
        detail: { quantity },
      }),
    );
  }

  $effect(() => {
    if (initial && !input.value) {
      input.value = `${initial} ${unit}`;
    }

    input.addEventListener("change", onchange);
    return () => {
      input.removeEventListener("change", onchange);
    };
  });
};
