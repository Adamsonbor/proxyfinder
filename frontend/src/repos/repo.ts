import { Config } from "../config";
import { IApiData } from "../utils/api/api"
import { Jwt } from "../utils/jwt/jwt";

export class Repo {
	_url: string;
	_jwt: Jwt;
	constructor(url: string, config: Config) {
		this._url = url
		this._jwt = new Jwt(config);
	}

	async GetAll(headers: object = {}, auth: boolean = false): Promise<IApiData> {
		headers = await this.authHeadersWrapper(headers, auth)
		return this.doRequest(this._url, {
			method: "GET",
			headers: {
				...headers,
			},
		})
	}

	async Get(params: string, headers: object = {}, auth: boolean = false): Promise<IApiData> {
		headers = await this.authHeadersWrapper(headers, auth)
		return this.doRequest(`${this._url}${params}`, {
			method: "GET",
			headers: {
				...headers,
			},
		})
	}

	async Create(data: object, headers: object = {}, auth: boolean = false): Promise<IApiData> {
		headers = await this.authHeadersWrapper(headers, auth)
		return this.doRequest(this._url, {
			method: "POST",
			headers: {
				"Content-Type": "application/json",
				...headers,
			},
			body: JSON.stringify(data),
		})
	}

	async Delete(params: string, headers: object = {}, auth: boolean = false): Promise<IApiData> {
		headers = await this.authHeadersWrapper(headers, auth)
		return this.doRequest(`${this._url}${params}`, {
			method: "DELETE",
			headers: {
				...headers,
			},
		})
	}

	async Update(params: string, data: object, headers: object = {}, auth: boolean = false): Promise<IApiData> {
		headers = await this.authHeadersWrapper(headers, auth)
		return this.doRequest(`${this._url}${params}`, {
			method: "PUT",
			headers: {
				"Content-Type": "application/json",
				...headers,
			},
			body: JSON.stringify(data),
		})
	}

	private async authHeadersWrapper(headers: object, auth: boolean): Promise<object> {
		if (auth) {
			headers = {
				...headers,
				Authorization: `Bearer ${await this._jwt.getAccessToken()}`,
			}
		}
		return headers
	}

	private async doRequest(url: string, options: object): Promise<IApiData> {
		return fetch(url, options)
			.then((res) => {
				if (!res.ok) {
					return Promise.reject(res)
				}
				return res.json()
			})
	}
}
