import ReactHeatMap from "@/components/ReactHeatMap/ReactHeatMap"
import Home from "@/pages/Home/Home"
import Maps from "@/pages/Maps/Maps"
import SingleMap from "@/pages/SingleMap/SingleMap"
export const routes = [
    {
      path: "/",
      component: Home,
      isPrivate: false,
    },
    {
      path: "/maps",
      component: Maps,
      isPrivate: false,
    },
    {
      path: "/maps/:id",
      component: SingleMap,
      isPrivate: false,
    },
    {
      path:"/sampleMap",
      component: ReactHeatMap,
      isPrivate: false,
    }
]