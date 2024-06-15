import { defineConfig } from "vite"
import { svelte } from "@sveltejs/vite-plugin-svelte"
import path from "path"

// https://vitejs.dev/config/
export default defineConfig({
    resolve: {
        alias: {
            "@lib": path.resolve(__dirname, "./src/lib"),
            "@routes": path.resolve(__dirname, "./src/routes"),
        },
    },

    plugins: [svelte()],
})
