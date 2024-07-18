import ReactDOM from 'react-dom/client'
import App from './App.tsx'
import './index.css'
import { ConfigProvider } from './config.tsx'

ReactDOM.createRoot(document.getElementById('root')!).render(
	<ConfigProvider >
		<App />
	</ConfigProvider>
)
