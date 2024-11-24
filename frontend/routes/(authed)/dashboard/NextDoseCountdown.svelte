<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";

  import * as e2 from "$lib/e2.svelte";
  import * as api from "$lib/api.svelte";
  import { type Snippet } from "svelte";
  import { DateTime, Duration } from "luxon";
  import { slide } from "svelte/transition";

  let {
    now,
    doses,
    onsubmit,
    footer,
  }: {
    now: DateTime;
    doses: {
      dosage?: api.Dosage;
      history?: e2.DosageHistory;
    };
    onsubmit: () => void;
    footer?: Snippet;
  } = $props();

  let dosage = $derived(doses?.dosage);
  let history = $derived(doses?.history ?? []);
  let loading = $derived(dosage === undefined || history === undefined);

  let lastDoseAt = $derived(
    (history?.length ?? 0) > 0 //
      ? history![history!.length - 1].takenAt
      : null,
  );
  let nextDoseAt = $derived(dosage && lastDoseAt && e2.timeUntilNextDose(dosage, lastDoseAt));
  let untilNextDose = $derived(nextDoseAt ? nextDoseAt.diff(now) : Duration.fromMillis(0));
  let isTimeForNextDose = $derived(untilNextDose.toMillis() <= 0);

  let submittingDose = $state(false);
  async function submitDose() {
    submittingDose = true;
    try {
      await api.recordDose();
      onsubmit();
    } finally {
      submittingDose = false;
    }
  }
</script>

{#snippet reminder()}
  <p class="dose-reminder primary text-center text-2xl font-bold mb-4" transition:slide>
    It's time!
  </p>
{/snippet}

{#snippet doseDue()}
  <p class="duration-display">
    <span>Your next dose was due</span>
    <mark class="duration primary">
      {e2.formatDuration(untilNextDose.negate())}
    </mark>
    <span>ago</span>
    {@render doseAux()}
  </p>
{/snippet}

{#snippet doseUntil()}
  <p class="duration-display">
    <span>Your next dose is in</span>
    <mark class="duration secondary">
      {e2.formatDuration(untilNextDose)}
    </mark>
    {@render doseAux()}
  </p>
{/snippet}

{#snippet doseAux()}
  <span class="duration-aux">
    <br />
    on {nextDoseAt?.toLocaleString({
      weekday: "long",
      hour: "numeric",
      minute: "numeric",
    })}
  </span>
{/snippet}

<div class="dose-countdown flex flex-col items-center">
  <div class="duration-container mb-4">
    {#if nextDoseAt}
      {#if isTimeForNextDose}
        {@render reminder()}
        {@render doseDue()}
      {:else}
        {@render doseUntil()}
      {/if}
    {:else}
      <p class="duration-display">You haven't taken a dose yet.</p>
    {/if}
  </div>

  <footer class="actions min-h-12 flex flex-row flex-wrap justify-center gap-2">
    <button
      class:outline={!isTimeForNextDose}
      onclick={() => submitDose()}
      disabled={loading || submittingDose}
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

    {@render footer?.()}
  </footer>
</div>

<style lang="scss">
  .duration-container {
    font-size: var(--font-size-2xl);
    line-height: 2rem;
  }

  .duration-display {
    text-align: center;
    margin: 0;

    & > * {
      white-space: nowrap;
    }
  }

  .duration {
    font-weight: bold;
  }

  .duration-aux {
    font-size: var(--font-size-xs);
  }
</style>
