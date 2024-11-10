import * as api from "$lib/api";
import { deliveryMethods } from "$lib/api";

// DeliveryMethods is a list of all available delivery methods.
// Command: GET /api/delivery-methods
export const DeliveryMethods: (api.DeliveryMethod & { patch?: true })[] = [
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

// Background task to assert that all delivery methods are up to date.
async () => {
  const updatedMethods = await deliveryMethods();

  for (const method of updatedMethods) {
    const knownMethod = DeliveryMethods.find((m) => m.id === method.id);
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
