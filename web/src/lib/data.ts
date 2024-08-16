import { appFetch } from "./appFetch"
import type { Post, Reply, ServerResp, User } from "./types"

export class AppData {
    private _apiUrl: string
    constructor(apiUrl: string) {
        this._apiUrl = apiUrl
    }

    public get apiUrl(): string {
        return this._apiUrl
    }

    // create like
    public async likePost(postId: string) {
        const fetchLike = await appFetch(`${this.apiUrl}/post/${postId}/like`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchLike.json()) as ServerResp<null>,

            // jika berhasil return 201
            status: fetchLike.status,
        }
    }

    // delete like
    public async deletelikePost(postId: string) {
        const fetchLike = await appFetch(`${this.apiUrl}/post/${postId}/cancel-like`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchLike.json()) as ServerResp<null>,

            // jika berhasil return 200
            status: fetchLike.status,
        }
    }

    // check is following -> http://localhost/v1/relation/{userId}/is-following
    public async checkIsFollowing(userId: string) {
        const fetchIsFollowing = await appFetch(`${this.apiUrl}/relation/${userId}/is-following`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchIsFollowing.json()) as ServerResp<{ isFollowing: boolean }>,
            status: fetchIsFollowing.status,
        }
    }

    public async checkIsFollower(userId: string) {
        const fetchIsFollowing = await appFetch(`${this.apiUrl}/relation/${userId}/is-follower`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchIsFollowing.json()) as ServerResp<{ isFollower: boolean }>,
            status: fetchIsFollowing.status,
        }
    }

    public async getUserByUsername(username: string) {
        const fetchUser = await appFetch(`${this.apiUrl}/user/username/${username}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchUser.json()) as ServerResp<{ user: User }>,
            status: fetchUser.status,
        }
    }

    public async getUserById(userId: string) {
        const fetchUser = await appFetch(`${this.apiUrl}/user/${userId}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchUser.json()) as ServerResp<{ user: User }>,
            status: fetchUser.status,
        }
    }

    public async createReply(postId: string, body: string) {
        const postReply = await appFetch(`${this.apiUrl}/reply/${postId}/create`, {
            method: "POST",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
            body: JSON.stringify({ body }),
        })

        return {
            res: (await postReply.json()) as ServerResp<null>,
            status: postReply.status,
        }
    }

    public async getReply(postId: string) {
        const fetchReply = await appFetch(`${this.apiUrl}/reply/post/${postId}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchReply.json()) as ServerResp<{ reply: Reply[]; meta: { cursor: number } }>,
            status: fetchReply.status,
        }
    }

    public async getPost(postId: string) {
        const fetchPost = await appFetch(`${this.apiUrl}/post/${postId}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchPost.json()) as ServerResp<{ post: Post }>,
            status: fetchPost.status,
        }
    }

    public async getForYou(trigger: boolean, cursor: number, limit: number) {
        throw new Error("method is not implemented yet!")
    }

    public async listPostByUser(trigger: boolean, cursor: number, limit: number, idUser: string) {
        const fetchPosts = await appFetch(`${this.apiUrl}/post/user/${idUser}?cursor=${cursor}&limit=${limit}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchPosts.json()) as ServerResp<{ posts: Post[]; meta: { cursor: number } }>,
            status: fetchPosts.status,
        }
    }

    public async getAllPost(trigger: boolean, cursor: number, limit: number) {
        const fetchPosts = await appFetch(`${this.apiUrl}/post/?cursor=${cursor}&limit=${limit}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        })

        return {
            res: (await fetchPosts.json()) as ServerResp<{ posts: Post[]; meta: { cursor: number } }>,
            status: fetchPosts.status,
        }
    }

    public async createPost(fd: FormData) {
        const fetchPosts = await appFetch(`${this.apiUrl}/post/`, {
            method: "POST",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
            body: fd,
        })

        return {
            res: (await fetchPosts.json()) as ServerResp<null>,
            status: fetchPosts.status,
        }
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
