<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import Dialog from "$lib/components/Dialog.svelte";
  import Tooltip from "$lib/components/popovers/Tooltip.svelte";
  import ErrorBox from "$lib/components/ErrorBox.svelte";
  import QRScanner from "$lib/components/QRScanner.svelte";
  import TextHorizontalRule from "$lib/components/TextHorizontalRule.svelte";

  import { fade } from "svelte/transition";
  import { user, auth, secretQRRegex } from "$lib/api.svelte";
  import { get } from "svelte/store";

  let {
    screen = $bindable(),
    promise = $bindable(),
  }: {
    screen: "login" | "register";
    promise: Promise<unknown>;
  } = $props();

  let error = $state<any>();
  let loginSecret = $state(get(user)?.secret ?? ""); // do not bind or the user will lose this!
  let showQRDialog = $state(false);

  async function submitLogin() {
    const secret = loginSecret.replace(/\s*/g, "").toUpperCase();
    try {
      await auth(secret);
    } catch (err) {
      error = err;
    }
  }
</script>

<article id="login" in:fade={{ duration: 200 }}>
  <h2>Login</h2>

  <div class="content">
    <div>
      <span>Scan the secret QR code:</span>
      <span style="float: right">
        <Tooltip tooltip={loginTooltip} selectable>
          <Icon name="info" />
        </Tooltip>
      </span>
    </div>
    <button class="secondary outline" onclick={() => (showQRDialog = true)}>
      Log in with QR <Icon name="qr-code-scanner" />
    </button>

    <TextHorizontalRule>or</TextHorizontalRule>

    <label class="main-input">
      <span>Input the secret manually:</span>
      <input
        type="password"
        name="secret"
        class="font-mono"
        placeholder="xxxxxxxxxxxxxxxxxxxx"
        bind:value={loginSecret}
      />
    </label>

    <ErrorBox {error} prefix="cannot log in" />
  </div>

  <div class="buttons">
    <button class="secondary" onclick={() => (screen = "register")}> Sign up </button>
    <button
      onclick={() => {
        promise = submitLogin();
      }}
    >
      Login <Icon name="arrow-forward" />
    </button>
  </div>
</article>

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

<Dialog wide dismissible bind:open={showQRDialog}>
  <header>
    <h4>Scan Secret QR</h4>
  </header>

  <QRScanner
    onscan={(result) => {
      const match = result.match(secretQRRegex);
      if (!match) {
        console.log("Discarding non-token QR", { result });
        return;
      }

      showQRDialog = false;
      loginSecret = match[1];
      promise = submitLogin();
    }}
  />

  <footer>
    <button aria-label="Cancel" onclick={() => (showQRDialog = false)}>
      Cancel <Icon name="close" />
    </button>
  </footer>
</Dialog>

<style lang="scss">
  @use "screen";

  article {
    @include screen.article;
  }

  .content {
    button {
      margin-top: calc(var(--pico-spacing) * 0.25);
    }
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
</style>
