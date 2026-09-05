# Hook Surface Robustness Sweep

## Problem

The Codex integration surfaced a failure *class*, not a one-off: hook
commands that die (or silently no-op) before `ctx` even runs, for
reasons the host turns into per-event failure noise — or worse, into
blocks. An adversarial audit of every hook surface (Claude manifest,
Codex manifests, Copilot CLI manifest + scripts, OpenCode plugin,
trace hook, dev reload script) confirmed 20 defects across seven
classes: pre-ctx anchor/env aborts, shell-dependent constructs,
BSD/GNU differences, `set -e` traps, unquoted expansions,
host-semantics mismatches (exit codes / stdout contracts), and
unbounded stdin reads.

## Fixes

- **Claude manifest** (`internal/assets/claude/hooks/hooks.json`,
  all 22 commands): uniform POSIX prologue —
  `command -v ctx || exit 0` (ctx is documented optional; absence must
  not spam 13 failures per prompt), then
  `[ -d "${CLAUDE_PROJECT_DIR:-}" ] || { echo remedy >&2; exit 1; }`
  (deterministic exit 1 on every shell; the previous `${VAR:?}` abort
  exits 2 under dash, which Claude Code interprets as a hard block),
  then `cd "$CLAUDE_PROJECT_DIR"`. Companions: parity-test anchor
  constant, `specs/cwd-anchored-context.md` hook-contract row.
- **Copilot CLI manifest**: every entry gains the schema-native
  `"cwd": "."` (resolved relative to the repository root), closing the
  same anchor gap fixed in the Codex manifests without embedding
  POSIX-isms in the cross-shell `command` slot.
- **Copilot CLI wrapper scripts**: both shipped script generations
  (`ctx-*` stdin-JSON and hyphenated argv contracts, 16 files) were
  dead code — the manifest invokes `ctx system ...` directly and
  references no script — and carried eight of the confirmed defects
  (`set -e`+jq aborts, unbounded `$(cat)`, string-built JSONL rot,
  cwd-relative audit paths, a deny contract using Claude field names
  on stderr that Copilot ignores, PowerShell mojibake/suppressed
  errors). Removed at the root: assets, embed globs, reader, deployer
  copy loop, constants, README rows.
- **Trace hook** (`prepare-commit-msg.sh`): amends (COMMIT_SOURCE
  `commit` + SHA arg) exit early and a same-key trailer guard makes
  the append idempotent — no more duplicate trailers per amend.
- **`hack/plugin-reload.sh`**: stages into `mktemp -d` and swaps
  atomically, so a mid-build failure can no longer destroy the old
  cache; mirrors the whole plugin root (closing the known `.mcp.json`
  dev-reload gap).
- **OpenCode plugin**: stdin-reading `ctx system` calls invoke with
  `< /dev/null` (a host-held pipe cost the 2 s stdin timeout per
  call); header documents the git-worktree requirement.
- **OpenCode MCP registration**: when `ctx` is not on PATH at setup
  time, fall back to `os.Executable()` instead of writing a bare name
  the OpenCode spawn cannot resolve.

## Deliberate non-fixes

- Kiro/Cursor/Cline MCP configs keep the bare `ctx` command:
  these files are project-scoped and committable, so embedding a
  machine-specific absolute path would break every other machine.
  Documented at each site.
- Copilot manifest commands stay guard-less (`command -v` is
  POSIX-only and the `command` slot is cross-shell); ctx-absent noise
  there requires ctx to have been uninstalled after setup. Tracked as
  a follow-up with the Windows work.

## Verification

Every finding and every fix was reproduced/validated by adversarial
verifier agents (shell matrices across sh/bash 3.2/bash 5/zsh/dash,
scratch marketplaces, stub binaries); the new Claude prologue was
re-verified on all four local shells (unset → exit 1 + remedy; set →
exit 0). Full lint/test/audit gates green.

## Non-Goals

Rewiring the Copilot manifest through scripts; changing ctx's own
Go-side exit-code contracts; Windows `commandWindows`/PowerShell
parity (tracked in TASKS).
