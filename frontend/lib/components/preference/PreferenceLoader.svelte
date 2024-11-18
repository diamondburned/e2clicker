<script lang="ts" module>
  export type Props = {
    loader: Promise<any> | Pick<api.AnyAsyncToOK, "promise">;
    valid?: boolean;
  };
</script>

<script lang="ts">
  import Icon from "$lib/components/Icon.svelte";
  import ErrorBox from "$lib/components/ErrorBox.svelte";

  import type * as api from "$lib/api.svelte";

  let { loader, valid }: Props = $props();
</script>

{#if valid == false}
  <span class="text-red">
    <Icon name="error" />
  </span>
{:else}
  {#await "promise" in loader ? loader.promise : loader}
    <span aria-busy="true"></span>
  {:then}
    <span class="text-green">Saved <Icon name="done" /></span>
  {:catch error}
    <ErrorBox tiny {error} />
  {/await}
{/if}
