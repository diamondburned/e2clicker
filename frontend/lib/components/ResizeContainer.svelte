<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    children,
    anchor = "top",
  }: {
    children: Snippet;
    anchor?: "middle" | "top" | "bottom";
  } = $props();

  let innerHeight = $state(0);
  let outerHeight = $state(0);
  let style = $derived(`height: ${innerHeight}px`);
  let equal = $derived(innerHeight === outerHeight);
</script>

<div class="resize-container {anchor}" class:equal {style} bind:offsetHeight={outerHeight}>
  <div class="child" bind:offsetHeight={innerHeight}>{@render children()}</div>
</div>

<style lang="scss">
  .resize-container {
    transition: height var(--pico-transition);

    /* display: flex; */
    /* flex-direction: column; */

    display: block;
    overflow: hidden;

    &.middle {
      justify-content: center;
    }

    &.top {
      justify-content: flex-start;
    }

    &.bottom {
      justify-content: flex-end;
    }

    &.equal {
      overflow: visible;
    }

    .child {
      display: flex;
      flex-direction: column;
    }
  }
</style>
