<script lang="ts">
  import type { Snippet } from "svelte";

  import Popover from "./Popover.svelte";

  let {
    children,
    tooltip,
    timeout = 200,
    selectable = false,
  }: {
    children: Snippet;
    tooltip?: Snippet;
    timeout?: number;
    selectable?: boolean;
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
    }, timeout);
  }
</script>

<div
  class="tooltip-container"
  class:selectable
  role="presentation"
  onfocus={() => selectable && open()}
  onmouseover={() => selectable && open()}
  onmouseleave={() => selectable && close()}
>
  <span
    class="children"
    class:has-tooltip={!!tooltip}
    bind:this={parent}
    role="presentation"
    onfocus={() => !selectable && open()}
    onmouseover={() => !selectable && open()}
    onmouseleave={() => !selectable && close()}
  >
    {@render children()}
  </span>

  {#if parent && tooltip}
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

    .children.has-tooltip {
      cursor: help;

      text-decoration: underline dashed;
      text-decoration-color: var(--pico-muted-color);
    }

    &:not(.selectable) :global .popover {
      cursor: default;
      pointer-events: none;
    }
  }
</style>
