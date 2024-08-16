<script lang="ts">
    import { type boxProps as Props } from "./index"
    import { cn, renderTwClass } from "@lib/utils"

    type $$Props = Props
    export let bg: $$Props["bg"] = undefined
    export let className: $$Props["class"] = undefined
    export let padding: $$Props["padding"] = undefined
    export let margin: $$Props["margin"] = undefined

    export { className as class }
</script>

<div class={cn([bg, renderTwClass("p", padding), renderTwClass("p", margin), className])} {...$$restProps}>
    <slot />
</div>
