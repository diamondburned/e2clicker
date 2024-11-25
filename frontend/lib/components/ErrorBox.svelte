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
    tiny?: boolean | "inline";
  } = $props();

  let errorThing = $derived.by(() => {
    const cleanError = (e: string) => e.replaceAll(/(^| )Error: /g, "$1").trim();
    switch (typeof error) {
      case "function":
        return error as Snippet;
      case "string":
        return cleanError(error);
      case "object":
        return cleanError(error.data?.message ?? error.message ?? `${error}`);
      default:
        return cleanError(`${error}`);
    }
  });

  let errorPrefix = $derived(prefix ? `${prefix}: ` : "");
</script>

{#snippet message()}
  {#if children && typeof errorThing == "string"}
    {@render children(errorThing)}
  {:else}
    {errorThing}
  {/if}
{/snippet}

{#if error}
  {#if tiny}
    <span class="error-span" transition:slide={{ duration: 200 }}>
      <span class="error-icon">
        <Icon name="error" />
      </span>
      {#if tiny == "inline"}
        <span class="error-text">
          {errorPrefix}
          {@render message()}
        </span>
      {:else}
        <Tooltip selectable>
          <span class="error-text">
            {prefix || "Error occured"}
          </span>
          {#snippet tooltip()}
            <b>{errorPrefix}</b>
            {@render message()}
          {/snippet}
        </Tooltip>
      {/if}
    </span>
  {:else}
    <blockquote
      class="error-box popping"
      class:tiny
      transition:slide={{ duration: 200 }}
      data-error-prefix={errorPrefix}
    >
      <span class="error-icon">
        <Icon name="error" />
      </span>
      <b class="error-text error-prefix">{errorPrefix}</b>
      <div class="error-text inline">{@render message()}</div>
    </blockquote>
  {/if}
{/if}

<style lang="scss" global>
  .error-box {
    --pico-primary: var(--pico-color-red);
  }

  .error-text,
  .error-icon {
    color: var(--pico-del-color);
  }
</style>
