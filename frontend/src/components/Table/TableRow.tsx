import { CellParams, Column } from "./types"

export interface TableRowProps<T> {
	columns: Column<T>[]
	value: T
	sx?: object
}

export function TableRow<T>({
	columns,
	value,
	sx = {},
}: TableRowProps<T>) {
	return (
		<span
			style={{
				display: "flex",
				alignItems: "center",
				...sx,
			}}>
			{
				columns.map((column: Column<T>, index: number) => (
					<div
						style={{
							flex: column.flex,
							display: "flex",
							alignItems: "center",
							minWidth: column.minWidth,
							maxWidth: column.maxWidth ?? "200px",
							whiteSpace: "nowrap",
							overflow: "hidden",
							padding: "0px 0px",
						}}
						key={index}>
						{column.renderCell({ row: value } as CellParams<T>)}
					</div>
				))
			}
		</span>
	)
}

