import { useTheme } from "@mui/material";
import SearchMultipleAutocomplete from "../SearchMultipleAutocomplete/SearchMultipleAutocomplete";
import { useState } from "react";

interface Props {
	label: string
	sx?: object
}

const values = [
	'Afghanistan', 'Albania', 'Algeria', 'Andorra', 'Angola', 'Antigua and Barbuda', 'Argentina', 'Armenia', 'Australia', 'Austria', 'Azerbaijan',
]

export default function SearchCard(props: Props) {
	const theme = useTheme()
	const [selectedValues, setSelectedValues] = useState<string[]>([]);

	return (
		<div
			className="d-flex flex-column"
			style={{
				...props.sx,
				width: '210px',
				minHeight: '224px',
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
			<SearchMultipleAutocomplete
				values={values}
				label="Country or Region"
				selectedValues={selectedValues}
				setSelectedValues={setSelectedValues}
				sx={{
					marginTop: "16px",
					width: "170px"
				}} />
			<button
			onClick={() => setSelectedValues([])}
			style={{
				color: theme.palette.blueFilterTextIcon,
				background: 'transparent',
				textAlign: 'left',
				paddingTop: '16px',
				marginTop: 'auto',
				border: 'none',
			}}>
			Clear all
		</button>
		</div >
	);
}
