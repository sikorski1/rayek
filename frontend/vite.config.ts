import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import path from "path";
export default defineConfig({
	plugins: [react()],
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
