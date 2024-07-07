<script lang="ts">
    import Input from "@lib/components/Input.svelte"
    import { ProgressRadial, getToastStore, type ToastSettings } from "@skeletonlabs/skeleton"
    import { handleRegister } from "./register"
    import { AuthError } from "@lib/types"
    import { ZodError } from "zod"
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

<div class="w-full h-full flex flex-col text-surface-200">
    <div class="flex flex-col m-auto w-full lg:w-1/2">
        <div class="flex flex-col gap-2.5 w-fit p-5">
            <h1 class="h1">Create your account</h1>
        </div>
        <form
            method="post"
            on:submit|preventDefault={async (e) => {
                isLoading = true
                try {
                    await handleRegister(e)
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
            <Input required={true} disabled={isLoading} label={"Name"} placeholder="enter your name" type={"text"} />
            <Input
                required={true}
                disabled={isLoading}
                label={"Password"}
                placeholder="enter your password"
                type={"password"}
                minLength="8"
            />
            <Input
                required={true}
                disabled={isLoading}
                label={"Confirm password"}
                placeholder="confirm your password"
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
                <span>Register</span>
            </button>
        </form>
        <div class="p-5 text-surface-400">
            <span>Already have an account? </span>
            <a class="text-primary-700" href="/#/login">Login</a>
        </div>
    </div>
</div>
