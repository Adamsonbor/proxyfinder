import { Box, useTheme } from "@mui/material";
import TableHeader from "./TableHeader";
import { TableRow } from "./TableRow";
import { CellParams, Column, Value } from "./types";
import { useEffect, useRef, useState } from "react";
import ProtocolTab from "../ProtocolTab/ProtocolTab";
import { IoStar, IoStarOutline } from "react-icons/io5";
import { Favorits, Proxy, User } from "../../types";
import { useQuery } from "../../repos/utils";
import { useConfig } from "../../config";
import { FavoritsRepo } from "../../repos/favorits/repo";
import { ProxyRepo } from "../../repos/proxy/repo";

interface Props {
	className?: string;
	sorts?: object;
	setSorts?: (sorts: object) => void;
	filter?: object;
	setFilter?: (filter: object) => void;
	setNotification?: (open: boolean) => void;
	user?: User;
	sx?: object;
}

function Table({
	className = '',
	sorts = {},
	setSorts = () => { },
	filter = {},
	setFilter = () => { },
	setNotification = () => { },
	user = undefined,
	sx = {},
}: Props) {
	const theme = useTheme();
	const tableBody = useRef<HTMLDivElement>(null);
	const config = useConfig();
	const proxyRepo = new ProxyRepo(config);
	const favoritsRepo = new FavoritsRepo(config);

	const [proxies, setProxies] = useState<Proxy[]>([]);
	const [favorits, setFavorits] = useState<Favorits[]>([]);
	const [values, setValues] = useState<Value[]>([]);
	const queryBuilder = useQuery();
	const favQB = useQuery();

	useEffect(() => {
		loadData()
	}, [sorts, filter]);

	useEffect(() => {
		setValues(proxies.map((proxy) => ({ ...proxy, isFavorite: favorits.some((favorit) => favorit.proxy_id === proxy.id) })))
	}, [proxies]);

	useEffect(() => {
		const handleScroll = (e: Event) => {
			const { scrollTop, scrollHeight, clientHeight } = e.target as HTMLDivElement;
			if (scrollTop + clientHeight == scrollHeight) {
				onScroll();
			}
		};

		tableBody.current?.addEventListener("scroll", handleScroll);

		return () => {
			tableBody.current?.removeEventListener("scroll", handleScroll); }
	}, [onScroll]);

	const columns: Column<Value>[] = [
		{
			field: 'ip',
			headerName: 'IP',
			minWidth: 100,
			flex: 1,
			renderCell: (params: CellParams<Value>) => (
				<>
					{params.row.ip}
				</>
			),
		},
		{
			field: 'port',
			headerName: 'PORT',
			minWidth: 100,
			flex: 1,
			renderCell: (params: CellParams<Value>) => (
				<>
					{params.row.port}
				</>
			),
		},
		{
			field: 'country_name',
			headerName: 'COUNTRY',
			minWidth: 150,
			flex: 1,
			renderCell: (params: CellParams<Value>) => (
				<>
					<span
						style={{ marginRight: "5px", minWidth: "20px" }}
						className={`fi fi-${params.row.country.code.toLowerCase()}`}>
					</span>
					{params.row.country.name}
				</>
			),
		},
		{
			field: 'protocol',
			headerName: 'PROTOCOL',
			minWidth: 100,
			flex: 1,
			renderCell: (params: CellParams<Value>) => (
				<ProtocolTab
					label={params.row.protocol.toUpperCase()}
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
			renderCell: (params: CellParams<Value>) => (
				<>
					{params.row.response_time}ms
				</>
			),
		},
		{
			field: 'updated_at',
			headerName: 'UPDATED',
			minWidth: 150,
			flex: 1,
			renderCell: (params: CellParams<Value>) => {
				const timeDuration = new Date().getTime() - new Date(params.row.updated_at).getTime();
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
			field: 'status_id',
			headerName: 'AVAILABLE',
			minWidth: 100,
			flex: 1,
			renderCell: (params: CellParams<Value>) => (
				<>
					<ProtocolTab
						label={params.row.status.name}
						sx={{
							backgroundColor: availableBgColor(params.row.status.name),
							color: availableTextColor(params.row.status.name),
						}}
					/>
				</>
			),
		},
		{
			field: 'favorits',
			headerName: 'FAVORITE',
			minWidth: 100,
			flex: 1,
			renderCell: (params: CellParams<Value>) => (
				<Box
					sx={{
						display: "flex",
						flexDirection: "row",
						alignItems: "center",
					}}>
					{params.row.isFavorite ? (
						<IoStar
							onClick={() => {
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

	return (
		<Box
			className={className}
			sx={{
				overflowX: "auto",
				...sx
			}}>
			<Box
				sx={{
					minWidth: "100%",
					width: "fit-content",
				}}>
				<TableHeader
					sx={{
						height: "52px",
					}}
					sorts={sorts}
					onHeaderClick={onHeaderClick}
					columns={columns} />
				<Box
					ref={tableBody}
					sx={{
						height: "84vh",
						overflowY: "auto",
					}}>
					{
						values.map((value: Value, index: number) => (
							<TableRow<Value>
								sx={{
									borderTop: "1px solid rgba(0, 0, 0, 0.1)",
									height: "52px",
								}}
								key={index}
								columns={columns}
								value={value} />
						))

					}
				</Box>
			</Box>
		</Box>
	)

	async function loadData(): Promise<void> {
		let page = filter["page"]
		if (page === undefined || page < 1) {
			filter["page"] = 1
		}

		// get proxy
		queryBuilder.setSorts(sorts);
		queryBuilder.setFilters(filter);
		queryBuilder.setFilter("page", page);
		proxyRepo.GetAll(queryBuilder).then((res) => {
			if (page === 1) {
				setProxies([...res]);
			} else {
				setProxies([...proxies, ...res]);
			}
		})

		// get favorits
		favQB.setFilter("page", page);
		favQB.setFilter("perPage", config.server.limit);
		favoritsRepo.GetAll(favQB, true).then((favoritsRes) => {
			if (page == 1) {
				setFavorits(favoritsRes);
			} else {
				setFavorits([...favorits, ...favoritsRes]);
			}
		})
	}

	async function favoriteHandler(proxyId: number, isFavorite: boolean) {
		if (!user) {
			setNotification(true);
			return
		}

		if (isFavorite) {
			await favoritsRepo.Delete(proxyId)
				.then(() => setFavorits(favorits.filter((favorit) => favorit.proxy_id !== proxyId)))
		} else {
			const newFav = {
				user_id: user.id,
				proxy_id: proxyId
			} as Favorits
			await favoritsRepo.Create(newFav, true)
				.then(() => setFavorits([...favorits, newFav]))
				.catch((err) => console.error(err))
		}

		setFilter({ ...filter, page: 1});
	}

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

	function onHeaderClick(col: Column<Value>) {
		if (!col) {
			return
		}

		if (!sorts[col.field]) {
			setSorts({ ...sorts, [col.field]: "asc" })
		} else if (sorts[col.field] === "asc") {
			setSorts({ ...sorts, [col.field]: "desc" })
		} else if (sorts[col.field] === "desc") {
			delete sorts[col.field]
			setSorts({ ...sorts })
		}

		setFilter({ ...filter, page: 1 });
	}

	function onScroll() {
		setFilter({ ...filter, page: filter["page"]? filter["page"] + 1: 1 });
	}
}

export default Table
