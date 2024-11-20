import { defineConfig } from "vite";
import { sveltekit } from "@sveltejs/kit/vite";
import { SvelteKitPWA as sveltekitPWA } from "@vite-pwa/sveltekit";
import pwaManifest from "./frontend/pwa-manifest";

const devVMAddress = process.env.BACKEND_HTTP_ADDRESS;

export default defineConfig({
  // Why the FUCK is clearScreen true by default? That is fucking stupid.
  clearScreen: false,
  plugins: [
    sveltekit(),
    sveltekitPWA({
      injectRegister: "inline",
      strategies: "injectManifest",
      srcDir: "frontend",
      filename: "service-worker.ts",
      manifest: pwaManifest,
      devOptions: {
        enabled: true,
        type: "module",
      },
    }),
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
    emptyOutDir: true,
    rollupOptions: {
      treeshake: true,
      output: {
        format: "esm",
      },
    },
    target: ["esnext"],
    minify: true,
    cssMinify: true,
    cssTarget: ["chrome100", "firefox100", "safari15"],
    cssCodeSplit: true,
    sourcemap: true,
    reportCompressedSize: true,
  },
  esbuild: {
    sourcemap: true,
    treeShaking: true,
  },
});
