import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';
import { useTheme } from '@mui/material';

interface Props {
	values: string[]
	defaultValue: string
	setValue: (value: string) => void
	sx?: object
}

export default function RadioButtonGroup(props: Props) {
	const theme = useTheme();

	return (
		<FormControl sx={props.sx}>
			<RadioGroup
				sx={{
					'&& .MuiSvgIcon-root': {
						height: theme.typography.fontSize,
						width: 'fit-content',
					},
					display: 'flex',
					flexDirection: 'column',
					gap: '10px',
					marginTop: '16px',
				}}
				aria-labelledby="demo-radio-buttons-group-label"
				defaultValue={props.defaultValue}
				name="radio-buttons-group"
			>
				{
					props.values.map((value, index) => (
						<FormControlLabel
							sx={{
								color: theme.palette.text.black,
								'&&.MuiFormControlLabel-root': {
									margin: 0,
								},
								'& .MuiFormControlLabel-label': {
									fontSize: theme.typography.uppercaseSize,
								}
							}}
							key={index}
							value={value}
							control={<Radio
								checkedIcon={
									<img
										src="radio-checked.svg"
										height={theme.inputs.height}
										width={theme.inputs.width}
									/>
								}
								sx={{
									'&& .MuiSvgIcon-root': {
										height: theme.inputs.height,
										width: theme.inputs.width,
									},
									padding: 0,
									marginRight: '10px',
								}} />}
							label={value}
							onClick={() => props.setValue(value)}
						/>
					))
				}
			</RadioGroup>
		</FormControl>
	);
}
