import { Box, Button } from "@mui/material"
import { Column } from "./types"
// import { useEffect, useState } from "react"
import { FaSortDown, FaSortUp } from "react-icons/fa"

export interface TableHeaderProps<T> {
	columns: Column<T>[]
	sorts?: object
	onHeaderClick?: (col: Column<T>) => void
	sx?: object
}

export default function TableHeader<T>({
	columns,
	sorts = {},
	onHeaderClick = () => { },
	sx = {},
}: TableHeaderProps<T>) {

	return (
		<span
			style={{
				display: "flex",
				...sx
			}} >
			{
				columns.map((column: Column<T>, index: number) => (
					<Button
						variant="text"
						style={{
							flex: column.flex,
							minWidth: column.minWidth,
							maxWidth: column.maxWidth ?? "200px",
							whiteSpace: "nowrap",
							overflow: "hidden",
							justifyContent: "left",
							position: "relative",
							padding: "0px 0px",
						}}
						onClick={() => {
							column.onClick || onHeaderClick(column)
						}}
						key={index}>
						{column.headerName}
						<Box
						sx={{
							display: "flex",
							flexDirection: "column",
						}}> 
							<FaSortUp
								style={{
									position: "absolute",
									opacity: sorts[column.field] == "asc" ? 1 : 0
								}} />
							<FaSortDown
								style={{
									opacity: sorts[column.field] == "desc" ? 1 : 0
								}} />
						</Box>
					</Button>
				))
			}
		</span>
	)
}
