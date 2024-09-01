<script lang="ts">
    import type { RouteDefinition } from "svelte-spa-router"
    import NotFound from "@routes/notFound.svelte"
    import AuthRouter from "@routes/auth/router.svelte"
    import wrap from "svelte-spa-router/wrap"
    import Router, { replace } from "svelte-spa-router"
    import AppRouter from "@routes/(protected)/router.svelte"
    import { JWT } from "@lib/appFetch"

    export let params: { username: string }

    const routes: RouteDefinition = {
        "/auth": wrap({
            component: AuthRouter,
        }),
        "/auth/*": wrap({
            component: AuthRouter,
        }),
        "/app/*": wrap({
            component: AppRouter,
            conditions: async () => {
                const jwt = new JWT()

                try {
                    if (!jwt.validateAccess()) {
                        if (!jwt.validateRefresh()) {
                            jwt.deleteToken()
                            return false
                        }

                        await jwt.getNewJWT()
                        return true
                    }
                } catch (error) {
                    jwt.deleteToken()
                    return false
                }

                return true
            },
        }),
        "*": NotFound,
    }

    function handleConditionsFailed() {
        replace("/auth/login")
    }
</script>

<Router {routes} on:conditionsFailed={handleConditionsFailed} />
