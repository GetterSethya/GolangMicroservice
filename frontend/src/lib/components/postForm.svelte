<script lang="ts">
    import Image from "@lib/components/svg/image.svelte"

    export let isLoading: boolean
    export let handleSubmit: (e: SubmitEvent) => Promise<void>

    let inputElement: HTMLInputElement
    let imagePreviewElement: HTMLImageElement

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
    method="post"
    on:submit|preventDefault={async (e) => {
        isLoading = true
        try {
            await handleSubmit(e)
        } catch (err) {
            console.error(err)
        }
    }}
    class="flex flex-col p-5 border-b border-surface-700"
>
    <div class="relative w-full max-h-[50vh] mb-4">
        <img alt="" bind:this={imagePreviewElement} class="w-full h-full object-cover" />
    </div>
    <textarea
        name="postBody"
        id="postBody"
        class="outline-none bg-transparent border-b border-surface-700 my-2.5"
        rows="4"
        placeholder="Apa yang sedang anda pikirkan"
    ></textarea>
    <div class="flex flex-row justify-between w-full items-center">
        <button
            on:click={() => {
                inputElement.click()
            }}
            class="bg-surface-800 border border-surface-700 p-2.5 rounded-lg fill-surface-400"
        >
            <Image />
        </button>
        <input bind:this={inputElement} type="file" name="inputFile" hidden on:input={handlePreview} />
        <button class="btn variant-filled-primary font-bold">
            <span>Post</span>
        </button>
    </div>
</form>
