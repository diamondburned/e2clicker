import type * as api from "#lib/api.ts";
import type { LineData, UTCTimestamp } from "lightweight-charts";
import { DateTime, Interval } from "luxon";
import { fillTargetRange } from "estrannaise/src/models";
import { deliveryMethod, type DeliveryMethod } from "./methods.svelte";

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

let colorNormalizeCtx: CanvasRenderingContext2D | null | undefined;

export function normalizeCssColor(cssColor: string): string {
  const c = cssColor.trim();
  if (!c) return "";

  const hex = /^#([0-9a-f]{3}|[0-9a-f]{6})$/i.exec(c);
  if (hex) {
    return c;
  }

  // Standard 0..255 comma-separated rgb/rgba is already supported by lightweight-charts.
  if (
    /^rgba?\(\s*\d+(?:\.\d+)?\s*,\s*\d+(?:\.\d+)?\s*,\s*\d+(?:\.\d+)?(?:\s*,\s*[\d.]+)?\s*\)$/i.test(
      c,
    )
  ) {
    return c;
  }

  // Convert percentage rgb(r%, g%, b%) / rgba(r%, g%, b%, a) (with commas or spaces).
  const pctMatch =
    /^rgba?\(\s*([\d.]+)%\s*[,\s]\s*([\d.]+)%\s*[,\s]\s*([\d.]+)%(?:\s*[,/]\s*([\d.]+%?))?\s*\)$/i.exec(
      c,
    );
  if (pctMatch) {
    const r = Math.round((parseFloat(pctMatch[1]) / 100) * 255);
    const g = Math.round((parseFloat(pctMatch[2]) / 100) * 255);
    const b = Math.round((parseFloat(pctMatch[3]) / 100) * 255);
    if (pctMatch[4] !== undefined) {
      const a = pctMatch[4].endsWith("%") ? parseFloat(pctMatch[4]) / 100 : parseFloat(pctMatch[4]);
      return `rgba(${r}, ${g}, ${b}, ${a})`;
    }
    return `rgb(${r}, ${g}, ${b})`;
  }

  // Fallback to Canvas 2D color normalization for CSS Color 4 formats (space rgb, oklch, hsl, etc.).
  if (typeof document !== "undefined") {
    if (colorNormalizeCtx === undefined) {
      const canvas = document.createElement("canvas");
      canvas.width = 1;
      canvas.height = 1;
      colorNormalizeCtx = canvas.getContext("2d");
    }
    if (colorNormalizeCtx) {
      colorNormalizeCtx.fillStyle = "#000000";
      colorNormalizeCtx.fillStyle = c;
      const normalized = colorNormalizeCtx.fillStyle;
      if (typeof normalized === "string" && normalized) {
        return normalized;
      }
    }
  }

  return c;
}

export function gatherStyles(element: HTMLElement | null): PlotStyles {
  if (element) {
    const styles = getComputedStyle(element);
    return mapValues(picoVariables, (key, variable) => {
      const raw = styles.getPropertyValue(variable);
      return key === "fontFamily" ? raw : normalizeCssColor(raw);
    });
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

export function withAlpha(cssColor: string, alpha: number): string {
  const c = normalizeCssColor(cssColor);
  if (!c) return "";
  const hex = /^#([0-9a-f]{3}|[0-9a-f]{6})$/i.exec(c);
  if (hex) {
    let digits = hex[1];
    if (digits.length === 3) {
      digits = digits
        .split("")
        .map((ch) => ch + ch)
        .join("");
    }
    const r = parseInt(digits.slice(0, 2), 16);
    const g = parseInt(digits.slice(2, 4), 16);
    const b = parseInt(digits.slice(4, 6), 16);
    return `rgba(${r}, ${g}, ${b}, ${alpha})`;
  }
  const rgb = /^rgba?\(\s*([\d.]+)\s*,\s*([\d.]+)\s*,\s*([\d.]+)(?:\s*,\s*[\d.]+)?\s*\)$/i.exec(c);
  if (rgb) {
    return `rgba(${rgb[1]}, ${rgb[2]}, ${rgb[3]}, ${alpha})`;
  }
  return c;
}

export function dataWithinInterval(
  data: LineData<UTCTimestamp>[],
  iv: Interval,
): LineData<UTCTimestamp>[] {
  return data.filter((d) => iv.contains(DateTime.fromSeconds(d.time)));
}

export type Dose = Omit<api.Dose, "deliveryMethod" | "takenAt" | "takenOffAt"> & {
  deliveryMethod: DeliveryMethod;
  takenAt: DateTime<true>;
  takenOffAt?: DateTime<true>;
};

export type DosageHistory = Dose[];

// convertDoseHistory converts the dosage history from the API to have proper
// DateTime objects, sorted chronologically ascending by takenAt.
export function convertDoseHistory(history: api.DosageHistory): DosageHistory {
  return history
    .map(
      (dose) =>
        ({
          ...dose,
          deliveryMethod: deliveryMethod(dose.deliveryMethod)!,
          takenAt: mustDateTimeFromISO(dose.takenAt),
          takenOffAt: dose.takenOffAt ? mustDateTimeFromISO(dose.takenOffAt) : undefined,
        }) as Dose,
    )
    .sort((a, b) => a.takenAt.toMillis() - b.takenAt.toMillis());
}

function mustDateTimeFromISO(iso: string): DateTime<true> {
  const dt = DateTime.fromISO(iso);
  if (!dt.isValid) {
    throw new Error(`invalid ISO date ${iso}: ${dt.invalidExplanation}`);
  }
  return dt;
}
