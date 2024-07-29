import { Datagrid, List, TextField } from "react-admin";

export const StatusList = () => {
	return (
		<List>
			<Datagrid>
				<TextField source="id" />
				<TextField source="name" />
				<TextField source="created_at" />
				<TextField source="updated_at" />
			</Datagrid>
		</List>
	);
};
