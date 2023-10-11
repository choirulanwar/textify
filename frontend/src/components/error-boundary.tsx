import { isRouteErrorResponse, useRouteError } from 'react-router-dom';

export function ErrorBoundary() {
	const error = useRouteError() as Error;

	if (!isRouteErrorResponse(error)) {
		return null;
	}

	return (
		<section className='flex justify-center items-center w-full h-screen'>
			<h1>Error</h1>

			<small>{error?.message}</small>
		</section>
	);
}
