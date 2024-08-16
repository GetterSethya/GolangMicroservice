<script lang="ts">
    import { AppData } from "@lib/data"
    import { FetchError, type Post, type User } from "@lib/types"
    import { ProgressRadial } from "@skeletonlabs/skeleton"
    import { getContext, onMount, setContext } from "svelte"
    import type { Writable } from "svelte/store"
    import { Post as PostComponent } from "@ui/content/"
    import * as Container from "@ui/container"
    import * as Button from "@ui/button"
    import LoadingSkeleton from "@lib/components/home/loadingSkeleton.svelte"
    import { profileStore as profile } from "@lib/store"

    const localUser = getContext<Writable<User | null>>("localUserStore")
    $:if ($localUser) {
        setContext("localUserId", $localUser?.id)
    }

    const appData = getContext<AppData>("appData")
    let fetchPostTrigger = false
    let limit = 10
    let posts: Post[] = []
    let errorState = {
        isError: false,
        message: "",
    }
    let showLoadMoreBtn = true
    let isLoading = false
    let meta = {
        cursor: 0,
    }

    onMount(async () => {
        await fetchPost(0)
    })

    async function fetchPost(cursor: number) {
        try {
            const { res, status } = await appData.listPostByUser(
                fetchPostTrigger,
                cursor,
                limit,
                $profile?.id as string
            )
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
            const { res, status } = await appData.listPostByUser(
                fetchPostTrigger,
                meta.cursor,
                limit,
                $profile?.id as string
            )
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

    $: console.log({ posts }, { errorState })
</script>

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
