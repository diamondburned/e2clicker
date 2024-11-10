import { defineConfig, presetUno } from "unocss";
import svelteExtractor from "@unocss/extractor-svelte";

export default defineConfig({
  presets: [presetUno({ dark: "media" })],
  extractors: [svelteExtractor()],
});
