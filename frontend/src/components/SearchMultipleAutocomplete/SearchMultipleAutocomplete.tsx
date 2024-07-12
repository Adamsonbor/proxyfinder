import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import Stack from '@mui/material/Stack';
import SearchTab from './SearchTab';
import { useTheme } from '@mui/material';
import { Search } from '@mui/icons-material';

const top100Films = [
	{ title: 'The Shawshank Redemption', year: 1994 },
	{ title: 'The Godfather', year: 1972 },
	{ title: 'The Godfather: Part II', year: 1974 },
	{ title: 'The Dark Knight', year: 2008 },
	{ title: '12 Angry Men', year: 1957 },
	{ title: "Schindler's List", year: 1993 },
	{ title: 'Pulp Fiction', year: 1994 },
]

interface Props {
	label: string
	values?: string[]
	selectedValues?: string[]
	setSelectedValues?: any

	sx?: any
}

export default function SearchMultipleAutocomplete(props: Props) {
	const theme = useTheme();
	// const [selectedValues, setSelectedValues] = useState<string[]>([]);
	const selectedValues = props.selectedValues || [];
	const setSelectedValues = props.setSelectedValues;
	const handleSelect = (value: string) => {
		if (selectedValues.includes(value)) {
			setSelectedValues(selectedValues.filter((val) => val !== value));
		} else {
			setSelectedValues([...selectedValues, value]);
		}
	};

	return (
		<Stack spacing={1} sx={{
			...props.sx,
		}}>
			<Autocomplete
				id="free-solo-demo"
				sx={{
					'&& label': {
						top: -8,
					},
				}}
				freeSolo
				filterSelectedOptions
				options={props.values || (top100Films.map((option) => option.title))}
				renderInput={(params) => (
					<TextField
						{...params}
						sx={{
							width: props.sx?.width,
							padding: 0,
							'& fieldset': {
								padding: 0,
							},
							'&&:hover fieldset': {
								borderWidth: "1px",
								borderColor: theme.palette.lightGray,
							},
							'&&:focus-within span': {
								display: 'none',
							},
							'&& span': {
								display: 'none',
							},
							'&& .MuiOutlinedInput-notchedOutline': {
								borderWidth: "1px",
								borderColor: theme.palette.stroke,
							},
							'&& .MuiInputBase-root': {
								padding: 0,
								height: '30px',
							},
							'&& .MuiInputBase-input': {
								padding: 0,
								color: theme.palette.text.black,
								fontSize: theme.typography.fontSize,
							}
						}}
						label={
							<div
								style={{
									width: 'fit-content',
									height: '100%',
									fontSize: '14px',
								}}>
								<Search
									sx={{
										fontSize: theme.typography.fontSize,
										padding: 0,
										margin: 0
									}} />
								{props.label}
							</div>
						} />
				)}
				renderOption={(props, option) => (
					<li
						{...props}
						style={{
							padding: '8px 16px',
							display: 'flex',
							alignItems: 'center',
							color: theme.palette.text.lightBlue,
						}}
						key={option}
						onClick={() => {
							handleSelect(option);
							console.log(selectedValues);
						}}>
						{option}
					</li>
				)}
				onChange={(_, value) => {
					if (!value) {
						setSelectedValues([]);
					}
				}}
			/>
			{
				selectedValues.map((value) => (
					<SearchTab
						label={value}
						onClick={() => handleSelect(value)}
						key={value} />
				))
			}
		</Stack >
	);
}
