import { writable } from "svelte/store"
import type { User, drawerData } from "./types"

export const localUser = writable<User | undefined>(undefined)
export const drawerStore = writable(false)
export const drawerDataStore = writable<drawerData[]>([])
