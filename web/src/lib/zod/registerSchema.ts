import { z } from "zod"
import { name, password, username } from "./baseSchema"

export const registerSchema = z
    .object({
        username,
        name,
        password,
        confirmPassword: z
            .string({ message: "Confirm password is not a string" })
            .min(8, { message: "Confirm password is too short, minimum is 8 character long" })
            .max(255, { message: "Confirm password is too long, maximum is 255 character long" }),
    })
    .refine(
        (d) => {
            if (d.password) {
                return d.password === d.confirmPassword
            }
        },
        { message: "Password and confirm password is different", path: ["confirmPassword"] }
    )
