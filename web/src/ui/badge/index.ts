import { colorVariant, type margin, type padding } from "@ui/types"
import type { HTMLAttributes } from "svelte/elements"
import Root from "./root.svelte"

export type BadgeVariants = "filled" | "ghost" | "soft" | "ringed" | "glass"

type BadgeProps = HTMLAttributes<HTMLDivElement> & {
    bg?: (typeof colorVariant)[number]
    variant?: BadgeVariants
    padding?: padding
    margin?: margin
}

export { Root, type BadgeProps }
