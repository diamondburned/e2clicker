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
      xs: "0.65em",
      sm: "0.8em",
      base: "1em",
      xl: "1.25em",
      "2xl": "1.563em",
      "3xl": "1.953em",
      "4xl": "2.441em",
      "5xl": "3.052em",
    },
  },
});
