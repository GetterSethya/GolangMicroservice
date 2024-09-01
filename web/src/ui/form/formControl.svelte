<script lang="ts">
    import { ValidationMessage } from "@felte/reporter-svelte"
    import { blur } from "svelte/transition"
    import { twMerge } from "tailwind-merge"

    export let name: string
    export let label: string
    export let labelClass: string | undefined = undefined
    export let errorClass: string | undefined = undefined
    export let formControlClass: string | undefined = undefined
</script>

<ValidationMessage for={name} let:messages>
    <div class={twMerge("flex flex-col gap-1.5", formControlClass)}>
        <label
            for={name}
            class={twMerge(
                "text-2xl transition-all ease-in-out font-bold",
                messages ? "text-error-500" : "",
                labelClass
            )}>{label}</label
        >
        <slot {messages}></slot>
        {#if messages}
            <div transition:blur class="flex flex-col gap-1 text-xs">
                {#each messages as msg}
                    <span transition:blur class={twMerge("text-error-500", errorClass)}>{msg}</span>
                {/each}
            </div>
        {/if}
    </div>
</ValidationMessage>
