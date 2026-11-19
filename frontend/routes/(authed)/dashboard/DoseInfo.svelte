<script lang="ts">
  import * as e2 from "#lib/e2.svelte.ts";
  import * as api from "#lib/api.ts";
  import Icon from "#lib/components/Icon.svelte";
  import { Duration } from "luxon";

  let {
    dosage,
    metrics,
  }: {
    dosage?: api.Dosage;
    metrics?: e2.LevelMetrics;
  } = $props();

  let delivery = $derived(dosage && e2.deliveryMethod(dosage.deliveryMethod));
  let cadenceText = $derived.by(() => {
    if (!metrics?.cadence) return undefined;
    const hours = Math.max(1, Math.round(metrics.cadence.avgIntervalDays * 24));
    return e2.formatDuration(Duration.fromObject({ hours }));
  });
</script>

<div class="header-row">
  <h3>Your Regimen</h3>
  <a
    href="/settings"
    class="secondary settings-link"
    aria-label="Edit regimen in settings"
    title="Edit regimen"
  >
    <Icon name="tune" />
  </a>
</div>

{#if dosage && delivery}
  <p class="dosage text-3xl leading-8">
    <span class="text-nowrap">
      {dosage.dose ?? ""}
      {delivery.units ?? ""}
    </span>
    {#if delivery.patch && dosage.concurrence && dosage.concurrence > 1}
      <span class="text-lg font-thin text-nowrap">(×{dosage.concurrence} patches)</span>
    {/if}
    <span class="font-thin">every</span>
    <span class="text-nowrap">
      {e2.roundDuration(Duration.fromObject({ days: dosage.interval }).rescale()).toHuman()}
    </span>
  </p>

  <div class="details text-sm muted">
    <p>{delivery.name}</p>
    {#if cadenceText}
      <p>Actual rhythm: every {cadenceText}</p>
    {/if}
  </div>
{/if}

<style lang="scss">
  .header-row {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .settings-link {
    line-height: 1;
  }

  .dosage {
    margin-bottom: calc(var(--pico-spacing) / 2);
  }

  .details p {
    margin: 0;
    line-height: 1.45;
  }
</style>
