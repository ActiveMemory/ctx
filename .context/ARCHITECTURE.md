# Architecture

## Overview

ctx is a CLI tool that creates and manages a `.context/` directory
containing structured markdown files. These files provide persistent,
token-budgeted, priority-ordered context for AI coding assistants
across sessions.

Design philosophy:

- **Markdown-centric**: all context is plain markdown; no databases,
  no proprietary formats. Files are human-readable and version-
  controlled alongside the code they describe.
- **Token-budgeteed**: context assembly respects configurable token
  limits so AI agents receive the most important information first
  without exceeding their context window.
- **Priority-ordered**: files are loaded in a deliberate sequence
  (rules before tasks, conventions before architecture) so agents
  internalize constraints before acting.
- **Convention over configuration**: sensible defaults with optional
  `.ctxrc` overrides. No config file required to get started.

For per-module deep dives (types, exported API, data flow, edge cases),
see [DETAILED_DESIGN.md](DETAILED_DESIGN.md).

## Package Dependency Graph

Entry point `cmd/ctx` → `bootstrap` (root Cobra command) → 23 CLI
command packages under `internal/cli/*`. Commands select from shared
packages: `context`, `drift`, `index`, `task`, `validation`,
`recall/parser`, `claude`, `notify`, `journal/state`, `crypto`,
`sysinfo`. Foundation packages (`config`, `assets`, `crypto`,
`sysinfo`) have zero internal dependencies — everything else builds
upward from them. The `rc` package mediates config resolution;
`context` depends on `rc` and `config`; `drift` depends on `context`,
`index`, and `rc`.

*Full dependency tree, matrix, and Mermaid graph:
[architecture-dia-dependencies.md](architecture-dia-dependencies.md)*

## Component Map

<!-- drift-check: ls -d internal/config internal/assets internal/crypto internal/sysinfo 2>/dev/null | wc -l -->
### Foundation Packages (zero internal dependencies)

| Package            | Purpose                                         | Key Exports                                       |
|--------------------|-------------------------------------------------|---------------------------------------------------|
| `internal/config`  | Constants, regex, file names, read order, perms | `FileReadOrder`, `RegExEntryHeader`, `FileType`   |
| `internal/assets`  | Embedded templates via `go:embed`               | `Template()`, `SkillContent()`, `PluginVersion()` |
| `internal/crypto`  | AES-256-GCM encryption (stdlib only)            | `Encrypt()`, `Decrypt()`, `GenerateKey()`         |
| `internal/sysinfo` | OS metrics with platform build tags             | `Collect()`, `Evaluate()`, `MaxSeverity()`        |

<!-- drift-check: ls -d internal/rc internal/context internal/drift internal/index internal/task internal/validation internal/recall/parser internal/claude internal/notify internal/journal/state 2>/dev/null | wc -l -->
### Core Packages

| Package                  | Purpose                                                | Key Exports                                |
|--------------------------|--------------------------------------------------------|--------------------------------------------|
| `internal/rc`            | Runtime config from `.ctxrc` + env + CLI flags         | `RC()`, `ContextDir()`, `TokenBudget()`    |
| `internal/context`       | Load `.context/` directory with token estimation       | `Load()`, `EstimateTokens()`, `Exists()`   |
| `internal/drift`         | Context quality validation (7 checks)                  | `Detect()`, `Report.Status()`              |
| `internal/index`         | Markdown index tables for DECISIONS/LEARNINGS          | `Update()`, `ParseEntryBlocks()`           |
| `internal/task`          | Task checkbox parsing                                  | `Completed()`, `Pending()`, `SubTask()`    |
| `internal/validation`    | Input sanitization and path boundary checks            | `SanitizeFilename()`, `ValidateBoundary()` |
| `internal/recall/parser` | Session transcript parsing (JSONL + Markdown)          | `ParseFile()`, `FindSessionsForCWD()`      |
| `internal/claude`        | Claude Code integration types and skill access         | `Skills()`, `SkillContent()`               |
| `internal/notify`        | Webhook notifications with encrypted URL storage       | `Send()`, `LoadWebhook()`, `SaveWebhook()` |
| `internal/journal/state` | Journal processing pipeline state (JSON)               | `Load()`, `Save()`, `Mark*()`              |

<!-- drift-check: ls -d internal/cli/*/ | wc -l -->
### Entry Point

| Package              | Purpose                                                    |
|----------------------|------------------------------------------------------------|
| `internal/bootstrap` | Create root Cobra command, register 23 subcommands         |

<!-- drift-check: ls -d internal/cli/*/ | wc -l -->
### CLI Commands (`internal/cli/*`)

| Command       | Purpose                                                                         |
|---------------|---------------------------------------------------------------------------------|
| `add`         | Append entries to context files (decisions, tasks, learnings, conventions)      |
| `agent`       | Generate AI-ready context packets with token budgeting                          |
| `compact`     | Archive completed tasks, clean up context files                                 |
| `complete`    | Mark tasks as done in TASKS.md                                                  |
| `decision`    | Manage DECISIONS.md (reindex)                                                   |
| `drift`       | Detect stale/invalid context and report issues                                  |
| `hook`        | Generate AI tool integration configs (Claude, Cursor, Aider, Copilot, Windsurf) |
| `initialize`  | Create `.context/` directory, deploy templates, merge settings                  |
| `journal`     | Export sessions; generate static site or Obsidian vault                         |
| `learnings`   | Manage LEARNINGS.md (reindex)                                                   |
| `load`        | Output assembled context in priority order                                      |
| `loop`        | Generate Ralph loop scripts for autonomous iteration                            |
| `notify`      | Send fire-and-forget webhook notifications                                      |
| `pad`         | Encrypted scratchpad CRUD with blob support and merge                           |
| `permissions` | Permission snapshot/restore (golden images) for Claude Code                     |
| `recall`      | Browse, export, lock/unlock AI session history                                  |
| `remind`      | Session-scoped reminders surfaced at start                                      |
| `serve`       | Serve static journal site locally via zensical                                  |
| `status`      | Display context health summary                                                  |
| `sync`        | Reconcile codebase changes with context documentation                           |
| `system`      | System diagnostics, resource monitoring, hook plumbing                          |
| `task`        | Task archival and snapshots                                                     |
| `watch`       | Monitor stdin for context-update tags and apply them                            |

## Data Flow Diagrams

Five core flows define how data moves through the system:

1. **`ctx init`**: User invokes → `cli/initialize` reads embedded
   templates from `assets` → creates `.context/` directory → writes
   all template files → generates AES-256 key → deploys hooks and
   skills → merges `settings.local.json` → writes/merges `CLAUDE.md`.

2. **`ctx agent`**: Agent invokes with `--budget N` → `cli/agent`
   queries `rc.TokenBudget()` → `context.Load()` reads all `.md`
   files → entries scored by recency and relevance → sorted and
   fitted to token budget → overflow entries listed as "Also Noted"
   → returns Markdown packet.

3. **`ctx drift`**: User invokes → `cli/drift` loads context →
   `drift.Detect()` runs 7 checks (path refs, staleness,
   constitution compliance, required files, file age, entry count,
   missing packages) → returns report with warnings and violations.

4. **`ctx recall export`**: User invokes with `--all` → `cli/recall`
   calls `parser.FindSessionsForCWD()` which scans
   `~/.claude/projects/` → parses JSONL transcripts → loads journal
   state → plans each session (new/regen/skip/locked) → formats as
   Markdown → writes to `.context/journal/` → marks exported in state.

<!-- drift-check: grep -c 'ctx system check-' internal/assets/claude/hooks/hooks.json -->
5. **Hook lifecycle**: Claude Code plugin fires hooks at 4 lifecycle
   points — `UserPromptSubmit` (9 checks: context size, ceremonies,
   persistence, journal, reminders, version, resources, knowledge,
   map staleness), `PreToolUse` (block-non-path-ctx for Bash,
   qa-reminder for Edit, context-load-gate for all tools),
   `PostToolUse` (post-commit for Bash), and `SessionEnd`
   (cleanup-tmp). All hooks execute synchronously; failures softened
   with `|| true` where appropriate.

*Full sequence diagrams:
[architecture-dia-data-flows.md](architecture-dia-data-flows.md)*

## State Diagrams

Five state machines govern lifecycle transitions:

1. **Context files**: Created → Populated (via `ctx init` templates)
   → Active (entries growing via `ctx add` / edits) → Stale (drift
   detected) → back to Active (via fixes) or Archived (via
   `ctx compact` / `ctx consolidate` to `.context/archive/`).

2. **Tasks**: Pending `- [ ]` → In-Progress (`#in-progress` label) /
   Done `- [x]` / Skipped `- [-]` with reason → Archivable (when no
   pending children remain) → Archived (via `ctx task archive` to
   `.context/archive/`).

3. **Journal pipeline**: Exported (JSONL→MD via `recall export`) →
   Enriched (YAML frontmatter, tags) → Normalized (soft-wrap, clean
   JSON) → Fences Verified (fence balance check) → Locked (prevent
   overwrite). Each stage tracked in `.context/journal/.state.json`;
   stages are idempotent; locked entries skip regeneration.

4. **Scratchpad encryption**: User input → `LoadKey()` reads 32-byte
   AES key → decrypt existing `scratchpad.enc` → append new entry →
   re-encrypt all with AES-256-GCM (12-byte random nonce, ciphertext,
   16-byte auth tag) → write `.enc` file.

5. **Config resolution**: CLI flags (highest) > environment variables
   (`CTX_DIR`, `CTX_TOKEN_BUDGET`) > `.ctxrc` (YAML) > hardcoded
   defaults in `internal/rc` → resolved once via `rc.RC()` with
   `sync.Once` singleton caching.

*Full state machine diagrams:
[architecture-dia-state-machines.md](architecture-dia-state-machines.md)*

## Security Architecture

Six defense layers protect the system (innermost to outermost):

- **Layer 0 — Encryption**: AES-256-GCM for scratchpad and webhook
  URLs; 12-byte random nonce + 16-byte authentication tag.
- **Layer 1 — File permissions**: Keys 0600, executables 0755,
  regular files 0644.
- **Layer 2 — Symlink rejection**: `.context/` directory and children
  must not be symlinks (defense against symlink attacks).
- **Layer 3 — Boundary validation**: `ValidateBoundary()` ensures
  resolved `.context/` path stays under project root (prevents path
  traversal).
- **Layer 4 — Permission deny list**: Blocks `sudo`, `rm -rf`,
  `curl`, `wget`, `go install`, force push via Claude Code settings.
- **Layer 5 — Plugin hooks**: `block-non-path-ctx` rejects bare
  `./ctx` or absolute-path invocations; `qa-reminder` gates commits.

**Secret detection** (drift check): scans for `.env`, `credentials*`,
`*secret*`, `*api_key*`, `*password*` — excludes `*.example`,
`*.sample`, and template markers.

*Full defense layer diagram:
[architecture-dia-security.md](architecture-dia-security.md)*

## Key Architectural Patterns

<!-- drift-check: sed -n '/^var FileReadOrder/,/^}/p' internal/config/file.go | grep -cP '^\t' -->
### Priority-Based File Ordering

Files load in a deliberate sequence defined by `config.FileReadOrder`:

1. CONSTITUTION (rules the agent must not violate)
2. TASKS (what to work on now)
3. CONVENTIONS (how to write code)
4. ARCHITECTURE (system structure)
5. DECISIONS (why things are this way)
6. LEARNINGS (gotchas and tips)
7. GLOSSARY (domain terms)
8. AGENT_PLAYBOOK (how to use this system)

Overridable via `priority_order` in `.ctxrc`.

### Token Budgeting

Token estimation uses a 4-characters-per-token heuristic
(see the context package). When the total context exceeds the
budget (default 8000, configurable via `CTX_TOKEN_BUDGET` or
`.ctxrc`), lower-priority files are truncated or omitted.
Higher-priority files always get included first.

### Structured Entry Headers

Decisions and learnings use timestamped headers for chronological
ordering and index generation:

```
## [2026-01-28-143022] Use PostgreSQL for primary database
```

The regex `config.RegExEntryHeader` parses these across the codebase.

### Runtime Config Hierarchy

Configuration resolution (highest priority wins):

1. CLI flags (`--context-dir`)
2. Environment variables (`CTX_DIR`, `CTX_TOKEN_BUDGET`)
3. `.ctxrc` file (YAML)
4. Hardcoded defaults in `internal/rc`

Managed by `internal/rc` with sync.Once singleton caching.

### Extensible Session Parsing

`internal/recall/parser` defines a `SessionParser` interface. Each
AI tool (Claude Code, potentially Aider, Cursor) registers its own
parser. Currently only Claude Code JSONL is implemented.
Session matching uses git remote URLs, relative paths, and exact
CWD matching.

<!-- drift-check: ls internal/assets/claude/skills/ | wc -l -->
### Template and Live Skill Dual-Deployment

Skills exist in two locations:

- **Templates** (`internal/assets/claude/skills/`): embedded in the
  binary, deployed on `ctx init`
- **Live** (`.claude/skills/`): project-local copies that the user
  and agent can edit

`ctx init` deploys templates to live. The `/update-docs` skill
checks for drift between them.

<!-- drift-check: cat internal/assets/claude/hooks/hooks.json | grep -c '"command"' -->
### Hook Architecture

The Claude Code plugin uses four lifecycle hooks defined in
`internal/assets/claude/hooks/hooks.json`: `UserPromptSubmit` (9
checks), `PreToolUse` (3 matchers), `PostToolUse` (1 matcher),
`SessionEnd` (1 cleanup). Hooks execute synchronously; failures
softened with `|| true` where appropriate.

<!-- drift-check: awk '/^require \(/{f=1;next}/^\)/{if(f)exit}f' go.mod | wc -l -->
## External Dependencies

Three direct Go dependencies: `fatih/color` (terminal coloring),
`spf13/cobra` (CLI framework), `gopkg.in/yaml.v3` (YAML parsing).
Optional external tools: `zensical` (static site generation for
journal and docs) and `gpg` (commit signing).

## Build & Release Pipeline

Local: `make build` (CGO_ENABLED=0, ldflags version), `make audit`
(gofmt, go vet, golangci-lint, lint scripts, tests), `make smoke`
(integration tests). Release: `hack/release.sh` bumps VERSION, syncs
plugin version, generates release notes, builds all targets, creates
signed git tag. CI: GitHub Actions runs test + lint on push; release
workflow triggers on `v*` tags producing 6 platform binaries
(darwin/linux/windows x amd64/arm64).

*Full build pipeline diagram:
[architecture-dia-build.md](architecture-dia-build.md)*

<!-- drift-check: ls -d cmd/ctx/ internal/ docs/ hack/ editors/vscode/ specs/ .context/ .claude/ -->
## File Layout

Top-level: `cmd/ctx/` (entry point), `internal/` (all packages),
`docs/` (site source), `site/` (generated static site), `hack/`
(build scripts), `editors/vscode/` (VS Code extension), `specs/`
(feature specs). Under `internal/`: `bootstrap/`, `claude/`,
`cli/` (23 command packages), `config/`, `context/`, `crypto/`,
`drift/`, `index/`, `journal/state/`, `notify/`, `rc/`,
`recall/parser/`, `sysinfo/`, `task/`, `assets/` (embedded
templates, hooks, skills), `validation/`. Project context lives
in `.context/` with its own journal, sessions, and archive
subdirectories. Claude Code integration in `.claude/` with
settings and 30 live skills.

*Full directory tree:
[architecture-dia-build.md](architecture-dia-build.md)*
