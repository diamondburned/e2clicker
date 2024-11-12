import type {
  LineData,
  IChartApi,
  LineSeriesPartialOptions,
  UTCTimestamp,
} from "lightweight-charts";
import { DateTime, Duration } from "luxon";
import { fillCurve } from "estrannaise/src/models";

export const picoVariables = {
  color: "--pico-color",
  primary: "--pico-primary",
  secondary: "--pico-secondary",
  muted: "--pico-muted-color",
  fontFamily: "--pico-font-family",
} as const;

export function gatherStyles(element: HTMLElement | null): {
  [key in keyof typeof picoVariables]: string;
} {
  if (element) {
    const styles = getComputedStyle(element);
    return mapValues(picoVariables, (_, variable) => styles.getPropertyValue(variable));
  } else {
    return mapValues(picoVariables, () => "");
  }
}

function mapValues<T extends Record<string, any>, U>(
  obj: T,
  fn: (key: keyof T, value: T[keyof T]) => U,
): Record<keyof T, U> {
  return Object.fromEntries(
    Object.entries(obj).map(([key, value]) => [key, fn(key, value)]),
  ) as Record<keyof T, U>;
}

export function ensureLine(
  data: LineData[],
  chart: IChartApi,
  options: LineSeriesPartialOptions,
): () => void {
  const line = chart.addLineSeries(options);
  line.setData(data);
  chart.timeScale().fitContent();
  return () => chart.removeSeries(line);
}

const plotPoints = 500;

export function fillE2Data(start: DateTime, end: DateTime, f: (t: number) => number): LineData[] {
  const xMin = 0;
  const xMax = end.diff(start).as("days");
  return fillCurve(f, xMin, xMax, plotPoints).map((p) => ({
    time: start.plus(Duration.fromObject({ days: p.Time })).toUnixInteger() as UTCTimestamp,
    value: p.E2,
  }));
}
