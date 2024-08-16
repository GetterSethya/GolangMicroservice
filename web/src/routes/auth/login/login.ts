import { AuthError, type AuthResp, type ServerResp, type User } from "@lib/types"
import { loginSchema } from "@lib/zod"
import { replace } from "svelte-spa-router"

export async function handleLogin(e: SubmitEvent) {
    const fd = new FormData(e.target as HTMLFormElement)
    const { username, password } = Object.fromEntries(fd) as Record<string, string>

    // simulasi server ***
    // await simulateLatency(2000)

    console.log({ username, password })

    loginSchema.parse({
        username,
        password,
    })

    const res = await fetch("http://localhost/v1/auth/login", {
        method: "POST",
        body: JSON.stringify({
            username,
            password,
        }),
    })

    const resJson = (await res.json()) as ServerResp<null | AuthResp>

    if (res.status != 200) {
        throw new AuthError(resJson.message)
    }

    // set access token localstorage
    localStorage.setItem("accessToken", resJson.data?.accessToken as string)

    // set refresh token localstorage
    localStorage.setItem("refreshToken", resJson.data?.refreshToken as string)

    // redirect ke halaman index
    replace("/home/")
}
