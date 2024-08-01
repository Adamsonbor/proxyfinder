import './App.css';
import "/node_modules/flag-icons/css/flag-icons.min.css";
import { useEffect, useState } from 'react'
import { Container, ThemeProvider } from '@mui/material'
import InfiniteTable from './components/Table/InfiniteTable';
import LeftPanel from './components/LeftPanel/LeftPanel';
import { ConfigProvider, useConfig } from './config';
import { useApi } from './utils/api/api';
import { Country, Proxy, ProxyRow, Status } from './types';
import { lightTheme } from './theme';

export default function App() {
	const config = useConfig();

	const [theme, setTheme] = useState(lightTheme);
	const [proxies, setProxies] = useState<ProxyRow[]>([]);
	const [fullProxies, setFullProxies] = useState<ProxyRow[]>([]);

	let countriesData: Country[] = useApi(`${config.apiUrl}/country`).data;
	let proxiesData: Proxy[] = useApi(`${config.apiUrl}/proxy`).data;
	let statusesData: Status[] = useApi(`${config.apiUrl}/status`).data;

	useEffect(() => {
		if (!proxiesData || !countriesData || !statusesData) {
			return
		}

		const out: ProxyRow[] = []

		for (const proxy of proxiesData) {
			out.push(proxyToProxyRow(proxy, countriesData, statusesData))
		}

		setProxies(out)
		setFullProxies(out)
	}, [proxiesData]);

	const body = document.getElementsByTagName('body')[0];
	body.style.backgroundColor = theme.palette.backgroundWhite;

	return (
		<>
			<ConfigProvider>
				<ThemeProvider theme={theme} >
					<Container maxWidth="xl" sx={{ color: theme.palette.textBlack }}>
						<div style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginTop: '10px' }} >
							<div
								style={{
									backgroundColor: theme.palette.backgroundWhite,
									width: "100%",
									position: "fixed",
									top: 0,
									left: 0,
									zIndex: 100,
									height: "50px",
									display: "flex",
									justifyContent: "center",
									alignItems: "center",
									boxShadow:
										theme.palette.mode === "dark" ?
										"0px 4px 4px 0px rgba(255, 255, 255, 0.05)":
										"0px 4px 4px 0px rgba(0, 0, 0, 0.05)",
								}}
								className="row">
								<img
									src={theme.palette.mode === "dark" ? "proxpro-night.svg" : "proxpro-day.svg"}
									height="18px"
									className="App-logo"
									alt="logo" />
							</div>
							<div className="row pt-5">
								<LeftPanel
									proxies={fullProxies}
									setProxies={setProxies}
									countries={countriesData}
									theme={theme}
									setTheme={setTheme}
									className="col-2" />
								<InfiniteTable
									proxies={proxies}
									countries={countriesData}
									className="col-10" />
							</div>
						</div>
					</Container>
				</ThemeProvider>
			</ConfigProvider>
		</>
	)

	function proxyToProxyRow(
		proxy: Proxy,
		countries: Country[],
		statuses: Status[],
	): ProxyRow {
		return {
			...proxy,
			status: statuses.find((status) => status.id === proxy.status_id)?.name || "Unknown",
			country_name: countries.find((country) => country.id === proxy.country_id)?.name || "Unknown",
			country_code: countries.find((country) => country.id === proxy.country_id)?.code || "Unknown",
			created_at_formatted: new Date(proxy.created_at),
			updated_at_formatted: new Date(proxy.updated_at),
		}
	}
}
