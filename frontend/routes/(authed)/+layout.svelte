<script lang="ts">
  import Header from "$lib/components/Header.svelte";
  import Footer from "$lib/components/Footer.svelte";

  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { get } from "svelte/store";
  import * as api from "$lib/api.svelte";
  import { isLoggedIn } from "$lib/api.svelte";
  import { updateDeliveryMethods } from "$lib/e2/methods.svelte";
  import { updatePushSubscription } from "$lib/notification";
  import LoadingPage from "$lib/components/LoadingPage.svelte";

  let { children } = $props();

  onMount(() => {
    if (!$isLoggedIn) {
      goto("/");
    }
  });

  async function init() {
    await Promise.all([updateDeliveryMethods(), api.updateUser()]);
    updatePushSubscription({
      accountSecret: get(api.user)?.secret,
      token: api.currentToken() ?? undefined,
    }).catch((err) => {
      console.warn("Failed to sync web push subscription on mount:", err);
    });
  }

  let promise = $state(new Promise(() => {}));
  onMount(() => {
    promise = init();
  });
</script>

<Header />

<main class="flex-1 container spaced-2">
  <LoadingPage {promise}>
    {@render children()}
  </LoadingPage>
</main>

<Footer />
