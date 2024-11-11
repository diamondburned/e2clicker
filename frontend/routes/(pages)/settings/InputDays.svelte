<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";
  import { Duration } from "luxon";
  import parseDuration from "parse-duration";

  let {
    initial,
    onchange,
    placeholder,
    ...attributes
  }: {
    initial?: number;
    placeholder?: string;
    onchange: (days: number) => void;
  } & Omit<HTMLInputAttributes, "value" | "onchange"> = $props();

  let value = $state(initial ? durationstr(Duration.fromObject({ days: initial })) : "");

  function durationstr(duration: Duration): string {
    return duration.normalize().rescale().toHuman({ listStyle: "narrow" });
  }
</script>

<input
  {...attributes}
  bind:value
  onchange={(ev) => {
    const input = ev.currentTarget;

    const parsed = parseDuration(input.value, "ms");
    if (!parsed) {
      value = "";

      input.setCustomValidity("Invalid duration.");
      input.reportValidity();
      return;
    }

    const duration = Duration.fromMillis(parsed);
    if (duration.as("hours") < 1 || duration.as("months") > 1) {
      value = "";

      input.setCustomValidity("Interval must be between 1 hour and 1 month.");
      input.reportValidity();
      return;
    }

    value = durationstr(duration);
    onchange(duration.as("days"));
  }}
/>
