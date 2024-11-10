<script lang="ts">
  import QRCode from "qrcode-svg";

  let {
    content = $bindable(),
    class: className,
    ...options
  }: QRCode.Options & { class: string } = $props();

  let qr = $derived(
    new QRCode({
      content,
      ...options,
    }),
  );
</script>

<div class="qr-code {className}">
  {@html qr.svg({ container: "svg-viewbox" })}
</div>

<style lang="scss">
  .qr-code {
    border-radius: var(--pico-border-radius);
    overflow: hidden;
    width: fit-content;

    :global svg {
      height: auto;
      width: 100%;
      aspect-ratio: 1/1;
    }
  }
</style>
