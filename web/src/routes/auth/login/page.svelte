<script lang="ts">
    import { ProgressRadial, type ToastSettings } from "@skeletonlabs/skeleton"
    import { z, ZodError } from "zod"
    import { AuthError } from "@lib/types"
    import { handleLogin } from "@routes/auth/login/login"
    import { getToastStore } from "@skeletonlabs/skeleton"
    import { AuthRepository } from "@lib/repository/auth"
    import { createForm } from "felte"
    import { validator } from "@felte/validator-zod"
    import { loginSchema } from "@lib/zod"
    import { reporter } from "@felte/reporter-svelte"
    import * as Input from "@ui/input/"

    const toastStore = getToastStore()
    const authRepo = AuthRepository.getCtx()

    const { form, isSubmitting } = createForm<z.infer<typeof loginSchema>>({
        onSubmit: async (data) => {
            let errMessage: string[] = []
            try {
                await handleLogin(authRepo, data)
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
            } finally {
                if (errMessage.length > 0) {
                    errMessage.map((d) => {
                        const t: ToastSettings = {
                            message: d,
                            background: "bg-error-500",
                            timeout: 3500,
                        }

                        toastStore.trigger(t)
                    })
                }
            }
        },
        extend: [validator({ schema: loginSchema }), reporter],
    })
</script>

<div class="w-full md:w-1/2 md:mx-auto h-full flex flex-col text-surface-200">
    <div class="flex flex-col m-auto w-full lg:w-1/2">
        <div class="flex flex-col gap-2.5 w-fit p-5">
            <h1 class="h1">Welcome back</h1>
            <span class="text-surface-400">Please enter your username & password to continue</span>
        </div>
        <form use:form method="post" class="p-5 flex flex-col gap-5">
            <Input.Text
                name={"username"}
                required={true}
                disabled={$isSubmitting}
                label={"Username"}
                placeholder="enter your username"
            />
            <Input.Password
                name={"password"}
                required={true}
                disabled={$isSubmitting}
                label={"Password"}
                placeholder="enter your password"
            />
            <button disabled={$isSubmitting} class="btn font-bold variant-filled-primary">
                {#if $isSubmitting}
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
