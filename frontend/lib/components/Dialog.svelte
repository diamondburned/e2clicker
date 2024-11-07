<script lang="ts">
  import { type Snippet } from "svelte";
  import { fade, fly } from "svelte/transition";

  let {
    open = $bindable(false),
    wide = false,
    dismissible = true,
    children,
  }: {
    open: boolean;
    wide?: boolean;
    dismissible?: boolean;
    children: Snippet;
  } = $props();
</script>

<svelte:window />

{#if open}
  <dialog
    role="presentation"
    open
    class="dialog-overlay"
    transition:fade={{ duration: 200 }}
    onmousedown={(ev) => {
      ev.stopPropagation();
      if (dismissible) {
        open = false;
      }
    }}
  >
    <article
      role="presentation"
      class="dialog"
      class:wide
      class:dismissible
      transition:fly={{ duration: 250, y: 50 }}
      onmousedown={(ev) => {
        ev.stopPropagation();
      }}
    >
      {@render children()}
    </article>
  </dialog>
{/if}

<style lang="scss">
  .dialog {
    position: relative;

    & > :global(header) {
      --pico-block-spacing-vertical: var(--pico-spacing);

      padding-bottom: var(--pico-block-spacing-vertical);
    }
  }
</style>
