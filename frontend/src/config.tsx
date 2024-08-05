import React from "react";

export interface Config {
	apiUrl: string;
	rabbitApi: string;
}

export const ConfigContext = React.createContext<Config>({} as Config);

export function ConfigProvider(props: any) {
	let defaultConfig: Config;

	if (import.meta.env.MODE === "development") {
		defaultConfig = {
			rabbitApi: "http://localhost:8080/rabbit",
			apiUrl: "http://localhost:8080/api/v1",
		}
	} else {
		defaultConfig = {
			rabbitApi: import.meta.env.VITE_API_URL + "/rabbit",
			apiUrl: import.meta.env.VITE_API_URL + "/api/v1",
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

