<script lang="ts">
  import { type Snippet } from "svelte";
  import { fade, fly } from "svelte/transition";

  let {
    open = $bindable(false),
    wide = false,
    class: className,
    dismissible = false,
    children,
  }: {
    open: boolean;
    wide?: boolean;
    class?: string;
    dismissible?: boolean;
    children: Snippet;
  } = $props();
</script>

<svelte:window />

{#if open}
  <dialog
    open
    role="presentation"
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
      class="dialog {className}"
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
  .dialog-overlay {
    margin: 0;
  }

  .dialog {
    position: relative;
    max-width: var(--max-width, initial);

    & > :global(header) {
      --pico-block-spacing-vertical: var(--pico-spacing);

      padding-bottom: var(--pico-block-spacing-vertical);
    }
  }
</style>
