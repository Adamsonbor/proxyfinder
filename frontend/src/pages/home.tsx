import "/node_modules/flag-icons/css/flag-icons.min.css";
import { useEffect, useState } from 'react'
import { Container, ThemeProvider } from '@mui/material'
import InfiniteTable from '../components/Table/InfiniteTable';
import LeftPanel from '../components/LeftPanel/LeftPanel';
import { ConfigProvider, useConfig } from '../config';
import { Country, Favorits, ProxyRow, ProxyV2, User } from '../types';
import { lightTheme } from '../theme';
import Header from '../components/Header/Header';
import { useApiV2 } from '../utils/api/apiv2';
import { GetCookie } from "../utils/utils";
import { IApiData } from "../utils/api/api";

export default function HomePage() {
	const config = useConfig();

	const [theme, setTheme] = useState(lightTheme);
	const [proxies, setProxies] = useState<ProxyRow[]>([]);
	const [fullProxies, setFullProxies] = useState<ProxyRow[]>([]);
	const [favorits, setFavorits] = useState<Favorits[]>([])
	const [user, setUser] = useState<User | undefined>(undefined);

	let countries: Country[] = []
	let proxiesData: ProxyV2[] = useApiV2(`/proxy?perPage=7000`).data;
	let favoritsData: Favorits[] = useApiV2(`/favorits`).data;

	useEffect(() => {
		if (proxiesData) {
			setProxies(proxiesData.map(proxyV2ToProxyRow));
			setFullProxies(proxiesData.map(proxyV2ToProxyRow));
			setFavorits(favoritsData);
			UserInfo();
			countries = proxiesData.map(proxy => proxy.country);
		}
	}, [proxiesData]);

	const body = document.getElementsByTagName('body')[0];
	body.style.backgroundColor = theme.palette.backgroundWhite;

	return (
		<>
			<ConfigProvider>
				<ThemeProvider theme={theme} >
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
									favoritHandler={(proxy_id, isFavorite) => { }}
									className="col-10" />
							</div>
						</div>
					</Container>
				</ThemeProvider>
			</ConfigProvider>
		</>
	);

	function UserInfo() {
		const access_token = GetCookie("access_token");

		if (!access_token) {
			return
		}

		fetch(
			`${config.server.apiV2Url}/user`,
			{
				method: "GET",
				headers: {
					"Authorization": `Bearer ${access_token}`,
				}
			})
			.then(res => res.json())
			.then((apiData: IApiData) => {
				if (apiData.data) {
					setUser(apiData.data)
				} else {
					console.log(apiData.error)
				}
			})
			.catch(err => console.log(err));
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
}
