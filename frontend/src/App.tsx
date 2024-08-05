import './App.css';
import "/node_modules/flag-icons/css/flag-icons.min.css";
import { useEffect, useState } from 'react'
import { Box, Button, Container, Input, Modal, ThemeProvider, Typography } from '@mui/material'
import InfiniteTable from './components/Table/InfiniteTable';
import LeftPanel from './components/LeftPanel/LeftPanel';
import { ConfigProvider, useConfig } from './config';
import { useApi } from './utils/api/api';
import { Country, Proxy, ProxyRow, Status } from './types';
import { lightTheme } from './theme';
import { RiMailSendLine } from "react-icons/ri";

export default function App() {
	const config = useConfig();

	const [theme, setTheme] = useState(lightTheme);
	const [proxies, setProxies] = useState<ProxyRow[]>([]);
	const [fullProxies, setFullProxies] = useState<ProxyRow[]>([]);
	const [modalOpen, setModalOpen] = useState(false);
	const [email, setEmail] = useState("");

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

	function sendEmail(email: string) {
		fetch(`${config.rabbitApi}/publish`, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
			},
			body: JSON.stringify({
				email: email
			})
		}).then((response) => {
			if (!response.ok) {
				throw new Error("Failed to send email");
			}
			console.log("Email sent successfully");
		}).catch((error) => {
			console.error(error);
		})
	}

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
									justifyContent: "space-between",
									alignItems: "center",
									boxShadow:
										theme.palette.mode === "dark" ?
											"0px 4px 4px 0px rgba(255, 255, 255, 0.05)" :
											"0px 4px 4px 0px rgba(0, 0, 0, 0.05)",
								}}
								className="column">
								<Container
									maxWidth="xl"
									sx={{
										display: "flex",
										gap: "10px",
										justifyContent: "space-between",
										alignItems: "center",
									}}>
									<div></div>
									<img
										src={theme.palette.mode === "dark" ? "proxpro-night.svg" : "proxpro-day.svg"}
										height="18px"
										className="App-logo"
										alt="logo" />
									<Button
										variant="text"
										sx={{
											padding: "8px 16px",
											marginRight: "40px",
											color: theme.palette.textBlack,
										}}
										onClick={() => {
											setModalOpen(true);
											console.log("test");
										}}>
										<RiMailSendLine size={24} />
									</Button>
									<Modal
										keepMounted
										open={modalOpen}
										onClose={() => { setModalOpen(false) }}
										aria-labelledby="modal-modal-title"
										aria-describedby="modal-modal-description" >
										<Box
											sx={{
												position: 'absolute',
												top: '50%',
												left: '50%',
												transform: 'translate(-50%, -50%)',
												width: '500px',
												bgcolor: theme.palette.backgroundWhite,
												borderRadius: "6px",
												boxShadow: 24,
												p: 4,
											}}>
											<Typography
												sx={{
													color: theme.palette.textBlack,
													fontWeight: 600,
												}}
												id="modal-modal-title"
												variant="h6"
												component="h2">
												Send proxy list!
											</Typography>
											<Typography
												id="modal-modal-description"
												sx={{
													mt: 2,
													color: theme.palette.textBlack,
													fontSize: "14px",
												}}>
												Enter your email and we will send you the proxy list to your email.
											</Typography>
											<Box
												sx={{
													display: "flex",
													flexDirection: "column",
													alignItems: "center",
													gap: "10px",
												}}>
												<Input
													value={email}
													onChange={(e) => { setEmail(e.target.value) }}
													sx={{
														mt: 2,
														width: "100%",
													}} />
												<Box
													sx={{
														mt: 2,
														display: "flex",
														alignItems: "left",
														width: "100%",
														gap: "10px",
													}}>
													<Button
														sx={{
															width: "fit-content",
															height: "42px",
															borderRadius: "21px",
															padding: '0px 40px',
															textTransform: "none",
															color: theme.palette.textGray,
															border: `1px solid ${theme.palette.textGray}`,
															"&:hover": {
																color: theme.palette.textBlack,
																border: `1px solid ${theme.palette.textBlack}`,
															}
														}}
														onClick={() => { setModalOpen(false) }}
														variant="outlined">
														Close
													</Button>
													<Button
														sx={{
															width: "fit-content",
															height: "42px",
															borderRadius: "21px",
															padding: '0px 40px',
															textTransform: "none",
														}}
														onClick={() => { sendEmail(email) }}
														variant="contained">
														Send proxy list
													</Button>
												</Box>
											</Box>
										</Box>
									</Modal>
								</Container>
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
	);

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
