<script lang="ts">
  let { children } = $props();
</script>

<div class="separator" role="separator">
  {@render children()}
</div>

<style lang="scss">
  .separator {
    margin: calc(var(--pico-spacing) / 2) 0;

    display: flex;
    align-items: center;
    text-align: center;

    &::before,
    &::after {
      content: "";
      flex: 1;
      border-bottom: 1px solid var(--pico-muted-color);
    }

    &:not(:empty)::before {
      margin-right: calc(var(--pico-spacing));
    }

    &:not(:empty)::after {
      margin-left: calc(var(--pico-spacing));
    }
  }
</style>
