<script lang="ts">
  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api";
  import { Duration } from "luxon";

  let {
    dosage,
  }: {
    dosage?: api.Dosage;
  } = $props();

  let delivery = $derived(dosage && e2.deliveryMethod(dosage.deliveryMethod));
</script>

<h3>Your Dosage</h3>

{#if dosage && delivery}
  <p class="dosage text-3xl leading-8">
    <span class="text-nowrap">
      {#if delivery.patch}
        1 patch
      {:else}
        {dosage.dose ?? ""}
        {delivery.units ?? ""}
      {/if}
    </span>
    <span class="font-thin">every</span>
    <span class="text-nowrap">
      {e2.roundDuration(Duration.fromObject({ days: dosage.interval }).rescale()).toHuman()}
    </span>
  </p>
  <p class="medication">
    {delivery.name}
  </p>
{/if}
