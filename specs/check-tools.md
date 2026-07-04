# Spec: tooling dependency version check

## Problem

The project depends on a spread of developer tooling — the Go
toolchain, golangci-lint, Node (for GitNexus and npx-launched MCP
servers), pipx + zensical (site builds), the `claude` CLI with
registered MCP servers (gitnexus, gemini-search) — but nothing
verifies any of it. Gaps surface mid-task instead of up front:
a doc-link fix session discovered zensical *and* pipx were both
missing only after the source edits were done and the site needed
rebuilding.

`make gitnexus-version` covers exactly one tool. There is no
single command that answers "is this machine ready to work on
ctx?"

## Design

Three pieces:

1. **`hack/tool-versions.txt`** — the manifest. One line per
   tool: name, type, requirement, minimum version. Version
   minimums live here, not in script logic, so bumps are
   one-line diffs.
2. **`hack/check-tools.sh`** — the checker. Reads the manifest,
   probes each tool, prints one aligned table with a verdict per
   row, exits non-zero only if a *required* tool fails.
3. **`make check-tools`** — the entry point developers actually
   remember.

### Tool types

| Type  | Probe                                                       |
|-------|-------------------------------------------------------------|
| `bin` | `command -v` + version extraction (`--version` / `go version`) |
| `npm` | installed version vs `npm view <pkg> version` (drift check) |
| `mcp` | registered in `claude mcp list` (config presence)           |

### Verdicts

- `OK` — present, and ≥ minimum when a minimum is declared
- `MISSING` — not found on PATH / not registered
- `OUTDATED` — below declared minimum, or (npm type) behind the
  registry's latest
- `SKIP` — probe impossible (no `claude` CLI for mcp rows, no
  network for npm registry) — reported, never fatal

### Requirement levels

`required` failures exit 1 (go, git, golangci-lint — the
build/lint core). `optional` failures print as warnings and do
not affect the exit code. `--strict` promotes optional failures
to fatal for CI-style use.

### Decisions

- **MCP checks are config-presence, not connectivity.** `claude
  mcp list` already health-checks; parsing its status marks is
  brittle across CLI versions. Presence answers the actual
  question ("did I register it?") deterministically.
- **Network probes are best-effort.** `npm view` runs under a
  timeout and degrades to `SKIP`-style notes; a checker that
  hangs offline would never get run.
- **FireCrawl MCP ships as a commented manifest line** — enabling
  it later is a one-character diff, per the original request.

## Non-Goals

- Auto-installing or upgrading anything (that stays with
  `make site-setup`, `make gitnexus-update`, `make register-mcp`)
- Verifying MCP server *liveness*
- Gating CI on optional tooling

## Acceptance

- [ ] `make check-tools` prints one row per manifest entry with
      verdict, version, and note
- [ ] Exit 0 on a machine with go/git/golangci-lint present, even
      when optional tools are missing
- [ ] Exit 1 when a required tool is missing or below minimum
- [ ] Manifest comment lines and blank lines are ignored
- [ ] Absence of `claude`, `npm`, or network yields SKIP notes,
      never a hang or crash
