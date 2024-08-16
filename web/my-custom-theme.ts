import type { CustomThemeConfig } from "@skeletonlabs/tw-plugin"

export const myCustomTheme: CustomThemeConfig = {
    name: "my-custom-theme",
    properties: {
        // =~= Theme Properties =~=
        "--theme-font-family-base": `ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace`,
        "--theme-font-family-heading": `ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, 'Liberation Mono', 'Courier New', monospace`,
        "--theme-font-color-base": "0 0 0",
        "--theme-font-color-dark": "255 255 255",
        "--theme-rounded-base": "6px",
        "--theme-rounded-container": "12px",
        "--theme-border-base": "1px",
        // =~= Theme On-X Colors =~=
        "--on-primary": "255 255 255",
        "--on-secondary": "255 255 255",
        "--on-tertiary": "255 255 255",
        "--on-success": "255 255 255",
        "--on-warning": "255 255 255",
        "--on-error": "255 255 255",
        "--on-surface": "255 255 255",
        // =~= Theme Colors  =~=
        // primary | #6b62df
        "--color-primary-50": "233 231 250", // #e9e7fa
        "--color-primary-100": "225 224 249", // #e1e0f9
        "--color-primary-200": "218 216 247", // #dad8f7
        "--color-primary-300": "196 192 242", // #c4c0f2
        "--color-primary-400": "151 145 233", // #9791e9
        "--color-primary-500": "107 98 223", // #6b62df
        "--color-primary-600": "96 88 201", // #6058c9
        "--color-primary-700": "80 74 167", // #504aa7
        "--color-primary-800": "64 59 134", // #403b86
        "--color-primary-900": "52 48 109", // #34306d
        // secondary | #3c3870
        "--color-secondary-50": "226 225 234", // #e2e1ea
        "--color-secondary-100": "216 215 226", // #d8d7e2
        "--color-secondary-200": "206 205 219", // #cecddb
        "--color-secondary-300": "177 175 198", // #b1afc6
        "--color-secondary-400": "119 116 155", // #77749b
        "--color-secondary-500": "60 56 112", // #3c3870
        "--color-secondary-600": "54 50 101", // #363265
        "--color-secondary-700": "45 42 84", // #2d2a54
        "--color-secondary-800": "36 34 67", // #242243
        "--color-secondary-900": "29 27 55", // #1d1b37
        // tertiary | #1d1b36
        "--color-tertiary-50": "221 221 225", // #dddde1
        "--color-tertiary-100": "210 209 215", // #d2d1d7
        "--color-tertiary-200": "199 198 205", // #c7c6cd
        "--color-tertiary-300": "165 164 175", // #a5a4af
        "--color-tertiary-400": "97 95 114", // #615f72
        "--color-tertiary-500": "29 27 54", // #1d1b36
        "--color-tertiary-600": "26 24 49", // #1a1831
        "--color-tertiary-700": "22 20 41", // #161429
        "--color-tertiary-800": "17 16 32", // #111020
        "--color-tertiary-900": "14 13 26", // #0e0d1a
        // success | #2f8816
        "--color-success-50": "224 237 220", // #e0eddc
        "--color-success-100": "213 231 208", // #d5e7d0
        "--color-success-200": "203 225 197", // #cbe1c5
        "--color-success-300": "172 207 162", // #accfa2
        "--color-success-400": "109 172 92", // #6dac5c
        "--color-success-500": "47 136 22", // #2f8816
        "--color-success-600": "42 122 20", // #2a7a14
        "--color-success-700": "35 102 17", // #236611
        "--color-success-800": "28 82 13", // #1c520d
        "--color-success-900": "23 67 11", // #17430b
        // warning | #c65d24
        "--color-warning-50": "246 231 222", // #f6e7de
        "--color-warning-100": "244 223 211", // #f4dfd3
        "--color-warning-200": "241 215 200", // #f1d7c8
        "--color-warning-300": "232 190 167", // #e8bea7
        "--color-warning-400": "215 142 102", // #d78e66
        "--color-warning-500": "198 93 36", // #c65d24
        "--color-warning-600": "178 84 32", // #b25420
        "--color-warning-700": "149 70 27", // #95461b
        "--color-warning-800": "119 56 22", // #773816
        "--color-warning-900": "97 46 18", // #612e12
        // error | #cf3a3a
        "--color-error-50": "248 225 225", // #f8e1e1
        "--color-error-100": "245 216 216", // #f5d8d8
        "--color-error-200": "243 206 206", // #f3cece
        "--color-error-300": "236 176 176", // #ecb0b0
        "--color-error-400": "221 117 117", // #dd7575
        "--color-error-500": "207 58 58", // #cf3a3a
        "--color-error-600": "186 52 52", // #ba3434
        "--color-error-700": "155 44 44", // #9b2c2c
        "--color-error-800": "124 35 35", // #7c2323
        "--color-error-900": "101 28 28", // #651c1c
        // surface | #2a2c32
        "--color-surface-50": "223 223 224", // #dfdfe0
        "--color-surface-100": "212 213 214", // #d4d5d6
        "--color-surface-200": "202 202 204", // #cacacc
        "--color-surface-300": "170 171 173", // #aaabad
        "--color-surface-400": "106 107 112", // #6a6b70
        "--color-surface-500": "42 44 50", // #2a2c32
        "--color-surface-600": "38 40 45", // #26282d
        "--color-surface-700": "32 33 38", // #202126
        "--color-surface-800": "25 26 30", // #191a1e
        "--color-surface-900": "21 22 25", // #151619
    },
}
