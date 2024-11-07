<script lang="ts">
  let { name }: { name: string } = $props();
</script>

<span class="material-symbols-rounded">{name.replaceAll("-", "_")}</span>

<style lang="scss">
  .material-symbols-rounded {
    font-variation-settings:
      "FILL" 1,
      "wght" 400,
      "GRAD" 0,
      "opsz" 24;
    vertical-align: text-bottom;
  }
</style>
