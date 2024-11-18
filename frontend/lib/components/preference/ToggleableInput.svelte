<script lang="ts" generics="ValueT, ZeroT">
  import { Snippet } from "svelte";

  let {
    value = $bindable(),
    zero,
    children,
  }: {
    value: ValueT | ZeroT; // the actual value of the input
    zero: ZeroT; // the value to use when the input is unchecked
    children: Snippet;
  } = $props();
</script>

<div class="toggleable-input"></div>
