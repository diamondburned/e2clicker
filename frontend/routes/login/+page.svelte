<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import Header from "$lib/components/Header.svelte";
  import Tooltip from "$lib/components/popovers/Tooltip.svelte";
  import ResizeContainer from "$lib/components/ResizeContainer.svelte";
  import TextHorizontalRule from "$lib/components/TextHorizontalRule.svelte";

  import { fade, slide } from "svelte/transition";

  // function submitForm(event: SubmitEvent) {
  //   event.preventDefault();
  //   const formData = new FormData(event.target! as HTMLFormElement);
  //   const body = Object.fromEntries(formData);
  //   console.log(body);
  // }

  let signUp = $state(false);
  let hideEstrannaise = $state(false);
</script>

<div class="outer-container">
  <Header />

  <div class="main-wrapper container">
    <main id="login">
      <ResizeContainer>
        {#if !signUp}
          <article id="login" class="spaced" in:fade={{ duration: 200 }}>
            <h2>Login</h2>

            <div class="content">
              <p>
                Scan the secret QR code:
                <span style="float: right">
                  <Tooltip tooltip={loginTooltip}>
                    <Icon name="info" />
                  </Tooltip>
                  {#snippet loginTooltip()}
                    <div class="login-tooltip">
                      <p>On the device that you're already logged in:</p>
                      <ul>
                        <li>Go to your Settings</li>
                        <li>Choose "Show secret QR code"</li>
                        <li>Scan the with this device.</li>
                      </ul>
                    </div>
                  {/snippet}
                </span>
              </p>
              <button class="secondary outline">
                Log in with QR <Icon name="qr-code-scanner" />
              </button>

              <TextHorizontalRule>or</TextHorizontalRule>

              <label class="main-input">
                <span>Input the secret manually:</span>
                <input type="password" name="secret" placeholder="xxxxxxxxxxxxxxxxxxxx" />
              </label>
            </div>

            <div class="buttons">
              <button class="secondary" onclick={() => (signUp = true)}> Sign up </button>
              <button>
                Login <Icon name="arrow-forward" />
              </button>
            </div>
          </article>
        {:else}
          <article id="signup" class="spaced" in:fade={{ duration: 200 }}>
            <h2>Sign up</h2>

            <div class="content">
              <label class="main-input">
                <span>Your preferred name:</span>
                <span style="float: right">
                  <Tooltip tooltip={preferredNameTooltip}>
                    <Icon name="info" />
                  </Tooltip>
                  {#snippet preferredNameTooltip()}
                    <div class="spaced">
                      <p class="preferred-name-tooltip">
                        This will only be used to address you in the app.
                        <b>It will not be shown to anyone else.</b>
                      </p>
                    </div>
                  {/snippet}
                </span>
                <input type="text" name="name" placeholder="Alice" />
              </label>
            </div>

            <div class="buttons">
              <button class="secondary" onclick={() => (signUp = false)}>Back</button>
              <button>
                Sign up <Icon name="arrow-forward" />
              </button>
            </div>
          </article>
        {/if}
      </ResizeContainer>

      {#if !hideEstrannaise}
        <footer class="estrannaise spaced" out:slide={{ duration: 200, axis: "y" }}>
          <h2>Need to calculate your dosage?</h2>
          <p>
            <a href="https://estrannai.se/">estrannai.se</a> lets you simulate estrogen dosages
            instantly. This is a free and open-source tool that
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

  .login-tooltip {
    --pico-typography-spacing-vertical: 0.35em;
    ul {
      margin-bottom: 0;
      padding-left: var(--pico-spacing);
      li {
        list-style-type: disc;
      }
    }
  }

  h2 {
    font-size: 1.2em;
  }

  main {
    --pico-block-spacing-vertical: clamp(var(--pico-spacing), 5vh, calc(var(--pico-spacing) * 3));
    --pico-block-spacing-horizontal: clamp(var(--pico-spacing), 5%, calc(var(--pico-spacing) * 3));
    --pico-block-spacing: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);

    border-radius: var(--pico-border-radius);
    background: var(--pico-card-background-color);
    box-shadow: var(--pico-card-box-shadow);

    width: 100%;
    height: auto;

    article {
      box-shadow: none;
      margin: 0;
      padding: var(--pico-block-spacing-vertical) var(--pico-block-spacing-horizontal);
      padding-bottom: var(--pico-block-spacing-horizontal);
      background: none;

      display: flex;
      flex-direction: column;
    }

    footer {
      padding: var(--pico-block-spacing);
      border-top: var(--pico-border-width) solid var(--pico-card-sectioning-background-color);
    }

    .content {
      width: 100%;

      & > button {
        width: 100%;
      }
    }

    .buttons {
      display: flex;
      justify-content: flex-end;
      gap: calc(var(--pico-spacing) / 2);

      button:last-child {
        justify-self: end;
        width: 7em;
      }
    }

    label.main-input {
      &,
      input {
        margin-bottom: 0;
      }
    }
  }
</style>
