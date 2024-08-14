import Autocomplete from '@mui/material/Autocomplete';
import TextField from '@mui/material/TextField';
import Stack from '@mui/material/Stack';
import SearchTab from './SearchTab';
import { useTheme } from '@mui/material';
import { Search } from '@mui/icons-material';
import { useRef, useState } from 'react';

interface Props {
	label: string
	values?: string[]
	selectedValues?: string[]
	setSelectedValues?: (values: string[]) => void

	sx?: any
}

export default function SearchMultipleAutocomplete({
	label = "Search...",
	values = [],
	selectedValues = [],
	setSelectedValues = () => { },
	sx = {},
}: Props) {

	const theme = useTheme();
	const inputRef = useRef<HTMLDivElement>(null);
	const [valueState, setValueState] = useState<string>("");

	const handleSelect = (value: string) => {
		if (selectedValues?.includes(value)) {
			setSelectedValues(selectedValues.filter((v) => v !== value));
		} else {
			setSelectedValues([...selectedValues, value]);
		}

		inputRef.current?.querySelector('input')?.blur();
		setValueState("");
	};

	const removeSelectedValue = (value: string) => {
		setSelectedValues(selectedValues.filter((v) => v !== value));
	};

	return (
		<Stack spacing={1} sx={{
			...sx,
		}}>
			<Autocomplete
				id="free-solo-demo"
				sx={{
					'&& label': {
						top: -8,
					},
				}}
				value={valueState}
				onChange={(_: any, newValue: string | null) => {
					handleSelect(newValue || "");
				}}
				freeSolo
				filterSelectedOptions
				options={values}
				clearOnEscape

				renderInput={(params) => (
					<TextField
						{...params}
						value={valueState}
						onChange={(e) => setValueState(e.target.value)}
						ref={inputRef}
						sx={{
							width: sx?.width,
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
								paddingLeft: '10px',
								height: '30px',
							},
							'&& .MuiInputBase-input': {
								padding: 0,
								color: theme.palette.textBlack,
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
								{label}
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
							color: theme.palette.lightBlue,
						}}
						key={option}
						onClick={() => { handleSelect(option) }}>
						{option}
					</li>
				)}
			/>
			{
				selectedValues?.map((value, index) => (
					<SearchTab
						label={value}
						onClick={() => removeSelectedValue(value)}
						key={index} />
				))
			}
		</Stack >
	);
}
