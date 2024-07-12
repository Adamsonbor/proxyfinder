import { useState } from "react";
import Card from "../Card/Card";
import { useTheme } from "@mui/material";
import SearchMultipleAutocomplete from "../SearchMultipleAutocomplete/SearchMultipleAutocomplete";
import CheckboxForm from "../CheckboxForm/CheckboxForm";
import LinkButton from "../LinkButton/LinkButton";
import RadioButtonGroup from "../RadioButtonGroup/RadioButtonGroup";

interface Props {
	className?: string
	sx?: object
}

const values = [
	'Afghanistan', 'Albania', 'Algeria', 'Andorra', 'Angola', 'Antigua and Barbuda', 'Argentina', 'Armenia', 'Australia', 'Austria', 'Azerbaijan',
]

const protocols = [
	{ label: "HTTP", checked: true },
	{ label: "HTTPS", checked: false },
	{ label: "SOCKS4", checked: false },
	{ label: "SOCKS5", checked: false },
]

const available = [
	"All",
	"Available",
	"Unavailable",
]

export default function LeftPanel(props: Props) {
	const theme = useTheme()
	const [selectedValues, setSelectedValues] = useState<string[]>([]);
	const [protocolStates, setProtocolStates] = useState(protocols);

	function protocolHandler(label: string) {
		const newProtocolStates = protocolStates.map((protocol) => {
			if (protocol.label === label) {
				return { ...protocol, checked: !protocol.checked };
			}
			return protocol;
		});
		setProtocolStates(newProtocolStates);
	}


	return (
		<div
			className={props.className}
			style={{
				...props.sx,
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
							...props.sx,
							height: '100%',
						}}>
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
							values={available}
							defaultValue={available[0]}
							setValue={() => { }}
						/>
					</div>
				} />
		</div >
	);
}
