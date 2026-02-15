<img src="https://ctx.ist/images/ctx-banner.png" />

# Context CLI v0.3.0

This release is a major evolution in how ctx works with AI agents. Slash commands are replaced by a full Agent Skills system, a new autonomous mode (`--ralph`) enables unattended operation, and a comprehensive suite of hooks keeps context healthy during long sessions. Under the hood, two consolidation sweeps eliminated magic strings, deshadowed variables, and split large files into focused modules.

## Highlights

- **Agent Skills System**: All 21 slash commands migrated to directory-based skills following the agentskills.io spec, each with frontmatter, quality gates, usage examples, and "When NOT to Use" triggers
- **Autonomous Loop Mode**: `ctx init --ralph` creates PROMPT.md and IMPLEMENTATION_PLAN.md configured for independent agent operation without clarifying questions
- **Context Health Hooks**: New `UserPromptSubmit` hooks for context size checkpoints, backup staleness warnings, and a deployed `context-watch.sh` monitor tool
- **Code Consolidation**: Two systematic sweeps replaced magic strings/numbers with constants, deshadowed variables across all packages, and split large files (journal/site.go into 12 files, recall/export.go into 4 files)

## Features

- Convert all `.claude/commands/*.md` to `.claude/skills/*/SKILL.md` directory structure with structured frontmatter
- Add `ctx init --ralph` flag for autonomous agent mode with dedicated PROMPT.md template
- Add context size checkpoint hook (`check-context-size.sh`) with adaptive reminder cadence (silent for 15 prompts, then every 5th, then every 3rd)
- Add `ctx-context-monitor` skill teaching agents how to respond to checkpoint signals
- Deploy `context-watch.sh` to `.context/tools/` via `ctx init` for all users
- Add backup staleness hook (`check-backup-age.sh`) warning when backups are >2 days old
- Add global backup support (`hack/backup-global.sh`) for `~/.claude/` with `make backup-global`
- Add `/consolidate` skill with 9 project-specific drift checks
- Add `/brainstorm` skill for design-before-implementation workflow
- Add `ctx agent --cooldown` and `--session` flags with tombstone debounce
- Add "Update When" triggers to all context file templates (CONSTITUTION, TASKS, CONVENTIONS, etc.)
- Add Anti-Patterns section to AGENT_PLAYBOOK.md (Stale Context, Context Sprawl, Implicit Context, etc.)
- Deploy `Makefile.ctx` template via `ctx init` (amend, never overwrite)
- Add journal site `/files/` index with popular/long-tail split
- Add journal site `/types/` pages grouping sessions by type
- Change `ctx recall export` default to update mode preserving YAML frontmatter (`--skip-existing` for old behavior, `--force` for full overwrite)
- Add 7 deterministic normalize scripts for journal fence/metadata repair
- Add `ctx-journal-normalize` skill for clean journal site rendering

## Bug Fixes

- Fix all 137 journal files: complete fence reconstruction (8 broken files with stray markers)
- Fix `consolidateToolRuns` root cause: `(xN)` on its own line creating broken fences
- Fix session export reliability and browser performance
- Resolve all golangci-lint v2 errcheck and staticcheck warnings
- Remove unnecessary `nl` parameter in `recall/fmt.go`
- Fix task insertion placement in TASKS.md
- Restore `release-notes` and `release` skills dropped during commands-to-skills migration

## Refactoring

- Replace magic strings/numbers with config constants (`ExtJSONL`, `IssueType*`, `DefaultSessionFilename`, `ClaudeField*`, session headings, template strings)
- Deshadow `err`/`ok` variables with descriptive names across drift, recall, session, task, context, and validation packages
- Split `journal/site.go` into 12 focused files; split `recall/export.go` into 4 files
- Extract error constructors to `err.go` files
- Move Claude raw types to `types.go` with project-standard godoc
- Extract hook matcher into `internal/claude/matcher.go`
- Simplify `defer file.Close()` patterns
- Replace `cmd.Printf("\n")` with `cmd.Println(fmt.Sprintf(...))`

## CI

- Upgrade to golangci-lint v2 for CI compatibility
- Fix goinstall mode removal and v2 security warnings

## Documentation

- Add architecture docs and remove DRIFT.md (superseded by `/consolidate` skill)
- Add quick-reference table to CLI reference
- Add copyright headers to normalize scripts
- Add 3 blog posts: "The Attention Budget", "You Can't Import Expertise", "The Anatomy of a Skill That Works"
- Add blog post topic frontmatter to all 6 existing posts
- Add journal pipeline docs (`session-journal.md`)

---

Full changelog: https://github.com/ActiveMemory/ctx/compare/v0.2.0...v0.3.0
