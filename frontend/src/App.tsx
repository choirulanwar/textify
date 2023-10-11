import { ErrorBoundary } from '@/components/error-boundary';
import { DashboardLayout } from '@/components/layouts/dashboard-layout';
import { DashboardPage } from '@/pages/dashboard/page';
import { KeywordTrendExplorerPage } from '@/pages/keyword-trend-explorer/[id]/page';
import { KeywordTrendExplorersPage } from '@/pages/keyword-trend-explorer/page';
import {
	Route,
	RouterProvider,
	createBrowserRouter,
	createRoutesFromElements
} from 'react-router-dom';
import { Toaster } from 'sonner';
import { ContentEditorPage } from './pages/content-editor/page';
import { DomainFinderPage } from './pages/domain-finder/page';
import { ProjectsPage } from './pages/projects/page';
import { SettingsPage } from './pages/settings/page';

const router = createBrowserRouter(
	createRoutesFromElements(
		<Route
			path='/'
			element={<DashboardLayout />}
			errorElement={<ErrorBoundary />}
		>
			<Route path='dashboard' index element={<DashboardPage />} />
			<Route path='project' element={<ProjectsPage />}>
				<Route path=':id' element={<ProjectsPage />} />
			</Route>
			<Route path='content-editor' element={<ContentEditorPage />}>
				<Route path=':id' element={<ContentEditorPage />} />
			</Route>
			<Route path='domain-finder' element={<DomainFinderPage />}>
				<Route path=':id' element={<DomainFinderPage />} />
			</Route>
			<Route
				path='keyword-trend-explorer'
				element={<KeywordTrendExplorersPage />}
			/>
			<Route
				path='keyword-trend-explorer/:id'
				element={<KeywordTrendExplorerPage />}
			/>
			<Route path='settings' element={<SettingsPage />} />
		</Route>
	)
);

export default function App() {
	return (
		<>
			<Toaster />
			<RouterProvider router={router} />
		</>
	);
}
