import react from "@vitejs/plugin-react";

// https://vite.dev/config/
export default {
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
};
