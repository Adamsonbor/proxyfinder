import { Config } from "../../config";
import { IApiData } from "../api/api";
import useCookies from "universal-cookie";

interface RefreshResponse {
	access_token: string,
	refresh_token: string,
	expires_in: number,
	expires_in_ref: number
}

export class Jwt {
	private _config: Config;
	private _cookies: any;

	constructor(config: Config) {
		this._config = config;
		this._cookies = new useCookies();
	}

	// Get new access token with refresh tokena using refresh token
	public async updateTokens(): Promise<void> {
		const refresh_token = this.refresh_token;
		if (refresh_token === undefined) {
			return Promise.reject("No refresh token");
		}

		await fetch(`${this._config.server.url}/auth/google/refresh?refresh_token=${refresh_token}`)
			.then((res) => {
				if (!res.ok) {
					return Promise.reject("Invalid refresh token");
				}

				return res.json();
			})
			.then((res: IApiData) => res.data)
			.then((res: RefreshResponse) => {
				console.log(`RefreshResponse: ${JSON.stringify(res)}`);
				this._cookies.set("access_token", res.access_token, { expires: new Date(res.expires_in * 1000) });
				this._cookies.set("refresh_token", res.refresh_token, { expires: new Date(res.expires_in_ref * 1000) });
			})
	}

	public get access_token() {
		return this._cookies.get("access_token");
	}

	public get refresh_token() {
		return this._cookies.get("refresh_token");
	}

	// Get access token if it exists
	// Else get access token with refresh token
	// Else return empty string
	public async getAccessToken(): Promise<string> {
		let access_token = this.access_token
		if (access_token != undefined && access_token.length > 0) {
			console.log("Found access token", access_token);
			return access_token;
		}

		console.log("No access token, refreshing...");
		return await this.updateTokens()
			.then(() => {
				return this.access_token
			})
	}
}
