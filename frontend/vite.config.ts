import react from "@vitejs/plugin-react";
import path from "path";
import { defineConfig } from "vite";
export default defineConfig({
	plugins: [react()],
	optimizeDeps: {
		include: ["mapbox-gl"],
	},
	server: {
		proxy: {
			"/api": "http://localhost:3001",
		},
	},
	resolve: {
		alias: {
			"@": path.resolve(__dirname, "./src"),
		},
	},
});
