import { z } from "zod"

export const name = z
    .string({ message: "Name is not a string" })
    .min(1, { message: "Name is too short, minimum name is 1 character long" })
    .max(50, { message: "Name is too long, maximum name is 50 character long" })

export const username = z
    .string({ message: "Username is not a string" })
    .min(6, { message: "Username is too short, minimum username is 6 character long" })
    .max(50, { message: "Username is too long, maximum username is 50 character long" })
    .regex(/^[a-zA-Z0-9]+$/, { message: "username cannot contain special character" })

export const password = z
    .string({ message: "Password is not a string" })
    .min(8, { message: "Password is too short, minimum password is 8 character long" })
    .max(255, { message: "Password is too long, maximum password is 255 character long" })
