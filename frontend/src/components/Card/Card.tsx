import { useTheme } from "@mui/material";

interface Props {
	label: string
	renderContent?: any
	sx?: object
}

export default function Card(props: Props) {
	const theme = useTheme()

	return (
		<div
			className="d-flex flex-column"
			style={{
				...props.sx,
				backgroundColor: theme.palette.shapeFilterLight,
				border: '1px solid ' + theme.palette.stroke,
				borderRadius: '6px',
				padding: '20px',
			}}>
			<span style={{
				color: theme.palette.textGray,
				margin: 0,
				padding: 0,
			}}>{props.label}</span>
			{props.renderContent}
		</div >
	);
}
