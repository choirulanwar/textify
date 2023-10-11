import * as React from 'react';

import { DotsHorizontalIcon } from '@radix-ui/react-icons';
import { type ColumnDef } from '@tanstack/react-table';

import AlertModal from '@/components/alert-modal';
import { DataTable } from '@/components/data-table/data-table';
import { DataTableColumnHeader } from '@/components/data-table/data-table-column-header';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
	DropdownMenu,
	DropdownMenuContent,
	DropdownMenuItem,
	DropdownMenuSeparator,
	DropdownMenuTrigger
} from '@/components/ui/dropdown-menu';
import { catchError } from '@/lib/utils';

import { DeleteKeywordTrendExplorer } from '@/../wailsjs/go/keyword_trend_explorer/Service';
import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';

export function KeywordTrendExplorersTable({
	data,
	pageCount
}: {
	data: any;
	pageCount: number;
}) {
	const [exportData, setExportData] = React.useState<any>([]);
	const navigate = useNavigate();
	const [isPending, startTransition] = React.useTransition();
	const [open, setOpen] = React.useState(false);
	const [selectedId, setSelectedId] = React.useState(0);

	function onDelete(id: number) {
		startTransition(() => {
			try {
				DeleteKeywordTrendExplorer(id);

				toast.success('Data deleted successfully.');
				setOpen(false);
			} catch (err) {
				catchError(err);
			}
		});
	}

	const columns = React.useMemo<ColumnDef<any, unknown>[]>(
		() => [
			{
				id: 'select',
				header: ({ table }) => (
					<Checkbox
						checked={table.getIsAllPageRowsSelected()}
						onCheckedChange={value => table.toggleAllPageRowsSelected(!!value)}
						aria-label='Select all'
						className='translate-y-[2px]'
					/>
				),
				cell: ({ row }) => (
					<Checkbox
						checked={row.getIsSelected()}
						onCheckedChange={value => row.toggleSelected(!!value)}
						aria-label='Select row'
						className='translate-y-[2px]'
					/>
				),
				enableSorting: false,
				enableHiding: false
			},
			{
				id: 'keyword',
				accessorKey: 'keyword',
				header: ({ column }) => (
					<DataTableColumnHeader column={column} title='Keyword' />
				),
				cell: ({ row }) => {
					return (
						<div className='flex space-x-2'>
							<span className='max-w-[100px] truncate text-sm'>
								{row.getValue('keyword')}
							</span>
						</div>
					);
				}
			},
			{
				id: 'country',
				accessorKey: 'country',
				header: ({ column }) => (
					<DataTableColumnHeader column={column} title='Country' />
				),
				cell: ({ row }) => {
					return (
						<div className='flex space-x-2'>
							<span className='max-w-[100px] truncate text-sm'>
								{row.getValue('country')}
							</span>
						</div>
					);
				}
			},
			{
				id: 'actions',
				cell: ({ row }) => (
					<DropdownMenu>
						<DropdownMenuTrigger asChild>
							<Button
								aria-label='Open menu'
								variant='ghost'
								className='flex h-8 w-8 p-0 data-[state=open]:bg-muted'
							>
								<DotsHorizontalIcon className='h-4 w-4' aria-hidden='true' />
							</Button>
						</DropdownMenuTrigger>
						<DropdownMenuContent align='end' className='w-[160px]'>
							<DropdownMenuItem
								onClick={() =>
									navigate(`/keyword-trend-explorer/${row.original.id}`)
								}
							>
								View
							</DropdownMenuItem>
							<DropdownMenuSeparator />
							<DropdownMenuItem
								onClick={() => {
									setSelectedId(row.original.id);
									setOpen(true);
								}}
							>
								Delete
							</DropdownMenuItem>
						</DropdownMenuContent>
					</DropdownMenu>
				)
			}
		],
		// eslint-disable-next-line react-hooks/exhaustive-deps
		[]
	);

	return (
		<>
			<AlertModal
				isOpen={open}
				onClose={() => setOpen(false)}
				onConfirm={() => onDelete(selectedId)}
				loading={isPending}
			/>
			<DataTable
				columns={columns}
				data={data}
				pageCount={pageCount}
				searchableColumns={[
					{
						id: 'keyword',
						title: 'By Keyword'
					}
				]}
			/>
		</>
	);
}
