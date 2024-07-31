import React from "react";

interface Config {
	apiUrl: string;
}

export const ConfigContext = React.createContext<Config>({} as Config);

export function ConfigProvider(props: any) {
	const defaultConfig: Config = {
		apiUrl: import.meta.env.VITE_API_URL + "/api/v1",
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

