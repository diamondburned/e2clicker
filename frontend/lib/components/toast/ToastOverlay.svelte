<script lang="ts">
  import { toasts } from "$lib/toasts";
  import Toast from "./Toast.svelte";
</script>

<div class="toast-overlay">
  {#each $toasts as toast}
    <Toast {toast} />
  {/each}
</div>

<style lang="scss">
  .toast-overlay {
    position: fixed;
    bottom: 0;
    width: 100%;

    padding: var(--pico-spacing);
    padding-bottom: 0;

    display: flex;
    flex-direction: column;
    align-items: center;

    z-index: 900;

    pointer-events: none;

    & > :global .toast {
      pointer-events: auto;
    }

    &:not(:empty) {
      background-image: linear-gradient(
        to top,
        color-mix(in srgb, var(--pico-primary-background), transparent 50%),
        transparent
      );
    }
  }
</style>
