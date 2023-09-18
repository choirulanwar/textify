import { ChevronRightIcon, SlashIcon } from '@radix-ui/react-icons';
import * as React from 'react';
import { Link } from 'react-router-dom';

import { cn, truncate } from '@/lib/utils';

interface BreadcrumbsProps {
	segments: {
		title: string;
		href: string;
	}[];
	separator?: React.ComponentType<{ className?: string }>;
	truncationLength?: number;
}

export function Breadcrumbs({
	segments,
	separator,
	truncationLength = 0
}: BreadcrumbsProps) {
	const SeparatorIcon = separator ?? ChevronRightIcon;

	return (
		<nav
			aria-label='breadcrumbs'
			className='flex items-center text-sm font-medium text-muted-foreground p-4'
		>
			{segments.map((segment, index) => {
				const isLastSegment = index === segments.length - 1;

				return (
					<React.Fragment key={segment.href}>
						<Link
							aria-current={isLastSegment ? 'page' : undefined}
							to={segment.href}
							className={cn(
								'truncate transition-colors hover:text-foreground text-xs',
								isLastSegment ? 'text-foreground' : 'text-muted-foreground'
							)}
						>
							{truncationLength > 0 && segment.title
								? truncate(segment.title, truncationLength)
								: segment.title}
						</Link>
						{!isLastSegment && (
							<SlashIcon className='mx-2 h-4 w-4 text-xs' aria-hidden='true' />
						)}
					</React.Fragment>
				);
			})}
		</nav>
	);
}
