<script lang="ts">
    import { AppData } from "@lib/data"
    import type { Post, Reply } from "@lib/types"
    import { getContext, onDestroy, onMount, setContext } from "svelte"
    import { isLoading } from "@lib/store"
    import { pop } from "svelte-spa-router"
    import { Post as PostComponent, Reply as ReplyComponent } from "@ui/content/"
    import { ReplyIcon } from "lucide-svelte"
    import { postDetailStore } from "@lib/store"
    import * as Button from "@ui/button/"
    import * as Icons from "@lib/components/svg"
    import * as Ui from "@lib/components/post/"
    import * as Container from "@ui/container"
    import * as jose from "jose"
    import { handleAddReply } from "./handler"
    import { getToastStore } from "@skeletonlabs/skeleton"

    export let params: { id: string }

    const appData = getContext<AppData>("appData")
    const decodedJwt = jose.decodeJwt(localStorage.getItem("accessToken") ?? "")
    const toastStore = getToastStore()

    setContext("localUserId", decodedJwt.sub)

    let post: Post | undefined = $postDetailStore
    let replies: Reply[] = []
    let cursor = 0

    async function fetchReplies() {
        try {
            const { status, res } = await appData.getReply(params.id as string)
            if (status === 200) {
                const serverResp = res
                return {
                    data: serverResp.data.reply,
                    meta: serverResp.data.meta,
                }
            }
        } catch (err) {
            console.error(err)
        }

        return { data: [], meta: { cursor: 0 } }
    }

    async function fetchPost() {
        try {
            const { status, res } = await appData.getPost(params.id as string)
            if (status === 200) {
                const serverResp = res
                return {
                    data: serverResp.data.post,
                }
            }
        } catch (err) {
            console.error(err)
        }
        return {
            data: undefined,
        }
    }

    onMount(async () => {
        isLoading.set(true)

        const resReplies = await fetchReplies()
        replies = resReplies.data
        cursor = resReplies.meta.cursor

        if (!post) {
            const resPost = await fetchPost()
            post = resPost.data
        }

        isLoading.set(false)
    })

    onDestroy(() => isLoading.set(false))
</script>

<Container.Box class="w-full border-b border-surface-700" padding={{ vertical: 2.5 }}>
    <Button.Root on:click={() => pop()} class="fill-surface-400 text-surface-400 font-medium">
        <svelte:fragment slot="leftElement">
            <Icons.LeftChevron />
        </svelte:fragment>
        <span>Back</span>
    </Button.Root>
</Container.Box>
{#if post}
    <PostComponent data={post} />
{/if}
<form class="w-full flex flex-row items-center justify-between" on:submit|preventDefault={() => {}}>
    <Container.Flex direction="row" gap={2.5} class="h-fit border-y border-surface-700" padding={2.5}>
        <form
            style="display:contents;"
            on:submit|preventDefault={async (e) => {
                await handleAddReply(e, toastStore, appData)
            }}
        >
            <input type="hidden" name="postId" value={params.id} />
            <input
                type="text"
                placeholder="Add your reply"
                name="body"
                required={true}
                class="outline-none p-2.5 bg-transparent w-full"
            />
            <Button.Root bg="primary" type="submit" variant="ghost" class="text-surface-200">
                <ReplyIcon size={20} />
            </Button.Root>
        </form>
    </Container.Flex>
</form>
{#if !$isLoading}
    <Container.Flex gap={5} class=" border-surface-700 divide-y divide-inherit overflow-y-auto">
        {#each replies as reply}
            <ReplyComponent data={reply} />
        {/each}
    </Container.Flex>
{:else}
    <Ui.LoadingSkeleton />
{/if}
