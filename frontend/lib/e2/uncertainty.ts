import type * as charts from "lightweight-charts";
import { e2Curve3C, e2Patch3C, terminalEliminationTime3C } from "estrannaise/src/models";
import {
  mcmcSamplesPK,
  PKParameters,
  type DeliveryMethod as EstrannaiseDeliveryMethod,
  type PKParameter,
} from "estrannaise/src/modeldata";
import { DateTime, type Interval } from "luxon";

const MS_PER_DAY = 86_400_000;

export type PreparedDose = {
  takenAtMillis: number;
  dose: number;
  delivery: EstrannaiseDeliveryMethod;
  wearDuration?: number;
};

export type ActualUncertaintySeries = {
  actual: charts.LineData<charts.UTCTimestamp>[];
  p05: charts.LineData<charts.UTCTimestamp>[];
  p50: charts.LineData<charts.UTCTimestamp>[];
  p95: charts.LineData<charts.UTCTimestamp>[];
};

const maxElimDaysCache = new Map<EstrannaiseDeliveryMethod, number>();

export function maxEliminationDays(delivery: EstrannaiseDeliveryMethod): number {
  const cached = maxElimDaysCache.get(delivery);
  if (cached !== undefined) return cached;

  let maxDays = terminalEliminationTime3C(...PKParameters[delivery], 8);
  for (const params of mcmcSamplesPK[delivery]) {
    const t = terminalEliminationTime3C(...params, 6);
    if (t > maxDays) maxDays = t;
  }
  maxElimDaysCache.set(delivery, maxDays);
  return maxDays;
}

export function evalDose(
  xDays: number,
  doseMgCf: number,
  params: PKParameter,
  wearDuration?: number,
): number {
  if (xDays < 0) return 0;
  const v =
    wearDuration !== undefined
      ? e2Patch3C(xDays, doseMgCf, ...params, wearDuration, false, 0.0)
      : e2Curve3C(xDays, doseMgCf, ...params, 0.0, 0.0, false, 0.0);
  return v > 0 ? v : 0;
}

export function sortedQuantile(sorted: Float64Array, q: number): number {
  if (sorted.length === 0) return 0;
  if (sorted.length === 1) return sorted[0];
  const pos = q * (sorted.length - 1);
  const lo = Math.floor(pos);
  const hi = Math.ceil(pos);
  if (lo === hi) return sorted[lo];
  const frac = pos - lo;
  return sorted[lo] * (1 - frac) + sorted[hi] * frac;
}

export function sampleSplitTimestamps(
  iv: Interval<true>,
  predictedIv: Interval<true>,
  totalPoints: number,
): { times: DateTime[]; splitTime: charts.UTCTimestamp } {
  const totalMs = Math.max(1, predictedIv.end.toMillis() - iv.start.toMillis());
  const pastRatio = Math.max(0, iv.end.toMillis() - iv.start.toMillis()) / totalMs;
  const count = Math.max(4, totalPoints);
  const pastCount = Math.max(2, Math.round(count * pastRatio));
  const predCount = Math.max(2, count - pastCount + 1);

  const times: DateTime[] = [];
  let lastTime = -Infinity;

  const pushInterval = (seg: Interval<true>, n: number) => {
    const startMs = seg.start.toMillis();
    const spanMs = seg.end.toMillis() - startMs;
    for (let i = 0; i < n; i++) {
      const t = DateTime.fromMillis(startMs + (spanMs * i) / (n - 1));
      const sec = t.toUnixInteger();
      if (sec > lastTime) {
        times.push(t);
        lastTime = sec;
      }
    }
  };

  pushInterval(iv, pastCount);
  pushInterval(predictedIv, predCount);

  return {
    times,
    splitTime: iv.end.toUnixInteger() as charts.UTCTimestamp,
  };
}

export function computeActualWithUncertainty({
  timestamps,
  doses,
  conversionFactor,
  includeUncertainty,
}: {
  timestamps: DateTime[];
  doses: PreparedDose[];
  conversionFactor: number;
  includeUncertainty: boolean;
}): ActualUncertaintySeries {
  const prepared = doses.map((d) => ({
    takenAtMillis: d.takenAtMillis,
    doseMgCf: conversionFactor * d.dose,
    mapParams: PKParameters[d.delivery],
    samples: mcmcSamplesPK[d.delivery],
    wearDuration: d.wearDuration,
    maxElimDays: maxEliminationDays(d.delivery),
  }));

  const numDraws =
    prepared.length > 0
      ? prepared.reduce((max, d) => Math.max(max, d.samples.length), 0)
      : mcmcSamplesPK["EV im"].length;
  const draws = includeUncertainty ? new Float64Array(numDraws) : undefined;

  const actual: charts.LineData<charts.UTCTimestamp>[] = [];
  const p05: charts.LineData<charts.UTCTimestamp>[] = [];
  const p50: charts.LineData<charts.UTCTimestamp>[] = [];
  const p95: charts.LineData<charts.UTCTimestamp>[] = [];

  const activeXDays: number[] = [];
  const activeIdx: number[] = [];

  for (let ti = 0; ti < timestamps.length; ti++) {
    const t = timestamps[ti];
    const tMillis = t.toMillis();
    const time = t.toUnixInteger() as charts.UTCTimestamp;

    activeXDays.length = 0;
    activeIdx.length = 0;

    let mapSum = 0;
    for (let di = 0; di < prepared.length; di++) {
      const d = prepared[di];
      const xDays = (tMillis - d.takenAtMillis) / MS_PER_DAY;
      if (xDays < 0 || xDays > d.maxElimDays) continue;
      mapSum += evalDose(xDays, d.doseMgCf, d.mapParams, d.wearDuration);
      if (draws) {
        activeXDays.push(xDays);
        activeIdx.push(di);
      }
    }

    actual.push({ time, value: mapSum });

    if (draws) {
      if (activeIdx.length === 0) {
        p05.push({ time, value: 0 });
        p50.push({ time, value: 0 });
        p95.push({ time, value: 0 });
      } else {
        draws.fill(0);
        for (let ai = 0; ai < activeIdx.length; ai++) {
          const d = prepared[activeIdx[ai]];
          const xDays = activeXDays[ai];
          const samples = d.samples;
          const nSamples = samples.length;
          for (let k = 0; k < numDraws; k++) {
            draws[k] += evalDose(xDays, d.doseMgCf, samples[k % nSamples], d.wearDuration);
          }
        }
        draws.sort();
        p05.push({ time, value: sortedQuantile(draws, 0.05) });
        p50.push({ time, value: sortedQuantile(draws, 0.5) });
        p95.push({ time, value: sortedQuantile(draws, 0.95) });
      }
    }
  }

  return { actual, p05, p50, p95 };
}
