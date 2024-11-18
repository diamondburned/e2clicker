import * as api from "$lib/api";
import {
  modelList,
  type DeliveryMethod as EstrannaiseDeliveryMethod,
} from "estrannaise/src/modeldata";
import { untrack } from "svelte";

export type DeliveryMethod = api.DeliveryMethod & {
  patch?: true;
};

let deliveryMethods = $state<DeliveryMethod[]>([]);

export const deliveryMethodsList = () => deliveryMethods;

// deliveryMethod returns the delivery method with the given ID.
export const deliveryMethod = (id: string) => {
  if (!deliveryMethods) {
    throw new Error("Delivery methods not loaded yet");
  }
  return deliveryMethods.find((m) => m.id === id);
};

export async function updateDeliveryMethods() {
  const current = untrack(() => deliveryMethods);

  // Restore cached delivery methods from local storage if we have one.
  const data = localStorage.getItem("e2clicker-deliveryMethods-v1");
  if (data) {
    const cached = JSON.parse(data) as DeliveryMethod[];
    if (!deliveryMethodsEq(current, cached)) {
      deliveryMethods = cached;
    }
  }

  try {
    const latest = await fetchDeliveryMethods();
    if (!deliveryMethodsEq(current, latest)) {
      localStorage.setItem("e2clicker-deliveryMethods-v1", JSON.stringify(latest));
      deliveryMethods = latest;
    }
  } catch (err) {
    console.warn("Failed to sync delivery methods with server:", err);
  }
}

async function fetchDeliveryMethods() {
  const methods = (await api.deliveryMethods()).map((method) => {
    return {
      ...method,
      patch: method.id.startsWith("patch") ? true : undefined,
    } as DeliveryMethod;
  });
  return methods;
}

function deliveryMethodsEq(a: DeliveryMethod[], b: DeliveryMethod[]) {
  const eq = (a: DeliveryMethod, b: DeliveryMethod) =>
    a.id === b.id && a.name === b.name && a.description === b.description && a.patch === b.patch;
  return a.length === b.length && a.every((m, i) => eq(m, b[i]));
}

export function estrannaiseDeliveryMethod(method: string): EstrannaiseDeliveryMethod {
  if (!modelList[method as EstrannaiseDeliveryMethod]) {
    throw new Error(`Unknown delivery method: ${method}`);
  }
  return method as EstrannaiseDeliveryMethod;
}
