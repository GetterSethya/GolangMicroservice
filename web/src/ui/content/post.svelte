<script lang="ts">
    import * as Container from "@ui/container/"
    import * as Button from "../button/"
    import Profile from "./profile.svelte"
    import Username from "./username.svelte"
    import Name from "./name.svelte"
    import Body from "./body.svelte"
    import Image from "./image.svelte"
    import HeartFill from "@lib/components/svg/heart-fill.svelte"
    import { drawerStore, drawerDataStore, postDetailStore } from "@lib/store"
    import DrawerMenu from "@lib/components/home/drawerMenu.svelte"
    import type { Post, drawerData } from "@lib/types"
    import { ChatBubble, Heart, Share, ThreeDotVertical } from "@lib/components/svg"
    import { UserRoundPlusIcon } from "lucide-svelte"
    import { getContext } from "svelte"
    import { push } from "svelte-spa-router"
    import type { AppData } from "@lib/data"
    import { UserRepository } from "@lib/repository/user"

    export let data: Post
    const appData = getContext<AppData>("appData")
    const userRepo = UserRepository.getCtx()
    const localUser = userRepo.getLocalUserCtx()
    const drawerMenu: drawerData<DrawerMenu> = {
        component: DrawerMenu,
        props: { post: { id: "", totalLike: 0, totalReply: 0 } },
        config: { height: "h-1/3" },
    }

    async function toggleLike() {
        try {
            if (data.isLiked === true) {
                const fetchCancelLike = await appData.deletelikePost(data.id)
                if (fetchCancelLike.status === 200) {
                    data.isLiked = false
                    data.totalLikes--
                } else {
                    data.isLiked = true
                }
            } else {
                const fetchLike = await appData.likePost(data.id)
                if (fetchLike.status === 201) {
                    data.isLiked = true
                    data.totalLikes++
                } else {
                    data.isLiked = true
                }
            }
        } catch (error) {
            console.error(error)
        }
    }
</script>

{#if !$localUser}
    <Container.Box class="w-full h-[35vh] p-5">
        <Container.Box class="bg-surface-800 h-full rounded-lg transition-all ease-in-out animate-pulse">
            <br />
        </Container.Box>
    </Container.Box>
{:else}
    <Container.Flex direction="col" gap={2.5} padding={{ horizontal: 5 }} justify="between">
        <Container.Flex direction="row" justify="between" padding={{ top: 5 }}>
            <Container.Anchor href={`/#/app/profile/${data.username}`}>
                <Container.Flex direction="row" gap={2.5}>
                    <Profile profileUrl={data.profile} />
                    <Container.Flex direction="col" gap={1}>
                        <Name name={data.name} />
                        <Username username={data.username} />
                    </Container.Flex>
                </Container.Flex>
            </Container.Anchor>
            {#if $localUser.id === data.idUser}
                <Button.Root
                    class="fill-surface-400"
                    on:click={() => {
                        $drawerStore = true
                        if (drawerMenu.props) {
                            drawerMenu.props.post = {
                                id: data.id,
                                totalReply: data.totalReplies,
                                totalLike: data.totalLikes,
                            }
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
        <Container.Anchor
            on:click={() => {
                postDetailStore.set(data)
                push(`/app/post/${data.id}`)
            }}
            href={undefined}
            class="w-full cursor-pointer"
        >
            <Container.Flex>
                {#if data.image !== ""}
                    <Image imageUrl={data.image} />
                {/if}
                <Body body={data.body} />
            </Container.Flex>
        </Container.Anchor>
    </Container.Flex>
    <Container.Flex direction="row" justify="between" alignItems={"center"} padding={{ horizontal: 5 }} gap={2.5}>
        <Container.Flex direction="row" gap={2.5} padding={{ vertical: 2.5 }}>
            <Button.Root padding={2.5} class="hover:bg-surface-500/20 fill-primary-500 text-surface-400 text-sm">
                <svelte:fragment slot="leftElement">
                    <ChatBubble />
                </svelte:fragment>
                <span>{data.totalReplies}</span>
            </Button.Root>
            <Button.Root
                on:click={toggleLike}
                padding={2.5}
                class="hover:bg-surface-500/20 fill-primary-500 text-surface-400 text-sm"
            >
                <svelte:fragment slot="leftElement">
                    {#if data.isLiked === true}
                        <HeartFill />
                    {:else}
                        <Heart />
                    {/if}
                </svelte:fragment>
                <span>{data.totalLikes}</span>
            </Button.Root>
        </Container.Flex>
        <Button.Root padding={2.5} class="hover:bg-surface-500/20 fill-primary-500 text-surface-400 text-sm">
            <svelte:fragment slot="leftElement">
                <Share />
            </svelte:fragment>
        </Button.Root>
    </Container.Flex>
{/if}
