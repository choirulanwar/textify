import { cn } from '@/lib/utils';
import {
	ArchiveIcon,
	BarChartIcon,
	DashboardIcon,
	GearIcon,
	Pencil1Icon
} from '@radix-ui/react-icons';
import { Link, useLocation } from 'react-router-dom';
import { Separator } from './ui/separator';

const routes = [
	{
		label: 'Dashboard',
		icon: DashboardIcon,
		href: '/dashboard'
	},
	{
		label: 'Projects',
		icon: ArchiveIcon,
		href: '/projects'
	},
	{
		label: 'Content Editor',
		icon: Pencil1Icon,
		href: '/content-editor'
	},
	{
		label: 'Keyword Research',
		icon: BarChartIcon,
		href: '/keyword-research'
	}
];

export function Sidebar() {
	const { pathname } = useLocation();

	return (
		<div className='flex pt-4 flex-col h-full bg-secondary justify-between'>
			<div className='flex-1 space-y-4'>
				<Link to='/dashboard' className='flex items-center pl-4'>
					<h1 className={cn('text-lg font-bold text-center')}>Textify</h1>
				</Link>

				<div className='space-y-1'>
					{routes.map(route => (
						<Link
							key={route.href}
							to={route.href}
							className={cn(
								'text-sm group flex py-2 w-full justify-start cursor-pointer hover:text-primary hover:bg-muted transition',
								pathname === route.href
									? 'text-primary bg-muted font-semibold'
									: 'text-foreground font-medium'
							)}
						>
							<div className='flex items-center flex-1 px-4'>
								<route.icon className={cn('h-4 w-4 mr-3')} />
								{route.label}
							</div>
						</Link>
					))}
				</div>
			</div>

			<Separator />

			<Link
				key={'settings'}
				to={'/settings'}
				className={cn(
					'text-sm group flex py-2 w-full justify-start cursor-pointer hover:text-primary hover:bg-muted transition',
					pathname === '/settings'
						? 'text-primary bg-muted font-semibold'
						: 'text-foreground font-medium'
				)}
			>
				<div className='flex items-center flex-1 px-4'>
					<GearIcon className={cn('h-4 w-4 mr-3')} />
					Settings
				</div>
			</Link>
		</div>
	);
}
