import { Chip, ChipProps } from "@mui/material"

export default function ProtocolTab(props: ChipProps) {
	return (
		<Chip
			{...props}
			size="small"
			sx={{
				...props.sx,
				height: "20px",
				padding: "10px 16px",
				width: "fit-content",
				'& .MuiChip-label': {
					fontSize: "10px",
				},
			}} />
	)
}
