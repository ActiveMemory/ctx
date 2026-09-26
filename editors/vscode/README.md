```text
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
```

## `ctx`: VS Code Chat Extension

A VS Code Chat Participant that brings [ctx](https://ctx.ist) (persistent
project context for AI coding sessions) directly into GitHub Copilot Chat.

Type `@ctx` in the Chat view for 36 slash commands: 27 run the ctx CLI,
and 9 run a canonical ctx skill (brainstorm, spec, next, wrap-up, ...)
through the chat model, grounded in live ctx output.

## Quick Start

1. Install the extension (or build from source; see [Development](#development))
2. Open a project in VS Code
3. Open Copilot Chat and type `@ctx /init`

The extension auto-downloads the ctx CLI binary if it isn't on your PATH.

## Slash Commands

### CLI-Backed

Each runs the `ctx` command shown and renders its output. A command that
exits non-zero is shown as a failure with the CLI's own message, never
as a normal result.

| Command | Runs |
|---------|------|
| `/init` | `ctx init --caller vscode`, then `ctx setup copilot --write` |
| `/status` | `ctx status` |
| `/agent [--budget N]` | `ctx agent` |
| `/drift` | `ctx drift` |
| `/recall [--limit N]`, `/recall show <id>` | `ctx journal source` |
| `/setup [tool] [preview]` | `ctx setup <tool> --write` (default tool: `copilot`) |
| `/add <type> <text> [flags]` | `ctx task\|decision\|learning\|convention add` |
| `/decision` | `ctx index .context/DECISIONS.md` |
| `/learning` | `ctx index .context/LEARNINGS.md` |
| `/load` | `ctx load` |
| `/compact` | `ctx compact` |
| `/sync` | `ctx sync` |
| `/task complete <ref>\|archive\|snapshot [name]` | `ctx task ...` |
| `/remind [add\|list\|dismiss]` | `ctx remind ...` |
| `/pad [add\|show\|rm\|edit\|mv\|resolve\|import\|export\|merge]` | `ctx pad ...` |
| `/notify test`, `/notify <message> --event <name>` | `ctx hook notify ...` |
| `/system resources\|stats\|bootstrap\|message` | `ctx sysinfo`, `ctx usage`, `ctx system bootstrap`, `ctx hook message ...` |
| `/memory sync\|status\|diff\|import\|publish\|unpublish` | `ctx memory ...` |
| `/journal site\|obsidian` | `ctx journal ...` |
| `/doctor` | `ctx doctor` |
| `/config switch <profile>\|status\|schema` | `ctx config ...` |
| `/why [document]` | `ctx why <document>` (default: `manifesto`) |
| `/change [--since D]` | `ctx change` |
| `/guide [--skills\|--commands]` | `ctx guide` |
| `/permission snapshot\|restore` | `ctx permission ...` |
| `/pause`, `/resume` | `ctx hook pause`, `ctx hook resume` |

`/add` fills in the provenance the CLI requires for tasks, decisions,
and learnings (`--session-id` from the VS Code session, `--branch` and
`--commit` from git). Everything else is passed through, so the
CLI's own rules apply: tasks and conventions need `--section`, decisions
need `--context`, `--rationale`, `--consequence`, and learnings need
`--context`, `--lesson`, `--application`. Quote multi-word values:

```text
@ctx /add decision Use PostgreSQL --context "Need a reliable DB" --rationale "ACID and JSON" --consequence "Ops training"
```

`/notify setup` points you at `ctx hook notify setup` in a terminal: the
webhook URL is a secret and does not belong in the chat history.

### Skill-Backed

`/<name>` runs the canonical `ctx-<name>` skill. The skill text is bundled
from `internal/assets/claude/skills/` at build time, and each request
hands the chat model the skill, the output of `ctx agent` (plus any
read-only ctx output the skill relies on), earlier turns of the same
skill conversation, and files you attach with `#file`.

| Command | Skill | Also reads |
|---------|-------|------------|
| `/brainstorm` | `ctx-brainstorm` | |
| `/spec` | `ctx-spec` | |
| `/implement` | `ctx-implement` | attach the plan with `#file` |
| `/next` | `ctx-next` | `ctx journal source --limit 3` |
| `/remember` | `ctx-remember` | `ctx journal source --limit 3` |
| `/reflect` | `ctx-reflect` | |
| `/wrap-up` | `ctx-wrap-up` | |
| `/blog` | `ctx-blog` | `ctx journal source --limit 10` |
| `/consolidate` | `ctx-consolidate` | `ctx drift --json` |

The model cannot run commands or edit files from here. Where a skill
says to persist something, it gives you the exact `ctx` or `@ctx`
command to run instead, and it never claims to have done it. A plain
reply right after a skill answer continues that skill, so multi-turn
workflows like `/brainstorm` keep their thread. Only skill exchanges
are sent back to the model: output of CLI commands such as `/pad` never
is.

Skills that must explore the repository or run commands on their own
(`ctx-architecture`, `ctx-link-check`, `ctx-worktree`,
`ctx-blog-changelog`) are left to agent integrations where ctx deploys
its skills (Claude Code, `ctx setup copilot-cli`).

## Background Behavior

| Trigger | What Happens |
|---------|--------------|
| **Extension activate** | Fires `ctx system session-event --type start` |
| **`/init` succeeds** | Fires the same session start (activation had no `.context/` yet) |
| **`.context/` file change, and every 5 minutes** | Refreshes the reminder status bar from `ctx remind list` (read-only) |
| **Extension deactivate** | Fires `ctx system session-event --type end` |

## Status Bar

A `$(bell) ctx` indicator appears in the status bar while `ctx remind
list` has pending reminders, and hides when the list is empty.

## Natural Language

Plain English after `@ctx` routes to a read-only command:

- "Do you remember?" → `/remember`
- "What should I work on next?" → `/next`
- "Time to wrap up" → `/wrap-up`
- "Show me the status" → `/status`
- "Check for drift" → `/drift`

A keyword match never changes context: it cannot add, complete, or
dismiss anything. Unmatched text shows the command list.

## Auto-Bootstrap

If the ctx CLI isn't found on PATH or at the configured path, the
extension automatically downloads the correct platform binary from
[GitHub Releases](https://github.com/ActiveMemory/ctx/releases):

1. Detects OS and architecture (darwin/linux/windows, amd64/arm64)
2. Fetches the latest release from the GitHub API
3. Downloads and verifies the matching binary
4. Caches it in VS Code's global storage directory

Subsequent sessions reuse the cached binary. To force a specific version,
set `ctx.executablePath` in your settings.

## Follow-Up Suggestions

After a command, Copilot Chat offers context-aware follow-ups. For
example:

- After `/init` → "Show context status" or "What should I work on next?"
- After `/drift` → "Sync context with codebase" or "Run health check"
- After `/brainstorm` → "Turn this into a spec"
- After `/reflect` → "Wrap up the session"

## Prerequisites

- VS Code 1.93+
- [GitHub Copilot Chat](https://marketplace.visualstudio.com/items?itemName=GitHub.copilot-chat) extension
- [ctx](https://ctx.ist) CLI on PATH, or let the extension auto-download it

## Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| `ctx.executablePath` | `ctx` | Path to the ctx CLI binary. Set this if ctx isn't on PATH and you don't want auto-download. |

## Development

```bash
cd editors/vscode
npm install
npm run watch   # Watch mode
npm run build   # Production build
npm test        # vitest
npm run lint    # eslint
```

### Architecture

The extension is a single-file implementation (`src/extension.ts`) that:

- Registers a `ChatParticipant` with `@ctx` as the handle
- Dispatches slash commands through two tables: `CLI_COMMANDS` (handlers
  that build a `ctx` argv) and `SKILLS` (command → canonical skill)
- Runs the ctx CLI via `execFile` **without a shell**, so prompt text
  reaches the binary as literal arguments, with stdin closed so no
  command can wait on a prompt
- Bundles the skill files with esbuild's text loader
  (`--loader:.md=text`); `vitest.config.ts` mirrors the loader

### Testing

- `src/extension.test.ts`: handler behavior against a mocked
  `execFile` and a VS Code API mock (`src/vscodeMock.ts`).
- `src/commandParity.test.ts`: `package.json` commands == dispatched
  commands; each skill-backed command bundles the skill it names; every
  follow-up and natural-language route targets a real command. It then
  drives a scenario per command branch through the chat handler and
  records every `ctx` argv in `src/ctx-cli-surface.json` (a file
  snapshot).
- `internal/bootstrap/vscode_surface_test.go` (Go, runs with `go test
  ./...`): parses every argv in that snapshot against the real cobra
  command tree, runs the `add` invocations in a scratch project, and
  checks every listed skill ships. A CLI rename that strands a chat
  command fails there.

After changing what a command runs, refresh the snapshot and review the
diff:

```bash
npx vitest run -u
git diff src/ctx-cli-surface.json
```

## Release

This extension is **published separately from the ctx Go binary**.
It does *not* ride along with `release.yml`. The release pipeline
is intentionally manual: a maintainer runs `vsce publish` from a
clean checkout against the `activememory` publisher account.

CI guardrails that protect this manual publish (`vscode-extension`
job in `.github/workflows/ci.yml`) run on every PR and push to
`main`:

- `npm ci`: clean dependency install from the committed lockfile.
- `npm run build`: esbuild bundles `src/extension.ts` (and the skill
  files it imports) to `dist/extension.js`.
- `npx tsc --noEmit -p tsconfig.ci.json`: type-checks the source and
  the tests.
- `npm run lint`: eslint.
- `npm test`: vitest, including the command-parity snapshot.
- `npx vsce package --no-dependencies`: packaging dry-run.

The Go `test` job runs `vscode_surface_test.go` against the same
snapshot.

Release checklist for a maintainer:

1. Bump `version` in `editors/vscode/package.json`.
2. Update `editors/vscode/CHANGELOG.md`.
3. Push to a branch, open PR. The `vscode-extension` CI job must
   pass on the PR head.
4. After merge, from a clean checkout of `main`:
   ```bash
   cd editors/vscode
   npm ci
   npm run build
   npx vsce package
   npx vsce publish      # requires VS Code Marketplace token
   ```
5. Tag the release commit and push the tag (the ctx-binary release
   workflow keys on `v*` tags; the extension's tag does not need
   to match, but keeping them in lockstep simplifies support).

## License

Apache-2.0
