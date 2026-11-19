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
  let lastPointerType = "mouse";

  function open() {
    if (closing) {
      clearTimeout(closing);
      closing = null;
    }
    opened = true;
  }

  function close() {
    if (closing) {
      clearTimeout(closing);
    }
    closing = setTimeout(() => {
      opened = false;
      closing = null;
    }, timeout);
  }

  function closeImmediately() {
    if (closing) {
      clearTimeout(closing);
      closing = null;
    }
    opened = false;
  }
</script>

<svelte:window
  onkeydown={() => {
    lastPointerType = "keyboard";
  }}
/>

<div
  class="tooltip-container"
  class:selectable
  role="presentation"
  onpointerdown={(ev) => {
    lastPointerType = ev.pointerType;
  }}
  onfocusin={() => selectable && lastPointerType !== "touch" && open()}
  onfocusout={() => selectable && close()}
  onpointerenter={(ev) => selectable && ev.pointerType !== "touch" && open()}
  onpointerleave={(ev) => selectable && ev.pointerType !== "touch" && close()}
>
  <span
    class="children"
    class:has-tooltip={!!tooltip}
    bind:this={parent}
    role="presentation"
    onmousedown={(ev) => {
      if (tooltip) {
        ev.stopPropagation();
      }
    }}
    onclick={() => {
      if (!tooltip) return;
      if (lastPointerType === "touch") {
        if (opened) {
          closeImmediately();
        } else {
          open();
        }
      } else {
        open();
      }
    }}
    onfocusin={() => !selectable && lastPointerType !== "touch" && open()}
    onfocusout={() => !selectable && close()}
    onpointerenter={(ev) => !selectable && ev.pointerType !== "touch" && open()}
    onpointerleave={(ev) => !selectable && ev.pointerType !== "touch" && close()}
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
