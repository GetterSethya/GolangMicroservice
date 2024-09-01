import { AppData } from "@lib/data"
import { FetchError } from "@lib/types"
import type { ToastStore } from "@skeletonlabs/skeleton"
import { ZodError, z } from "zod"
const replySchema = z.object({
    id: z.string({ message: "Invalid post id" }),
    body: z
        .string({ message: "Invalid reply" })
        .min(1, { message: "Reply is too short" })
        .max(1025, { message: "Reply is too long" }),
})

export async function handleAddReply(e: SubmitEvent, t: ToastStore, appData: AppData) {
    const formElement = e.target as HTMLFormElement
    const fd = new FormData(formElement)
    let toastMessage = "Unknown error when creating reply"
    let toastBg = "bg-error-500"

    const { postId, body } = Object.fromEntries(fd) as Record<string, string>
    try {
        replySchema.parse({ id: postId, body })

        const { res, status } = await appData.createReply(postId, body)
        if (status !== 201) {
            throw new FetchError(res.message)
        }

        formElement.reset()
        toastMessage = "Reply created"
        toastBg = "bg-success-500"
    } catch (err) {
        console.error(err)
        if (err instanceof ZodError) {
            toastBg = "bg-warning-500"
            const errMessage = err.errors
                .map((e) => {
                    return e.message
                })
                .join("\n")
            toastMessage = errMessage
        }

        if (err instanceof FetchError) {
            toastMessage = err.message
        }
    }

    t.trigger({ message: toastMessage, background: toastBg })
}
