<script lang="ts">
  import { cubicInOut as easeFade } from "svelte/easing";
  import { fade } from "svelte/transition";

  let {
    promise = Promise.resolve(),
  }: {
    promise?: Promise<unknown>;
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
{:catch error}
  <div class="loading-screen error" transition:fade={{ duration: 350, easing: easeFade }}>
    <article class="loading-error spaced">
      <h3>ou nyow :(</h3>
      <pre>{error.message}</pre>
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
    max-width: 400px;

    padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);

    h3 {
      user-select: none;
    }

    pre {
      background: none;
      border-radius: 0;
    }
  }
</style>
