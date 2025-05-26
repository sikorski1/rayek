import { ElementType } from "react";
import { Route, Routes, useLocation } from "react-router-dom";
import { AnimatePresence } from "framer-motion";
import { routes } from "./routes";
const AnimatedRoutes = () => {
	const location = useLocation();
	return (
		<AnimatePresence mode="wait" initial={false}>
			<Routes location={location} key={location.pathname}>
				{routes?.map((route: { path: string; component: ElementType }) => (
					<Route path={route.path} element={<route.component />} key={route.path} />
				))}
			</Routes>
		</AnimatePresence>
	);
};
export default AnimatedRoutes;
