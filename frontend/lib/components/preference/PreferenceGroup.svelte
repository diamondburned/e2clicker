<script lang="ts">
  import PreferenceLoader from "./PreferenceLoader.svelte";
  import type { Props as LoaderProps } from "./PreferenceLoader.svelte";

  import type { Snippet } from "svelte";

  let {
    name,
    description,
    children,
    header,
    misc,
    loader,
  }: {
    name: string;
    description?: Snippet;
    children?: Snippet;
    header?: Snippet;
    misc?: Snippet;
    loader?: LoaderProps;
  } = $props();
</script>

<section class="preference-group">
  <header>
    <div class="top-header">
      <hgroup class="m-0">
        <h2>{name}</h2>
        {#if description}
          <p>{@render description()}</p>
        {/if}
      </hgroup>
      <div class="misc">
        {@render misc?.()}
        {#if loader}
          <PreferenceLoader {...loader} />
        {/if}
      </div>
    </div>
    {@render header?.()}
  </header>

  <div class="preferences" role="list">
    {@render children?.()}
  </div>
</section>

<style lang="scss">
  @use "sass:map";
  @use "@picocss/pico/scss/settings" as *;

  .preference-group {
    --pico-block-spacing-vertical: calc(1.5 * var(--pico-spacing));

    @media (max-width: map.get(map.get($breakpoints, "sm"), "viewport")) {
      --pico-block-spacing-vertical: var(--pico-spacing);
    }

    .top-header {
      display: flex;
      align-items: center;
      gap: calc(var(--pico-spacing) / 2);

      margin-bottom: var(--pico-block-spacing-vertical);

      hgroup {
        flex: 1;
      }
    }
  }

  .preferences {
    display: flex;
    flex-direction: column;

    :global .preference {
      margin-bottom: var(--pico-block-spacing-vertical);

      &:last-child {
        margin-bottom: 0;
      }
    }
  }
</style>
