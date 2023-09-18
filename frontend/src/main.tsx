import { ThemeProvider } from '@/components/theme-provider';
import React from 'react';
import ReactDOM from 'react-dom/client';
import { Toaster } from 'sonner';
import App from './App';
import './styles/globals.css';

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
	<React.StrictMode>
		<Toaster />
		<ThemeProvider>
			<App />
		</ThemeProvider>
	</React.StrictMode>
);
