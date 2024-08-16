<script lang="ts">
    import type { Post } from "@lib/types"
    import type { AppData } from "@lib/data"
    import { getContext, onMount, setContext } from "svelte"
    import { FetchError } from "@lib/types"
    import PostForm from "@lib/components/home/postForm.svelte"
    import LoadingSkeleton from "@lib/components/home/loadingSkeleton.svelte"
    import * as Button from "@ui/button/"
    import * as Container from "@ui/container/"
    import { Post as PostComponent } from "@ui/content/"
    import { ProgressRadial } from "@skeletonlabs/skeleton"
    import * as jose from "jose"

    const appData = getContext<AppData>("appData")
    const decodedJwt = jose.decodeJwt(localStorage.getItem("accessToken") ?? "")

    setContext("localUserId", decodedJwt.sub)
    let name = "orang"
    let isLoading: boolean = false
    let fetchPostTrigger = false
    let posts: Post[] = []
    let meta = {
        cursor: 0,
    }
    let limit = 10
    let errorState = {
        isError: false,
        message: "",
    }
    let showLoadMoreBtn = true

    onMount(async () => {
        await fetchPost(0)
    })

    async function fetchPost(cursor: number) {
        try {
            const { res, status } = await appData.getAllPost(fetchPostTrigger, cursor, limit)
            const { data, message } = res

            if (status !== 200) {
                throw new FetchError(message)
            }
            posts = data.posts
            meta = data.meta
            errorState.isError = false
        } catch (err) {
            console.error(err)
            errorState.isError = true
            if (err instanceof FetchError) {
                errorState.message = err.message
            } else {
                errorState.message = "something went wrong"
            }
        }
    }

    async function handleLoadMore() {
        try {
            const { res, status } = await appData.getAllPost(fetchPostTrigger, meta.cursor, limit)
            const { data, message } = res

            if (data.posts.length === 0) {
                showLoadMoreBtn = false
            }

            if (status !== 200) {
                throw new FetchError(message)
            }

            posts = [...posts, ...data.posts]
            posts = posts // trigger svelte reactivity

            if (data.meta.cursor !== 0) {
                meta = data.meta
            }

            errorState.isError = false
        } catch (err) {
            console.error(err)
            errorState.isError = true
            if (err instanceof FetchError) {
                errorState.message = err.message
            } else {
                errorState.message = "something went wrong"
            }
        }
    }
</script>

<Container.Flex class="h-[calc(100%-3rem)] md:h-full  mx-auto border-e border-surface-700">
    <PostForm
        {isLoading}
        onAfterSubmit={async () => {
            fetchPost(0)
        }}
    />
    {#if posts.length > 0 && !errorState.isError}
        <Container.Flex class="overflow-y-scroll divide-y divide-inherit border-surface-700">
            {#each posts as post}
                <div class="w-full">
                    <PostComponent data={post} />
                </div>
            {/each}
            {#if showLoadMoreBtn}
                <Container.Flex justify="center" alignItems="center" padding={2.5}>
                    <Button.Root
                        on:click={async () => {
                            isLoading = true
                            await handleLoadMore()
                            isLoading = false
                        }}
                        class="text-primary-500 hover:text-primary-400"
                    >
                        {#if isLoading}
                            <ProgressRadial width="w-4" meter="stroke-primary-500" />
                        {/if}
                        Load more...
                    </Button.Root>
                </Container.Flex>
            {/if}
        </Container.Flex>
    {:else if errorState.isError}
        <Container.Flex justify="center" alignItems="center" gap={2.5} class="w-1/2 m-auto">
            <span class="text-center">{errorState.message}</span>
            <Button.Root
                bg="error"
                on:click={async () => {
                    await fetchPost(meta.cursor)
                }}
            >
                <span>Refresh</span>
            </Button.Root>
        </Container.Flex>
    {:else}
        <LoadingSkeleton />
    {/if}
</Container.Flex>
