# `tools/typecheck/pi/`

Type-check gate for the embedded Pi extension.

## What this is

The Pi extension source at
`internal/assets/integrations/pi/extension/ctx.ts` is shipped inside
the ctx binary via `//go:embed` and deployed to the user's
`.pi/extensions/ctx.ts` at install time (see
`internal/assets/README.md` for the embed contract).

Without a type-check gate, a typo or a drift from the
`@earendil-works/pi-coding-agent` extension API would ship as bytes
and fail only when Pi loads the extension on a user's machine.

This directory holds the tooling that gates that risk:

- `package.json`: declares dependencies on
  `@earendil-works/pi-coding-agent` (for the `ExtensionAPI` /
  `ExtensionContext` types), `@types/node` (for the
  `node:child_process` import), and `typescript`.
- `tsconfig.json`: `noEmit: true`, strict, with `include`
  pointing at the embedded TS file via relative path.
- `package-lock.json`: committed; pinned for reproducibility.

The directory sits **outside** `internal/assets/` deliberately: it
is *about* the embedded payload, not part of it. If it lived
alongside the `.ts` source, it would either bloat the embed (the
file-by-file `//go:embed` directive does not currently glob this
dir, but the proximity is misleading) or invite the question every
time.

## Run locally

```sh
cd tools/typecheck/pi
npm ci               # or: bun install
npx tsc --noEmit     # or: bunx tsc --noEmit
```

Either toolchain works; `tsc` is the same compiler under both.
CI uses `npm ci` (matching the `editors/vscode/` and
`tools/typecheck/opencode/` conventions and the committed
`package-lock.json`); local contributors may use whichever they
have installed. The Pi extension itself runs under Bun at the
consumer's machine, but the type-check tool only needs `tsc`, the
`@earendil-works/pi-coding-agent` type declarations, and
`@types/node` for the Node built-ins.

## What this does **not** check

- **Runtime behavior:** `tsc --noEmit` is a *static* check. It
  catches type errors, not logic bugs. Runtime issues still
  surface only when Pi loads the deployed extension.
- **The other embedded TypeScript assets:** the OpenCode plugin is
  covered by `tools/typecheck/opencode/`. If new `.ts` assets are
  added under `internal/assets/integrations/pi/`, extend the
  `include` glob in `tsconfig.json` to cover them.
- **Embed coverage:** that lives in `internal/assets/embed_test.go`.
  The two checks are complementary: this verifies the bytes are
  valid TypeScript; that verifies the bytes are actually
  embedded.

## Maintenance

Bump `@earendil-works/pi-coding-agent` when Pi releases a version
that changes extension event or type signatures. The extension
source itself documents the Pi API surface it targets in its own
header comment. Keep that comment and this dependency in sync.
