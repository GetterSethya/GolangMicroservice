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

    const appData = getContext<AppData>("appData")
    let isLoading: boolean = false
    let fetchPostTrigger = false
    let posts: Post[] = []
    let meta = {
        cursor: 0,
    }
    let limit = 2
    let isError = false

    onMount(async () => {
        const access = localStorage.getItem("accessToken")
        if (!access) {
            push("/login")
        }
        await fetchPost()
    })

    async function fetchPost() {
        try {
            const { data } = await appData.getAllPost(fetchPostTrigger, meta.cursor, limit)
            posts = data.posts
            meta = data.meta
            isError = false
        } catch (err) {
            console.error(err)
            isError = true
        }
    }
</script>

<div
    class="w-full h-full text-surface-400 justify-between flex flex-col-reverse md:flex-row w-2/3 mx-auto border-e border-surface-700"
>
    <Sidebar />
    <div class="flex flex-col w-full h-full">
        <div class="flex flex-col">
            <Header />
            <PostForm
                {isLoading}
                callBack={async () => {
                    await fetchPost()
                }}
            />
        </div>
        <div class="h-[62vh] md:h-full flex flex-col overflow-y-scroll">
            {#if posts.length > 0 && !isError}
                <div class="w-full flex flex-col divide-y divide-inherit border-surface-700">
                    {#each posts as post}
                        <PostCard {post} />
                    {/each}
                </div>
            {:else if isError}
                <div class="w-1/2 justify-center items-center gap-2.5 flex flex-col m-auto">
                    <span class="text-center">Something went wrong, please try again later</span>
                    <button
                        on:click={async () => {
                            await fetchPost()
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
            <div class="p-2.5 flex flex-col border-t border-surface-700">
                <button
                    on:click={async () => {
                        isLoading = true
                        try {
                            const { data } = await appData.getAllPost(fetchPostTrigger, meta.cursor, limit)
                            posts = [...posts, ...data.posts]
                            posts = posts // trigger svelte reactivity
                            if (data.meta.cursor !== 0) {
                                meta = data.meta
                            }
                            isError = false
                        } catch (err) {
                            console.error(err)
                            isError = true
                        }

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
        </div>
    </div>
</div>
