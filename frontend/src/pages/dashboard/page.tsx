import { Breadcrumbs } from '@/components/breadcrumbs';
import { Separator } from '@/components/ui/separator';

export function DashboardPage() {
	return (
		<div className='h-screen p-4 flex flex-col space-y-4'>
			<div>
				<Breadcrumbs
					segments={[
						{
							title: 'Textify',
							href: '/dashboard'
						}
					]}
				/>
				<div className='px-4 pb-4'>
					<h1 className='font-bold text-xl'>Dashboard</h1>
				</div>

				<Separator />
			</div>
		</div>
	);
}
