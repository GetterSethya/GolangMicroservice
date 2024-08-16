<script lang="ts">
    import { fade } from "svelte/transition"
    import type { fadeProps as Props } from "./index"
    import { cn } from "@lib/utils"

    type $$Props = Props
    export let config: $$Props["config"] = undefined
    export let className: $$Props["class"] = ""
    export { className as class }
</script>

<div transition:fade={config} class={cn(className)}>
    <slot />
</div>
