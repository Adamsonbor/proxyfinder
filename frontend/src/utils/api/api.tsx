import React from "react"

export enum ApiStatus {
	// API request is being made
	Loading,
	// API call was successful
	Success,
	// API call resulted in an unauthorized error even after attempting
	// a token refresh
	ErrorUnauthorized,
	// API resulted in an error
	Error,
	// The initial request failed and we are attempting to refresh an
	// access token
	RefreshingToken,
	// We have new access token and will attempt to make a request
	// again. Note: if the retry fails the status will be `Error`.
	Retrying,
}

export interface IApiData {
	status: ApiStatus
	error: any,
	data: any,
}

export const useApi = (url: string) => {
	const [data, setData] = React.useState<IApiData>({ status: ApiStatus.Loading, error: null, data: null });

	const fetchData = async () => {
		try {
			const response = await fetch(url);
			const json = await response.json();
			setData({ status: ApiStatus.Success, error: null, data: json });
		} catch (error) {
			setData({ status: ApiStatus.Error, error: error, data: null });
		}
	};

	React.useEffect(() => {
		fetchData();
	}, []);

	return data;
}
