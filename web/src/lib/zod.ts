import { z } from "zod";

export const loginSchema = z.object({
    username:z.string({message:"username cant be blank"})
        .min(8,{message:"username is too short, minimum 8 characters"})
        .max(100,{message:"username is too long, maximum 100 characters"}),
    password:z.string({message:"password cant be blank"})
        .min(8, {message:"password is too short, minimum 8 characters"})
        .max(100, {message:"password is too long, maximum 100 characters"})
})
export type LoginSchema = typeof loginSchema

export const registerSchema = z.object({
    username:z.string({message:"username cant be blank"})
        .min(8,{message:"username is too short, minimum 8 characters"})
        .max(100,{message:"username is too long, maximum 100 characters"}),
    name:z.string({message:"name cant be blank"})
        .min(1,{message:"name cant be blank"}),
    password:z.string({message:"password cant be blank"})
        .min(8, {message:"password is too short, minimum 8 characters"})
        .max(100, {message:"password is too long, maximum 100 characters"}),
    confirmPassword:z.string({message:"confirm password cant be blank"})
        .min(8, {message:"confirm password is too short, minimum 8 characters"})
        .max(100, {message:"confirm password is too long, maximum 100 characters"})
}).refine(d=>{
    if (d.password !== d.confirmPassword) {
        return false
    }

    return true
},{message:"invalid confirm password",path:["confirmPassword"]})
export type RegisterSchema = typeof registerSchema
