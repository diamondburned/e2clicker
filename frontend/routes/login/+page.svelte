<script lang="ts">
  import Header from "$lib/components/Header.svelte";

  import PasswordIcon from "svelte-google-materialdesign-icons/Password.svelte";
  import KeyIcon from "svelte-google-materialdesign-icons/Key.svelte";
  import PersonAddIcon from "svelte-google-materialdesign-icons/Person_add.svelte";
  import ArrowBackIcon from "svelte-google-materialdesign-icons/Arrow_back_ios.svelte";
  import { fly } from "svelte/transition";

  function submitForm(event: SubmitEvent) {
    event.preventDefault();
    const formData = new FormData(event.target! as HTMLFormElement);
    const body = Object.fromEntries(formData);
    console.log(body);
  }

  let loginMethod = $state<"password" | "passkeys" | "signup" | null>(null);
  let flyDirection = $derived(loginMethod ? 1 : -1);

  let flyIn = $derived({ delay: 300, duration: 200, x: 10 * flyDirection });
  let flyOut = $derived({ duration: 200, x: -10 * flyDirection });
</script>

<div class="outer-container">
  <Header fixed />

  <div class="main-wrapper">
    <main id="login" class="container">
      {#snippet backHeader(method: string)}
        <header>
          <button class="back" onclick={() => (loginMethod = null)}>
            <div class="icon"><ArrowBackIcon size="24" /></div>
            Login Methods
          </button>
          <span class="method">{method}</span>
        </header>
      {/snippet}

      {#if loginMethod === null}
        <article in:fly={flyIn} out:fly={flyOut}>
          <header>
            <h2>Choose a login method</h2>
          </header>
          <section id="login-methods" class="buttons">
            <button class="outline" onclick={() => (loginMethod = "password")}>
              <div class="icon"><PasswordIcon size="48" /></div>
              <div>
                <h4>Password</h4>
                <p>Log in with your email and password</p>
              </div>
            </button>
            <button class="outline" title="Currently unsupported!" disabled>
              <div class="icon"><KeyIcon size="48" /></div>
              <div>
                <h4>Passkeys</h4>
                <p>Log in with your passkeys</p>
              </div>
            </button>
            <button class="outline contrast" onclick={() => (loginMethod = "signup")}>
              <div class="icon"><PersonAddIcon size="48" /></div>
              <div>
                <h4>Sign Up</h4>
                <p>Create a new account</p>
              </div>
            </button>
          </section>
        </article>
      {/if}

      {#if loginMethod === "password"}
        <article in:fly={flyIn} out:fly={flyOut}>
          {@render backHeader("Password")}
          <section>
            <form onsubmit={submitForm}>
              <fieldset>
                <label>
                  Email
                  <input type="email" name="email" required />
                </label>
                <label>
                  Password
                  <input type="password" name="password" required />
                </label>
              </fieldset>

              <input type="submit" value="Log in" />
            </form>
          </section>
        </article>
      {/if}

      {#if loginMethod === "passkeys"}{/if}

      {#if loginMethod === "signup"}
        <article in:fly={flyIn} out:fly={flyOut}>
          {@render backHeader("Sign up")}
          <section>
            <form onsubmit={submitForm}>
              <fieldset>
                <label>
                  Email
                  <input type="email" name="email" required />
                </label>
                <label>
                  Password
                  <input type="password" name="password" required />
                </label>
                <label>
                  Confirm Password
                  <input type="password" name="passwordConfirm" required />
                </label>
              </fieldset>
              <input type="submit" value="Sign up" />
            </form>
          </section>
        </article>
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

  h2 {
    font-size: 1.2em;
  }

  main {
    border-radius: var(--pico-border-radius);
    background: var(--pico-card-background-color);
    box-shadow: var(--pico-card-box-shadow);

    height: min(100%, 500px);
    padding: max(var(--pico-spacing), calc(5%));

    width: calc(100% - 2 * var(--pico-spacing));
    margin: var(--pico-spacing) auto;

    article {
      min-height: 100%;
      box-shadow: none;

      display: flex;
      flex-direction: column;
    }

    header,
    footer {
      height: 2em;
      background: none;

      * {
        margin: 0;
      }
    }

    header {
      display: flex;
      flex-direction: row;
      align-items: baseline;

      button.back {
        --pico-color: var(--pico-contrast);

        background: none;
        border: none;
        padding: 0;

        &:hover {
          --pico-color: var(--pico-primary-hover);
        }

        .icon {
          display: inline-block;
          vertical-align: bottom;
        }
      }

      .method {
        color: var(--pico-primary);
        font-weight: 600;

        border-left: 1px solid var(--pico-muted-color);
        margin-left: 0.5em;
        padding-left: 0.5em;
      }
    }

    section {
      flex: 1;

      display: flex;
      flex-direction: column;
      justify-content: center;
    }
  }

  #login-methods {
    display: flex;
    flex-direction: column;
    gap: var(--pico-spacing);

    button {
      display: flex;
      flex-direction: row;
      gap: var(--pico-spacing);
      text-align: left;

      div:nth-child(2) {
        display: flex;
        flex-direction: column;
        justify-content: center;

        h4,
        p {
          --pico-color: inherit;
          --pico-font-size: 1em;
          margin: 0;
        }

        p {
          --pico-line-height: 1;
        }
      }
    }
  }
</style>
