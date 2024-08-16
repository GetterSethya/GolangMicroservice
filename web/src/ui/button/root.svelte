<script lang="ts">
    import { type ButtonProps as Props, type ButtonEvents as Events } from "./index"
    import { cn } from "@lib/utils"
    import Flex from "@ui/container/flex.svelte"
    import { hoverClasses } from "@ui/types"

    type $$Props = Props
    type $$Events = Events

    export let className: $$Props["class"] = undefined
    export let bg: $$Props["bg"] = undefined
    export let variant: $$Props["variant"] = undefined
    export let type: $$Props["type"] = "button"
    export let padding: $$Props["padding"] = 2.5
    export let margin: $$Props["margin"] = 0
    let defaultClass = "text-white rounded-lg w-fit"

    export { className as class }
</script>

<button
    on:click
    on:keydown
    {type}
    class={cn([defaultClass, bg ? `variant-${variant ?? "filled"}-${bg}` : "", bg ? hoverClasses[bg] : "", className])}
    {...$$restProps}
>
    <Flex direction="row" gap={2.5} alignItems="center" {padding} {margin}>
        {#if $$slots.leftElement}
            <slot name="leftElement" />
        {/if}
        <slot />
        {#if $$slots.rightElement}
            <slot name="rightElement" />
        {/if}
    </Flex>
</button>
