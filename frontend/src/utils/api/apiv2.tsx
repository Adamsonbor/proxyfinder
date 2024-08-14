import React from "react"
import { useConfig } from "../../config";
import { Jwt } from "../jwt/jwt";

interface IApiData {
	status: string
	data: any,
	error: any,
}

export const useApiV2 = (url: string, auth?: boolean) => {
	const config = useConfig();
	const jwt = new Jwt(config);
	const [data, setData] = React.useState<IApiData>({ status: "Loading", error: null, data: null });

	const fetchData = async () => {
		try {
			const authToken: string = auth ? await jwt.getAccessToken() : "";
			const response = await fetch(`${config.serverUrl}/api/v2${url}`, {
				headers: {
					"Authorization": `Bearer ${authToken}`,
				}
			});
			const json = await response.json();
			setData(json);
		} catch (error) {
			setData({ status: "Error", error: error, data: null });
		}
	};

	React.useEffect(() => {
		fetchData();
	}, []);

	return data;
}
