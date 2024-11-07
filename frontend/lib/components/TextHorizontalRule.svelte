<script lang="ts">
  import type { Snippet } from "svelte";

  let { children, thin = false }: { children: Snippet; thin?: boolean } = $props();
</script>

<div class="separator" role="separator" class:thin>
  {@render children()}
</div>

<style lang="scss">
  .separator {
    margin: calc(var(--pico-spacing) / 2) 0;

    &.thin {
      margin-top: 0;
      margin-bottom: 0;
    }

    display: flex;
    align-items: center;
    text-align: center;

    color: var(--pico-muted-color);

    &::before,
    &::after {
      content: "";
      flex: 1;
      border-bottom: var(--pico-border-width) solid var(--pico-muted-border-color);
    }

    &:not(:empty)::before {
      margin-right: calc(var(--pico-spacing));
    }

    &:not(:empty)::after {
      margin-left: calc(var(--pico-spacing));
    }
  }
</style>
