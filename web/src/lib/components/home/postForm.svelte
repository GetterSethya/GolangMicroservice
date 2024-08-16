<script lang="ts">
    import Image from "@lib/components/svg/image.svelte"
    import type { AppData } from "@lib/data"
    import { handleSubmitPost } from "@routes/home/handler"
    import { getContext } from "svelte"
    import { getToastStore } from "@skeletonlabs/skeleton"
    import * as Button from "@ui/button"
    import { Pen } from "lucide-svelte"
    import Box from "@ui/container/box.svelte"
    import Flex from "@ui/container/flex.svelte"
    import X from "@lib/components/svg/x.svelte"

    const toastStore = getToastStore()

    export let isLoading: boolean
    export let onAfterSubmit = () => {}

    const appData = getContext<AppData>("appData")

    let inputElement: HTMLInputElement
    let imagePreviewElement: HTMLImageElement
    let formElement: HTMLFormElement
    let showImage = false

    function handlePreview() {
        if (inputElement?.files?.[0]) {
            const reader = new FileReader()

            reader.onload = (e) => {
                if (imagePreviewElement) {
                    imagePreviewElement.src = e.target?.result as string
                    showImage = true
                }
            }

            reader.readAsDataURL(inputElement.files[0])
        }
    }
</script>

<Flex direction="col" class="min-h-1/4">
    <form
        bind:this={formElement}
        method="post"
        on:submit|preventDefault={async (e) => {
            isLoading = true
            const { success } = await handleSubmitPost(e, appData, toastStore)
            if (success) {
                formElement.reset()
            }
            isLoading = false
            onAfterSubmit()
        }}
        class="flex w-full flex-col p-5 border-b border-surface-700"
    >
        <Box class="relative max-h-[50vh]" margin={{ bottom: 4 }}>
            {#if inputElement?.files?.[0]}
                <Button.Root
                    on:click={() => {
                        inputElement.value = ""
                        inputElement.files = null
                        imagePreviewElement.src = ""
                        imagePreviewElement.removeAttribute("src")
                        showImage = false
                    }}
                    class="absolute right-5 top-5 fill-error-500 rounded-full bg-error-100 flex flex-row items-center justify-center h-8 w-8"
                >
                    <X />
                </Button.Root>
            {/if}

            <img
                alt=""
                src=""
                bind:this={imagePreviewElement}
                class={`w-full h-full object-cover ${showImage ? "block" : "hidden"}`}
            />
        </Box>
        <textarea
            name="reqBody"
            id="postBody"
            class="outline-none bg-transparent border-b border-surface-700 my-2.5"
            rows="3"
            placeholder="Apa yang sedang anda pikirkan"
        ></textarea>
        <Flex direction="row" justify="between" alignItems="center">
            <Button.Root
                variant="ghost"
                bg="surface"
                class="fill-surface-400"
                on:click={() => {
                    inputElement.click()
                }}
            >
                <Image />
            </Button.Root>
            <input bind:this={inputElement} type="file" name="reqImage" hidden on:input={handlePreview} />
            <Button.Root type="submit" bg="primary" variant="ghost" class="text-surface-200">
                <Pen slot="leftElement" size={18} />
                <span class="font-bold text-sm">Post</span>
            </Button.Root>
        </Flex>
    </form>
</Flex>
