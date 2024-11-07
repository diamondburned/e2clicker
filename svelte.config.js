import adapter from "@sveltejs/adapter-node";
import { sveltePreprocess } from "svelte-preprocess";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  preprocess: [sveltePreprocess()],
  kit: {
    adapter: adapter({
      out: "dist/frontend",
      precompress: false,
    }),
    files: {
      assets: "assets",
      lib: "frontend/lib",
      params: "frontend/params",
      routes: "frontend/routes",
      serviceWorker: "frontend/worker",
      appTemplate: "frontend/app.html",
      errorTemplate: "frontend/error.html",
    },
  },
  compilerOptions: {
    runes: true,
  },
};

export default config;
