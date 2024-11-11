<script lang="ts">
  import Header from "$lib/components/Header.svelte";

  import { onNavigate } from "$app/navigation";

  let { children } = $props();

  onNavigate((navigation) => {
    if (!document.startViewTransition) return;

    return new Promise((resolve) => {
      document.startViewTransition(async () => {
        resolve();
        await navigation.complete;
      });
    });
  });
</script>

<Header />

<main class="container spaced-2">
  {@render children()}
</main>
