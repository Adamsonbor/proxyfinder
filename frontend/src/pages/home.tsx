import "/node_modules/flag-icons/css/flag-icons.min.css";
import { useEffect, useState } from 'react'
import { Container, Snackbar, ThemeProvider } from '@mui/material'
import LeftPanel from '../components/LeftPanel/LeftPanel';
import { useConfig } from '../config';
import { User } from '../types';
import { lightTheme } from '../theme';
import Header from '../components/Header/Header';
import Table from "../components/Table/Table";

export default function HomePage() {
	const config = useConfig();

	const [theme, setTheme] = useState(lightTheme);
	const [user, setUser] = useState<User | undefined>(undefined);
	const [openNotification, setOpenNotification] = useState(false);
	const [sorts, setSorts] = useState<object>({});
	const [filter, setFilter] = useState<object>({})

	useEffect(() => {
		setFilter({
			"page": 1,
			"perPage": config.server.limit
		})
	}, []);

	console.log("render")

	const body = document.getElementsByTagName('body')[0];
	body.style.backgroundColor = theme.palette.backgroundWhite;

	return (
		<>
			<ThemeProvider theme={theme}>
				<Header
					user={user}
					setUser={setUser}
					setModalOpen={() => { }} />
				<Container maxWidth="xl" sx={{ color: theme.palette.textBlack }}>
					<div style={{ display: 'flex', flexDirection: 'column', gap: '10px', marginTop: '10px' }} >
						<div className="row">
							<LeftPanel
								className="col-2"
								filter={filter}
								setFilter={setFilter}
								theme={theme}
								setTheme={setTheme} />
							<Table
								sx={{
									'& .MuiButtonBase-root': {
										color: theme.palette.textGray,
										fontSize: theme.typography.uppercaseSize,
										fontWeight: theme.typography.fontWeightMedium,

									},
									height: '93vh',
								}}
								className="col-10"
								user={user}
								sorts={sorts}
								setSorts={setSorts}
								filter={filter}
								setFilter={setFilter} />
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
}
