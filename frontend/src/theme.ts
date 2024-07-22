import { createTheme } from "@mui/material"

declare module '@mui/material' {
	interface ThemeOptions {
		inputs: {
			width: number
			height: number
		}
	}
	interface Theme extends ThemeOptions { }

	interface PaletteOptions {
		blue?: string
		gray?: string
		lightGray?: string
		lightBlue?: string

		backgroundWhite: string
		shapeFilterLight: string
		textBlack: string
		textGray: string
		stroke: string
		blueFilterTab: string
		blueFilterTextIcon: string
		greenTabAvailable: string
		greenTextAvailable: string
		redTabUnavailable: string
		redTextUnavailable: string
		grayTabProtocol: string
		grayTextProtocol: string
	}
	interface Palette extends PaletteOptions { }

	interface TypographyVariantsOptions {
		uppercaseSize: number
	}
	interface TypographyVariants extends TypographyVariantsOptions { }

	interface TypeText {
		black: string
		secondary: string
		link: string
		blue: string
		lightBlue: string
		red: string
		green: string
		purple: string
	}
	interface TypeBackground {
		blue: string
		red: string
		green: string
		purple: string
		lightBlue: string
	}
}


const lightTheme = createTheme({
	palette: {
		mode: 'light',
		blue: "#0479D6",
		gray: "#B8B8B8",
		lightGray: "#DDDDDD",
		lightBlue: "#4B86FA",

		backgroundWhite: "#FFFFFF",
		shapeFilterLight: "#F9F9F9",
		textBlack: "#141414",
		textGray: "#B8B8B8",
		stroke: "#ECECEC",
		blueFilterTab: "#E4EDF7",
		blueFilterTextIcon: "#5191DF",
		grayTabProtocol: "#F4F4F4",
		grayTextProtocol: "#A5A5A5",
		greenTabAvailable: "#E8FEE7E5",
		greenTextAvailable: "#026F00E5",
		redTabUnavailable: "#FFE2E2E5",
		redTextUnavailable: "#6F0000E5",
	},
	typography: {
		fontSize: 16,
		uppercaseSize: 14
	},
	inputs: {
		width: 16,
		height: 16,
	},
})

const darkTheme = createTheme({
	palette: {
		mode: 'dark',
		blue: "#0479D6",
		gray: "#B8B8B8",
		lightGray: "#DDDDDD",
		lightBlue: "#4B86FA",

		backgroundWhite: "#202125",
		shapeFilterLight: "#292A2E",
		textBlack: "#FFFFFF",
		textGray: "#B8B8B8",
		stroke: "#3E3E3E",
		blueFilterTab: "#E4EDF7F2",
		blueFilterTextIcon: "#1B85D9",
		grayTabProtocol: "#484848",
		grayTextProtocol: "#D3D3D3",
		greenTabAvailable: "#E8FEE7CC",
		greenTextAvailable: "#026F00E5",
		redTabUnavailable: "#FFE2E2E5",
		redTextUnavailable: "#6F0000E5",

	},
	typography: {
		fontSize: 16,
		uppercaseSize: 14
	},
	inputs: {
		width: 16,
		height: 16,
	},
})

export { lightTheme, darkTheme }
