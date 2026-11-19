<script lang="ts">
  import * as e2 from "#lib/e2.svelte.ts";
  import Icon from "#lib/components/Icon.svelte";
  import { DateTime } from "luxon";

  let {
    now,
    metrics,
  }: {
    now: DateTime;
    metrics: e2.LevelMetrics;
  } = $props();

  let trendLabel = $derived.by(() => {
    const trend = metrics.trend;
    if (!trend) return undefined;
    switch (trend.direction) {
      case "rising":
        return { icon: "trending-up", text: "Rising", highlight: true };
      case "falling":
        return { icon: "trending-down", text: "Falling", highlight: false };
      case "steady":
        return { icon: "trending-flat", text: "Steady", highlight: false };
    }
  });

  let peakEta = $derived(
    metrics.trend?.direction === "rising" && metrics.trend.peakAt
      ? e2.formatETA(metrics.trend.peakAt, now)
      : undefined,
  );

  let nextTroughEta = $derived(
    metrics.nextDoseAt && !metrics.isOverdue ? e2.formatETA(metrics.nextDoseAt, now) : undefined,
  );

  let troughDiffersFromSteady = $derived(
    metrics.nextTrough !== undefined &&
      metrics.ideal.trough > 0 &&
      Math.abs(metrics.nextTrough - metrics.ideal.trough) / metrics.ideal.trough > 0.05,
  );
</script>

<div class="header-row">
  <h3>Current Level</h3>
  {#if trendLabel}
    <span
      class="trend text-xs"
      class:secondary={trendLabel.highlight}
      class:muted={!trendLabel.highlight}
    >
      <Icon name={trendLabel.icon} />
      {trendLabel.text}
    </span>
  {/if}
</div>

{#if metrics.current !== undefined}
  <p class="current-levels text-3xl leading-8">
    <span class="text-nowrap">{e2.fixed(metrics.current)}</span>
    <span class="font-thin text-lg">{e2.units()}</span>
  </p>

  <div class="details text-sm muted">
    {#if metrics.trend?.direction === "rising"}
      {#if metrics.trend.peakValue !== undefined && peakEta}
        <p>
          Peaks ~<b>{e2.fixed(metrics.trend.peakValue)}</b>
          {e2.units()} in <b>{peakEta}</b>
        </p>
      {:else}
        <p>
          Rising ~<b>{e2.fixed(metrics.trend.slopePerDay)}</b>
          {e2.units()}/day
        </p>
      {/if}
    {:else if metrics.trend?.direction === "falling"}
      <p>
        Falling ~<b>{e2.fixed(Math.abs(metrics.trend.slopePerDay))}</b>
        {e2.units()}/day
      </p>
    {:else if metrics.trend}
      <p>Holding steady</p>
    {/if}

    {#if !metrics.isOverdue && metrics.nextTrough !== undefined}
      <p>
        Next trough: ~<b>{e2.fixed(metrics.nextTrough)}</b>
        {e2.units()}{#if nextTroughEta}
          {" "}in <b>{nextTroughEta}</b>{/if}
        {#if troughDiffersFromSteady}
          <span>(steady ~<b>{e2.fixed(metrics.ideal.trough)}</b>)</span>
        {/if}
      </p>
    {:else if metrics.isOverdue}
      {@const deficit = metrics.ideal.trough - metrics.current}
      {#if deficit > 2 * e2.conversionFactor()}
        <p>
          <b>{e2.fixed(deficit)}</b>
          {e2.units()} below steady trough
        </p>
      {:else}
        <p>Dose due now</p>
      {/if}
    {/if}
  </div>
{:else}
  <p class="current-levels empty">No doses logged</p>
  <p class="text-sm muted mb-0">
    Expected steady trough: ~<b>{e2.fixed(metrics.ideal.trough)}</b>
    {e2.units()}
  </p>
{/if}

<style lang="scss">
  .header-row {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    gap: 0.5rem;
  }

  .trend {
    display: inline-flex;
    align-items: center;
    gap: 0.25rem;
    font-weight: 600;
    white-space: nowrap;
  }

  .current-levels {
    margin-bottom: calc(var(--pico-spacing) / 2);

    &.empty {
      color: var(--pico-muted-color);
    }
  }

  .details p {
    margin: 0;
    line-height: 1.45;
  }
</style>
