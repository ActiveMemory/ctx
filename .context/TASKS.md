# Tasks

<!--
STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently — never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers
TASK STATUS LABELS:
- `[ ]` — pending
- `[x]` — completed
- `[-]` — skipped (with reason)
- `#in-progress` — currently being worked on (add inline, don't move task)
-->

### Phase -1: Quality Verification




- [-] internal/claude/hooks/registry.go -> — truncated stub, intent unknown; registry is at internal/assets/hooks/messages/ and appears complete

### Phase GK: Global Encryption Key — Spec: `specs/global-encryption-key.md`

- [-] GK.6: Update ARCHITECTURE.md and DETAILED_DESIGN.md for new key resolution model — no references to old paths found in either file #added:2026-03-02-114146










### Phase -2: Housekeeping (Clean Before Renovating)

No broken windows. These fix structural issues in state management,
directory layout, and agent hygiene before adding new features.

Spec: `specs/user-level-dir-relocation.md`, `specs/state-consolidation.md`,
`specs/task-completion-nudge.md`. Read the specs before starting any P-2 task.

**Init guard and state consolidation:**



**User-level directory relocation:**

- [-] P-2.3: Relocate user-level dir from ~/.local/ctx to ~/.ctx — superseded by Phase GK (global encryption key at ~/.ctx/.ctx.key)
  Spec: `specs/user-level-dir-relocation.md`
  #priority:high #added:2026-03-01

- [-] P-2.4: Update docs for ~/.ctx key path — superseded by GK.5 and GK.6
  Spec: `specs/user-level-dir-relocation.md`
  #priority:high #added:2026-03-01

**Task completion nudge:**


### Phase -0.5: Hack Script Absorption

Absorb remaining `hack/` scripts into Go subcommands. Eliminates shell
dependencies, improves portability, and makes the skill layer call `ctx`
directly instead of `make` targets.

**Remaining candidates (from review):**


- [-] P-0.5.2: Evaluate `hack/context-watch.sh` for absorption as `ctx watch` or
  `ctx system watch` — deleted instead; heartbeat now includes token telemetry
  (tokens, context_window, usage_pct) making the watch script redundant.
  #priority:low #added:2026-03-01 #done:2026-03-01

### Phase 0.9: Suppress Nudges After Wrap-Up

Spec: `specs/suppress-nudges-after-wrap-up.md`. Read the spec before starting
any P0.9 task.

**Phase 3 — Skill integration:**


- [-] P0.9.2: Split cli-reference.md — moved to Future
  #added:2026-02-24-204208

- [-] P0.9.3: Investigate proactive content suggestions — moved to Future
  #added:2026-02-24-185754

### Phase 0.8: RSS/Atom Feed Generation (`ctx site feed`)

Spec: `specs/rss-feed.md`. Read the spec before starting any P0.8 task.

**Phase 4 — Tests and integration:**

- [-] P0.8.2: Investigate converting UserPromptSubmit hooks to JSON output —
  Skipped: VERBATIM boxes ARE the feature (human-readable nudges injected into
  agent prompt). JSON would make them less useful. External tooling already gets
  structured JSON via webhooks. #added:2026-02-22-194446



### Phase 0.4: Hook Message Templates

Spec: `specs/future-complete/hook-message-templates.md`. Read the spec before
starting any P0.4 task.

**Phase 2 — Discoverability + documentation:**

Spec: `specs/future-complete/hook-message-customization.md`.

### Phase 0.4.9: Injection Oversize Nudge

Spec: `specs/injection-oversize-nudge.md`. Read the spec before starting
any P0.4.9 task.

### Phase 0.4.10: Context Window Token Usage

Spec: `specs/context-window-usage.md`. Read the spec before starting any
P0.4.10 task.

### Phase 0.6: Plugin Enablement Gap

Ref: `ideas/plugin-enablement-gap.md`. Local-installed plugins get registered
in `installed_plugins.json` but not auto-added to `enabledPlugins`, so slash
commands are invisible in non-ctx projects.

### Prompting Guide — Canonical Reference

- [-] PG.1: Agent/tool compatibility matrix — moved to Future
      #priority:medium #added:2026-02-25

- [-] PG.2: Versioning/stability note — moved to Future
      #priority:low #added:2026-02-25

### Phase 0: Ideas (drift markers)

- [-] P0.1: Standardize drift-check comment format — moved to Future. AI parses
  ad-hoc markers fine; standardization benefits tooling/CLI but not urgent.
  #priority:medium #added:2026-02-28

### Phase 0: Ideas (from competitive analysis)




### Phase 0: Ideas

**User-Facing Documentation** (from `ideas/done/REPORT-7-documentation.md`):
Docs are feature-organized, not problem-organized. Key structural improvements:

- [x] P0.6: Use-case page: "My AI Keeps Making the Same Mistakes" — problem-first
      page showcasing DECISIONS.md and CONSTITUTION.md. Partially covered in
      about.md but deserves standalone treatment as the #2 pain point.
      #priority:medium #source:report-7 #added:2026-02-17 #done:2026-03-05

- [x] P0.7: Use-case page: "Joining a ctx Project" — team onboarding guide. What
      to read first, how to check context health, starting your first session,
      adding context, session etiquette, common pitfalls. Currently
      undocumented. #priority:medium #source:report-7 #added:2026-02-17 #done:2026-03-05

- [x] P0.8: Use-case page: "Keeping AI Honest" — unique ctx differentiator.
      Covers confabulation problem, grounded memory via context files,
      anti-hallucination rules in AGENT_PLAYBOOK, verification loop,
      ctx drift for detecting stale context. #priority:medium
      #source:report-7 #added:2026-02-17 #done:2026-03-05

- [x] P0.9: Expand comparison page with specific tool comparisons: .cursorrules,
      Aider --read, Copilot @workspace, Cline memory, Windsurf rules.
      Current page positions against categories but not the specific tools
      users are evaluating. #priority:low #source:report-7 #added:2026-02-17 #done:2026-03-05

- [x] P0.10: FAQ page: collect answers to common questions currently scattered
      across docs — Why markdown? Does it work offline? What gets committed?
      How big should my token budget be? Why not a database?
      #priority:low #source:report-7 #added:2026-02-17 #done:2026-03-05

- [x] P0.11: Enhance security page for team workflows: code review for .context/
      files, gitignore patterns, team conventions for context management,
      multi-developer sharing. #priority:low #source:report-7 #added:2026-02-17 #done:2026-03-05

- [x] P0.12: Version history changelog summaries: each version entry should have
      2-3 bullet points describing key changes, not just a link to the
      source tree. #priority:low #source:report-7 #added:2026-02-17 #done:2026-03-05

**Agent Team Strategies** (from `ideas/REPORT-8-agent-teams.md`):
8 team compositions proposed. Reference material, not tasks. Key takeaways:

- [x] P0.13: Document agent team recipes in `hack/` or `.context/`: team
      compositions for feature dev (3 agents), consolidation sprint
      (3-4 agents), release prep (2 agents), doc sprint (3 agents).
      Include coordination patterns and anti-patterns. #priority:low #source:report-8 #done:2026-03-05

### Phase S-0: Memory Bridge Groundwork

Prerequisites that unblocked the memory bridge phases.

- [x] Investigate Claude Code project directory naming: examine `~/.claude/projects/` to understand the path encoding scheme — full findings in `ideas/claude-code-project-directory-structure.md` #done:2026-03-05
- [x] Design brainstorm and spec split — foundation in `specs/memory-bridge.md`, future phases in `specs/memory-import.md` and `specs/memory-publish.md` #done:2026-03-05

### Phase MB: Memory Bridge Foundation (`ctx memory`)

Spec: `specs/memory-bridge.md`. Read the spec before starting any MB task.

Bridge Claude Code's auto memory (MEMORY.md) into `.context/` with discovery,
mirroring, and drift detection. Foundation for future import/publish phases.

**MB.1 — Config constants and directory setup:**

- [x] MB.1.1: Add `DirMemory = "memory"` and `DirMemoryArchive = "memory/archive"` to `internal/config/dir.go`
      DoD: constants compile, referenced by at least one other file
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.1.2: Add `FileMemoryMirror = "mirror.md"` and `FileMemoryState = "memory-import.json"` to `internal/config/file.go`
      DoD: constants compile, referenced by at least one other file
      #priority:high #added:2026-03-05 #done:2026-03-05

**MB.2 — Core package `internal/memory/`:**

- [x] MB.2.1: Create `internal/memory/doc.go` with package documentation
      DoD: `go build ./internal/memory/` succeeds
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.2.2: Implement `discover.go` — `DiscoverMemoryPath(projectRoot string) (string, error)`
      Slug encoding: replace `/` with `-`, prefix with `-`. Resolve via `~/.claude/projects/<slug>/memory/MEMORY.md`.
      Handle edge cases: missing file (return error), symlinks, different home dirs.
      DoD: function returns correct path for known project root; returns error when MEMORY.md absent
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.2.3: Write `discover_test.go` — unit tests for slug encoding roundtrip, various home dirs (HOME isolation), missing MEMORY.md
      DoD: `go test ./internal/memory/ -run Discover` passes with 3+ test cases
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.2.4: Implement `mirror.go` — `Sync(contextDir, sourcePath string) error`, `Archive(contextDir string) error`, `Diff(contextDir, sourcePath string) (string, error)`
      Sync: copy MEMORY.md to `.context/memory/mirror.md`, create dirs if needed.
      Archive: snapshot current mirror to `archive/mirror-<timestamp>.md` before overwrite.
      Diff: unified diff between mirror.md and current MEMORY.md.
      DoD: Sync creates mirror, Archive creates timestamped snapshot, Diff returns unified diff string
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.2.5: Write `mirror_test.go` — unit tests: first sync (no prior mirror), sync with archive, diff with changes, empty MEMORY.md
      DoD: `go test ./internal/memory/ -run Mirror` passes with 4+ test cases
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.2.6: Implement `state.go` — sync state tracking (load/save `memory-import.json` with `last_sync`, `last_import`, `last_publish`, `imported_hashes`)
      DoD: state round-trips through JSON; missing file returns zero-value state
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.2.7: Write `state_test.go` — unit tests: load/save roundtrip, missing file defaults, corrupt JSON error
      DoD: `go test ./internal/memory/ -run State` passes with 3+ test cases
      #priority:high #added:2026-03-05 #done:2026-03-05

**MB.3 — CLI commands `internal/cli/memory/`:**

- [x] MB.3.1: Create parent command `ctx memory` in `internal/cli/memory/memory.go`
      DoD: `ctx memory --help` shows subcommands
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.3.2: Register `memory` command in `internal/bootstrap/bootstrap.go`
      DoD: `ctx memory` is accessible from the built binary
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.3.3: Implement `ctx memory sync` in `sync.go`
      Calls Discover → Archive (if mirror exists) → Sync → update state. Reports line counts and drift.
      Exit 0 on success, exit 1 if MEMORY.md not found.
      `--dry-run` flag shows plan without writing.
      DoD: running `ctx memory sync` creates `.context/memory/mirror.md` matching source; `--dry-run` writes nothing
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.3.4: Implement `ctx memory status` in `status.go`
      Shows source path, mirror path, last sync time, line counts, drift indicator, archive count.
      Exit 0 no drift, exit 1 MEMORY.md not found, exit 2 drift detected.
      DoD: output matches spec format; exit codes are correct for each scenario
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [x] MB.3.5: Implement `ctx memory diff` in `diff.go`
      Shows unified diff between mirror and current MEMORY.md.
      DoD: diff output shows added/removed lines; no output when identical
      #priority:medium #added:2026-03-05 #done:2026-03-05

**MB.4 — Hook integration:**

- [x] MB.4.1: Implement `ctx system check-memory-drift` in `internal/cli/system/memory_drift.go`
      Discover MEMORY.md → compare mtime against last sync → output nudge box if drifted.
      Session tombstone at `.context/state/memory-drift-nudged` suppresses repeat nudges.
      Skip silently if MEMORY.md doesn't exist.
      DoD: hook outputs nudge box when drift detected; silent on no drift or missing source; nudge fires once per session
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [x] MB.4.2: Register `check-memory-drift` in `internal/assets/claude/hooks/hooks.json` under `UserPromptSubmit`
      DoD: hook fires on prompt submit; `ctx system check-memory-drift` is callable
      #priority:medium #added:2026-03-05 #done:2026-03-05

**MB.5 — Integration and docs:**

- [x] MB.5.1: Run `make lint && make test` — all existing + new tests pass, no lint errors
      DoD: clean `make lint` and `make test` output (golangci-lint not installed — go vet + gofmt clean)
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MB.5.2: Update ARCHITECTURE.md — add `internal/memory` to Core Packages table, add `memory` to CLI Commands table, update component counts in drift-check comments
      DoD: `ctx drift` does not flag new package as missing; counts match
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [x] MB.5.3: Update DETAILED_DESIGN.md with `internal/memory` module deep dive
      DoD: new section covers types, exports, data flow, edge cases
      #priority:low #added:2026-03-05 #done:2026-03-05

- [x] MB.5.4: Update cli-reference.md with `ctx memory` commands and add memory-bridge recipe
      DoD: cli-reference.md has sync/status/diff entries; recipe covers discovery + sync + drift workflow
      Note: site/ not rebuilt (zensical not installed on this machine)
      #priority:medium #added:2026-03-05 #done:2026-03-05

### Phase MI: Memory Import Pipeline (`ctx memory import`)

Spec: `specs/memory-import.md`. Read the spec before starting any MI task.

Import entries from Claude Code's MEMORY.md into structured `.context/` files
using heuristic classification. Builds on Phase MB foundation (discover, mirror, state).

**MI.1 — Entry parser:**

- [x] MI.1.1: Implement `internal/memory/parse.go` — `ParseEntries(content string) []Entry`
      Parse MEMORY.md into discrete entries. Boundaries: headers (##, ###),
      blank-line-separated paragraphs, list items (-, *). Each Entry has Text, StartLine, Type (header/paragraph/list).
      DoD: parser splits a mixed MEMORY.md into correct entry boundaries
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MI.1.2: Write `internal/memory/parse_test.go` — table-driven tests: headers, paragraphs, list items, mixed content, empty input
      DoD: `go test ./internal/memory/ -run Parse` passes with 5+ test cases (7 tests)
      #priority:high #added:2026-03-05 #done:2026-03-05

**MI.2 — Classifier:**

- [x] MI.2.1: Implement `internal/memory/classify.go` — `Classify(entry Entry) Classification`
      Heuristic keyword matching: conventions (always/prefer/never/standard), decisions (decided/chose/trade-off/approach),
      learnings (gotcha/learned/watch out/bug/caveat), tasks (todo/need to/follow up). Case-insensitive.
      Priority order: conventions > decisions > learnings > tasks > skip.
      Classification has Target (file type) and Confidence (matched keywords).
      DoD: classifier assigns correct targets for representative entries
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MI.2.2: Write `internal/memory/classify_test.go` — table-driven tests: one test per target type, ambiguous entry, skip case
      DoD: `go test ./internal/memory/ -run Classify` passes with 6+ test cases (14 tests)
      #priority:high #added:2026-03-05 #done:2026-03-05

**MI.3 — Deduplication:**

- [x] MI.3.1: Implement hash-based dedup in `internal/memory/state.go` — `EntryHash(text string) string`, `(*State).Imported(hash string) bool`, `(*State).MarkImported(hash, target string)`
      Hash: SHA-256 truncated to 16 hex chars. Check against ImportedHashes before promoting.
      DoD: duplicate entries are skipped; new entries pass through
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MI.3.2: Write dedup tests in `internal/memory/state_test.go` — hash roundtrip, imported check, mark and re-check
      DoD: `go test ./internal/memory/ -run Dedup` passes with 3+ test cases
      #priority:high #added:2026-03-05 #done:2026-03-05

**MI.4 — Promotion and CLI:**

- [x] MI.4.1: Implement `internal/memory/promote.go` — `Promote(entry Entry, classification Classification) error`
      Reuses `add.WriteEntry()` for decisions/learnings/tasks/conventions. Add "Source: auto-memory import" annotation.
      DoD: promoted entry appears in correct .context/ file with proper formatting
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MI.4.2: Wire `ctx memory import` in `internal/cli/memory/import.go`
      Discover → read source → parse entries → classify → dedup → promote. Report counts per target + skipped.
      `--dry-run` flag shows plan without writing.
      DoD: `ctx memory import --dry-run` shows classification plan; without flag, entries appear in .context/ files
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MI.4.3: Write `internal/memory/promote_test.go` — unit tests: promote decision, learning, task, convention to temp .context/
      DoD: `go test ./internal/memory/ -run Promote` passes with 4+ test cases (5 tests)
      #priority:medium #added:2026-03-05 #done:2026-03-05

**MI.5 — Integration and docs:**

- [x] MI.5.1: Integration test with fixture MEMORY.md — end-to-end: parse → classify → dedup → promote → verify files
      DoD: test creates temp .context/, imports fixture, verifies entries landed in correct files
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [x] MI.5.2: Run `go vet ./... && gofmt -l . && make test` — all tests pass, no formatting issues
      DoD: clean output
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MI.5.3: Update docs — add `ctx memory import` to cli/tools.md, update memory-bridge recipe with import workflow
      DoD: docs cover import command, dry-run, and classification heuristics
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [-] MI.future: `--interactive` mode for agent-assisted classification — skipped: `--dry-run` covers review; agents can use `ctx add` directly for overrides; interactive CLI prompts don't compose with agent workflows

### Phase S-3: Blog Post — "Agent Memory is Infrastructure"

Spec: `specs/blog-agent-memory-infrastructure.md`.

- [x] S-3.1: Draft blog post "Agent Memory is Infrastructure" #done:2026-03-04
- [x] S-3.2: Review tone: generous toward Anthropic, concrete, honest about gaps #done:2026-03-04
- [x] S-3.3: Add "The Arc" section connecting to blog series #done:2026-03-04
- [x] S-3.4: Cross-link with companion posts #done:2026-03-04
- [x] S-3.5: Publish after at least one memory feature ships #done:2026-03-05

### Phase MP: Memory Publish (`ctx memory publish`)

Spec: `specs/memory-publish.md`. Read the spec before starting any MP task.

Push curated context from `.context/` into Claude Code's MEMORY.md so the agent
sees structured project context on session start without needing hooks.

**MP.1 — Content selection and formatting:**

- [x] MP.1.1: Implement `internal/memory/publish.go` — `SelectContent(contextDir string, budget int) (string, error)`
      Select pending tasks (max 10), recent decisions (7 days, max 5), key conventions (max 10), recent learnings (7 days, max 5).
      Format as Markdown sections. Trim from bottom (learnings → conventions → decisions) if over budget.
      DoD: returns formatted Markdown within line budget
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MP.1.2: Implement marker-based merge — `MergePublished(existing, published string) string`
      Wrap published block in `<!-- ctx:published -->` / `<!-- ctx:end -->` markers.
      Replace existing marker block if present. Append if markers missing (recovery).
      DoD: merge preserves Claude-owned content outside markers; replaces inside
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MP.1.3: Write `internal/memory/publish_test.go` — marker insertion (empty file), marker replacement, marker stripping recovery, budget trimming, content selection priority
      DoD: `go test ./internal/memory/ -run Publish` passes with 5+ test cases (7 tests)
      #priority:high #added:2026-03-05 #done:2026-03-05

**MP.2 — CLI command:**

- [x] MP.2.1: Wire `ctx memory publish` in `internal/cli/memory/publish.go`
      Discover MEMORY.md → select content → merge → write. Report published line counts.
      `--budget` flag (default 80). `--dry-run` shows plan without writing.
      DoD: `ctx memory publish --dry-run` shows what would be published; without flag, MEMORY.md is updated
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MP.2.2: Wire `ctx memory unpublish` in `internal/cli/memory/unpublish.go`
      Remove the `<!-- ctx:published -->` marker block from MEMORY.md, preserving Claude-owned content.
      DoD: marker block removed, Claude content intact
      #priority:medium #added:2026-03-05 #done:2026-03-05

**MP.3 — Integration and docs:**

- [x] MP.3.1: Integration test — covered by publish_test.go TestSelectContent (end-to-end with fixture .context/)
      DoD: test round-trips publish → read → verify content between markers
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [x] MP.3.2: Run `go vet ./... && gofmt -l . && make test` — all tests pass
      DoD: clean output
      #priority:high #added:2026-03-05 #done:2026-03-05

- [x] MP.3.3: Update docs — add publish/unpublish to cli/tools.md, update memory-bridge recipe, update parent command help
      DoD: docs cover publish workflow, budget flag, marker format
      #priority:medium #added:2026-03-05 #done:2026-03-05

### Phase 9: Context Consolidation Skill `#priority:medium`

**Context**: `/ctx-consolidate` skill that groups overlapping entries by keyword
similarity and merges them with user approval. Originals archived, not deleted.
Spec: `specs/context-consolidation.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 3)

### Phase 10: Architecture Mapping Skill (`/ctx-map`)

**Context**: Skill that incrementally builds and maintains ARCHITECTURE.md
and DETAILED_DESIGN.md. Coverage tracked in map-tracking.json.
Spec: `specs/ctx-map.md`

### Docs: Knowledge Health

- [x] DK.1: Create recipe for knowledge health flow: nudge detection → review →
      `/ctx-consolidate` → archive originals. The old `knowledge-scaling.md`
      recipe was deleted; this replaces it with the nudge-based approach.
      #priority:medium #added:2026-02-21 #done:2026-03-05
- [x] DK.2: Add consolidation cross-link to `knowledge-capture.md` "See also"
      section. #priority:low #added:2026-02-21 — already present #done:2026-03-05

### Phase WC: Write Consolidation

Baseline commit: `4ec5999` (Auto-prune state directory on session start).
Goal: consolidate user-facing messages into `internal/write/` as the central
output package. All CLI commands should route printed output through this package.

- [x] WC.1: Add godoc docstrings to all functions in `internal/write/`, add `doc.go` #added:2026-03-06 #done:2026-03-06
- [x] Move add command example strings from core/example.go to assets — user-facing text for i18n #added:2026-03-06-191651

- [ ] Localization inventory: audit config constants, write package templates, and assets YAML for i18n mapping — low priority, most users are English-first developers #added:2026-03-06-192419

- [ ] Consider indexing tasks and conventions in TASKS.md and CONVENTIONS.md (currently only decisions and learnings have index tables) #added:2026-03-06-190225

- [x] Move internal/cli/add/core/err.go error constructors to internal/err — 12 functions, update all callers #added:2026-03-06-190220

- [ ] Remove FlagNoColor and fatih/color dependency — replace with stdlib terminal coloring or plain output #added:2026-03-06-182831

- [ ] Validate .ctxrc against ctxrc.schema.json at load time — schema is embedded but never enforced, doctor does field-level checks without using it #added:2026-03-06-174851

- [ ] Fix 3 CI compliance issues from PR #27 after merge: missing copyright header on internal/mcp/server_test.go, missing doc.go for internal/cli/mcp/, literal newlines in internal/mcp/resources.go and tools.go #added:2026-03-06-141508

- [ ] Add PostToolUse session event capture. Append lightweight event records (tool name, files touched, timestamp) to .context/state/session-events.jsonl on significant PostToolUse events (file edits, git operations, errors). Not SQLite — just JSONL append. This feeds the PreCompact snapshot hook with richer input so it can report what the agent was actively working on, not just static file state. #added:2026-03-06-185126

- [ ] Add next-step hints to ctx agent and ctx status output. Append actionable suggestions based on context health (e.g. stale tasks, high completion ratio, drift findings). Pattern learned from GitNexus self-guiding agent workflows. #added:2026-03-06-184829

- [ ] Implement PreCompact and SessionStart hooks for session continuity across compaction. Wire ctx agent --budget 4000 to both events: PreCompact outputs context packet before compaction so compactor preserves key info; SessionStart re-injects context packet so fresh/post-compact sessions start oriented. Two thin ctx system subcommands, two entries in hooks.json. See ideas/gitnexus-contextmode-analysis.md for design rationale. #added:2026-03-06-184825

- [ ] Audit fatih/color removal across ~35 files — removed from recall/run.go, recall/lock.go, write/validate.go; ~30 files remain. Separate consolidation pass. #added:2026-03-06-050140

- [ ] Audit remaining 2006-01-02 usages across codebase — 5+ files still use the literal instead of config.DateFormat. Incremental migration. #added:2026-03-06-050140

- [ ] WC.2: Audit CLI packages for direct fmt.Print/Println usage — candidates for migration #added:2026-03-06

## Later

### Phase PR: State Pruning (`ctx system prune`)

Clean stale per-session state files from `.context/state/`. Files with UUID
session ID suffixes accumulate ~6-8 per session with no cleanup. Strategy:
age-based — prune files older than N days (default 7).

- [x] PR.1: Implement `ctx system prune` in `internal/cli/system/prune.go`
      Scan `.context/state/` for files with UUID session ID suffixes.
      Delete files older than `--days` (default 7). `--dry-run` shows plan.
      Preserve global files (no UUID suffix): events.jsonl, memory-import.json, etc.
      DoD: prunes old session files, reports count; `--dry-run` writes nothing
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [x] PR.2: Write `internal/cli/system/prune_test.go` — age-based pruning, dry-run, global file preservation
      DoD: `go test ./internal/cli/system/ -run Prune` passes with 3+ test cases
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [x] PR.3: Run `go vet && make test`, update docs
      DoD: clean build, system.go help updated
      #priority:medium #added:2026-03-05 #done:2026-03-05

- [ ] Regenerate site/ for state-maintenance recipe (docs/recipes/state-maintenance.md added but site not rebuilt) #added:2026-03-05-205425

- [ ] Audit remaining global tombstones for session-scoping: backup-reminded, ceremony-reminded, check-knowledge, journal-reminded, version-checked, ctx-wrapped-up all have the same cross-session suppression bug as memory-drift-nudged #added:2026-03-05-205425

- [ ] F.2: ctx recall import — import Claude Code session JSONLs from local or remote (~/.claude/projects/) into local ~/.claude/projects/. Pure Go: local copy with os.CopyFS-style walk, remote via os/exec ssh+scp (no rsync dependency). --source flag accepts local path or user@host. --dry-run shows what would be copied. Skips existing files (content-addressed by UUID filenames). Enables journal export from sessions that ran on other machines. #added:2026-03-05-141912

- [ ] P0.5: Blog: "Building a Claude Code Marketplace Plugin" — narrative from session
      history, journals, and git diff of feat/plugin-conversion branch.
      Covers: motivation (shell hooks to Go subcommands), plugin directory
      layout, marketplace.json, eliminating make plugin, bugs found during
      dogfooding (hooks creating partial .context/), and the fix. Use
      /ctx-blog-changelog with branch diff as source material. #added:2026-02-16-111948
- [ ] P9.2: Test manually on this project's LEARNINGS.md (20+ entries).
      #priority:medium #added:2026-02-19
- [ ] P0.8.1: Install golangci-lint on the integration server #for-human
      #priority:medium #added:2026-02-23 #added:2026-02-23-170213
- [x] PM.1: Add topic-based navigation to blog when post count reaches 15+ — grouped 24 posts into 6 topic sections #priority:low #added:2026-02-07-015054 #done:2026-03-05
- [x] PM.2: Revisit Recipes nav structure when count reaches ~25 — already grouped in zensical.toml (6 sections); added missing prompt-templates.md; 29 recipes across 6 groups
      into sub-sections (Sessions, Knowledge, Security, Advanced) to reduce
      sidebar crowding. Currently at 18. #priority:low #added:2026-02-20
- [ ] PM.3: Review hook diagnostic logs after a long session. Check
      `.context/logs/check-persistence.log` and
       `.context/logs/check-context-size.log` to verify hooks fire correctly.
       Tune nudge frequency if needed. #priority:medium #added:2026-02-09
- [ ] PM.4: Run `/consolidate` to address codebase drift. Considerable drift has
      accumulated (predicate naming, magic strings, hardcoded permissions,
      godoc style). #priority:medium #added:2026-02-06
- [x] PM.5: Add `--since`/`--until` flags to `ctx recall list` for date range filtering (YYYY-MM-DD, both inclusive)
      #priority:low #added:2026-02-09 #done:2026-03-05
- [x] PM.6: Enhance CONTRIBUTING.md: added "How To Add Things" section to docs/home/contributing.md — new CLI command, new session parser, new bundled skill, test expectations. Updated project layout with memory/. Root CONTRIBUTING.md already links to the full guide.
      #priority:medium #source:report-6 #added:2026-02-17 #done:2026-03-05
- [ ] PM.7: Aider/Cursor parser implementations: the recall architecture was
      designed for extensibility (tool-agnostic Session type with
      tool-specific parsers). Adding basic Aider and Cursor parsers would
      validate the parser interface, broaden the user base, and fulfill
      the "works with any AI tool" promise. Aider format is simpler than
      Claude Code's. #priority:medium #source:report-6 #added:2026-02-17

## Future

- [ ] P0.8.5: Enable webhook notifications in worktrees. Currently `ctx notify`
      silently fails because `.context.key` is gitignored and absent in
      worktrees. For autonomous runs with opaque worktree agents, notifications
      are the one feature that would genuinely be useful. Possible approaches:
      resolve the key via `git rev-parse --git-common-dir` to find the main
      checkout, or copy the key into worktrees at creation time (ctx-worktree
      skill). #priority:medium #added:2026-02-22
- [ ] P0.9.2: Split cli-reference.md (1633 lines) into command group pages:
  cli-overview, cli-init-status, cli-context, cli-recall, cli-tools, cli-system —
  each page covers a natural command group with its subcommands and flags
  #added:2026-02-24-204208
- [ ] P0.9.3: Investigate proactive content suggestions: docs/recipes/publishing.md claims
  agents suggest blog posts and journal rebuilds at natural moments, but no hook
  or playbook mechanism exists to trigger this — either wire it up (e.g.
  post-task-completion nudge) or tone down the docs to match reality
  #added:2026-02-24-185754
- [ ] PG.1: Add agent/tool compatibility matrix to prompting guide — document which
      patterns degrade gracefully when agents lack file access, CLI tools, or
      ctx integration. Treat as a "works best with / degrades to" table.
      #priority:medium #added:2026-02-25
- [ ] PG.2: Add versioning/stability note to prompting guide — "these principles are
      stable; examples evolve" + doc date in frontmatter. Needed once the guide
      becomes canonical and people start quoting it. #priority:low #added:2026-02-25
- [ ] P0.1: Brainstorm: Standardize drift-check comment format and integrate with
  `/ctx-drift` — formalize ad-hoc `<!-- drift-check: ... -->` markers, teach
  drift skill to parse/execute them, publish pattern in docs/recipes. Benefits
  tooling/CLI but AI handles ad-hoc fine for now. #priority:medium #added:2026-02-28
- [ ] F.1: MCP server integration: expose context as tools/resources via Model
  Context Protocol. Would enable deep integration with any
  MCP-compatible client. #priority:low #source:report-6
