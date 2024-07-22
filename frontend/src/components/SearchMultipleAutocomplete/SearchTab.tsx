import { Close } from "@mui/icons-material";
import { useTheme } from "@mui/material";

interface Props {
	label: string;
	onClick: () => void;
	sx?: any;
}

export default function SearchTab(props: Props) {
	const theme = useTheme();

	return (
		<div
			style={{
				padding: '8px 16px',
				display: 'flex',
				alignItems: 'center',
				borderRadius: '24px',
				width: 'fit-content',
				lineHeight: '1px',
				backgroundColor: theme.palette.blueFilterTab,
				...props.sx
			}} >
			<span
				style={{
					fontSize: '12px',
					color: theme.palette.blueFilterTextIcon,
					marginRight: '1rem'
				}}>{props.label}</span>
			<Close
				sx={{
					color: theme.palette.blue,
					width: theme.inputs.width,
					height: theme.inputs.height,
					cursor: 'pointer'
				}}
				onClick={props.onClick} />
		</div>
	);
}
