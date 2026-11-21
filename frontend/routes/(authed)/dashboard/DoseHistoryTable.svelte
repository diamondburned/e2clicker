<script lang="ts">
  import ResizeContainer from "#lib/components/ResizeContainer.svelte";
  import Tooltip from "#lib/components/popovers/Tooltip.svelte";
  import Icon from "#lib/components/Icon.svelte";

  import * as e2 from "#lib/e2.svelte.ts";
  import * as api from "#lib/api.ts";
  import { DateTime } from "luxon";
  import { logErrorToast } from "#lib/toasts.ts";

  let {
    now,
    doses,
    update,
    oneditdose,
    onlogdose,
  }: {
    now: DateTime;
    doses: {
      dosage?: api.Dosage;
      history?: e2.DosageHistory;
    };
    update: () => void;
    oneditdose: (dose: e2.Dose) => void;
    onlogdose: () => void;
  } = $props();

  let dosage = $derived(doses?.dosage);
  let history = $derived(doses?.history ?? []);

  let searchQuery = $state("");
  let pageSize = $state(10);
  let selectedIds = $state<string[]>([]);
  let confirmBulkDelete = $state(false);
  let deletingBulk = $state(false);

  let reversedHistory = $derived(history.toReversed());

  let filteredDoses = $derived.by(() => {
    const q = searchQuery.trim().toLowerCase();
    if (!q) return reversedHistory;
    return reversedHistory.filter((dose) => {
      const dateStr = dose.takenAt.toLocaleString(DateTime.DATETIME_MED).toLowerCase();
      return (
        dose.deliveryMethod.name.toLowerCase().includes(q) ||
        dose.deliveryMethod.id.toLowerCase().includes(q) ||
        String(dose.dose).includes(q) ||
        (dose.comment?.toLowerCase().includes(q) ?? false) ||
        dateStr.includes(q)
      );
    });
  });

  let visibleDoses = $derived(filteredDoses.slice(0, pageSize));
  let hasMore = $derived(filteredDoses.length > visibleDoses.length);

  // Drop stale selected IDs when history updates.
  $effect(() => {
    const valid = new Set(history.map((d) => d.id));
    if (selectedIds.some((id) => !valid.has(id))) {
      selectedIds = selectedIds.filter((id) => valid.has(id));
      confirmBulkDelete = false;
    }
  });

  function toggleSelected(id: string) {
    confirmBulkDelete = false;
    if (selectedIds.includes(id)) {
      selectedIds = selectedIds.filter((x) => x !== id);
    } else {
      selectedIds = [...selectedIds, id];
    }
  }

  async function deleteSelected() {
    if (selectedIds.length === 0 || deletingBulk) return;
    if (!confirmBulkDelete) {
      confirmBulkDelete = true;
      return;
    }

    try {
      deletingBulk = true;
      await api.forgetDoses(selectedIds);
      selectedIds = [];
      confirmBulkDelete = false;
      update();
    } catch (err) {
      logErrorToast("Failed to delete selected doses", err);
    } finally {
      deletingBulk = false;
    }
  }
</script>

<div class="history-toolbar flex flex-wrap items-center justify-between gap-2 mb-3">
  <input
    type="search"
    class="search-input m-0"
    placeholder="Filter by method, dose, note, or date..."
    bind:value={searchQuery}
    oninput={() => {
      pageSize = 10;
    }}
  />

  <div class="toolbar-actions flex flex-wrap items-center gap-2">
    {#if selectedIds.length > 0}
      <button
        type="button"
        class="outline"
        disabled={deletingBulk}
        onclick={() => {
          selectedIds = [];
          confirmBulkDelete = false;
        }}
      >
        Clear ({selectedIds.length})
      </button>
      <button type="button" disabled={deletingBulk} onclick={deleteSelected}>
        <Icon name={confirmBulkDelete ? "check" : "delete"} />
        {confirmBulkDelete
          ? `Really delete ${selectedIds.length}?`
          : `Delete (${selectedIds.length})`}
      </button>
    {/if}

    <button type="button" class="secondary outline" onclick={onlogdose}>
      <Icon name="add" /> Log dose
    </button>
  </div>
</div>

{#if visibleDoses.length === 0}
  <p class="text-[var(--pico-muted-color)] text-center my-4">
    {history.length === 0 ? "No doses recorded in this range yet." : "No matching doses found."}
  </p>
{:else}
  <ResizeContainer>
    <table id="dose-history-table">
      <tbody>
        <tr>
          <th data-column="Select"></th>
          <th data-column="When">When</th>
          <th data-column="Dose">Dose</th>
          <th data-column="Misc"></th>
        </tr>
        {#each visibleDoses as dose (dose.id)}
          {@const selected = selectedIds.includes(dose.id)}
          <tr class="dose-row" class:selected onclick={() => oneditdose(dose)}>
            <td
              data-column="Select"
              onclick={(ev) => {
                ev.stopPropagation();
              }}
            >
              <input
                type="checkbox"
                aria-label="Select dose from {dose.takenAt.toLocaleString(DateTime.DATETIME_MED)}"
                checked={selected}
                onchange={() => toggleSelected(dose.id)}
              />
            </td>
            <td data-column="When">
              <div>{e2.formatDoseTime(dose, now)}</div>
              <small class="exact-time">
                {dose.takenAt.toLocaleString(DateTime.DATETIME_MED)}
                {#if dose.takenOffAt}
                  · off {dose.takenOffAt.toLocaleString(DateTime.DATETIME_MED)}
                {/if}
              </small>
            </td>
            <td data-column="Dose">
              {dose.dose}
              {dose.deliveryMethod.units}
              {#if dose.deliveryMethod.id != dosage?.deliveryMethod}
                <small class="delivery">({dose.deliveryMethod.name})</small>
              {/if}
            </td>
            <td data-column="Misc">
              <div class="misc-actions">
                {#if dose.comment}
                  <span
                    role="presentation"
                    onclick={(ev) => {
                      ev.stopPropagation();
                    }}
                  >
                    <Tooltip>
                      <Icon name="comment" />
                      {#snippet tooltip()}
                        <span class="comment">{dose.comment}</span>
                      {/snippet}
                    </Tooltip>
                  </span>
                {/if}
                <button
                  type="button"
                  class="minimal inline"
                  aria-label="Edit dose"
                  onclick={(ev) => {
                    ev.stopPropagation();
                    oneditdose(dose);
                  }}
                >
                  <Icon name="edit" />
                </button>
              </div>
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </ResizeContainer>
{/if}

{#if hasMore || (pageSize > 10 && filteredDoses.length > 10)}
  <div class="paginator flex flex-wrap items-center justify-center gap-2 mt-2">
    {#if hasMore}
      <button
        type="button"
        class="outline"
        onclick={() => {
          pageSize += 10;
        }}
      >
        Show older ({filteredDoses.length - visibleDoses.length} more)
      </button>
    {/if}
    {#if pageSize > 10}
      <button
        type="button"
        class="outline secondary"
        onclick={() => {
          pageSize = 10;
        }}
      >
        Collapse
      </button>
    {/if}
  </div>
{/if}
<div id="dose-history-anchor"></div>

<style lang="scss">
  @use "sass:map";
  @use "@picocss/pico/scss/settings" as *;

  .history-toolbar {
    input,
    button {
      margin-bottom: 0;
    }
  }

  .search-input {
    flex: 1 1 14rem;
    max-width: 24rem;
  }

  .exact-time {
    display: block;
    color: var(--pico-muted-color);
  }

  .misc-actions {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
  }

  #dose-history-table {
    tbody {
      display: grid;
      grid-template-columns: auto 1fr 1fr auto;

      tr {
        display: contents;
      }

      th,
      td {
        background: none;
      }

      tr.dose-row {
        cursor: pointer;

        &:hover td {
          background-color: color-mix(in srgb, var(--pico-primary) 6%, transparent);
        }

        &.selected td {
          background-color: color-mix(in srgb, var(--pico-primary) 10%, transparent);
        }
      }

      th {
        font-weight: bold;
        border-width: calc(2 * var(--pico-border-width));
      }

      td[data-column="Select"] input {
        margin: 0;
      }

      @media (max-width: map.get(map.get($breakpoints, "md"), "breakpoint")) {
        grid-template-columns: 1fr;

        tr {
          display: grid;
          grid-template-columns: auto 1fr auto;
          grid-template-rows: auto auto;
        }

        th {
          display: none;
        }

        td[data-column="Select"] {
          grid-column: 1;
          grid-row: span 2;
          display: flex;
          align-items: center;
          padding-right: 0.25rem;
        }

        td[data-column="When"] {
          grid-column: 2;
          grid-row: 1;
          border-bottom: none;
          padding-bottom: 0;
        }

        td[data-column="Dose"] {
          grid-column: 2;
          grid-row: 2;
          padding-top: 0;
          font-weight: bold;

          .delivery {
            font-weight: normal;
          }
        }

        td[data-column="Misc"] {
          grid-column: 3;
          grid-row: span 2;
          display: flex;
          align-items: center;
        }
      }
    }
  }
</style>
