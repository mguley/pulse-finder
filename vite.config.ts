import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  root: "public",
  plugins: [react()],
  server: {
    port: 3000,
    host: true, // This will bind the server to 0.0.0.0
  },
  build: {
    outDir: "../dist", // This will place the build output in the project root 'dist' folder
    emptyOutDir: true, // Ensures that the output directory is cleared before building
  },
  base: "/pulse-finder/", // Set the base path for GitHub Pages
});
