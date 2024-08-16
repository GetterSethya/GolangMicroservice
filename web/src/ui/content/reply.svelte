<script lang="ts">
    import * as Container from "@ui/container/"
    import Profile from "./profile.svelte"
    import Username from "./username.svelte"
    import Name from "./name.svelte"
    import Body from "./body.svelte"
    import type { Reply, drawerData } from "@lib/types"
    import { getContext } from "svelte"
    import DrawerMenu from "@lib/components/post/drawerMenu.svelte"
    import * as Button from "@ui/button"
    import { drawerStore,drawerDataStore } from "@lib/store"
    import { ThreeDotVertical } from "@lib/components/svg"
    import { UserRoundPlusIcon } from "lucide-svelte"

    export let data: Reply
    const localUserId = getContext<string>("localUserId")
    const drawerMenu: drawerData<DrawerMenu> = {
        component: DrawerMenu,
        config: { height: "h-1/3" },
    }
</script>

<Container.Flex direction="col" gap={2.5} padding={{ horizontal: 5 }} justify="between">
    <Container.Flex direction="row" justify="between" padding={{ top: 5 }}>
        <Container.Anchor href={`/#/profile/${data.username}`}>
            <Container.Flex direction="row" gap={2.5}>
                <Profile profileUrl={data.profile} />
                <Container.Flex direction="col" gap={1}>
                    <Name name={data.name} />
                    <Username username={data.username} />
                </Container.Flex>
            </Container.Flex>
        </Container.Anchor>
        {#if localUserId === data.idUser}
            <Button.Root
                class="fill-surface-400"
                on:click={() => {
                    $drawerStore = true
                    if (drawerMenu.props) {
                    }
                    $drawerDataStore.drawerData = drawerMenu
                }}
            >
                <svelte:fragment slot="leftElement">
                    <ThreeDotVertical />
                </svelte:fragment>
            </Button.Root>
        {:else}
            <Button.Root variant="ghost" bg="primary" class="text-surface-200" on:click={() => {}}>
                <svelte:fragment slot="leftElement">
                    <UserRoundPlusIcon size={18} />
                </svelte:fragment>
            </Button.Root>
        {/if}
    </Container.Flex>
    <Container.Anchor href={`/#/post/${data.id}`} class="w-full">
        <Container.Flex>
            <Body body={data.body} />
        </Container.Flex>
    </Container.Anchor>
</Container.Flex>
