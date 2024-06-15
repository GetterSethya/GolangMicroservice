import { z } from "zod"
import { name, password, username } from "./baseSchema"

export const registerSchema = z.object({
    username,
    name,
    password,
})
