

export interface IColors {
    main: string
    secondary: string
}

export interface ITheme {
    colors: IColors
}

export const theme: ITheme = {
    colors: {
        main: "red",
        secondary: "blue"
    }
}