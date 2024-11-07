<script lang="ts">
  import { slide } from "svelte/transition";

  let { error, prefix }: { error: any; prefix?: string } = $props();

  let message = $derived(
    typeof error == "object" && !!error.data?.message ? error.data.message : `${error}`,
  );
</script>

{#if error}
  <blockquote
    class="error"
    transition:slide={{ duration: 200 }}
    data-error-prefix={prefix ? `Error: ${prefix}: ` : "Error: "}
  >
    {message}
  </blockquote>
{/if}

<style lang="scss">
  .error {
    border: var(--pico-border-width) solid var(--pico-color-red);
    border-radius: var(--pico-border-radius);

    color: var(--pico-del-color);
    background-color: color-mix(in srgb, var(--pico-color-red), transparent 90%);

    &::before {
      content: attr(data-error-prefix);
      font-weight: bold;
    }
  }
</style>
