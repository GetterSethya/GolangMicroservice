import { colorVariant, type margin, type padding } from "@ui/types"
import type { HTMLButtonAttributes } from "svelte/elements"
import Root from "./root.svelte"

export type ButtonVariants = "filled" | "ghost" | "soft" | "ringed" | "glass"

type ButtonProps = HTMLButtonAttributes & {
    bg?: (typeof colorVariant)[number]
    variant?: ButtonVariants
    padding?: padding
    margin?: margin
}

type ButtonEvents = {
    click: (
        e: MouseEvent & {
            currentTarget: EventTarget & HTMLButtonElement
        }
    ) => void
    keypress: (
        e: KeyboardEvent & {
            currentTarget: EventTarget & HTMLButtonElement
        }
    ) => void
}

export { Root, type ButtonProps, type ButtonEvents }
