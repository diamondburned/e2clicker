<script lang="ts">
  import Header from "$lib/components/Header.svelte";
  import LoadingScreen from "$lib/components/LoadingScreen.svelte";
  import ResizeContainer from "$lib/components/ResizeContainer.svelte";

  import LoginScreen from "./LoginScreen.svelte";
  import RegisterScreen from "./RegisterScreen.svelte";

  import { slide } from "svelte/transition";

  let screen = $state<"login" | "register">("login");
  let promise = $state(Promise.resolve());
  let hideEstrannaise = $state(false);
</script>

<LoadingScreen {promise} />

<div class="outer-container">
  <Header />

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
          <a
            href="#hide-estrannaise"
            class="hide-estrannaise"
            onclick={() => {
              hideEstrannaise = true;
            }}
          >
            <small>Hide this message</small>
          </a>
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
