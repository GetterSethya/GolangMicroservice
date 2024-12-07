import { superValidate } from "sveltekit-superforms";
import type { PageLoad } from "./$types";
import {zod} from "sveltekit-superforms/adapters"
import { loginSchema } from "@/zod";

export const load:PageLoad = async()=>{

    return {
        form:await superValidate(zod(loginSchema))
    }
}
