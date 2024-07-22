import { useTheme } from "@mui/material";
import { DataGrid } from "@mui/x-data-grid";
import ProtocolTab from "../ProtocolTab/ProtocolTab";
import { Country, ProxyRow } from "../../types";

interface Props {
	className?: string;
	proxies?: ProxyRow[];
	countries?: Country[];
	sx?: object;
}

export default function InfiniteTable({
	className = '',
	proxies = [],
	countries = [],
	sx = {},
}: Props) {
	const theme = useTheme();

	const columns = [
		{ field: 'Ip', headerName: 'IP', minWidth: 100, flex: 1 },
		{ field: 'Port', headerName: 'PORT', minWidth: 100, flex: 1 },
		{
			field: 'CountryName',
			headerName: 'COUNTRY',
			minWidth: 150,
			flex: 1,
			renderCell: (params: any) => (
				<>
					<span
						style={{ marginRight: "5px" }}
						className={`fi fi-${params.row.CountryCode.toLowerCase()}`}>
					</span>
					{params.row.CountryName}
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
			field: 'Status',
			headerName: 'AVAILABLE',
			minWidth: 150,
			flex: 1,
			renderCell: (params: any) => (
				<>
					<ProtocolTab
						label={params.row.Status}
						sx={{
							backgroundColor: availableBgColor(params.row.Status),
							color: availableTextColor(params.row.Status),
						}}
					/>
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
			autoHeight
			hideFooter
			className={className}
			sx={{
				...sx,
				border: "none",
				fontSize: '14px',
				color: theme.palette.textBlack,
				'& .MuiDataGrid-columnHeaders': {
					color: theme.palette.textGray,
				},
				'& .MuiDataGrid-topContainer::after': {
					display: "none",
				},
				'& .MuiDataGrid-columnHeader:focus': {
					outline: "none",
					border: "none",
				},
				'& [role="row"]': {
					backgroundColor: 'transparent !important',
				},
				'& [role="rowgroup"]': {
					height: "100%",
					overflowY: "scroll",
				},
				'& .MuiDataGrid-virtualScrollerContent': {
					height: "calc(100vh - 120px) !important",
					// {
						// xs: "calc(100vh - 120px) !important",
						// sm: "calc(100vh - 120px) !important",
						// md: "calc(100vh - 120px) !important",
						// lg: "calc(100vh - 120px) !important",
						// xl: "calc(100vh - 120px) !important",
					// }
				},
				'& .MuiDataGrid-cell': {
					border: "none",
					outline: "none",
					fontSize: theme.typography.fontSize,
					borderTop: "1px solid " + theme.palette.stroke,
				},
				'& .MuiDataGrid-cell:focus': {
					outline: "none",
					border: "none",
				}
			}}
			disableColumnMenu
		/>
	);

	function availableBgColor(name: string): string {
		switch (name) {
			case "Available":
				return theme.palette.greenTabAvailable;
			case "Unavailable":
				return theme.palette.redTabUnavailable;
			default:
				return theme.palette.redTabUnavailable;
		}
	}

	function availableTextColor(name: string): string {
		switch (name) {
			case "Available":
				return theme.palette.greenTextAvailable;
			case "Unavailable":
				return theme.palette.redTextUnavailable;
			default:
				return theme.palette.redTextUnavailable;
		}
	}

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
}
