import { z } from "zod"
import { APIResponseSchema, BaseSchema, type APIResponse, type Base } from "./base"
import { appFetch } from "@lib/appFetch"
import { API_URL } from "@lib/constant"
import { getContext, setContext } from "svelte"

export const PostSchema = z.object({
    profile: z.string().min(1, { message: "Invalid profile" }),
    name: z
        .string({ message: "Name is not a string" })
        .min(1, { message: "Name is too short, minimum name is 1 character long" })
        .max(50, { message: "Name is too long, maximum name is 50 character long" }),

    username: z
        .string({ message: "Username is not a string" })
        .min(6, { message: "Username is too short, minimum username is 6 character long" })
        .max(50, { message: "Username is too long, maximum username is 50 character long" })
        .regex(/^[a-zA-Z0-9]+$/, { message: "username cannot contain special character" }),
    idUser: z.string(),
    image: z.string(),
    body: z.string(),
    totalLikes: z.number().int(),
    totalReplies: z.number().int(),
    isLiked: z.boolean(),
})

export const MetaSchema = z.object({
    cursor: z.number(),
})

export type Post = z.infer<typeof PostSchema> & Base

export type GetPostArgs = {
    postId: string
}

export type GetForYouArgs = {
    trigger: boolean
    cursor: number
    limit: number
}

export type ListPostByUserArgs = {
    idUser: string
    trigger: boolean
    cursor: number
    limit: number
}

export type GetAllPostArgs = {
    trigger: boolean
    cursor: number
    limit: number
}

export type CreatePostArgs = {
    fd: FormData
}

export type EditPostArgs = {}

export type DeletePostArgs = {}

interface IPost {
    getPost(args: GetPostArgs): Promise<APIResponse<{ post: Post } | null>>
    getForYou(args: GetForYouArgs): Promise<APIResponse<{ posts: Post[]; meta: { cursor: number } }>>
    listPostByUser(args: ListPostByUserArgs): Promise<APIResponse<{ posts: Post[]; meta: { cursor: number } }>>
    getAllPost(args: GetAllPostArgs): Promise<APIResponse<{ posts: Post[]; meta: { cursor: number } }>>
    createPost(args: CreatePostArgs): Promise<APIResponse<null>>
    editPost(args: EditPostArgs): Promise<void>
    deletePost(args: DeletePostArgs): Promise<void>
}

export class PostRepository implements IPost {
    public async getPost(args: GetPostArgs): Promise<APIResponse<{ post: Post } | null>> {
        let status = 500
        const fetchPost = await appFetch(`${API_URL}/post/${args.postId}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        try {
            const postData = z.object({ post: PostSchema.merge(BaseSchema) }).parse(fetchPost.data)

            return {
                res: { data: postData, message: fetchPost.message },
                status: status,
            }
        } catch (err) {
            console.error(err)
        }

        return {
            res: { data: null, message: fetchPost.message },
            status: status,
        }
    }

    public async listPostByUser(
        args: ListPostByUserArgs
    ): Promise<APIResponse<{ meta: { cursor: number }; posts: Post[] }>> {
        let status = 500
        const fetchPosts = await appFetch(
            `${API_URL}/post/user/${args.idUser}?cursor=${args.cursor}&limit=${args.limit}`,
            {
                method: "GET",
                headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
            }
        ).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        try {
            const postsData = z
                .object({ posts: z.array(PostSchema.merge(BaseSchema)), meta: MetaSchema })
                .parse(fetchPosts.data)

            return {
                res: { message: fetchPosts.message, data: postsData },
                status: status,
            }
        } catch (error) {
            console.error(error)
        }

        return {
            res: { message: fetchPosts.message, data: { posts: [], meta: { cursor: 0 } } },
            status: status,
        }
    }

    public async getAllPost(args: GetAllPostArgs): Promise<APIResponse<{ posts: Post[]; meta: { cursor: number } }>> {
        let status = 500
        const fetchPosts = await appFetch(`${API_URL}/post/?cursor=${args.cursor}&limit=${args.limit}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        try {
            const postsData = z
                .object({ meta: MetaSchema, posts: z.array(PostSchema.merge(BaseSchema)) })
                .parse(fetchPosts.data)

            return {
                res: { message: fetchPosts.message, data: postsData },
                status: status,
            }
        } catch (error) {
            console.error(error)
        }

        return {
            res: { message: fetchPosts.message, data: { posts: [], meta: { cursor: 0 } } },
            status: status,
        }
    }

    public async createPost(args: CreatePostArgs): Promise<APIResponse<null>> {
        let status = 500
        const fetchPosts = await appFetch(`${API_URL}/post/`, {
            method: "POST",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
            body: args.fd,
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        return {
            res: { message: fetchPosts.message, data: null },
            status: status,
        }
    }

    public static setCtx() {
        return setContext(POST_REPO_KEY, new this())
    }

    public static getCtx() {
        return getContext<ReturnType<typeof this.setCtx>>(POST_REPO_KEY)
    }

    public async getForYou(args: GetForYouArgs): Promise<APIResponse<{ posts: Post[]; meta: { cursor: number } }>> {
        throw new Error("Method not implemented.")
    }

    public async editPost(args: EditPostArgs): Promise<void> {
        throw new Error("Method not implemented.")
    }

    public async deletePost(args: DeletePostArgs): Promise<void> {
        throw new Error("Method not implemented.")
    }
}

const POST_REPO_KEY = Symbol("post_repo_key")
