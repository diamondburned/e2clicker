import { defineConfig } from "vite";
import { sveltekit } from "@sveltejs/kit/vite";

// import * as path from "path";
// const root = new URL(".", import.meta.url).pathname;

export default defineConfig({
  // Why the FUCK is clearScreen true by default? That is fucking stupid.
  clearScreen: false,
  plugins: [sveltekit()],
  server: {
    port: 5001,
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
  // https://github.com/vitejs/vite/issues/7385#issuecomment-1286606298
  // resolve: {
  //   alias: {
  //     "#/libdb.so/e2clicker": root,
  //   },
  // },
});

// if (import.meta.hot) {
//   // always reload the page on change because v86 is fragile
//   import.meta.hot.accept(() => import.meta.hot!.invalidate());
// }
