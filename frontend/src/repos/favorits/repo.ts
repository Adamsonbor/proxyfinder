import { Config } from "../../config";
import { Favorits } from "../../types";
import { Jwt } from "../../utils/jwt/jwt";
import { CRUDMixin } from "../repo";

export class FavoritsRepo extends CRUDMixin<Favorits> {
	constructor(config: Config) {
		super(`${config.server.apiUrl}/favorits`, new Jwt(config));
	}
}
