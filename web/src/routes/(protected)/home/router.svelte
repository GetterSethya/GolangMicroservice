<script lang="ts">
    import Router, { type RouteDefinition } from "svelte-spa-router"
    import Layout from "../../layout/Baselayout.svelte"
    import wrap from "svelte-spa-router/wrap"

    const prefix = "/app/home"
    const routes: RouteDefinition = {
        "/": wrap({
            asyncComponent: () => import("@routes/(protected)/home/page.svelte"),
        }),
    }
</script>

<Layout>
    <Router {prefix} {routes} />
</Layout>
