import AnimatedRoutes from "@/router/AnimatedRoutes";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { BrowserRouter } from "react-router-dom";
const queryClient = new QueryClient();
function App() {
	return (
		<QueryClientProvider client={queryClient}>
			<BrowserRouter>
				<AnimatedRoutes />
			</BrowserRouter>
		</QueryClientProvider>
	);
}

export default App;
