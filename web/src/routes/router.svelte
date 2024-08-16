<script lang="ts">
    import type { RouteDefinition } from "svelte-spa-router"
    import NotFound from "@routes/notFound.svelte"
    import AuthRouter from "@routes/auth/router.svelte"
    import PostRouter from "@routes/post/router.svelte"
    import ProfileRouter from "@routes/profile/router.svelte"
    import HomeRouter from "@routes/home/router.svelte"
    import wrap from "svelte-spa-router/wrap"
    import Router, { replace } from "svelte-spa-router"

    export let params: { username: string }

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
        "/auth": wrap({
            component: AuthRouter,
        }),
        "/auth/*": wrap({
            component: AuthRouter,
        }),
        "*": NotFound,
    }

    function handleConditionsFailed() {
        localStorage.removeItem("accessToken")
        localStorage.removeItem("refreshToken")
        replace("/auth/login")
    }
</script>

<Router {routes} on:conditionsFailed={handleConditionsFailed} />
