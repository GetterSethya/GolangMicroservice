<script lang="ts">
    import * as Button from "@ui/button/"
    import { Flex } from "@ui/container"
    import * as Icons from "../svg/index"
    import { MessageSquareReply } from "lucide-svelte"

    export let reply = {
        id: "",
        totalResponse: 0,
    }

    async function handleDeleteReply(id: string) {
        console.log("hit delete", id)
    }

    async function handleEditReply(id: string) {
        console.log("hit edit", id)
    }
</script>

<Flex gap={5} class="divide-y divide-inherit border-surface-700 ">
    <Flex gap={2.5} padding={{ vertical: 2.5 }}>
        <Button.Root class="text-primary-500 rounded-lg hover:bg-surface-500/20 cursor-default" padding={5}>
            <svelte:fragment slot="leftElement">
                <MessageSquareReply size={20} />
            </svelte:fragment>
            <span class="text-primary-200">{reply.totalResponse} Responses</span>
        </Button.Root>
    </Flex>
    <Flex gap={2.5} padding={{ vertical: 2.5 }}>
        <Button.Root
            on:click={async () => {
                await handleEditReply(reply.id)
            }}
            class="fill-warning-500 rounded-lg hover:bg-surface-500/20 w-full"
            padding={5}
        >
            <svelte:fragment slot="leftElement">
                <Icons.PencilSquare />
            </svelte:fragment>
            <span class="text-warning-200">Edit reply</span>
        </Button.Root>
        <Button.Root
            on:click={async () => {
                await handleDeleteReply(reply.id)
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
