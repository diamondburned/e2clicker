<script lang="ts">
  import { type AssignedToast, closeToast } from "$lib/toast.svelte";
  import { type Snippet } from "svelte";
  import Icon from "../Icon.svelte";

  let {
    toast,
  }: {
    toast: AssignedToast;
  } = $props();
</script>

<article class="toast flex">
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
