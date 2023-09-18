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

import { useNavigate } from 'react-router-dom';
import { toast } from 'sonner';

export type Payment = {
	id: string;
	amount: number;
	status: 'pending' | 'processing' | 'success' | 'failed';
	email: string;
};

interface PaymentTableProps {
	data: Payment[];
	pageCount: number;
}

export function PaymentTable({ data, pageCount }: PaymentTableProps) {
	const navigate = useNavigate();
	const [isPending, startTransition] = React.useTransition();
	const [open, setOpen] = React.useState(false);
	const [selectedId, setSelectedId] = React.useState('');
	const [userId, setUserId] = React.useState('');

	function onDelete(id: string, userId: string) {
		startTransition(() => {
			try {
				// await deletePaymentAction({
				// 	id: id,
				// 	userId: userId
				// });

				toast.success('Payment deleted successfully.');
				setOpen(false);
			} catch (err) {
				catchError(err);
			}
		});
	}

	const columns = React.useMemo<ColumnDef<Payment, unknown>[]>(
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
				id: 'email',
				accessorKey: 'email',
				header: ({ column }) => (
					<DataTableColumnHeader column={column} title='Email' />
				),
				cell: ({ row }) => {
					return (
						<div className='flex space-x-2'>
							<span className='max-w-[500px] truncate font-medium'>
								{row.getValue('email')}
							</span>
						</div>
					);
				}
			},
			{
				accessorKey: 'amount',
				header: () => <div className='text-right text-xs'>Amount</div>,
				cell: ({ row }) => {
					const amount = parseFloat(row.getValue('amount'));
					const formatted = new Intl.NumberFormat('en-US', {
						style: 'currency',
						currency: 'USD'
					}).format(amount);

					return <div className='text-right font-medium'>{formatted}</div>;
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
								onClick={() => navigate(`/payment/${row.original.id}`)}
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
				onConfirm={() => onDelete(selectedId, userId)}
				loading={isPending}
			/>
			<DataTable
				columns={columns}
				data={data}
				pageCount={pageCount}
				searchableColumns={[
					{
						id: 'email',
						title: 'By Email'
					}
				]}
			/>
		</>
	);
}
