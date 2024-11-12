<script lang="ts">
  import ResizeContainer from "$lib/components/ResizeContainer.svelte";
  import Tooltip from "$lib/components/popovers/Tooltip.svelte";
  import Icon from "$lib/components/Icon.svelte";

  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api";
  import { fade, slide } from "svelte/transition";
  import { DateTime } from "luxon";

  let {
    now,
    dosage,
    history,
    editing = $bindable(false),
  }: {
    now: DateTime;
    dosage: api.Dosage;
    history: api.DosageHistory;
    editing?: boolean;
  } = $props();

  let lastDoseAt = $derived(
    history.length > 0 //
      ? DateTime.fromISO(history![history!.length - 1].takenAt)
      : null,
  );

  let visiblePastDay = $state(0);

  let visibleTotalDays = $derived.by(() => {
    const daysList = history.map((d) => DateTime.fromISO(d.takenAt).startOf("day"));
    const days = new Set(daysList);
    return days.size;
  });

  let visibleDoseTime = $derived(
    (lastDoseAt ?? DateTime.now()).minus({ days: visiblePastDay }).startOf("day"),
  );

  let visibleDoses = $derived.by(() => {
    const start = visibleDoseTime;
    const end = visibleDoseTime.endOf("day");
    return history.filter((d) => {
      const takenAt = DateTime.fromISO(d.takenAt);
      return takenAt >= start && takenAt <= end;
    });
  });
</script>

<ResizeContainer>
  <table id="dose-history-table">
    <tbody>
      <tr>
        <th data-column="When">When</th>
        <th data-column="Dose">Dose</th>
        <th data-column="Misc"></th>
      </tr>
      {#each visibleDoses.toReversed() as dose}
        {@const delivery = e2.deliveryMethod(dose.deliveryMethod)}
        <tr>
          <td data-column="When">{e2.formatDoseTime(dose, now)} ago</td>
          <td data-column="Dose">
            {dose.dose}
            {delivery?.units}
            {#if delivery && delivery.id != dosage?.deliveryMethod}
              <small class="delivery">({delivery.name})</small>
            {/if}
          </td>
          <td data-column="Misc">
            {#if dose.comment}
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
    onclick={() => visiblePastDay++}
    disabled={visiblePastDay >= visibleTotalDays}
  >
    <Icon name="chevron-left" />
  </button>
  <p class="mx-4 my-0">
    Showing <b>{visibleDoseTime.toLocaleString({ dateStyle: "long" })}</b>
  </p>
  <button class="p-1 outline xx" onclick={() => visiblePastDay--} disabled={visiblePastDay <= 0}>
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
