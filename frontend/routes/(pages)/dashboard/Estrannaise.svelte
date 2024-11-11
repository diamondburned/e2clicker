<!-- estrannai.se as a svelte component -->

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

  let plotDiv: HTMLDivElement | null = null;
  let styles = $derived(gatherStyles(plotDiv));

  let chart = $derived.by(() => {
    if (!plotDiv) {
      return null;
    }

    return charts.createChart(plotDiv, {
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
      timeScale: {
        lockVisibleTimeRangeOnResize: true,
        borderVisible: false,
        timeVisible: true,
        secondsVisible: false,
        fixLeftEdge: true,
        fixRightEdge: true,
      },
    });
  });

  let seriesOptions = $derived.by(() => {
    if (!plotDiv || !delivery) {
      return null;
    }

    return {
      lineType: charts.LineType.Curved,
      lineWidth: 2,
      baseLineVisible: false,
      priceLineVisible: false,
      priceFormat: {
        type: "custom",
        formatter: (price: charts.BarPrice) => price.toFixed(precision()) + " " + units,
      },
      autoscaleInfoProvider: (autoscale: () => charts.AutoscaleInfo) => {
        const scale = autoscale();
        if (scale != null) {
          // Ensure min Y is always 0.
          scale.priceRange.minValue = 0;
          scale.priceRange.maxValue *= 1.25;
        }
        return scale;
      },
    } as charts.LineSeriesPartialOptions;
  });

  let pks = $derived(functions(dosage));

  // let average = $derived(
  //   delivery && dosage
  //     ? e2.e2ssAverage3C(conversionFactor() * dosage.dose, dosage.interval, ...pks.pkParams)
  //     : 0,
  // )styles && ;
  // let trough = $derived(pks.pk(0));

  // Draw the steady line.
  $effect(() => {
    if (chart && seriesOptions) {
      const steadyData = fillE2Data(startTime, endTime, pks.pk);
      return ensureLine(steadyData, chart, {
        ...seriesOptions,
        title: "Ideal Levels",
        color: styles.secondary,
        lineStyle: charts.LineStyle.Dashed,
      });
    }
  });

  // Draw the actual line.
  $effect(() => {
    if (chart && seriesOptions) {
      const actualData = [] as charts.LineData[];
      return ensureLine(actualData, chart, {
        ...seriesOptions,
        title: "Current Levels",
        color: styles.primary,
        lineStyle: charts.LineStyle.Solid,
      });
    }
  });
</script>

<div class="estrannaise spaced">
  <div class="estrannaise-plot" bind:this={plotDiv}></div>

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

    height: clamp(250px, 30vh, 500px);
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
