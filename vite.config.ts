import { defineConfig } from "vite";
import { sveltekit } from "@sveltejs/kit/vite";

// import * as path from "path";
// const root = new URL(".", import.meta.url).pathname;

const backendHTTPAddress = process.env.BACKEND_HTTP_ADDRESS;

export default defineConfig({
  // Why the FUCK is clearScreen true by default? That is fucking stupid.
  clearScreen: false,
  plugins: [sveltekit()],
  server: {
    port: 8000,
    proxy: (() => {
      if (backendHTTPAddress) {
        console.log("Enabling backend reverse proxy in Vite.");
        console.log("  /api ->", backendHTTPAddress);
        return {
          "/api": backendHTTPAddress,
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
  },
  esbuild: {
    sourcemap: true,
  },
});

// if (import.meta.hot) {
//   // always reload the page on change because v86 is fragile
//   import.meta.hot.accept(() => import.meta.hot!.invalidate());
// }
