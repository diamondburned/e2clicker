import { defineConfig, presetUno } from "unocss";
import svelteExtractor from "@unocss/extractor-svelte";

const filteredRules = [/^__container$/, "outline"];
const filterRule = (m: string | RegExp) => filteredRules.every((r) => r != m);

const unoPreset = presetUno({ dark: "media" });
// Delete the rules inside filteredRules.
unoPreset.rules = unoPreset.rules!.filter(([m]) => filterRule(m));
// Shortcuts only contains the container, so delete all of it.
unoPreset.shortcuts = [];

export default defineConfig({
  presets: [unoPreset],
  extractors: [svelteExtractor()],
  theme: {
    breakpoints: {
      sm: "576px",
      md: "768px",
      lg: "1024px",
      xl: "1280px",
      xxl: "1536px",
    },
    fontSize: {
      xs: ["var(--font-size-xs)", "var(--line-height-xs)"],
      sm: ["var(--font-size-sm)", "var(--line-height-sm)"],
      base: ["var(--font-size-base)", "var(--line-height-base)"],
      xl: ["var(--font-size-xl)", "var(--line-height-xl)"],
      "2xl": ["var(--font-size-2xl)", "var(--line-height-2xl)"],
      "3xl": ["var(--font-size-3xl)", "var(--line-height-3xl)"],
      "4xl": ["var(--font-size-4xl)", "var(--line-height-4xl)"],
      "5xl": ["var(--font-size-5xl)", "var(--line-height-5xl)"],
    },
  },
});
