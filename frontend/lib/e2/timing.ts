import * as api from "$lib/api";
import { DateTime, Duration, type DateTimeMaybeValid, type DurationObjectUnits } from "luxon";

// timeUntilNextDose calculates the time at which the next dose should be taken
// based on the last dose and the dosage schedule.
export function timeUntilNextDose(dosage: api.Dosage, lastDose?: DateTimeMaybeValid): DateTime {
  if (!lastDose || !lastDose.isValid) {
    return DateTime.now();
  }
  const interval = Duration.fromObject({ days: dosage.interval });
  return lastDose.plus(interval);
}

export function formatDuration(duration: Duration): string {
  duration = duration.rescale();
  duration = durationTrimLower(duration);
  return duration.toHuman({
    listStyle: "long",
    unitDisplay: "long",
    roundingMode: "ceil",
    roundingIncrement: 1,
    roundingPriority: "lessPrecision",
    maximumFractionDigits: 0,
  });
}

export const relevantUnits: (keyof DurationObjectUnits)[] = [
  "years",
  "months",
  "weeks",
  "days",
  "hours",
  "minutes",
  "seconds",
] as const;

// Trim the lower bound of a duration to the nearest whole unit.
// The number of units to keep is specified by the numUnits argument.
// For example, if a duration is 1 year, 2 months, 3 days, and 4 hours,
// and numUnits = 2, the result will be 1 year and 2 months.
export function durationTrimLower(duration: Duration, numUnits = 2): Duration {
  const obj = duration.toObject();
  for (let i = 0; i < relevantUnits.length; i++) {
    const units = relevantUnits.slice(i, Math.min(i + numUnits, relevantUnits.length));
    const values = units.map((unit) => obj[unit] ?? 0);
    if (values[0] > 0) {
      // This slice of units is fully not empty, so we can keep it.
      const zipped = units.map((unit, i) => [unit, values[i]]);
      const trimmed = Duration.fromObject(Object.fromEntries(zipped));
      return trimmed;
    }
  }
  return duration;
}

// Rescale duration so that if there is 2 or more units, it tries to rescale it
// to only 1 unit.
export function roundDuration(
  duration: Duration,
  {
    allowDecimals = false,
  }: {
    allowDecimals?: boolean;
  } = {},
): Duration {
  const dobj = duration.toObject();
  const keys = Object.keys(dobj) as (keyof DurationObjectUnits)[];
  if (keys.length < 2) {
    return duration; // nothing to do
  }
  return duration.shiftTo(allowDecimals ? keys[0] : keys[keys.length - 1]);
}

// Format the given dose time for display. The returned string will be quite
// long.
export function formatDoseTime(dose: api.DosageObservation, now = DateTime.now()): string {
  const then = DateTime.fromISO(dose.takenAt);
  return formatDuration(durationTrimLower(now.diff(then).rescale(), 2));
}
