<script lang="ts">
  import type { Snippet } from "svelte";

  let {
    name,
    description,
    children,
    misc,
  }: {
    name: string;
    description?: Snippet;
    children?: Snippet;
    misc?: Snippet;
  } = $props();
</script>

<section class="preference-group">
  <header>
    <hgroup>
      <h2>{name}</h2>
      {#if description}
        <p>{@render description()}</p>
      {/if}
    </hgroup>
    {#if misc}
      <div class="misc">
        {@render misc()}
      </div>
    {/if}
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

    header {
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
