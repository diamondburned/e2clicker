import { defineConfig } from "vite";
import { sveltekit } from "@sveltejs/kit/vite";

// import * as path from "path";
// const root = new URL(".", import.meta.url).pathname;

const devVMAddress = process.env.BACKEND_HTTP_ADDRESS;

export default defineConfig({
  // Why the FUCK is clearScreen true by default? That is fucking stupid.
  clearScreen: false,
  plugins: [sveltekit()],
  server: {
    host: "0.0.0.0",
    port: 8000,
    watch: {
      ignored: [".direnv/**", ".svelte-kit/**", "dist/**"],
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
      output: {
        format: "esm",
      },
    },
    target: "esnext",
    sourcemap: true,
    reportCompressedSize: true,
    // Fix estrannaise using require() syntax.
    commonjsOptions: { transformMixedEsModules: true },
  },
  esbuild: {
    sourcemap: true,
  },
});

// if (import.meta.hot) {
//   // always reload the page on change because v86 is fragile
//   import.meta.hot.accept(() => import.meta.hot!.invalidate());
// }
