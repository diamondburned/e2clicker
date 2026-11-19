import * as api from "#lib/api.ts";
import type * as charts from "lightweight-charts";
import { LineStyle, ColorType, LineType, CrosshairMode } from "lightweight-charts";
import {
  availableUnits,
  e2ssAverage3C,
  e2Patch3C,
  randomMCMCSample,
  PKFunctions,
  PKParameters,
  PKRandomFunctions,
  terminalEliminationTime3C,
} from "estrannaise/src/models";
import { deliveryMethod, estrannaiseDeliveryMethod } from "./methods.svelte";
import { timeUntilNextDose } from "./timing";
import type { Dose, PlotStyles } from "./plot";
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

// fixed formats a number to a fixed precision and units.
export function fixed<XType extends number | undefined>(
  x: XType,
  units_?: string,
): XType extends undefined ? string | undefined : string {
  if (x === undefined) {
    return undefined as any;
  }
  return x.toFixed(precision()) + (units_ ? " " + units_ : "");
}

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
      labelVisible: false,
      labelBackgroundColor: styles.color,
    },
    horzLine: {
      visible: false,
      // color: styles.muted,
      // style: LineStyle.Dashed,
      // labelVisible: true,
      // labelBackgroundColor: styles.color,
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

export function defaultPlotPastDays(intervalDays: number): number {
  const safeInterval = Math.max(0.5, intervalDays);
  // Show 11 dosing intervals in the past window (plus predictAhead) so >10 peaks
  // are visible by default regardless of injection or patch frequency.
  return 11 * safeInterval;
}

export function defaultPlotPredictDays(intervalDays: number): number {
  const safeInterval = Math.max(0.5, intervalDays);
  // Predict ~1 cycle ahead, clamped to [2, 21] days so daily/patch curves stay
  // balanced while long esters still show a full upcoming cycle.
  return Math.max(2, Math.min(21, safeInterval));
}

// compoundingBurnInDays returns how many days of prior dosing history are needed
// before the plotted curve reaches ~steady-state compounding (~5 half-lives, or
// >=4 half-lives when constrained by the backend's 1-year history limit).
export function compoundingBurnInDays(dosage: api.Dosage): number {
  const method = estrannaiseDeliveryMethod(dosage.deliveryMethod);
  const isPatch = deliveryMethod(dosage.deliveryMethod)?.patch;
  const wearDays = isPatch ? dosage.interval * effectiveConcurrence(dosage) : 0;
  const fullBurnIn = Math.ceil(terminalEliminationTime3C(method) + wearDays);
  const visiblePast = Math.ceil(defaultPlotPastDays(dosage.interval));
  if (fullBurnIn + visiblePast <= 360) {
    return fullBurnIn;
  }
  const minBurnIn = Math.ceil(terminalEliminationTime3C(method, 4) + wearDays);
  return Math.max(minBurnIn, 360 - visiblePast);
}

export function fillCurve(
  iv: Interval<true>,
  f: (t: DateTime, tprev?: DateTime) => number,
  points = plotPoints(),
): charts.LineData<charts.UTCTimestamp>[] {
  const startMillis = iv.start.toMillis();
  const endMillis = iv.end.toMillis();
  const count = Math.max(2, points);
  const stepMillis = (endMillis - startMillis) / (count - 1);

  let pts = new Array<ReturnType<typeof fillCurve>[number]>(count);
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

function curvePointsForDosage(iv: Interval<true>, dosage: api.Dosage): number {
  const cycles = iv.length("days") / Math.max(0.5, dosage.interval);
  const perCycle = plotPreferences.plotHighDensity ? 24 : 16;
  return Math.min(1600, Math.max(plotPoints(), Math.ceil(cycles * perCycle)));
}

function effectiveConcurrence(dosage?: api.Dosage): number {
  return dosage?.concurrence && dosage.concurrence > 0 ? dosage.concurrence : 1;
}

export function patchWearDuration(
  dose: Pick<Dose, "deliveryMethod" | "takenAt" | "takenOffAt">,
  dosage?: api.Dosage,
): number {
  if (dose.takenOffAt) {
    return Math.max(0, dose.takenOffAt.diff(dose.takenAt).as("days"));
  }
  if (dosage && dose.deliveryMethod.id === dosage.deliveryMethod) {
    return dosage.interval * effectiveConcurrence(dosage);
  }
  return dose.deliveryMethod.id === "patch ow" ? 7.0 : 3.5;
}

function makePKFunctions(
  dosage: api.Dosage,
  conversionFactor: number,
  {
    random = false,
    interval = dosage.interval * effectiveConcurrence(dosage),
  }: {
    random?: boolean | number;
    interval?: number;
  } = {},
) {
  const cf = conversionFactor;
  return !random
    ? {
        ...PKFunctions(cf),
        "patch tw": (t: number, dose: number, steady = false, T = 0.0) =>
          e2Patch3C(t, cf * dose, ...PKParameters["patch tw"], interval, steady, T),
        "patch ow": (t: number, dose: number, steady = false, T = 0.0) =>
          e2Patch3C(t, cf * dose, ...PKParameters["patch ow"], interval, steady, T),
      }
    : {
        ...PKRandomFunctions(cf),
        "patch tw": (t: number, dose: number, steady = false, T = 0.0, ri?: number) =>
          e2Patch3C(t, cf * dose, ...randomMCMCSample("patch tw", ri), interval, steady, T),
        "patch ow": (t: number, dose: number, steady = false, T = 0.0, ri?: number) =>
          e2Patch3C(t, cf * dose, ...randomMCMCSample("patch ow", ri), interval, steady, T),
      };
}

// fillE2IdealData fills a curve with E2 data from the given E2 function.
export function fillE2IdealData(
  iv: Interval<true>,
  dosage: api.Dosage,
  conversionFactor: number,
  {
    takenAt,
  }: {
    takenAt?: DateTime;
  } = {},
) {
  const offset = takenAt ? dateToX(iv.start, takenAt) : 0;
  const delivery = estrannaiseDeliveryMethod(dosage.deliveryMethod);

  const pkFunctions = makePKFunctions(dosage, conversionFactor);
  const pk = pkFunctions[delivery];

  return fillCurve(
    iv,
    (t) => pk(dateToX(iv.start, t) - offset, dosage.dose, true, dosage.interval),
    curvePointsForDosage(iv, dosage),
  );
}

function makeActualEvaluator(
  dosage: api.Dosage | undefined,
  history: Dose[],
  conversionFactor_: number,
  random: boolean | number = false,
): (tMillis: number) => number {
  const prepared = history.map((dose) => {
    const delivery = estrannaiseDeliveryMethod(dose.deliveryMethod.id);
    return {
      takenAtMillis: dose.takenAt.toMillis(),
      dose: dose.dose,
      delivery,
      wearDuration:
        delivery === "patch tw" || delivery === "patch ow"
          ? patchWearDuration(dose, dosage)
          : undefined,
    };
  });
  const pkFunctions = random
    ? PKRandomFunctions(conversionFactor_)
    : PKFunctions(conversionFactor_);
  const randomIdx = typeof random === "number" && random > 0 ? random : undefined;

  return (tMillis: number): number => {
    let sum = 0;
    for (let i = 0; i < prepared.length; i++) {
      const dose = prepared[i];
      const x = (tMillis - dose.takenAtMillis) / 86_400_000;
      if (x < 0) {
        continue;
      }
      if (dose.wearDuration !== undefined) {
        const params = random
          ? randomMCMCSample(dose.delivery, randomIdx)
          : PKParameters[dose.delivery];
        sum += e2Patch3C(
          x,
          conversionFactor_ * dose.dose,
          ...params,
          dose.wearDuration,
          false,
          0.0,
        );
      } else {
        sum += pkFunctions[dose.delivery](x, dose.dose, false, 0.0, randomIdx);
      }
    }
    return sum;
  };
}

// evalActualE2 evaluates the compounded E2 concentration at time t from all
// logged doses in history.
export function evalActualE2(
  t: DateTime,
  dosage: api.Dosage | undefined,
  history: Dose[],
  conversionFactor_ = conversionFactor(),
  random: boolean | number = false,
): number {
  return makeActualEvaluator(dosage, history, conversionFactor_, random)(t.toMillis());
}

// fillE2ActualData fills a curve with E2 data from the given dosage observations.
export function fillE2ActualData(
  iv: Interval<true>,
  dosage: api.Dosage,
  history: Dose[],
  conversionFactor: number,
  {
    random = false,
  }: {
    random?: boolean | number;
  } = {},
) {
  const evalAt = makeActualEvaluator(dosage, history, conversionFactor, random);
  return fillCurve(iv, (t) => evalAt(t.toMillis()), curvePointsForDosage(iv, dosage));
}

function sampleIdealCycle(
  dosage: api.Dosage,
  conversionFactor_ = 1.0,
  samples = 120,
): { trough: number; peak: number; average: number } {
  const delivery = estrannaiseDeliveryMethod(dosage.deliveryMethod);
  const pk = makePKFunctions(dosage, conversionFactor_)[delivery];
  const T = dosage.interval;
  if (T <= 0) {
    const v = pk(0, dosage.dose, true, T);
    return { trough: v, peak: v, average: v };
  }

  let trough = Infinity;
  let peak = -Infinity;
  let area = 0;
  let prevV = pk(0, dosage.dose, true, T);
  trough = Math.min(trough, prevV);
  peak = Math.max(peak, prevV);

  const dt = T / samples;
  for (let i = 1; i <= samples; i++) {
    const x = i === samples ? T - 1e-9 : i * dt;
    const v = pk(x, dosage.dose, true, T);
    trough = Math.min(trough, v);
    peak = Math.max(peak, v);
    area += 0.5 * (prevV + v) * dt;
    prevV = v;
  }

  const average = deliveryMethod(dosage.deliveryMethod)?.patch
    ? area / T
    : e2ssAverage3C(conversionFactor_ * dosage.dose, T, ...PKParameters[delivery]);

  return { trough, peak, average };
}

export function idealE2Trough(dosage: api.Dosage, conversionFactor = 1.0): number {
  return sampleIdealCycle(dosage, conversionFactor).trough;
}

export function idealE2Peak(dosage: api.Dosage, conversionFactor = 1.0): number {
  return sampleIdealCycle(dosage, conversionFactor).peak;
}

export function idealE2Average(dosage: api.Dosage, conversionFactor = 1.0): number {
  return sampleIdealCycle(dosage, conversionFactor).average;
}

export type LevelTrend =
  | { direction: "rising"; slopePerDay: number; peakValue?: number; peakAt?: DateTime }
  | { direction: "falling"; slopePerDay: number }
  | { direction: "steady"; slopePerDay: number };

export type LevelMetrics = {
  ideal: {
    trough: number;
    peak: number;
    average: number;
  };
  current?: number;
  trend?: LevelTrend;
  nextDoseAt?: DateTime;
  nextTrough?: number;
  isOverdue: boolean;
  window?: {
    start: DateTime;
    end: DateTime;
    days: number;
    clampedToFirstDose: boolean;
    trough: number;
    peak: number;
  };
  cadence?: {
    avgIntervalDays: number;
    sampleCount: number;
  };
};

export function computeLevelMetrics({
  now,
  dosage,
  history,
  conversionFactor: cf = conversionFactor(),
}: {
  now: DateTime;
  dosage: api.Dosage;
  history: Dose[];
  conversionFactor?: number;
}): LevelMetrics {
  const ideal = sampleIdealCycle(dosage, cf);
  if (history.length === 0) {
    return { ideal, isOverdue: false };
  }

  const sorted = [...history].sort((a, b) => a.takenAt.toMillis() - b.takenAt.toMillis());
  const evalAt = makeActualEvaluator(dosage, sorted, cf);

  const nowMillis = now.toMillis();
  const current = evalAt(nowMillis);

  const lastDose = sorted[sorted.length - 1];
  const nextDoseAt = timeUntilNextDose(dosage, lastDose.takenAt);
  const nextDoseMillis = nextDoseAt.toMillis();
  const isOverdue = nowMillis >= nextDoseMillis;
  const nextTrough = isOverdue ? current : evalAt(nextDoseMillis);

  // Forward difference over 15 minutes avoids straddling a discontinuous slope jump at takenAt.
  const dtMinutes = 15;
  const dtMillis = dtMinutes * 60_000;
  const slopePerDay = (evalAt(nowMillis + dtMillis) - current) / (dtMinutes / (24 * 60));
  const slopeThreshold = 2 * cf;

  let trend: LevelTrend;
  if (slopePerDay > slopeThreshold) {
    const scanEndMillis = isOverdue
      ? nowMillis + Math.max(1, dosage.interval) * 86_400_000
      : nextDoseMillis;
    const scanSteps = 96;
    let bestValue = current;
    let bestMillis = nowMillis;
    if (scanEndMillis > nowMillis) {
      for (let i = 1; i <= scanSteps; i++) {
        const t = nowMillis + ((scanEndMillis - nowMillis) * i) / scanSteps;
        const v = evalAt(t);
        if (v > bestValue) {
          bestValue = v;
          bestMillis = t;
        }
      }
    }
    trend =
      bestMillis > nowMillis && bestValue > current * 1.01
        ? {
            direction: "rising",
            slopePerDay,
            peakValue: bestValue,
            peakAt: DateTime.fromMillis(bestMillis),
          }
        : { direction: "rising", slopePerDay };
  } else if (slopePerDay < -slopeThreshold) {
    trend = { direction: "falling", slopePerDay };
  } else {
    trend = { direction: "steady", slopePerDay };
  }

  // Cycle-aware window peak & trough: max(7d, dosage.interval), clamped to first dose.
  const windowDays = Math.max(7, dosage.interval);
  const rawStartMillis = nowMillis - windowDays * 86_400_000;
  const firstDoseMillis = sorted[0].takenAt.toMillis();
  const clampedToFirstDose = firstDoseMillis > rawStartMillis;
  const windowStartMillis = clampedToFirstDose ? firstDoseMillis : rawStartMillis;

  const sampleMillis = new Set<number>();
  const uniformSteps = 168;
  for (let i = 0; i <= uniformSteps; i++) {
    sampleMillis.add(
      Math.round(windowStartMillis + ((nowMillis - windowStartMillis) * i) / uniformSteps),
    );
  }
  for (const dose of sorted) {
    const t = dose.takenAt.toMillis();
    if (t >= windowStartMillis && t <= nowMillis) {
      sampleMillis.add(t);
      const preT = t - 60_000;
      if (preT >= windowStartMillis) {
        sampleMillis.add(preT);
      }
    }
  }
  sampleMillis.add(nowMillis);

  let windowTrough = Infinity;
  let windowPeak = -Infinity;
  for (const t of sampleMillis) {
    const v = evalAt(t);
    if (v < windowTrough) windowTrough = v;
    if (v > windowPeak) windowPeak = v;
  }

  // Actual dosing cadence across recent intervals, ignoring split doses (< 0.25 * interval).
  const minGapDays = 0.25 * dosage.interval;
  const gaps: number[] = [];
  for (let i = 1; i < sorted.length; i++) {
    const gapDays = (sorted[i].takenAt.toMillis() - sorted[i - 1].takenAt.toMillis()) / 86_400_000;
    if (gapDays >= minGapDays) {
      gaps.push(gapDays);
    }
  }
  const recentGaps = gaps.slice(-6);
  const cadence =
    recentGaps.length > 0
      ? {
          avgIntervalDays: recentGaps.reduce((a, b) => a + b, 0) / recentGaps.length,
          sampleCount: recentGaps.length,
        }
      : undefined;

  return {
    ideal,
    current,
    trend,
    nextDoseAt,
    nextTrough,
    isOverdue,
    window: {
      start: DateTime.fromMillis(windowStartMillis),
      end: now,
      days: windowDays,
      clampedToFirstDose,
      trough: windowTrough,
      peak: windowPeak,
    },
    cadence,
  };
}

// convert a date to an x value
function dateToX(start: DateTime, date: DateTime): number {
  return date.diff(start).as("days");
}
