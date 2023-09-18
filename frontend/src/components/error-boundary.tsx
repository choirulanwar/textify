import { isRouteErrorResponse, useRouteError } from 'react-router-dom';

export function ErrorBoundary() {
	const error = useRouteError() as Error;

	if (!isRouteErrorResponse(error)) {
		return null;
	}

	return (
		<section>
			<h1>Error Boundary</h1>

			<small>{error?.message}</small>
		</section>
	);
}
