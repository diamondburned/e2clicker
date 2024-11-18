<script lang="ts">
  import type { Snippet } from "svelte";
  import type { HTMLLabelAttributes } from "svelte/elements";

  let {
    name,
    description,
    children,
    ...attrs
  }: {
    name: string | Snippet;
    description?: Snippet;
    children?: Snippet;
  } & HTMLLabelAttributes = $props();
</script>

<label class="preference" {...attrs}>
  <hgroup class="preference-info">
    <h4 class="title">
      {#if typeof name === "string"}
        {@html name}
      {:else}
        {@render name()}
      {/if}
    </h4>
    {#if description}
      <p class="description"><small>{@render description()}</small></p>
    {/if}
  </hgroup>
  <div class="preference-content">
    {@render children?.()}
  </div>
</label>

<style lang="scss">
  @use "sass:map";
  @use "@picocss/pico/scss/settings" as *;

  .preference {
    display: flex;
    align-items: center;
    gap: var(--pico-spacing);

    .preference-info {
      flex: 1;
      margin: 0;

      .title {
        color: inherit;
      }

      .description {
        margin-top: calc(var(--pico-spacing) / 4);
        max-width: 600px;
      }
    }

    .preference-content {
      display: flex;
      flex-direction: row;
      justify-content: center;

      max-width: 40%;
      min-width: 10%;
      margin: 0;

      :global {
        * {
          margin-bottom: 0;
        }
      }
    }

    @media (max-width: map.get(map.get($breakpoints, "sm"), "breakpoint")) {
      flex-direction: column;
      align-items: stretch;
      gap: calc(var(--pico-spacing) / 2);

      .preference-info,
      .preference-content {
        max-width: initial;
      }
    }
  }
</style>
