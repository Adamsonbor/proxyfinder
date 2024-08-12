import React from "react"
import { useConfig } from "../../config";

interface IApiData {
	status: string
	data: any,
	error: any,
}

export const useApiV2 = (url: string, authToken?: string) => {
	const config = useConfig();
	const [data, setData] = React.useState<IApiData>({ status: "Loading", error: null, data: null });

	const fetchData = async () => {
		try {
			const response = await fetch(`${config.serverUrl}/api/v2${url}`, {
				headers: {
					"Authorization" : `Bearer ${authToken}`,
				},
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
