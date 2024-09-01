import type { ServerResp } from "@lib/types"
import { z } from "zod"

export const BaseSchema = z.object({
    id: z.string(),
    createdAt: z.number(),
    updatedAt: z.number(),
})

export type Base = z.infer<typeof BaseSchema>
export type Init = Omit<RequestInit, "body" | "method" | "headers">

export type APIResponse<T> = {
    res: ServerResp<T>
    status: number
}

export const APIResponseSchema = z.object({
    message: z.string(),
    data: z.unknown(),
})
