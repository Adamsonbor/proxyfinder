import "/node_modules/flag-icons/css/flag-icons.min.css";
import { useEffect, useState } from 'react'
import { Box, Container, Snackbar, ThemeProvider } from '@mui/material'
import LeftPanel from '../components/LeftPanel/LeftPanel';
import { useConfig } from '../config';
import { Favorits, Proxy, User } from '../types';
import { lightTheme } from '../theme';
import Header from '../components/Header/Header';
import { FavoritsRepo } from "../repos/favorits/repo";
import { ProxyRepo } from "../repos/proxy/repo";
import { UserRepo } from "../repos/user/repo";
import Table from "../components/Table/Table";
import { CellParams, Column } from "../components/Table/types";
import ProtocolTab from "../components/ProtocolTab/ProtocolTab";
import { IoStar, IoStarOutline } from "react-icons/io5";
import { useQuery } from "../repos/utils";

interface Value extends Proxy {
	isFavorite: boolean;
}

export default function HomePage() {
	const config = useConfig();
	const userRepo = new UserRepo(config);
	const favoritsRepo = new FavoritsRepo(config);
	const proxyRepo = new ProxyRepo(config);
	const queryBuilder = useQuery();
	const favQB = useQuery();

	const [theme, setTheme] = useState(lightTheme);
	const [proxies, setProxies] = useState<Proxy[]>([]);
	const [values, setValues] = useState<Value[]>([]);
	const [favorits, setFavorits] = useState<Favorits[]>([])
	const [user, setUser] = useState<User | undefined>(undefined);
	const [openNotification, setOpenNotification] = useState(false);
	const [page, setPage] = useState<number>(1);
	const [sorts, setSorts] = useState({
		"favorits": "desc"
	});
	const [filter, setFilter] = useState<object>({
		"page": page,
		"perPage": config.server.limit
	});

	useEffect(() => {
		loadData()
	}, [page, sorts, filter]);

	useEffect(() => {
		setValues(proxies.map((proxy) => ({ ...proxy, isFavorite: favorits.some((favorit) => favorit.proxy_id === proxy.id) })))
	}, [proxies, favorits]);


	const body = document.getElementsByTagName('body')[0];
	body.style.backgroundColor = theme.palette.backgroundWhite;

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
						style={{ marginRight: "5px" }}
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
		<>
			<ThemeProvider theme={theme}>
				<Header
					user={user}
					setModalOpen={() => { }} />
				<Container maxWidth="xl" sx={{ color: theme.palette.textBlack }}>
					<div style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginTop: '10px' }} >
						<div className="row">
							<LeftPanel
								className="col-2"
								theme={theme}
								countries={proxies?.map((proxy) => proxy.country.name) || []}
								filter={filter}
								setFilter={setFilter}
								setTheme={setTheme} />
							<Table<Value>
								sx={{
									'& .MuiButtonBase-root': {
										color: theme.palette.textGray,
										fontSize: theme.typography.uppercaseSize,
										fontWeight: theme.typography.fontWeightMedium,

									},
									height: '93vh',
								}}
								bodyStyle={{
									height: '93%',
								}}
								className="col-10"
								values={values}
								sorts={{ ...sorts }}
								onScroll={onScroll}
								onHeaderClick={onHeaderClick}
								columns={columns} />
						</div>
					</div>
					<Snackbar
						anchorOrigin={{ vertical: 'top', horizontal: 'center' }}
						autoHideDuration={3000}
						open={openNotification}
						onClose={() => setOpenNotification(false)}
						message="Login required" />
				</Container>
			</ThemeProvider>
		</>
	);

	function loadData() {
		if (page < 1) {
			return;
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

		// get user
		userRepo.GetBy({}, {}, true).then((userRes) => {
			if (userRes) {
				setUser(userRes);
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

	function onScroll() {
		setPage(page + 1);
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

		setPage(1)
	}

	async function favoriteHandler(proxyId: number, isFavorite: boolean) {
		if (!user) {
			await userRepo.GetBy({}, {}, true)
				.then(user => setUser(user))
				.catch((err) => {
					setOpenNotification(true)
					console.log(err)
				})
			if (!user) { return }

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
}
