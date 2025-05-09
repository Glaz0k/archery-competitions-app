import pluginQuery from "@tanstack/eslint-plugin-query";
import mantine from "eslint-config-mantine";
import prettierConfig from "eslint-config-prettier";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import tseslint from "typescript-eslint";

export default tseslint.config(
  ...mantine,
  ...pluginQuery.configs["flat/recommended"],
  { ignores: ["dist"] },
  {
    files: ["**/*.{js,jsx}", "**/*.{mjs,cjs,js,d.ts,d.mts}"],
    plugins: {
      "react-hooks": reactHooks,
      "react-refresh": reactRefresh,
    },
    rules: { "no-console": "off" },
  },
  prettierConfig
);
