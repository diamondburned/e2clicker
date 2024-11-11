<script lang="ts" module>
  import { DateTime, Duration } from "luxon";

  // The rate at which to poll the server for new data.
  const apiUpdateRate = Duration.fromObject({ minutes: 5 });

  // The rate at which to update the clock used for relative times.
  const durationUpdateRate = Duration.fromObject({ second: 1 });

  function mountTimer(duration: Duration, callback: () => void): () => void {
    const v = setInterval(callback, duration.as("milliseconds"));
    return () => clearInterval(v);
  }
</script>

<script lang="ts">
  import LoadingPage from "$lib/components/LoadingPage.svelte";
  import Dialog from "$lib/components/Dialog.svelte";
  import Icon from "$lib/components/Icon.svelte";
  import DosagePlot from "./DosagePlot.svelte";

  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api.svelte";
  import { onMount } from "svelte";
  import ResizeContainer from "$lib/components/ResizeContainer.svelte";

  // Update this every second.
  // We'll use this just to render relative times.
  let now = $state(DateTime.now());
  onMount(() =>
    mountTimer(durationUpdateRate, () => {
      now = DateTime.now();
    }),
  );

  let endTime = $state(DateTime.now());
  let startTime = $derived(endTime.minus({ month: 1 }));
  onMount(() =>
    mountTimer(apiUpdateRate, () => {
      endTime = DateTime.now();
    }),
  );

  // Update all the information on this page
  // by updating all the time inputs.
  const update = () => {
    now = DateTime.now();
    endTime = DateTime.now();
  };

  let dosageLoader = new api.AsyncToOK(api.dosage, {
    // Allow future update calls without triggering LoadingPage.
    firstPromiseOnly: true,
  });

  $effect(() => {
    dosageLoader.load({
      historyStart: startTime.toISO(),
      historyEnd: endTime.toISO(),
    });
  });

  let doses = $derived(dosageLoader.value);
  let dosage = $derived(doses?.dosage);
  let delivery = $derived(dosage && e2.deliveryMethod(dosage.deliveryMethod));

  let lastDoseAt = $derived(
    (doses?.history?.length ?? 0) > 0
      ? DateTime.fromISO(doses!.history![doses!.history!.length - 1].takenAt)
      : null,
  );
  let nextDoseAt = $derived(dosage && lastDoseAt && e2.timeUntilNextDose(dosage, lastDoseAt));
  let untilNextDose = $derived(nextDoseAt ? nextDoseAt.diff(now) : Duration.fromMillis(0));
  let isTimeForNextDose = $derived(untilNextDose.toMillis() <= 0);

  let visiblePastDay = $state(0);
  let visibleTotalDays = $derived.by(() => {
    if (!doses?.history) {
      return 0;
    }
    const daysList = doses.history.map((d) => DateTime.fromISO(d.takenAt).startOf("day"));
    const days = new Set(daysList);
    return days.size;
  });

  let visibleDoseTime = $derived(
    (lastDoseAt ?? DateTime.now()).minus({ days: visiblePastDay }).startOf("day"),
  );
  let visibleDoses = $derived.by(() => {
    const start = visibleDoseTime;
    const end = visibleDoseTime.endOf("day");
    return doses?.history?.filter((d) => {
      const takenAt = DateTime.fromISO(d.takenAt);
      return takenAt >= start && takenAt <= end;
    });
  });

  let submittingDose = $state(false);
  async function submitDose() {
    submittingDose = true;
    try {
      await api.recordDose({ takenAt: DateTime.now().toISO() });
      update();
    } finally {
      submittingDose = false;
    }
  }
</script>

<svelte:head>
  <title>Dashboard - e2clicker</title>
</svelte:head>

<LoadingPage promise={dosageLoader.promise} />

<section class="dashboard-grid">
  <div id="next-dose" class="flex flex-col items-center spaced-2">
    {#if nextDoseAt}
      {#if isTimeForNextDose}
        <p class="dose-reminder primary text-center text-3xl font-bold mb-4">
          It's time to take your next dose!
        </p>
        <p class="duration-display text-2xl text-center leading-10">
          <span>Your next dose was due</span>
          <mark class="duration font-bold text-nowrap">
            {e2.formatDuration(untilNextDose.negate())}
          </mark>
          <span>ago</span>
        </p>
      {:else}
        <p class="duration-display text-2xl text-center leading-10">
          <span>Your next dose is in</span>
          <mark class="duration font-bold text-nowrap">
            {e2.formatDuration(untilNextDose)}
          </mark>
        </p>
      {/if}
    {:else}
      <p class="text-xl">You haven't taken a dose yet.</p>
    {/if}

    <footer class="actions min-h-12">
      {#if dosageLoader.loading}
        <span aria-busy="true">Loading...</span>
      {:else}
        <button
          class:outline={!isTimeForNextDose}
          onclick={() => submitDose()}
          disabled={dosageLoader.loading || submittingDose}
        >
          <Icon name="medication" />
          {#if lastDoseAt}
            {#if isTimeForNextDose}
              I took my dose!
            {:else}
              I took it early
            {/if}
          {:else}
            I took my <b>first</b> dose!
          {/if}
        </button>

        {#if lastDoseAt}
          <button class="secondary ml-2" onclick={() => {}}>
            <Icon name="edit" />
            Edit doses
          </button>
        {/if}
      {/if}
    </footer>
  </div>

  <article id="estrannaise-plot">
    <h3>Estrogen Levels</h3>
    <DosagePlot {doses} {startTime} {endTime} />
  </article>

  <article id="dose-info">
    <h3>Your Dosage</h3>
    {#if doses?.dosage && delivery}
      <p class="dosage text-3xl leading-10">
        <span class="text-nowrap">
          {#if delivery.patch}
            1 patch
          {:else}
            {doses.dosage.dose ?? ""}
            {delivery.units ?? ""}
          {/if}
        </span>
        <span class="font-thin">every</span>
        <span class="text-nowrap">
          {e2
            .roundDuration(Duration.fromObject({ days: doses.dosage.interval }).rescale())
            .toHuman()}
        </span>
      </p>
      <p class="medication">
        {delivery.name}
      </p>
    {/if}
  </article>

  <article id="levels-info">
    <h3>Current Levels</h3>
  </article>
</section>

<section id="dose-history" class="as-card">
  <h2 class="no-fat-padding">Dose History</h2>
  <ResizeContainer>
    <table id="dose-history-table">
      <tbody>
        <tr>
          <th data-column="When">When</th>
          <th data-column="Dose">Dose</th>
          <th data-column="Comment"></th>
        </tr>
        {#each (visibleDoses ?? []).toReversed() as dose}
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
            <td data-column="Comment">
              {#if true}
                <Icon name="comment" />
              {/if}
            </td>
          </tr>
        {/each}
      </tbody>
    </table>
  </ResizeContainer>
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
</section>

{#if doses && !doses.dosage}
  <Dialog open>
    <h3>Further setup required</h3>
    <p>
      You don't currently have a dose schedule set up yet.
      <br />
      <span class="brand">Let's get that set up now!</span>
    </p>
    <footer>
      <a href="/settings" role="button">
        Head to settings <Icon name="arrow-forward" />
      </a>
    </footer>
  </Dialog>
{/if}

<style lang="scss">
  @use "sass:map";
  @use "@picocss/pico/scss/settings" as *;

  article {
    margin-bottom: 0;

    h3:not(.no-fat-padding) {
      margin-bottom: calc(1.5 * var(--pico-typography-spacing-vertical));
    }
  }

  section.as-card {
    padding: 0 var(--pico-block-spacing-horizontal);
    @media (max-width: map.get(map.get($breakpoints, "sm"), "breakpoint")) {
      padding: 0;
    }
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    grid-gap: var(--pico-spacing);

    @mixin name-grid($id) {
      ##{$id} {
        grid-area: $id;
      }
    }

    @include name-grid(next-dose);
    @include name-grid(estrannaise-plot);
    @include name-grid(dose-info);
    @include name-grid(levels-info);

    grid-template-areas:
      "next-dose next-dose next-dose"
      "estrannaise-plot estrannaise-plot estrannaise-plot"
      "dose-info levels-info levels-info"
      "dose-history dose-history dose-history";

    @media (max-width: map.get(map.get($breakpoints, "md"), "breakpoint")) {
      grid-template-areas:
        "next-dose"
        "estrannaise-plot"
        "levels-info"
        "dose-info"
        "dose-history";
    }

    & > * {
      grid-area: attr(data-grid);
    }
  }

  #next-dose {
    --y-margin: clamp(
      var(--pico-block-spacing-vertical),
      10vh,
      calc(6 * var(--pico-block-spacing-vertical))
    );

    margin-top: var(--y-margin);
    margin-bottom: var(--y-margin);

    font-size: clamp(1em, 5vw, 1.15em);
  }

  #dose-history-table {
    tbody {
      display: grid;
      grid-template-columns: 1fr 1fr auto;

      tr {
        display: contents;
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

        td[data-column="Comment"] {
          grid-column: 2;
          grid-row: span 2;

          display: flex;
          align-items: center;
        }
      }
    }
  }
</style>
