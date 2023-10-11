import type {
	DataTableFilterableColumn,
	DataTableSearchableColumn
} from '@/types/table';
import {
	flexRender,
	getCoreRowModel,
	getFacetedRowModel,
	getFacetedUniqueValues,
	getFilteredRowModel,
	getPaginationRowModel,
	getSortedRowModel,
	useReactTable,
	type ColumnDef,
	type ColumnFiltersState,
	type PaginationState,
	type SortingState,
	type VisibilityState
} from '@tanstack/react-table';
import * as React from 'react';
import { useLocation, useNavigate, useSearchParams } from 'react-router-dom';

import { DataTablePagination } from '@/components/data-table/data-table-pagination';
import { DataTableToolbar } from '@/components/data-table/data-table-toolbar';
import {
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableHeader,
	TableRow
} from '@/components/ui/table';
import { useDebounce } from '@/hooks/use-debounce';

interface DataTableProps<TData, TValue> {
	columns: ColumnDef<TData, TValue>[];
	data: TData[];
	pageCount: number;
	filterableColumns?: DataTableFilterableColumn<TData>[];
	searchableColumns?: DataTableSearchableColumn<TData>[];
	toolbar?: React.ReactNode;
}

export function DataTable<TData, TValue>({
	columns,
	data,
	pageCount,
	filterableColumns = [],
	searchableColumns = [],
	toolbar
}: DataTableProps<TData, TValue>) {
	const navigate = useNavigate();
	const { pathname } = useLocation();
	const [searchParams] = useSearchParams();

	// Search params
	const page = searchParams?.get('page') ?? '1';
	const perPage = searchParams?.get('perPage') ?? '100';
	const sort = searchParams?.get('sort');
	const [column, order] = sort?.split('.') ?? [];

	// Create query string
	const createQueryString = React.useCallback(
		(params: Record<string, string | number | null>) => {
			const newSearchParams = new URLSearchParams(searchParams?.toString());

			for (const [key, value] of Object.entries(params)) {
				if (value === null) {
					newSearchParams.delete(key);
				} else {
					newSearchParams.set(key, String(value));
				}
			}

			return newSearchParams.toString();
		},
		[searchParams]
	);

	// Table states
	const [rowSelection, setRowSelection] = React.useState({});
	const [columnVisibility, setColumnVisibility] =
		React.useState<VisibilityState>({});
	const [columnFilters, setColumnFilters] = React.useState<ColumnFiltersState>(
		[]
	);

	// Handle server-side pagination
	const [{ pageIndex, pageSize }, setPagination] =
		React.useState<PaginationState>({
			pageIndex: Number(page) - 1,
			pageSize: Number(perPage)
		});

	const pagination = React.useMemo(
		() => ({
			pageIndex,
			pageSize
		}),
		[pageIndex, pageSize]
	);

	React.useEffect(() => {
		setPagination({
			pageIndex: Number(page) - 1,
			pageSize: Number(perPage)
		});
	}, [page, perPage]);

	React.useEffect(() => {
		navigate(
			`${pathname}?${createQueryString({
				page: pageIndex + 1,
				perPage: pageSize
			})}`
		);

		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [pageIndex, pageSize]);

	// Handle server-side sorting
	const [sorting, setSorting] = React.useState<SortingState>([
		{
			id: column ?? '',
			desc: order === 'desc'
		}
	]);

	React.useEffect(() => {
		navigate(
			`${pathname}?${createQueryString({
				page,
				sort: sorting[0]?.id
					? `${sorting[0]?.id}.${sorting[0]?.desc ? 'desc' : 'asc'}`
					: null
			})}`
		);

		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [sorting]);

	// Handle server-side filtering
	const debouncedSearchableColumnFilters = JSON.parse(
		useDebounce(
			JSON.stringify(
				columnFilters.filter(filter => {
					return searchableColumns.find(column => column.id === filter.id);
				})
			),
			500
		)
	) as ColumnFiltersState;

	const filterableColumnFilters = columnFilters.filter(filter => {
		return filterableColumns.find(column => column.id === filter.id);
	});

	React.useEffect(() => {
		for (const column of debouncedSearchableColumnFilters) {
			if (typeof column.value === 'string') {
				navigate(
					`${pathname}?${createQueryString({
						page: 1,
						[column.id]: typeof column.value === 'string' ? column.value : null
					})}`
				);
			}
		}

		if (searchParams) {
			for (const key of searchParams.keys()) {
				if (
					searchableColumns.find(column => column.id === key) &&
					!debouncedSearchableColumnFilters.find(column => column.id === key)
				) {
					navigate(
						`${pathname}?${createQueryString({
							page: 1,
							[key]: null
						})}`
					);
				}
			}
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [JSON.stringify(debouncedSearchableColumnFilters)]);

	React.useEffect(() => {
		for (const column of filterableColumnFilters) {
			if (typeof column.value === 'object' && Array.isArray(column.value)) {
				navigate(
					`${pathname}?${createQueryString({
						page: 1,
						[column.id]: column.value.join('.')
					})}`
				);
			}
		}

		if (searchParams) {
			for (const key of searchParams.keys()) {
				if (
					filterableColumns.find(column => column.id === key) &&
					!filterableColumnFilters.find(column => column.id === key)
				) {
					navigate(
						`${pathname}?${createQueryString({
							page: 1,
							[key]: null
						})}`
					);
				}
			}
		}
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [JSON.stringify(filterableColumnFilters)]);

	const table = useReactTable({
		data,
		columns,
		pageCount: pageCount ?? -1,
		state: {
			pagination,
			sorting,
			columnVisibility,
			rowSelection,
			columnFilters
		},
		enableRowSelection: true,
		onRowSelectionChange: setRowSelection,
		onPaginationChange: setPagination,
		onSortingChange: setSorting,
		onColumnFiltersChange: setColumnFilters,
		onColumnVisibilityChange: setColumnVisibility,
		getCoreRowModel: getCoreRowModel(),
		getFilteredRowModel: getFilteredRowModel(),
		getPaginationRowModel: getPaginationRowModel(),
		getSortedRowModel: getSortedRowModel(),
		getFacetedRowModel: getFacetedRowModel(),
		getFacetedUniqueValues: getFacetedUniqueValues(),
		manualPagination: true,
		manualSorting: true,
		manualFiltering: true
	});

	return (
		<div className='w-full space-y-4 overflow-auto'>
			<DataTableToolbar
				table={table}
				customToolbar={toolbar}
				filterableColumns={filterableColumns}
				searchableColumns={searchableColumns}
			/>
			<div className='rounded-md'>
				<Table>
					<TableHeader className='bg-secondary'>
						{table.getHeaderGroups().map(headerGroup => (
							<TableRow key={headerGroup.id}>
								{headerGroup.headers.map(header => {
									return (
										<TableHead key={header.id}>
											{header.isPlaceholder
												? null
												: flexRender(
														header.column.columnDef.header,
														header.getContext()
												  )}
										</TableHead>
									);
								})}
							</TableRow>
						))}
					</TableHeader>
					<TableBody>
						{table.getRowModel().rows?.length ? (
							table.getRowModel().rows.map(row => (
								<TableRow
									key={row.id}
									data-state={row.getIsSelected() && 'selected'}
								>
									{row.getVisibleCells().map(cell => (
										<TableCell key={cell.id}>
											{flexRender(
												cell.column.columnDef.cell,
												cell.getContext()
											)}
										</TableCell>
									))}
								</TableRow>
							))
						) : (
							<TableRow>
								<TableCell
									colSpan={columns.length}
									className='h-24 text-center'
								>
									No results.
								</TableCell>
							</TableRow>
						)}
					</TableBody>
				</Table>
			</div>
			<DataTablePagination table={table} />
		</div>
	);
}
