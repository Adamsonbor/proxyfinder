import { Container, Button, useTheme, Box, Menu, MenuItem } from "@mui/material";
import { FcGoogle } from "react-icons/fc";
import { SiMaildotru } from "react-icons/si";
import { FaGithub } from "react-icons/fa";
import { User } from "../../types";
import UserBlock from "./UserBlock";
import { useConfig } from "../../config";
import { SetCookie } from "../../utils/utils";
import { useEffect, useRef, useState } from "react";
import { UserRepo } from "../../repos/user/repo";

interface Props {
	// User or undefined
	user?: User
	setUser?: (user: User | undefined) => void
	setModalOpen?: any
}

export default function Header({
	user = undefined,
	setUser = () => { },
}: Props) {
	const config = useConfig();
	const theme = useTheme();
	const userRepo = new UserRepo(config);

	const [openUserMenu, setOpenUserMenu] = useState(false);
	const menuButton = useRef(null);

	useEffect(() => {
		if (user === undefined) {
			userRepo.GetBy({}, {}, true).then((user) => {
				if (user !== undefined) {
					setUser(user);
				}
			});
		}
	}, [])

	function login() {
		// redirect
		window.location.href = `${config.server.url}/auth/google/login`;
	}

	return (
		<div
			style={{
				backgroundColor: theme.palette.backgroundWhite,
				width: "100%",
				top: 0,
				left: 0,
				zIndex: 100,
				height: "50px",
				display: "flex",
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
					gap: "10px",
					display: "grid",
					gridTemplateColumns: "1fr 1fr 1fr",
					alignItems: "center",
				}}>
				<div></div>
				<img
					src={theme.palette.mode === "dark" ? "proxpro-night.svg" : "proxpro-day.svg"}
					style={{ margin: "auto" }}
					height="18px"
					className="App-logo"
					alt="logo" />
				<Box
					sx={{
						marginLeft: "auto",
						marginRight: "40px",
					}}>
					{!user && <Box
						sx={{
							color: theme.palette.textBlack,
							display: "flex",
							gap: "4px",
							alignItems: "center",
						}}>
						<Button
							variant="text"
							sx={{
								minWidth: "40px",
							}}
							onClick={() => {
								login()
							}}>
							<FcGoogle width={18} size={18} />
						</Button>
						<Button
							sx={{
								minWidth: "40px",
							}}
							disabled>
							<SiMaildotru size={18} />
						</Button>
						<Button
							sx={{
								minWidth: "40px",
							}}
							disabled>
							<FaGithub size={18} />
						</Button>
					</Box>}
					{user &&
						<Box>
							<Button
								sx={{
									textTransform: "none",
									color: theme.palette.textGray,
								}}
								ref={menuButton}
								id="basic-button"
								variant="text"
								aria-controls={openUserMenu ? 'basic-menu' : undefined}
								aria-haspopup="true"
								aria-expanded={openUserMenu ? 'true' : undefined}
								onClick={() => setOpenUserMenu(!openUserMenu)} >
								<UserBlock user={user} />
							</Button>
							<Menu
								anchorEl={menuButton.current}
								open={openUserMenu}
								onClose={() => setOpenUserMenu(false)}
								MenuListProps={{
									'aria-labelledby': 'basic-button',
								}} >
								<MenuItem
									onClick={() => {
										SetCookie("access_token", "", 0);
										SetCookie("refresh_token", "", 0);
										window.location.reload();
									}} >
									Logout
								</MenuItem>

							</Menu>
						</Box>
					}
				</Box>
			</Container>
		</div>
	)
}
