<script lang="ts">
  import QRScanner from "qr-scanner";
  import { onDestroy, onMount } from "svelte";

  let {
    onscan,
  }: {
    onscan: (data: string) => void | Promise<void>;
  } = $props();

  let video: HTMLVideoElement;
  let loading = $state<Promise<void> | null>(null);
  let scanner = $state<QRScanner | null>(null);
  let visible = $state(false);
  let destroyed = $state(false);
  let processing = $state(false);

  async function initialize(): Promise<void> {
    if (!(await QRScanner.hasCamera())) {
      throw new Error("No camera found.");
    }

    scanner = new QRScanner(
      video,
      async (result) => {
        if (processing) {
          return;
        }

        processing = true;
        await onscan(result.data);

        // Only resume scanning after 3 seconds.
        setTimeout(() => (processing = false), 3000);
      },
      {
        returnDetailedScanResult: true,
        highlightScanRegion: true,
        highlightCodeOutline: true,
      },
    );

    await scanner.start();
  }

  function destroy(scanner: QRScanner | null) {
    if (scanner) {
      scanner.stop();
      scanner.destroy();
    }
  }

  onMount(async () => {
    loading = initialize();
    try {
      await loading;
      visible = true;
    } finally {
      if (destroyed) {
        destroy(scanner);
      }
    }
  });

  onDestroy(() => {
    destroyed = true;
    destroy(scanner);
  });
</script>

<div class="qr-scanner-box">
  {#await loading}
    <p class="status loading" aria-busy="true">Loading...</p>
  {:catch error}
    <p class="status error">Failed to load camera: {error.message}</p>
  {/await}

  <div class="qr-scanner" class:visible>
    <video class="view" bind:this={video}>
      <track kind="captions" />
    </video>
  </div>
</div>

<style lang="scss">
  video {
    width: auto;
    height: auto;
    max-width: 100%;
    max-height: calc(100vh - 5em);

    margin: 0 auto;
    display: block;

    scroll-margin-top: 1em;

    border-radius: var(--pico-border-radius);
    box-shadow: var(--pico-card-box-shadow);
  }

  .status {
    margin: var(--pico-typography-spacing-vertical) 0;
    margin-bottom: 0;
    text-align: center;

    &.error {
      color: var(--pico-color-red);
    }

    &.error::before {
      content: "Error: ";
      font-weight: bold;
    }
  }

  .qr-scanner {
    &:not(.visible) {
      display: none;
    }

    :global(.scan-region-highlight-svg) {
      stroke: var(--pico-primary) !important;
    }

    :global(.code-outline-highlight) {
      stroke: var(--pico-secondary) !important;
    }
  }
</style>
