<script lang="ts">
    import Router, { replace, type RouteDefinition } from "svelte-spa-router"
    import wrap from "svelte-spa-router/wrap"
    import Layout from "@routes/layout/Baselayout.svelte"
    import ProfileLayout from "./layout.svelte"

    let internalParams: { username: string | null } = { username: null }
    const prefix = "/app/profile"
    const routes: RouteDefinition = {
        "/:username": wrap({
            asyncComponent: () => import("@routes/(protected)/profile/page.svelte"),
            conditions: (d) => {
                console.log("/:username")
                if (d.params && d.params.username) {
                    internalParams.username = d.params.username
                }
                replace(`${prefix}/${d.params?.username}/post`)
                return true
            },
        }),
        "/:username/post": wrap({
            asyncComponent: () => import("@routes/(protected)/profile/page.svelte"),
            conditions: (d) => {
                console.log("/:username/post")
                if (d.params && d.params.username) {
                    console.log(d.params.username)
                    internalParams.username = d.params.username
                }
                return true
            },
        }),
        "/:username/like": wrap({
            asyncComponent: () => import("@routes/(protected)/profile/page.svelte"),
            conditions: (d) => {
                if (d.params && d.params.username) {
                    internalParams.username = d.params.username
                }
                return true
            },
        }),
    }
</script>

<Layout>
    <ProfileLayout bind:params={internalParams} {prefix}>
        <Router {prefix} {routes} />
    </ProfileLayout>
</Layout>
