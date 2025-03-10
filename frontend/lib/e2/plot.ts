import type * as api from "$lib/api";
import type { LineData, UTCTimestamp } from "lightweight-charts";
import { DateTime, Interval } from "luxon";
import { fillTargetRange } from "estrannaise/src/models";
import { deliveryMethod } from "./methods.svelte";

// style variables used in the plot.
export const picoVariables = {
  color: "--pico-color",
  primary: "--pico-primary",
  secondary: "--pico-secondary",
  muted: "--pico-muted-color",
  mutedBorder: "--pico-muted-border-color",
  fontFamily: "--pico-font-family",
} as const;

export type PlotStyles = {
  [key in keyof typeof picoVariables]: string;
};

export function gatherStyles(element: HTMLElement | null): PlotStyles {
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

// wpathRange returns the range of the WPATH standards for estradiol levels.
export function wpathRange(conversionFactor: number): {
  lower: number;
  upper: number;
} {
  const filled = fillTargetRange(0, 1, conversionFactor);
  return {
    lower: filled[0].lower,
    upper: filled[0].upper,
  };
}

export function dataWithinInterval(
  data: LineData<UTCTimestamp>[],
  iv: Interval,
): LineData<UTCTimestamp>[] {
  return data.filter((d) => iv.contains(DateTime.fromSeconds(d.time)));
}

export type Dose = Omit<api.Dose, "deliveryMethod" | "takenAt" | "takenOffAt"> & {
  deliveryMethod: api.DeliveryMethod;
  takenAt: DateTime<true>;
  takenOffAt?: DateTime<true>;
  _takenAt: string;
  _takenOffAt?: string;
};

export type DosageHistory = Dose[];

// convertDoseHistory converts the dosage history from the API to have proper
// DateTime objects.
export function convertDoseHistory(history: api.DosageHistory): DosageHistory {
  return history.map(
    (dose) =>
      ({
        ...dose,
        deliveryMethod: deliveryMethod(dose.deliveryMethod)!,
        takenAt: mustDateTimeFromISO(dose.takenAt),
        takenOffAt: dose.takenOffAt ? mustDateTimeFromISO(dose.takenOffAt) : undefined,
        _takenAt: dose.takenAt,
        _takenOffAt: dose.takenOffAt,
      }) as Dose,
  );
}

function mustDateTimeFromISO(iso: string): DateTime<true> {
  const dt = DateTime.fromISO(iso);
  if (!dt.isValid) {
    throw new Error(`invalid ISO date ${iso}: ${dt.invalidExplanation}`);
  }
  return dt;
}
