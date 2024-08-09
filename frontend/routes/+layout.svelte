<script lang="ts">
  import "$lib/styles/styles.scss";
  import { setContext } from "svelte";

  let { children } = $props();

  setContext("iconCtx", {
    variation: "round",
  });
</script>

{@render children()}
