<script lang="ts">
  import Icon from "./Icon.svelte";

  import * as api from "$lib/api";
  import { isLoggedIn } from "$lib/api";

  let { fixed = false } = $props();

  let header: HTMLElement;
  let scrolled = false;

  let me = $state<api.User | null>(null);

  $effect(() => {
    if ($isLoggedIn) {
      api.currentUser().then((u) => {
        me = u;
      });
    }
  });
</script>

<svelte:window
  onscroll={() => {
    header.classList.toggle("scrolled", window.scrollY > 0);
  }}
/>

<header id="header" bind:this={header} class:scrolled class:fixed>
  <nav class="container">
    <ul>
      <li class="brand">
        <img class="logo" src="/logo.svg" alt="logo" />
        <strong><a href={$isLoggedIn ? "/dashboard" : "/"}>e2clicker</a></strong>
      </li>
    </ul>
    <ul>
      {#if me}
        <li class="current-user">
          <a href="/settings">
            <span class="name">{me.name}</span>
            <div class="avatar">
              <Icon name="person" />
            </div>
          </a>
        </li>
      {/if}
    </ul>
  </nav>
</header>

<style lang="scss">
  header {
    position: sticky;
    width: 100%;
    top: 0;

    z-index: 50;
    user-select: none;

    transition: box-shadow 100ms var(--pico-transition-easing);

    &.scrolled {
      background-color: var(--pico-background-color);
      box-shadow:
        0 1px var(--pico-contrast-focus),
        0 -2px 4px 4px rgba(0, 0, 0, 0.07);
    }

    a:hover {
      text-decoration: none;
    }
  }

  .brand {
    height: 1.5rem; // prevent jiggling when loaded

    padding-top: 0;
    padding-bottom: 0;

    .logo {
      width: 1.5rem;
      height: 1.5rem;
      vertical-align: text-bottom;
      margin-right: 0.5em;
    }

    a {
      padding-top: 0;
      padding-bottom: 0;
    }
  }

  .current-user {
    .avatar {
      display: inline-block;
    }
  }
</style>
