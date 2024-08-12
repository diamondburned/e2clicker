<script lang="ts">
  let { fixed = false } = $props();

  let header: HTMLElement;
  let scrolled = false;
</script>

<svelte:window
  onscroll={() => {
    header.classList.toggle("scrolled", window.scrollY > 0);
  }}
/>

<header id="header" bind:this={header} class:scrolled class:fixed>
  <nav class="container">
    <ul>
      <li>
        <img class="logo" src="/favicon.png" alt="logo" />
        <strong><a href="/">e2clicker</a></strong>
      </li>
    </ul>
    <ul></ul>
  </nav>
</header>

<style lang="scss">
  header {
    position: sticky;
    width: 100%;
    top: 0;

    transition: box-shadow 100ms var(--pico-transition-easing);

    &.scrolled {
      background-color: var(--pico-background-color);
      box-shadow:
        0 1px var(--pico-contrast-focus),
        0 -2px 4px 4px rgba(0, 0, 0, 0.07);
    }

    &.fixed {
      position: fixed;
    }
  }

  img.logo {
    width: 1.5em;
    height: auto;
    vertical-align: text-bottom;
    margin-right: 0.5em;
  }
</style>
