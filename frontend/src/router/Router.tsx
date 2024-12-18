import { ElementType } from "react";
import { BrowserRouter, Route, Routes } from "react-router-dom";
// import { BrowserRouter, Route, Routes } from "react-router-dom";
import { routes } from "./routes";
const Router = () => {
    return (
      <BrowserRouter>
        <Routes>
          {routes?.map((route: { path: string; component: ElementType }) => (
            <Route
              path={route.path}
              element={<route.component />}
              key={route.path}
            />
          ))}
        </Routes>
      </BrowserRouter>
    );
  };
export default Router;