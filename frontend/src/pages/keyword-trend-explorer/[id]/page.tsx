import {
	GetKeywordTrendExplorer,
	GetKeywordsByExplorers
} from '@/../wailsjs/go/keyword_trend_explorer/Service';
import { pagination } from '@/../wailsjs/go/models';
import { EventsOff, EventsOn } from '@/../wailsjs/runtime/runtime';
import { Breadcrumbs } from '@/components/breadcrumbs';
import { Card, CardContent } from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { useEffect, useState } from 'react';
import { useParams, useSearchParams } from 'react-router-dom';
import { toast } from 'sonner';
import { KeywordTrendExplorerTable } from './components/data-table';

interface ProgressEventPayload {
	progress: number;
}

interface ProgressEventProps {
	payload: ProgressEventPayload;
}

export function KeywordTrendExplorerPage() {
	const { id } = useParams();
	const [data, setData] = useState<any>();
	const [rows, setRows] = useState<any[]>([]);
	const [pageInfo, setPageInfo] = useState<pagination.Pagination>({});

	const [searchParams] = useSearchParams();
	const paramsPage = searchParams?.get('page');
	const perPage = searchParams?.get('perPage');
	const sort = searchParams?.get('sort');
	const title = searchParams?.get('title');
	const [column, order] = sort?.split('.') ?? ['created_at', 'desc'];

	const limit = typeof perPage === 'string' ? parseInt(perPage) : 100;
	const page = typeof paramsPage === 'string' ? parseInt(paramsPage) : 1;

	useEffect(() => {
		const fetchKeywordByExplorers = async () => {
			try {
				if (typeof id === 'string') {
					const keywordsByExplorers = await GetKeywordsByExplorers(
						parseInt(id),
						{
							page: page,
							limit: limit,
							sort: `${column} ${order}`
						}
					);
					const { rows, ...pageInfo } = keywordsByExplorers.data;
					setRows(rows);
					setPageInfo(pageInfo);
				}
			} catch (error) {
				toast.error('Error fetching data');
			}
		};

		const fetchKeywordTrendExplorer = async () => {
			try {
				if (typeof id === 'string') {
					const keywordTrendExplorer = await GetKeywordTrendExplorer(
						parseInt(id)
					);
					const { data } = keywordTrendExplorer;
					setData(data);
				}
			} catch (error) {
				toast.error('Error fetching data');
			}
		};

		fetchKeywordTrendExplorer();
		fetchKeywordByExplorers();

		const eventListener = () => fetchKeywordByExplorers();
		EventsOn(`keyword_trend_explorer:update:${id}`, eventListener);
		EventsOn('keyword:update', eventListener);
		return () =>
			EventsOff(`keyword_trend_explorer:update:${id}`, 'keyword:update');
	}, [id, page, limit, sort]);

	const pageCount = pageInfo && pageInfo.total_pages ? pageInfo.total_pages : 0;

	const [busy, setBusy] = useState<boolean>(false);
	const [progress, setProgress] = useState<number>(0);

	useEffect(() => {
		EventsOn('progress:update', function (e: ProgressEventPayload) {
			setProgress(e.progress);
			if (e.progress > 0) setBusy(true);
			if (e.progress === 100) {
				toast.success('Success');
				setBusy(false);
			}
		});
		return () => EventsOff('progress:update');
	}, []);

	return (
		<div className='h-screen p-4 flex flex-col space-y-4'>
			<div className=''>
				<Breadcrumbs
					segments={[
						{
							title: 'Textify',
							href: '/dashboard'
						},
						{
							title: 'Keyword Trend Explorer',
							href: '/keyword-trend-explorer'
						}
					]}
				/>
				{data && (
					<div className='px-4 pb-4 flex justify-between'>
						<span className='text-sm font-bold'>Keyword: {data.keyword}</span>
						<span className='text-xs font-medium'>
							Results: {pageInfo.total_rows || 0}
						</span>
					</div>
				)}

				<Separator />
			</div>

			{!busy && (
				<div>
					<Card>
						<CardContent className='py-4'>
							<KeywordTrendExplorerTable pageCount={pageCount} data={rows} />
						</CardContent>
					</Card>
				</div>
			)}
		</div>
	);
}
