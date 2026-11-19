import assert from "node:assert/strict";
import test from "node:test";
import { DateTime, Interval } from "luxon";
import { PKParameters, terminalEliminationTime3C } from "estrannaise/src/models";
import type { DeliveryMethod } from "estrannaise/src/modeldata";
import {
  computeActualWithUncertainty,
  maxEliminationDays,
  sampleSplitTimestamps,
  sortedQuantile,
  type PreparedDose,
} from "./uncertainty.ts";

const ALL_MODELS = Object.keys(PKParameters) as DeliveryMethod[];

test("terminalEliminationTime3C and maxEliminationDays are finite positive for all 8 models", () => {
  assert.equal(ALL_MODELS.length, 8);
  for (const model of ALL_MODELS) {
    const t5 = terminalEliminationTime3C(...PKParameters[model]);
    const tMax = maxEliminationDays(model);
    assert.ok(Number.isFinite(t5) && t5 > 0, `${model} terminal5HL=${t5}`);
    assert.ok(Number.isFinite(tMax) && tMax >= t5, `${model} maxElim=${tMax}`);
  }
});

test("sortedQuantile handles edge cases and interpolation", () => {
  assert.equal(sortedQuantile(new Float64Array([]), 0.5), 0);
  assert.equal(sortedQuantile(new Float64Array([42]), 0.05), 42);
  const arr = new Float64Array([10, 20, 30, 40, 50]);
  assert.equal(sortedQuantile(arr, 0), 10);
  assert.equal(sortedQuantile(arr, 0.5), 30);
  assert.equal(sortedQuantile(arr, 0.25), 20);
  assert.equal(sortedQuantile(arr, 1), 50);
});

test("computeActualWithUncertainty handles empty history", () => {
  const end = DateTime.fromISO("2026-11-19T00:00:00Z");
  const start = end.minus({ days: 7 });
  const iv = Interval.fromDateTimes(start, end) as Interval<true>;
  const pred = Interval.fromDateTimes(end, end.plus({ days: 7 })) as Interval<true>;
  const { times } = sampleSplitTimestamps(iv, pred, 40);

  const res = computeActualWithUncertainty({
    timestamps: times,
    doses: [],
    conversionFactor: 1.0,
    includeUncertainty: true,
  });

  assert.equal(res.actual.length, times.length);
  assert.equal(res.p05.length, times.length);
  assert.equal(res.p50.length, times.length);
  assert.equal(res.p95.length, times.length);
  for (let i = 0; i < times.length; i++) {
    assert.equal(res.actual[i].value, 0);
    assert.equal(res.p05[i].value, 0);
    assert.equal(res.p50[i].value, 0);
    assert.equal(res.p95[i].value, 0);
  }
});

test("computeActualWithUncertainty is deterministic, ordered, and captures EUn heavy tail + patches", () => {
  const now = DateTime.fromISO("2026-11-19T00:00:00Z");

  // EUn im 5 mg every 7 days over 220 days -> steady state around MAP ~148, P95 ~450.
  const eunDoses: PreparedDose[] = [];
  for (let d = 220; d >= 0; d -= 7) {
    eunDoses.push({
      takenAtMillis: now.minus({ days: d }).toMillis(),
      dose: 5,
      delivery: "EUn im",
    });
  }

  const iv = Interval.fromDateTimes(now.minus({ days: 14 }), now) as Interval<true>;
  const pred = Interval.fromDateTimes(now, now.plus({ days: 7 })) as Interval<true>;
  const { times, splitTime } = sampleSplitTimestamps(iv, pred, 80);

  // Strict ascending timestamps and exact boundary inclusion.
  for (let i = 1; i < times.length; i++) {
    assert.ok(times[i].toUnixInteger() > times[i - 1].toUnixInteger());
  }
  assert.ok(times.some((t) => t.toUnixInteger() === splitTime));

  const run1 = computeActualWithUncertainty({
    timestamps: times,
    doses: eunDoses,
    conversionFactor: 1.0,
    includeUncertainty: true,
  });
  const run2 = computeActualWithUncertainty({
    timestamps: times,
    doses: eunDoses,
    conversionFactor: 1.0,
    includeUncertainty: true,
  });

  assert.deepEqual(run1, run2);

  for (let i = 0; i < times.length; i++) {
    assert.equal(run1.actual[i].time, run1.p05[i].time);
    assert.equal(run1.actual[i].time, run1.p95[i].time);
    assert.ok(run1.p05[i].value >= 0);
    assert.ok(run1.p05[i].value <= run1.p50[i].value + 1e-9);
    assert.ok(run1.p50[i].value <= run1.p95[i].value + 1e-9);
  }

  const lastIdx = times.length - 1;
  assert.ok(run1.p95[lastIdx].value > run1.actual[lastIdx].value * 2);

  // Patch tw with explicit wearDuration.
  const patchDoses: PreparedDose[] = [
    {
      takenAtMillis: now.minus({ days: 7 }).toMillis(),
      dose: 100,
      delivery: "patch tw",
      wearDuration: 3.5,
    },
    {
      takenAtMillis: now.minus({ days: 3.5 }).toMillis(),
      dose: 100,
      delivery: "patch tw",
      wearDuration: 3.5,
    },
  ];
  const patchRes = computeActualWithUncertainty({
    timestamps: times,
    doses: patchDoses,
    conversionFactor: 1.0,
    includeUncertainty: true,
  });
  for (let i = 0; i < times.length; i++) {
    assert.ok(patchRes.p05[i].value >= 0);
    assert.ok(patchRes.p05[i].value <= patchRes.p50[i].value + 1e-9);
    assert.ok(patchRes.p50[i].value <= patchRes.p95[i].value + 1e-9);
  }
});
