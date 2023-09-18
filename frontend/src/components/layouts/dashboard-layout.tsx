import { Sidebar } from '@/components/sidebar';
import { Outlet } from 'react-router-dom';

export function DashboardLayout() {
	return (
		<div className='h-full relative bg-background'>
			<div className='hidden h-full md:flex md:w-64 md:flex-col md:fixed md:inset-y-0 z-80'>
				<Sidebar />
			</div>
			<main className='md:pl-64'>
				<div className='bg-background border shadow-sm'>
					<Outlet />
				</div>
			</main>
		</div>
	);
}
