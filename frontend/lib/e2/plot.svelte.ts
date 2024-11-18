import * as api from "$lib/api";
import type * as charts from "lightweight-charts";
import { LineStyle, ColorType, LineType, CrosshairMode } from "lightweight-charts";
import {
  availableUnits,
  e2ssAverage3C,
  PKFunctions,
  PKParameters,
  PKRandomFunctions,
} from "estrannaise/src/models";
import { deliveryMethod, estrannaiseDeliveryMethod } from "./methods.svelte";
import type { DosageObservation, PlotStyles } from "./plot";
import type { LevelUnits } from "estrannaise/src/modeldata";
import { DateTime, Duration, Interval } from "luxon";

export * from "./plot";

class PlotPreferences {
  units = $state<LevelUnits>("pg/mL");
  plotSide = $state<"left" | "right">("left");
  showIdealLevels = $state<boolean>(false);
  plotHighDensity = $state<boolean>(false);
  // how far to predict ahead for the actual E2 curve.
  predictAhead = $state(Duration.fromObject({ week: 1 }));
}

// Preferences.
export let plotPreferences = new PlotPreferences();

// Derived properties.
export let units = () => plotPreferences.units;
export let plotSide = () => plotPreferences.plotSide;
export let showIdealLevels = () => plotPreferences.showIdealLevels;
export let predictAhead = () => plotPreferences.predictAhead;
export let precision = () => availableUnits[units()].precision;
export let conversionFactor = () => availableUnits[units()].conversionFactor;
// how many points to plot in total. the more points, the smoother the curve.
export let plotPoints = () => (plotPreferences.plotHighDensity ? 800 : 200);

export let chartOptions = (styles: PlotStyles): charts.DeepPartial<charts.ChartOptions> => ({
  autoSize: true,
  handleScale: true,
  handleScroll: false,
  grid: {
    vertLines: {
      color: styles.mutedBorder,
      style: LineStyle.Solid,
      visible: true,
    },
    horzLines: {
      color: styles.mutedBorder,
      style: LineStyle.Solid,
      visible: true,
    },
  },
  crosshair: {
    mode: CrosshairMode.Normal,
    vertLine: {
      color: styles.muted,
      style: LineStyle.Dashed,
      labelVisible: true,
      labelBackgroundColor: styles.color,
    },
    horzLine: {
      color: styles.muted,
      style: LineStyle.Dashed,
      labelVisible: true,
      labelBackgroundColor: styles.color,
    },
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
  rightPriceScale: {
    visible: plotSide() == "right",
    borderColor: styles.muted,
    ticksVisible: true,
    entireTextOnly: true,
    scaleMargins: {
      top: 0.5,
      bottom: 0,
    },
  },
  leftPriceScale: {
    visible: plotSide() == "left",
    borderColor: styles.muted,
    ticksVisible: true,
    entireTextOnly: true,
    scaleMargins: {
      top: 0,
      bottom: 0,
    },
  },
  timeScale: {
    borderColor: styles.muted,
    lockVisibleTimeRangeOnResize: false,
    rightBarStaysOnScroll: true,
    ticksVisible: true,
    uniformDistribution: true,
  },
});

export let lineSeriesOptions = (): charts.DeepPartial<
  charts.LineStyleOptions & charts.SeriesOptionsCommon
> => ({
  lineType: LineType.Simple,
  baseLineVisible: false,
  priceScaleId: plotSide(),
  priceFormat: {
    type: "custom",
    formatter: (price: charts.BarPrice) => price.toFixed(precision()) + " " + units(),
  },
  autoscaleInfoProvider: (autoscale: () => charts.AutoscaleInfo) => {
    const scale = autoscale();
    if (scale != null) {
      // Ensure min Y is always 0.
      scale.priceRange.minValue *= 0;
      scale.priceRange.maxValue *= 1.5;
    }
    return scale;
  },
});

export function fillCurve(
  iv: Interval<true>,
  f: (t: DateTime, tprev?: DateTime) => number,
): charts.LineData<charts.UTCTimestamp>[] {
  const startMillis = iv.start.toMillis();
  const endMillis = iv.end.toMillis();
  const stepMillis = (endMillis - startMillis) / (plotPoints() - 1);

  let pts = new Array<ReturnType<typeof fillCurve>[number]>(plotPoints());
  let i = 0;

  let tprev: DateTime | undefined;
  for (let tms = startMillis; tms <= endMillis; tms += stepMillis) {
    const t = DateTime.fromMillis(tms);
    const v = f(t, tprev);
    tprev = t;

    if (!isNaN(v)) {
      pts[i] = {
        time: t.toUnixInteger() as charts.UTCTimestamp,
        value: v,
      };
    }

    i++;
  }

  pts = pts.slice(0, i);

  return pts;
}

// fillE2IdealData fills a curve with E2 data from the given E2 function.
export function fillE2IdealData(
  iv: Interval<true>,
  dosage: api.Dosage,
  conversionFactor = 1.0,
  takenAt?: DateTime,
) {
  const offset = takenAt ? dateToX(iv.start, takenAt) : 0;
  const delivery = estrannaiseDeliveryMethod(dosage.deliveryMethod);
  const pk = PKFunctions(conversionFactor)[delivery];
  return fillCurve(iv, (t) =>
    pk(dateToX(iv.start, t) - offset, dosage.dose, true, dosage.interval),
  );
}

// fillE2ActualData fills a curve with E2 data from the given dosage observations.
export function fillE2ActualData(
  iv: Interval<true>,
  history: DosageObservation[],
  conversionFactor = 1.0,
  random: boolean | number = false,
) {
  const deliveryMethods = history.map((o) => estrannaiseDeliveryMethod(o.deliveryMethod.id));
  const pkFactory = !random ? PKFunctions(conversionFactor) : PKRandomFunctions(conversionFactor);
  const randomParam =
    typeof random == "number" && random > 0 //
      ? random
      : undefined;

  return fillCurve(iv, (t, tprev) => {
    if (!tprev) {
      return NaN;
    }

    const xstart = dateToX(iv.start, t);
    const tdiff = t.diff(tprev).as("days");

    return history.reduce((sum, dose, i) => {
      const x = xstart - dateToX(iv.start, dose.takenAt);
      return sum + pkFactory[deliveryMethods[i]](x, dose.dose, false, tdiff, randomParam);
    }, 0);
  });
}

export function idealE2Trough(dosage: api.Dosage, conversionFactor = 1.0): number {
  const delivery = estrannaiseDeliveryMethod(dosage.deliveryMethod);
  const pk = PKFunctions(conversionFactor)[delivery];
  return pk(0, dosage.dose, true, dosage.interval);
}

export function idealE2Average(dosage: api.Dosage, conversionFactor = 1.0): number | undefined {
  if (deliveryMethod(dosage.deliveryMethod)?.patch) {
    return undefined;
  }
  const delivery = estrannaiseDeliveryMethod(dosage.deliveryMethod);
  const pkParams = PKParameters[delivery];
  return e2ssAverage3C(conversionFactor * dosage.dose, dosage.interval, ...pkParams);
}

// convert a date to an x value
function dateToX(start: DateTime, date: DateTime): number {
  return date.diff(start).as("days");
}
