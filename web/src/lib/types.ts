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

export type ServerResp<T> = {
    message: string
    data: T
}

export type AuthResp = {
    accessToken: string
    refreshToken: string
}
