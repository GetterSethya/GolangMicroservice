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
    image: string
    body: string
    totalReplies: number
    totalLikes: number
}

export class AuthError extends Error {
    constructor(msg: string) {
        super(msg)
    }
}
