<script lang="ts">
  import { fly } from "svelte/transition";
  import type { FlyProps } from "$lib/components/popovers/Popover.svelte";
  import { onMount } from "svelte";

  let footerIndex = $state(0);
  let pauseFooter = $state(false);

  const footerIn: FlyProps = { duration: 200, x: 20 };
  const footerOut: FlyProps = { duration: 200, x: -20 };

  onMount(() => {
    const v = setInterval(() => {
      if (!pauseFooter) {
        footerIndex = (footerIndex + 1) % 4;
      }
    }, 3000);
    return () => clearInterval(v);
  });
</script>

<footer
  class="secondary"
  class:focus={pauseFooter}
  onfocus={() => (pauseFooter = true)}
  onblur={() => (pauseFooter = false)}
  onmouseover={() => (pauseFooter = true)}
  onmouseleave={() => (pauseFooter = false)}
>
  <div class="footer-grid container">
    {#if footerIndex == 0}
      <p in:fly={footerIn} out:fly={footerOut}>
        <a href="https://github.com/diamondburned/e2clicker" target="_blank">
          Source code on <b>GitHub</b>
        </a>
      </p>
    {/if}
    {#if footerIndex == 1}
      <p in:fly={footerIn} out:fly={footerOut}>
        <a href="https://estrannai.se" target="_blank">
          Calculations sourced from
          <b>estrannai.se</b>
        </a>
      </p>
    {/if}
    {#if footerIndex == 2}
      <p in:fly={footerIn} out:fly={footerOut}>
        <a href="https://www.gnu.org/licenses/gpl-3.0.en.html" target="_blank">
          Licensed under the <b>GNU GPL v3</b>, free and open-source forever
        </a>
      </p>
    {/if}
    {#if footerIndex == 3}
      <p in:fly={footerIn} out:fly={footerOut}>🏳️‍⚧️ Trans rights are human rights</p>
    {/if}
  </div>
</footer>

<style lang="scss">
  footer {
    width: 100%;
    opacity: 0.75;
    font-size: 0.8em;

    display: flex;
    align-items: center;
    flex-direction: column;

    padding: var(--pico-block-spacing-vertical) 0;

    &.focus {
      opacity: 1;
    }

    opacity: 0.5;
    transition: all var(--pico-transition);
    &:hover {
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
