import { Avatar, Box } from "@mui/material";
import { User } from "../../types";

interface Props {
	user?: User
}

export default function UserBlock({
	user,
}: Props) {

	console.log(user);

	return (
		<>
			<Box
				sx={{
					display: "flex",
					flexDirection: "row",
					alignItems: "center",
				}}>
				<Avatar
					sx={{
						width: 40,
						height: 40,
						marginRight: "10px",
					}}
					src={user?.photo_url}>
					{user?.name}
				</Avatar>
				<span>{user && user.name.charAt(0).toUpperCase() + user.name.slice(1)}</span>
			</Box>
		</>
	);
}
