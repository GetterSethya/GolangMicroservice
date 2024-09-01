<script lang="ts">
    import Router, { type RouteDefinition } from "svelte-spa-router"
    import wrap from "svelte-spa-router/wrap"
    import Layout from "@routes/layout/Baselayout.svelte"

    const prefix = "/app/post"
    const routes: RouteDefinition = {
        "/:id": wrap({
            asyncComponent: () => import("@routes/(protected)/post/page.svelte"),
            conditions:()=>{
                console.log("hit /:id")
                return true
            }
        }),
    }
</script>

<Layout>
    <Router {prefix} {routes} />
</Layout>
