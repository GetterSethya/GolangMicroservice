import { twSize, type bgTheme, type margin, type padding } from "@ui/types"
import type { HTMLAnchorAttributes, HTMLAttributes } from "svelte/elements"
import Anchor from "./anchor.svelte"
import Flex from "./flex.svelte"
import Box from "./box.svelte"

type boxProps = HTMLAttributes<HTMLDivElement> & {
    padding?: padding
    margin?: margin
    bg?: bgTheme
}

type anchorProps = HTMLAnchorAttributes

type flexDirection = "row" | "col"
type flexJustify = "between" | "center" | "start" | "end" | "evenly"
type flexItems = "start" | "end" | "center" | "baseline" | "stretch"

type flexProps = boxProps & {
    direction?: flexDirection | `${flexDirection}-reverse`
    justify?: flexJustify
    alignItems?: flexItems
    gap?: (typeof twSize)[number]
}

export { type boxProps, type flexProps, type anchorProps, Anchor, Flex, Box }
