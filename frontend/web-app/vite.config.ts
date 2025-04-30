import path from "path";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";

export default defineConfig({
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
    },
  },
  server: {
    allowedHosts: [".trycloudflare.com"],
    proxy: {
      "/api": {
        target: "http://localhost:8080",
        changeOrigin: true,
      },
    },
  },
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          react: ["react", "react-dom", "react-router"],

          data: ["@tanstack/react-query"],

          mantine: [
            "@mantine/core",
            "@mantine/dates",
            "@mantine/form",
            "@mantine/hooks",
            "@tabler/core",
            "@tabler/icons-react",
            "postcss-preset-mantine",
          ],

          dates: ["date-fns", "date-fns-tz", "dayjs"],

          devtools: [
            "@eslint/js",
            "eslint",
            "prettier",
            "@ianvs/prettier-plugin-sort-imports",
            "prettier-plugin-organize-imports",
          ],
        },
      },
    },
  },
});
