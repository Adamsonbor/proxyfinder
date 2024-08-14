import { Config } from "../../config";
import { User } from "../../types";
import { Repo } from "../repo";

export class UserRepo {
	_repo: Repo;

	constructor(config: Config) {
		this._repo = new Repo(`${config.server.apiV2Url}/user`, config);
	}

	async Get(): Promise<User> {
		return this._repo.Get("", {}, true)
		.then((res) => res.data as User)
		.catch((err) => Promise.reject(err))
	}
}
