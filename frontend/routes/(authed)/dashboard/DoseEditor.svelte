<script lang="ts">
  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api.svelte";

  import Dialog from "$lib/components/Dialog.svelte";
  import Icon from "$lib/components/Icon.svelte";
  import InputQuantity from "$lib/components/inputs/InputQuantity.svelte";
  import { DateTime, type DurationLike } from "luxon";
  import { logErrorToast } from "$lib/toasts";

  let {
    open = $bindable(false),
    dose,
    update,
  }: {
    open: boolean;
    dose: e2.Dose;
    // update is called after saving the dose to refresh the parent component.
    update: () => void;
  } = $props();

  let busy = $state(false);
  let confirmDelete = $state(false);

  async function saveDose() {
    try {
      await api.editDose(dose.oldTakenAt, {
        ...dose,
        deliveryMethod: dose.deliveryMethod.id,
        takenAt: dose.takenAt.toISO(),
        takenOffAt: dose.takenOffAt?.toISO(),
      });
      update();
    } catch (err) {
      logErrorToast("Failed to save dose", err);
    }
  }

  async function forgetDose() {
    try {
      await api.forgetDose(dose.oldTakenAt);
      update();
    } catch (err) {
      logErrorToast("Failed to forget dose", err);
    }
  }

  function addToDoseTime(duration: DurationLike) {
    dose.takenAt = dose.takenAt.plus(duration);
  }
</script>

<Dialog bind:open dismissible>
  <header>
    <h2>Edit Dose</h2>
  </header>
  <form>
    <label>
      <span>Taken at</span>
      <input
        type="datetime-local"
        class="mb-2"
        value={dose.takenAt.set({ second: 0, millisecond: 0 }).toISO({ includeOffset: false })}
        onchange={(e) => {
          const newDate = DateTime.fromISO(e.currentTarget.value, { zone: dose!.takenAt.zone });
          if (newDate.isValid) {
            dose.takenAt = newDate;
          }
        }}
      />
      <div class="time-actions flex gap-2 justify-end text-sm">
        <button onclick={() => addToDoseTime({ hour: -1 })}>-1h</button>
        <button onclick={() => addToDoseTime({ hour: -6 })}>-6h</button>
      </div>
    </label>
    <label>
      <span>Amount</span>
      <InputQuantity
        unit={dose.deliveryMethod.units}
        initial={dose.dose}
        onchange={(qty) => {
          if (qty) {
            dose!.dose = qty;
          }
        }}
      />
    </label>
  </form>
  <footer class="mt-0 flex justify-end">
    <button
      aria-label="Cancel"
      class="outline"
      disabled={busy}
      onclick={() => {
        open = false;
      }}
    >
      Cancel <Icon name="close" />
    </button>

    <div class="flex-1"></div>

    <button
      aria-label="Forget dose"
      disabled={busy}
      onclick={async () => {
        if (!confirmDelete) {
          confirmDelete = true;
          return;
        }

        try {
          busy = true;
          await forgetDose();
        } finally {
          confirmDelete = false;
          busy = false;
          open = false;
        }
      }}
    >
      {#if confirmDelete}
        Really? <Icon name="check" />
      {:else}
        Forget <Icon name="delete" />
      {/if}
    </button>

    <button
      aria-label="Save dose"
      class="secondary"
      disabled={busy}
      onclick={async () => {
        try {
          busy = true;
          await saveDose();
        } finally {
          busy = false;
          open = false;
        }
      }}
    >
      Save <Icon name="save" />
    </button>
  </footer>
</Dialog>

<style>
  form label span {
    display: block;
    margin-bottom: calc(var(--pico-spacing) / 4);
  }
</style>
