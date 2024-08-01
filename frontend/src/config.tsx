import React from "react";

interface Config {
	apiUrl: string;
}

export const ConfigContext = React.createContext<Config>({} as Config);

export function ConfigProvider(props: any) {
	let defaultConfig: Config;

	if (import.meta.env.MODE === "development") {
		defaultConfig = {
			apiUrl: "http://localhost:8080/api/v1",
		}
	} else {
		defaultConfig = {
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

