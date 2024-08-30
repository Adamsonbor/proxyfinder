import { Config } from "../../config";
import { Country } from "../../types";
import { Jwt } from "../../utils/jwt/jwt";
import { CRUDMixin } from "../repo";

export class CountryRepo extends CRUDMixin<Country> {
	constructor(config: Config) {
		super(`${config.server.apiUrl}/countries`, new Jwt(config));
	}
}
