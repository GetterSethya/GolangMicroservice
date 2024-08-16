<script lang="ts">
    import { AppData } from "@lib/data"
    import { getContext, onMount, setContext } from "svelte"
    import * as jose from "jose"
    import { writable, type Writable } from "svelte/store"
    import type { User } from "@lib/types"

    const decodedJwt = jose.decodeJwt(localStorage.getItem("accessToken") ?? "")
    const apiUrl = import.meta.env.VITE_API_URL as string
    const appData = new AppData(apiUrl)

    setContext<Writable<User | null>>("localUserStore", writable(null))
    const userStore = getContext<Writable<User | null>>("localUserStore")

    onMount(async () => {
        const { res } = await appData.getUserById(decodedJwt.sub as string)
        userStore.set(res.data.user)
    })
</script>

<slot></slot>
