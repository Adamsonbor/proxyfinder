import { Datagrid, List, TextField } from "react-admin";

export const ProxyList = () => {
	return (
		<List>
			<Datagrid>
				<TextField source="id" />
				<TextField source="ip" />
				<TextField source="port" />
				<TextField source="protocol" />
				<TextField source="response_time" />
				<TextField source="status.name" />
				<TextField source="country.name" />
			</Datagrid>
		</List>
	);
};
