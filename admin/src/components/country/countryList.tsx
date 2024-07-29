import { Datagrid, List, TextField } from "react-admin";

export const CountryList = () => {
	return (
		<List>
			<Datagrid>
				<TextField source="id" />
				<TextField source="name" />
				<TextField source="code" />
				<TextField source="created_at" />
				<TextField source="updated_at" />
			</Datagrid>
		</List>
	);
};
