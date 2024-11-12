<script lang="ts">
  import "$lib/styles/styles.scss";
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

<svelte:head>
  <!-- No render blocking -->
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Material+Symbols+Rounded:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200"
  />
  <link
    rel="stylesheet"
    href="https://fonts.googleapis.com/css2?family=Inconsolata:wght@200..900&family=Inter:ital,opsz,wght@0,14..32,100..900;1,14..32,100..900&family=Reddit+Sans:ital,wght@0,200..900;1,200..900&display=swap"
  />
</svelte:head>

{@render children()}
