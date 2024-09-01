<script lang="ts">
    import * as Input from "@ui/input"
    import { ProgressRadial, getToastStore, type ToastSettings } from "@skeletonlabs/skeleton"
    import { handleRegister } from "./register"
    import { AuthError } from "@lib/types"
    import { z, ZodError } from "zod"
    import { AuthRepository } from "@lib/repository/auth"
    import { createForm } from "felte"
    import { validator } from "@felte/validator-zod"
    import { registerSchema } from "@lib/zod"
    import { reporter } from "@felte/reporter-svelte"

    const toastStore = getToastStore()
    const authRepo = AuthRepository.getCtx()

    const { form, isSubmitting } = createForm<z.infer<typeof registerSchema>>({
        onSubmit: async (data) => {
            let errMessage: string[] = []
            try {
                await handleRegister(authRepo, data)
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
        extend: [validator({ schema: registerSchema }), reporter],
    })
</script>

<div class="w-full md:w-1/2 md:mx-auto h-full flex flex-col text-surface-200">
    <div class="flex flex-col m-auto w-full lg:w-1/2">
        <div class="flex flex-col gap-2.5 w-fit p-5">
            <h1 class="h1">Create your account</h1>
        </div>
        <form use:form method="post" class="p-5 flex flex-col gap-5">
            <Input.Text
                required={true}
                disabled={$isSubmitting}
                label={"Username"}
                name="username"
                placeholder="enter your username"
            />
            <Input.Text
                required={true}
                disabled={$isSubmitting}
                label={"Name"}
                name="name"
                placeholder="enter your name"
                type={"text"}
            />
            <Input.Password
                disabled={$isSubmitting}
                required={true}
                name="password"
                label={"Password"}
                placeholder="enter your password"
                type={"password"}
            />
            <Input.Password
                disabled={$isSubmitting}
                required={true}
                label={"Confirm password"}
                placeholder="confirm your password"
                name="confirmPassword"
                type={"password"}
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
                <span>Register</span>
            </button>
        </form>
        <div class="p-5 text-surface-400">
            <span>Already have an account? </span>
            <a class="text-primary-700" href="#/auth/login">Login</a>
        </div>
    </div>
</div>
