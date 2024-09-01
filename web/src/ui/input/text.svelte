<script lang="ts">
    import FormControl from "@ui/form/formControl.svelte"
    import { twMerge } from "tailwind-merge"

    export let label: string
    export let required: boolean
    export let disabled: boolean
    export let name: string
    export let className: string | undefined = undefined
</script>

<FormControl {label} {name} let:messages>
    <input
        {name}
        type="text"
        class={twMerge(
            `p-2.5 outline-primary-500 bg-transparent border ${messages ? "border-error-500" : "border-surface-700"} rounded-lg`,
            className
        )}
        {...$$restProps}
        {required}
        {disabled}
    />
</FormControl>
