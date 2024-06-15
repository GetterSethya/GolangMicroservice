<script lang="ts">
    export let label: string = ""
    export let onClick: () => void = () => {}
    export let type: "wide" | "narrow" = "narrow"
</script>

<div class:min-w-[80px]={type === "wide"} class="fill-primary-500 flex flex-row gap-2.5 items-center">
    <button class="hover:bg-surface-700 p-2.5 rounded-lg" on:click={onClick}>
        <slot />
    </button>
    {#if label !== ""}
        <span class="font-bold text-left">{label}</span>
    {/if}
</div>
