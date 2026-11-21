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
    onchange: (quantity: number | undefined) => void;
  } & Omit<
    HTMLInputAttributes,
    "pattern" | "value" | "inputmode" | "placeholder" | "onchange"
  > = $props();

  const pattern = $derived(`^([0-9]+(?:\\.[0-9]+)?) *(${unit})$`);

  function formatQuantity(qty: number | undefined, currentUnit: string) {
    return qty !== undefined && !isNaN(qty) && qty > 0 ? `${qty} ${currentUnit}` : "";
  }

  let inputEl: HTMLInputElement | null = null;
  let value = $state("");
  $effect(() => {
    value = formatQuantity(initial, unit);
    inputEl?.setCustomValidity("");
  });
</script>

<input
  bind:this={inputEl}
  {pattern}
  {...attributes}
  bind:value
  inputmode="decimal"
  placeholder="{placeholder || initial || 0} {unit}"
  onchange={(ev) => {
    const input = ev.currentTarget;
    let qty: number | undefined;

    const match = input.value.trim().match(pattern);
    if (match) {
      qty = parseFloat(match[1] ?? "");
    } else {
      // Check if the user didn't include the unit.
      const parsed = parseFloat(input.value);
      if (isNaN(parsed) || parsed <= 0) {
        value = "";
        qty = undefined;
      } else {
        // Yes! We can manually add the unit.
        qty = parsed;
        value = `${qty} ${unit}`;
      }
    }

    if (qty !== undefined && qty <= 0) {
      input.setCustomValidity("Dose must be greater than 0.");
      input.reportValidity();
      qty = undefined;
    } else {
      input.setCustomValidity("");
    }

    onquantitychange(qty);
  }}
/>
