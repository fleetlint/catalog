// ESLint flat config (see https://fleetlint.org/stacks/#nodejs--typescript). Written by fleetlint fix; adjust freely.
// Dependencies: eslint, typescript-eslint, @eslint-community/eslint-plugin-eslint-comments, eslint-plugin-unicorn.
import eslint from "@eslint/js";
import tseslint from "typescript-eslint";
import comments from "@eslint-community/eslint-plugin-eslint-comments/configs";
import unicorn from "eslint-plugin-unicorn";

export default tseslint.config(
  { ignores: ["dist/", "coverage/", "node_modules/"] },
  eslint.configs.recommended,
  ...tseslint.configs.strictTypeChecked,
  ...tseslint.configs.stylisticTypeChecked,
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
      "@typescript-eslint/explicit-module-boundary-types": "error",
      "@typescript-eslint/no-floating-promises": "error",
      "complexity": ["error", 15],
      "max-lines-per-function": ["error", { max: 50, skipBlankLines: true, skipComments: true }],
    },
  },
);
