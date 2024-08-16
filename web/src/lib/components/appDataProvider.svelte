<script lang="ts">
    import { AppData } from "@lib/data"
    import { setContext } from "svelte"

    const apiUrl = import.meta.env.VITE_API_URL as string
    const appData = new AppData(apiUrl)
    setContext<AppData>("appData", appData)
</script>

<slot />
