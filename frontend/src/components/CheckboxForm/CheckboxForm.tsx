import { FormControlLabel, useTheme } from '@mui/material';
import Checkbox from '@mui/material/Checkbox';

interface Props extends React.ComponentProps<typeof Checkbox> {
	label: string;
	chacked?: boolean;
	sx?: object;
}

export default function CheckboxForm({ label, checked, sx, ...other }: Props) {
	const theme = useTheme();

	return (
		<>
			<FormControlLabel
				sx={{
					...sx,
					'& .MuiSvgIcon-root': {
						width: theme.inputs.width,
						height: theme.inputs.height,
					},
					'& .MuiFormControlLabel-label': {
						fontSize: theme.typography.uppercaseSize,
					},
					margin: 0,
					color: theme.palette.textBlack,
				}}
				control={
					<Checkbox
						sx={{
							padding: 0,
							marginRight: '10px',
							color: theme.palette.gray,
						}}
						checked={checked}
						checkedIcon={<>
							<img src="checkbox-checked.svg"
								style={{
									width: theme.inputs.width,
									height: theme.inputs.height,
								}} />
						</>}
						{...other}
					/>
				}
				label={label} />
		</>
	);
}
