<script lang="ts">
  import { slide } from "svelte/transition";
  import Tooltip from "./popovers/Tooltip.svelte";
  import Icon from "./Icon.svelte";
  import type { Snippet } from "svelte";

  let {
    error,
    children,
    prefix,
    tiny = false,
  }: {
    error: any | Snippet;
    children?: Snippet<[string]>;
    prefix?: string;
    tiny?: boolean;
  } = $props();

  let errorString = $derived.by(() => {
    switch (typeof error) {
      case "function":
        return error as Snippet;
      case "string":
        return error;
      case "object":
        return error.data?.message ?? error.message ?? `${error}`;
      default:
        return `${error}`;
    }
  });

  let errorPrefix = $derived(prefix ? `Error: ${prefix}: ` : "Error: ");
</script>

{#snippet message()}
  {#if children}
    {@render children(errorString)}
  {:else}
    {errorString}
  {/if}
{/snippet}

{#if error}
  {#if tiny}
    <span class="error-span" transition:slide={{ duration: 200 }}>
      <Tooltip selectable>
        <span class="error-text">
          {prefix || "Error occured"}
        </span>
        {#snippet tooltip()}
          <b>{errorPrefix}</b>
          {@render message()}
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
      <span class="error-text">{@render message()}</span>
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
