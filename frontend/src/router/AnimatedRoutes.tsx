import { AnimatePresence } from "framer-motion";
import { ElementType } from "react";
import { Route, Routes, useLocation } from "react-router-dom";
import { routes } from "./routes";
import PrivateRoute from "@/components/PrivateRoute/PrivateRoute";
const AnimatedRoutes = () => {
	const location = useLocation();
	return (
		<AnimatePresence mode="wait" initial={false}>
			<Routes location={location} key={location.pathname}>
				{routes.map((route: { path: string; component: ElementType; isPrivate?: boolean }) => {
					const Component = route.component;
					return (
						<Route
							key={route.path}
							path={route.path}
							element={
								route.isPrivate ? (
									<PrivateRoute>
										<Component />
									</PrivateRoute>
								) : (
									<Component />
								)
							}
						/>
					);
				})}
			</Routes>
		</AnimatePresence>
	);
};
export default AnimatedRoutes;
