<script lang="ts">
    import { cn } from "@lib/utils"
    import type { anchorProps as Props } from "./index"

    type $$Props = Props
    type $$Events = {
        click: (
            a: MouseEvent & {
                currentTarget: EventTarget & HTMLAnchorElement
            }
        ) => {}
    }

    export let href: $$Props["href"] = undefined
    export let className: $$Props["class"] = undefined
    export { className as class }
</script>

<a {href} on:click {...$$restProps} class={cn([className])}>
    <slot />
</a>
