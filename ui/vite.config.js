import { defineConfig } from "vite";
import { VitePWA } from "vite-plugin-pwa";

const icons = [
    {
        src: "assets/icons/pwa-64x64.png",
        sizes: "64x64",
        type: "image/png",
    },
    {
        src: "assets/icons/pwa-192x192.png",
        sizes: "192x192",
        type: "image/png",
    },
    {
        src: "assets/icons/pwa-512x512.png",
        sizes: "512x512",
        type: "image/png",
    },
    {
        src: "assets/icons/maskable-icon-512x512.png",
        sizes: "512x512",
        type: "image/png",
        purpose: "maskable",
    },
];

const screenshots = [];

/** @type {import("vite-plugin-pwa").VitePWAOptions} */
const manifest = {
    strategies: "generateSW",
    registerType: "prompt",
    includeAssets: ["assets/fonts/bootstrap-icons.woff", "assets/fonts/bootstrap-icons.woff2"],
    manifest: {
        name: "PicoW LED",
        short_name: "picow-led",
        icons: icons,
        screenshots: screenshots,
        theme_color: "#09090b",
        background_color: "#09090b",
        display: "standalone",
        scope: ".",
        start_url: "./",
        //publicPath: "/picow-led.github.io",
    },
};

export default defineConfig({
    clearScreen: false,
    plugins: [VitePWA(manifest)],
});
