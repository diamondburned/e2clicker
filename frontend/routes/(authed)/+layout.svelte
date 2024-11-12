<script lang="ts">
  import Header from "$lib/components/Header.svelte";
  import Footer from "$lib/components/Footer.svelte";

  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { isLoggedIn } from "$lib/api";

  let { children } = $props();

  onMount(() => {
    if (!$isLoggedIn) {
      goto("/");
    }
  });
</script>

<Header />

<main class="container spaced-2">
  {@render children()}
</main>

<Footer />
