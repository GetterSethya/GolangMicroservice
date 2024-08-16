<script lang="ts">
    import { drawerStore, drawerDataStore } from "@lib/store"
    import type { drawerData } from "@lib/types"
    import type { Action } from "svelte/action"
    import * as Container from "@ui/container/"
    import * as Transition from "@ui/transition/"
    import { cn } from "@lib/utils"

    let drawerData: drawerData | undefined = undefined

    drawerDataStore.subscribe((d) => {
        drawerData = d.drawerData
        drawerData = drawerData
    })

    type Attributes = {
        "on:outside"?: (event: CustomEvent) => void
    }

    type clickOutsideAction = Action<HTMLElement, any, Attributes>

    function closeModal() {
        drawerStore.set(false)
    }

    const clickOutside: clickOutsideAction = (element) => {
        function handleClick(event: MouseEvent) {
            const targetEl = event.target as HTMLElement

            if (element && !element.contains(targetEl)) {
                const clickOutsideEvent = new CustomEvent("outside")
                element.dispatchEvent(clickOutsideEvent)
            }
        }

        document.addEventListener("click", handleClick, true)

        return {
            destroy() {
                document.removeEventListener("click", handleClick, true)
            },
        }
    }
</script>

{#if $drawerStore}
    <Transition.Slide class="fixed bottom-0 left-0 z-[9999] w-screen h-screen">
        <Container.Flex class="h-full">
            <div
                on:outside={closeModal}
                use:clickOutside
                class={cn([
                    "bg-surface-900 w-full shadow-lg flex flex-col gap-2.5 p-5 mt-auto border-t border-surface-700",
                    drawerData?.config?.height ?? "h-1/4",
                ])}
            >
                {#if drawerData}
                    <svelte:component this={drawerData.component} {...drawerData.props} />
                {/if}
            </div>
        </Container.Flex>
    </Transition.Slide>
    <Transition.Fade class="bg-black/20 z-50 w-screen h-screen absolute top-0">
        <br />
    </Transition.Fade>
{/if}
