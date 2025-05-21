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
  },
  plugins: [react()],
  build: {
    rollupOptions: {
      output: {
        manualChunks: {
          react: ["react", "react-dom", "react-router"],

          api: ["@tanstack/react-query", "axios", "zod", "mantine-form-zod-resolver"],

          mantine: [
            "@mantine/core",
            "@mantine/dates",
            "@mantine/form",
            "@mantine/hooks",
            "@mantine/notifications",
            "postcss-preset-mantine",
          ],

          icons: ["@tabler/core", "@tabler/icons-react"],

          dates: ["date-fns", "date-fns-tz", "dayjs"],
        },
      },
    },
  },
});
