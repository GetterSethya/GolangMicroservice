<script lang="ts">
    import { UserRepository } from "@lib/repository/user"
    import { PostRepository } from "@lib/repository/post"
    import { onDestroy, onMount } from "svelte"

    PostRepository.setCtx()
    const userRepo = UserRepository.setCtx()
    const local = userRepo.setLocalUserCtx(null)
    const abort = new AbortController()
    const signal = abort.signal

    onMount(async () => {
        const { res } = await userRepo.getLocalUserData({ init: { signal } })
        if (res.data?.user) {
            local.set(res.data.user)
        }
    })

    onDestroy(() => abort.abort())
</script>

<slot />
