<script lang="ts">
    import { AuthRepository } from "@lib/repository/auth"
    import Router, { replace, type RouteDefinition } from "svelte-spa-router"
    import wrap from "svelte-spa-router/wrap"

    AuthRepository.setCtx()

    const prefix = "/auth"
    const routes: RouteDefinition = {
        "/login": wrap({
            asyncComponent: () => import("@routes/auth/login/page.svelte"),
            conditions: () => {
                const token = localStorage.getItem("accessToken")
                if (token && token !== "") {
                    return false
                }

                return true
            },
        }),
        "/register": wrap({
            asyncComponent: () => import("@routes/auth/register/page.svelte"),
            conditions: () => {
                const token = localStorage.getItem("accessToken")
                if (token && token !== "") {
                    return false
                }

                return true
            },
        }),
    }

    function handleConditionsFailed() {
        replace("/home")
    }
</script>

<Router {prefix} {routes} on:conditionsFailed={handleConditionsFailed} />
