import { AuthError } from "@lib/types"
import { registerSchema } from "@lib/zod"
import { push } from "svelte-spa-router"

export async function handleRegister(e: SubmitEvent) {
    const fd = new FormData(e.target as HTMLFormElement)
    const { username, name, password, confirmpassword } = Object.fromEntries(fd) as Record<string, string>

    if (password !== confirmpassword) {
        throw new AuthError("Password missmatch")
    }

    registerSchema.parse({
        username,
        name,
        password,
    })

    const res = await fetch("http://localhost/v1/auth/register", {
        method: "POST",
        body: JSON.stringify({
            username,
            name,
            password,
        }),
    })

    if (res.status !== 201) {
        throw new AuthError("Register failed, try again later")
    }

    // redirect ke halaman login
    push("/auth/login")
}
