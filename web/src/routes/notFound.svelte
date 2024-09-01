<script lang="ts">
    import {  replace } from "svelte-spa-router"

    const token = localStorage.getItem("accessToken")

    if (token) {
        replace("/app/home")
    } else {
        replace("/auth/login")
    }
</script>
