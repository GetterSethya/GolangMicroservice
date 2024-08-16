<script lang="ts">
    import Router, { type RouteDefinition } from "svelte-spa-router"
    import wrap from "svelte-spa-router/wrap"
    import Layout from "@routes/layout/Baselayout.svelte"

    const prefix = "/post"
    const routes: RouteDefinition = {
        "/:id": wrap({
            asyncComponent: () => import("@routes/post/page.svelte"),
        }),
    }
</script>

<Layout>
    <Router {prefix} {routes} />
</Layout>
