<script lang="ts">
    import Router, { type RouteDefinition } from "svelte-spa-router"
    import wrap from "svelte-spa-router/wrap"
    import PostRouter from "@routes/(protected)/post/router.svelte"
    import ProfileRouter from "@routes/(protected)/profile/router.svelte"
    import HomeRouter from "@routes/(protected)/home/router.svelte"
    import RepositoryProvider from "@lib/components/repositoryProvider.svelte"

    const prefix = "/app"
    const routes: RouteDefinition = {
        "/home": wrap({
            component: HomeRouter,
            conditions: () => {
                const access = localStorage.getItem("accessToken")
                if (!access) {
                    return false
                }

                return true
            },
        }),
        "/home/*": wrap({
            component: HomeRouter,
            conditions: () => {
                const access = localStorage.getItem("accessToken")
                if (!access) {
                    return false
                }

                return true
            },
        }),
        // /post
        "/post": wrap({
            component: PostRouter,
            conditions: () => {
                console.log("hit /post")
                const access = localStorage.getItem("accessToken")
                if (!access) {
                    return false
                }

                return true
            },
        }),
        "/post/*": wrap({
            component: PostRouter,
            conditions: () => {
                console.log("hit /post/*")
                const access = localStorage.getItem("accessToken")
                if (!access) {
                    return false
                }
                return true
            },
        }),

        // /profile
        "/profile": wrap({
            component: ProfileRouter,
            conditions: () => {
                const access = localStorage.getItem("accessToken")
                if (!access) {
                    return false
                }

                return true
            },
        }),
        "/profile/*": wrap({
            component: ProfileRouter,
            conditions: () => {
                const access = localStorage.getItem("accessToken")
                if (!access) {
                    return false
                }

                return true
            },
        }),
    }
</script>

<RepositoryProvider>
    <Router {prefix} {routes} />
</RepositoryProvider>
