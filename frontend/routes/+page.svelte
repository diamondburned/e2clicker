<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import { isLoggedIn } from "$lib/api.js";

  onMount(() => {
    if (!$isLoggedIn) {
      goto("/login");
    }
  });
</script>
