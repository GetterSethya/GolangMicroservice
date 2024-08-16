import type { twSize, twSizeObj } from "@ui/types"
import clsx, { type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"

export function cn(...inputs: ClassValue[]) {
    return twMerge(clsx(inputs))
}

export function renderTwClass(prefix: string, size: (typeof twSize)[number] | twSizeObj | undefined) {
    if (typeof size === "object") {
        let start = size.start ? prefix + "s-" + size.start : undefined
        let end = size.end ? prefix + "e-" + size.end : undefined
        let top = size.top ? prefix + "t-" + size.top : undefined
        let bottom = size.bottom ? prefix + "b-" + size.bottom : undefined
        let horizontal = size.horizontal ? prefix + "x-" + size.horizontal : undefined
        let vertical = size.vertical ? prefix + "y-" + size.vertical : undefined

        return [start, end, top, bottom, horizontal, vertical]
            .filter((value): value is string => value !== undefined)
            .join(" ")
    }
    if (typeof size === "number") {
        return prefix + "-" + size
    }

    return undefined
}
