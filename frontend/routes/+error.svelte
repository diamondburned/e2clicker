<script lang="ts">
  import LoadingPage from "$lib/components/LoadingPage.svelte";

  import { page } from "$app/stores";

  let error = $derived(`Status ${$page.status}` + ($page.error ? `: ${$page.error.message}` : ""));
  let promise = $derived(new Promise((_, reject) => reject(error)));
</script>

<main class="container">
  <LoadingPage {promise}>
    {#snippet errorHeader()}
      <div class="image-container">
        <a href="https://x.com/rainbowemuweeps/status/946963954391179264" target="_blank">
          <img class="astolfo" alt="astolfo smug" src="/astolfo-smug.png" />
        </a>
      </div>
    {/snippet}
  </LoadingPage>
</main>

<style lang="scss">
  .image-container {
    display: flex;
    flex-direction: column;
    align-items: center;

    img {
      max-width: 300px;
      border-radius: var(--pico-border-radius);
    }

    img.astolfo {
      filter: // put outline around astolfo
        drop-shadow(2px 2px 0 var(--pico-primary)) drop-shadow(-2px 2px 0 var(--pico-primary))
        drop-shadow(2px -2px 0 var(--pico-primary)) drop-shadow(-2px -2px 0 var(--pico-primary))
        drop-shadow(2px 1px 0 var(--pico-primary)) drop-shadow(1px 2px 0 var(--pico-primary))
        drop-shadow(2px -1px 0 var(--pico-primary)) drop-shadow(-1px 2px 0 var(--pico-primary));
    }
  }
</style>
