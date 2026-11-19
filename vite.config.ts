import { defineConfig } from "vite";
import adapter from "@sveltejs/adapter-node";
import unoCSS from "@unocss/svelte-scoped/preprocess";
import { sveltePreprocess } from "svelte-preprocess";
import { withPwa } from "@vite-pwa/sveltekit";
import pwaManifest from "./frontend/pwa-manifest.ts";

const devVMAddress = process.env.BACKEND_HTTP_ADDRESS;

export default defineConfig({
  // Why the FUCK is clearScreen true by default? That is fucking stupid.
  clearScreen: false,
  plugins: [
    withPwa(
      {
        adapter: adapter({
          out: "dist/frontend",
          precompress: true,
        }),
        files: {
          assets: "assets",
          routes: "frontend/routes",
          appTemplate: "frontend/app.html",
          errorTemplate: "frontend/error.html",
        },
        serviceWorker: {
          register: false,
        },
        preprocess: [
          sveltePreprocess(),
          unoCSS({
            combine: false,
            classPrefix: "uno-",
          }),
        ],
        compilerOptions: {
          runes: true,
        },
      },
      {
        strategies: "buildSW",
        buildSW: {
          swSrc: "frontend/service-worker.ts",
        },
        manifest: pwaManifest,
        devOptions: {
          enabled: true,
          type: "module",
          navigateFallback: "/",
        },
      },
    ),
  ],
  server: {
    host: "0.0.0.0",
    port: 8000,
    watch: {
      ignored: ["**/.direnv/**", "**/dist/**", "**/result/**"],
    },
    proxy: (() => {
      if (devVMAddress) {
        console.log("Enabling backend reverse proxy in Vite.");
        console.log("  /api ->", devVMAddress);
        return {
          "/api": devVMAddress,
        };
      }
    })(),
  },
  build: {
    assetsDir: "static",
    target: ["esnext"],
    minify: true,
    cssMinify: true,
    cssTarget: ["chrome100", "firefox100", "safari15"],
    cssCodeSplit: true,
    sourcemap: true,
    reportCompressedSize: true,
  },
});
