import type { PostRepository } from "@lib/repository/post"
import { FetchError } from "@lib/types"
import type { ToastSettings, ToastStore } from "@skeletonlabs/skeleton"

export async function handleSubmitPost(e: SubmitEvent, postRepo: PostRepository, toastStore: ToastStore) {
    const form = e.target as HTMLFormElement
    const fd = new FormData(form)

    try {
        const { res, status } = await postRepo.createPost({ fd })
        const { message } = res

        if (status !== 201) {
            throw new FetchError(message)
        }
        const t: ToastSettings = {
            message: "success creating new posts",
            background: "bg-success-500",
            timeout: 3500,
        }
        toastStore.trigger(t)
    } catch (err) {
        console.error(err)
        const t: ToastSettings = {
            message: "error creating post",
            background: "bg-error-500",
            timeout: 3500,
        }
        if (err instanceof FetchError) {
            t.message = err.message
        }

        toastStore.trigger(t)

        return { success: false }
    }

    return { success: true }
}
