import { Config } from "../../config";
import { User } from "../../types";
import { Jwt } from "../../utils/jwt/jwt";
import { CRUDMixin } from "../repo";

export class UserRepo extends CRUDMixin<User> {
	constructor(config: Config) {
		super(`${config.server.apiUrl}/user`, new Jwt(config));
	}
}
