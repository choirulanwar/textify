import { KeywordTrendExplorerRequest } from '@/../wailsjs/go/keyword_trend_explorer/Service';
import { EventsEmit } from '@/../wailsjs/runtime/runtime';
import { Button } from '@/components/ui/button';
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormMessage
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { ScrollArea } from '@/components/ui/scroll-area';
import {
	Select,
	SelectContent,
	SelectGroup,
	SelectItem,
	SelectTrigger,
	SelectValue
} from '@/components/ui/select';
import { Separator } from '@/components/ui/separator';
import { catchError } from '@/lib/utils';
import {
	KeywordTrendExplorerInputs,
	KeywordTrendExplorerSchema
} from '@/lib/validations/keyword-trend-explorer';
import { zodResolver } from '@hookform/resolvers/zod';
import {
	CrossCircledIcon,
	MagnifyingGlassIcon,
	ReloadIcon
} from '@radix-ui/react-icons';
import {
	AR,
	AU,
	BR,
	CA,
	CL,
	CO,
	DE,
	ES,
	FR,
	GB,
	ID,
	IT,
	MX,
	NL,
	PL,
	PT,
	SE,
	US
} from 'country-flag-icons/react/3x2';
import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { Toaster } from 'sonner';

type Props = {
	isBusy?: boolean;
	isSubmitting: (v: boolean) => void;
};

export function KeywordTrendExplorerForm({
	isBusy = false,
	isSubmitting
}: Props) {
	const navigate = useNavigate();
	const [loading, setLoading] = useState<boolean>(false);

	const form = useForm<KeywordTrendExplorerInputs>({
		resolver: zodResolver(KeywordTrendExplorerSchema),
		defaultValues: {
			country: 'US',
			period: '3'
		}
	});

	function onSubmit(data: KeywordTrendExplorerInputs) {
		try {
			setLoading(true);
			isSubmitting(true);
			setTimeout(async () => {
				const newExplorers = await KeywordTrendExplorerRequest({
					query: data.query,
					country: data.country,
					language: data.language,
					period: data.period,
					include_serp: false
				});

				// await ProgressTracker();
				setLoading(false);

				navigate(`/keyword-trend-explorer/${newExplorers.data.id}`);
				isSubmitting(false);
			}, 1000);
		} catch (error) {
			catchError(error);
			isSubmitting(false);
		}
	}

	async function abortProcess() {
		try {
			await EventsEmit('progress:stop');
		} catch (error) {
			console.error('[+] Error when aborting process.');
		} finally {
			setLoading(false);
		}
	}

	return (
		<>
			<Toaster />
			<Form {...form}>
				<form onSubmit={form.handleSubmit(onSubmit)}>
					<div className='space-y-3'>
						<FormField
							control={form.control}
							name='query'
							render={({ field }) => (
								<FormItem>
									<FormControl>
										<Input placeholder='e.g. best SEO tool' {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<Separator />

						<div className='flex flex-row justify-between items-center'>
							<div className='flex items-center space-x-2 text-sm'>
								<span className='text-xs'>Results for</span>

								<FormField
									control={form.control}
									name='country'
									render={({ field }) => (
										<FormItem>
											<Select
												{...field}
												onValueChange={field.onChange}
												defaultValue={field.value}
											>
												<FormControl>
													<SelectTrigger className='w-[250px]'>
														<SelectValue placeholder='Select a country' />
													</SelectTrigger>
												</FormControl>
												<SelectContent>
													<SelectGroup>
														<ScrollArea className='h-72'>
															<SelectItem value='US'>
																<div className='flex flex-row space-x-2'>
																	<US width={17} />
																	<span className='text-xs'>United States</span>
																</div>
															</SelectItem>
															<SelectItem value='ID'>
																<div className='flex flex-row space-x-2'>
																	<ID width={17} />
																	<span className='text-xs'>Indonesia</span>
																</div>
															</SelectItem>
															<SelectItem value='CA'>
																<div className='flex flex-row space-x-2'>
																	<CA width={17} />
																	<span className='text-xs'>Canada</span>
																</div>
															</SelectItem>
															<SelectItem value='IT+ES+PT+GR+MT'>
																<div className='flex flex-row space-x-2'>
																	<PT width={17} />
																	<span className='text-xs'>South Europe</span>
																</div>
															</SelectItem>
															<SelectItem value='IT'>
																<div className='flex flex-row space-x-2'>
																	<IT width={17} />
																	<span className='text-xs'>Italia</span>
																</div>
															</SelectItem>
															<SelectItem value='ES'>
																<div className='flex flex-row space-x-2'>
																	<ES width={17} />
																	<span className='text-xs'>Spain</span>
																</div>
															</SelectItem>
															<SelectItem value='DE+AT+CH'>
																<div className='flex flex-row space-x-2'>
																	<DE width={17} />
																	<span className='text-xs'>Germany</span>
																</div>
															</SelectItem>
															<SelectItem value='GB+IE'>
																<div className='flex flex-row space-x-2'>
																	<GB width={17} />
																	<span className='text-xs'>
																		Great Britain and Ireland
																	</span>
																</div>
															</SelectItem>
															<SelectItem value='FR'>
																<div className='flex flex-row space-x-2'>
																	<FR width={17} />
																	<span className='text-xs'>France</span>
																</div>
															</SelectItem>
															<SelectItem value='SE+DK+FI+NO'>
																<div className='flex flex-row space-x-2'>
																	<SE width={17} />
																	<span className='text-xs'>
																		Nordic Country
																	</span>
																</div>
															</SelectItem>
															<SelectItem value='NL+BE+LU'>
																<div className='flex flex-row space-x-2'>
																	<NL width={17} />
																	<span className='text-xs'>Benelux</span>
																</div>
															</SelectItem>
															<SelectItem value='PL+RO+HU+SK+CZ'>
																<div className='flex flex-row space-x-2'>
																	<PL width={17} />
																	<span className='text-xs'>East Europe</span>
																</div>
															</SelectItem>
															<SelectItem value='MX+AR+CO+CL'>
																<div className='flex flex-row space-x-2'>
																	<CL width={17} />
																	<span className='text-xs'>
																		Latin America Hispanic
																	</span>
																</div>
															</SelectItem>
															<SelectItem value='CO'>
																<div className='flex flex-row space-x-2'>
																	<CO width={17} />
																	<span className='text-xs'>Colombia</span>
																</div>
															</SelectItem>
															<SelectItem value='AR'>
																<div className='flex flex-row space-x-2'>
																	<AR width={17} />
																	<span className='text-xs'>Argentina</span>
																</div>
															</SelectItem>
															<SelectItem value='MX'>
																<div className='flex flex-row space-x-2'>
																	<MX width={17} />
																	<span className='text-xs'>Mexico</span>
																</div>
															</SelectItem>
															<SelectItem value='BR'>
																<div className='flex flex-row space-x-2'>
																	<BR width={17} />
																	<span className='text-xs'>Brazil</span>
																</div>
															</SelectItem>
															<SelectItem value='AU+NZ'>
																<div className='flex flex-row space-x-2'>
																	<AU width={17} />
																	<span className='text-xs'>Australasia</span>
																</div>
															</SelectItem>
														</ScrollArea>
													</SelectGroup>
												</SelectContent>
											</Select>
											<FormMessage />
										</FormItem>
									)}
								/>

								<FormField
									control={form.control}
									name='period'
									render={({ field }) => (
										<FormItem>
											<Select
												{...field}
												onValueChange={field.onChange}
												defaultValue={field.value}
											>
												<FormControl>
													<SelectTrigger className='w-[190px]'>
														<SelectValue placeholder='Select a trends period' />
													</SelectTrigger>
												</FormControl>
												<SelectContent>
													<SelectGroup>
														<SelectItem value='0'>
															<span className='text-xs'>
																Top daily trends (Google)
															</span>
														</SelectItem>
														<SelectItem value='1'>
															<span className='text-xs'>
																Top monthly trends
															</span>
														</SelectItem>
														<SelectItem value='2'>
															<span className='text-xs'>Top yearly trends</span>
														</SelectItem>
														<SelectItem value='3'>
															<span className='text-xs'>Growing trends</span>
														</SelectItem>
														<SelectItem value='4'>
															<span className='text-xs'>Seasonal trends</span>
														</SelectItem>
													</SelectGroup>
												</SelectContent>
											</Select>
											<FormMessage />
										</FormItem>
									)}
								/>
							</div>

							{!isBusy && !loading ? (
								<Button
									type='submit'
									className='w-[100px]'
									size={'sm'}
									variant={'default'}
									disabled={isBusy && loading}
								>
									{isBusy && loading ? (
										<ReloadIcon className='mr-2 h-4 w-4 animate-spin text-primary-foreground' />
									) : (
										<MagnifyingGlassIcon className='mr-2 h-4 w-4 text-primary-foreground' />
									)}
									Explorer
								</Button>
							) : (
								<Button
									className='w-[100px]'
									size={'sm'}
									variant={'destructive'}
									disabled={!isBusy && !loading}
									onClick={abortProcess}
								>
									<CrossCircledIcon className='mr-2 h-4 w-4 text-white' />
									Cancel
								</Button>
							)}
						</div>
					</div>
				</form>
			</Form>
		</>
	);
}
