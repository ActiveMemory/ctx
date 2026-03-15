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

### Phase -2: Further Cleanup

* Human: internal/recall/parser requires a serious refactoring; for example
  the parser object and its private and public methods need to go to its own
  package and other helper functions need to go to a different adjacent package.
* Human: internal/notify/notify.go requires refactoring (all functions bagged in
  one file; types need to go to types.go per convention etc etc)
* Human: split err package into sub packages.

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

### Phase 10: Architecture Mapping Skill (`/ctx-architecture`)

**Context**: Skill that incrementally builds and maintains ARCHITECTURE.md
and DETAILED_DESIGN.md. Coverage tracked in map-tracking.json.
Spec: `specs/ctx-architecture.md`

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

### Phase SP: Configurable Session Prefixes

Spec: `specs/session-prefixes.md`. Read the spec before starting any SP task.

Replace hardcoded `session_prefix` / `session_prefix_alt` pair with a
user-extensible `session_prefixes` list in `.ctxrc`. Parser vocabulary
is not i18n text — it belongs in runtime config.

**SP.1 — Config and defaults:**

- [x] SP.1.1: Add `DefaultSessionPrefixes` to `internal/config/parser/` — `[]string{"Session:", "Oturum:"}` #priority:high #added:2026-03-14 #done:2026-03-14
- [x] SP.1.2: Add `SessionPrefixes []string` field to `CtxRC` in `internal/rc/types.go` (`yaml:"session_prefixes"`) #priority:high #added:2026-03-14 #done:2026-03-14
- [x] SP.1.3: Add `SessionPrefixes()` accessor to `internal/rc/rc.go` — returns rc list, falls back to `parser.DefaultSessionPrefixes` if empty/nil #priority:high #added:2026-03-14 #done:2026-03-14

**SP.2 — Parser migration:**

- [x] SP.2.1: Refactor `isSessionHeader()` and `parseSessionHeader()` in `markdown.go` to use `rc.SessionPrefixes()` instead of `assets.TextDesc()` #priority:high #added:2026-03-14 #done:2026-03-14
- [x] SP.2.2: Remove `TextDescKeyParserSessionPrefix` and `TextDescKeyParserSessionPrefixAlt` from `internal/assets/embed.go` #priority:high #added:2026-03-14 #done:2026-03-14
- [x] SP.2.3: Remove `parser.session_prefix` and `parser.session_prefix_alt` entries from `text.yaml` #priority:high #added:2026-03-14 #done:2026-03-14

**SP.3 — Tests:**

- [x] SP.3.1: Update `markdown_test.go` — keep Turkish cases, add custom-prefix test case (e.g., Japanese "セッション:") #priority:high #added:2026-03-14 #done:2026-03-14
- [x] SP.3.2: Add `rc_test.go` test for `SessionPrefixes()` — default fallback and override behavior #priority:high #added:2026-03-14 #done:2026-03-14

**SP.4 — Documentation:**

- [x] SP.4.1: Add `session_prefixes` to `.ctxrc` reference in `docs/cli/index.md` #priority:medium #added:2026-03-14 #done:2026-03-14
- [x] SP.4.2: Add multilingual session parsing section to `docs/recipes/multi-tool-setup.md` #priority:medium #added:2026-03-14 #done:2026-03-14
- [x] SP.4.3: Update "Extensible Session Parsing" in ARCHITECTURE.md to mention configurable prefixes #priority:medium #added:2026-03-14 #done:2026-03-14
- [x] SP.4.4: Verify `docs/home/contributing.md` mentions prefix extensibility in "Add a Session Parser" section #priority:low #added:2026-03-14 #done:2026-03-14

**SP.5 — Validation:**

- [x] SP.5.1: Run `make lint && make test` — all tests pass, no lint errors #priority:high #added:2026-03-14 #done:2026-03-14

### Phase EH: Error Handling Audit

Systematic audit of silently discarded errors across the codebase.
Many call sites use `_ =` or `_, _ =` to discard errors without
any feedback. Some are legitimate (best-effort cleanup), most are
lazy escapes that hide failures.

- [ ] Add drift check: MCP tool coverage vs CLI commands — programmatic check that compares registered MCP tool names (config/mcp/tool) against ctx CLI subcommands to detect newly added CLI commands without MCP equivalents. Could be a drift detector check or a compliance test. @CoderMungan #priority:medium #added:2026-03-15-120116

- [ ] MCP v0.3: expose additional CLI commands as MCP tools — candidates: ctx_load (full context packet), ctx_agent (token-budgeted packet), ctx_reindex (rebuild indices), ctx_sync (reconcile docs/code), ctx_doctor (health check). Evaluate which provide value over the protocol vs requiring terminal interaction. @CoderMungan #priority:medium #added:2026-03-15-120025

- [ ] Make MCP defaults configurable via .ctxrc — add mcp_recall_limit, mcp_truncate_len, mcp_truncate_content_len, mcp_min_word_len, mcp_min_word_overlap fields to .ctxrc schema; expose via rc.MCP*() with fallback to config/mcp/cfg defaults; update tools.go to read from rc instead of cfg constants. @CoderMungan #priority:medium #added:2026-03-15-114700

- [ ] MCP tools.go cleanup pass: magic strings, duplicated fragments, nested templates. Lines: 461:481 + 186:196 duplicated code; 335 magic number; 382:385 nested TextDescs → single template; 390+851 magic time literal; 443+499+800 magic words; 557+892+902 magic numbers; 590+638 nested TextDesc templating; 820 prefixed %s; 854 suffix %s #priority:high #added:2026-03-15-110429

- [ ] EH.0: Create central warning sink — `internal/log/warn.go` with
      `var Sink io.Writer = os.Stderr` and `func Warn(format string, args ...any)`.
      All stderr warnings (`fmt.Fprintf(os.Stderr, ...)`) route through this
      function. The `fmt.Fprintf` return error is handled once, centrally.
      The sink is swappable (tests use `io.Discard`, future: syslog, file).
      EH.2–EH.4 should use `log.Warn()` instead of raw `fmt.Fprintf`.
      DoD: `grep -rn 'fmt.Fprintf(os.Stderr' internal/` returns zero hits
      #priority:high #added:2026-03-15

- [ ] EH.1: Catalogue all silent error discards — recursive walk of `internal/`
      for patterns: `_ = `, `_, _ = `, `//nolint:errcheck`, bare `return` after
      error-producing calls. Group by category:
      (a) file close in defer — often legitimate but should log on failure
      (b) file write/read — data loss risk, must surface
      (c) os.Remove/Rename — state corruption risk
      (d) fmt.Fprint to stderr — truly best-effort, acceptable
      Commands: `grep -rn '_ =' internal/`, `grep -rn 'nolint:errcheck' internal/`
      Output: spreadsheet in `.context/` with file, line, expression, category,
      and recommended action (log-stderr, return-error, acceptable-as-is).
      DoD: every `_ =` in the codebase is categorised and has a recommended action
      #priority:high #added:2026-03-14

- [ ] EH.2: Address category (b) — file write/read discards. These risk silent
      data loss. Fix: return the error, or at minimum emit to stderr with
      `fmt.Fprintf(os.Stderr, "ctx: ...: %v\n", err)` following the pattern
      established in `internal/log/event.go`.
      DoD: no write/read error is silently discarded
      #priority:high #added:2026-03-14

- [ ] EH.3: Address category (a) — file close in defer. Most are `defer func()
      { _ = f.Close() }()`. For read-only files, close errors are rare but
      should still surface. For write/append files, close can fail if the
      final flush fails — these are data loss. Fix: `if err := f.Close();
      err != nil { fmt.Fprintf(os.Stderr, "ctx: close %s: %v\n", path, err) }`.
      DoD: all defer-close sites log failures to stderr
      #priority:medium #added:2026-03-14

- [ ] EH.4: Address category (c) — os.Remove/Rename discards. These are state
      operations (rotation, pruning, temp file cleanup). Silent failure leaves
      stale state. Fix: stderr warning at minimum; for rotation/rename, consider
      returning the error.
      DoD: no Remove/Rename error is silently discarded
      #priority:medium #added:2026-03-14

- [ ] EH.5: Validate — `grep -rn '_ =' internal/` returns only category (d)
      entries (fmt.Fprint to stderr) and entries explicitly annotated as
      acceptable. Run `make lint && make test` to confirm no regressions.
      DoD: grep output is clean or fully annotated; CI green
      #priority:high #added:2026-03-14

### Phase ET: Error Package Taxonomy (`internal/err/`)

`errors.go` is 1995 lines with 188 functions in a single file. Split into
domain-grouped files. No API changes — same package, same function signatures,
just file reorganization.

Taxonomy (from prefix analysis):

| File             | Prefixes / Domain                                     | ~Count |
|------------------|-------------------------------------------------------|--------|
| `memory.go`      | Memory*, Discover*                                    | 17     |
| `parser.go`      | Parser*                                               | 7      |
| `crypto.go`      | Crypto*, Encrypt*, Decrypt*, GenerateKey, SaveKey, LoadKey, NoKeyAt | 14     |
| `task.go`        | Task*, NoTaskSpecified, NoTaskMatch, NoCompletedTasks | 8      |
| `journal.go`     | LoadJournalState*, SaveJournalState*, ReadJournalDir, NoJournalDir, NoJournalEntries, ScanJournal, UnknownStage, StageNotSet | 10 |
| `session.go`     | Session*, FindSessions, NoSessionsFound, All*, Ambiguous* | 8 |
| `pad.go`         | Edit*, Blob*, ReadScratchpad, OutFlagRequiresBlob, NoConflict*, Resolve* | 10 |
| `recall.go`      | Reindex*, Stats*, EventLog*                           | 6      |
| `fs.go`          | Read*, Write*, Open*, Stat*, File*, Mkdir*, CreateDir, DirNotFound, NotDirectory, Boundary* | 30 |
| `backup.go`      | Backup*, CreateBackup*, CreateArchive*                | 6      |
| `prompt.go`      | Prompt*, NoPromptTemplate, ListTemplates, ReadTemplate, NoTemplate | 7 |
| `hook.go`        | Embedded*, Override*, UnknownHook, UnknownVariant, MarkerNotFound | 6 |
| `skill.go`       | Skill*                                                | 2      |
| `config.go`      | UnknownProfile, ReadProfile, UnknownFormat, UnknownProjectType, InvalidTool, UnsupportedTool, NotInitialized, ContextNotInitialized, ContextDirNotFound, FlagRequires* | 12 |
| `errors.go`      | Remaining general-purpose: WorkingDirectory, CtxNotInPath, ReadInput, InvalidDate*, Reminder*, Drift*, Git*, Webhook*, etc. | ~25 |

- [x] ET.1: Create the 14 domain files with copyright headers, move functions #done:2026-03-14
      verify no function is duplicated or lost. Leave `errors.go` with only
      the uncategorised remainder (~25 functions).
      Validation: `grep -c '^func ' internal/err/*.go | awk -F: '{s+=$2} END {print s}'` equals 188.
      DoD: `go build ./...` and `make test` pass; function count preserved
      #priority:medium #added:2026-03-14

- [x] ET.2: Verify all callers compile #done:2026-03-14

- [x] ET.3: Split remaining errors.go (~31 functions) into final domain files: #done:2026-03-14
      state.go (ReadingStateDir, LoadState, SaveState),
      reminder.go (ReadReminders, ParseReminders, InvalidID, ReminderNotFound, ReminderIDRequired),
      date.go (InvalidDateValue, InvalidDate),
      init.go (NotInitialized, ContextNotInitialized, HomeDir, ReadProjectReadme, ReadInitTemplate, CreateMakefile, DetectReferenceTime),
      git.go (GitNotFound, NotInGitRepo),
      notify.go (WebhookEmpty, SaveWebhook, LoadWebhook, SendNotification — move MarshalPayload from config.go here),
      validation.go (FlagRequired, ArgRequired, DriftViolations, NoInput, ReadInput, ReadInputStream),
      site.go (NoSiteConfig, ZensicalNotFound).
      After split, errors.go should be empty and deleted.
      DoD: `grep -c '^func ' internal/err/*.go` sums to 188; no errors.go remains; `make lint && make test` green
      #priority:high #added:2026-03-14 — `go build ./...` with no changes
      outside `internal/err/`. Since it's the same package, no import changes
      needed. Run `make test` to confirm.
      DoD: CI green
      #priority:medium #added:2026-03-14

- [ ] Add freshness_files to .ctxrc defaults seeded by ctx init — currently the freshness config is only in the gitignored .ctxrc, so new clones don't get it. Consider a .ctxrc.defaults pattern or seeding via ctx init template. #priority:medium #added:2026-03-14-105143

- [ ] SEC.1: Security-sensitive file change hook — PostToolUse on Edit/Write matching security-critical paths (.claude/settings.local.json, .claude/settings.json, CLAUDE.md, .claude/CLAUDE.md, .context/CONSTITUTION.md). Three actions: (1) nudge user in-session, (2) relay to webhook for out-of-band alerting (autonomous loops), (3) append to dedicated security log (.context/state/security-events.jsonl) for forensics. Separate from general event log. Spec needed. #priority:high #added:2026-03-13

- [ ] O.5: Session timeline view — add --sessions flag to ctx system events. Per-session breakdown of eval/fired counts with hook list. See ideas/spec-hook-observability.md Phase 5 #added:2026-03-12-145401

- [ ] O.4: Doctor hook health check — surface hook activity in ctx doctor output (active/evaluated-never-fired/never-evaluated). See ideas/spec-hook-observability.md Phase 4 #added:2026-03-12-145401

- [ ] O.3: Skip reason logging — add eventlog.Skip() with standard reason constants (paused, throttled, condition-not-met). Instrument 19 hook early-exit paths. See ideas/spec-hook-observability.md Phase 3 #added:2026-03-12-145401

- [ ] O.2: Event summary view — add --summary flag to ctx system events. Aggregates eval/fired counts per hook, shows last-eval/last-fired timestamps, lists never-evaluated hooks. See ideas/spec-hook-observability.md Phase 2 #added:2026-03-12-145401

- [ ] O.1: Hook eval logging — wrap hook cobra commands to log 'eval' events on every invocation. Refactor Run() signatures from os.Stdin to io.Reader (peek+replay pattern). Adds eventlog.Eval(), EventTypeEval constant. See ideas/spec-hook-observability.md Phase 1 #added:2026-03-12-145401

- [ ] Companion intelligence recommendation: implement spec from ideas/spec-companion-intelligence.md — ctx doctor companion detection, ctx init recommendation tip, ctx agent awareness in packets #added:2026-03-12-133008

- [ ] Add configurable assets layer: allow users to plug their own YAML files for localization (language selection, custom text overrides). Currently all user-facing text is hardcoded in commands.yaml; need a mechanism to load user-provided YAML that overlays or replaces built-in text. This enables i18n without forking. #priority:low #added:2026-03-07-233756

- [ ] Cleanup internal/cli/system/core/persistence.go: move 10 (base for ParseInt) to config constant #priority:low #added:2026-03-07-220825

- [ ] Cleanup internal/cli/system/core/session_tokens.go: move SessionStats from state.go to types.go #priority:low #added:2026-03-07-220825

- [ ] Cleanup internal/cli/system/core/wrapup.go: line 18 constant should go to config; make WrappedUpExpiry configurable via ctxrc #priority:low #added:2026-03-07-220825

- [ ] Cleanup internal/cli/system/core/version.go: line 81 newline should come from config #priority:low #added:2026-03-07-220819

- [ ] Add taxonomy to internal/cli/system/core/ — currently an unstructured bag of files; group by domain (backup, hooks, session, knowledge, etc.) #priority:medium #added:2026-03-07-220819

- [ ] Cleanup internal/cli/system/core/version_drift.go: line 53 string formatting should use assets #priority:medium #added:2026-03-07-220819

- [ ] Cleanup internal/cli/system/core/state.go: magic permissions (0o750), magic strings ('Context: ' prefix, etc.) #priority:medium #added:2026-03-07-220819

- [ ] Cleanup internal/cli/system/core/smb.go: errors should come from internal/err; lines 101, 116, 111 need assets text #priority:medium #added:2026-03-07-220819

- [ ] Make AutoPruneStaleDays configurable via ctxrc. Currently hardcoded to 7 days in config.AutoPruneStaleDays; add a ctxrc key (e.g., auto_prune_days) and fallback to the default. #priority:low #added:2026-03-07-220512

- [ ] Refactor check_backup_age/run.go: move consts (lines 23-24) to config, magic directories (line 59) to config, symbolic constants for strings (line 72), messages to assets (lines 79, 90-91), extract non-Run functions to system/core, fix docstrings #priority:medium #added:2026-03-07-180020

- [ ] Add ctxrc support for recall.list.limit to make the default --limit for recall list configurable. Currently hardcoded as config.DefaultRecallListLimit (20). #priority:low #added:2026-03-07-164342

- [ ] Extract journal/core into a standalone journal parser package — functionally isolated enough for its own package rather than remaining as core/ #added:2026-03-07-093815

- [ ] Move PluginInstalled/PluginEnabledGlobally/PluginEnabledLocally from initialize to internal/claude — these are Claude Code plugin detection functions, not init-specific #added:2026-03-07-091656

- [ ] Move guide/cmd/root/run.go text to assets, listCommands to separate file + internal/write #added:2026-03-07-090322

- [ ] Move drift/core/sanitize.go strings to assets #added:2026-03-07-090322

- [ ] Move drift/core/out.go output functions to internal/write per convention #added:2026-03-07-090322

- [ ] Move drift/core/fix.go fmt.Sprintf strings to assets — user-facing output text for i18n #added:2026-03-07-090322

- [ ] Move drift/cmd/root/run.go cmd.Print* output strings to internal/write per convention #added:2026-03-07-084152

- [ ] Extract doctor/core/checks.go strings — 105 inline Name/Category/Message values to assets (i18n) and config (Name/Category constants) #added:2026-03-07-083428

- [ ] Split deps/core builders into per-ecosystem packages — go.go, node.go, python.go, rust.go are specific enough for their own packages under deps/core/ or deps/builders/ #added:2026-03-07-082827

- [ ] Audit git graceful degradation — verify all exec.Command(git) call sites degrade gracefully when git is absent, per project guide recommendation #added:2026-03-07-081625

- [ ] Fix 19 doc.go quality issues: system (13 missing subcmds), agent (phantom refs), load/loop (header typo), claude (stale migration note), 13 minimal descriptions (pause, resume, task, notify, decision, learnings, remind, context, eventlog, index, rc, recall/parser, task/core) #added:2026-03-07-075741

- [ ] Move cmd.Print* output strings in compact/cmd/root/run.go to internal/write per convention #added:2026-03-07-074737

- [ ] Extract changes format.go rendering templates to assets — headings, labels, and format strings are user-facing text for i18n #added:2026-03-07-074719

- [ ] Lift HumanAgo and Pluralize to a common package — reusable time formatting, used by changes and potentially status/recall #added:2026-03-07-074649

- [ ] Extract isAlnum predicate for localization — currently ASCII-only in agent keyword extraction (score.go:141) #added:2026-03-07-073900

- [ ] Make stopwords configurable via .ctxrc — currently embedded in assets, domain users need custom terms #added:2026-03-07-073900

- [ ] Make recency scoring thresholds and relevance match cap configurable via .ctxrc — currently hardcoded in config (7/30/90 days, cap 3) #added:2026-03-07-073900

- [ ] Make DefaultAgentCooldown configurable via .ctxrc — currently hardcoded at 10 minutes in config #added:2026-03-07-073106

- [ ] Make TaskBudgetPct and ConventionBudgetPct configurable via .ctxrc — currently hardcoded at 0.40 and 0.20 in config #added:2026-03-07-072714

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
