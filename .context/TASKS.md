# Tasks

<!--
STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently — never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers
-->

### Phase -1: Quality Verification

- [x] Investigate PreToolUse:Edit hook not firing — qa-reminder (hooks.json:10) should fire on every Edit call but zero PreToolUse:Edit entries appeared in session JSONL despite 6+ Edit calls. Installed ctx v0.6.4 matches assets. Check: hooks sync to active config, matcher dispatch for Edit, session 8b80f8e6. #priority:high #added:2026-02-20-150011

- [ ] Human: there are still minor render issues in journal site.
      Collaborate with AI to sort them out.

- [x] Update enrich-all skill: remove normalize-first prerequisite #priority:medium #added:2026-02-20-062259

- [ ] Move 'ctx journal mark' to 'ctx system mark-journal' #priority:medium #added:2026-02-20-015721

- [ ] Human: Ensure the new journal creation /ctx-journal-normalize and
      /ctx-journal-enrich-all works.
- [ ] Human: Ensure the new ctx files consolidation /ctx-consolidate works.
- [ ] AI: human renamed a skill as "audit" docs may need update
- [ ] AI: ctx-borrow project skill is confusing as `ctx-` prefix implies a
      ctx skill; needs rename.
- [ ] AI: verify and archive completed tasks in TASK.md; the file has gotten
      crowded. Verify each task individually before archiving.

### Phase 0: Ideas

- [ ] Blog: "Building a Claude Code Marketplace Plugin" — narrative from session history, journals, and git diff of feat/plugin-conversion branch. Covers: motivation (shell hooks to Go subcommands), plugin directory layout, marketplace.json, eliminating make plugin, bugs found during dogfooding (hooks creating partial .context/), and the fix. Use /ctx-blog-changelog with branch diff as source material. #added:2026-02-16-111948

- [ ] MCP server integration: expose context as tools/resources via Model
      Context Protocol. Would enable deep integration with any
      MCP-compatible client. #priority:low #source:report-6

**User-Facing Documentation** (from `ideas/done/REPORT-7-documentation.md`):
Docs are feature-organized, not problem-organized. Key structural improvements:

- [ ] Investigate why this PR is closed, is there anything we can leverage
      from it: https://github.com/ActiveMemory/ctx/pull/17

- [ ] Use-case page: "My AI Keeps Making the Same Mistakes" — problem-first
      page showcasing DECISIONS.md and CONSTITUTION.md. Partially covered in
      about.md but deserves standalone treatment as the #2 pain point.
      #priority:medium #source:report-7 #added:2026-02-17

- [ ] Use-case page: "Joining a ctx Project" — team onboarding guide. What
      to read first, how to check context health, starting your first session,
      adding context, session etiquette, common pitfalls. Currently
      undocumented. #priority:medium #source:report-7 #added:2026-02-17

- [ ] Use-case page: "Keeping AI Honest" — unique ctx differentiator.
      Covers confabulation problem, grounded memory via context files,
      anti-hallucination rules in AGENT_PLAYBOOK, verification loop,
      ctx drift for detecting stale context. #priority:medium
      #source:report-7 #added:2026-02-17

- [ ] Expand comparison page with specific tool comparisons: .cursorrules,
      Aider --read, Copilot @workspace, Cline memory, Windsurf rules.
      Current page positions against categories but not the specific tools
      users are evaluating. #priority:low #source:report-7 #added:2026-02-17

- [ ] FAQ page: collect answers to common questions currently scattered
      across docs — Why markdown? Does it work offline? What gets committed?
      How big should my token budget be? Why not a database?
      #priority:low #source:report-7 #added:2026-02-17

- [ ] Enhance security page for team workflows: code review for .context/
      files, gitignore patterns, team conventions for context management,
      multi-developer sharing. #priority:low #source:report-7 #added:2026-02-17

- [ ] Version history changelog summaries: each version entry should have
      2-3 bullet points describing key changes, not just a link to the
      source tree. #priority:low #source:report-7 #added:2026-02-17

**Agent Team Strategies** (from `ideas/REPORT-8-agent-teams.md`):
8 team compositions proposed. Reference material, not tasks. Key takeaways:

- [ ] Document agent team recipes in `hack/` or `.context/`: team
      compositions for feature dev (3 agents), consolidation sprint
      (3-4 agents), release prep (2 agents), doc sprint (3 agents).
      Include coordination patterns and anti-patterns. #priority:low #source:report-8

### Phase 1: Journal Site Improvements `#priority:high`

**Context**: Enriched journal files now have YAML frontmatter (topics, type, outcome,
technologies, key_files). The site generator should leverage this metadata for
better discovery and navigation.

**Features (priority order):**


**Design note**: Topics and key-files both need content search across journal
entries (not just filename matching). Consider FTS5 or similar full-text
indexing — grep across 100+ journal files won't scale.

**Design status**: Understanding confirmed, ready for design approaches.

### Phase 5: Knowledge Scaling `#priority:medium`

**Context**: DECISIONS.md and LEARNINGS.md grow monotonically with no archival path.
Tasks have `ctx tasks archive` / `ctx compact --archive`, but decisions and learnings
accumulate forever. Long-lived projects will hit token budget pressure and
signal-to-noise decay. Spec: `specs/knowledge-scaling.md`

- [ ] P5.1: Extract reusable entry parser for decisions and learnings. The
      archive commands need to parse entries, filter by date, and write back.
      Reuse existing reindex parsing logic where possible.
      #priority:medium #added:2026-02-18

- [ ] P5.2: `ctx decisions archive` command — move entries older than N days
      to `.context/archive/decisions-YYYY-MM-DD.md`. Flags: `--days` (default
      90), `--dry-run`, `--all`, `--keep`. Rebuild index after archival.
      #priority:medium #added:2026-02-18

- [ ] P5.3: `ctx learnings archive` command — same pattern as P5.2 but for
      learnings. Write to `.context/archive/learnings-YYYY-MM-DD.md`.
      #priority:medium #added:2026-02-18

- [ ] P5.4: Extend `ctx compact --archive` to include decisions and learnings
      alongside existing task archival. Same age threshold.
      #priority:medium #added:2026-02-18

- [ ] P5.5: Superseded entry convention — document `~~Superseded by [...]~~`
      marker for decisions. `compact --archive` archives superseded entries
      regardless of age. Add to CONVENTIONS.md.
      #priority:low #added:2026-02-18

- [ ] P5.6: `.contextrc` configuration — `archive_knowledge_after_days` (default
      90) and `archive_keep_recent` (default 5). Wire into archive commands and
      compact.
      #priority:low #added:2026-02-18

- [ ] P5.7: Documentation — update cli-reference.md, context-files.md,
      recipes/context-health.md. Update `/consolidate` skill to suggest
      archival when files are large.
      #priority:low #added:2026-02-18

### Phase 9: Context Consolidation Skill `#priority:medium`

**Context**: `/ctx-consolidate` skill that groups overlapping entries by keyword
similarity and merges them with user approval. Originals archived, not deleted.
Spec: `specs/context-consolidation.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 3)

- [ ] P9.2: Test manually on this project's LEARNINGS.md (20+ entries).
      #priority:medium #added:2026-02-19

### Phase 2: Recall Export Safety `#priority:medium`

**Context**: `ctx recall export --all` silently regenerates every existing
journal body — destructive by default with no confirmation. Users need:
safe-by-default export (new sessions only), explicit `--regenerate` opt-in,
confirmation prompts, lock protection, and clearer flag names.
Spec: `specs/recall-export-safety.md` (supersedes `specs/export-update-mode.md`)

- [x] T2.1: Safe-by-default export — change `--all` to only export new
      sessions (skip existing files). Deprecate `--skip-existing` (now the
      default). Single-session `export <id>` always writes (explicit intent).
      Files: `internal/cli/recall/cmd.go`, `run.go`
      #priority:high #added:2026-02-21 #done:2026-02-21

- [x] T2.2: `--regenerate` flag — add flag to opt in to re-exporting
      existing sessions. `--regenerate` without `--all` is an error.
      `--keep-frontmatter=false` implies `--regenerate`.
      Files: `internal/cli/recall/cmd.go`, `run.go`
      #priority:high #added:2026-02-21 #done:2026-02-21

- [x] T2.3: Confirmation prompt — before destructive writes, compute plan
      (new/regenerate/locked counts), print summary, prompt `proceed? [y/N]`.
      Add `--yes`/`-y` flag to bypass. New-only exports skip confirmation.
      `--dry-run` prints summary and exits (no prompt, no writes).
      Files: `internal/cli/recall/run.go`, `cmd.go`
      #priority:high #added:2026-02-21 #done:2026-02-21

- [ ] T2.4: Lock/unlock state layer — add `Locked` field to `FileState`,
      `MarkLocked()`, `ClearLocked()`, `IsLocked()`, add `"locked"` to
      `Mark()` and `ValidStages`. Tests: mark/clear/round-trip/no-op.
      Files: `internal/journal/state/state.go`, `state_test.go`
      #priority:medium #added:2026-02-20

- [ ] T2.5: Lock/unlock CLI commands — `ctx recall lock <pattern>` and
      `ctx recall unlock <pattern>` with `--all`. Reuse slug/date/id
      matching from export. Multi-part: locking base locks all `-pN` parts.
      Frontmatter: insert/remove `locked: true  # managed by ctx`.
      `.state.json` is source of truth; frontmatter for human visibility.
      File: new `internal/cli/recall/lock.go`, register in `recall.go`
      #priority:medium #added:2026-02-20

- [ ] T2.6: Export respects locks — skip locked files with log line and
      counter. Neither `--regenerate` nor `--force` overrides locks (require
      explicit unlock). Add `locked` counter to confirmation summary.
      File: `internal/cli/recall/run.go`
      #priority:medium #added:2026-02-20

- [ ] T2.7: Replace `--force` with `--keep-frontmatter` — add
      `--keep-frontmatter` flag (bool, default `true`). Keep `--force` as
      deprecated alias via `cmd.Flags().MarkDeprecated`. Logic:
      `discardFrontmatter := !keepFrontmatter || force`.
      `--keep-frontmatter=false` implies `--regenerate`.
      Files: `internal/cli/recall/cmd.go`, `run.go`
      #priority:medium #added:2026-02-20

- [x] T2.8: Ergonomics — bare `ctx recall export` prints help (not error).
      Files: `internal/cli/recall/cmd.go`, `run.go`
      #priority:medium #added:2026-02-20 #done:2026-02-21

- [x] T2.9: Tests — safe-default (--all skips existing), --regenerate
      re-exports, confirmation prompt fires on regenerate, --yes bypasses,
      --dry-run shows summary without writing. 8 new tests, 6 updated.
      Lock/export and --keep-frontmatter tests deferred to Phase 2.
      Files: `internal/cli/recall/run_test.go`
      #priority:medium #added:2026-02-21 #done:2026-02-21

- [x] T2.10: Documentation — updated cli-reference, session-journal,
      session-archaeology recipe, recall skill, common-workflows, publishing
      recipe. Documented new flags (--regenerate, --yes, --dry-run),
      deprecation of --skip-existing. Lock/unlock and --keep-frontmatter
      docs deferred to Phase 2.
      #priority:low #added:2026-02-21 #done:2026-02-21

### Maintenance

- [ ] Recipes section needs human review. For example, certain workflows can
  be autonomously done by asking AI "can you record our learnings?" but
  from the documenation it's not clear. Spend as much time as necessary
  on every single recipe.

- [ ] Investigate ctx init overwriting user-generated content in .context/ 
      files. Commit a9df9dd wiped 18 decisions from DECISIONS.md, replacing with 
      empty template. Need guard to prevent reinit from destroying user data 
     (decisions, learnings, tasks). Consider: skip existing files, merge strategy, 
      or --force-only overwrite. #added:2026-02-06-182205
- [ ] Add ctx help command; use-case-oriented cheat sheet for lazy CLI users. 
      Should cover: (1) core CLI commands grouped by workflow (getting started, tracking decisions, browsing history, AI context), (2) available slash-command skills with one-line descriptions, (3) common workflow recipes showing how commands and skills combine. One screen, no scrolling. Not a skill; a real CLI command. #added:2026-02-06-184257
- [ ] Add topic-based navigation to blog when post count reaches 15+ #priority:low #added:2026-02-07-015054
- [ ] Revisit Recipes nav structure when count reaches ~25 — consider grouping into sub-sections (Sessions, Knowledge, Security, Advanced) to reduce sidebar crowding. Currently at 18. #priority:low #added:2026-02-20
- [ ] Review hook diagnostic logs after a long session. Check `.context/logs/check-persistence.log` and `.context/logs/check-context-size.log` to verify hooks fire correctly. Tune nudge frequency if needed. #priority:medium #added:2026-02-09
- [ ] Run `/consolidate` to address codebase drift. Considerable drift has
      accumulated (predicate naming, magic strings, hardcoded permissions,
      godoc style). #priority:medium #added:2026-02-06
- [ ] `/ctx-journal-enrich-all` should handle export-if-needed: check for
      unexported sessions before enriching and export them automatically,
      so the user can say "process the journal" and the skill handles the
      full pipeline (export → normalize → enrich). #priority:medium #added:2026-02-09
- [ ] Add `--date` or `--since`/`--until` flags to `ctx recall list` for
      date range filtering. Currently the agent eyeballs dates from the
      full list output, which works but is inefficient for large session
      histories. #priority:low #added:2026-02-09
- [ ] Enhance CONTRIBUTING.md: add architecture overview for contributors
      (package map), how to add a new command (pattern to follow), how to
      add a new parser (interface to implement), how to create a skill
      (directory structure), and test expectations per package. Lowers the
      contribution barrier. #priority:medium #source:report-6 #added:2026-02-17
- [ ] Aider/Cursor parser implementations: the recall architecture was
      designed for extensibility (tool-agnostic Session type with
      tool-specific parsers). Adding basic Aider and Cursor parsers would
      validate the parser interface, broaden the user base, and fulfill
      the "works with any AI tool" promise. Aider format is simpler than
      Claude Code's. #priority:medium #source:report-6 #added:2026-02-17

## Reference

**Task Status Labels**:
- `[ ]` — pending
- `[x]` — completed
- `[-]` — skipped (with reason)
- `#in-progress` — currently being worked on (add inline, don't move task)
