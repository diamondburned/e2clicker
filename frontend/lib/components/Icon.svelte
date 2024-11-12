<script lang="ts">
  let { name }: { name: string } = $props();
</script>

<span class="material-symbols-rounded">{name.replaceAll("-", "_")}</span>

<style lang="scss">
  .material-symbols-rounded {
    vertical-align: text-bottom;
  }
</style>
