import Radio from '@mui/material/Radio';
import RadioGroup from '@mui/material/RadioGroup';
import FormControlLabel from '@mui/material/FormControlLabel';
import FormControl from '@mui/material/FormControl';
import { useTheme } from '@mui/material';

interface Props {
	values: string[]
	defaultValue?: string
	setValue?: any
	sx?: object
}

export default function RadioButtonGroup({
	values,
	defaultValue = values[0],
	setValue = () => { },
	sx = {},
}: Props) {
	const theme = useTheme();

	return (
		<FormControl sx={sx}>
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
				defaultValue={defaultValue}
				name="radio-buttons-group"
			>
				{
					values.map((value, index) => (
						<FormControlLabel
							sx={{
								color: theme.palette.textBlack,
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
							onClick={() => setValue(value)}
						/>
					))
				}
			</RadioGroup>
		</FormControl>
	);
}
