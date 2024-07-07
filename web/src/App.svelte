<script lang="ts">
    import Router from "svelte-spa-router"
    import Home from "@routes/home/home.svelte"
    import Login from "@routes/login/login.svelte"
    import Register from "@routes/register/register.svelte"
    import Profile from "@routes/profile/profile.svelte"
    import AppDataProvider from "@lib/components/appDataProvider.svelte"
    import { Toast, initializeStores } from "@skeletonlabs/skeleton"
    import Drawer from "@lib/components/drawer.svelte"

    initializeStores()

    const routes = {
        "/": Home,
        "/login": Login,
        "/register": Register,
        "/profile": Profile,
    }
</script>

<AppDataProvider>
    <Toast />
    <Drawer />
    <Router {routes} />
</AppDataProvider>
