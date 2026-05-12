# `tools/typecheck/opencode/`

Type-check gate for the embedded OpenCode plugin.

## What this is

The OpenCode plugin source at
`internal/assets/integrations/opencode/plugin/index.ts` is shipped
inside the ctx binary via `//go:embed` and deployed to the user's
`.opencode/plugins/ctx.ts` at install time (see
`internal/assets/README.md` for the embed contract).

Without a type-check gate, a typo or a drift from the
`@opencode-ai/plugin` SDK would ship as bytes and fail only when
Bun loads the plugin on a user's machine.

This directory holds the tooling that gates that risk:

- `package.json`: declares dependencies on `@opencode-ai/plugin`
  (for the `Plugin` type), `@types/bun` (for `Bun.$` / BunShell
  globals), and `typescript`.
- `tsconfig.json`: `noEmit: true`, strict, with `include`
  pointing at the embedded TS file via relative path.
- `package-lock.json`: committed; pinned for reproducibility.

The directory sits **outside** `internal/assets/` deliberately: it
is *about* the embedded payload, not part of it. If it lived
alongside the `.ts` source, it would either bloat the embed (the
file-by-file `//go:embed` directive does not currently glob this
dir, but the proximity is misleading) or invite the question
every time.

## Run locally

```sh
cd tools/typecheck/opencode
npm ci               # or: bun install
npx tsc --noEmit     # or: bunx tsc --noEmit
```

Either toolchain works; `tsc` is the same compiler under both.
CI uses `npm ci` (matching the `editors/vscode/` convention and
the committed `package-lock.json`); local contributors may use
whichever they have installed. The OpenCode plugin itself runs
under Bun at the consumer's machine, but the type-check tool
only needs `tsc` and the `@types/bun` declarations.

## What this does **not** check

- **Runtime behavior:** `tsc --noEmit` is a *static* check. It
  catches type errors, not logic bugs. Runtime issues still
  surface only when OpenCode loads the deployed plugin.
- **The other embedded TypeScript assets:** there are none
  today. If new `.ts` assets are added to `internal/assets/`,
  extend the `include` glob in `tsconfig.json` to cover them.
- **Embed coverage:** that lives in `internal/assets/embed_test.go`.
  The two checks are complementary: this verifies the bytes are
  valid TypeScript; that verifies the bytes are actually
  embedded.

## Maintenance

Bump `@opencode-ai/plugin` when the OpenCode SDK releases a
version that changes hook signatures. The plugin source itself
documents the SDK version it targets in its own header comment.
Keep that comment and this dependency in sync.
