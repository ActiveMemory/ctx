# cli-command-surface

## [2026-05-30-114429] Name the add JSON-ingest flag --json-file, not --json

**Status**: Accepted

**Context**: The CLI-FIX spec specified the literal flag --json <file>, but
--json is already a bool output-format flag across the CLI (ctx
status/drift/doctor/bootstrap --json all mean 'emit machine-readable output').

**Decision**: Name the add JSON-ingest flag --json-file, not --json

**Rationale**: Overloading --json as a string input-path on the add commands
would break that cross-command convention and confuse muscle memory. --json-file
is unambiguous, parallels the existing --file/-f source flag, and leaves -j
free. Pushed back on the spec's literal wording rather than satisfice.

**Consequence**: The add commands intentionally diverge from the spec's literal
--json; the spec was updated to reflect --json-file. Any future JSON-input flag
elsewhere should follow the --json-file naming, reserving --json for bool
output.

---

## [2026-04-16-011520] Deprecate and remove ctx backup

**Status**: Accepted

**Context**: ctx backup is environment-specific (SMB/GVFS), fires nag hooks for
unconfigured users, and solves a problem that belongs to the OS layer. ctx hub
already handles cross-machine knowledge persistence.

**Decision**: Deprecate and remove ctx backup

**Rationale**: Hub handles persistence, backup is env-specific, wrong layer for
ctx to own. No external users depend on it. An internal mirror issue and the
GVFS Linux-only dependency add maintenance burden.

**Consequence**: Need backup-strategy runbook before removal. Maintainer must
set up replacement cron job. About 60 files to remove across CLI, config, hooks,
docs, skills. Spec: specs/deprecate-ctx-backup.md

---

## [2026-04-14-010205] Bootstrap stays under ctx system bootstrap (reverted experimental top-level promotion)

**Status**: Accepted

**Context**: Mid-session promoted ctx bootstrap to top-level to make a stale
CLAUDE.md instruction work. User reverted it and reaffirmed the original design.

**Decision**: Bootstrap stays under ctx system bootstrap (reverted experimental
top-level promotion)

**Rationale**: The ctx system namespace is for agent and hook plumbing the user
does not type by hand. Bootstrap is invoked by AI agents at session start;
surfacing it at top-level pollutes ctx --help for humans without benefit.

**Consequence**: internal/bootstrap/group.go reverted;
internal/config/embed/cmd/system.go header now correctly states bootstrap is
intentionally not promoted. The CLAUDE.md template across the repo (and the
workspace copy) updated to reference ctx system bootstrap as canonical.

---

## [2026-04-01-074416] Rename ctx hook → ctx setup to disambiguate from the hook system

**Status**: Accepted

**Context**: PR #45 contributor assumed hook meant the setup command, causing
naming collisions with the PreToolUse/PostToolUse hook system

**Decision**: Rename ctx hook → ctx setup to disambiguate from the hook system

**Rationale**: hook has a specific meaning in ctx; setup accurately describes
generating AI tool integration configs

**Consequence**: CLI breaking change. All docs, specs, TypeScript extension, and
YAML assets updated. Released specs left as historical.

---

## [2026-03-30-075927] Flags-not-subcommands for journal source: list and show are view modes on a noun, not independent entities

**Status**: Accepted

**Context**: During the journal-recall merge, recall had separate list and show
subcommands. Merging them into journal created a design choice: source list +
source show (three levels) vs source --show (two levels).

**Decision**: Flags-not-subcommands for journal source: list and show are view
modes on a noun, not independent entities

**Rationale**: Keeps CLI nesting to two levels max. Default behavior (bare
source) lists sessions; --show switches to inspect mode. When two operations
differ only in how they view the same data, make them flags on one command.

**Consequence**: journal source dispatches via --show flag rather than
positional subcommand. Future view-mode toggles should follow this pattern.

---

## [2026-03-30-003756] Journal consumed recall — recall CLI package deleted

**Status**: Accepted

**Context**: ctx recall was never registered in bootstrap; ctx journal had all
the same subcommands

**Decision**: Journal consumed recall — recall CLI package deleted

**Rationale**: One dead command group creates confusion in docs and skills.
Journal is the canonical command group.

**Consequence**: internal/cli/recall/ deleted, 19 doc files updated,
docs/cli/recall.md renamed to journal.md, zensical.toml updated. MCP tool
ctx_recall rename tasked separately (API contract)

---

## [2026-03-18-193623] Singular command names for all CLI entities

**Status**: Accepted

**Context**: ctx add used learning (singular) but ctx learnings was plural.
Inconsistency across 6 commands.

**Decision**: Singular command names for all CLI entities

**Rationale**: Less headache for i18n; one rule (singular = entity); developers
think in OOP. Use field values come from DescKey constants for
single-source-of-truth renaming.

**Consequence**: All commands singular: task, decision, learning, change,
permission, dep. YAML keys, desc constants, directory names, and 50+ files
updated.

---

## [2026-03-16-022635] Rename --consequences flag to --consequence for singular consistency

**Status**: Accepted

**Context**: All other CLI flags (context, rationale, lesson, application) are
singular nouns. consequences was the only plural.

**Decision**: Rename --consequences flag to --consequence for singular
consistency

**Rationale**: Singular form matches the pattern. Consistency wins over natural
language preference.

**Consequence**: 75+ files updated. Breaking change for --consequences users.

---



