<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";
  import { Duration } from "luxon";
  import parseDuration from "parse-duration";

  let {
    initial,
    delimiter,
    placeholder,
    minDuration = Duration.fromObject({ hours: 1 }),
    maxDuration = Duration.fromObject({ months: 1 }),
    onchange,
    onchangemultiple,
    ...attributes
  }: {
    initial?: number | number[];
    delimiter?: string;
    placeholder?: string;
    minDuration?: Duration;
    maxDuration?: Duration;
    onchange?: (days: number | undefined) => void;
    onchangemultiple?: (days: number[] | undefined) => void;
  } & Omit<
    HTMLInputAttributes,
    "initial" | "delimiter" | "value" | "onchange" | "onchangemultiple"
  > = $props();

  let value = $state("");
  $effect(() => {
    if (initial) {
      value = (Array.isArray(initial) ? initial : [initial])
        .map((p) => daysstr(p))
        .join(delimiter ? `${delimiter.trim()} ` : "")
        .trim();
    }
  });

  function daysstr(days: number): string {
    return durationstr(Duration.fromObject({ days }));
  }

  function durationstr(duration: Duration): string {
    return duration.normalize().rescale().toHuman({ listStyle: "narrow" });
  }

  function setValidity(input: HTMLInputElement, message: string | null) {
    input.setCustomValidity(message || "");
    input.reportValidity();
  }

  function setInvalid(input: HTMLInputElement, message: string | null) {
    setValidity(input, message);
  }

  function setValid(input: HTMLInputElement) {
    setValidity(input, null);
  }
</script>

<input
  {...attributes}
  bind:value
  disabled={(!onchange && !onchangemultiple) || initial == undefined}
  onchange={(ev) => {
    const input = ev.currentTarget;
    const parts = delimiter ? input.value.split(delimiter).filter((v) => !!v) : [input.value];

    const durations = parts
      .map((part) => parseDuration(part.trim(), "ms"))
      .map((ms) => (ms ? Duration.fromMillis(ms) : undefined))
      .filter((d) => d != undefined)
      .filter((d) => d.isValid);
    console.debug({ durations, parts });
    if (durations.length != parts.length) {
      setInvalid(input, delimiter ? "Invalid durations." : "Invalid duration.");
      return;
    }

    if (durations.some((p) => p.toMillis() < minDuration.toMillis())) {
      setInvalid(input, `Interval must be at least ${minDuration.toHuman()}.`);
      return;
    }

    if (durations.some((p) => p.toMillis() > maxDuration.toMillis())) {
      setInvalid(input, `Interval must be at most ${maxDuration.toHuman()}.`);
      return;
    }

    value = durations.map(durationstr).join(delimiter);

    const days = durations.map((d) => d.as("days"));
    onchange?.(days.length > 0 ? days[0] : undefined);
    onchangemultiple?.(days);

    setValid(input);
  }}
/>
