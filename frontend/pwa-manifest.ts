import type { ManifestOptions } from "vite-plugin-pwa";

const manifest: Partial<ManifestOptions> = {
  name: "e2clicker",
  id: "/",
  scope: "/",
  start_url: "/dashboard",
  display: "standalone",
  theme_color: "#c72259",
  background_color: "#ffffff",
  icons: [
    {
      src: "/favicon.ico",
      sizes: "48x48",
    },
    {
      src: "/logo.png",
      type: "image/png",
      sizes: "512x512",
      purpose: "any",
    },
    {
      src: "/logo@2x.png",
      type: "image/png",
      sizes: "1024x1024",
      purpose: "any",
    },
    {
      src: "/logo.svg",
      type: "image/svg+xml",
      sizes: "any",
      purpose: "any",
    },
    {
      src: "/logo-maskable.png",
      type: "image/png",
      sizes: "512x512",
      purpose: "maskable",
    },
    {
      src: "/logo-maskable@2x.png",
      type: "image/png",
      sizes: "1024x1024",
      purpose: "maskable",
    },
    {
      src: "/logo-maskable.svg",
      type: "image/svg+xml",
      sizes: "any",
      purpose: "maskable",
    },
  ],
  screenshots: [
    {
      src: "/screenshots/dashboard-light.png",
      type: "image/png",
      sizes: "1223x707",
      label: "e2clicker dashboard",
      platform: "web",
      form_factor: "wide",
    },
    {
      src: "/screenshots/dashboard-dark.png",
      type: "image/png",
      sizes: "1223x707",
      label: "e2clicker dashboard (dark mode)",
      platform: "web",
      form_factor: "wide",
    },
    {
      src: "/screenshots/dashboard-light-mobile.png",
      type: "image/png",
      sizes: "750x1334",
      label: "e2clicker dashboard on mobile",
      platform: "android",
      form_factor: "narrow",
    },
    {
      src: "/screenshots/dashboard-dark-mobile.png",
      type: "image/png",
      sizes: "750x1334",
      label: "e2clicker dashboard on mobile (dark mode)",
      platform: "android",
      form_factor: "narrow",
    },
  ],
};

export default manifest;
