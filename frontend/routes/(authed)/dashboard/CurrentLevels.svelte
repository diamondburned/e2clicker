<script lang="ts">
  import { Data as DosageData } from "./DosagePlot.svelte";
  import { fixed, units } from "$lib/e2.svelte";

  let {
    data,
  }: {
    data: DosageData;
  } = $props();

  let currentLevel = $derived(data.actualData.at(-1));
</script>

{#if currentLevel}
  <p class="current-levels text-3xl leading-8">
    <span class="text-nowrap">
      {fixed(currentLevel.value)}
    </span>
    <span class="font-thin">{units()}</span>
  </p>
  <!-- <p class="medication"> -->
  <!--   {delivery.name} -->
  <!-- </p> -->
{:else}
  <p class="current-levels empty">No data</p>
{/if}

<style lang="scss">
  .current-levels {
    &.empty {
      color: var(--pico-muted-color);
    }
  }
</style>
