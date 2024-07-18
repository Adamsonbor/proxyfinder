import { useTheme } from "@mui/material";
import { DataGrid } from "@mui/x-data-grid";
import ProtocolTab from "../ProtocolTab/ProtocolTab";
import { Country, Proxy } from "../../utils/api/api";

interface Props {
	className?: string;
	proxies?: Proxy[];
	countries?: Country[];
}

export default function InfiniteTable({
	className = '',
	proxies = [],
	countries = [],
}: Props) {
	const theme = useTheme();

	const columns = [
		{ field: 'Ip', headerName: 'IP', minWidth: 100, flex: 1 },
		{ field: 'Port', headerName: 'PORT', minWidth: 100, flex: 1 },
		{
			field: 'CountryId',
			headerName: 'COUNTRY',
			minWidth: 150,
			flex: 1,
			renderCell: (params: any) => (
				<>
				<span 
				style={{marginRight: "5px"}}
				className={`fi fi-${getCountryCode(countries, params.row.CountryId).toLowerCase()}`}>
				</span>
					{getCountryName(countries, params.row.CountryId)}
				</>
			),
		},
		{
			field: 'Protocol',
			headerName: 'PROTOCOL',
			minWidth: 100,
			flex: 1,
			renderCell: (params: any) => (
				<ProtocolTab
					label={params.row.Protocol}
					sx={{
						backgroundColor: theme.palette.grayTabProtocol,
						color: theme.palette.grayTextProtocol,
					}}
				/>
			),
		},
		{
			field: 'ResponseTime',
			headerName: 'RESPONSE',
			minWidth: 100, flex: 1,
			renderCell: (params: any) => (
				<>
					{params.row.ResponseTime}ms
				</>
			)
		},
		{ field: 'UpdatedAt', headerName: 'UPDATED', minWidth: 150, flex: 1 },
		{
			field: 'StatusId',
			headerName: 'AVAILABLE',
			minWidth: 150,
			flex: 1,
			renderCell: (params: any) => (
				<>
					{params.row.StatusId == 2 ?
						<ProtocolTab
							label="Available"
							sx={{
								backgroundColor: theme.palette.background.green,
								color: theme.palette.text.green,
							}}
						/>
						:
						<ProtocolTab
							label="Unavailable"
							sx={{
								backgroundColor: theme.palette.background.red,
								color: theme.palette.text.red,
							}}
						/>
					}
				</>
			)
		},
	]
	if (!proxies || !countries) {
		return <></>;
	}
	return (
		<DataGrid
			rows={proxies.map((proxy, index) => ({ ...proxy, id: index }))}
			columns={columns}
			hideFooter
			autoHeight
			hideFooterPagination
			hideFooterSelectedRowCount
			className={className}
			sx={{
				border: "none",
				fontSize: '14px',
				'& .MuiDataGrid-columnHeaders': {
					color: theme.palette.text.secondary,
				},
				'& .MuiDataGrid-columnHeader:focus': {
					outline: "none",
					border: "none",
				},
				'& .MuiDataGrid-cell': {
					border: "none",
					outline: "none",
					fontSize: theme.typography.fontSize,
				},
				'& .MuiDataGrid-cell:focus': {
					outline: "none",
					border: "none",
				}
			}}
			disableColumnMenu
		/>
	);

	// function getBackgroundColor(protocol: string) {
	// 	switch (protocol) {
	// 		case "HTTP":
	// 			return theme.palette.background.blue;
	// 		case "HTTPS":
	// 			return theme.palette.background.purple;
	// 		case "SOCKS4":
	// 			return theme.palette.background.green;
	// 		case "SOCKS5":
	// 			return theme.palette.background.red;
	// 		default:
	// 			return theme.palette.blue;
	// 	}
	// }
	// function getTextColor(protocol: string) {
	// 	switch (protocol) {
	// 		case "HTTP":
	// 			return theme.palette.text.blue;
	// 		case "HTTPS":
	// 			return theme.palette.text.purple;
	// 		case "SOCKS4":
	// 			return theme.palette.text.green;
	// 		case "SOCKS5":
	// 			return theme.palette.text.red;
	// 		default:
	// 			return theme.palette.blue;
	// 	}
	// }

	function getCountryCode(countries: Country[], countryId: number) {
		return countries.find((country: Country) => country.Id == countryId)?.Code || "Unknown"
	}

	function getCountryName(countries: Country[], countryId: number) {
		return countries.find((country: Country) => country.Id == countryId)?.Name || "Unknown"
	}
}
