<script lang="ts">
    export let href: string | null | undefined
    export let label: string
    export let isActive: boolean = false
</script>

<a
    {href}
    {...$$restProps}
    class:bg-surface-800={isActive}
    class:text-primary-500={isActive}
    class:fill-primary-500={isActive}
    class="flex
    flex-row
    gap-2.5
    hover:outline
    outline-[1px]
    outline-surface-700
    px-5
    py-2.5 hover:text-primary-500 hover:fill-primary-500 w-fit rounded-lg items-center hover:bg-surface-800"
>
    <slot />
    <span class="font-bold text-lg hidden md:block"> {label} </span>
</a>
