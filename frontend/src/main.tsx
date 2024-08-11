import ReactDOM from 'react-dom/client'
import './index.css'
import { ConfigProvider } from './config.tsx'
import App from './App.tsx'

ReactDOM.createRoot(document.getElementById('root')!).render(
	<ConfigProvider >
		<App />
	</ConfigProvider>
)
