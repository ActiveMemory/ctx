<img src="https://ctx.ist/images/ctx-banner.png" />

# Context CLI v0.8.0

374 commits since v0.6.0, 2255 files changed, and the most ambitious architectural overhaul in ctx's history. This release adds an MCP server for tool-agnostic AI integration, a memory bridge connecting Claude Code's auto-memory to `.context/`, and a VS Code extension with 15 commands. Under the hood, every CLI package was restructured into a `cmd/ + core/` taxonomy, all user-facing strings were externalized to YAML for i18n readiness, and the sole third-party UI dependency (`fatih/color`) was removed.

## Canonical Release Narrative

(coming soon) https://ctx.ist/blog/

## Highlights

- **MCP Server**: Full Model Context Protocol v0.2 implementation (`ctx mcp serve`) with 8 tools, 4 prompts, resource subscriptions, and session state. Any MCP-compatible AI tool can now read and write `.context/` files without ctx-specific integration.
- **Memory Bridge**: `ctx memory sync/import/diff` connects Claude Code's auto-memory (`MEMORY.md`) to `.context/` files, turning ephemeral tool memory into structured project context.
- **Total String Externalization**: All command descriptions, flag descriptions, and user-facing text routed through embedded YAML assets (`commands.yaml`, `flags.yaml`, `text/*.yaml`) with `DescKey` constant lookups. Foundation for future localization.
- **Architecture Overhaul**: Every CLI package restructured into `cmd/root/` (cobra wiring) + `core/` (logic and types). Output functions moved to `internal/write/` packages. Cross-cutting types consolidated in `internal/entity`. Errors split into 22 domain files in `internal/err/`.
- **VS Code Extension v0.8.0**: 15 new commands covering complete, remind, tasks, pad, notify, and system operations; auto-bootstrap of ctx binary on first use.

## Features

### MCP Server
- Implement MCP v0.1 spec: JSON-RPC 2.0 over stdin/stdout with tools for add, recall, status, drift, compact, and watch
- Implement MCP v0.2: add prompts (agent packet, constitution, tasks review), resource subscriptions, and session state tracking
- Extract routes, catalog, defs, dispatch, and response builder packages for clean separation

### Memory Bridge
- Add `ctx memory sync` to mirror Claude Code MEMORY.md into `.context/memory/`
- Add `ctx memory import` with `--dry-run` for promoting auto-memory entries into decisions, learnings, or conventions
- Add `ctx memory diff` to show divergence between auto-memory and context files
- Add `check-memory-drift` hook to nudge when MEMORY.md changes

### Webhook Notifications
- Add `ctx notify` with fire-and-forget webhook delivery
- Encrypted webhook URL storage using AES-256-GCM
- Thread session ID through all system hooks for webhook attribution
- Pass hook output as webhook detail for all system hooks

### CLI Commands
- Add `ctx guide` command for onboarding and help
- Add `ctx dep` for multi-ecosystem dependency graphs (Go, Node.js, Python, Rust)
- Add `ctx system bootstrap` for AI agent context-dir discovery
- Add `ctx system stats` for session token usage telemetry
- Add `ctx site feed` for Atom 1.0 blog feed generation
- Add `ctx pad import`, `ctx pad export`, and `ctx pad merge`
- Add `ctx recall sync` for frontmatter-to-state lock synchronization
- Add `ctx change` for codebase change detection
- Add `ctx loop` for generating autonomous iteration scripts

### Hooks
- Add `context-load-gate` v2: auto-inject context content instead of directing agent to read files
- Add `check-freshness` hook with per-file review URL configuration
- Add `check-knowledge`, `check-map-staleness`, `specs-nudge`, and `post-commit` hooks
- Add `check-memory-drift`, `check-version`, and `check-task-completion` hooks
- Strengthen `qa-reminder` with hard gate and anti-deferral language
- Make notify events opt-in (no config = no notifications)
- Replace configurable session prefix pair with configurable list

### VS Code Extension
- Add 15 commands: complete, remind, tasks, pad show/edit, notify, system bootstrap/stats/resources/backup, recall export/list, drift, status
- Auto-bootstrap ctx binary on first use
- Fix task handler and broken documentation links

### Security
- Centralize file I/O with system path deny-list
- Add `SafePost` to centralize HTTP client security policy
- Add `SafeCreateFile` and `SafeAppendFile` with permission enforcement
- Move encryption key to global `~/.ctx/.ctx.key` (replaces per-project slug keys)
- Add Markdown session parser for tool-agnostic session ingestion
- Safe-by-default recall export with lock/unlock and `--keep-frontmatter` flag

### System Monitoring
- Add `sysinfo` package with platform build tags for OS metrics
- Add `ctx system resources` for memory, swap, disk, and load display
- Add `ctx doctor` with configurable health checks
- Add heartbeat token telemetry with conditional fields
- Auto-prune state directory on session start

## Bug Fixes

- Fix `resourcesList` returning only 1 MCP resource; deduplicate subscribe/unsubscribe handlers
- Fix recall export `--force` to properly discard enriched frontmatter
- Fix decision insertion inside HTML comment blocks
- Fix memory drift tombstone scoping (global tombstones were suppressing hooks across all sessions)
- Fix `lint-drift.sh` false positives: `Use*` constants incorrectly checked against commands.yaml, wrong exclusion filenames, cross-namespace duplicate check flagging intentional key reuse
- Fix journal `consolidateToolRuns` root cause: `(xN)` on its own line creating broken fences
- Resolve all golangci-lint v2 errcheck and staticcheck warnings
- Eliminate all `nolint:errcheck` directives in favor of explicit error handling

## Refactoring

### Architecture
- Restructure all 24 CLI packages into `cmd/root/` + `core/` taxonomy with `parent.go` wiring
- Move all output functions from `core/` and `cmd/` to `internal/write/` domain packages
- Consolidate cross-cutting types in `internal/entity` (session, parser, export types)
- Split `internal/err` into 22 domain files replacing monolithic `errors.go`
- Extract `internal/entry` for shared entry domain API
- Extract `internal/inspect` for string predicates
- Extract `internal/format` for shared formatting utilities
- Extract `internal/io` for safe file operations (`TouchFile`, `SafeCreateFile`)
- Move `Cmd()` from parent packages into `cmd/root/cmd.go` across all CLI packages
- Pure-logic `CompactContext` with no I/O; callers own file writes and reporting

### String Externalization
- Externalize all 105 command descriptions to `commands.yaml`
- Externalize all flag descriptions to `flags.yaml`
- Split `text.yaml` into 6 domain files loaded via `loadYAMLDir`
- Add exhaustive `TextDescKey` test verifying all 879 constants resolve to non-empty YAML values
- Replace inline flag name strings with `cFlag` constants
- Replace Unicode escape sequences with `config/token` icon constants
- Eliminate hardcoded pluralization; use explicit singular/plural text key pairs

### Code Quality
- Remove `fatih/color` dependency; Unicode symbols are sufficient for terminal output
- Add `doc.go` to all 75 packages missing package documentation
- Standardize import aliases: Yoda-style, camelCase, domain-specific
- Singularize all CLI command directory names
- Replace typographic em-dashes with ASCII equivalents across codebase
- Add `Use*` constants for all 35 system subcommands
- Centralize errors in `internal/err`, eliminating per-package `err.go` broken-window pattern
- Delete legacy key migration code (5 callers, test coverage, zero users)
- Eager `Init()` for static embedded data instead of per-accessor `sync.Once`

### Config
- Split `commands.yaml` into 4 domain files
- Add composite directory path constants for multi-segment paths
- Consolidate all session state to `.context/state/`

## CI

- Add DCO (Developer Certificate of Origin) workflow
- Fix missing `get-pr-commits` step in DCO check

## Dependencies

- Remove `fatih/color` (terminal coloring); replaced by Unicode symbols
- Remaining direct dependencies: `spf13/cobra`, `gopkg.in/yaml.v3`

## Documentation

- Add release operations runbook (`docs/operations/release.md`)
- Add ctx Manifesto as docs landing page
- Restructure Getting Started into focused pages with next-up chains
- Comprehensive docs style audit: fix typography, stale refs, hook diagrams, skill lists
- Add CODEOWNERS file with 4 contributors
- 6 blog posts: "Code Is Cheap, Judgment Is Not", "The 3:1 Ratio", "Merge Debt", and others

## Contributors

- @parlakisik
- @codermungan
- @bilersan
- @hamzaerbay

---

Full changelog: https://github.com/ActiveMemory/ctx/compare/v0.6.0...v0.8.0
