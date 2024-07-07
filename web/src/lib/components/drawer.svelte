<script lang="ts">
    import { fade, slide } from "svelte/transition"
    import { drawerStore, drawerDataStore } from "@lib/store"
    import type { drawerData } from "@lib/types"
    import type { Action } from "svelte/action"

    let drawerData: drawerData[] = []

    drawerDataStore.subscribe((d) => {
        drawerData = d
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
    <div role="none" transition:slide class="fixed bottom-0 left-0 z-[9999] w-screen h-screen flex flex-col">
        <div
            on:outside={closeModal}
            use:clickOutside
            class="bg-surface-900 shadow-lg flex h-1/4 flex-col gap-2.5 p-5 mt-auto border-t border-surface-700"
        >
            {#if drawerData}
                {#each drawerData as data}
                    <div class="flex flex-col">
                        <button
                            on:click={data.onClick}
                            class={`${data.icon.fill} active:bg-surface-700 flex flex-row gap-2.5 border border-surface-700 hover:bg-surface-800 p-2.5 rounded-lg items-center`}
                        >
                            <div class={`p-2.5 ${data.icon.class}`}>
                                <svelte:component this={data.icon.component} />
                            </div>
                            <span class={`font-bold text-lg ${data.labelClass}`}>{data.label}</span>
                        </button>
                    </div>
                {/each}
            {/if}
        </div>
    </div>
    <span transition:fade class="bg-black/20 z-50 w-screen h-screen absolute top-0">
        <br />
    </span>
{/if}
