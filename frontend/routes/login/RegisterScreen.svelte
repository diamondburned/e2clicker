<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import Tooltip from "$lib/components/popovers/Tooltip.svelte";

  import { fade } from "svelte/transition";

  let {
    screen = $bindable(),
    promise = $bindable(),
  }: {
    screen: "login" | "register";
    promise: Promise<unknown>;
  } = $props();

  let registerName = $state("");

  async function submitRegister() {}
</script>

<article id="register" class="spaced" in:fade={{ duration: 200 }}>
  <h2>Create an Account</h2>

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
      <input type="text" name="name" placeholder="Alice" bind:value={registerName} />
    </label>
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
