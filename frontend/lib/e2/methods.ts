import * as api from "$lib/api";
import {
  type DeliveryMethod as EstrannaiseDeliveryMethod,
  type PKParameter,
} from "estrannaise/src/modeldata";
import { PKFunctions, PKParameters, PKRandomFunctions } from "estrannaise/src/models";

// DeliveryMethods is a list of all available delivery methods.
// Command: GET /api/delivery-methods
export const deliveryMethods: (api.DeliveryMethod & { patch?: true })[] = [
  {
    id: "EB im",
    name: "Estradiol Benzoate, Intramuscular",
    units: "mg",
  },
  {
    id: "EV im",
    name: "Estradiol Valerate, Intramuscular",
    units: "mg",
  },
  {
    id: "EEn im",
    name: "Estradiol Enanthate, Intramuscular",
    units: "mg",
  },
  {
    id: "EC im",
    name: "Estradiol Cypionate, Intramuscular",
    units: "mg",
  },
  {
    id: "EUn im",
    name: "Estradiol Undecylate, Intramuscular",
    units: "mg",
  },
  {
    id: "EUn casubq",
    name: "Estradiol Undecylate in Castor oil, Subcutaneous",
    units: "mg",
  },
  {
    id: "patch",
    name: "Patch",
    units: "mcg/day",
    patch: true,
  },
] as const;

// deliveryMethod returns the delivery method with the given ID.
export const deliveryMethod = (id: string) => deliveryMethods.find((m) => m.id === id);

type OurDeliveryMethod = (typeof deliveryMethods)[number]["id"];

export const estrannaiseDeliveryMethods: {
  [key in OurDeliveryMethod]: EstrannaiseDeliveryMethod;
} = {
  "EB im": "EB im",
  "EV im": "EV im",
  "EEn im": "EEn im",
  "EC im": "EC im",
  "EUn im": "EUn im",
  "EUn casubq": "EUn casubq",
  patch: "patch ow",
};

export function getFunctions(
  dosage?: api.Dosage,
  conversionFactor = 1.0,
): {
  pk: (t: number) => number;
  pkRandom: (t: number) => number;
  pkParams: PKParameter;
} {
  const estrannaiseDelivery = dosage && estrannaiseDeliveryMethods[dosage.deliveryMethod];
  if (!estrannaiseDelivery) {
    return {
      pk: () => 0,
      pkRandom: () => 0,
      pkParams: [0, 0, 0, 0],
    };
  }

  const pk = PKFunctions(conversionFactor)[estrannaiseDelivery];
  const pkRandom = PKRandomFunctions(conversionFactor)[estrannaiseDelivery];
  const pkParams = PKParameters[estrannaiseDelivery];

  return {
    pk: (t: number) => pk(t, dosage.dose, true, dosage.interval),
    pkRandom: (t: number) => pkRandom(t, dosage.dose, true, dosage.interval),
    pkParams,
  };
}

// Background task to assert that all delivery methods are up to date.
async () => {
  const updatedMethods = await api.deliveryMethods();

  for (const method of updatedMethods) {
    const knownMethod = deliveryMethods.find((m) => m.id === method.id);
    if (!knownMethod) {
      console.warn(
        `e2conversions: Delivery method ${method.id} is not in the list of known methods.`,
      );
      continue;
    }

    if (knownMethod.name !== method.name || knownMethod.units !== method.units) {
      console.warn(`e2conversions: Delivery method ${method.id} has changed upstream.`);
    }
  }
};
