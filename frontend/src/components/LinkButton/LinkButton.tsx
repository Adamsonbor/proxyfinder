import { useTheme } from "@mui/material";

interface Props {
	label?: string
	onClick?: () => void
}

export default function LinkButton(props: Props) {
	const theme = useTheme()

	return (
		<button
			onClick={props.onClick}
			style={{
				color: theme.palette.blueFilterTextIcon,
				background: 'transparent',
				textAlign: 'left',
				paddingTop: '16px',
				marginTop: 'auto',
				border: 'none',
			}}>
			{props.label || "button"}
		</button>
	);
}
