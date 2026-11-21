<script lang="ts">
  import * as api from "#lib/api.ts";
  import * as e2 from "#lib/e2.svelte.ts";
  import type * as charts from "lightweight-charts";
  import { ColorType, CrosshairMode, LineStyle, LineType, createChart } from "lightweight-charts";
  import { DateTime, Interval } from "luxon";
  import { untrack } from "svelte";

  type Props =
    | {
        mode: "steady";
        dosage: api.Dosage;
      }
    | {
        mode: "impact";
        dosage?: api.Dosage;
        beforeHistory: e2.Dose[];
        afterHistory: e2.Dose[];
        focusTime: DateTime;
      };

  let props: Props = $props();

  let plotDiv = $state<HTMLDivElement | null>(null);
  let styles = $derived(e2.gatherStyles(plotDiv));

  let steadySummary = $derived(
    props.mode === "steady" ? e2.sampleIdealCycle(props.dosage, e2.conversionFactor()) : undefined,
  );

  let previewData = $derived.by(() => {
    const cf = e2.conversionFactor();
    const wpath = e2.wpathRange(cf);

    if (props.mode === "steady") {
      const start = DateTime.now().startOf("day");
      const spanDays = Math.max(3 * Math.max(0.5, props.dosage.interval), 14);
      const end = start.plus({ milliseconds: Math.round(spanDays * 86_400_000) });
      const iv = Interval.fromDateTimes(start, end) as Interval<true>;
      const primary = e2.fillE2IdealData(iv, props.dosage, cf, { points: 160 });
      return {
        iv,
        wpath,
        primary,
        secondary: [] as charts.LineData<charts.UTCTimestamp>[],
        focusUnix: undefined as charts.UTCTimestamp | undefined,
      };
    }

    const futureDays = Math.max((props.dosage?.interval ?? 7) * 2, 10);
    const start = props.focusTime.minus({ days: 3 });
    const end = props.focusTime.plus({ milliseconds: Math.round(futureDays * 86_400_000) });
    const iv = Interval.fromDateTimes(start, end) as Interval<true>;
    const primary = e2.fillE2ActualData(iv, props.dosage, props.afterHistory, cf, { points: 160 });
    const secondary = e2.fillE2ActualData(iv, props.dosage, props.beforeHistory, cf, {
      points: 160,
    });
    const focusUnix = props.focusTime.toUnixInteger() as charts.UTCTimestamp;
    return { iv, wpath, primary, secondary, focusUnix };
  });

  let chart = $derived.by(() => {
    if (!plotDiv) {
      return;
    }

    const div = plotDiv;
    return untrack(() => {
      const chart = createChart(div, {
        autoSize: true,
        handleScale: false,
        handleScroll: false,
        crosshair: {
          mode: CrosshairMode.Hidden,
        },
      });

      const wpathFake = chart.addLineSeries({
        visible: true,
        lineVisible: false,
        pointMarkersVisible: false,
        lastValueVisible: false,
        priceLineVisible: false,
      }) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

      const initialMuted = e2.withAlpha(styles.muted, 0.7) || "#888888";

      const wpathLower = wpathFake.createPriceLine({
        price: 0,
        color: initialMuted,
        lineWidth: 1,
        lineStyle: LineStyle.Dashed,
        axisLabelVisible: false,
        title: "",
      });

      const wpathUpper = wpathFake.createPriceLine({
        price: 0,
        color: initialMuted,
        lineWidth: 1,
        lineStyle: LineStyle.Dashed,
        axisLabelVisible: false,
        title: "",
      });

      const secondary = chart.addLineSeries({
        lineType: LineType.Simple,
        lineWidth: 2,
        lineStyle: LineStyle.Dashed,
        priceLineVisible: false,
        lastValueVisible: false,
        crosshairMarkerVisible: false,
      }) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

      const primary = chart.addLineSeries({
        lineType: LineType.Simple,
        lineWidth: 2,
        lineStyle: LineStyle.Solid,
        priceLineVisible: false,
        lastValueVisible: false,
        crosshairMarkerVisible: false,
      }) as charts.ISeriesApi<"Line", charts.UTCTimestamp>;

      return Object.assign(chart, {
        wpathFake,
        wpathLower,
        wpathUpper,
        secondary,
        primary,
      });
    });
  });

  $effect(() => {
    const c = chart;
    return () => c?.remove();
  });

  $effect(() => {
    if (!chart) {
      return;
    }

    const baseOpts = e2.chartOptions(styles);
    chart.applyOptions({
      ...baseOpts,
      handleScale: false,
      handleScroll: false,
      crosshair: {
        mode: CrosshairMode.Hidden,
      },
      layout: {
        background: {
          type: ColorType.Solid,
          color: "transparent",
        },
        textColor: styles.color,
        fontFamily: styles.fontFamily,
        // @ts-ignore
        attributionLogo: false,
      },
      timeScale: {
        ...baseOpts.timeScale,
        ticksVisible: false,
      },
    });

    const previewLineOpts = {
      ...e2.lineSeriesOptions(),
      priceLineVisible: false,
      lastValueVisible: false,
      autoscaleInfoProvider: (autoscale: () => charts.AutoscaleInfo) => {
        const scale = autoscale();
        if (scale != null) {
          scale.priceRange.minValue = 0;
          scale.priceRange.maxValue *= 1.15;
          scale.margins = { above: 7, below: 7 };
        }
        return scale;
      },
    };

    chart.wpathFake.applyOptions({
      ...previewLineOpts,
      visible: true,
      lineVisible: false,
    });
    const mutedLineColor = e2.withAlpha(styles.muted, 0.7) || "#888888";
    chart.wpathLower.applyOptions({ color: mutedLineColor });
    chart.wpathUpper.applyOptions({ color: mutedLineColor });

    chart.secondary.applyOptions({
      ...previewLineOpts,
      color: e2.withAlpha(styles.muted, 0.85) || "#888888",
    });
    chart.primary.applyOptions({
      ...previewLineOpts,
      color: styles.primary || "#c27c88",
    });
  });

  $effect(() => {
    if (!chart) {
      return;
    }

    const { iv, wpath, primary, secondary, focusUnix } = previewData;
    chart.wpathFake.setData([
      { time: iv.start.toUnixInteger() as charts.UTCTimestamp },
      { time: iv.end.toUnixInteger() as charts.UTCTimestamp },
    ]);
    chart.wpathLower.applyOptions({ price: wpath.lower });
    chart.wpathUpper.applyOptions({ price: wpath.upper });

    chart.secondary.setData(secondary);
    chart.primary.setData(primary);

    if (focusUnix !== undefined && primary.length > 0) {
      chart.primary.setMarkers([
        {
          time: focusUnix,
          position: "belowBar",
          color: styles.primary,
          shape: "arrowUp",
        },
      ]);
    } else {
      chart.primary.setMarkers([]);
    }

    chart.timeScale().fitContent();
  });
</script>

<div class="levels-preview">
  {#if props.mode === "steady" && steadySummary}
    <div class="preview-header">
      <span class="preview-title">Estimated Steady-State Levels</span>
      <div class="summary-pills">
        <span class="pill">
          <small>Trough</small>
          <strong>{e2.fixed(steadySummary.trough, e2.units())}</strong>
        </span>
        <span class="pill">
          <small>Average</small>
          <strong>{e2.fixed(steadySummary.average, e2.units())}</strong>
        </span>
        <span class="pill">
          <small>Peak</small>
          <strong>{e2.fixed(steadySummary.peak, e2.units())}</strong>
        </span>
      </div>
    </div>
  {:else if props.mode === "impact"}
    <div class="preview-header">
      <span class="preview-title">Levels Preview</span>
      <div class="impact-legend">
        <span class="legend-item before">Before</span>
        <span class="legend-item after">After</span>
      </div>
    </div>
  {/if}

  <div class="preview-plot" bind:this={plotDiv}></div>
</div>

<style lang="scss">
  .levels-preview {
    display: flex;
    flex-direction: column;
    gap: calc(var(--pico-spacing) / 2);
  }

  .preview-header {
    display: flex;
    flex-wrap: wrap;
    align-items: center;
    justify-content: space-between;
    gap: 0.4rem 0.75rem;
  }

  .preview-title {
    font-size: 0.85em;
    color: var(--pico-muted-color);
    font-weight: 600;
  }

  .summary-pills {
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;

    .pill {
      display: inline-flex;
      align-items: baseline;
      gap: 0.35rem;
      padding: 0.15rem 0.55rem;
      border-radius: 999px;
      background: color-mix(in srgb, var(--pico-primary) 12%, transparent);
      border: var(--pico-border-width) solid
        color-mix(in srgb, var(--pico-primary) 28%, transparent);
      font-size: 0.8em;

      small {
        color: var(--pico-muted-color);
      }
    }
  }

  .impact-legend {
    display: inline-flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 0.8em;

    .legend-item {
      display: inline-flex;
      align-items: center;
      gap: 0.35rem;

      &::before {
        content: "";
        width: 1rem;
        border-top: 2px solid var(--pico-primary);
      }

      &.before {
        color: var(--pico-muted-color);
        &::before {
          border-top-style: dashed;
          border-top-color: var(--pico-muted-color);
        }
      }
    }
  }

  .preview-plot {
    height: 150px;
    width: 100%;
    position: relative;

    :global {
      th,
      td {
        border: initial;
      }
    }
  }
</style>
