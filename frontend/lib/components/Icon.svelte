<script lang="ts">
  import type { SvelteHTMLElements } from "svelte/elements";

  let {
    name,
    class: klass = "",
    ...attrs
  }: {
    name: string;
    class?: string;
  } & SvelteHTMLElements["span"] = $props();
</script>

<span class="material-symbols-rounded {klass}" {...attrs}>{name.replaceAll("-", "_")}</span>

<style lang="scss">
  .material-symbols-rounded {
    vertical-align: bottom;
    max-width: 1em;
    overflow: hidden;
    font-size: 1.25em;
  }
</style>
