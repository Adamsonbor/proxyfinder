import './App.css';

import { RouterProvider } from 'react-router-dom';
import { router } from './router';
import { ConfigProvider } from './config';

export default function App() {
	return (
		<ConfigProvider>
			<RouterProvider router={router} />
		</ConfigProvider>
	);
}
