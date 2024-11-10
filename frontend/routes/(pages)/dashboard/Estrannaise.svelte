<!-- estrannai.se as a svelte component -->

<script lang="ts">
  import * as e2 from "$lib/e2";
  import * as api from "$lib/api";
  import * as models from "estrannaise/src/models";

  let {
    doseHistory,
  }: {
    doseHistory: api.DosageHistory["history"];
  } = $props();

  let conversionFactor = $state(1.0); // never changed for now
  let pkFunctions = $derived(models.PKFunctions(conversionFactor));
  let pkRandomFunctions = $derived(models.PKRandomFunctions(conversionFactor));

  $effect(() => {
    console.log({
      conversionFactor,
      pkFunctions,
      pkRandomFunctions,
    });
  });
</script>
