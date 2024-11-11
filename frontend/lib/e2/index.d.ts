declare module "estrannaise/src/models" {
  import { modelList, PKParameters, availableUnits } from "estrannaise/src/modeldata";
  export { modelList, PKParameters, availableUnits };

  import type { DeliveryMethod } from "estrannaise/src/modeldata";

  // Generate the "curve" for target mean levels for transfeminine HRT,
  // based on WPATH SOC 8 + Endocrine Society Guidelines.
  export function fillTargetRange(
    xMin: number,
    xMax: number,
    // Conversion factor between units
    conversionFactor?: number,
  ): { time: number; lower: number; upper: number }[];

  // Convenience function for when iterating over model functions
  export function fillCurve(
    func: (i: number) => number,
    xMin: number,
    xMax: number,
    numPoints: number,
  ): { Time: number; E2: number }[];

  export function PKFunctions(conversionFactor?: number): {
    [key in DeliveryMethod]: (t: number, dose: number, steadystate?: boolean, T?: number) => number;
  };

  export function PKRandomFunctions(conversionFactor?: number): {
    [key in DeliveryMethod]: (
      t: number, // time offset for dose calculation
      dose: number, // Dose amounts, in mg
      steadystate?: boolean,
      T?: number,
      idx?: number,
    ) => number;
  };

  export function e2MultiDose3C(
    t: number,
    doses?: number[],
    times?: number[],
    models?: string[],
    cf?: number,
    random?: boolean | number,
    intervals?: boolean,
  ): number;

  export function e2ssAverage3C(
    dose: number,
    T: number,
    d: number,
    k1: number,
    k2: number,
    k3: number,
  ): number;

  // This is an approximation, but it's good enough for our purposes
  export function terminalEliminationTime3C(model: string, nbHalfLives?: number): number;

  // The underlying function for getPKQuantities.
  export function getPKQuantities3C(
    d: number,
    k1: number,
    k2: number,
    k3: number,
  ): { Tmax: number; Cmax: number; halfLife: number };

  // Get the PK quantities for a given model.
  export function getPKQuantities(model: string): ReturnType<typeof getPKQuantities3C>;
}

declare module "estrannaise/src/modeldata" {
  export type DeliveryMethod =
    | "EV im"
    | "EEn im"
    | "EC im"
    | "EUn im"
    | "EUn casubq"
    | "EB im"
    | "patch tw"
    | "patch ow";

  export type LevelUnits = "pg/mL" | "pmol/L" | "ng/L" | "FFF";

  export const menstrualCycleData: {
    t: number[]; // day in month, [0..29]
    E2: number[]; // same length
    E2p5: number[]; // same length
    E2p95: number[]; // same length
  };

  export const availableUnits: {
    [key in LevelUnits]: {
      units: string;
      conversionFactor: number;
      precision: number;
    };
  };

  export const modelList: {
    [key in DeliveryMethod]: {
      units: string;
      description: string;
    };
  };

  export type PKParameter = [number, number, number, number];

  export const PKParameters: {
    [key in DeliveryMethod]: PKParameter;
  };

  export let mcmcSamplesPK: {
    [key in DeliveryMethod]: number[][];
  };
}

declare module "estrannaise/src/plotting" {
  export const WONG_PALETTE: string[];

  export type PlottingOptions = {
    menstrualCycleVisible: boolean;
    targetRangeVisible: boolean;
    units: string;
    strokeWidth: number;
    numberOfLinePoints: number;
    numberOfCloudPoints: number;
    pointCloudSize: number;
    pointCloudOpacity: number;
    currentColorscheme: string;
    backgroundColor: string;
    strongForegroundColor: string;
    softForegroundColor: string;
    fontSize: string;
    aspectRatio: number;
  };

  export function wongPalette(n: number): string;

  export function generatePlottingOptions(options?: Partial<PlottingOptions>): PlottingOptions;

  export function plotCurves(
    dataset: any,
    options?: ReturnType<typeof generatePlottingOptions>,
    returnSVG?: boolean,
  ): string | any;
}
