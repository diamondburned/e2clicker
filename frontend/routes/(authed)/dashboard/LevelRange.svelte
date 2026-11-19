<script lang="ts">
  import * as e2 from "#lib/e2.svelte.ts";
  import * as api from "#lib/api.ts";
  import { Duration } from "luxon";

  let {
    metrics,
    dosage,
  }: {
    metrics: e2.LevelMetrics;
    dosage: api.Dosage;
  } = $props();

  let title = $derived.by(() => {
    if (!metrics.window) {
      return "Steady-State Range";
    }
    if (metrics.window.clampedToFirstDose) {
      return "Since First Dose";
    }
    if (dosage.interval <= 7) {
      return "Last 7 Days";
    }
    const label = e2
      .roundDuration(Duration.fromObject({ days: dosage.interval }).rescale())
      .toHuman();
    return `Last ${label}`;
  });

  let displayTrough = $derived(metrics.window ? metrics.window.trough : metrics.ideal.trough);
  let displayPeak = $derived(metrics.window ? metrics.window.peak : metrics.ideal.peak);
</script>

<h3>{title}</h3>

<div class="range-grid">
  <div class="stat">
    <span class="stat-label">Trough</span>
    <p class="stat-value text-3xl leading-8">
      <span class="text-nowrap">{e2.fixed(displayTrough)}</span>
      <span class="font-thin text-lg">{e2.units()}</span>
    </p>
  </div>

  <div class="stat">
    <span class="stat-label">Peak</span>
    <p class="stat-value text-3xl leading-8">
      <span class="text-nowrap">{e2.fixed(displayPeak)}</span>
      <span class="font-thin text-lg">{e2.units()}</span>
    </p>
  </div>
</div>

<p class="sub-line text-sm">
  {#if metrics.window}
    Steady state: {e2.fixed(metrics.ideal.trough)}–{e2.fixed(metrics.ideal.peak)} {e2.units()}
  {:else}
    Projected cycle — log a dose to track actual levels
  {/if}
</p>

<style lang="scss">
  .range-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: var(--pico-spacing);
    margin-bottom: calc(var(--pico-spacing) / 2);
  }

  .stat {
    display: flex;
    flex-direction: column;
    gap: 0.15rem;
  }

  .stat-label {
    font-size: var(--font-size-xs);
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .stat-value {
    margin: 0;
  }

  .sub-line {
    margin: 0;
  }
</style>
