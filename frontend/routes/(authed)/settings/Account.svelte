<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import Dialog from "$lib/components/Dialog.svelte";
  import QRCode from "$lib/components/QRCode.svelte";
  import PreferenceItem from "$lib/components/preference/PreferenceItem.svelte";
  import PreferenceGroup from "$lib/components/preference/PreferenceGroup.svelte";

  import * as api from "$lib/api.svelte";
  import { user } from "$lib/api.svelte";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";

  onMount(async () => {
    await api.updateUser();
  });

  let showSecretDialog = $state(false);

  function prettySecret(secret: string) {
    return secret.replace(/(.{4})/g, "$1 ").trim();
  }
</script>

<PreferenceGroup name="Account">
  {#snippet description()}
    Change your account settings.
  {/snippet}

  <PreferenceItem name="Name">
    {#snippet description()}
      Your name as it appears on the site.
    {/snippet}
    <input type="text" autocomplete="off" value={$user?.name ?? ""} disabled />
  </PreferenceItem>

  <PreferenceItem name="Secret">
    {#snippet description()}
      Your account secret. This is used to identify your account apart from others. It can not be
      changed.
      <b>It is strongly recommended that you store this in a safe place.</b>
    {/snippet}
    <button onclick={() => (showSecretDialog = true)} disabled={$user == null}>
      Reveal Secret <Icon name="visibility" />
    </button>
  </PreferenceItem>
</PreferenceGroup>

{#if $user}
  <Dialog bind:open={showSecretDialog} dismissible --max-width="400px">
    <header>
      <h3 class="text-center">Your Account Secret</h3>
    </header>

    <section class="flex flex-col items-center spaced m-0" transition:fade={{ duration: 200 }}>
      <QRCode class="max-w-56" content={api.secretQRData($user.secret)} padding={2} />
      <p class="secret font-mono">{prettySecret($user.secret)}</p>
    </section>
  </Dialog>
{/if}
