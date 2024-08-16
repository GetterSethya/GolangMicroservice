<script lang="ts">
    import AppDataProvider from "@lib/components/appDataProvider.svelte"
    import { Toast, initializeStores } from "@skeletonlabs/skeleton"
    import Drawer from "@lib/components/drawer.svelte"
    import Router from "svelte-spa-router"
    import MainRouter from "@routes/router.svelte"

    initializeStores()

    const routes = {
        "*": MainRouter,
    }
</script>

<AppDataProvider>
    <Toast />
    <Drawer />
    <Router {routes} />
</AppDataProvider>
