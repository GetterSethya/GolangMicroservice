<script lang="ts">
    import Router from "svelte-spa-router"
    import Home from "@routes/home/home.svelte"
    import Login from "@routes/login/login.svelte"
    import Register from "@routes/register/register.svelte"
    import Profile from "@routes/profile/profile.svelte"
    import { Toast, initializeStores } from "@skeletonlabs/skeleton"

    initializeStores()

    const routes = {
        "/": Home,
        "/login": Login,
        "/register": Register,
        "/profile": Profile,
    }
</script>

<Toast />
<Router {routes} />
