<!-- estrannai.se as a svelte component -->

<script lang="ts" module>
  function ptInElement(pt: { x: number; y: number }, elem: HTMLElement) {
    return pt.x > 0 && pt.y > 0 && pt.x < elem.clientWidth && pt.y < elem.clientHeight;
  }
</script>

<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";

  import * as e2 from "$lib/e2";
  import * as api from "$lib/api";
  import * as charts from "lightweight-charts";
  import { DateTime } from "luxon";
  import {
    units,
    precision,
    gatherStyles,
    ensureLine,
    fillE2Data,
    functions,
    conversionFactor,
  } from "$lib/e2/plot.svelte";

  let {
    startTime = $bindable(),
    endTime = $bindable(),
    doses,
  }: {
    startTime: DateTime;
    endTime: DateTime;
    doses?: {
      dosage?: api.Dosage;
      history?: api.DosageHistory;
    };
  } = $props();

  let dosage = $derived(doses?.dosage);
  let history = $derived(doses?.history ?? []);
  let delivery = $derived(dosage && e2.deliveryMethod(dosage.deliveryMethod));

  let pks = $derived(functions(dosage));

  let trough = $derived(pks.pk(0));
  let average = $derived(
    dosage
      ? e2.e2ssAverage3C(conversionFactor() * dosage.dose, dosage.interval, ...pks.pkParams)
      : 0,
  );

  let plotDiv: HTMLDivElement | null = null;
  let styles = $derived(gatherStyles(plotDiv));

  type PlotTooltipData = {
    left: number;
    top: number;
    time: string;
    data: [string, string][];
  };

  let plotTooltipDiv: HTMLDivElement | null = null;
  let plotTooltip = $state<PlotTooltipData | null>(null);

  let chart = $derived.by(() => {
    if (!plotDiv) {
      return;
    }

    const chart = charts.createChart(plotDiv);
    const idealLevels = chart.addLineSeries({
      title: "Ideal Levels",
      color: styles.secondary,
      lineStyle: charts.LineStyle.Dashed,
    });
    const currentLevels = chart.addLineSeries({
      title: "Current Levels",
      color: styles.primary,
      lineStyle: charts.LineStyle.Solid,
    });

    const tooltipMargin = 12;
    chart.subscribeCrosshairMove((ev) => {
      if (!plotDiv || !plotTooltipDiv || !ev.point || !ev.time || !ptInElement(ev.point, plotDiv)) {
        plotTooltip = null;
        return;
      }

      const { clientWidth: tooltipWidth, clientHeight: tooltipHeight } = plotTooltipDiv;
      const y = ev.point.y;

      let left = ev.point.x + tooltipMargin;
      if (left > plotDiv.clientWidth - tooltipWidth) {
        left = ev.point.x - tooltipMargin - tooltipWidth;
      }

      let top = y + tooltipMargin;
      if (top > plotDiv.clientHeight - tooltipHeight) {
        top = y - tooltipHeight - tooltipMargin;
      }

      const fixed = (x: number) => x.toFixed(precision());

      const time = DateTime.fromSeconds(ev.time as charts.UTCTimestamp);
      const idealData = ev.seriesData.get(idealLevels) as charts.LineData;

      plotTooltip = {
        left,
        top,
        time: time.toLocaleString({ dateStyle: "short", timeStyle: "long" }),
        data: [
          ["Average", `${fixed(average)} ${units}`],
          ["Trough", `${fixed(trough)} ${units}`],
          ["Ideal", `${fixed(idealData?.value ?? 0)} ${units}`],
        ],
      };
    });

    return Object.assign(chart, {
      idealLevels,
      currentLevels,
    });
  });

  // Watch changes to styles and other styling-related properties and
  // automatically apply them to the entire chart.
  $effect(() => {
    if (!chart) {
      return;
    }

    chart.applyOptions({
      autoSize: true,
      handleScale: false,
      handleScroll: false,
      grid: {
        vertLines: {
          visible: false,
        },
        horzLines: {
          color: styles.muted,
          style: charts.LineStyle.Solid,
          visible: true,
        },
      },
      layout: {
        background: {
          type: charts.ColorType.Solid,
          color: "transparent",
        },
        textColor: styles.color,
        fontFamily: styles.fontFamily,
        // @ts-ignore
        attributionLogo: false,
      },
      rightPriceScale: {
        borderColor: styles.muted,
        ticksVisible: true,
        scaleMargins: {
          top: 0,
          bottom: 0,
        },
      },
      overlayPriceScales: {},
      timeScale: {
        lockVisibleTimeRangeOnResize: true,
        borderVisible: false,
        timeVisible: true,
        secondsVisible: false,
        fixLeftEdge: true,
        fixRightEdge: true,
      },
    });

    [chart.idealLevels, chart.currentLevels].forEach((series) =>
      series.applyOptions({
        lineType: charts.LineType.Curved,
        lineWidth: 2,
        baseLineVisible: false,
        priceLineVisible: false,
        lastValueVisible: true,
        lastPriceAnimation: charts.LastPriceAnimationMode.OnDataUpdate,
        priceFormat: {
          type: "custom",
          formatter: (price: charts.BarPrice) => price.toFixed(precision()) + " " + units,
        },
        autoscaleInfoProvider: (autoscale: () => charts.AutoscaleInfo) => {
          const scale = autoscale();
          if (scale != null) {
            // Ensure min Y is always 0.
            scale.priceRange.minValue *= 0;
            scale.priceRange.maxValue *= 1.25;
          }
          return scale;
        },
      }),
    );
  });

  $effect(() => {
    if (chart) {
      chart.idealLevels.setData(fillE2Data(startTime, endTime, pks.pk));
      chart.timeScale().fitContent();
    }
  });
</script>

<div class="estrannaise spaced">
  <div class="estrannaise-plot" bind:this={plotDiv}>
    <div
      class="estrannaise-plot-tooltip text-sm"
      class:show={plotTooltip != null}
      bind:this={plotTooltipDiv}
      style="top: {plotTooltip?.top ?? 0}px; left: {plotTooltip?.left ?? 0}px"
    >
      <ul class="tooltip-data list-none">
        {#each plotTooltip?.data ?? [] as [label, value]}
          <li>
            <strong>{label}</strong>: {value}
          </li>
        {/each}
      </ul>
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
  }

  .estrannaise-plot-tooltip {
    position: absolute;
    padding: calc(var(--pico-spacing) / 2);

    width: 20ch;
    z-index: 10;
    top: 12px;
    left: 12px;
    pointer-events: none;

    border: var(--pico-border-width) solid var(--pico-primary-border);
    border-radius: var(--pico-border-radius);
    background: var(--pico-primary-background);
    color: var(--pico-primary-inverse);

    display: none;

    &.show {
      display: block;
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
