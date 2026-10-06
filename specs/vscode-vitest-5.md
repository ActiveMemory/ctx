# VS Code Extension: vitest 5

Dependabot PR #170 bumps the extension's `vitest` from `^4.0.18`
to `^5.0.0`. On its own the bump fails CI at `npm ci`, before any
code is built. This spec records the companion changes the bump
needs, so it can land as one commit.

## Problem

vitest 5 raises its floor in two places the extension did not
meet:

1. **Peer range.** vitest 5 declares `@types/node` as an optional
   peer at `^22.0.0 || >=24.0.0`. The extension pins
   `@types/node@^20.0.0`, and npm's resolver rejects the tree with
   `ERESOLVE` (optional peers still conflict when present).
2. **Runtime.** vitest 5 requires Node.js 22. The CI
   `vscode-extension` job ran on `node-version: '20'`.

Vite 8, pulled in by vitest 5, also warns that
`vitest.config.ts` uses ESM syntax in a file loaded as CommonJS
(the package has no `"type": "module"`), and says a future Vite
major will make that an error.

## Solution

- `editors/vscode/package.json`: `@types/node` `^20.0.0` →
  `^22.0.0`; lockfile regenerated.
- `.github/workflows/ci.yml`: the `vscode-extension` job's
  `node-version` `'20'` → `'22'`. The other Node jobs (OpenCode
  plugin, ctx-desktop) are unaffected by this bump and stay on 20.
- Rename `editors/vscode/vitest.config.ts` →
  `vitest.config.mts`, so it is loaded as ESM and the warning is
  gone. `.vscodeignore` named the old file explicitly; its
  `**/*.ts` glob does not match `.mts`, so the entry is renamed
  too to keep the config out of the `.vsix`.

## Trade-off

`engines.vscode: ^1.93.0` means the extension can run on a VS Code
whose Electron ships Node 20. With `@types/node@22`, `tsc` no
longer flags a Node-22-only API used in `src/`. The extension's
Node surface is `child_process`, `fs`, `https`, `os`, and `path`,
all long-stable modules, so the gap is accepted rather than
pinning vitest to 4.

## Verification

Under Node 22, from a clean `npm ci` in `editors/vscode`, every
step of the CI job passes: `npm run build`,
`npx tsc --noEmit -p tsconfig.ci.json`, `npm run lint`,
`npm test` (53/53, no Vite config warning), and
`npx vsce package --no-dependencies` (10 files, config excluded).

## See Also

- `specs/fix-vscode-extension-tests.md`: the spec that put
  `npm test` and the vsce dry-run into the CI job.
