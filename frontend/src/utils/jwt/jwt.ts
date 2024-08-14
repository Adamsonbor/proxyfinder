import { Config } from "../../config";
import { IApiData } from "../api/api";

interface RefreshResponse {
	access_token: string,
	refresh_token: string
	expires_in: number
}

export class Jwt {
	private _config: Config;

	constructor(config: Config) {
		this._config = config;
	}

	// Get new access token with refresh tokena using refresh token
	public async updateTokens(): Promise<void> {
		const refresh_token = this.refresh_token;
		if (refresh_token === "") {
			return Promise.reject("No refresh token");
		}

		await fetch(`${this._config.server.url}/auth/google/refresh?refresh_token=${refresh_token}`)
			.then((res) => {
				if (res.ok) {
					return res.json()
				}
				return Promise.reject("No refresh token");
			})
			.then((data: IApiData) => data.data as RefreshResponse)
			.then((data: RefreshResponse) => {
				// console.log(`data: ${JSON.stringify(data)}`)
				this.setCookie("access_token", data.access_token, data.expires_in);
				this.setCookie("refresh_token", data.refresh_token, 1 * 24 * 60 * 60);
			})
	}

	// Get access token if it exists
	// Else get access token with refresh token
	// Else return empty string
	public async getAccessToken(): Promise<string> {
		let access_token = this.access_token;
		if (access_token !== "") {
			// console.log("Found access token: ", access_token);
			return access_token;
		}

		console.log("No access token, refreshing...");
		return await this.updateTokens()
		.then(() => {
			return this.access_token
		})
	}

	public get access_token(): string {
		return this.getCookie("access_token");
	}

	public get refresh_token(): string {
		return this.getCookie("refresh_token");
	}
	// Split cookie into name and value
	// If name matches, return value
	// Else return empty string
	private getCookie(name: string): string {
		const value = `; ${document.cookie}`;
		const parts = value.split(`; ${name}=`);
		if (parts.length === 2) return parts.pop()!.split(';').shift() || "";
		else return "";
	}

	private setCookie(name: string, value: string, expireSeconds: number = 3600) {
		const date = new Date(Date.now() + expireSeconds * 1000);
		document.cookie = `${name}=${value}; expires=${date.toUTCString()}; path=/`;
	}
}
