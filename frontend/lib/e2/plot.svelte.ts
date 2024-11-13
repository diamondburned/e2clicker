import * as api from "$lib/api";
import type { LevelUnits } from "estrannaise/src/modeldata";
import { availableUnits } from "estrannaise/src/models";
import { getFunctions } from "./methods";

export * from "./plot";

export let units = $state<LevelUnits>("pg/mL");
export let precision = () => availableUnits[units].precision;
export let conversionFactor = () => availableUnits[units].conversionFactor;
export let functions = (dosage?: api.Dosage) => getFunctions(dosage, conversionFactor());
