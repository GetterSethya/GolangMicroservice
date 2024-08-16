<script lang="ts">
    import Input from "@lib/components/Input.svelte"
    import { ProgressRadial, type ToastSettings } from "@skeletonlabs/skeleton"
    import { ZodError } from "zod"
    import { AuthError } from "@lib/types"
    import { handleLogin } from "@routes/auth/login/login"
    import { getToastStore } from "@skeletonlabs/skeleton"
    let isLoading = false
    let errMessage: string[] = []

    const toastStore = getToastStore()

    $: if (errMessage.length > 0) {
        errMessage.map((d) => {
            const t: ToastSettings = {
                message: d,
                background: "bg-error-500",
                timeout: 3500,
            }

            toastStore.trigger(t)
        })
    }
</script>

<div class="w-full md:w-1/2 md:mx-auto h-full flex flex-col text-surface-200">
    <div class="flex flex-col m-auto w-full lg:w-1/2">
        <div class="flex flex-col gap-2.5 w-fit p-5">
            <h1 class="h1">Welcome back</h1>
            <span class="text-surface-400">Please enter your username & password to continue</span>
        </div>
        <form
            method="post"
            on:submit|preventDefault={async (e) => {
                isLoading = true
                try {
                    await handleLogin(e)
                } catch (err) {
                    console.error(err)
                    switch (true) {
                        case err instanceof ZodError:
                            errMessage = err.errors.map((e) => e.message)
                            break

                        case err instanceof AuthError:
                            errMessage = [err.message]
                            break

                        default:
                            errMessage = ["Something went wrong"]
                            break
                    }
                }

                isLoading = false
            }}
            class="p-5 flex flex-col gap-5"
        >
            <Input
                required={true}
                disabled={isLoading}
                label={"Username"}
                placeholder="enter your username"
                type={"text"}
                minLength="6"
            />
            <Input
                required={true}
                disabled={isLoading}
                label={"Password"}
                placeholder="enter your password"
                type={"password"}
                minLength="8"
            />
            <button disabled={isLoading} class="btn font-bold variant-filled-primary">
                {#if isLoading}
                    <ProgressRadial
                        stroke={100}
                        width="w-5"
                        meter="stroke-surface-500"
                        track="stroke-surface-500/30"
                        strokeLinecap="butt"
                    />
                {/if}
                <span>Login</span>
            </button>
        </form>
        <div class="p-5 text-surface-400">
            <span>Dont have an account yet? </span>
            <a href="#/auth/register" class="text-primary-700">Register</a>
        </div>
    </div>
</div>
