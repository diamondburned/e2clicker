<script lang="ts">
  import PreferenceItem from "$lib/components/preference/PreferenceItem.svelte";
  import PreferenceGroup from "$lib/components/preference/PreferenceGroup.svelte";
  import InputDays from "$lib/components/inputs/InputDays.svelte";
  import InputQuantity from "$lib/components/inputs/InputQuantity.svelte";

  import * as api from "$lib/api.svelte";
  import * as e2 from "$lib/e2.svelte";
  import { slide } from "svelte/transition";
  import { sineInOut } from "svelte/easing";
  import { onMount } from "svelte";

  let dosage = $state<Partial<api.Dosage> | null>(null);
  let initialPromise = $state<Promise<void>>();

  onMount(() => {
    initialPromise = (async () => {
      const response = await api.dosage();
      dosage = response.dosage ?? {};
    })();
  });

  let deliveryMethod = $derived(
    dosage?.deliveryMethod //
      ? e2.deliveryMethod(dosage.deliveryMethod)
      : null,
  );

  let dosageIsValid = $derived(
    !!dosage && !!dosage.deliveryMethod && !!dosage.dose && !!dosage.interval,
  );

  let setDosage = new api.AsyncToOK(api.discard(api.setDosage), {
    debounce: true,
    initial: null,
  });
  function saveDosage() {
    if (dosage && dosageIsValid) {
      setDosage.do(dosage as api.Dosage);
    }
  }
</script>

<PreferenceGroup
  name="Dosage"
  loader={{
    loader: Promise.all([initialPromise, setDosage.promise]),
    valid: dosageIsValid,
  }}
>
  {#snippet description()}
    Fine tune your dosage schedule.
  {/snippet}

  {#snippet header()}
    <blockquote class="popping">
      <h4>Please only use this if you already have a regimen.</h4>
      <span class="text-sm">
        If you are unsure, please consult with your doctor or refer to various resources online such
        as
        <a href="https://transfemscience.org/">Transfeminine Science</a> or
        <a href="https://estrannai.se">estrannaise</a>.
        <b>
          This app is not a substitute for medical advice, professional or not. We are not
          responsible for any harm.
        </b>
        Please take care of yourself! 💖
      </span>
    </blockquote>

    {#if !deliveryMethod}
      <blockquote
        class="no-delivery-method popping error-box"
        transition:slide={{ duration: 250, axis: "y", easing: sineInOut }}
      >
        <h4>A delivery method is required.</h4>
        <span class="text-sm">
          In order for <span class="brand">e2clicker</span> to provide any predictive functionality,
          you need to specify how you are taking your estradiol. Please select a delivery method from
          the dropdown below.
        </span>
      </blockquote>
    {/if}
  {/snippet}

  {#if dosage}
    <PreferenceItem>
      {#snippet name()}
        <span class:brand={!deliveryMethod}>Delivery Method</span>
      {/snippet}
      {#snippet description()}
        How you want to take your medication. This is
        <b class="brand">required</b> for any functionality to work.
      {/snippet}
      <select
        name="delivery-method"
        class="text-ellipsis"
        required
        bind:value={dosage.deliveryMethod}
        onchange={() => saveDosage()}
      >
        <option value={undefined}>- Unspecified -</option>
        {#each e2.deliveryMethodsList() as method}
          <option value={method.id}>{method.name}</option>
        {/each}
      </select>
    </PreferenceItem>
  {/if}

  {#if dosage && deliveryMethod}
    <div transition:slide={{ duration: 200, axis: "y" }}>
      <PreferenceItem name="Dose">
        {#snippet description()}
          How much medication you are taking for each dose.
        {/snippet}
        <InputQuantity
          unit={deliveryMethod.units}
          initial={dosage.dose}
          onchange={(qty) => {
            dosage!.dose = qty;
            saveDosage();
          }}
        />
      </PreferenceItem>

      <PreferenceItem name="Interval">
        {#snippet description()}
          The time between each dose in hours, days or weeks.
        {/snippet}
        <InputDays
          initial={dosage.interval}
          onchange={(days) => {
            dosage!.interval = days;
            saveDosage();
          }}
          placeholder="1 week"
        />
      </PreferenceItem>

      {#if deliveryMethod.patch}
        <div transition:slide={{ duration: 200, axis: "y" }}>
          <PreferenceItem name="Patch Change">
            {#snippet description()}
              How many patches you have on at once before changing the oldest one. If set to
              <em>0</em>, the system will not automatically mark the last patch as "taken off".
            {/snippet}
            <input
              type="number"
              min="0"
              max="7"
              bind:value={dosage.concurrence}
              onchange={() => saveDosage()}
            />
          </PreferenceItem>
        </div>
      {/if}
    </div>
  {/if}
</PreferenceGroup>
