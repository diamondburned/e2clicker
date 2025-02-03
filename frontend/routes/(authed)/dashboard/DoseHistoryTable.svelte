<script lang="ts">
  import ResizeContainer from "$lib/components/ResizeContainer.svelte";
  import Tooltip from "$lib/components/popovers/Tooltip.svelte";
  import Dialog from "$lib/components/Dialog.svelte";
  import Icon from "$lib/components/Icon.svelte";

  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api";
  import { slide } from "svelte/transition";
  import { DateTime } from "luxon";
  import { logErrorToast } from "$lib/toasts";

  let {
    now,
    doses,
    update = () => {},
    editing = $bindable(false),
  }: {
    now: DateTime;
    doses: {
      dosage?: api.Dosage;
      history?: e2.DosageHistory;
    };
    update?: () => void;
    editing?: boolean;
  } = $props();

  let dosage = $derived(doses?.dosage);
  let history = $derived(doses?.history ?? []);

  let lastDoseAt = $derived(
    history.length > 0 //
      ? history![history!.length - 1].takenAt
      : null,
  );

  let visiblePastWeek = $state(0);

  let visibleTotalWeeks = $derived.by(() => {
    const daysList = history.map((d) => d.takenAt.startOf("week"));
    const days = new Set(daysList);
    return days.size;
  });

  let visibleDoseTime = $derived(
    (lastDoseAt ?? DateTime.now()).minus({ weeks: visiblePastWeek }).startOf("week"),
  );

  let visibleDoses = $derived.by(() => {
    const start = visibleDoseTime;
    const end = visibleDoseTime.endOf("week");
    return history.filter((d) => {
      return d.takenAt >= start && d.takenAt <= end;
    });
  });

  let deleteDoseOpen = $state(false);
  let deletingDose = $state<e2.Dose | null>(null);
  let deletingDoseBusy = $state(false);
  $effect(() => {
    if (!deleteDoseOpen) {
      deletingDose = null;
    }
  });
</script>

{#snippet doseDisplay_when(dose: e2.Dose)}
  {e2.formatDoseTime(dose, now)} ago
{/snippet}

{#snippet doseDisplay_dose(dose: e2.Dose)}
  {dose.dose}
  {dose.deliveryMethod.units}
  {#if dose.deliveryMethod.id != dosage?.deliveryMethod}
    <small class="delivery">({dose.deliveryMethod.name})</small>
  {/if}
{/snippet}

<ResizeContainer>
  <table id="dose-history-table">
    <tbody>
      <tr>
        <th data-column="When">When</th>
        <th data-column="Dose">Dose</th>
        <th data-column="Misc"></th>
      </tr>
      {#each visibleDoses.toReversed() as dose (dose.takenAt)}
        <tr>
          <td data-column="When">{@render doseDisplay_when(dose)}</td>
          <td data-column="Dose">{@render doseDisplay_dose(dose)}</td>
          <td data-column="Misc">
            {#if editing}
              <button
                class="minimal inline"
                onclick={() => {
                  deletingDose = dose;
                  deleteDoseOpen = true;
                }}
              >
                <Icon name="delete" />
              </button>
            {:else if dose.comment}
              <Tooltip>
                <Icon name="comment" />
                {#snippet tooltip()}
                  <span class="comment">{dose.comment}</span>
                {/snippet}
              </Tooltip>
            {/if}
          </td>
        </tr>
      {/each}
    </tbody>
  </table>
</ResizeContainer>

{#if editing && deletingDose}
  <Dialog bind:open={deleteDoseOpen} dismissible class="delete-dose-confirmation">
    <h3>Are you sure you want to delete this dose?</h3>

    <ul class="px-4">
      <li>
        <b>When:</b>
        {@render doseDisplay_when(deletingDose)}
      </li>
      <li>
        <b>Dose:</b>
        {@render doseDisplay_dose(deletingDose)}
      </li>
    </ul>

    <footer>
      <button
        class="secondary outline"
        aria-label="Cancel"
        onclick={() => {
          deleteDoseOpen = false;
        }}
      >
        Cancel
        <Icon name="close" />
      </button>
      <button
        aria-label="Delete"
        disabled={deletingDoseBusy}
        onclick={async () => {
          if (!deletingDose) return;
          try {
            deletingDoseBusy = true;
            await api.forgetDoses([deletingDose._takenAt]);
            update();
          } catch (err) {
            logErrorToast("Failed to delete dose", err);
          } finally {
            deletingDoseBusy = false;
            deleteDoseOpen = false;
          }
        }}
      >
        Delete
        {#if deletingDoseBusy}
          <span aria-busy="true" class="spinner"></span>
        {:else}
          <Icon name="delete" />
        {/if}
      </button>
    </footer>
  </Dialog>
{/if}

{#if editing}
  <p
    class="text-[var(--pico-muted-color)] text-center"
    transition:slide={{
      axis: "y",
      duration: 200,
    }}
  >
    – editing mode –
    <button class="minimal inline" onclick={() => (editing = false)}>stop</button>
    –
  </p>
{/if}

<div class="paginator flex items-center justify-center">
  <button
    class="p-1 outline"
    onclick={() => visiblePastWeek++}
    disabled={visiblePastWeek >= visibleTotalWeeks}
  >
    <Icon name="chevron-left" />
  </button>
  <p class="mx-4 my-0">
    Showing <b>{visibleDoseTime.toLocaleString({ dateStyle: "long" })}</b>
  </p>
  <button class="p-1 outline xx" onclick={() => visiblePastWeek--} disabled={visiblePastWeek <= 0}>
    <Icon name="chevron-right" />
  </button>
</div>
<div id="dose-history-anchor"></div>

<style lang="scss">
  @use "sass:map";
  @use "@picocss/pico/scss/settings" as *;

  #dose-history-table {
    tbody {
      display: grid;
      grid-template-columns: 1fr 1fr auto;

      tr {
        display: contents;
      }

      th,
      td {
        background: none;
      }

      th {
        font-weight: bold;
        border-width: calc(2 * var(--pico-border-width));
      }

      @media (max-width: map.get(map.get($breakpoints, "md"), "breakpoint")) {
        grid-template-columns: 1fr;

        tr {
          display: grid;

          grid-template-columns: 1fr auto;
          grid-template-rows: auto auto;
          /* grid-gap: var(--pico-spacing); */
        }

        th {
          display: none;
        }

        td[data-column="When"] {
          grid-column: 1;
          grid-row: 1;

          border-bottom: none;
          padding-bottom: 0;
        }

        td[data-column="Dose"] {
          grid-column: 1;
          grid-row: 2;

          padding-top: 0;

          font-weight: bold;
          .delivery {
            font-weight: normal;
          }
        }

        td[data-column="Misc"] {
          grid-column: 2;
          grid-row: span 2;

          display: flex;
          align-items: center;
        }
      }
    }
  }
</style>
