import { defineConfig } from "vite"
import { svelte } from "@sveltejs/vite-plugin-svelte"
import path from "path"

// https://vitejs.dev/config/
export default defineConfig({
    server: { port: 1437 },
    resolve: {
        alias: {
            "@lib": path.resolve(__dirname, "./src/lib"),
            "@ui": path.resolve(__dirname, "./src/ui"),
            "@routes": path.resolve(__dirname, "./src/routes"),
        },
    },

    plugins: [svelte()],
})
