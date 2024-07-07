<script lang="ts">
    import Image from "@lib/components/svg/image.svelte"
    import type { AppData } from "@lib/data"
    import { handleSubmitPost } from "@routes/home/home"
    import { getContext } from "svelte"
    import { getToastStore } from "@skeletonlabs/skeleton"

    const toastStore = getToastStore()

    export let isLoading: boolean
    export let onAfterSubmit = () => {}

    const appData = getContext<AppData>("appData")

    let inputElement: HTMLInputElement
    let imagePreviewElement: HTMLImageElement
    let formElement: HTMLFormElement

    function handlePreview() {
        if (inputElement?.files?.[0]) {
            const reader = new FileReader()

            reader.onload = (e) => {
                if (imagePreviewElement) {
                    imagePreviewElement.src = e.target?.result as string
                }
            }

            reader.readAsDataURL(inputElement.files[0])
        }
    }
</script>

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
    class="flex flex-col p-5 border-b border-surface-700"
>
    <div class="relative w-full max-h-[50vh] mb-4">
        <img alt="" bind:this={imagePreviewElement} class="w-full h-full object-cover" />
    </div>
    <textarea
        name="reqBody"
        id="postBody"
        class="outline-none bg-transparent border-b border-surface-700 my-2.5"
        rows="3"
        placeholder="Apa yang sedang anda pikirkan"
    ></textarea>
    <div class="flex flex-row justify-between w-full items-center">
        <button
            type="button"
            on:click|stopPropagation={() => {
                inputElement.click()
            }}
            class="bg-surface-800 border border-surface-700 p-2.5 rounded-lg fill-surface-400"
        >
            <Image />
        </button>
        <input bind:this={inputElement} type="file" name="reqImage" hidden on:input={handlePreview} />
        <button class="btn variant-filled-primary font-bold">
            <span>Post</span>
        </button>
    </div>
</form>
