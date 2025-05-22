import js from "@eslint/js";
import globals from "globals";
import tseslint from "typescript-eslint";
import pluginReact from "eslint-plugin-react";
import { defineConfig } from "eslint/config";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";

export default defineConfig([
  // Main source files (type-aware linting)
  {
    files: ["src/**/*.{js,mjs,cjs,ts,mts,cts,jsx,tsx}"],
    plugins: { js, "react-refresh": reactRefresh, "react-hooks": reactHooks },
    extends: [
      js.configs.recommended,
      ...tseslint.configs.strictTypeChecked,
      pluginReact.configs.flat.recommended,
    ],
    languageOptions: {
      globals: globals.browser,
      parserOptions: {
        project: "./tsconfig.app.json",
      },
    },
    rules: {
      ...reactHooks.configs.recommended.rules,
      "react/react-in-jsx-scope": "off",
      "react-refresh/only-export-components": [
        "warn",
        { allowConstantExport: true },
      ],
    },
  },
  // Config files (basic linting, no type-aware rules)
  {
    files: ["*.config.{js,ts,mjs}"],
    extends: [js.configs.recommended],
  },
]);
