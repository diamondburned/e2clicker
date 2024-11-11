<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";

  import { fly } from "svelte/transition";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { isLoggedIn } from "$lib/api.js";
  import type { FlyProps } from "$lib/components/popovers/Popover.svelte";

  onMount(() => {
    if ($isLoggedIn) {
      // goto("/dashboard");
    }
  });

  let footerIndex = $state(0);
  let pauseFooter = $state(false);

  const footerIn: FlyProps = { duration: 200, x: 50 };
  const footerOut: FlyProps = { duration: 200, x: -50 };

  onMount(() => {
    const v = setInterval(() => {
      if (!pauseFooter) {
        footerIndex = (footerIndex + 1) % 4;
      }
    }, 5000);
    return () => clearInterval(v);
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

  <section>
    {#if $isLoggedIn}
      <a href="/dashboard" role="button">
        Go to Dashboard <Icon name="arrow-forward" />
      </a>
    {:else}
      <a href="/login" role="button">
        Get Started <Icon name="arrow-forward" />
      </a>
    {/if}
  </section>
</main>

<footer
  class="container secondary"
  class:focus={pauseFooter}
  onfocus={() => (pauseFooter = true)}
  onblur={() => (pauseFooter = false)}
  onmouseover={() => (pauseFooter = true)}
  onmouseleave={() => (pauseFooter = false)}
>
  <div class="footer-grid">
    {#if footerIndex == 0}
      <p in:fly={footerIn} out:fly={footerOut}>
        <a href="https://github.com/diamondburned/e2clicker">
          Source code on <b>GitHub</b>.
        </a>
      </p>
    {/if}
    {#if footerIndex == 1}
      <p in:fly={footerIn} out:fly={footerOut}>
        <a href="https://estrannai.se">
          Calculations sourced from
          <b>estrannai.se</b>.
        </a>
      </p>
    {/if}
    {#if footerIndex == 2}
      <p in:fly={footerIn} out:fly={footerOut}>
        <a href="https://www.gnu.org/licenses/gpl-3.0.en.html">
          Licensed under the <b>GNU GPL v3</b>, free and open-source forever.
        </a>
      </p>
    {/if}
    {#if footerIndex == 3}
      <p in:fly={footerIn} out:fly={footerOut}>🏳️‍⚧️ Trans rights are human rights.</p>
    {/if}
  </div>
</footer>

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
        box-shadow: var(--pico-box-shadow-thick);
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

  footer {
    width: 100%;
    opacity: 0.75;
    font-size: 0.8em;

    display: flex;
    align-items: center;
    flex-direction: column;

    padding-bottom: var(--pico-block-spacing-vertical);

    &.focus {
      opacity: 1;
    }

    .footer-grid {
      display: grid;
      grid-template-rows: 1fr;
      grid-template-columns: 1fr;

      > * {
        // Make them all overlap
        grid-row: 1;
        grid-column: 1;
      }
    }

    p {
      margin: 0;

      text-align: center;
      user-select: none;
    }
  }
</style>
