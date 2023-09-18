import { ErrorBoundary } from '@/components/error-boundary';
import { DashboardLayout } from '@/components/layouts/dashboard-layout';
import { DashboardPage } from '@/pages/dashboard/page';
import { KeywordResearchPage } from '@/pages/keyword-research/page';
import {
	Route,
	RouterProvider,
	createBrowserRouter,
	createRoutesFromElements
} from 'react-router-dom';
import { Toaster } from 'sonner';
import { ContentEditorPage } from './pages/content-editor/page';
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
			<Route path='projects' element={<ProjectsPage />}>
				<Route path=':projectId' element={<ProjectsPage />} />
			</Route>
			<Route path='content-editor' element={<ContentEditorPage />}>
				<Route path=':contentId' element={<ContentEditorPage />} />
			</Route>
			<Route path='keyword-research' element={<KeywordResearchPage />}>
				<Route path=':researchId' element={<KeywordResearchPage />} />
			</Route>
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
