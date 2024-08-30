import { IApiData } from "../utils/api/api"
import { Jwt } from "../utils/jwt/jwt";
import { QueryBuilder } from "./utils";

export interface QBuilder {
	[key: string]: any

	build(): string
}

export interface Options {
	[key: string]: [string, any]
}

export class UrlMixin {
	_url: string;

	constructor(url: string) {
		this._url = url
	}

	public get url(): string {
		return this._url
	}

	public set url(url: string) {
		this._url = url
	}
}

export class CRUDMixin<T> extends UrlMixin {
	private _jwt: Jwt
	private _qb: QueryBuilder

	constructor(url: string, jwt: Jwt, filters: object = {}, sorts: object = {}) {
		super(url)
		this._jwt = jwt
		this._qb = new QueryBuilder()
		this._qb.setFilters(filters)
		this._qb.setSorts(sorts)
	}

	async GetAll(qb?: QueryBuilder | object, auth: boolean = false): Promise<T[]> {
		if (qb instanceof QueryBuilder) {
			return this._getData(qb, {}, auth)
		} else if (qb) {
			this._qb.setFilters(qb)
			return this._getData(this._qb, {}, auth)
		}

		return this._getData(this._qb, {}, auth)
	}

	async Get(id: number, auth: boolean = false): Promise<T> {
		const headers = await authHeadersWrapper(this._jwt, {}, auth)
		return this.doRequest(`${this._url}/${id}`, {
			method: "GET",
			headers: headers,
		}).then((res) => res.data as T)
	}

	async GetBy(options: Options = {}, headers: object = {}, auth: boolean = false): Promise<T> {
		const db = new QueryBuilder()
		db.setFilters(options)

		headers = await authHeadersWrapper(this._jwt, headers, auth)
		const opts = {
			method: "GET",
			headers: headers,
		}

		return this.doRequest(`${this._url}${db.build()}`, opts).then((res) => res.data as T)
	}

	async Create(data: T, auth: boolean = true): Promise<IApiData> {
		const headers = await authHeadersWrapper(this._jwt, {}, auth)
		headers["Content-Type"] = "application/json"
		return this.doRequest(this._url, {
			method: "POST",
			headers: headers,
			body: JSON.stringify(data),
		})
	}

	async Update(id: number, data: T, auth: boolean = true): Promise<IApiData> {
		const headers = await authHeadersWrapper(this._jwt, {}, auth)
		headers["Content-Type"] = "application/json"
		return this.doRequest(`${this._url}/${id}`, {
			method: "PUT",
			headers: headers,
			body: JSON.stringify(data),
		})
	}

	async Delete(id: number, options: Options = {}, auth: boolean = true): Promise<IApiData> {
		const qb = new QueryBuilder()
		qb.setFilters(options)

		const headers = await authHeadersWrapper(this._jwt, {}, auth)
		return this.doRequest(`${this._url}/${id}${qb.build()}`, {
			method: "DELETE",
			headers: headers,
		})
	}

	private async _getData(builder: QueryBuilder, headers: object = {}, auth: boolean = false): Promise<T[]> {
		headers = await authHeadersWrapper(this._jwt, headers, auth)
		const opts = {
			method: "GET",
			headers: headers,
		}
		return this.doRequest(`${this._url}${builder.build()}`, opts)
			.then((res) => {
				if (!res.data) {
					return []
				}

				return res.data as T[]
			})
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

async function authHeadersWrapper(jwt: Jwt, headers: object, auth: boolean): Promise<object> {
	if (auth) {
		headers = {
			...headers,
			Authorization: `Bearer ${await jwt.getAccessToken()}`,
		}
	}
	return headers
}
