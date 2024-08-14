import { Config } from "../../config";
import { Favorits } from "../../types";
import { IApiData } from "../../utils/api/api";
import { Repo } from "../repo";

export class FavoritsRepo {
	_repo: Repo;

	constructor(config: Config) {
		this._repo = new Repo(`${config.server.apiV2Url}/favorits`, config);
	}

	async GetAll(): Promise<Favorits[]> {
		return await this._repo.GetAll({}, true).then((res) => res.data)
	}

	async Create(userId: number, proxyId: number): Promise<IApiData> {
		const headers = { "Content-Type": "application/json" }
		const data = { "user_id": userId, "proxy_id": proxyId }
		return this._repo.Create(data, headers, true)
	}

	async Delete(proxyId: number): Promise<IApiData> {
		const headers = { "Content-Type": "application/json", }
		const params = `/${proxyId}`
		return await this._repo.Delete(params, headers, true)
	}
}
