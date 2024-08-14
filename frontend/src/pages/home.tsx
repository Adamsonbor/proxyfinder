import "/node_modules/flag-icons/css/flag-icons.min.css";
import { useEffect, useState } from 'react'
import { Container, ThemeProvider } from '@mui/material'
import InfiniteTable from '../components/Table/InfiniteTable';
import LeftPanel from '../components/LeftPanel/LeftPanel';
import { useConfig } from '../config';
import { Favorits, ProxyRow, ProxyV2, User } from '../types';
import { lightTheme } from '../theme';
import Header from '../components/Header/Header';
import { FavoritsRepo } from "../repos/favorits/repo";
import { ProxyV2Repo } from "../repos/proxy/repo";
import { UserRepo } from "../repos/user/repo";
import { IApiData } from "../utils/api/api";

export default function HomePage() {
	const config = useConfig();
	const userRepo = new UserRepo(config);
	const proxyRepo = new ProxyV2Repo(config);
	const favoritsRepo = new FavoritsRepo(config);

	const [theme, setTheme] = useState(lightTheme);
	const [proxies, setProxies] = useState<ProxyRow[]>([]);
	const [fullProxies, setFullProxies] = useState<ProxyRow[]>([]);
	const [favorits, setFavorits] = useState<Favorits[]>([])
	const [user, setUser] = useState<User | undefined>(undefined);

	useEffect(() => {
		proxyRepo.GetAll().then((proxies) => {
			setProxies(proxies.map(proxyV2ToProxyRow));
			setFullProxies(proxies.map(proxyV2ToProxyRow));
		})

		favoritsRepo.GetAll().then((favorits) => {
			setFavorits(favorits);
		})

		userRepo.Get().then((user) => {
			setUser(user);
		})
	}, []);

	const body = document.getElementsByTagName('body')[0];
	body.style.backgroundColor = theme.palette.backgroundWhite;

	return (
		<>
			<ThemeProvider theme={theme}>
				<Container maxWidth="xl" sx={{ color: theme.palette.textBlack }}>
					<div style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginTop: '10px' }} >
						<Header
							user={user}
							setModalOpen={() => { }} />
						<div className="row pt-5">
							<LeftPanel
								proxies={fullProxies}
								setProxies={setProxies}
								theme={theme}
								setTheme={setTheme}
								className="col-2" />
							<InfiniteTable
								proxies={proxies}
								favorits={favorits}
								favoriteHandler={favoriteHandler}
								className="col-10" />
						</div>
					</div>
				</Container>
			</ThemeProvider>
		</>
	);

	// 
	async function favoriteHandler(proxyId: number, isFavorite: boolean) {
		if (!user) {
			await userRepo.Get().then(user => setUser(user))
			if (!user) { console.error("Error getting user from server"); return }
		}
		if (isFavorite) {
			await favoritsRepo.Delete(proxyId);
			setFavorits(favorits.filter((favorit) => favorit.proxy_id !== proxyId));
		} else {
			await favoritsRepo.Create(user!.id, proxyId)
				.then((data: IApiData) => setFavorits([...favorits, data.data]))
				.catch((err) => console.log(err))
		}
	}
}

function proxyV2ToProxyRow(proxy: ProxyV2): ProxyRow {
	return {
		id: proxy.id,
		ip: proxy.ip,
		port: proxy.port,
		protocol: proxy.protocol,
		response_time: proxy.response_time,
		status_id: proxy.status.id,
		country_id: proxy.country.id,
		created_at: proxy.created_at,
		updated_at: proxy.updated_at,
		status: proxy.status.name,
		country_code: proxy.country.code,
		country_name: proxy.country.name,
		created_at_formatted: new Date(proxy.created_at),
		updated_at_formatted: new Date(proxy.updated_at),
	}
}
