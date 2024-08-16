<script lang="ts">
    import { cn, renderTwClass } from "@lib/utils"
    import { type flexProps as Props } from "./index"

    type $$Props = Props
    export let bg: $$Props["bg"] = undefined
    export let className: $$Props["class"] = ""
    export let padding: $$Props["padding"] = undefined
    export let margin: $$Props["margin"] = undefined
    export let direction: $$Props["direction"] = "col"
    export let alignItems: $$Props["alignItems"] = "start"
    export let justify: $$Props["justify"] = "start"
    export let gap: Props["gap"] = undefined

    export { className as class }
</script>

<div
    class={cn([
        bg,
        renderTwClass("p", padding),
        renderTwClass("m", margin),
        "flex w-full",
        `flex-${direction}`,
        gap ? `gap-${gap}` : "",
        `items-${alignItems}`,
        `justify-${justify}`,
        "w-full",
        className,
    ])}
    {...$$restProps}
>
    <slot />
</div>
