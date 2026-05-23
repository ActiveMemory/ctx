// ESLint flat config for the ctx VS Code extension.
// Correctness-only ruleset: no style rules, no formatting
// preferences. The intent is to catch real bugs (unused
// imports, ts-ignore drift, unreachable code) without
// bikeshedding layout. Editor + Prettier-ish defaults
// handle that.
//
// `no-explicit-any` is intentionally disabled — the VS Code
// API surface (and our vitest mocks of it) leans on `any`
// extensively as a deliberate type-loosening for boundary
// shims. Tightening this would conflate eslint with a
// codebase-wide refactor.

const js = require("@eslint/js");
const tseslint = require("typescript-eslint");

module.exports = tseslint.config(
  {
    ignores: ["dist/**", "node_modules/**", "*.vsix"],
  },
  js.configs.recommended,
  ...tseslint.configs.recommended,
  {
    rules: {
      "@typescript-eslint/no-explicit-any": "off",
      "@typescript-eslint/no-unused-vars": [
        "error",
        {
          argsIgnorePattern: "^_",
          varsIgnorePattern: "^_",
          caughtErrorsIgnorePattern: "^_",
        },
      ],
    },
  },
);
