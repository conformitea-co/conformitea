import { defineConfig, searchForWorkspaceRoot } from "vite";
import tailwindcss from "@tailwindcss/vite";
import { reactRouter } from "@react-router/dev/vite";
import path from "path";

export default defineConfig({
  plugins: [tailwindcss(), reactRouter()],
  server: {
    fs: {
      allow: [searchForWorkspaceRoot(process.cwd()), "../node_modules"],
    },
  },
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
});
