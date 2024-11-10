<script lang="ts">
  import Estrannaise from "./Estrannaise.svelte";

  import * as api from "$lib/api";
  import { fly } from "svelte/transition";
  import { DateTime } from "luxon";
  import LoadingScreen from "$lib/components/LoadingScreen.svelte";
  import Dialog from "$lib/components/Dialog.svelte";
  import Icon from "$lib/components/Icon.svelte";
  import { onMount } from "svelte";

  let historyEnd = $state(DateTime.now());
  let historyStart = $derived(historyEnd.minus({ month: 1 }));

  let loadingPromise = $state<
    Promise<{
      doseHistory: api.DosageHistory["history"];
      doseSchedule: api.DosageSchedule | null;
    }>
  >(new Promise(() => {}));

  onMount(() => {
    loadingPromise = Promise.all([
      api.doseHistory(historyStart.toISO(), historyEnd.toISO()),
      api.dosageSchedule(),
    ]).then(([doseHistory, doseSchedule]) => ({
      doseHistory: doseHistory.history,
      doseSchedule: doseSchedule.schedule ?? null,
    }));
  });
</script>

<svelte:head>
  <title>e2clicker</title>
</svelte:head>

{#await loadingPromise}
  <LoadingScreen promise={loadingPromise} />
{:then { doseHistory, doseSchedule }}
  {#if doseSchedule}
    <main class="container" transition:fly={{ duration: 200, y: 100 }}>
      <h1>Dashboard</h1>

      <section class="dashboard-grid">
        <article id="estrannaise-plot">
          <h2>Estrogen Levels</h2>
          <Estrannaise {doseHistory} />
        </article>

        <article id="dose-history">
          <h2>Dose History</h2>
          <ul>
            {#each doseHistory as dose}
              <li>{dose}</li>
            {/each}
          </ul>
        </article>

        <article id="next-dose">
          <b>Take your next dose in</b>
        </article>

        <article id="actions">
          <button> I took my dose! </button>
        </article>
      </section>
    </main>
  {:else}
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
{:catch}
  <LoadingScreen promise={loadingPromise} />
{/await}

<style lang="scss">
  main {
    flex: 1;

    h1 {
      margin-top: var(--pico-typography-spacing-top);
      margin-bottom: calc(2 * var(--pico-typography-spacing-vertical));
    }
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    grid-gap: var(--pico-spacing);

    @mixin name-grid($id) {
      ##{$id} {
        grid-area: $id;

        margin: 0;
      }
    }

    @include name-grid(estrannaise-plot);
    @include name-grid(dose-history);
    @include name-grid(next-dose);
    @include name-grid(actions);

    grid-template-areas:
      "estrannaise-plot estrannaise-plot estrannaise-plot"
      "next-dose dose-history dose-history"
      "actions actions actions";

    @media (max-width: 800px) {
      grid-template-areas:
        "estrannaise-plot"
        "next-dose"
        "dose-history"
        "actions";
    }

    & > * {
      grid-area: attr(data-grid);
    }
  }
</style>
