<!-- estrannai.se as a svelte component -->

<script lang="ts" module>
</script>

<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";

  import type * as charts from "lightweight-charts";
  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api";
  import { createChart, LineStyle, LastPriceAnimationMode } from "lightweight-charts";
  import { Interval } from "luxon";
  import { fade } from "svelte/transition";

  let {
    interval,
    doses,
  }: {
    interval: Interval<true>;
    doses: {
      dosage?: api.Dosage;
      history?: e2.DosageHistory;
    };
  } = $props();

  let dosage = $derived(doses?.dosage);
  let history = $derived(doses?.history ?? []);

  let predictedInterval = $derived(
    Interval.fromDateTimes(
      // Predict ahead by the configured interval.
      interval.end,
      interval.end.plus(e2.predictAhead()),
    ) as Interval<true>,
  );
  const totalInterval = $derived(interval.union(predictedInterval) as Interval<true>);

  let wpathRange = $derived(e2.wpathRange(e2.conversionFactor()));
  // let idealTrough = $derived(dosage && e2.idealE2Trough(dosage, e2.conversionFactor()));
  // let idealAverage = $derived(dosage && e2.idealE2Average(dosage, e2.conversionFactor()));

  let plotDiv: HTMLDivElement | null = null;
  let plotWidth = $state(0);

  let styles = $derived(e2.gatherStyles(plotDiv));

  let plotTooltipDiv: HTMLDivElement | null = null;
  let plotTooltip = $state<e2.PlotTooltipData | null>(null);

  let chart = $derived.by(() => {
    if (!plotDiv) {
      return;
    }

    const chart = createChart(plotDiv);

    const wpathFake = chart.addLineSeries({
      visible: true,
    }) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

    const wpathLower = wpathFake.createPriceLine({
      price: 0,
      color: styles.muted,
      lineWidth: 2,
      lineStyle: LineStyle.Dashed,
      axisLabelVisible: true,
      title: "Lower Target",
    });

    const wpathUpper = wpathFake.createPriceLine({
      price: 0,
      color: styles.muted,
      lineWidth: 2,
      lineStyle: LineStyle.Dashed,
      axisLabelVisible: true,
      title: "Upper Target",
    });

    const idealLevel = chart.addLineSeries({
      color: styles.secondary,
      lineWidth: 2,
      lineStyle: LineStyle.Dashed,
      priceLineVisible: false,
      lastValueVisible: false,
    }) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

    const currentLevel = chart.addLineSeries({
      title: "Current Level",
      color: styles.primary,
      lineWidth: 2,
      lineStyle: LineStyle.Solid,
      priceLineVisible: e2.plotSide() == "left",
      lastValueVisible: true,
      lastPriceAnimation: LastPriceAnimationMode.OnDataUpdate,
    }) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

    const currentLevelPrediction = chart.addLineSeries({
      color: styles.primary,
      lineWidth: 1,
      lineStyle: LineStyle.LargeDashed,
      priceLineVisible: false,
      lastValueVisible: false,
    }) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

    chart.subscribeCrosshairMove((ev) => {
      const idealData = e2.lineDataFromSeriesData(ev.seriesData, idealLevel);
      const actualData = e2.lineDataFromSeriesData(ev.seriesData, currentLevel);
      const predictedData = e2.lineDataFromSeriesData(ev.seriesData, currentLevelPrediction);

      plotTooltip = e2.renderPlotTooltip(plotTooltip, ev, plotDiv, plotTooltipDiv, [
        { name: "Ideal", value: idealData?.value ?? NaN },
        { name: "Current", value: actualData?.value ?? NaN },
        { name: "Predicted", value: predictedData?.value ?? NaN },
      ]);
    });

    return Object.assign(chart, {
      wpathFake,
      wpathLower,
      wpathUpper,
      idealLevel,
      currentLevel,
      currentLevelPrediction,
    });
  });

  // Watch changes to styling.
  $effect(() => {
    if (!chart) {
      return;
    }

    chart.applyOptions(e2.chartOptions(styles));
    chart.wpathFake.applyOptions({
      visible: e2.showIdealLevels(),
    });

    Object.values(chart)
      .filter((series): series is charts.ISeriesApi<"Line"> => series?.seriesType?.() == "Line")
      .forEach((series) => series.applyOptions(e2.lineSeriesOptions()));
  });

  // Watch changes to data.
  $effect(() => {
    if (!chart || !dosage) {
      return;
    }

    // Enable high-density mode for larger plots.
    e2.plotPreferences.plotHighDensity = plotWidth > 512;

    chart.wpathFake.setData([
      { time: interval.start.toUnixInteger() as charts.UTCTimestamp },
      { time: interval.end.toUnixInteger() as charts.UTCTimestamp },
    ]);
    chart.wpathLower.applyOptions({ price: wpathRange.lower });
    chart.wpathUpper.applyOptions({ price: wpathRange.upper });

    const lastDoseAt = history[history.length - 1]?.takenAt ?? interval.end;
    const idealData = e2.fillE2IdealData(totalInterval, dosage, e2.conversionFactor(), lastDoseAt);
    const actualData = e2.fillE2ActualData(totalInterval, history, e2.conversionFactor());

    chart.idealLevel.setData(idealData);
    chart.currentLevel.setData(e2.dataWithinInterval(actualData, interval));
    chart.currentLevelPrediction.setData(e2.dataWithinInterval(actualData, predictedInterval));

    chart.currentLevel.setMarkers(
      history.map((dose) => ({
        time: dose.takenAt.toUnixInteger() as charts.UTCTimestamp,
        position: "belowBar",
        color: styles.primary,
        shape: "arrowUp",
      })),
    );

    chart.timeScale().setVisibleRange({
      from: interval.end.minus({ weeks: 2 }).toUnixInteger() as charts.UTCTimestamp,
      to: predictedInterval.end.toUnixInteger() as charts.UTCTimestamp,
    });
  });
</script>

<div class="estrannaise spaced">
  <div class="estrannaise-plot" bind:this={plotDiv} bind:clientWidth={plotWidth}>
    <div
      class="estrannaise-plot-tooltip-container"
      style="top: {plotTooltip?.top ?? 0}px; left: {plotTooltip?.left ?? 0}px"
      bind:this={plotTooltipDiv}
    >
      {#if plotTooltip?.visible}
        <div
          class="estrannaise-plot-tooltip text-sm"
          transition:fade={{
            duration: 200,
          }}
        >
          <time class="text-sm" datetime={plotTooltip.time.toISO()}>
            <span class="font-bold">
              {plotTooltip.time.toLocaleString({
                month: "numeric",
                day: "numeric",
                weekday: "long",
              })}
            </span>
            <span>
              {plotTooltip.time.toLocaleString({
                hour: "numeric",
              })}
            </span>
          </time>
          <ul class="tooltip-data p-0">
            {#each plotTooltip.data! as { name, value }}
              <li>
                <span class="label font-bold">{name}</span>
                <span class="value text-right">{value}</span>
              </li>
            {/each}
          </ul>
        </div>
      {/if}
    </div>
  </div>

  <ul class="legends">
    <li class="actual primary">
      <Icon name="show-chart" /> Current Levels
    </li>
    <li class="steady secondary">
      <Icon name="auto-graph" /> Ideal Levels
    </li>
  </ul>
</div>

<style lang="scss">
  .estrannaise-plot {
    --pico-background-color: var(--pico-card-background-color);

    height: clamp(250px, 20vh, 500px);
    position: relative;

    // Reset pico.css' table styling.
    :global {
      th,
      td {
        border: initial;
      }
    }
  }

  .estrannaise-plot-tooltip-container {
    position: absolute;
    z-index: 10;
    pointer-events: none;

    .estrannaise-plot-tooltip {
      padding: calc(var(--pico-spacing) / 2);
      max-width: 40ch;

      border: var(--pico-border-width) solid var(--pico-primary-border);
      border-radius: var(--pico-border-radius);
      background: var(--pico-primary-background);
      color: var(--pico-primary-inverse);

      opacity: 0.75;
      backdrop-filter: blur(4px);
    }

    time {
      text-align: center;
      margin-bottom: calc(var(--pico-spacing) / 4);

      display: flex;
      flex-direction: row;
      justify-content: space-between;
      gap: calc(var(--pico-spacing) / 2);
    }

    ul {
      margin: 0;

      display: grid;
      grid-template-columns: 1fr auto;
      grid-column-gap: calc(var(--pico-spacing) / 2);

      li {
        display: contents; // grid time
      }
    }
  }

  .legends {
    width: 100%;
    padding: 0;

    display: flex;
    flex-direction: row;
    justify-content: center;
    gap: calc(2 * var(--pico-spacing));

    li {
      list-style: none;
    }
  }
</style>
