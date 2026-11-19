<!-- estrannai.se as a svelte component -->

<script lang="ts" module>
  import type * as charts from "lightweight-charts";
  import * as e2 from "#lib/e2.svelte.ts";
  import { Interval } from "luxon";

  export class Data {
    private interval_ = $state<Interval<true>>();
    private doses = $state<{
      dosage?: api.Dosage;
      history?: e2.DosageHistory;
    }>({});

    readonly dosage = $derived(this.doses?.dosage);
    readonly history = $derived(this.doses?.history ?? []);
    readonly interval = $derived(this.interval_!);

    readonly predictedInterval = $derived(
      Interval.fromDateTimes(
        // Predict ahead proportional to the dosing interval (falling back to preference).
        this.interval.end,
        this.interval.end.plus(
          this.dosage
            ? {
                milliseconds: Math.round(
                  e2.defaultPlotPredictDays(this.dosage.interval) * 86_400_000,
                ),
              }
            : e2.predictAhead(),
        ),
      ) as Interval<true>,
    );
    readonly totalInterval = $derived(
      this.interval.union(this.predictedInterval) as Interval<true>,
    );

    readonly wpathRange = $derived(e2.wpathRange(e2.conversionFactor()));

    private readonly actualWhole = $derived(
      this.dosage
        ? e2.fillE2ActualWithUncertainty(
            this.interval,
            this.predictedInterval,
            this.dosage,
            this.history,
            e2.conversionFactor(),
          )
        : undefined,
    );

    readonly idealData = $derived(
      this.dosage
        ? e2.fillE2IdealData(this.totalInterval, this.dosage, e2.conversionFactor(), {
            takenAt: this.history[this.history.length - 1]?.takenAt ?? this.interval.end,
            timestamps: this.actualWhole?.timestamps,
          })
        : [],
    );

    readonly idealTrough = $derived(
      this.dosage && e2.idealE2Trough(this.dosage, e2.conversionFactor()),
    );
    readonly idealPeak = $derived(
      this.dosage && e2.idealE2Peak(this.dosage, e2.conversionFactor()),
    );
    readonly idealAverage = $derived(
      this.dosage && e2.idealE2Average(this.dosage, e2.conversionFactor()),
    );

    readonly actualData = $derived(
      this.actualWhole
        ? this.actualWhole.actual.filter((d) => d.time <= this.actualWhole!.splitTime)
        : [],
    );

    readonly actualPredictionData = $derived(
      this.actualWhole
        ? this.actualWhole.actual.filter((d) => d.time >= this.actualWhole!.splitTime)
        : [],
    );

    readonly actualP05Data = $derived(
      this.actualWhole
        ? this.actualWhole.p05.filter((d) => d.time <= this.actualWhole!.splitTime)
        : [],
    );

    readonly actualP95Data = $derived(
      this.actualWhole
        ? this.actualWhole.p95.filter((d) => d.time <= this.actualWhole!.splitTime)
        : [],
    );

    readonly predictedP05Data = $derived(
      this.actualWhole
        ? this.actualWhole.p05.filter((d) => d.time >= this.actualWhole!.splitTime)
        : [],
    );

    readonly predictedP95Data = $derived(
      this.actualWhole
        ? this.actualWhole.p95.filter((d) => d.time >= this.actualWhole!.splitTime)
        : [],
    );

    constructor(init: {
      interval: NonNullable<typeof Data.prototype.interval>;
      doses: NonNullable<typeof Data.prototype.doses>;
    }) {
      if (init) {
        this.interval_ = init.interval;
        this.doses = init.doses;
      }
    }
  }
</script>

<script lang="ts">
  import Icon from "#lib/components/Icon.svelte";
  import Tooltip from "#lib/components/popovers/Tooltip.svelte";

  import * as api from "#lib/api.ts";
  import { createChart, LineStyle, LastPriceAnimationMode } from "lightweight-charts";
  import { untrack } from "svelte";
  import { fade } from "svelte/transition";

  let {
    data,
  }: {
    data: Data;
  } = $props();

  let plotDiv = $state<HTMLDivElement | null>(null);
  let plotWidth = $state(0);

  let styles = $derived(e2.gatherStyles(plotDiv));

  let plotTooltipDiv: HTMLDivElement | null = null;
  let plotTooltip = $state<e2.PlotTooltipData | null>(null);

  let chart = $derived.by(() => {
    if (!plotDiv) {
      return;
    }

    const div = plotDiv;
    return untrack(() => {
      const chart = createChart(div);

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

      const uncertaintyFillColor = e2.withAlpha(styles.primary, 0.2);
      const uncertaintyPredictedFillColor = e2.withAlpha(styles.primary, 0.12);

      const actualLower = chart.addLineSeries(
        e2.uncertaintyLineSeriesOptions(),
      ) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

      const actualUpper = chart.addLineSeries(
        e2.uncertaintyLineSeriesOptions(),
      ) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

      const actualBand = new e2.UncertaintyBandPrimitive({
        color: uncertaintyFillColor,
        visible: e2.showUncertainty(),
      });
      actualUpper.attachPrimitive(actualBand);

      const predictedLower = chart.addLineSeries(
        e2.uncertaintyLineSeriesOptions(),
      ) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

      const predictedUpper = chart.addLineSeries(
        e2.uncertaintyLineSeriesOptions(),
      ) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

      const predictedBand = new e2.UncertaintyBandPrimitive({
        color: uncertaintyPredictedFillColor,
        visible: e2.showUncertainty(),
      });
      predictedUpper.attachPrimitive(predictedBand);

      chart.subscribeCrosshairMove((ev) => {
        const idealData = e2.lineDataFromSeriesData(ev.seriesData, idealLevel);
        const actualData = e2.lineDataFromSeriesData(ev.seriesData, currentLevel);
        const predictedData = e2.lineDataFromSeriesData(ev.seriesData, currentLevelPrediction);
        const showUnc = e2.showUncertainty();
        const actualP05 = showUnc
          ? e2.lineDataFromSeriesData(ev.seriesData, actualLower)
          : undefined;
        const actualP95 = showUnc
          ? e2.lineDataFromSeriesData(ev.seriesData, actualUpper)
          : undefined;
        const predP05 = showUnc
          ? e2.lineDataFromSeriesData(ev.seriesData, predictedLower)
          : undefined;
        const predP95 = showUnc
          ? e2.lineDataFromSeriesData(ev.seriesData, predictedUpper)
          : undefined;

        plotTooltip = e2.renderPlotTooltip(plotTooltip, ev, plotDiv, plotTooltipDiv, [
          { name: "Ideal", value: idealData?.value ?? NaN },
          {
            name: "Current",
            value: actualData?.value ?? NaN,
            range: actualP05 && actualP95 ? [actualP05.value, actualP95.value] : undefined,
          },
          {
            name: "Predicted",
            value: predictedData?.value ?? NaN,
            range: !actualData && predP05 && predP95 ? [predP05.value, predP95.value] : undefined,
          },
        ]);
      });

      return Object.assign(chart, {
        wpathFake,
        wpathLower,
        wpathUpper,
        idealLevel,
        currentLevel,
        currentLevelPrediction,
        actualLower,
        actualUpper,
        actualBand,
        predictedLower,
        predictedUpper,
        predictedBand,
      });
    });
  });

  $effect(() => {
    const c = chart;
    return () => c?.remove();
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

    const baseLineOpts = e2.lineSeriesOptions();
    for (const series of [
      chart.wpathFake,
      chart.idealLevel,
      chart.currentLevel,
      chart.currentLevelPrediction,
    ]) {
      series.applyOptions(baseLineOpts);
    }

    const showUnc = e2.showUncertainty();
    const uncOpts = e2.uncertaintyLineSeriesOptions(showUnc);
    chart.actualLower.applyOptions(uncOpts);
    chart.actualUpper.applyOptions(uncOpts);
    chart.actualBand.setOptions({
      color: e2.withAlpha(styles.primary, 0.2),
      visible: showUnc,
    });
    chart.predictedLower.applyOptions(uncOpts);
    chart.predictedUpper.applyOptions(uncOpts);
    chart.predictedBand.setOptions({
      color: e2.withAlpha(styles.primary, 0.12),
      visible: showUnc,
    });
  });

  // Watch changes to data.
  $effect(() => {
    chart;
    data;

    if (!chart) {
      return;
    }

    // Enable high-density mode for larger plots.
    e2.plotPreferences.plotHighDensity = plotWidth > 512;

    chart.wpathFake.setData([
      { time: data.interval.start.toUnixInteger() as charts.UTCTimestamp },
      { time: data.interval.end.toUnixInteger() as charts.UTCTimestamp },
    ]);
    chart.wpathLower.applyOptions({ price: data.wpathRange.lower });
    chart.wpathUpper.applyOptions({ price: data.wpathRange.upper });

    chart.idealLevel.setData(data.idealData);
    chart.currentLevel.setData(data.actualData);
    chart.currentLevelPrediction.setData(data.actualPredictionData);
    chart.actualLower.setData(data.actualP05Data);
    chart.actualUpper.setData(data.actualP95Data);
    chart.actualBand.setData(data.actualP05Data, data.actualP95Data);
    chart.predictedLower.setData(data.predictedP05Data);
    chart.predictedUpper.setData(data.predictedP95Data);
    chart.predictedBand.setData(data.predictedP05Data, data.predictedP95Data);

    chart.currentLevel.setMarkers(
      data.history
        .filter((dose) => dose.takenAt >= data.interval.start && dose.takenAt <= data.interval.end)
        .map((dose) => ({
          time: dose.takenAt.toUnixInteger() as charts.UTCTimestamp,
          position: "belowBar",
          color: styles.primary,
          shape: "arrowUp",
        })),
    );

    const pastDays = data.dosage ? e2.defaultPlotPastDays(data.dosage.interval) : 14;
    const desiredFrom = data.interval.end.minus({
      milliseconds: Math.round(pastDays * 86_400_000),
    });
    const visibleFrom = desiredFrom < data.interval.start ? data.interval.start : desiredFrom;

    chart.timeScale().setVisibleRange({
      from: visibleFrom.toUnixInteger() as charts.UTCTimestamp,
      to: data.predictedInterval!.end.toUnixInteger() as charts.UTCTimestamp,
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
            {#each plotTooltip.data! as { name, value, range }}
              <li>
                <span class="label font-bold">{name}</span>
                <span class="value text-right">
                  {value}
                  {#if range}
                    <small class="range">{range}</small>
                  {/if}
                </span>
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
    <li class="uncertainty">
      <button
        type="button"
        class="legend-toggle"
        class:active={e2.showUncertainty()}
        aria-pressed={e2.showUncertainty()}
        onclick={() => {
          e2.plotPreferences.showUncertainty = !e2.plotPreferences.showUncertainty;
        }}
      >
        <Icon name="stacked-line-chart" /> Uncertainty
      </button>
      <Tooltip>
        <button type="button" class="legend-help" aria-label="About uncertainty">
          <Icon name="help" />
        </button>
        {#snippet tooltip()}
          5–95% confidence interval across individual pharmacokinetic variation.
        {/snippet}
      </Tooltip>
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

      .value {
        display: flex;
        flex-direction: column;
        align-items: flex-end;
      }

      .range {
        font-size: 0.8em;
        opacity: 0.85;
        white-space: nowrap;
      }
    }
  }

  .legends {
    width: 100%;
    padding: 0;

    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    justify-content: center;
    gap: var(--pico-spacing) calc(2 * var(--pico-spacing));

    li {
      list-style: none;
      display: inline-flex;
      align-items: center;
      gap: 0.35em;
    }

    .legend-toggle,
    .legend-help {
      background: none;
      border: none;
      box-shadow: none;
      padding: 0;
      margin: 0;
      text-decoration: none;
      display: inline-flex;
      align-items: center;
      gap: 0.3em;
    }

    .legend-toggle {
      color: var(--pico-muted-color);
      opacity: 0.65;
      cursor: pointer;

      &:hover {
        opacity: 0.9;
      }

      &.active {
        color: color-mix(in srgb, var(--pico-primary) 75%, var(--pico-color));
        opacity: 1;
      }
    }

    .legend-help {
      color: var(--pico-muted-color);
      opacity: 0.55;
      cursor: help;
      line-height: 1;
      vertical-align: center;

      &:hover,
      &:focus-visible {
        opacity: 0.95;
        color: var(--pico-color);
      }

      :global(.material-symbols-rounded) {
        font-size: 1.05em;
      }
    }
  }
</style>
