<script lang="ts">
    import { type RouteDefinition } from "svelte-spa-router"
    import wrap from "svelte-spa-router/wrap"
    import ProfilePostPage from "./post/page.svelte"
    import ProfileLikePage from "./like/page.svelte"
    import Router from "svelte-spa-router"

    export let params: { username: string }

    const prefix = "/profile/" + params.username
    const routes: RouteDefinition = {
        "/post": wrap({
            component: ProfilePostPage,
        }),
        "/like": wrap({
            component: ProfileLikePage,
        }),
    }
</script>

<Router {prefix} {routes} />
