<script lang="ts">
  let {
    "no-expand": noExpand = false,
  }: {
    "no-expand"?: boolean;
  } = $props();

  let footer: HTMLElement;
  let footerHover = $state(false);
  let scrollHeight = $state(0);

  function scrollIntoView() {
    footer!.scrollIntoView({
      behavior: "instant",
      block: "end",
    });
  }

  $effect(() => {
    scrollHeight && footerHover && scrollIntoView();
  });
</script>

<svelte:body bind:clientHeight={scrollHeight} />

<footer
  class="text-xs"
  class:no-expand={noExpand}
  bind:this={footer}
  role="presentation"
  onmouseover={() => (footerHover = true)}
  onmouseleave={() => (footerHover = false)}
  onfocus={() => (footerHover = true)}
  onblur={() => (footerHover = false)}
>
  <div class="container">
    <p>
      <span>©</span>
      <a href="https://github.com/diamondburned/e2clicker" target="_blank">e2clicker</a>
      / license
      <a href="https://www.gnu.org/licenses/gpl-3.0.en.html" target="_blank">GPLv3</a>
      /
      <a href="/attribution" target="_blank">attributions</a>
    </p>
    <p>trans rights ❤️ 🏳️‍⚧️</p>
  </div>
</footer>

<style lang="scss">
  footer {
    border-top: var(--pico-border-width) solid var(--pico-muted-border-color);
    padding: var(--pico-block-spacing-vertical) var(--pico-spacing);

    opacity: 0.5;
    transition: all var(--pico-transition);
    &:hover {
      opacity: 1;
      &:not(.no-expand) {
        padding-top: calc(3 * var(--pico-block-spacing-vertical));
      }
    }

    .container {
      display: flex;
      flex-direction: row;
      justify-content: space-between;
    }

    p {
      margin: 0;
    }
  }
</style>
