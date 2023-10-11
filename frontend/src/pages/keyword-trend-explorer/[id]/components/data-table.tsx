import { IconBrandGoogle, IconBrandPinterest } from '@tabler/icons-react';
import { type ColumnDef } from '@tanstack/react-table';
import * as React from 'react';
import { Area, AreaChart, ResponsiveContainer } from 'recharts';

import AlertModal from '@/components/alert-modal';
import { DataTable } from '@/components/data-table/data-table';
import { DataTableColumnHeader } from '@/components/data-table/data-table-column-header';
import { Checkbox } from '@/components/ui/checkbox';
import { catchError, convertToMonthlyData } from '@/lib/utils';

import { OpenSaveFileDialog } from '@/../wailsjs/go/file/Service';
import {
	DeleteKeyword,
	GetKeywordsByExplorers
} from '@/../wailsjs/go/keyword/Service';
import { Button } from '@/components/ui/button';
import type { TrendsData } from '@/lib/utils';
import { Parser } from '@json2csv/plainjs';
import { DownloadIcon } from '@radix-ui/react-icons';
import { toast } from 'sonner';

export function KeywordTrendExplorerTable({
	data,
	pageCount
}: {
	data: any;
	pageCount: number;
}) {
	const [exportData, setExportData] = React.useState<string>('');
	const [isPending, startTransition] = React.useTransition();
	const [open, setOpen] = React.useState(false);
	const [selectedId, setSelectedId] = React.useState(0);

	async function fetchKeywords(keywordTrendExplorerID: number) {
		try {
			const keywords = await GetKeywordsByExplorers(keywordTrendExplorerID);
			const parser = new Parser();
			const csv = parser.parse(keywords.data);

			setExportData(csv);
		} catch (error) {
			toast.error('failed to get data');
		}
	}

	function onDelete(id: number) {
		startTransition(() => {
			try {
				DeleteKeyword(id);

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
						<div className='flex space-x-2 flex-wrap'>
							<span className='max-w-[200px] text-sm'>
								{/* <Badge>{row.getValue('search_intent')}</Badge> */}
								{row.getValue('keyword')}
							</span>
						</div>
					);
				}
			},
			{
				id: 'trends',
				accessorKey: 'trends',
				header: ({ column }) => (
					<DataTableColumnHeader column={column} title='Trends' />
				),
				cell: ({ row }) => {
					const source = row.getValue('source');
					const trendsData = row.getValue<TrendsData[]>('trends');

					return (
						<div className='w-16 h-6'>
							{Array.isArray(trendsData) && (
								<ResponsiveContainer width='100%' height='100%'>
									<AreaChart
										data={
											source === 'pinterest'
												? convertToMonthlyData(trendsData)
												: trendsData
										}
									>
										<Area
											dataKey='value'
											fill='#1784ff'
											stroke='#1271e6'
											type='monotone'
										/>
									</AreaChart>
								</ResponsiveContainer>
							)}
						</div>
					);
				}
			},
			{
				id: 'volume',
				accessorKey: 'volume',
				header: ({ column }) => (
					<DataTableColumnHeader column={column} title='Volume' />
				),
				cell: ({ row }) => {
					return (
						<div className='flex justify-center'>
							<span className='max-w-[100px] truncate text-center text-sm'>
								{row.getValue('volume')}
							</span>
						</div>
					);
				}
			},
			{
				id: 'cpc',
				accessorKey: 'cpc',
				header: ({ column }) => (
					<DataTableColumnHeader column={column} title='CPC' />
				),
				cell: ({ row }) => {
					return (
						<div className='flex justify-center'>
							<span className='max-w-[100px] truncate text-center text-sm'>
								{row.getValue('cpc')}
							</span>
						</div>
					);
				}
			},
			{
				id: 'seo_difficulty',
				accessorKey: 'seo_difficulty',
				header: ({ column }) => (
					<DataTableColumnHeader column={column} title='KD' />
				),
				cell: ({ row }) => {
					return (
						<div className='flex justify-center'>
							<span className='max-w-[100px] truncate text-center text-sm'>
								{row.getValue('seo_difficulty')}
							</span>
						</div>
					);
				}
			},
			{
				id: 'source',
				accessorKey: 'source',
				header: ({ column }) => (
					<DataTableColumnHeader column={column} title='Source' />
				),
				cell: ({ row }) => {
					const source = row.getValue('source');
					return (
						<div className='flex justify-center items-center'>
							{source === 'google' ? (
								<IconBrandGoogle
									size={20}
									className='text-primary text-center'
								/>
							) : (
								<IconBrandPinterest
									size={20}
									className='text-destructive text-center'
								/>
							)}
						</div>
					);
				}
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
				toolbar={
					<>
						<Button
							size='sm'
							className='ml-auto h-8'
							onClick={async () => {
								try {
									await fetchKeywords(data[0].keyword_trend_explorer_id);

									if (exportData) {
										await OpenSaveFileDialog('*.csv', 'output.csv', exportData);
									}
								} catch (error) {}
							}}
						>
							<DownloadIcon className='mr-2 h-4 w-4' />
							Export
						</Button>
					</>
				}
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
