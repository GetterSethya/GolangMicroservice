import { writable } from "svelte/store"
import type { Post, drawerData } from "./types"
import type { User } from "./repository/user"

export const localUser = writable<User | undefined>(undefined)
export const drawerStore = writable(false)
export const drawerDataStore = writable<{
    drawerData: drawerData | undefined
}>({ drawerData: undefined })
export let isLoading = writable(false)
export let postDetailStore = writable<Post | undefined>(undefined)
export let profileStore = writable<User | null>(null)
