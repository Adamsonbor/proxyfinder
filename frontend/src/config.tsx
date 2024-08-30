import React from "react";

export interface Config {
	server: Server
}

interface Server {
	limit: number;
	url: string;
	apiUrl: string;
}

export const ConfigContext = React.createContext<Config>({} as Config);

export function ConfigProvider(props: any) {
	let defaultConfig: Config;

	if (import.meta.env.MODE === "development") {
		defaultConfig = {
			server: {
				limit: 20,
				url: "http://localhost:8080",
				apiUrl: "http://localhost:8080/api/v1",
			}
		}
	} else {
		defaultConfig = {
			server: {
				limit: 20,
				url: import.meta.env.VITE_API_URL,
				apiUrl: import.meta.env.VITE_API_URL + "/api/v1",
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

