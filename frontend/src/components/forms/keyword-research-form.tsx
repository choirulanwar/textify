import { Button } from '@/components/ui/button';
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
	KeywordResearchInputs,
	KeywordResearchSchema
} from '@/lib/validations/keyword-research';
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
import { v4 as uuidv4 } from 'uuid';
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormMessage
} from '../ui/form';

type Props = {
	isSubmitting: (v: boolean) => void;
};

export function KeywordResearchForm({ isSubmitting }: Props) {
	const navigate = useNavigate();
	const [loading, setLoading] = useState<boolean>(false);

	const form = useForm<KeywordResearchInputs>({
		resolver: zodResolver(KeywordResearchSchema),
		defaultValues: {
			country: 'US',
			period: '3'
		}
	});

	function onSubmit(data: KeywordResearchInputs) {
		try {
			const id = uuidv4();

			setLoading(true);
			isSubmitting(true);
			setTimeout(async () => {
				setLoading(false);
				navigate('/keyword-research', {
					state: {
						researchId: id
					}
				});
				isSubmitting(false);
			}, 1000);
			console.log('[+] Data', data);
		} catch (error) {
			catchError(error);
			isSubmitting(false);
		}
	}

	async function abortProcess() {
		try {
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
													<SelectTrigger className='w-[170px]'>
														<SelectValue placeholder='Select a trends period' />
													</SelectTrigger>
												</FormControl>
												<SelectContent>
													<SelectGroup>
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

							{!loading ? (
								<Button
									type='submit'
									className='w-[100px]'
									size={'sm'}
									variant={'default'}
									disabled={loading}
								>
									{loading ? (
										<ReloadIcon className='mr-2 h-4 w-4 animate-spin text-primary-foreground' />
									) : (
										<MagnifyingGlassIcon className='mr-2 h-4 w-4 text-primary-foreground' />
									)}
									Explore
								</Button>
							) : (
								<Button
									className='w-[100px]'
									size={'sm'}
									variant={'destructive'}
									disabled={!loading}
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
