<script lang="ts" module>
  import { DateTime, Duration, Interval } from "luxon";

  // The rate at which to poll the server for new data.
  const apiUpdateRate = Duration.fromObject({ minutes: 5 });

  // The rate at which to update the clock used for relative times.
  const durationUpdateRate = Duration.fromObject({ second: 1 });

  function mountTimer(duration: Duration, callback: () => void): () => void {
    const v = setInterval(callback, duration.as("milliseconds"));
    return () => clearInterval(v);
  }

  function intervalUntilNow(duration: Duration) {
    const iv = Interval.fromDateTimes(DateTime.now().minus(duration), DateTime.now());
    if (!iv.isValid) throw new Error(iv.invalidReason);
    return iv;
  }
</script>

<script lang="ts">
  import LoadingPage from "$lib/components/LoadingPage.svelte";
  import Icon from "$lib/components/Icon.svelte";

  import NextDoseCountdown from "./NextDoseCountdown.svelte";
  import CurrentLevels from "./CurrentLevels.svelte";
  import DoseHistoryTable from "./DoseHistoryTable.svelte";
  import DosagePlot from "./DosagePlot.svelte";
  import DoseInfo from "./DoseInfo.svelte";

  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api.svelte";
  import { Data as DosageData } from "./DosagePlot.svelte";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";

  // Update this every second.
  // We'll use this just to render relative times.
  let now = $state(DateTime.now());
  onMount(() =>
    mountTimer(durationUpdateRate, () => {
      now = DateTime.now();
    }),
  );

  const historyDuration = Duration.fromObject({ month: 1 });
  let historyInterval = $state(intervalUntilNow(historyDuration));
  onMount(() =>
    mountTimer(apiUpdateRate, () => {
      historyInterval = intervalUntilNow(historyDuration);
    }),
  );

  // Update all the information on this page
  // by updating all the time inputs.
  const update = () => {
    now = DateTime.now();
    historyInterval = intervalUntilNow(historyDuration);
  };

  let dosageLoader = new api.AsyncToOK(api.dosage, {
    // Allow future update calls without triggering LoadingPage.
    firstPromiseOnly: true,
  });

  $effect(() => {
    dosageLoader.load({
      start: historyInterval.start.toISO(),
      end: historyInterval.end.toISO(),
    });
  });

  let doses = $derived.by(() => {
    if (!dosageLoader.value) {
      return {};
    }

    const { dosage, history } = dosageLoader.value;
    return {
      dosage,
      history: history && e2.convertDoseHistory(history),
    };
  });

  let dosageData = $derived(
    new DosageData({
      doses,
      interval: historyInterval,
    }),
  );

  let editingDoses = $state(false);
</script>

<svelte:head>
  <title>Dashboard - e2clicker</title>
</svelte:head>

<LoadingPage promise={dosageLoader.promise}>
  {#if doses?.dosage}
    <section class="dashboard-grid">
      <div id="next-dose">
        <NextDoseCountdown {now} {doses} onsubmit={() => update()}>
          {#snippet footer()}
            <button
              class="secondary outline"
              onclick={() => {
                editingDoses = true;
                goto("#dose-history-anchor");
              }}
              disabled={editingDoses}
            >
              <Icon name="edit" />
              Edit doses
            </button>
          {/snippet}
        </NextDoseCountdown>
      </div>

      <article id="estrannaise-plot">
        <h3>Estrogen Levels</h3>
        <DosagePlot data={dosageData} />
      </article>

      <article id="dose-info">
        <DoseInfo dosage={doses?.dosage} />
      </article>

      <article id="levels-info">
        <h3>Current Levels</h3>
        <CurrentLevels data={dosageData} />
      </article>
    </section>

    <section id="dose-history" class="as-card">
      <h2>Dose History</h2>
      <DoseHistoryTable {now} {doses} {update} bind:editing={editingDoses} />
    </section>
  {:else}
    <section id="setup-required" class="w-fit mx-auto">
      <h2>Further setup required</h2>
      <p>
        You don't currently have a dose schedule set up yet.
        <br />
        <span class="brand">Let's get that set up now!</span>
      </p>
      <footer class="text-center my-6">
        <a href="/settings" role="button">
          Head to settings <Icon name="arrow-forward" />
        </a>
      </footer>
    </section>
  {/if}
</LoadingPage>

<style lang="scss">
  @use "sass:map";
  @use "@picocss/pico/scss/settings" as *;

  section:last-child {
    margin-bottom: var(--pico-block-spacing-vertical);
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

    article {
      margin-bottom: 0;
    }
  }

  #next-dose,
  #setup-required {
    --y-margin: clamp(
      var(--pico-block-spacing-vertical),
      10vh,
      calc(6 * var(--pico-block-spacing-vertical))
    );

    margin-top: var(--y-margin);
    margin-bottom: var(--y-margin);

    font-size: clamp(1em, 5vw, 1.15em);
  }
</style>
