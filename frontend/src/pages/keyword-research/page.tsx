import { Breadcrumbs } from '@/components/breadcrumbs';
import { KeywordResearchForm } from '@/components/forms/keyword-research-form';
import { Button } from '@/components/ui/button';
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
import { useLocation, useSearchParams } from 'react-router-dom';
import { toast } from 'sonner';
import { Payment, PaymentTable } from './components/data-table';

interface ProgressEventPayload {
	progress: number;
}

interface ProgressEventProps {
	payload: ProgressEventPayload;
}

interface Props {
	id?: string;
	searchParams: {
		[key: string]: string | string[] | undefined;
	};
}

type SearchParams = {
	[key: string]: string | string[] | undefined;
};

function getData(): Payment[] {
	// Fetch data from your API here.
	return [
		{
			id: '728ed52f',
			amount: 100,
			status: 'pending',
			email: 'm@example.com'
		},
		{
			id: 'asdsadsadas',
			amount: 100,
			status: 'pending',
			email: 'a@example.com'
		},
		{
			id: 'zxcxzcxzczx',
			amount: 100,
			status: 'pending',
			email: 'c@example.com'
		}
	];
}

export function KeywordResearchPage() {
	const { state } = useLocation();
	const [researchId, setResearchId] = useState<string>('');

	useEffect(() => {
		if (state?.researchId) {
			setResearchId(state.researchId);
		}
	}, [state]);

	const [searchParams] = useSearchParams();

	const page = searchParams?.get('page');
	const perPage = searchParams?.get('perPage');
	const sort = searchParams?.get('sort');
	const title = searchParams?.get('title');
	const [column, order] = sort?.split('.') ?? [];

	const limit = typeof perPage === 'string' ? parseInt(perPage) : 10;
	const offset = typeof page === 'string' ? (parseInt(page) - 1) * limit : 0;

	// const aisTransaction = await getAisAction({
	// 	title,
	// 	limit,
	// 	offset,
	// 	sort:
	// 		typeof sort === 'string'
	// 			? (sort as `${keyof Ai | 'createdAt'}.${'asc' | 'desc'}`)
	// 			: 'createdAt.desc',
	// 	userId: user.id
	// });

	let data: Payment[];
	if (title) {
		data = getData().filter(d => d.email.includes(title.toString()));
	}
	data = getData();

	const pageCount = Math.ceil(data.length / limit);

	const [busy, setBusy] = useState<boolean>(false);
	const [progress, setProgress] = useState<number>(0);
	/*
	useEffect(() => {
		const z = uuidv4();

		const unListen = listen('PROGRESS', (e: ProgressEventProps) => {
			setProgress(e.payload.progress);

			if (e.payload.progress === 100) {
				toast.success('Success');
			}
		});

		return () => {
			unListen.then(f => f());
		};
	}, []);
*/

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
							title: 'Keyword Research',
							href: '/keyword-research'
						}
					]}
				/>
				<div className='px-4 pb-4'>
					<h1 className='font-bold text-xl'>Keyword Research</h1>
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
							<KeywordResearchForm isSubmitting={(v: boolean) => setBusy(v)} />
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

			{researchId && (
				<>
					<Card>
						<CardContent className='py-4'>
							<PaymentTable pageCount={pageCount} data={data} />
						</CardContent>
					</Card>

					<div className='w-full flex justify-center  flex-col space-y-2'>
						<div>
							<Button
								size={'sm'}
								variant={'default'}
								disabled={busy}
								onClick={() => {
									setBusy(true);
									setTimeout(async () => {
										const date = new Date();
										date.setDate(date.getDate() - 7);

										try {
										} catch (error) {
											toast.error(JSON.stringify(error, null, 2));
										}

										setBusy(false);
									}, 1000);
								}}
							>
								Trends
							</Button>
						</div>
						<div>
							<Button
								size={'sm'}
								variant={'default'}
								disabled={busy}
								onClick={() => {
									setBusy(true);
									setTimeout(async () => {
										const date = new Date();
										date.setDate(date.getDate() - 7);

										try {
										} catch (error) {
											toast.error(JSON.stringify(error));
										}

										setBusy(false);
									}, 1000);
								}}
							>
								Related Terms
							</Button>
						</div>
						<div>
							<Button
								size={'sm'}
								variant={'default'}
								disabled={busy}
								onClick={() => {
									setBusy(true);
									setTimeout(async () => {
										const date = new Date();
										date.setDate(date.getDate() - 7);

										try {
										} catch (error) {
											toast.error(JSON.stringify(error));
										}

										setBusy(false);
									}, 1000);
								}}
							>
								Metrics Data
							</Button>
						</div>
					</div>
				</>
			)}
		</div>
	);
}
