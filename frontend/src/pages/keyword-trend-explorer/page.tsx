import { GetKeywordTrendExplorers } from '@/../wailsjs/go/keyword_trend_explorer/Service';
import { pagination } from '@/../wailsjs/go/models';
import { EventsOff, EventsOn } from '@/../wailsjs/runtime/runtime';
import { Breadcrumbs } from '@/components/breadcrumbs';
import { KeywordTrendExplorerForm } from '@/components/forms/KeywordTrendExplorer';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle
} from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { Separator } from '@/components/ui/separator';
import { ReloadIcon } from '@radix-ui/react-icons';
import { useEffect, useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { toast } from 'sonner';
import { KeywordTrendExplorersTable } from './components/data-table';

interface ProgressEventPayload {
	progress: number;
}

interface ProgressEventProps {
	payload: ProgressEventPayload;
}

export function KeywordTrendExplorersPage() {
	const navigate = useNavigate();
	const [rows, setRows] = useState([]);
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
		const fetchKeywordTrendExplorers = async () => {
			try {
				const keywordTrendExplorers = await GetKeywordTrendExplorers({
					page: page,
					limit: limit,
					sort: `${column} ${order}`
				});
				const { rows, ...pageInfo } = keywordTrendExplorers.data;

				setRows(rows);
				setPageInfo(pageInfo);
			} catch (error) {
				toast.error('Error fetching data');
			}
		};
		fetchKeywordTrendExplorers();

		const eventListener = () => fetchKeywordTrendExplorers();
		EventsOn('keyword_trend_explorer:update', eventListener);

		return () => EventsOff('keyword_trend_explorer:update');
	}, []);

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
			<div>
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
				<div className='px-4 pb-4'>
					<h1 className='font-bold text-xl'>Keyword Trend Explorer</h1>
				</div>

				<Separator />
			</div>

			<div className='space-y-6'>
				<div className='flex flex-col items-center justify-center [&>div]:w-full space-y-2'>
					<Card>
						<CardHeader>
							<CardTitle>
								Unlock the power of data-driven content creation.
							</CardTitle>
							<CardDescription>
								Save time and effort by accessing Google Trends and Pinterest
								Trends in one place.
							</CardDescription>
						</CardHeader>
						<CardContent className='space-y-3'>
							<KeywordTrendExplorerForm
								isBusy={busy}
								isSubmitting={(v: boolean) => setBusy(v)}
							/>
						</CardContent>
					</Card>

					{busy && (
						<div className='flex items-center space-x-2'>
							<Progress value={progress} />
							{busy && <ReloadIcon className='mr-2 h-4 w-4 animate-spin' />}
							<span className='text-xs font-bold'>{progress}%</span>
						</div>
					)}
				</div>
			</div>

			<Separator />

			{!busy && (
				<div>
					<Card>
						<CardContent className='py-4'>
							<KeywordTrendExplorersTable pageCount={pageCount} data={rows} />
						</CardContent>
					</Card>
				</div>
			)}
		</div>
	);
}
