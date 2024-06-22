import { appFetch } from "./appFetch"
import type { Post, ServerResp } from "./types"

export class AppData {
    private _apiUrl: string
    constructor(apiUrl: string) {
        this._apiUrl = apiUrl
    }

    public get apiUrl(): string {
        return this._apiUrl
    }

    public async getForYou(trigger: boolean, cursor: number, limit: number) {
        throw new Error("method is not implemented yet!")
    }

    public async getAllPost(trigger: boolean, cursor: number, limit: number) {
        const fetchPosts = await appFetch(`${this.apiUrl}/post/?cursor=${cursor}&limit=${limit}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return fetchPosts.json() as Promise<ServerResp<{ posts: Post[]; meta: { cursor: number } }>>
    }

    public async createPost(fd: FormData) {
        const fetchPosts = await appFetch(`${this.apiUrl}/post/`, {
            method: "POST",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
            body: fd,
        })
        return fetchPosts.json() as Promise<ServerResp<null>>
    }

    public async editPost() {
        throw new Error("method is not implemented yet!")
    }

    public async deletePost() {
        throw new Error("method is not implemented yet!")
    }

    public async getProfile() {
        throw new Error("method is not implemented yet!")
    }

    public async editProfile() {
        throw new Error("method is not implemented yet!")
    }

    public async updatePassword() {
        throw new Error("method is not implemented yet!")
    }
}
