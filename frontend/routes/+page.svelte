<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { isLoggedIn } from "$lib/api.js";
  import Icon from "$lib/components/Icon.svelte";

  onMount(() => {
    if ($isLoggedIn) {
      // goto("/dashboard");
    }
  });
</script>

<main class="container spaced">
  <header class="big-brand">
    <img class="logo" src="/favicon.png" alt="logo" />
    <h1 class="brand">e2clicker</h1>
  </header>

  <section class="big-screenshot">
    <img alt="Screenshot of the e2clicker dashboard" src="/e2clicker-legacy.png" />
  </section>

  <p class="yap">
    <span class="avoid-wrap">
      Easily track your estrogen levels and dosages with <span class="brand">e2clicker</span>,</span
    >
    <span class="avoid-wrap">free and open-source forever for everyone!</span>
  </p>

  <footer>
    <a href="/login#register" role="button">
      Get Started <Icon name="arrow-forward" />
    </a>
  </footer>
</main>

<style lang="scss">
  main {
    --margin-y: clamp(
      var(--pico-block-spacing-vertical),
      6vh,
      calc(4 * var(--pico-block-spacing-vertical))
    );

    font-size: clamp(1em, 5vw, 1.25em);
    max-width: 800px;
    padding: var(--margin-y) var(--pico-spacing);

    display: grid;
    grid-template-rows: auto 1fr auto;

    & > * {
      --pico-block-spacing-vertical: var(--margin-y);

      margin-left: auto;
      margin-right: auto;
    }

    .big-brand {
      font-size: 1.35em;

      width: fit-content;

      display: flex;
      align-items: center;
      gap: calc(var(--pico-spacing) / 2);

      h1 {
        margin: 0;
      }

      img.logo {
        width: 3em;
      }
    }

    .big-screenshot {
      img {
        border: var(--pico-border-width) solid var(--pico-primary);
        border-radius: var(--pico-border-radius);
        box-shadow: var(--pico-card-box-shadow);
      }
    }

    .yap {
      text-align: center;

      .avoid-wrap {
        display: inline-block;
      }
    }
  }
</style>
