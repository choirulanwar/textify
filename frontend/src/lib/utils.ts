import { clsx, type ClassValue } from 'clsx';
import { toast } from 'sonner';
import { twMerge } from 'tailwind-merge';
import { z } from 'zod';

export function cn(...inputs: ClassValue[]) {
	return twMerge(clsx(inputs));
}

export function truncate(str: string, length: number) {
	return str.length > length ? `${str.substring(0, length)}...` : str;
}

export function catchError(err: unknown) {
	if (err instanceof z.ZodError) {
		const errors = err.issues.map(issue => {
			return issue.message;
		});
		return toast(errors.join('\n'));
	} else if (err instanceof Error) {
		return toast(err.message);
	} else {
		return toast('Something went wrong, please try again later.');
	}
}

export type TrendsData = {
	time: number;
	value: number;
};

export function convertToMonthlyData(trendsData: TrendsData[]): TrendsData[] {
	const monthlyData: TrendsData[] = [];
	const months: { [month: string]: number[] } = {};

	for (const data of trendsData) {
		const date = new Date(data.time * 1000);
		const month = `${date.getMonth() + 1}-${date.getFullYear()}`;

		if (!months[month]) {
			months[month] = [];
		}
		months[month].push(data.value);
	}

	for (const month in months) {
		const values = months[month];
		const average =
			values.reduce((sum, value) => sum + value, 0) / values.length;

		const [m, y] = month.split('-');
		const unixTime = new Date(`${y}-${m}-01`).getTime() / 1000;

		monthlyData.push({ time: unixTime, value: average });
	}

	return monthlyData;
}
