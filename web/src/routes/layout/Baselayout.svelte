<script lang="ts">
    import Header from "@lib/components/header.svelte"
    import Sidebar from "@lib/components/sidebar.svelte"
    import Flex from "@ui/container/flex.svelte"
    import { location } from "svelte-spa-router"
</script>

<Flex class="h-full text-surface-400 md:flex-row-reverse lg:w-3/4 md:mx-auto">
    <Flex class="h-[calc(100%-7rem)] md:h-full  mx-auto border-e border-surface-700">
        {#if $location.startsWith("/app/home")}
            <Header />
        {/if}
        <slot />
    </Flex>
    <Sidebar />
</Flex>
