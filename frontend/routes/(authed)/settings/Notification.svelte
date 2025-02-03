<script lang="ts" module>
  const neverPromise = new Promise<any>(() => {});
</script>

<script lang="ts">
  import ErrorBox from "$lib/components/ErrorBox.svelte";
  import Icon from "$lib/components/Icon.svelte";
  import PreferenceGroup from "$lib/components/preference/PreferenceGroup.svelte";
  import PreferenceItem from "$lib/components/preference/PreferenceItem.svelte";

  import * as notification from "$lib/notification";
  import { onMount } from "svelte";
  import { fade } from "svelte/transition";

  let pushReason = $state<notification.NotificationError | undefined>();
  let promise = $state<Promise<notification.UpdatePushSubscriptionResult>>(neverPromise);
  onMount(() => {
    promise = notification.updatePushSubscription();
  });
</script>

<PreferenceGroup name="Notification">
  {#snippet description()}
    Configure how you want to be notified.
  {/snippet}

  <!--
  <PreferenceItem name="Email Address">
    {#snippet description()}
      Receive notifications via email, if enabled.
    {/snippet}
  </PreferenceItem>
  -->

  <PreferenceItem name="Push Notifications">
    {#snippet description()}
      Receive notifications via your browser. This will only work for this particular device.
      <!-- <b>This method may be unreliable! Prefer other methods if possible.</b> -->
      <b>Warning! Notifications are not implemented yet!</b>
      You will NOT receive any notifications yet.
    {/snippet}

    {#await promise}
      <span aria-busy="true"></span>
    {:then { enabled, available, reason }}
      {#if available}
        <button
          class:outline={enabled}
          onclick={() => {
            promise = notification.updatePushSubscription({ toggle: !enabled });
          }}
          in:fade={{ duration: 200 }}
        >
          {#if enabled}
            <Icon name="notifications_off" />
            Disable
          {:else}
            <Icon name="notifications" />
            Activate
          {/if}
        </button>
      {:else if reason}
        <ErrorBox error={reason} prefix="Not available" tiny />
      {:else}
        <span class="error-text"> Not available </span>
      {/if}
    {/await}
  </PreferenceItem>
</PreferenceGroup>
