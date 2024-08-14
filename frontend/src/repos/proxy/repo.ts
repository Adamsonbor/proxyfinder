import { Config } from "../../config";
import { ProxyV2 } from "../../types";
import { Repo } from "../repo";

export class ProxyV2Repo {
	_repo: Repo;

	constructor(config: Config) {
		this._repo = new Repo(`${config.server.apiV2Url}/proxy`, config);
	}

	async GetAll(): Promise<ProxyV2[]> {
		return this._repo.Get("?perPage=7000", {}, true).then((res) => res.data as ProxyV2[])
	}
}
