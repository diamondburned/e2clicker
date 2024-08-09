<script lang="ts">
  import { goto } from "$app/navigation";
  import { isLoggedIn } from "$lib/api.svelte.js";
  import { onMount } from "svelte";

  onMount(() => {
    if (!isLoggedIn()) {
      goto("/login");
    }
  });
</script>
