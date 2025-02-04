<script lang="ts">
  import ErrorBox from "$lib/components/ErrorBox.svelte";
  import PreferenceGroup from "$lib/components/preference/PreferenceGroup.svelte";
  import PreferenceItem from "$lib/components/preference/PreferenceItem.svelte";

  import * as api from "$lib/api.svelte";
  import Icon from "$lib/components/Icon.svelte";
  import { addToast, logErrorToast } from "$lib/toasts";

  let preferences = $state<api.NotificationPreferences>();
  let preferencesPromise = api.userNotificationPreferences().then((prefs) => {
    preferences = prefs;
    return prefs;
  });

  // let setPreferences = new api.AsyncToOK(api.discard(api.userUpdateNotificationPreferences), {
  //   debounce: true,
  //   initial: null,
  // });

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

  let testingNotification = $state(false);
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
    {/await}
  </PreferenceItem>

  <PreferenceItem name="Web Push">
    {#snippet description()}
      Receive notifications via your browser. This will only work for this particular device.
    {/snippet}

    <ErrorBox error="Not yet implemented!" prefix="Not available" tiny />
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
      disabled={testingNotification}
    >
      Test Notification
      <Icon name="send" />
    </button>
  </PreferenceItem>
</PreferenceGroup>
