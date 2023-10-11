import { models } from '@/../wailsjs/go/models';
import { SetGeneralData } from '@/../wailsjs/go/setting/Service';
import { Button } from '@/components/ui/button';
import { Checkbox } from '@/components/ui/checkbox';
import {
	Form,
	FormControl,
	FormField,
	FormItem,
	FormLabel,
	FormMessage
} from '@/components/ui/form';
import { Input } from '@/components/ui/input';
import { Separator } from '@/components/ui/separator';
import { Textarea } from '@/components/ui/textarea';
import { catchError } from '@/lib/utils';
import { SettingInputs, SettingSchema } from '@/lib/validations/setting';
import { zodResolver } from '@hookform/resolvers/zod';
import { UpdateIcon } from '@radix-ui/react-icons';
import { useEffect, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { Toaster, toast } from 'sonner';

type Props = {
	generalInfo?: models.Setting;
	isBusy?: boolean;
	isSubmitting: (v: boolean) => void;
};

export function SettingForm({
	generalInfo,
	isBusy = false,
	isSubmitting
}: Props) {
	const navigate = useNavigate();
	const [loading, setLoading] = useState<boolean>(false);

	const form = useForm<SettingInputs>({
		resolver: zodResolver(SettingSchema),
		defaultValues: generalInfo && {
			browserPath: generalInfo.browser_path,
			browserVisible: generalInfo.browser_visible,
			sessionGoogle: generalInfo.session_google,
			sessionPinterest: generalInfo.session_pinterest,
			proxyUrl: generalInfo.proxy_url
		}
	});

	useEffect(() => {
		if (generalInfo) {
			form.reset({
				browserPath: generalInfo.browser_path,
				browserVisible: generalInfo.browser_visible,
				sessionGoogle: generalInfo.session_google,
				sessionPinterest: generalInfo.session_pinterest,
				proxyUrl: generalInfo.proxy_url
			});
		}
	}, [generalInfo]);

	function onSubmit(data: SettingInputs) {
		try {
			setLoading(true);
			isSubmitting(true);
			setTimeout(async () => {
				await SetGeneralData({
					id: 1,
					browser_path: data?.browserPath || '',
					browser_visible: data.browserVisible,
					session_google: data?.sessionGoogle || '',
					session_pinterest: data.sessionPinterest,
					proxy_url: data?.proxyUrl || ''
				});

				setLoading(false);

				navigate(`/settings`);
				isSubmitting(false);
				toast.success('Settings updated');
			}, 1000);
		} catch (error) {
			catchError(error);
			isSubmitting(false);
		}
	}

	return (
		<>
			<Toaster />
			<Form {...form}>
				<form onSubmit={form.handleSubmit(onSubmit)}>
					<div className='space-y-5'>
						<FormField
							control={form.control}
							name='sessionGoogle'
							render={({ field }) => (
								<FormItem>
									<div className=''>
										<FormLabel>Google Session</FormLabel>
									</div>
									<FormControl>
										<Textarea placeholder='google session/cookies' {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name='sessionPinterest'
							render={({ field }) => (
								<FormItem>
									<div className=''>
										<FormLabel>Pinterest Session</FormLabel>
									</div>
									<FormControl>
										<Textarea
											placeholder='pinterest session/cookies'
											{...field}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name='proxyUrl'
							render={({ field }) => (
								<FormItem>
									<div className=''>
										<FormLabel>Proxy URL</FormLabel>
									</div>
									<FormControl>
										<Input placeholder='https://johndoe:8080' {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name='browserPath'
							render={({ field }) => (
								<FormItem>
									<div className=''>
										<FormLabel>Browser path location</FormLabel>
									</div>
									<FormControl>
										<Input placeholder='chrome.exe' {...field} />
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<FormField
							control={form.control}
							name='browserVisible'
							render={({ field }) => (
								<FormItem className='flex flex-row space-x-2 items-center space-y-0'>
									<div className=''>
										<FormLabel>Show browser</FormLabel>
									</div>
									<FormControl>
										<Checkbox
											checked={field.value}
											onCheckedChange={field.onChange}
										/>
									</FormControl>
									<FormMessage />
								</FormItem>
							)}
						/>

						<Separator />

						<div className='flex flex-row justify-end items-center'>
							<Button
								type='submit'
								className='w-[100px]'
								size={'sm'}
								variant={'default'}
								disabled={isBusy && loading}
							>
								{isBusy && loading ? (
									<UpdateIcon className='mr-2 h-4 w-4 animate-spin text-primary-foreground' />
								) : (
									<UpdateIcon className='mr-2 h-4 w-4 text-primary-foreground' />
								)}
								Save
							</Button>
						</div>
					</div>
				</form>
			</Form>
		</>
	);
}
