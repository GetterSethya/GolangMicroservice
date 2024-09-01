import type { AuthRepository } from "@lib/repository/auth"
import { AuthError } from "@lib/types"
import { replace } from "svelte-spa-router"

export async function handleLogin(authRepo: AuthRepository, data: { username: string; password: string }) {
    const { res, status } = await authRepo.login({ username: data.username, password: data.password })

    if (status != 200) {
        throw new AuthError(res.message)
    }

    // set access token localstorage
    localStorage.setItem("accessToken", res.data?.accessToken as string)

    // set refresh token localstorage
    localStorage.setItem("refreshToken", res.data?.refreshToken as string)

    // redirect ke halaman index
    replace("/app/home/")
}
