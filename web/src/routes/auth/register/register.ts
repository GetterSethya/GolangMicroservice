import type { AuthRepository } from "@lib/repository/auth"
import { AuthError } from "@lib/types"
import { push } from "svelte-spa-router"

export async function handleRegister(
    authRepo: AuthRepository,
    data: { password: string; confirmPassword: string; username: string; name: string }
) {
    const { status } = await authRepo.register({
        username: data.username,
        name: data.name,
        password: data.password,
        confirmPassword: data.confirmPassword,
    })

    if (status !== 201) {
        throw new AuthError("Register failed, try again later")
    }

    // redirect ke halaman login
    push("/auth/login")
}
