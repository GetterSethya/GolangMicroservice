<script lang="ts">
    import * as Button from "@ui/button/"
    import { Flex } from "@ui/container"
    import ChatBubble from "../svg/chatBubble.svelte"
    import * as Icons from "../svg/index"

    export let post = {
        id: "",
        totalReply: 0,
        totalLike: 0,
    }

    async function handleDeletePost(id: string) {
        console.log("hit delete", id)
    }

    async function handleEditPost(id: string) {
        console.log("hit edit", id)
    }
</script>

<Flex gap={5} class="divide-y divide-inherit border-surface-700 ">
    <Flex gap={2.5} direction="row">
        <Button.Root class="fill-primary-500 rounded-lg hover:bg-surface-500/20 cursor-default" padding={5}>
            <svelte:fragment slot="leftElement">
                <ChatBubble />
            </svelte:fragment>
            <span class="text-primary-200">{post.totalReply} Reply</span>
        </Button.Root>
        <Button.Root class="fill-primary-500 rounded-lg hover:bg-surface-500/20 cursor-default" padding={5}>
            <svelte:fragment slot="leftElement">
                <Icons.Heart />
            </svelte:fragment>
            <span class="text-primary-200">{post.totalLike} Love</span>
        </Button.Root>
    </Flex>
    <Flex gap={2.5} padding={{ vertical: 2.5 }}>
        <Button.Root
            on:click={async () => {
                await handleEditPost(post.id)
            }}
            class="fill-warning-500 rounded-lg hover:bg-surface-500/20 w-full"
            padding={5}
        >
            <svelte:fragment slot="leftElement">
                <Icons.PencilSquare />
            </svelte:fragment>
            <span class="text-warning-200">Edit post</span>
        </Button.Root>
        <Button.Root
            on:click={async () => {
                await handleDeletePost(post.id)
            }}
            class="fill-error-500 rounded-lg hover:bg-surface-500/20 w-full"
            padding={5}
        >
            <svelte:fragment slot="leftElement">
                <Icons.Trash />
            </svelte:fragment>
            <span class="text-error-200">Delete post</span>
        </Button.Root>
    </Flex>
</Flex>
