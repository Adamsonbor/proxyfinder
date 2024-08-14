import { Box, useTheme } from "@mui/material";
import { DataGrid } from "@mui/x-data-grid";
import ProtocolTab from "../ProtocolTab/ProtocolTab";
import { Favorits, ProxyRow } from "../../types";
import { IoStar, IoStarOutline } from "react-icons/io5";

interface Props {
	className?: string;
	proxies?: ProxyRow[];
	favorits?: Favorits[];
	favoriteHandler?: (proxy_id: number, isFavorite: boolean) => void;
	sx?: object;
}

export default function InfiniteTable({
	className = '',
	proxies = [],
	favorits = [],
	favoriteHandler = () => { },
	sx = {},
}: Props) {
	const theme = useTheme();

	const columns = [
		{ field: 'ip', headerName: 'IP', minWidth: 100, flex: 1 },
		{ field: 'port', headerName: 'PORT', minWidth: 100, flex: 1 },
		{
			field: 'country_name',
			headerName: 'COUNTRY',
			minWidth: 150,
			flex: 1,
			renderCell: (params: any) => (
				<>
					<span
						style={{ marginRight: "5px" }}
						className={`fi fi-${params.row.country_code.toLowerCase()}`}>
					</span>
					{params.row.country_name}
				</>
			),
		},
		{
			field: 'protocol',
			headerName: 'PROTOCOL',
			minWidth: 100,
			flex: 1,
			renderCell: (params: any) => (
				<ProtocolTab
					label={params.row.protocol}
					sx={{
						backgroundColor: theme.palette.grayTabProtocol,
						color: theme.palette.grayTextProtocol,
					}}
				/>
			),
		},
		{
			field: 'response_time',
			headerName: 'RESPONSE',
			minWidth: 100, flex: 1,
			renderCell: (params: any) => (
				<>
					{params.row.response_time}ms
				</>
			)
		},
		{
			field: 'updated_at_formatted',
			headerName: 'UPDATED',
			minWidth: 150,
			flex: 1,
			renderCell: (params: any) => {
				const timeDuration = new Date().getTime() - params.row.updated_at_formatted.getTime();
				let timeFormatted: string;

				if (timeDuration < 60 * 1000) {
					timeFormatted = `${Math.floor(timeDuration / 1000)}s`;
				} else if (timeDuration < 60 * 60 * 1000) {
					timeFormatted = `${Math.floor(timeDuration / (60 * 1000))}m`;
				} else if (timeDuration < 24 * 60 * 60 * 1000) {
					timeFormatted = `${Math.floor(timeDuration / (60 * 60 * 1000))}h`;
				} else {
					timeFormatted = `${Math.floor(timeDuration / (24 * 60 * 60 * 1000))}d`;
				}
				return (
					<>
						{timeFormatted} ago
					</>
				)
			}
		},
		{
			field: 'status',
			headerName: 'AVAILABLE',
			minWidth: 150,
			flex: 1,
			renderCell: (params: any) => (
				<>
					<ProtocolTab
						label={params.row.status}
						sx={{
							backgroundColor: availableBgColor(params.row.status),
							color: availableTextColor(params.row.status),
						}}
					/>
				</>
			)
		},
		{
			field: 'isFavorite',
			headerName: 'FAVORITE',
			minWidth: 100,
			flex: 0.5,
			renderCell: (params: any) => (
				<Box
					sx={{
						display: "flex",
						flexDirection: "row",
						justifyContent: "right",
						alignItems: "center",
					}}>
					{params.row.isFavorite ? (
						<IoStar
							onClick={() => {
								// console.log(params.row);
								favoriteHandler(params.row.id, params.row.isFavorite);
							}}
							style={{
								width: "40px",
								height: "40px",
								margin: "0px 20px",
								marginLeft: "auto",
								padding: "10px 10px",
								color: "#FFD700",
								cursor: "pointer",
							}} />
					) : (
						<IoStarOutline
							onClick={() => {
								// console.log(params.row);
								favoriteHandler(params.row.id, params.row.isFavorite);
							}}
							style={{
								width: "40px",
								height: "40px",
								margin: "0px 20px",
								marginLeft: "auto",
								padding: "10px 10px",
								color: theme.palette.textGray,
								cursor: "pointer",
							}} />
					)}
				</Box>
			)
		}

	]
	if (!proxies) {
		return <></>;
	}
	return (
		<DataGrid
			rows={proxies.map((proxy, _) => ({
				...proxy,
				isFavorite: favorits?.find((favorit) => favorit?.proxy_id === proxy?.id) ? true : false,
			}))}
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
				},
				'& .MuiDataGrid-cell': {
					border: "none",
					outline: "none !important",
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
}
