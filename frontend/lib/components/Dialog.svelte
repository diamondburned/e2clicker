<script lang="ts">
  import { Snippet } from "svelte";
  import { fade } from "svelte/transition";

  let {
    open = $bindable(false),
    wide = false,
    dismissible = true,
    ondismiss = () => {
      open = false;
    },
    header,
    footer,
    content,
  }: {
    open: boolean;
    wide: boolean;
    dismissible: boolean;
    ondismiss: () => void;
    header?: Snippet;
    footer?: Snippet;
    content: Snippet;
  } = $props();
</script>

{#if open}
  <!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
  <div
    class="dialog-overlay"
    transition:fade
    onclick={() => {
      if (dismissible) {
        ondismiss();
      }
    }}
  >
    <dialog open class:wide></dialog>
  </div>
{/if}
