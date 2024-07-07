import type { ComponentType } from "svelte"

export type drawerData = {
    icon: {
        component: ComponentType
        class?: string
        fill: string
        fillClass?: string
    }
    label: string
    labelClass?: string
    onClick: () => void
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
}

export type User = {
    id: string
    name: string
    username: string
    createdAt: string
    updatedAt: string
    totalFollower: number
    totlaFollowing: number
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
