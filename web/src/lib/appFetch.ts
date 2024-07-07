import { push } from "svelte-spa-router"
import * as jose from "jose"
import type { AuthResp, ServerResp } from "./types"
import { AuthError } from "@lib/types"

export class JWT {
    private readonly _access: string = "accessToken"
    private readonly _refresh: string = "refreshToken"

    public get access(): string {
        const access = localStorage.getItem(this._access)
        if (!access) throw new AuthError("jwt token missing")

        return access
    }

    public get refresh(): string {
        const refresh = localStorage.getItem(this._refresh)
        if (!refresh) throw new AuthError("jwt token missing")

        return refresh
    }

    public set access(v: string) {
        localStorage.setItem(this._access, v)
    }

    public set refresh(v: string) {
        localStorage.setItem(this._refresh, v)
    }

    public deleteAccessToken() {
        localStorage.removeItem(this._access)
    }

    public deleteRefreshToken() {
        localStorage.removeItem(this._refresh)
    }

    public deleteToken() {
        this.deleteAccessToken()
        this.deleteRefreshToken()
    }

    public validateAccess(): boolean {
        const access = this.access
        if (!access) return false

        const decodedAccess = jose.decodeJwt(access)
        if (!decodedAccess.exp) return false

        return decodedAccess.exp > new Date().getTime() / 1000
    }

    public validateRefresh(): boolean {
        const refresh = this.refresh
        if (!refresh) return false

        const decodedRefresh = jose.decodeJwt(refresh)
        if (!decodedRefresh.exp) return false

        return decodedRefresh.exp > new Date().getTime() / 1000
    }

    public async getNewJWT() {
        const fetchNewJWT = await fetch("http://localhost/v1/auth/refresh", {
            method: "POST",
            body: JSON.stringify({ refreshToken: this.refresh }),
        })

        if (fetchNewJWT.status !== 200) {
            this.deleteToken()
            push("/login")

            throw new AuthError("invalid token")
        }

        //jika token baru berhasil di fetch, set token lagi
        const resp = (await fetchNewJWT.json()) as ServerResp<AuthResp>
        this.access = resp.data.accessToken
        this.refresh = resp.data.refreshToken
    }
}

export async function appFetch(input: URL | RequestInfo, init?: RequestInit | undefined) {
    const jwt = new JWT()

    if (!jwt.validateAccess()) {
        console.log("token tidak valid, mencoba refresh jwt")
        if (!jwt.validateRefresh()) {
            jwt.deleteToken()
            push("/login")

            throw new AuthError("invalid token")
        }

        await jwt.getNewJWT()
        if (!init) {
            init = {}
        }
        if (!init.headers) {
            init.headers = new Headers()
        }

        const newHeader = new Headers(init.headers)
        newHeader.set("Authorization", jwt.access)
        init.headers = newHeader
        console.log("refresh jwt berhasil")

        return await fetch(input, init)
    }
    console.log("token valid")
    const originalReq = await fetch(input, init)

    return originalReq
}
