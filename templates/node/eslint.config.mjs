// ESLint flat config (see https://fleetlint.org/stacks/#nodejs--typescript). Written by fleetlint fix; adjust freely.
// Dependencies: eslint, @eslint/js, typescript-eslint, @eslint-community/eslint-plugin-eslint-comments, eslint-plugin-unicorn.
import { defineConfig, globalIgnores } from "eslint/config";
import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import comments from "@eslint-community/eslint-plugin-eslint-comments/configs";
import unicorn from "eslint-plugin-unicorn";

export default defineConfig(
  globalIgnores(["dist/", "coverage/", "node_modules/"]),
  eslint.configs.recommended,
  tseslint.configs.strictTypeChecked,
  tseslint.configs.stylisticTypeChecked,
  comments.recommended,
  unicorn.configs.recommended,
  {
    languageOptions: { parserOptions: { projectService: true, tsconfigRootDir: import.meta.dirname } },
    rules: {
      // Slop guards: a disable needs a reason, a TODO needs an owner or a deadline, no debugging output.
      "@eslint-community/eslint-comments/require-description": "error",
      "@eslint-community/eslint-comments/no-unlimited-disable": "error",
      "unicorn/expiring-todo-comments": ["error", { allowWarningComments: false }],
      "no-console": ["error", { allow: ["error", "warn"] }],
      "no-debugger": "error",
      "@typescript-eslint/no-explicit-any": "error",
      "@typescript-eslint/no-floating-promises": "error",
      "complexity": ["error", 15],
      "max-lines-per-function": ["error", { max: 50, skipBlankLines: true, skipComments: true }],
    },
  },
  {
    files: ["**/*.{ts,mts,cts,tsx}"],
    rules: { "@typescript-eslint/explicit-module-boundary-types": "error" },
  },
  // JavaScript files (this config among them) are outside the TypeScript project: lint them without type information.
  {
    files: ["**/*.{js,mjs,cjs,jsx}"],
    extends: [tseslint.configs.disableTypeChecked],
  },
);
