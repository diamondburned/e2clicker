import { defineConfig, presetMini } from "unocss";
import svelteExtractor from "@unocss/extractor-svelte";

const filteredRules = ["container", "outline"];
const mini = presetMini({ dark: "media" });

// Delete these so they don't conflict with pico.css's.
mini.rules = mini.rules!.filter(([match]) => filteredRules.every((rule) => rule !== match));
delete mini.theme!.container;
delete mini.theme!.containers;

export default defineConfig({
  presets: [mini],
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
