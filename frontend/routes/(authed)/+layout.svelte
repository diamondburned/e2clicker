<script lang="ts">
  import Header from "$lib/components/Header.svelte";
  import Footer from "$lib/components/Footer.svelte";

  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { isLoggedIn } from "$lib/api.svelte";
  import { updateDeliveryMethods } from "$lib/e2/methods.svelte";
  import LoadingPage from "$lib/components/LoadingPage.svelte";

  let { children } = $props();

  onMount(() => {
    if (!$isLoggedIn) {
      goto("/");
    }
  });

  async function init() {
    await updateDeliveryMethods();
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
