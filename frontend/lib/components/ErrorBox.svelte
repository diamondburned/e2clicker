<script lang="ts">
  import { slide } from "svelte/transition";
  import Tooltip from "./popovers/Tooltip.svelte";
  import Icon from "./Icon.svelte";

  let {
    error,
    prefix,
    tiny = false,
  }: {
    error: any;
    prefix?: string;
    tiny?: boolean;
  } = $props();

  let message = $derived(
    typeof error == "object"
      ? !!error.data?.message //
        ? error.data.message
        : error.message
      : `${error}`,
  );

  let errorPrefix = $derived(prefix ? `Error: ${prefix}: ` : "Error: ");
</script>

{#if error}
  {#if tiny}
    <span class="error-span" transition:slide={{ duration: 200 }}>
      <Tooltip selectable>
        <span class="error-text">
          {prefix || "Error occured"}
        </span>
        {#snippet tooltip()}
          <b>{errorPrefix}</b>{message}
        {/snippet}
      </Tooltip>
      <span class="error-icon">
        <Icon name="error" />
      </span>
    </span>
  {:else}
    <blockquote
      class="error-box"
      class:tiny
      transition:slide={{ duration: 200 }}
      data-error-prefix={errorPrefix}
    >
      <span class="error-icon">
        <Icon name="error" />
      </span>
      <b class="error-text error-prefix">{errorPrefix}</b>
      <span class="error-text">{message}</span>
    </blockquote>
  {/if}
{/if}

<style lang="scss" global>
  .error-box {
    border: var(--pico-border-width) solid var(--pico-color-red);
    border-radius: var(--pico-border-radius);
    background-color: color-mix(in srgb, var(--pico-color-red), transparent 90%);
  }

  .error-text,
  .error-icon {
    color: var(--pico-del-color);
  }
</style>
