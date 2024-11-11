<script lang="ts" generics="T">
  import { cubicInOut as easeFade } from "svelte/easing";
  import { fade, fly } from "svelte/transition";
  import ErrorBox from "./ErrorBox.svelte";
  import type { Snippet } from "svelte";

  let {
    main = false,
    mainClass,
    promise = new Promise<T>(() => {}),
    children,
  }: {
    main?: boolean;
    mainClass?: string;
    promise?: Promise<T>;
    children?: Snippet<[T]> | Snippet;
  } = $props();
</script>

{#await promise}
  <div
    class="loading-screen loading"
    aria-busy="true"
    transition:fade={{ duration: 350, easing: easeFade }}
  >
    Loading...
  </div>
{:then value}
  {#if children}
    {#if main}
      <main class="container {mainClass ?? ''}" transition:fly={{ duration: 200, y: 100 }}>
        {@render children(value)}
      </main>
    {:else}
      {@render children(value)}
    {/if}
  {/if}
{:catch error}
  <div class="loading-screen error" transition:fade={{ duration: 350, easing: easeFade }}>
    <article class="loading-error spaced">
      <h3>ou nyow :(</h3>
      <ErrorBox {error} />
    </article>
  </div>
{/await}

<style lang="scss">
  .loading-screen {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    z-index: 1000;

    color: var(--pico-color);
    background-color: var(--pico-modal-overlay-background-color);

    display: flex;
    justify-content: center;
    align-items: center;

    user-select: none;

    &.loading {
      cursor: wait;
    }

    &.error {
      cursor: not-allowed;
    }
  }

  .loading-error {
    cursor: initial;

    width: 100%;
    max-width: 600px;

    padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);

    h3 {
      user-select: none;
    }
  }
</style>
