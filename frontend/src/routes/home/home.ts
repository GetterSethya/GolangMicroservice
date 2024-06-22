import type { AppData } from "@lib/data"
import type { ToastSettings, ToastStore } from "@skeletonlabs/skeleton"

export async function handleSubmitPost(e: SubmitEvent, appData: AppData, toastStore: ToastStore) {
    const fd = new FormData(e.target as HTMLFormElement)

    try {
        await appData.createPost(fd)
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
        toastStore.trigger(t)

        return { success: false }
    }

    return { success: true }
}
