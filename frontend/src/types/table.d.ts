export interface DataTableSearchableColumn<TData> {
	id: keyof TData;
	title: string;
}

export interface DataTableFilterableColumn<TData>
	extends DataTableSearchableColumn<TData> {
	options: Option[];
}
