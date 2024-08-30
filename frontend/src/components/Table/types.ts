export interface CellParams<T> {
	row: T;
}

export interface Column<T> {
	field: string;
	headerName: string;
	minWidth: number;
	width?: number;
	maxWidth?: number;
	flex: number;
	headerStyle?: object;
	order?: string;
	renderCell: (params: CellParams<T>) => JSX.Element;
	renderHeaderCell?: (params: CellParams<T>) => JSX.Element;
	onClick?: (col: Column<T>) => void;
	sorts?: object;
}

