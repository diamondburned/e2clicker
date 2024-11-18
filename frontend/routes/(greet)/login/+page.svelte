<script lang="ts">
  import LoadingPage from "$lib/components/LoadingPage.svelte";
  import ResizeContainer from "$lib/components/ResizeContainer.svelte";

  import LoginScreen from "./LoginScreen.svelte";
  import RegisterScreen from "./RegisterScreen.svelte";

  import { slide } from "svelte/transition";
  import { page } from "$app/stores";
  import { isLoggedIn } from "$lib/api.svelte";
  import { goto } from "$app/navigation";

  let screen = $state<"login" | "register">(
    $page.url.hash == "#register" // set by homepage
      ? "register"
      : "login",
  );
  let promise = $state(Promise.resolve());
  let hideEstrannaise = $state(false);

  $effect(() => {
    if ($isLoggedIn) {
      goto("/dashboard");
    }
  });
</script>

<svelte:head>
  <title>Login - e2clicker</title>
</svelte:head>

<LoadingPage {promise} />

<div class="outer-container">
  <div class="main-wrapper container">
    <main id="login">
      <ResizeContainer>
        {#if screen == "login"}
          <LoginScreen bind:screen bind:promise />
        {:else}
          <RegisterScreen bind:screen bind:promise />
        {/if}
      </ResizeContainer>

      {#if !hideEstrannaise}
        <footer class="estrannaise spaced" out:slide={{ duration: 200, axis: "y" }}>
          <h2>Need to calculate your dosage?</h2>
          <p>
            <a href="https://estrannai.se/" target="_blank">estrannai.se</a> lets you simulate
            estrogen dosages instantly. This is a free and open-source tool that
            <span class="brand">e2clicker</span> relies on, so please consider supporting them!
          </p>
          <button
            class="hide-estrannaise minimal"
            onclick={() => {
              hideEstrannaise = true;
            }}
          >
            <small>Hide this message</small>
          </button>
        </footer>
      {/if}
    </main>
  </div>
</div>

<style lang="scss">
  @use "screen";

  .outer-container {
    height: 100dvh;

    display: flex;
    flex-direction: column;
  }

  .main-wrapper {
    flex: 1;

    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  h2 {
    font-size: 1.2em;
  }

  main {
    border-radius: var(--pico-border-radius);
    background: var(--pico-card-background-color);
    box-shadow: var(--pico-card-box-shadow);

    width: 100%;
    height: auto;

    footer {
      @include screen.spacings;

      padding: var(--pico-block-spacing);
      border-top: var(--pico-border-width) solid var(--pico-card-sectioning-background-color);
    }
  }
</style>
