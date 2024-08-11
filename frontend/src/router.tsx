import { createBrowserRouter } from "react-router-dom";
import App from "./App";
import HomePage from "./pages/home";
import LoginPage from "./pages/login";

export const router = createBrowserRouter([
	{
		path: "/",
		element: <HomePage />,
	},
	{
		path: "/login",
		element: <LoginPage />,
	},
	{
		path: "*",
		element: <App />,
	},
]);
