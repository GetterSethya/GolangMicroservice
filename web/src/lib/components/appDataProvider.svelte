<script lang="ts">
    import { AppData } from "@lib/data"
    import { setContext } from "svelte"
    import * as jose from "jose"

    const apiUrl = import.meta.env.VITE_API_URL as string
    const appData = new AppData(apiUrl)

    const decodedJwt = jose.decodeJwt(localStorage.getItem("accessToken") as string)

    setContext<AppData>("appData", appData)
    setContext("localUserId", decodedJwt.sub)
</script>

<slot />
