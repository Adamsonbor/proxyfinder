import './App.css';
import "/node_modules/flag-icons/css/flag-icons.min.css";
import { darkTheme, lightTheme } from './theme';
import ThemeSwitch from './components/ThemeSwitch/ThemeSwitch';

import React from 'react'
import { Container, ThemeProvider } from '@mui/material'
import InfiniteTable from './components/Table/InfiniteTable';
import LeftPanel from './components/LeftPanel/LeftPanel';
import { ConfigProvider } from './config';

export default function App() {
	const [theme, setTheme] = React.useState(lightTheme)
	const toggleTheme = () => {
		setTheme(theme === lightTheme ? darkTheme : lightTheme);
	};
	return (
		<>
			<ConfigProvider>
				<ThemeProvider theme={theme} >
					<Container maxWidth="xl">
						<div style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginTop: '10px' }} >
							<ThemeSwitch onChange={toggleTheme} className="position-fixed bottom-0 mb-5" />
							<div
								style={{
									backgroundColor: theme.palette.backgroundWhite,
									width: "100%",
									position: "fixed",
									top: 0,
									left: 0,
									zIndex: 100,
									height: 50,
									display: "flex",
									justifyContent: "center",
									alignItems: "center",
									boxShadow: "0px 4px 4px 0px rgba(0, 0, 0, 0.05)",
								}}
								className="row">
								<img src="proxpro-day.svg" height="18px" className="App-logo" alt="logo" />
							</div>
							<div className="row pt-5">
								<LeftPanel className="col-2" />
								<InfiniteTable className="col-10" />
							</div>
						</div>
					</Container>
				</ThemeProvider>
			</ConfigProvider>
		</>
	)
}
