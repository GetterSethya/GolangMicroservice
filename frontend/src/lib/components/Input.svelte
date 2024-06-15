<script lang="ts">
    export let label: string
    export let required: boolean
    export let disabled: boolean
</script>

<div class="flex flex-col gap-2.5">
    <label for={label.toLowerCase().split(" ").join("")} class="text-2xl font-bold" {...$$restProps}>{label}</label>
    <input
        name={label.toLowerCase().split(" ").join("")}
        id={label.toLowerCase().split(" ").join("")}
        class="p-2.5 outline-primary-500 bg-transparent border border-surface-700 rounded-lg"
        {...$$restProps}
        {required}
        {disabled}
    />
</div>
