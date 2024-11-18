<script lang="ts">
  import { type AssignedToast, closeToast } from "$lib/toasts";
  import { type Snippet } from "svelte";
  import Icon from "../Icon.svelte";
  import { slide } from "svelte/transition";

  let {
    toast,
  }: {
    toast: AssignedToast;
  } = $props();

  let isError = $derived(toast.urgency == "error");
</script>

<article
  class="toast flex items-center popping"
  class:error-box={isError}
  transition:slide={{ duration: 200 }}
>
  {#if toast.urgency == "error"}
    <Icon name="error" class="mr-2 text-red" />
  {/if}
  <div class="body flex-1">
    {@render render(toast.message)}
    {#if toast.description}
      <div class="toast-description text-sm">
        {@render render(toast.description)}
      </div>
    {/if}
  </div>
  <aside class="flex flex-col items-center">
    <button class="close minimal" aria-label="Close toast" onclick={() => closeToast(toast)}>
      <Icon name="close" />
    </button>
  </aside>
</article>

{#snippet render(snippet: Snippet | string)}
  {#if typeof snippet == "function"}
    {@render snippet()}
  {:else}
    {snippet}
  {/if}
{/snippet}

<style lang="scss">
  .toast {
    width: min(40ch, 100%);
    margin-top: 0;
  }
</style>
