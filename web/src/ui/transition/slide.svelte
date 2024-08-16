<script lang="ts">
    import { slide } from "svelte/transition"
    import type { slideProps as Props } from "./index"
    import { cn } from "@lib/utils"

    type $$Props = Props
    export let config: $$Props["config"] = undefined
    export let className: $$Props["class"] = ""
    export { className as class }
</script>

<div transition:slide={config} class={cn(className)}>
    <slot />
</div>
