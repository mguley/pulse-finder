import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  root: "public",
  plugins: [react()],
  server: {
    port: 3000,
    host: true, // This will bind the server to 0.0.0.0
  },
});
