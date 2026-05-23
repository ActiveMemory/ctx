# Scaffold ESLint for editors/vscode/

The `lint` script in `editors/vscode/package.json` referenced
`eslint src --ext ts` but no `.eslintrc*` or `eslint.config.js`
was checked in, so the script crashed today. The previous
task (`specs/fix-vscode-extension-tests.md`) deferred this as
out-of-scope and logged it as a follow-up.

## Problem

A `lint` script that crashes is worse than no script — it
gives the impression that linting *exists* and might be
caught by anyone running `npm run lint` locally. The
follow-up was tagged `#priority:low` but doing it now,
immediately after re-enabling the test suite, avoids a
second round of latent drift.

## Decisions (taken in dialog with user)

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Preset | `@typescript-eslint/recommended` | Fast (no type info), no parser-graph false positives, adequate for a tiny TS surface. |
| Style rules | None | Correctness only. No quotes/semicolons/indent. Editor handles layout. |
| Config scope | Per-package | `editors/vscode/eslint.config.js` only. `tools/typecheck/opencode/` is a separate decision. |
| Rule override | `no-explicit-any: off` | The VS Code API shim and vitest mocks lean on `any` deliberately; enforcing it would conflate eslint with a codebase-wide refactor. |
| Rule override | `no-unused-vars: error` with `^_` ignore pattern | Catches dead imports / leftover params. Underscore prefix is the documented opt-out. |

## Solution

1. Added devDeps: `eslint@^9` + `typescript-eslint@^8` (the
   combined umbrella package; supersedes the separate
   `@typescript-eslint/parser` + `@typescript-eslint/eslint-plugin`
   packages under flat config).
2. Created `editors/vscode/eslint.config.js` (flat config —
   ESLint 9's default since ESLint 8 went EOL Oct 2024).
   Composes `js.configs.recommended` +
   `tseslint.configs.recommended` + the two rule overrides.
3. Updated the `lint` script: `eslint src --ext ts` →
   `eslint src`. Flat config infers extensions from the
   matched files; the legacy `--ext` flag is rejected by
   ESLint 9.
4. First-run violation:
   `extension.ts:278  let disposable …  prefer-const`.
   Attempted a `const disposable = token?.onCancellationRequested(…)`
   refactor — declared after `child = execFile(…)` so the
   listener can close over `child`. Real Node defers
   `execFile` callbacks to `process.nextTick`, but vitest's
   mock fires them **synchronously**, which hit a TDZ on
   `disposable` from inside the callback. Reverted to `let`
   with a `// eslint-disable-next-line prefer-const`
   comment explaining the test-mock constraint.
5. Wired `npm run lint` into the `vscode-extension` CI job
   between typecheck and test.

## Verification

- `npm run lint` from `editors/vscode/` — zero errors,
  zero warnings.
- `npx tsc --noEmit -p tsconfig.ci.json` — passes.
- `npx vitest run` — 53/53 pass.
- `npx vsce package --no-dependencies` — produces a clean
  vsix (devDeps not bundled into marketplace package).
- `npm audit` — zero vulnerabilities (one transient
  brace-expansion CVE from the install was auto-fixed by
  `npm audit fix`).

## Out of Scope

- ESLint for `tools/typecheck/opencode/`. That package is
  a CI-only typecheck harness, ~zero application code,
  and would need its own per-package config. Track as
  separate task if desired.
- Tightening `no-explicit-any`. Would require a sweeping
  refactor of the VS Code API mocks; not in scope for
  re-enabling lint.
- Migration to `--max-warnings 0` posture. Currently the
  config produces zero warnings, so there's nothing to
  tighten. If future rules are added as `warn`, revisit
  the CI posture then.
