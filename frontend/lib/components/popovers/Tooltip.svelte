<script lang="ts">
  import type { Snippet } from "svelte";

  import Popover from "./Popover.svelte";

  let {
    children,
    tooltip,
  }: {
    children: Snippet;
    tooltip: Snippet;
  } = $props();

  let parent = $state<HTMLElement | null>(null);
  let opened = $state(false);
  let closing = $state<ReturnType<typeof setTimeout> | null>(null);

  function open() {
    if (closing) {
      clearTimeout(closing);
    }
    opened = true;
  }

  function close() {
    closing = setTimeout(() => {
      opened = false;
      closing = null;
    }, 500);
  }
</script>

<div
  class="tooltip-container"
  role="presentation"
  onfocus={() => open()}
  onmouseover={() => open()}
  onmouseleave={() => close()}
>
  <span class="children" bind:this={parent}>
    {@render children()}
  </span>

  {#if parent}
    <Popover bind:open={opened} direction="top" {parent}>
      {@render tooltip()}
    </Popover>
  {/if}
</div>

<style lang="scss">
  .tooltip-container {
    position: relative;
    display: inline-block;
    color: var(--pico-contrast);

    &:hover {
      color: var(--pico-contrast-hover);
    }

    .children {
      cursor: help;
      font-size: 0.85em;

      text-decoration: underline dashed;
      text-decoration-color: var(--pico-muted-color);
    }
  }
</style>
