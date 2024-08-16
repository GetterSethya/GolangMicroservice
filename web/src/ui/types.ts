export const colorVariant = ["primary", "secondary", "tertiary", "success", "warning", "error", "surface"] as const

export const numberToken = [50, 100, 200, 300, 400, 500, 600, 700, 800, 900] as const

export type bgTheme = `bg-${(typeof colorVariant)[number]}-${(typeof numberToken)[number]}`

export const hoverClasses = {
    primary: "hover:bg-primary-500/70",
    secondary: "hover:bg-secondary-500/70",
    tertiary: "hover:bg-tertiary-500/70",
    success: "hover:bg-success-500/70",
    warning: "hover:bg-warning-500/70",
    error: "hover:bg-error-500/70",
    surface: "hover:bg-surface-500/70",
}

export const twSize = [
    0, 1, 1.5, 2, 2.5, 3, 3.5, 4, 5, 6, 7, 8, 9, 10, 11, 12, 14, 16, 20, 24, 28, 32, 36, 40, 44, 48, 52, 56, 60, 64, 72,
    80, 96,
] as const

export const twFractions = [
    "1/2",
    "1/3",
    "1/4",
    "1/5",
    "1/6",
    "1/7",
    "1/8",
    "1/9",
    "1/10",
    "1/11",
    "1/12",
    "2/3",
    "2/4",
    "2/5",
    "2/6",
    "2/7",
    "2/8",
    "2/9",
    "2/10",
    "2/11",
    "2/12",
    "3/4",
    "3/5",
    "3/6",
    "3/7",
    "3/8",
    "3/9",
    "3/10",
    "3/11",
    "3/12",
    "4/5",
    "4/6",
    "4/7",
    "4/8",
    "4/9",
    "4/10",
    "4/11",
    "4/12",
    "5/6",
    "5/7",
    "5/8",
    "5/9",
    "5/10",
    "5/11",
    "5/12",
    "6/7",
    "6/8",
    "6/9",
    "6/10",
    "6/11",
    "6/12",
    "7/8",
    "7/9",
    "7/10",
    "7/11",
    "7/12",
    "8/9",
    "8/10",
    "8/11",
    "8/12",
    "9/10",
    "9/11",
    "9/12",
    "10/11",
    "10/12",
    "11/12",
] as const

export const twDivider = [10, 20, 30, 40, 50, 60, 70, 80, 90] as const
export type twSizeObj = {
    start?: (typeof twSize)[number]
    end?: (typeof twSize)[number]
    top?: (typeof twSize)[number]
    bottom?: (typeof twSize)[number]
    vertical?: (typeof twSize)[number]
    horizontal?: (typeof twSize)[number]
}

export type padding = (typeof twSize)[number] | twSizeObj

export type margin = (typeof twSize)[number] | twSizeObj
