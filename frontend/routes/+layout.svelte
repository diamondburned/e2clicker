<script lang="ts">
  import "$lib/styles/google-fonts.css";
  import "$lib/styles/styles.scss";

  import { onNavigate } from "$app/navigation";

  let { children } = $props();

  onNavigate((navigation) => {
    if (!document.startViewTransition) return;

    return new Promise((resolve) => {
      document.startViewTransition(async () => {
        resolve();
        await navigation.complete;
      });
    });
  });
</script>

<!--
  Force the browser to preload this font very early on.
  If we don't do this, then the font will only be loaded once our <Icon />
  component has been loaded by JavaScript code, which is way slower!
-->

<span class="material-symbols-rounded" style="position: fixed; top: -100px; left: -100px"
  >sentiment_satisfied</span
>

{@render children()}
