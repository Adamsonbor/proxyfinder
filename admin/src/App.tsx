import {
  Admin,
  Resource,
} from "react-admin";
import { Layout } from "./Layout";
import { authProvider } from "./authProvider";
import { proxyfinderDataProvider } from "./myDataProvider";
import { ProxyList } from "./components/proxy/proxylist";
import { StatusList } from "./components/status/statuslist";
import { CountryList } from "./components/country/countryList";

export const App = () => (
  <Admin
    layout={Layout}
    dataProvider={proxyfinderDataProvider}
    authProvider={authProvider}
  >
	<Resource 
		name="proxy"
		list={ProxyList}
	/>
	<Resource 
		name="status"
		list={StatusList}
	/>
	<Resource 
		name="country"
		list={CountryList}
	/>
  </Admin>
);
