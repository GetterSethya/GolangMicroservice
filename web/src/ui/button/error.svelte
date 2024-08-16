<script lang="ts">
    export let onClick: () => void = () => {}
    export let btnType: "button" | "submit" | "reset" = "button"
</script>

<button
    on:click={onClick}
    type={btnType}
    class="transition-all ease-in-out bg-error-500 text-white font-bold h-fit px-5 rounded-lg py-1.5 hover:bg-error-500/90"
>
    <slot />
</button>
