<script lang="ts">
    import { handleSubmitPost } from "@routes/home/home"
    import Sidebar from "@lib/components/sidebar.svelte"
    import Header from "@lib/components/header.svelte"
    import { onMount } from "svelte"
    import { push } from "svelte-spa-router"
    import PostForm from "@lib/components/postForm.svelte"
    import { appFetch } from "@lib/appFetch"
    import type { ServerResp, Post } from "@lib/types"
    import PostCard from "@lib/components/postCard.svelte"
    import PostCardSkeleton from "@lib/components/postCardSkeleton.svelte"
    import ArrowClockwise from "@lib/components/svg/arrowClockwise.svelte"

    let isLoading: boolean = false
    let fetchPostTrigger = false

    onMount(async () => {
        const access = localStorage.getItem("accessToken")
        if (!access) {
            push("/login")
        }
    })

    async function getPosts(trigger?: boolean) {
        try {
            const fetchPosts = await appFetch("http://localhost/v1/post/", {
                method: "GET",
                headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
            })

            await new Promise((res) => {
                setTimeout(() => {
                    res("done")
                }, 1000)
            })

            return fetchPosts.json() as Promise<ServerResp<{ posts: Post[] }>>
        } catch (err) {
            console.error(err)
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
            <PostForm {isLoading} handleSubmit={handleSubmitPost} />
        </div>
        <div class="h-[62vh] md:h-full flex overflow-y-scroll">
            {#await getPosts(fetchPostTrigger)}
                <div class="w-full flex flex-col divide-y divide-inherit border-surface-700">
                    <PostCardSkeleton />
                    <PostCardSkeleton />
                    <PostCardSkeleton />
                </div>
            {:then res}
                {#if res && res.data}
                    <div class="w-full flex flex-col divide-y divide-inherit border-surface-700">
                        {#each res.data.posts as post}
                            <PostCard {post} />
                        {/each}
                    </div>
                {/if}
            {:catch _}
                <div class="w-1/2 justify-center items-center gap-2.5 flex flex-col m-auto">
                    <span class="text-center">Something went wrong, please try again later</span>
                    <button
                        on:click={() => {
                            fetchPostTrigger = !fetchPostTrigger
                        }}
                        class="btn fill-white variant-filled-error w-fit m-auto"
                    >
                        <ArrowClockwise />
                        <span>refresh</span>
                    </button>
                </div>
            {/await}
        </div>
    </div>
</div>
