# Fix VS Code Extension Tests + Re-Enable in CI

The vitest suite at `editors/vscode/src/extension.test.ts`
was excluded from CI's TypeScript typecheck (via
`tsconfig.ci.json`'s `**/*.test.ts` exclude) and never run
under `npm test`. Two distinct breakages had accumulated.

## Problem

**Bug 1 — handler name drift.** The test suite imported
`handleComplete` and `handleTasks` from `extension.ts`. These
handlers were merged into a single dispatching `handleTask`
that takes a subcommand string (`complete <ref>`, `archive`,
`snapshot [name]`). The test file kept the old imports and
call shapes. TypeScript would have caught this at strict
typecheck time, but the carve-out in `tsconfig.ci.json`
silenced it.

**Bug 2 — fakeToken signature.** The `fakeToken` helper
declared `onCancellationRequested: vi.fn((cb: () => void) =>
…)`. VS Code's `CancellationToken.onCancellationRequested`
is typed `Event<any>`, whose listener contract is
`(e: any) => any`. The zero-arg listener trips strict TS as
an under-arity assignment.

**Bug 3 (latent, surfaced once the suite ran) — every
handler also gained a trailing `args.push("--no-color")`
after the original test was written.** 18 assertions across
`handleRemind`, `handlePad`, `handleNotify`, `handleSystem`,
and the new `handleTask` shape did not include
`"--no-color"` in their expected argv arrays. None of these
fired until the suite was actually executed under vitest —
they all passed the carve-out typecheck because
`expect.anything()` and `expect.any(Function)` admit any
shape. This is the cost of "we'll fix the tests later" —
the underlying drift accreted across multiple unrelated
commits.

## Solution

1. Drop the `**/*.test.ts` carve-out in
   `editors/vscode/tsconfig.ci.json`. The CI typecheck now
   covers the test file too.
2. Rewrite the imports + 4 `handleComplete` cases as
   `handleTask(..., "complete <ref>", ...)`; rewrite the 5
   `handleTasks` cases as `handleTask(..., "archive"|"snapshot
   [name]", ...)`. Rename the two describe blocks
   accordingly (`handleTask complete`, `handleTask
   archive/snapshot`). The `metadata.command` value is now
   `"task"` for both (the merged handler emits a single
   command tag).
3. Widen `fakeToken`'s listener type to match `Event<any>`:
   `type Listener = (e: unknown) => unknown`. The `_fire`
   helper now passes `undefined` to mirror VS Code's "no
   payload" cancellation event.
4. Append `"--no-color"` to every argv assertion that
   exercises a handler. Specific lines: see commit diff.
5. Extend the CI `vscode-extension` job
   (`.github/workflows/ci.yml:92`):
   - `npm test` (vitest)
   - `npx vsce package --no-dependencies --out
     /tmp/ctx-extension.vsix` (catches the marketplace-side
     packaging contract — bad manifest, missing files, etc.
     — without publishing).

## Out of Scope

- `npm run lint`. The script exists in `package.json` but no
  `.eslintrc*` config is checked in, so the script crashes.
  Scaffolding a typescript-eslint config is a separate
  decision (which preset, which rules, style vs. correctness
  only) and lands as a follow-up task. Logged inline in
  `.context/TASKS.md` alongside the closeout for this item.
- Splitting `extension.test.ts` (702 lines) into per-handler
  files. The current monolith works; restructuring would
  conflate this fix with refactor noise.

## Verification

- `npx tsc --noEmit -p tsconfig.ci.json` from
  `editors/vscode/` — passes (test file now included).
- `npx vitest run` — 53 / 53 pass.
- `npx vsce package --no-dependencies --out
  /tmp/ctx-extension-dryrun.vsix` — produces a 9-file,
  26.65 KB vsix.
