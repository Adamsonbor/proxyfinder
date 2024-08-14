import React from "react";

export interface Config {
	apiUrl: string;
	serverUrl: string;
	server: Server
}

interface Server {
	url: string;
	apiUrl: string;
	apiV2Url: string;
}

export const ConfigContext = React.createContext<Config>({} as Config);

export function ConfigProvider(props: any) {
	let defaultConfig: Config;

	if (import.meta.env.MODE === "development") {
		defaultConfig = {
			apiUrl: "http://localhost:8080/api/v1",
			serverUrl: "http://localhost:8080",
			server: {
				url: "http://localhost:8080",
				apiUrl: "http://localhost:8080/api/v1",
				apiV2Url: "http://localhost:8080/api/v2",
			}
		}
	} else {
		defaultConfig = {
			apiUrl: import.meta.env.VITE_API_URL + "/api/v1",
			serverUrl: import.meta.env.VITE_API_URL,
			server: {
				url: import.meta.env.VITE_API_URL,
				apiUrl: import.meta.env.VITE_API_URL + "/api/v1",
				apiV2Url: import.meta.env.VITE_API_URL + "/api/v2",
			}
		}
	}

	return (
		<ConfigContext.Provider
			value={defaultConfig}>
			{props.children}
		</ConfigContext.Provider>
	);
}

export function useConfig(): Config {
	return React.useContext(ConfigContext);
}

