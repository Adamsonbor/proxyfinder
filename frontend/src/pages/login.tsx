import { GetCookie, SetCookie } from "../utils/utils";
import { useConfig } from "../config";
import { useEffect, useState } from "react";
import { User } from "../types";
import { IApiData } from "../utils/api/api";

export default function LoginPage() {
	const config = useConfig();
	const [user, setUser] = useState<User | null>(null)

	useEffect(() => {
		UserInfo();
		if (user) {
			SetCookie("user", JSON.stringify(user), 100000)
		}
	}, []);

	return (
		<div>
			<h1>Login</h1>
		</div>
	);

	function UserInfo() {
		const access_token = GetCookie("access_token");

		fetch(
			`${config.server.apiUrl}/user`,
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
}
