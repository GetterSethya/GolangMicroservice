<script lang="ts">
    import * as Container from "@ui/container"
    import * as Button from "@ui/button"
    import * as Icons from "@lib/components/svg/"
    import Badge from "@ui/badge/root.svelte"
    import { Name, Profile, Username } from "@ui/content"
    import { UserRoundCog, UserRoundMinusIcon, UserRoundPlusIcon } from "lucide-svelte"
    import { location, pop, replace } from "svelte-spa-router"
    import { getContext, onMount } from "svelte"
    import { type Writable } from "svelte/store"
    import { AuthError, type User } from "@lib/types"
    import type { AppData } from "@lib/data"
    import { getToastStore, type ToastSettings } from "@skeletonlabs/skeleton"
    import {profileStore as profile} from "@lib/store"

    export let prefix: string
    export let params: { username: string | null }

    const toastStore = getToastStore()
    const localUser = getContext<Writable<User | null>>("localUserStore")
    const appData = getContext<AppData>("appData")
    let isFollower = false
    let isFollowing = false

    onMount(async () => {
        let userId = ""
        try {
            const { res, status } = await appData.getUserByUsername(params.username as string)
            if (status === 200) {
                profile.set(res.data.user)
                userId = res.data.user.id
            }
        } catch (err) {
            let t: ToastSettings = {
                message: "Something went wrong",
                background: "bg-error-500",
            }

            if (err instanceof AuthError) {
                t.message = err.message
            }

            toastStore.trigger(t)
        }

        if (userId !== "") {
            const getIsFollower = appData.checkIsFollower(userId)
            const getIsFollowing = appData.checkIsFollowing(userId)
            const [resIsFollower, resIsFollowing] = await Promise.allSettled([getIsFollower, getIsFollowing])

            if (resIsFollower.status === "fulfilled" && resIsFollower.value.status === 200) {
                isFollower = resIsFollower.value.res.data.isFollower
            }

            if (resIsFollowing.status === "fulfilled" && resIsFollowing.value.status === 200) {
                isFollowing = resIsFollowing.value.res.data.isFollowing
            }
        }
    })
</script>

{#if params && params.username && $profile}
    <Container.Box class="w-full border-b border-surface-700" padding={{ vertical: 2.5 }}>
        <Button.Root on:click={() => pop()} class="fill-surface-400 text-surface-400 font-medium">
            <svelte:fragment slot="leftElement">
                <Icons.LeftChevron />
            </svelte:fragment>
            <span>Back</span>
        </Button.Root>
    </Container.Box>
    <Container.Flex
        direction="col"
        justify="between"
        gap={2.5}
        class="p-5 border-b border-surface-700"
        alignItems="center"
    >
        <Container.Flex direction="row" gap={2.5} alignItems="center">
            <Profile width="w-[80px]" height="h-[80px]" profileUrl={$profile.profile} />
            <Container.Flex direction="col" gap={1} class="text-3xl text-surface-200">
                <span class=" line-clamp-1">
                    <Name name={$profile.name} />
                </span>
                <Container.Flex direction="row" alignItems="center" gap={2.5}>
                    <span class="text-lg text-surface-400 truncate no-wrap line-clamp-1">
                        <Username username={$profile.username} />
                    </span>
                    {#if isFollower}
                        <Badge bg="surface" variant="ghost" class="text-surface-200 cursor-default">
                            <span> Following you </span>
                        </Badge>
                    {/if}
                </Container.Flex>
            </Container.Flex>
        </Container.Flex>
        <Container.Flex direction="row" gap={2.5} alignItems="center" justify="between">
            <Container.Flex direction="row" justify="start" class="w-fit" gap={10}>
                <Container.Anchor href={`/#${prefix}/${params.username}/following`}>
                    <Container.Flex direction="row" gap={2.5}>
                        <span>Following</span>
                        <span class="text-surface-200 font-bold">{$profile.totalFollowing}</span>
                    </Container.Flex>
                </Container.Anchor>
                <Container.Anchor href={`/#${prefix}/${params.username}/follower`}>
                    <Container.Flex direction="row" gap={2.5}>
                        <span>Follower</span>
                        <span class="text-surface-200 font-bold">{$profile.totalFollower}</span>
                    </Container.Flex>
                </Container.Anchor>
            </Container.Flex>
            {#if $localUser?.username === params.username}
                <Button.Root variant="ghost" bg="warning" class="text-surface-200" on:click={() => {}}>
                    <svelte:fragment slot="leftElement">
                        <UserRoundCog size={18} />
                    </svelte:fragment>
                    <span>Update</span>
                </Button.Root>
            {:else if isFollowing}
                <Button.Root variant="ghost" bg="primary" class="text-surface-200" on:click={() => {}}>
                    <svelte:fragment slot="leftElement">
                        <UserRoundMinusIcon size={18} />
                    </svelte:fragment>
                    <span>Unfollow</span>
                </Button.Root>
            {:else}
                <!-- sementara -->
                <Button.Root variant="ghost" bg="primary" class="text-surface-200" on:click={() => {}}>
                    <svelte:fragment slot="leftElement">
                        <UserRoundPlusIcon size={18} />
                    </svelte:fragment>
                    <span>Follow</span>
                </Button.Root>
            {/if}
        </Container.Flex>
    </Container.Flex>
    <Container.Flex direction="row" class="divide-x border-surface-700 divide-inherit border-b" alignItems="center">
        <Container.Flex justify="center" alignItems="center" class="w-full text-center h-full hover:bg-surface-800">
            <Button.Root
                class="w-full text-surface-200 font-bold"
                on:click={() => {
                    replace(`/profile/${params.username}/post/`)
                }}
            >
                <span class="mx-auto" class:text-primary-500={$location.startsWith(`/profile/${params.username}/post`)}
                    >Post</span
                >
            </Button.Root>
        </Container.Flex>
        <Container.Flex justify="center" alignItems="center" class="w-full text-center h-full hover:bg-surface-800">
            <Button.Root
                class="w-full text-surface-200 font-bold"
                on:click={() => {
                    replace(`/profile/${params.username}/like/`)
                }}
            >
                <span class="mx-auto" class:text-primary-500={$location.startsWith(`/profile/${params.username}/like`)}
                    >Like</span
                >
            </Button.Root>
        </Container.Flex>
    </Container.Flex>
{/if}
<slot {profile}></slot>
