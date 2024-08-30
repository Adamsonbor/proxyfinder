import { Config } from "../../config";
import { Proxy } from "../../types";
import { Jwt } from "../../utils/jwt/jwt";
import { CRUDMixin } from "../repo";

export class ProxyRepo extends CRUDMixin<Proxy> {
	constructor(config: Config) {
		super(`${config.server.apiUrl}/proxy`, new Jwt(config));
	}
}
