<script lang="ts" module>
  // Hard-coded size of the popover arrow.
  const arrowSize = 8;

  // Enforce that there has to be at least 8px of margin between the popover
  // and the edge of the screen.
  const popoverEdgeMargin = 8;

  type Pt = { x: number; y: number };
  const zeroPt: Pt = { x: 0, y: 0 };

  type Size = { width: number; height: number };
  const zeroSize: Size = { width: 0, height: 0 };

  export type FlyProps = {
    x?: number;
    y?: number;
    duration?: number;
  };
</script>

<script lang="ts">
  import { onMount, tick, type Snippet } from "svelte";
  import { fly } from "svelte/transition";

  let {
    children,
    parent,

    direction: wantedDirection = "top",
    open = $bindable(true),
    fly: flyProps2,
  }: {
    children: Snippet;
    parent: HTMLElement;

    // direction is the direction that the popover will be shown.
    direction?: "top" | "bottom";
    // open is whether the popover is open.
    open?: boolean;
    // The settings for the fly transition.
    fly?: FlyProps;
  } = $props();

  let direction = $state(wantedDirection);

  let flyProps = $derived({
    x: flyProps2?.x,
    y: flyProps2?.y ?? (direction === "top" ? 10 : -10),
    duration: flyProps2?.duration ?? 200,
  });

  let winsize = $state(new DOMRect());
  onMount(() => {
    const parent = document.body;
    const update = () => {
      winsize = parent.getBoundingClientRect();
    };
    update();

    const observer = new ResizeObserver(update);
    observer.observe(parent);
    return () => observer.disconnect();
  });

  let popover = $state<HTMLDivElement | undefined>(undefined);
  let psize = $state(zeroSize); // uses bind:*
  let arrowX = $state(0);
  let offset = $state(zeroPt);

  $effect(() => {
    winsize;
    popover;
    open;

    if (!popover || !open) {
      direction = wantedDirection;
    }

    (async () => {
      await tick();
      if (popover && open) {
        const r = popover.getBoundingClientRect();
        offset = {
          x: Math.max(0, r.right - winsize.left - winsize.width + popoverEdgeMargin),
          y: Math.max(0, r.bottom - winsize.top - winsize.height + popoverEdgeMargin),
        };
      } else {
        offset = zeroPt;
      }

      await tick();
      if (popover && open) {
        const r = popover.getBoundingClientRect();
        direction =
          r.top < 0 // check for top overflow
            ? "bottom"
            : wantedDirection;
        arrowX = -popover.offsetLeft + parent.clientWidth / 2 + arrowSize / 2;
      } else {
        direction = wantedDirection;
        arrowX = 0;
      }
    })();
  });
</script>

<!--
  Close the popover if the user clicks outside of it.
  This only works if the popover also has an on:mousedown|stopPropagation
  so that it swallows mouse clicks inside the popover.
-->
<svelte:window
  onmousedown={(ev) => {
    ev.stopPropagation();
    if (open) {
      open = false;
    }
  }}
/>

{#if open}
  <div
    role="presentation"
    class="popover"
    onmousedown={(ev) => ev.stopPropagation()}
    bind:this={popover}
    bind:offsetWidth={psize.width}
    bind:offsetHeight={psize.height}
    class:popover-top={direction == "top"}
    class:popover-bottom={direction == "bottom"}
    style="
      --arrow-x: {arrowX}px;
      --arrow-size: {arrowSize}px;
      --width: {psize.width}px;
      --height: {psize.height}px;
      --overflow-x: {offset.x}px;
      --overflow-y: {offset.y}px;
    "
    transition:fly={flyProps}
  >
    {@render children()}
  </div>
{/if}

<style lang="scss">
  .popover {
    --pos-y: calc(-1 * var(--height) - var(--offset-y, 0px) - var(--arrow-y, 0));

    --arrow-x: 0;
    --arrow-y: calc(var(--arrow-size));
    --arrow-border-y: var(--arrow-size) solid var(--pico-dropdown-background-color);

    z-index: 10;

    border: 1px solid rgba(0 0 0 / 14%);
    background: var(--pico-dropdown-background-color);
    box-shadow:
      0 1px 5px 1px rgba(0 0 0 / 9%),
      0 2px 14px 3px rgba(0 0 0 / 5%);
    padding: var(--pico-spacing);
    border-radius: var(--pico-border-radius);

    width: max-content;
    min-width: 100px;
    max-width: min(250px, calc(100vw - 32px));

    font-size: 0.85em;

    &::after {
      content: " ";

      position: absolute;
      left: calc(var(--arrow-x) - 2px);

      border-left: var(--arrow-size) solid transparent;
      border-right: var(--arrow-size) solid transparent;
    }

    position: absolute;
    left: calc(50% - var(--width) / 2 - var(--overflow-x, 0px));

    &-top {
      top: var(--pos-y);

      &::after {
        bottom: calc(-1 * var(--arrow-y));
        border-top: var(--arrow-border-y);
      }
    }

    &-bottom {
      bottom: var(--pos-y);

      &::after {
        top: calc(-1 * var(--arrow-y));
        border-bottom: var(--arrow-border-y);
      }
    }
  }
</style>
