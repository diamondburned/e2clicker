<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import Dialog from "$lib/components/Dialog.svelte";
  import QRCode from "$lib/components/QRCode.svelte";
  import ErrorBox from "$lib/components/ErrorBox.svelte";
  import LoadingPage from "$lib/components/LoadingPage.svelte";

  import PreferenceGroup from "./PreferenceGroup.svelte";
  import Preference from "./Preference.svelte";
  import InputDays from "./InputDays.svelte";
  import InputQuantity from "./InputQuantity.svelte";

  import * as api from "$lib/api.svelte";
  import * as e2 from "$lib/e2";
  import { slide } from "svelte/transition";
  import { sineInOut } from "svelte/easing";

  let user = $state<api.User | null>(null);
  let secret = $state<api.UserSecret | null>(null);
  let dosage = $state<Partial<api.Dosage> | null>(null);

  let promise = $state(
    Promise.all([
      (async () => (user = await api.currentUser()))(),
      (async () => (secret = (await api.currentUserSecret()).secret))(),
      (async () => (dosage = (await api.dosage()).dosage ?? {}))(),
    ]),
  );

  let showSecretDialog = $state(false);
  let deliveryMethod = $derived(e2.deliveryMethods.find((m) => m.id === dosage?.deliveryMethod));

  let dosageIsValid = $derived(
    // [dosage, dosage?.deliveryMethod, dosage?.dose, dosage?.interval].every((v) => !!v),
    !!dosage && !!dosage.deliveryMethod && !!dosage.dose && !!dosage.interval,
  );
  $inspect({ dosage, dosageIsValid });

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

<LoadingPage {promise} />

{#if !deliveryMethod}
  <article
    class="no-delivery-method error-box"
    transition:slide={{ duration: 250, axis: "y", easing: sineInOut }}
  >
    <h3>A delivery method is required.</h3>
    <p class="m-0 text-sm">
      In order for <span class="brand">e2clicker</span> to provide any predictive functionality, you
      need to specify how you are taking your estradiol. Please select a delivery method from the dropdown
      below.
    </p>
  </article>
{/if}

<PreferenceGroup name="Dosage">
  {#snippet misc()}
    {@render preferencesLoader(setDosage)}
  {/snippet}

  {#snippet description()}
    Fine tune your dosage schedule.
  {/snippet}

  {#if dosage}
    <Preference>
      {#snippet name()}
        <span class:brand={!deliveryMethod}>Delivery Method</span>
      {/snippet}
      {#snippet description()}
        How you want to take your medication. This is
        <b class="brand">required</b> for any functionality to work. For details, refer to
        <a href="https://estrannai.se">estrannai.se</a>.
      {/snippet}
      <select
        name="delivery-method"
        class="text-ellipsis"
        required
        bind:value={dosage.deliveryMethod}
        onchange={() => saveDosage()}
      >
        <option value={undefined}>- Unspecified -</option>
        {#each e2.deliveryMethods as method}
          <option value={method.id}>{method.name}</option>
        {/each}
      </select>
    </Preference>
  {/if}

  {#if dosage && deliveryMethod}
    <div transition:slide={{ duration: 200, axis: "y" }}>
      <Preference name="Dose">
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
      </Preference>

      <Preference name="Interval">
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
      </Preference>

      {#if deliveryMethod.patch}
        <div transition:slide={{ duration: 200, axis: "y" }}>
          <Preference name="Patch Change">
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
          </Preference>
        </div>
      {/if}
    </div>
  {/if}
</PreferenceGroup>

<PreferenceGroup name="Notification">
  {#snippet description()}
    Configure how you want to be notified.
  {/snippet}
</PreferenceGroup>

<PreferenceGroup name="Account">
  {#snippet description()}
    Change your account settings.
  {/snippet}

  <Preference name="Name">
    {#snippet description()}
      Your name as it appears on the site.
    {/snippet}
    {#if user}
      <input type="text" autocomplete="off" bind:value={user.name} disabled />
    {/if}
  </Preference>

  <Preference name="Secret">
    {#snippet description()}
      Your account secret. This is used to identify your account apart from others. It can not be
      changed.
      <b>It is strongly recommended that you store this in a safe place.</b>
    {/snippet}
    <button onclick={() => (showSecretDialog = true)}>
      Reveal Secret <Icon name="visibility" />
    </button>
  </Preference>
</PreferenceGroup>

<Dialog bind:open={showSecretDialog} dismissible --max-width="400px">
  <header>
    <h3 class="text-center">Your Account Secret</h3>
  </header>

  {#if secret}
    <section class="flex flex-col items-center spaced m-0">
      <QRCode class="max-w-56" content={api.secretQRData(secret)} padding={2} />
      <p class="secret font-mono">{secret}</p>
    </section>
  {:else}
    <span aria-busy="true"></span>
  {/if}
</Dialog>

{#snippet preferencesLoader(asyncToOK: api.AnyAsyncToOK)}
  {#await asyncToOK.promise}
    <span aria-busy="true"></span>
  {:then}
    <span class="text-green">Saved <Icon name="done" /></span>
  {:catch error}
    <ErrorBox tiny {error} />
  {/await}
{/snippet}
