<script lang="ts">
  import Header from "$lib/components/Header.svelte";

  import { fly } from "svelte/transition";
  import { page } from "$app/stores";
  import { cubicIn, cubicOut } from "svelte/easing";

  let { children } = $props();

  const prefersReducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)").matches;
  const duration = prefersReducedMotion ? 0 : 300;
  const delay = duration + (prefersReducedMotion ? 0 : 100);
  const y = 10;

  const transitionIn = { easing: cubicOut, y, duration, delay };
  const transitionOut = { easing: cubicIn, y: -y, duration };
</script>

<Header />

{#key $page.url.pathname}
  <div in:fly={transitionIn} out:fly={transitionOut}>
    <main class="container spaced-2">
      {@render children()}
    </main>
  </div>
{/key}

<style lang="scss">
  div {
    width: 100%;
    height: 100%;

    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
  }
</style>
