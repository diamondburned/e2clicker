<script lang="ts">
  import * as e2 from "#lib/e2.svelte.ts";
  import * as api from "#lib/api.svelte.ts";

  import Dialog from "#lib/components/Dialog.svelte";
  import Icon from "#lib/components/Icon.svelte";
  import InputQuantity from "#lib/components/inputs/InputQuantity.svelte";
  import LevelsPreview from "#lib/components/LevelsPreview.svelte";
  import { DateTime, type DurationLike } from "luxon";
  import { logErrorToast } from "#lib/toasts.ts";
  import { untrack } from "svelte";

  let {
    open = $bindable(false),
    mode = "edit",
    dose,
    initialTakenAt,
    doses,
    update,
    onselectdose,
  }: {
    open: boolean;
    mode?: "create" | "edit";
    dose?: e2.Dose;
    initialTakenAt?: DateTime;
    doses: {
      dosage?: api.Dosage;
      history?: e2.DosageHistory;
    };
    // update is called after saving or forgetting a dose to refresh the parent component.
    update: (savedTakenAt?: DateTime) => void;
    onselectdose?: (dose: e2.Dose) => void;
  } = $props();

  type DraftDose = {
    id?: string;
    deliveryMethod: e2.DeliveryMethod;
    dose: number | undefined;
    takenAt: DateTime<true>;
    takenOffAt?: DateTime<true>;
    comment: string;
  };

  function defaultDeliveryMethod(): e2.DeliveryMethod {
    return (
      (doses.dosage?.deliveryMethod && e2.deliveryMethod(doses.dosage.deliveryMethod)) ||
      e2.deliveryMethodsList()[0]
    );
  }

  function createInitialDraft(): DraftDose {
    if (mode === "edit" && dose) {
      return {
        id: dose.id,
        deliveryMethod: dose.deliveryMethod,
        dose: dose.dose,
        takenAt: dose.takenAt,
        takenOffAt: dose.takenOffAt,
        comment: dose.comment ?? "",
      };
    }

    const baseTime = initialTakenAt && initialTakenAt.isValid ? initialTakenAt : DateTime.now();
    return {
      deliveryMethod: defaultDeliveryMethod(),
      dose: doses.dosage?.dose ?? 4,
      takenAt: baseTime.set({ second: 0, millisecond: 0 }) as DateTime<true>,
      takenOffAt: undefined,
      comment: "",
    };
  }

  let draft = $state<DraftDose>(untrack(createInitialDraft));
  let busy = $state(false);
  let confirmDelete = $state(false);

  $effect(() => {
    if (!open) return;
    mode;
    dose?.id;
    initialTakenAt;
    untrack(() => {
      draft = createInitialDraft();
      confirmDelete = false;
    });
  });

  const sortedHistory = $derived(doses.history ?? []);
  const currentIndex = $derived(
    mode === "edit" && draft.id ? sortedHistory.findIndex((d) => d.id === draft.id) : -1,
  );
  const prevDose = $derived(currentIndex > 0 ? sortedHistory[currentIndex - 1] : undefined);
  const nextDose = $derived(
    currentIndex >= 0 && currentIndex < sortedHistory.length - 1
      ? sortedHistory[currentIndex + 1]
      : undefined,
  );

  const showTakenOffAt = $derived(
    Boolean(draft.deliveryMethod?.patch || draft.takenOffAt !== undefined),
  );

  const validationError = $derived.by(() => {
    if (!draft.dose || !isFinite(draft.dose) || draft.dose <= 0) {
      return "Amount must be greater than 0.";
    }
    if (draft.takenOffAt && draft.takenOffAt < draft.takenAt) {
      return "Taken off time cannot be before taken time.";
    }
    return undefined;
  });

  const draftPreviewDose = $derived.by((): e2.Dose | undefined => {
    if (!draft.dose || !isFinite(draft.dose) || draft.dose <= 0) {
      return undefined;
    }
    return {
      id: draft.id ?? "draft",
      deliveryMethod: draft.deliveryMethod,
      dose: draft.dose,
      takenAt: draft.takenAt,
      takenOffAt: draft.takenOffAt,
      comment: draft.comment || undefined,
    };
  });

  const beforeHistory = $derived(sortedHistory);
  const afterHistory = $derived.by(() => {
    if (!draftPreviewDose) return beforeHistory;
    const base =
      mode === "edit" && draft.id ? beforeHistory.filter((d) => d.id !== draft.id) : beforeHistory;
    return [...base, draftPreviewDose].sort((a, b) => a.takenAt.toMillis() - b.takenAt.toMillis());
  });

  function formatDatetimeLocal(dt: DateTime): string {
    return dt.set({ second: 0, millisecond: 0 }).toISO({ includeOffset: false }) ?? "";
  }

  function addToDoseTime(duration: DurationLike) {
    draft.takenAt = draft.takenAt.plus(duration) as DateTime<true>;
  }

  async function saveDose() {
    if (validationError || !draft.dose) return;

    const trimmedComment = draft.comment.trim();
    if (mode === "create") {
      await api.recordDose({
        deliveryMethod: draft.deliveryMethod.id,
        dose: draft.dose,
        takenAt: draft.takenAt.toISO(),
        takenOffAt: draft.takenOffAt?.toISO(),
        comment: trimmedComment || undefined,
      });
    } else if (draft.id) {
      await api.editDose(draft.id, {
        deliveryMethod: draft.deliveryMethod.id,
        dose: draft.dose,
        takenAt: draft.takenAt.toISO(),
        takenOffAt: draft.takenOffAt?.toISO(),
        comment: trimmedComment,
      });
    }
    update(draft.takenAt);
  }

  async function forgetDose() {
    if (!draft.id) return;
    await api.forgetDose(draft.id);
    update();
  }
</script>

<Dialog bind:open wide dismissible>
  <header class="flex items-center justify-between gap-2">
    <h2 class="m-0">{mode === "create" ? "Log Dose" : "Edit Dose"}</h2>

    {#if mode === "edit" && (prevDose || nextDose)}
      <div class="nav-buttons flex items-center gap-1">
        <button
          type="button"
          class="outline secondary p-1"
          aria-label="Previous dose"
          disabled={busy || !prevDose}
          onclick={() => prevDose && onselectdose?.(prevDose)}
        >
          <Icon name="chevron-left" />
        </button>
        <button
          type="button"
          class="outline secondary p-1"
          aria-label="Next dose"
          disabled={busy || !nextDose}
          onclick={() => nextDose && onselectdose?.(nextDose)}
        >
          <Icon name="chevron-right" />
        </button>
      </div>
    {/if}
  </header>

  <form
    onsubmit={async (ev) => {
      ev.preventDefault();
      if (validationError || busy) return;
      try {
        busy = true;
        await saveDose();
        open = false;
      } catch (err) {
        logErrorToast("Failed to save dose", err);
      } finally {
        busy = false;
      }
    }}
  >
    <div class="editor-grid">
      <label>
        <span>Delivery method</span>
        <select
          value={draft.deliveryMethod.id}
          onchange={(e) => {
            const nextMethod = e2.deliveryMethod(e.currentTarget.value);
            if (!nextMethod) return;
            draft.deliveryMethod = nextMethod;
            if (!nextMethod.patch) {
              draft.takenOffAt = undefined;
            }
          }}
        >
          {#each e2.deliveryMethodsList() as method}
            <option value={method.id}>{method.name}</option>
          {/each}
        </select>
      </label>

      <label>
        <span>Amount</span>
        <InputQuantity
          unit={draft.deliveryMethod.units}
          initial={draft.dose}
          onchange={(qty) => {
            draft.dose = qty;
          }}
        />
      </label>
    </div>

    <label>
      <span>Taken at</span>
      <input
        type="datetime-local"
        class="mb-2"
        value={formatDatetimeLocal(draft.takenAt)}
        onchange={(e) => {
          const newDate = DateTime.fromISO(e.currentTarget.value, { zone: draft.takenAt.zone });
          if (newDate.isValid) {
            draft.takenAt = newDate as DateTime<true>;
          }
        }}
      />
      <div class="time-actions flex flex-wrap gap-2 justify-end text-sm">
        <button
          type="button"
          class="outline secondary"
          onclick={() => {
            draft.takenAt = DateTime.now().set({ second: 0, millisecond: 0 }) as DateTime<true>;
          }}
        >
          Now
        </button>
        <button type="button" class="outline" onclick={() => addToDoseTime({ minutes: -15 })}>
          -15m
        </button>
        <button type="button" class="outline" onclick={() => addToDoseTime({ hour: -1 })}>
          -1h
        </button>
        <button type="button" class="outline" onclick={() => addToDoseTime({ hour: -6 })}>
          -6h
        </button>
        <button type="button" class="outline" onclick={() => addToDoseTime({ day: -1 })}>
          -1d
        </button>
      </div>
    </label>

    {#if showTakenOffAt}
      <label>
        <span>Taken off at <small class="muted">(optional)</small></span>
        <input
          type="datetime-local"
          class="mb-2"
          value={draft.takenOffAt ? formatDatetimeLocal(draft.takenOffAt) : ""}
          onchange={(e) => {
            if (!e.currentTarget.value) {
              draft.takenOffAt = undefined;
              return;
            }
            const newDate = DateTime.fromISO(e.currentTarget.value, {
              zone: draft.takenAt.zone,
            });
            if (newDate.isValid) {
              draft.takenOffAt = newDate as DateTime<true>;
            }
          }}
        />
        <div class="time-actions flex flex-wrap gap-2 justify-end text-sm">
          <button
            type="button"
            class="outline secondary"
            onclick={() => {
              draft.takenOffAt = DateTime.now().set({
                second: 0,
                millisecond: 0,
              }) as DateTime<true>;
            }}
          >
            Take off now
          </button>
          <button
            type="button"
            class="outline"
            disabled={!draft.takenOffAt}
            onclick={() => {
              draft.takenOffAt = undefined;
            }}
          >
            Clear
          </button>
        </div>
      </label>
    {/if}

    <label>
      <span>Comment <small class="muted">(optional)</small></span>
      <textarea rows="2" placeholder="Optional note..." bind:value={draft.comment}></textarea>
    </label>

    {#if validationError}
      <p class="validation-error text-sm">{validationError}</p>
    {/if}

    <div class="preview-box">
      <LevelsPreview
        mode="impact"
        dosage={doses.dosage}
        {beforeHistory}
        {afterHistory}
        focusTime={draft.takenAt}
      />
    </div>
  </form>

  <footer class="mt-4 flex flex-wrap items-center justify-end gap-2">
    <button
      type="button"
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

    {#if mode === "edit"}
      <button
        type="button"
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
            open = false;
          } catch (err) {
            logErrorToast("Failed to forget dose", err);
          } finally {
            confirmDelete = false;
            busy = false;
          }
        }}
      >
        {#if confirmDelete}
          Really? <Icon name="check" />
        {:else}
          Forget <Icon name="delete" />
        {/if}
      </button>
    {/if}

    <button
      type="button"
      aria-label={mode === "create" ? "Log dose" : "Save dose"}
      class="secondary"
      disabled={busy || !!validationError}
      onclick={async () => {
        if (validationError || busy) return;
        try {
          busy = true;
          await saveDose();
          open = false;
        } catch (err) {
          logErrorToast("Failed to save dose", err);
        } finally {
          busy = false;
        }
      }}
    >
      {mode === "create" ? "Log dose" : "Save"}
      <Icon name="save" />
    </button>
  </footer>
</Dialog>

<style lang="scss">
  @use "sass:map";
  @use "@picocss/pico/scss/settings" as *;

  form label span {
    display: block;
    margin-bottom: calc(var(--pico-spacing) / 4);
  }

  .editor-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 0 var(--pico-spacing);

    @media (max-width: map.get(map.get($breakpoints, "sm"), "breakpoint")) {
      grid-template-columns: 1fr;
    }
  }

  .time-actions button {
    padding: 0.25rem 0.6rem;
    font-size: 0.85em;
  }

  .validation-error {
    color: var(--pico-del-color);
    margin-bottom: calc(var(--pico-spacing) / 2);
  }

  .preview-box {
    margin-top: calc(var(--pico-spacing) / 2);
    padding: calc(var(--pico-spacing) * 0.75);
    border: var(--pico-border-width) solid var(--pico-muted-border-color);
    border-radius: var(--pico-border-radius);
  }

  .muted {
    color: var(--pico-muted-color);
    font-weight: normal;
  }
</style>
