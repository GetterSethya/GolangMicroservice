<script lang="ts">
    import Sidebar from "@lib/components/sidebar.svelte"
    import Header from "@lib/components/header.svelte"
    import { getContext, onMount } from "svelte"
    import { push } from "svelte-spa-router"
    import PostForm from "@lib/components/postForm.svelte"
    import type { Post } from "@lib/types"
    import PostCard from "@lib/components/postCard.svelte"
    import PostCardSkeleton from "@lib/components/postCardSkeleton.svelte"
    import ArrowClockwise from "@lib/components/svg/arrowClockwise.svelte"
    import type { AppData } from "@lib/data"
    import { ProgressRadial } from "@skeletonlabs/skeleton"
    import { FetchError } from "@lib/types"

    const appData = getContext<AppData>("appData")
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
        const access = localStorage.getItem("accessToken")
        if (!access) {
            push("/login")
        }
        await fetchPost(0)
    })

    async function fetchPost(cursor: number) {
        try {
            const { res, status } = await appData.getAllPost(fetchPostTrigger, cursor, limit)
            const { data, message } = await res

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
            const { data, message } = await res

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

<div class="flex w-full h-full text-surface-400 md:flex-row-reverse flex-col w-full md:w-3/4 md:mx-auto h-full">
    <div class="w-full h-[calc(100%-7rem)] md:h-full flex flex-col md:flex-col mx-auto border-e border-surface-700">
        <div class="flex flex-col w-full min-h-1/4">
            <Header />
            <PostForm
                {isLoading}
                onAfterSubmit={async () => {
                    fetchPost(0)
                }}
            />
        </div>
        {#if posts.length > 0 && !errorState.isError}
            <div class="w-full flex overflow-y-scroll flex-col divide-y divide-inherit border-surface-700">
                {#each posts as post}
                    <PostCard {post} />
                {/each}
                {#if showLoadMoreBtn}
                    <div class="p-2.5 flex flex-col border-t border-surface-700">
                        <button
                            on:click={async () => {
                                isLoading = true
                                await handleLoadMore()
                                isLoading = false
                            }}
                            class="text-primary-500 flex flex-row gap-2.5 justify-center items-center"
                        >
                            {#if isLoading}
                                <ProgressRadial
                                    width="w-4"
                                    heigth="h-4"
                                    meter="stroke-primary-500"
                                    track="stroke-primary-500/30"
                                />
                            {/if}
                            <span>Load more...</span>
                        </button>
                    </div>
                {/if}
            </div>
        {:else if errorState.isError}
            <div class="w-1/2 justify-center items-center gap-2.5 flex flex-col m-auto">
                <span class="text-center">{errorState.message}</span>
                <button
                    on:click={async () => {
                        await fetchPost(meta.cursor)
                    }}
                    class="btn fill-white variant-filled-error w-fit m-auto"
                >
                    <ArrowClockwise />
                    <span>refresh</span>
                </button>
            </div>
        {:else}
            <div class="w-full flex flex-col divide-y divide-inherit border-surface-700">
                <PostCardSkeleton />
                <PostCardSkeleton />
                <PostCardSkeleton />
            </div>
        {/if}
    </div>
    <Sidebar />
</div>
