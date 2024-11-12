<script lang="ts" generics="T">
  import { cubicInOut as easeFade } from "svelte/easing";
  import { fade, fly } from "svelte/transition";
  import ErrorBox from "./ErrorBox.svelte";
  import { onDestroy, type Snippet } from "svelte";

  let {
    main = false,
    mainClass,
    promise = new Promise<T>(() => {}),
    children,
    errorHeader,
    errorFooter,
  }: {
    main?: boolean;
    mainClass?: string;
    promise?: Promise<T>;
    children?: Snippet<[T]> | Snippet;
    errorHeader?: Snippet;
    errorFooter?: Snippet;
  } = $props();
</script>

{#await promise}
  <div
    class="loading-screen loading"
    aria-busy="true"
    out:fade={{ duration: 400, easing: easeFade }}
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
  <div class="loading-screen error" transition:fade={{ duration: 400, easing: easeFade }}>
    <article class="loading-error spaced">
      {#if errorHeader}
        {@render errorHeader?.()}
      {:else}
        <h3>ou nyow :(</h3>
      {/if}
      <ErrorBox {error} />
      {@render errorFooter?.()}
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

      opacity: 0;
      @keyframes fadeIn {
        0% {
          opacity: 0;
        }
        100% {
          opacity: 1;
        }
      }

      animation: fadeIn 150ms var(--pico-transition-easing);
      animation-delay: 500ms;
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
