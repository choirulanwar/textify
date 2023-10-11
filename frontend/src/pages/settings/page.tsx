import { models } from '@/../wailsjs/go/models';
import { GetGeneralInfo } from '@/../wailsjs/go/setting/Service';
import { Breadcrumbs } from '@/components/breadcrumbs';
import { SettingForm } from '@/components/forms/Setting';
import {
	Card,
	CardContent,
	CardDescription,
	CardHeader,
	CardTitle
} from '@/components/ui/card';
import { Separator } from '@/components/ui/separator';
import { useEffect, useState } from 'react';
import { toast } from 'sonner';

export function SettingsPage() {
	const [generalInfo, setGeneralInfo] = useState<models.Setting>();
	const [busy, setBusy] = useState<boolean>(false);

	useEffect(() => {
		const fetchGeneralInfo = async () => {
			try {
				const generalInfo = await GetGeneralInfo();
				setGeneralInfo(generalInfo.data);
			} catch (error) {
				toast.error('error fetching data');
			}
		};
		fetchGeneralInfo();
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
							title: 'Settings',
							href: '/settings'
						}
					]}
				/>
				<div className='px-4 pb-4'>
					<h1 className='font-bold text-xl'>Settings</h1>
				</div>

				<Separator />
			</div>

			<div className='space-y-6'>
				<div className='flex flex-col items-center justify-center [&>div]:w-full space-y-2'>
					<Card>
						<CardHeader>
							<CardTitle>Settings title.</CardTitle>
							<CardDescription>Settings description.</CardDescription>
						</CardHeader>
						<CardContent>
							<SettingForm
								generalInfo={generalInfo}
								isBusy={busy}
								isSubmitting={(v: boolean) => setBusy(v)}
							/>
						</CardContent>
					</Card>
				</div>
			</div>
		</div>
	);
}
