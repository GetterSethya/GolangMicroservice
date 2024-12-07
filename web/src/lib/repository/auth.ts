import { z } from "zod"
import { APIResponseSchema, type APIResponse } from "./base"
import { API_URL } from "@/constant"
import { getContext, setContext } from "svelte"

export const AuthSchema = z.object({
    accessToken: z.string(),
    refreshToken: z.string(),
})

export type Auth = z.infer<typeof AuthSchema>

export type LoginArgs = {
    username: string
    password: string
}

export type RegisterArgs = {
    username: string
    name: string
    password: string
    confirmPassword: string
}

interface IAuth {
    login(args: LoginArgs): Promise<APIResponse<Auth | null>>
    register(args: RegisterArgs): Promise<APIResponse<null>>
    refresh(): Promise<APIResponse<Auth | null>>
}

export class AuthRepository implements IAuth {
    public async login(args: LoginArgs): Promise<APIResponse<Auth | null>> {
        let status = 500

        const req = await fetch(`${API_URL}/auth/login`, {
            method: "POST",
            body: JSON.stringify({ ...args }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        try {
            const authData = AuthSchema.parse(req.data)
            return {
                status: status,
                res: { data: authData, message: req.message },
            }
        } catch (error) {
            console.error(error)
        }

        return {
            status: status,
            res: { data: null, message: req.message },
        }
    }

    public async register(args: RegisterArgs): Promise<APIResponse<null>> {
        let status = 500

        const req = await fetch(`${API_URL}/auth/register`, {
            method: "POST",
            body: JSON.stringify({ ...args }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        return {
            status: status,
            res: { data: null, message: req.message },
        }
    }

    public async refresh(): Promise<APIResponse<Auth | null>> {
        let status = 500
        const refreshToken = localStorage.getItem("refreshToken")

        const req = await fetch(`${API_URL}/auth/refresh`, {
            method: "POST",
            body: JSON.stringify({ refreshToken }),
        }).then(async (r) => {
            status = r.status
            return APIResponseSchema.parse(await r.json())
        })

        try {
            const authData = AuthSchema.parse(req.data)
            return {
                status: status,
                res: { data: authData, message: req.message },
            }
        } catch (error) {
            console.error(error)
        }

        return {
            status: status,
            res: { data: null, message: req.message },
        }
    }

    public static setCtx() {
        return setContext(AUTH_REPO_KEY, new this())
    }

    public static getCtx() {
        return getContext<ReturnType<typeof this.setCtx>>(AUTH_REPO_KEY)
    }
}

const AUTH_REPO_KEY = Symbol("auth_repo_key")
