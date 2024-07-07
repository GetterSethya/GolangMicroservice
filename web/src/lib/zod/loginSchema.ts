import { z } from "zod"
import { password, username } from "./baseSchema"

export const loginSchema = z.object({
    username,
    password,
})
