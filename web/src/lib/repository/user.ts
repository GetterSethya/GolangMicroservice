import { getContext, setContext } from "svelte"
import { appFetch } from "@/appFetch"
import { API_URL } from "@/constant"
import * as Jose from "jose"
import { z } from "zod"
import { APIResponseSchema, BaseSchema, type APIResponse, type Base, type Init } from "./base"
import { writable } from "svelte/store"

export const UserSchema = z.object({
    username: z.string(),
    name: z.string(),
    profile: z.string(),
    totalFollower: z.number(),
    totalFollowing: z.number(),
})

export type User = z.infer<typeof UserSchema> & Base

type GetUserByIdArgs = {
    id: string
    init?: Init
}

type GetUserByUsernameArgs = {
    username: string
    init?: Init
}

type UpdateUserByJWTArgs = {
    id: string
    name?: string
    profile?: string
    init?: Init
}

type UpdateUserPassword = {
    currentPassword: string
    newPassword: string
    confirmNewPassword: string
    init?: Init
}

type GetLocalUserDataArgs = {
    init?: Init
}

interface IUser {
    getUserById(args: GetUserByIdArgs): Promise<APIResponse<{ user: User } | null>>
    getUserByUsername(args: GetUserByUsernameArgs): Promise<APIResponse<{ user: User } | null>>
    getLocalUserData(args?: GetLocalUserDataArgs): Promise<APIResponse<{ user: User } | null>>
    updateUserByJWT(args: UpdateUserByJWTArgs): Promise<APIResponse<null>>
    updateUserPassword(args: UpdateUserPassword): Promise<APIResponse<null>>
}

export class UserRepository implements IUser {
    public async getLocalUserData(args?: GetLocalUserDataArgs): Promise<APIResponse<{ user: User } | null>> {
        const accessToken = localStorage.getItem("accessToken")
        const jwt = Jose.decodeJwt(accessToken!)

        let status = 500

        const req = await appFetch(`${API_URL}/user/${jwt.sub}`, {
            ...args?.init,
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        try {
            const userData = z.object({ user: UserSchema.merge(BaseSchema) }).parse(req.data)
            return {
                res: { message: req.message, data: userData },
                status: status,
            }
        } catch (error) {
            console.error(error)
        }

        return {
            res: { message: req.message, data: null },
            status: status,
        }
    }

    public async getUserById(args: GetUserByIdArgs): Promise<APIResponse<{ user: User } | null>> {
        let status = 500

        const req = await appFetch(`${API_URL}/user/${args.id}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        try {
            const userData = z.object({ user: UserSchema.merge(BaseSchema) }).parse(req.data)
            return {
                res: { message: req.message, data: userData },
                status: status,
            }
        } catch (error) {
            console.error(error)
        }

        return {
            res: { message: req.message, data: null },
            status: status,
        }
    }

    public async getUserByUsername(args: GetUserByUsernameArgs): Promise<APIResponse<{ user: User } | null>> {
        let status = 500

        const req = await appFetch(`${API_URL}/user/username/${args.username}`, {
            method: "GET",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        try {
            const userData = z.object({ user: UserSchema.merge(BaseSchema) }).parse(req.data)
            return {
                res: { message: req.message, data: userData },
                status: status,
            }
        } catch (err) {
            console.error(err)
        }

        return {
            res: { message: req.message, data: null },
            status: status,
        }
    }

    public async updateUserByJWT(args: UpdateUserByJWTArgs): Promise<APIResponse<null>> {
        const fd = new FormData()

        let status = 500
        fd.set("id", args.id)
        args.name && fd.set("name", args.name)
        args.profile && fd.set("profile", args.profile)

        const fetchUpdate = await appFetch(`${API_URL}/user/update`, {
            method: "POST",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
            body: fd,
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        return {
            res: { message: fetchUpdate.message, data: null },
            status: status,
        }
    }

    public async updateUserPassword(args: UpdateUserPassword): Promise<APIResponse<null>> {
        let status = 500

        const fetchUpdate = await appFetch(`${API_URL}/user/update_password`, {
            method: "POST",
            headers: new Headers({ Authorization: localStorage.getItem("accessToken") as string }),
            body: JSON.stringify({ ...args }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        return {
            res: { message: fetchUpdate.message, data: null },
            status: status,
        }
    }

    public setLocalUserCtx(userData: User | null) {
        const store = writable(userData)
        return setContext(LOCAL_USER, store)
    }

    public getLocalUserCtx() {
        return getContext<ReturnType<typeof this.setLocalUserCtx>>(LOCAL_USER)
    }

    public static setCtx() {
        return setContext(USER_REPO_KEY, new this())
    }

    public static getCtx() {
        return getContext<ReturnType<typeof this.setCtx>>(USER_REPO_KEY)
    }
}

const LOCAL_USER = Symbol("local_user")
const USER_REPO_KEY = Symbol("user_repo_key")
