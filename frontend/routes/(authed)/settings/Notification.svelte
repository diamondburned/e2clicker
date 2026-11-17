<script lang="ts">
  import ErrorBox from "$lib/components/ErrorBox.svelte";
  import PreferenceGroup from "$lib/components/preference/PreferenceGroup.svelte";
  import PreferenceItem from "$lib/components/preference/PreferenceItem.svelte";

  import * as api from "$lib/api.svelte";
  import { user } from "$lib/api.svelte";
  import Icon from "$lib/components/Icon.svelte";
  import { updatePushSubscription, type UpdatePushSubscriptionResult } from "$lib/notification";
  import { addToast, logErrorToast } from "$lib/toasts";
  import { onMount } from "svelte";

  let preferences = $state<api.NotificationPreferences>();
  let supportedMethods = $state<api.NotificationMethodSupports>([]);

  let preferencesPromise = Promise.all([
    api.userNotificationPreferences(),
    api.supportedNotificationMethods(),
  ]).then(([prefs, methods]) => {
    preferences = prefs;
    supportedMethods = methods;
    return prefs;
  });

  let setPreferences = new api.AsyncToOK(
    async (newPreferences: api.NotificationPreferences) => {
      await api.userUpdateNotificationPreferences(newPreferences);
      preferences = newPreferences;
    },
    {
      debounce: true,
      initial: null,
    },
  );

  let pushState = $state<UpdatePushSubscriptionResult | null>(null);
  let pushLoadError = $state<unknown>(null);
  let togglingPush = $state(false);
  let testingNotification = $state(false);

  onMount(async () => {
    try {
      pushState = await updatePushSubscription({
        accountSecret: $user?.secret,
        token: api.currentToken() ?? undefined,
      });
    } catch (err) {
      pushLoadError = err;
    }
  });

  async function togglePush() {
    if (!pushState || togglingPush) return;
    togglingPush = true;
    pushLoadError = null;
    try {
      const next = !pushState.enabled;
      pushState = await updatePushSubscription({
        toggle: next,
        accountSecret: $user?.secret,
        token: api.currentToken() ?? undefined,
      });
      if (pushState.reason) {
        addToast({
          urgency: "error",
          message: "Web Push unavailable",
          description: pushState.reason,
        });
      } else {
        addToast({
          message: pushState.enabled
            ? "Web Push notifications enabled on this device."
            : "Web Push notifications disabled on this device.",
        });
      }
    } catch (err) {
      pushLoadError = err;
      logErrorToast("Failed to update Web Push subscription", err);
    } finally {
      togglingPush = false;
    }
  }
</script>

<PreferenceGroup
  name="Notification"
  loader={{
    loader: Promise.all([preferencesPromise, setPreferences.promise]),
  }}
>
  {#snippet description()}
    Configure how you want to be notified.
  {/snippet}

  <PreferenceItem name="Email">
    {#snippet description()}
      Receive notifications via email. This is the most reliable method of notification.
    {/snippet}

    {#await preferencesPromise}
      <span aria-busy="true"></span>
    {:then preferences}
      {#if supportedMethods.includes("email")}
        <input
          type="email"
          value={(preferences.notificationConfigs.email ?? [])
            .map((email) => email.address)
            .join(", ")}
          onchange={(e) => {
            const emails = e.currentTarget.value
              .split(",")
              .map((email) => email.trim())
              .filter((email) => email.length > 0)
              .map((email) => ({ address: email }) as api.EmailSubscription);
            preferences.notificationConfigs.email = emails;
            setPreferences.do(preferences);
          }}
        />
      {:else}
        <ErrorBox
          error="Email notifications are not configured on this server."
          prefix="Not available"
          tiny
        />
      {/if}
    {/await}
  </PreferenceItem>

  <PreferenceItem name="Web Push">
    {#snippet description()}
      Receive notifications via your browser. This will only work for this particular device.
    {/snippet}

    {#await preferencesPromise}
      <span aria-busy="true"></span>
    {:then}
      {#if !supportedMethods.includes("webPush")}
        <ErrorBox error="Server does not support this feature." prefix="Not available" tiny />
      {:else if pushLoadError}
        <ErrorBox error={pushLoadError} prefix="Web Push error" tiny />
      {:else if !pushState}
        <span aria-busy="true"></span>
      {:else if !pushState.available}
        <ErrorBox
          error={pushState.reason ?? "Web Push is unavailable."}
          prefix="Not available"
          tiny
        />
      {:else}
        <div class="flex flex-col items-end gap-1">
          <button
            class={pushState.enabled ? "secondary outline" : "outline"}
            aria-busy={togglingPush}
            disabled={togglingPush}
            onclick={togglePush}
          >
            {#if pushState.enabled}
              Disable Web Push
              <Icon name="notifications_off" />
            {:else}
              Enable Web Push
              <Icon name="notifications_active" />
            {/if}
          </button>
          {#if pushState.reason}
            <ErrorBox error={pushState.reason} prefix="Blocked" tiny />
          {/if}
        </div>
      {/if}
    {/await}
  </PreferenceItem>

  <PreferenceItem name="Test Notification">
    {#snippet description()}
      Send a test notification to verify that your settings are correct.
    {/snippet}
    <button
      class="outline ml-2"
      onclick={async () => {
        testingNotification = true;
        try {
          await api.sendTestNotification();
          addToast({ message: "Test notification sent!" });
        } catch (err) {
          logErrorToast("Failed to send test notification", err);
        } finally {
          setTimeout(() => {
            testingNotification = false;
          }, 5000);
        }
      }}
      disabled={testingNotification || supportedMethods.length === 0}
    >
      Test Notification
      <Icon name="send" />
    </button>
  </PreferenceItem>
</PreferenceGroup>
