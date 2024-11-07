<script lang="ts">
  import "$lib/styles/styles.scss";
  import { onMount } from "svelte";
  import * as api from "$lib/api.svelte.js";

  let { children } = $props();

  onMount(api.init);
</script>

{@render children()}
