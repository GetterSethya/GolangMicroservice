<script lang="ts">
    import Router, { type RouteDefinition } from "svelte-spa-router"
    import Layout from "../layout/Baselayout.svelte"
    import wrap from "svelte-spa-router/wrap"

    const prefix = "/home"
    const routes: RouteDefinition = {
        "/": wrap({
            asyncComponent: () => import("@routes/home/page.svelte"),
        }),
    }
</script>

<Layout>
    <Router {prefix} {routes} />
</Layout>
