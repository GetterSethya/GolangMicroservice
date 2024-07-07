<script lang="ts">
    import ChatBubble from "@lib/components/svg/chatBubble.svelte"
    import Heart from "@lib/components/svg/heart.svelte"
    import Share from "@lib/components/svg/share.svelte"
    import PostButton from "@lib/components/postButton.svelte"
    import type { Post, drawerData } from "@lib/types"
    import ThreeDotVertical from "@lib/components/svg/threeDotVertical.svelte"
    import { getContext } from "svelte"
    import { drawerStore, drawerDataStore } from "@lib/store"
    import Trash from "@lib/components/svg/trash.svelte"
    import Pencil from "@lib/components/svg/pencilSquare.svelte"

    export const onAfterDelete: () => void = () => {}
    const localUserId = getContext<string>("localUserId")
    const drawerMenu: drawerData[] = [
        {
            icon: { component: Trash, fill: "fill-error-500" },
            onClick: handleDeletePost,
            label: "Delete post",
            labelClass: "text-error-400",
        },
        {
            icon: { component: Pencil, fill: "fill-surface-400" },
            onClick: handleEditPost,
            label: "Edit post",
            labelClass: "text-surface-400",
        },
    ]

    async function handleDeletePost() {
        console.log("delete post")
    }
    async function handleEditPost() {
        console.log("edit post")
    }

    export let post: Post
</script>

<div>
    <div class="w-full h-fit px-5 pt-5 flex flex-col gap-2.5">
        <div class="flex flex-row justify-between">
            <a href={`/#/profile/${post.username}`} class="flex flex-row gap-2.5">
                <img
                    src={"http://localhost/v1/image/thumbnail/" + post.profile}
                    alt=""
                    class="w-[50px] h-[50px] rounded-full object-cover"
                />
                <div class="flex flex-col gap-1">
                    <span class="font-bold"
                        >{post.name.length > 20 ? post.name.substring(0, 17) + "..." : post.name}</span
                    >
                    <span>{post.username}</span>
                </div>
            </a>
            {#if localUserId === post.idUser}
                <button
                    on:click={() => {
                        $drawerStore = true
                        $drawerDataStore = drawerMenu
                    }}
                    type="button"
                    class="transition-all ease-in-out fill-surface-400 hover:bg-surface-800 px-1 rounded-lg"
                >
                    <ThreeDotVertical />
                </button>
            {/if}
        </div>
        <a href={`/#/post/${post.id}`} class="flex flex-col gap-2.5">
            {#if post.image !== ""}
                <img
                    class="rounded-lg max-h-[50vh] w-full object-cover"
                    src={"http://localhost/v1/image/thumbnail/" + post.image}
                    alt=""
                />
            {/if}
            <span class="text-surface-100 pt-2.5">{post.body}</span>
        </a>
    </div>
    <div class="px-5 py-5 flex justify-evenly flex-col gap-2.5 border-surface-700">
        <div class="flex flex-row justify-between">
            <div class="flex flex-row gap-5">
                <PostButton type="wide" label={post.totalReplies.toString()}>
                    <ChatBubble />
                </PostButton>
                <PostButton type="wide" label={post.totalLikes.toString()}>
                    <Heart />
                </PostButton>
            </div>
            <PostButton>
                <Share />
            </PostButton>
        </div>
    </div>
</div>
