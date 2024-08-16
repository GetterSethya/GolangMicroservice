import { join } from "path"
import { myCustomTheme } from "./my-custom-theme"
import type { Config } from "tailwindcss"

// 1. Import the Skeleton plugin
import { skeleton } from "@skeletonlabs/tw-plugin"

const config = {
    // 2. Opt for dark mode to be handled via the class method
    darkMode: "class",
    content: [
        "./src/lib/**/*.{html,js,svelte,ts}",
        "./src/routes/**/*.{html,js,svelte,ts}",
        "./src/*.{html,js,svelte,ts}",
        "./src/ui/**/*.{html,js,svelte,ts}",
        "./*.html",
        // 3. Append the path to the Skeleton package
        join(require.resolve("@skeletonlabs/skeleton"), "../**/*.{html,js,svelte,ts}"),
    ],
    theme: {
        extend: {},
    },
    safelist: [
        { pattern: /^m-/ },
        { pattern: /^ms-/ },
        { pattern: /^me-/ },
        { pattern: /^mt-/ },
        { pattern: /^mb-/ },
        { pattern: /^mx-/ },
        { pattern: /^my-/ },
        { pattern: /^p-/ },
        { pattern: /^ps-/ },
        { pattern: /^pe-/ },
        { pattern: /^pt-/ },
        { pattern: /^pb-/ },
        { pattern: /^px-/ },
        { pattern: /^py-/ },
        { pattern: /^justify-/ },
        { pattern: /^items-/ },
        { pattern: /^flex-/ },
        { pattern: /^variant-/ },
        { pattern: /^gap-/ },
    ],
    plugins: [
        // 4. Append the Skeleton plugin (after other plugins)
        skeleton({
            themes: { custom: [myCustomTheme] },
        }),
    ],
} satisfies Config

export default config
