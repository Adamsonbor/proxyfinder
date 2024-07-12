import { useTheme } from "@mui/material";
import { DataGrid } from "@mui/x-data-grid";
import ProtocolTab from "../ProtocolTab/ProtocolTab";

const rows = [
	{ 'ip': '1.1.1.1', 'port': '80', 'country': 'United States', 'protocols': ['HTTPS', "HTTP"], 'response': '200', 'updated': '2021-01-01', 'available': true },
	{ 'ip': '1.1.1.1', 'port': '80', 'country': 'United States', 'protocols': ['HTTPS', "HTTP"], 'response': '200', 'updated': '2021-01-01', 'available': false },
	{ 'ip': '1.1.1.1', 'port': '80', 'country': 'United States', 'protocols': ['HTTPS', "HTTP"], 'response': '200', 'updated': '2021-01-01', 'available': true },
	{ 'ip': '1.1.1.1', 'port': '80', 'country': 'United States', 'protocols': ['HTTPS', "HTTP", "SOCKS4", "SOCKS5"], 'response': '200', 'updated': '2021-01-01', 'available': true },
]

export default function TableHeader() {
	const theme = useTheme();
	const getBackgroundColor = (protocol: string) => {
		switch (protocol) {
			case "HTTP":
				return theme.palette.background.blue;
			case "HTTPS":
				return theme.palette.background.purple;
			case "SOCKS4":
				return theme.palette.background.green;
			case "SOCKS5":
				return theme.palette.background.red;
			default:
				return theme.palette.blue;
		}
	}
	const getTextColor = (protocol: string) => {
		switch (protocol) {
			case "HTTP":
				return theme.palette.text.blue;
			case "HTTPS":
				return theme.palette.text.purple;
			case "SOCKS4":
				return theme.palette.text.green;
			case "SOCKS5":
				return theme.palette.text.red;
			default:
				return theme.palette.blue;
		}
	}
	const getProtocol = (protocol: string, index: number) => {
		return (
			<>
				<ProtocolTab
					key={index}
					label={protocol}
					sx={{
						backgroundColor: getBackgroundColor(protocol),
						color: getTextColor(protocol),
					}}
				/>
			</>
		)
	}

	const columns = [
		{ field: 'ip', headerName: 'IP', minWidth: 100, flex: 1 },
		{ field: 'port', headerName: 'PORT', minWidth: 100, flex: 1 },
		{ field: 'country', headerName: 'COUNTRY', minWidth: 150, flex: 1 },
		{
			field: 'protocols',
			headerName: 'PROTOCOLS',
			minWidth: 100,
			flex: 4,
			renderCell: (params: any) => (
				<>
					{params.row.protocols.map((protocol: string, index: number) => (
						getProtocol(protocol, index)
					))}
				</>
			),
		},
		{ field: 'response', headerName: 'RESPONSE', minWidth: 100, flex: 1 },
		{ field: 'updated', headerName: 'UPDATED', minWidth: 150, flex: 1 },
		{
			field: 'available',
			headerName: 'AVAILABLE',
			minWidth: 150,
			flex: 1,
			renderCell: (params: any) => (
				<>
					{params.row.available ?
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
	return (
		<DataGrid
			rows={rows.map((row, index) => ({ ...row, id: index }))}
			columns={columns}
			hideFooter
			autoHeight
			hideFooterPagination
			hideFooterSelectedRowCount
			sx={{
				border: "none",
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
				},
				'& .MuiDataGrid-cell:focus': {
					outline: "none",
					border: "none",
				}
			}}
			disableColumnMenu
		/>
	);
}
