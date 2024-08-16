import type { twFractions, twSize } from "@ui/types"
import type { ComponentProps, ComponentType, SvelteComponent } from "svelte"

export type drawerData<T extends SvelteComponent = SvelteComponent> = {
    component: ComponentType
    props?: ComponentProps<T>
    config?: {
        height?: `h-${(typeof twSize)[number]}` | `h-${(typeof twFractions)[number]}`
    }
}

export type ServerResp<T> = {
    message: string
    data: T
}

export type AuthResp = {
    accessToken: string
    refreshToken: string
}

export type Post = {
    id: string
    profile: string
    name: string
    username: string
    idUser: string
    image: string
    body: string
    totalReplies: number
    totalLikes: number
    isLiked:boolean

    createdAt: number
    updatedAt: number
}

export type Reply = {
    id: string
    body: string
    idPost: string
    idUser: string
    name: string
    username: string
    profile: string
    totalChild: number
    parentId: string | null

    createdAt: number
    updatedAt: number
}

export type User = {
    id: string
    name: string
    username: string
    profile: string
    createdAt: string
    updatedAt: string
    totalFollower: number
    totalFollowing: number
}

export class AuthError extends Error {
    constructor(msg: string) {
        super(msg)
    }
}

export class FetchError extends Error {
    constructor(msg: string) {
        super(msg)
    }
}
