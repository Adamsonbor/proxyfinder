import { useEffect, useState } from "react";
import Card from "../Card/Card";
import SearchMultipleAutocomplete from "../SearchMultipleAutocomplete/SearchMultipleAutocomplete";
import CheckboxForm from "../CheckboxForm/CheckboxForm";
import LinkButton from "../LinkButton/LinkButton";
import RadioButtonGroup from "../RadioButtonGroup/RadioButtonGroup";
import { darkTheme, lightTheme } from "../../theme";
import ThemeSwitch from "../ThemeSwitch/ThemeSwitch";
import { Theme } from "@mui/material/styles";
import { CountryRepo } from "../../repos/country/repo";
import { useConfig } from "../../config";
import { Country } from "../../types";

interface Props {
	className?: string
	sx?: object
	theme?: Theme
	setTheme?: (theme: Theme) => void
	filter?: object
	setFilter?: (filter: object) => void
}

interface ProtocolState {
	label: string
	checked: boolean
	name: string
}

const protocols: ProtocolState[] = [
	{ label: "HTTP", checked: false, name: "http" },
	{ label: "HTTPS", checked: false, name: "https" },
	{ label: "SOCKS4", checked: false, name: "socks4" },
	{ label: "SOCKS5", checked: false, name: "socks5" },
]

const availableStatuses = [
	"All",
	"Available",
	"Unavailable",
]

export default function LeftPanel({
	className = "",
	sx = {},
	filter = {},
	setFilter = () => { },
	theme = lightTheme,
	setTheme = () => { },
}: Props) {
	const config = useConfig();
	const countryRepo = new CountryRepo(config);

	const [countries, setCountries] = useState<Country[]>([]);
	const [selectedCountries, setSelectedCountries] = useState<string[]>([]);
	const [protocolStates, setProtocolStates] = useState<ProtocolState[]>(protocols);
	const [selectedStatus, setSelectedStatus] = useState<string>(availableStatuses[0]);

	const toggleTheme = () => {
		setTheme(theme === lightTheme ? darkTheme : lightTheme);
	};


	useEffect(() => {
		countryRepo.GetAll({
			"page": 1,
			"perPage": 300,
		}).then((res) => {
			setCountries(res);
		})
	}, []);

	useEffect(() => {
		if (selectedCountries.length === 0) {
			delete filter["country_name"];
		} else {
			filter["country_name"] = selectedCountries;
		}

		if (protocolStates.filter((protocol: ProtocolState) => protocol.checked).length === 0) {
			delete filter["protocol"];
		} else {
			filter["protocol"] = protocolStates
				.filter((protocol: ProtocolState) => protocol.checked)
				.map((protocol: ProtocolState) => protocol.name);
		}

		if (selectedStatus === "All") {
			delete filter["status_name"];
		} else {
			filter["status_name"] = selectedStatus;
		}

		filter["page"] = 1;

		setFilter({ ...filter });
	}, [selectedCountries, protocolStates, selectedStatus]);


	return (
		<div
			className={className}
			style={{
				...sx,
				display: 'flex',
				flexDirection: 'column',
				gap: '10px',
			}}
		>
			<Card
				sx={{ minHeight: '224px', }}
				label="COUNTRY"
				renderContent={
					<div
						className="d-flex flex-column"
						style={{
							...sx,
							height: '100%',
						}}>
						<SearchMultipleAutocomplete
							values={countries?.map((country) => country.name) ?? []}
							label="Country or Region"
							selectedValues={selectedCountries}
							setSelectedValues={setSelectedCountries}
							sx={{
								marginTop: "16px",
								// width: "170px"
							}} />
						<button
							onClick={() => setSelectedCountries([])}
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
				} />
			<Card
				sx={{ minHeight: '224px', }}
				label="PROTOCOLS"
				renderContent={
					<div style={{
						display: 'flex',
						flexDirection: 'column',
						gap: '10px',
						marginTop: '16px',
					}}>
						{protocolStates.map((protocol, index) => (
							<CheckboxForm
								key={index}
								label={protocol.label}
								onClick={() => protocolHandler(protocol.label)}
								checked={protocol.checked}
							/>
						))}
						<LinkButton
							onClick={() => {
								setProtocolStates(protocols.map((protocol) => ({ ...protocol, checked: false })));
							}}
							label="Clear all" />
					</div>
				} />
			<Card
				label="AVAILABLE"
				renderContent={
					<div style={{}}>
						<RadioButtonGroup
							values={availableStatuses}
							defaultValue={availableStatuses[0]}
							setValue={setSelectedStatus}
						/>
					</div>
				} />
			<div style={{ marginTop: 'auto', marginBottom: '40px' }}>
				<div
					style={{
						display: 'flex',
						alignItems: 'center',
						gap: '10px',
						marginBottom: '21px',
					}}>
					<span style={{ color: theme.palette.textGray }}>Theme</span>
					<ThemeSwitch onChange={toggleTheme} />
				</div>
				<a
					style={{
						color: theme.palette.textGray,
					}}
					href="https://github.com/Adamsonbor">Project developer</a>
			</div>
		</div >
	);

	function protocolHandler(label: string) {
		const newProtocolStates = protocolStates.map((protocol) => {
			if (protocol.label === label) {
				return { ...protocol, checked: !protocol.checked };
			}
			return protocol;
		});

		setProtocolStates(newProtocolStates);
	}
}
