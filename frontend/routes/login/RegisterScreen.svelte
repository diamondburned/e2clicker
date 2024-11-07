<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import Tooltip from "$lib/components/popovers/Tooltip.svelte";
  import { auth, register } from "$lib/openapi.gen";
  import { setToken } from "$lib/api";

  import { fade } from "svelte/transition";
  import ErrorBox from "$lib/components/ErrorBox.svelte";

  let {
    screen = $bindable(),
    promise = $bindable(),
  }: {
    screen: "login" | "register";
    promise: Promise<unknown>;
  } = $props();

  let registerName = $state("");
  let error = $state<any>();

  async function submitRegister() {
    try {
      const { secret } = await register({ name: registerName });
      const { token } = await auth({ secret });
      setToken(token);
    } catch (err) {
      error = err;
    }
  }
</script>

<article id="register" in:fade={{ duration: 200 }}>
  <h2>Create an Account</h2>

  <div class="content">
    <p>
      We only need a preferred name to address you. This name could be anything you want!
      <b>It will never be publicly visible.</b>
    </p>
    <label class="main-input">
      <span>Your preferred name:</span>
      <input type="text" name="name" placeholder="Alice" bind:value={registerName} />
    </label>

    <ErrorBox {error} prefix="cannot register" />
  </div>

  <div class="buttons">
    <button class="secondary" onclick={() => (screen = "login")}>Back</button>
    <button
      onclick={() => {
        promise = submitRegister();
      }}
    >
      Sign up <Icon name="arrow-forward" />
    </button>
  </div>
</article>

<style lang="scss">
  @use "screen";

  article {
    @include screen.article;
  }
</style>
