<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import Footer from "$lib/components/Footer.svelte";

  import { isLoggedIn } from "$lib/api.js";
  import { dev } from "$app/environment";
</script>

<svelte:head>
  <meta property="og:title" content="e2clicker" />
  <meta
    property="og:description"
    content="Easily track your estrogen levels and dosages with e2clicker, free and open-source forever for everyone!"
  />
  <meta property="og:image" content="/screenshots/dashboard-light.png" />
  <meta property="og:image:width" content="1223" />
  <meta property="og:image:height" content="707" />
  <meta property="og:image:alt" content="Screenshot of the e2clicker dashboard" />
</svelte:head>

<main class="container spaced">
  <header class="big-brand">
    <img class="logo" src="/logo.svg" alt="logo" />
    <h1 class="brand">e2clicker</h1>
  </header>

  <section class="big-screenshot">
    <img
      class="light"
      alt="Screenshot of the e2clicker dashboard"
      src="/screenshots/dashboard-light.png"
    />
    <img
      class="dark"
      alt="Screenshot of the e2clicker dashboard"
      src="/screenshots/dashboard-dark.png"
    />
  </section>

  <p class="yap">
    <span class="avoid-wrap">
      Easily track your estrogen levels and dosages with <span class="brand">e2clicker</span>,</span
    >
    <span class="avoid-wrap">free and open-source forever for everyone!</span>
  </p>

  <section>
    {#if dev}
      {#if $isLoggedIn}
        <a href="/dashboard" role="button">
          Go to Dashboard <Icon name="arrow-forward" />
        </a>
      {:else}
        <a href="/login" role="button">
          Get Started <Icon name="arrow-forward" />
        </a>
      {/if}
    {:else}
      <h3>Coming soon!</h3>
    {/if}
  </section>
</main>

<Footer no-expand />

<style lang="scss">
  main {
    --margin-y: clamp(
      var(--pico-block-spacing-vertical),
      6vh,
      calc(4 * var(--pico-block-spacing-vertical))
    );

    padding: var(--margin-y) var(--pico-spacing);

    font-size: clamp(1em, 5vw, 1.25em);
    max-width: 800px;
    flex: 1;

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
        box-shadow: var(--pico-box-shadow-thick);
      }

      .dark {
        display: none;
      }

      @media (prefers-color-scheme: dark) {
        .light {
          display: none;
        }
        .dark {
          display: block;
        }
      }
    }

    .yap {
      text-align: center;
      line-height: 1.5;

      .avoid-wrap {
        display: inline-block;
      }
    }
  }
</style>
