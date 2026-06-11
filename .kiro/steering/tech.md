---
name: tech
description: Technology stack, constraints, and dependencies
mode: always
---

# Technology Stack

## Primary

- **Go 1.26+**, statically linked (`CGO_ENABLED=0`). The `ctx`
  binary is the entire deliverable for the core; everything else
  ships as embedded bytes inside it.
- **Cobra** for the CLI command surface.
- **`embed.FS`** for shipping foreign-language assets (TypeScript,
  Bash, PowerShell, Markdown, JSON, YAML) inside the Go binary.
  See `internal/assets/README.md` for the embed contract; the
  hard `//go:embed` no-`../` rule shapes the directory layout.

## Separately-published

- **VS Code extension** at `editors/vscode/` ships as a `.vsix`
  to the VS Code Marketplace under publisher `activememory`. It
  is NOT embedded; it has its own `package.json`, `tsconfig.json`,
  and CI guardrails (`vscode-extension` job).
- The embedded **OpenCode plugin** at
  `internal/assets/integrations/opencode/plugin/index.ts` has its
  type-check tooling outside the embed tree at
  `tools/typecheck/opencode/`.

## Hard constraints

- **No runtime dependencies.** No package manager, no network
  fetch on install. If a feature needs a daemon or a service,
  it's the wrong feature.
- **No CGO.** Build must succeed with `CGO_ENABLED=0` on every
  supported platform (Linux/macOS/Windows × amd64/arm64).
- **No network calls during normal operation.** Tests included.
  Operations that genuinely need network (e.g. GitHub release
  download in the VS Code extension auto-bootstrap) are scoped
  and opt-in.
- **Foreign-language assets ship embedded, not at install time.**
  TypeScript / Bash / PowerShell that integrates with external
  tools is baked into the Go binary at compile time and written
  out to the user's filesystem by `ctx setup <tool>`.

## Companion tooling

- **GitNexus** (`mcp__gitnexus__*`) — code intelligence MCP
  server for impact analysis, route maps, and shape checks.
- **Gemini Search** — preferred over built-in web search for
  faster, more accurate results.

## Build / test / lint

- `make build`, `make test`, `make lint` are the canonical
  entrypoints. CI runs the same.
- `make site` rebuilds `site/` from `docs/` via zensical.
- The TS type-check for embedded OpenCode plugin lives at
  `tools/typecheck/opencode/`; `npx tsc --noEmit` is the gate.
- The VS Code extension gate runs `npm ci && npm run build &&
  npx tsc --noEmit -p tsconfig.ci.json` in CI.
