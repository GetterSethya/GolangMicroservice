import type { SlideParams, FadeParams, ScaleParams } from "svelte/transition"
import Slide from "./slide.svelte"
import Fade from "./fade.svelte"
import type { HTMLAttributes } from "svelte/elements"

type slideProps = HTMLAttributes<HTMLDivElement> & {
    config?: SlideParams
}
type fadeProps = HTMLAttributes<HTMLDivElement> & {
    config?: FadeParams
}
type scaleProps = HTMLAttributes<HTMLDivElement> & {
    config?: ScaleParams
}

export { type slideProps, type fadeProps, type scaleProps, Slide, Fade }
