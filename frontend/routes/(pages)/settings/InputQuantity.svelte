<script lang="ts">
  import type { HTMLInputAttributes } from "svelte/elements";

  let {
    unit,
    initial,
    onchange: onquantitychange,
    placeholder,
    ...attributes
  }: {
    unit: string;
    initial?: number;
    placeholder?: number;
    onchange: (quantity: number) => void;
  } & Omit<
    HTMLInputAttributes,
    "pattern" | "value" | "inputmode" | "placeholder" | "onchange"
  > = $props();

  const pattern = $derived(`^([0-9]+(?:\\.[0-9]+)?) *(${unit})$`);
  let value = $state(initial ? `${initial} ${unit}` : "");
</script>

<input
  {pattern}
  {...attributes}
  bind:value
  inputmode="decimal"
  placeholder="{placeholder || initial || 0} {unit}"
  onchange={(ev) => {
    const input = ev.currentTarget;
    let qty: number;

    let match = input.value.match(pattern);
    if (match) {
      qty = parseFloat(match?.[1] ?? "");
    } else {
      // Check if the user didn't include the unit.
      qty = parseFloat(input.value);
      if (isNaN(qty)) {
        value = "";
        return;
      }
      // Yes! We can manually add the unit.
      value = `${qty} ${unit}`;
    }

    onquantitychange(qty);
  }}
/>
