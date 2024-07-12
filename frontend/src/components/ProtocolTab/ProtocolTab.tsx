import { Chip, ChipProps } from "@mui/material"

export default function ProtocolTab(props: ChipProps) {
	return (
		<Chip
			{...props}
			size="small"
			sx={{
				...props.sx,
				height: "20px",
				padding: "1px 8px",
				width: "fit-content"
			}} />
	)
}
