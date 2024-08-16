import type { AppData } from "@lib/data"
import { FetchError } from "@lib/types"
import type { ToastSettings, ToastStore } from "@skeletonlabs/skeleton"

export async function handleSubmitPost(e: SubmitEvent, appData: AppData, toastStore: ToastStore) {
    const fd = new FormData(e.target as HTMLFormElement)

    try {
        const { res, status } = await appData.createPost(fd)
        const { message } = await res

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
