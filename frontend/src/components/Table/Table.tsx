import { Box } from "@mui/material";
import TableHeader from "./TableHeader";
import { TableRow } from "./TableRow";
import { Column } from "./types";
import { useEffect, useRef } from "react";

interface Props<Value> {
	className?: string;
	values?: Value[];
	sorts?: object;
	setValues?: (values: Value[]) => void;
	columns?: Column<Value>[];
	headerStyle?: object;
	bodyStyle?: object;
	onScroll?: () => void;
	onHeaderClick?: (col: Column<Value>) => void;
	sx?: object;
}

export default function Table<Value>({
	className = '',
	values = [],
	columns = [],
	sorts = {},
	onScroll = () => { },
	onHeaderClick = () => { },
	headerStyle = {},
	bodyStyle = {},
	sx = {},
}: Props<Value>) {
	const tableBody = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const handleScroll = (e: Event) => {
			const { scrollTop, scrollHeight, clientHeight } = e.target as HTMLDivElement;
			if (scrollTop + clientHeight == scrollHeight) {
				// console.log(values.length);
				onScroll();
				// console.log(values.length);
			}
		};

		tableBody.current?.addEventListener("scroll", handleScroll);

		return () => {
			tableBody.current?.removeEventListener("scroll", handleScroll);
		}
	}, [onScroll]);

	return (
		<Box
			className={className}
			sx={{
				overflowX: "auto",
				...sx
			}}>
			<Box
				sx={{
					minWidth: "100%",
					width: "fit-content",
					...bodyStyle
				}}>
				<TableHeader
					sx={{
						height: "52px",
						...headerStyle,
					}}
					sorts={sorts}
					onHeaderClick={onHeaderClick}
					columns={columns} />
				<Box
					ref={tableBody}
					sx={{
						height: "100%",
						overflowY: "auto",
					}}>
					{
						values.map((value: Value, index: number) => (
							<TableRow<Value>
								sx={{
									borderTop: "1px solid rgba(0, 0, 0, 0.1)",
									height: "52px",
								}}
								key={index}
								columns={columns}
								value={value} />
						))

					}
				</Box>
			</Box>
		</Box>
	)

}
