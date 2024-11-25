<script lang="ts">
  import "$lib/styles/google-fonts.css";
  import "$lib/styles/styles.scss";

  import LoadingPage from "$lib/components/LoadingPage.svelte";

  import { onNavigate } from "$app/navigation";
  import ToastOverlay from "$lib/components/toast/ToastOverlay.svelte";

  let { children } = $props();

  onNavigate((navigation) => {
    if (!document.startViewTransition) return;

    return new Promise((resolve) => {
      document.startViewTransition(async () => {
        resolve();
        await navigation.complete;
      });
    });
  });

  let errorPromise = $state<Promise<any> | null>(null);

  function setPromise(promise: Promise<any> | null) {
    if (!errorPromise) {
      errorPromise = promise;
    }
  }

  function setError(error: any) {
    setPromise(new Promise((_, reject) => reject(error)));
  }
</script>

<svelte:window
  onerror={(ev) => {
    ev.preventDefault();
    setError("An unknown browser error occured :(");
  }}
  onunhandledrejection={(ev) => {
    ev.preventDefault();
    try {
      // insert exception checks here.
      if (ev.reason.status == 401) return;
    } catch (e) {
      // do nothing
    }
    setError(ev.reason);
    console.error("An unhandled exception occured:", { reason: ev.reason });
  }}
/>

{#if errorPromise}
  <LoadingPage important promise={errorPromise} />
{/if}

<!--
  Force the browser to preload this font very early on.
  If we don't do this, then the font will only be loaded once our <Icon />
  component has been loaded by JavaScript code, which is way slower!
-->

<span class="material-symbols-rounded" style="position: fixed; top: -100px; left: -100px"
  >sentiment_satisfied</span
>

<ToastOverlay />

{@render children()}
