import adapter from "@sveltejs/adapter-node";

/** @type {import('@sveltejs/kit').Config} */
const config = {
  kit: {
    adapter: adapter({
      out: "dist/frontend",
      precompress: false,
    }),
    files: {
      assets: "frontend/public",
      lib: "frontend/lib",
      params: "frontend/params",
      routes: "frontend/routes",
      serviceWorker: "frontend/worker",
      appTemplate: "frontend/app.html",
      errorTemplate: "frontend/error.html",
    },
  },
  compilerOptions: {},
};

export default config;
