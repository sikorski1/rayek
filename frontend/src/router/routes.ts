import Home from "@/pages/Home/Home";
import Info from "@/pages/Info/Info";
import Login from "@/pages/Login/Login";
import Maps from "@/pages/Maps/Maps";
import SingleMap from "@/pages/SingleMap/SingleMap";
export const routes = [
	{
		path: "/",
		component: Home,
		isPrivate: false,
	},
	{
		path: "/maps",
		component: Maps,
		isPrivate: true,
	},
	{
		path: "/info",
		component: Info,
		isPrivate: false,
	},
	{
		path: "/maps/:id",
		component: SingleMap,
		isPrivate: true,
	},
	{
		path: "/login",
		component: Login,
		isPrivate: false,
	},
];
