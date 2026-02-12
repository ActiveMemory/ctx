# Tasks

<!--
STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently — never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers
-->

### Phase 0: Ideas

**Extractable Patterns** (from `ideas/REPORT-1-extractable-patterns.md`):
Analysis of 69 sessions found 8 recurring workflow patterns. 7 have automation gaps worth addressing.

- [x] `/ctx-remember` slash command: single command that triggers the full
      "reload and present" session-start ritual. Runs `ctx agent --budget 4000`,
      lists last 3 session files, presents structured readback template.
      42% of sessions start with a memory check; current behavior is
      instruction-dependent and sometimes fails. #priority:high #source:report-1
      #done:2026-02-08

- [x] Enhanced `/ctx-save`: produce richer session summaries instead of
      empty auto-save stubs. Should include: key topics discussed, files
      modified, decisions/learnings added, commits created, suggested tasks
      for next session. 2025-era manual saves show what rich summaries look
      like vs current empty stubs. #priority:high #source:report-1
      #done:2026-02-08

- [x] `/ctx-next` slash command: "what should I work on next?" Reads TASKS.md,
      identifies unblocked high-priority items, cross-references recent session
      activity to avoid re-doing work, suggests 1-3 next actions with rationale.
      49 task-review prompts found across sessions. #priority:medium #source:report-1
      #done:2026-02-08

- [x] `/ctx-commit` slash command: commit-with-context workflow. Runs standard
      commit workflow, then prompts "Any decisions or learnings to record?"
      and optionally runs `/ctx-reflect`. Connects the 120 commit sequences
      with context persistence. #priority:medium #source:report-1
      #done:2026-02-08

- [x] `make check` target: single entry point for build + test + lint cycle.
      311 build-then-test sequences found; agent frequently runs raw `go build`
      and `go test` instead of using Makefile. Alternatively, enhance `/qa`
      to include the build step. #priority:medium #source:report-1
      **Done**: `make check` = `build` + `audit` (fmt, vet, lint, drift,
      docs, test). #done:2026-02-08

- [x] `/ctx-implement` slash command: accepts a plan document (inline or file
      path), breaks it into checkable subtasks, executes each step with
      build/test verification between steps, reports progress. 31 plan-
      implementation prompts found, but agents handle this well already
      with raw prompts. #priority:low #source:report-1
      #done:2026-02-08

- [x] Post-commit doc drift reminder: hook or convention that suggests
      `/update-docs` after commits touching source code. Existing `/update-docs`
      covers the need, just needs prompting integration — could be part of
      `/ctx-commit` workflow. #priority:low #source:report-1
      **Done**: Added as step 4 in `/ctx-commit` — conditionally suggests
      `/update-docs` when source files are committed. #done:2026-02-08

- [x] Update user-facing documentation for all of the above skills and use cases.

- [x] create "use-case-based" user-facing documentation (the problem, how
      ctx solves it, typical workflow, best practices, gotchas, etc.)

- [ ] Recipes section needs human review. For example, certain workflows can
      be autonomously done by asking AI "can you record our learnings?" but
      from the documenation it's not clear. Spend as much time as necessary
      on every single recipe.

**Documentation Drift** (from `ideas/REPORT-2-documentation-drift.md`):
Overall drift severity LOW. 14 existing doc.go files are accurate. Key gaps below.

- [x] Fix `internal/cli/recall/doc.go`: replace stale `serve` subcommand
      with `export`, remove "(Phase 3)" annotation. Godoc shows a command
      that doesn't exist. #priority:high #source:report-2 #done:2026-02-11

- [x] Fix `internal/claude/doc.go`: change stale `tpl/commands/*.md` to
      `claude/skills/*/SKILL.md`, add `prompt-coach.sh` and
      `check-context-size.sh` to embedded assets list. #priority:medium #source:report-2
      #done:2026-02-11

- [ ] Fix `internal/bootstrap/bootstrap.go:44-46` godoc: add 5 missing
      subcommands (decision, learnings, recall, journal, serve) or replace
      with "all ctx subcommands" to avoid future drift. #priority:medium #source:report-2

- [ ] Fix `internal/cli/hook/run.go` Claude Code example: update stale
      `preToolCall` JSON key to `PreToolUse` with nested matcher
      structure. #priority:medium #source:report-2

- [ ] Add doc.go files for packages lacking documentation: `status`,
      `sync`, `watch`, `serve`, `config`, `validation`. #priority:low #source:report-2

- [ ] Fix package comment mismatch: `internal/cli/decision/decision.go:7`
      says "Package decisions" but package is `decision`. #priority:low #source:report-2

**Maintainability** (from `ideas/REPORT-3-maintainability.md`):
Score 7.5/10. Strong structural discipline. Key refactoring opportunities below.

- [x] Add `Context.File(name string)` method to eliminate 10+ linear scan
      boilerplate instances across 5 packages. Low effort, high
      impact. #priority:high #source:report-3 #done:2026-02-11

- [ ] Extract `findSessions(allProjects bool)` helper in `recall/run.go`:
      identical 20-line block duplicated in 3 functions. #priority:high #source:report-3

- [ ] Extract `deployHookScript()` helper in `initialize/hook.go`:
      4 nearly identical copy-paste blocks for hook script creation.
      Reduces 103 lines to ~30. #priority:medium #source:report-3

- [ ] Move utility functions from `recall/run.go` to `recall/fmt.go`:
      formatDuration, formatTokens, stripLineNumbers, extractSystemReminders,
      normalizeCodeFences, formatToolUse. Brings run.go from 649 to
      ~490 lines. #priority:medium #source:report-3

- [ ] Convert `formatToolUse` switch (recall/run.go:596) to dispatch map:
      45-line switch becomes 10-line `map[string]string`
      lookup. #priority:low #source:report-3

- [ ] Split `config/tpl.go` (401 lines) by feature area: tpl_entry.go,
      tpl_journal.go, tpl_recall.go, tpl_session.go,
      tpl_loop.go. #priority:medium #source:report-3

- [ ] Unify task archiving logic: 3 separate implementations in
      task/run.go, compact/task.go, drift/fix.go. Create shared
      archive helper. #priority:medium #source:report-3

**Security** (from `ideas/REPORT-4-security.md`):
Overall risk LOW. No critical/high findings. 3 medium, 5 low.

- [ ] M-1: Add path boundary validation on `--context-dir` / `CTX_DIR`.
      No guard prevents operations outside project root. In AI-agent
      context with auto-approved commands, a prompt-injected agent could
      write to sensitive locations. Add optional boundary check with
      `--allow-outside-cwd` escape hatch. #priority:medium #source:report-4

- [ ] M-2: Add symlink detection before file read/write in `.context/`.
      No `Lstat()` or `EvalSymlinks()` calls anywhere. A malicious
      `.context/` with symlinks could cause reads/writes outside project
      boundary. #priority:medium #source:report-4

- [ ] M-3: Use secure temp file patterns for cooldown tombstones and
      counter files. Predictable `/tmp/ctx-*` paths with `$PPID` are
      vulnerable to symlink race. Use `os.CreateTemp()` or user-specific
      subdirectory. #priority:low #source:report-4

- [ ] L-3: Fix prompt-coach.sh session tracking: uses `$$` (hook PID)
      instead of session ID, so counter state never persists across
      prompts. Rate limiting is broken. #priority:low #source:report-4

- [ ] L-5: Remove or secure debug logging in `block-git-push.sh`:
      appends all intercepted git commands to world-readable
      `/tmp/claude-hook-debug.log`. #priority:low #source:report-4

**Blog Themes** (from `ideas/REPORT-5-blog-themes.md`):
8 post proposals with narrative sequencing. See report for full details.

- [x] Blog: "ctx v0.3.0: The Discipline Release" -- release post anchoring
      timeline, covers commands-to-skills migration, skill sweep,
      consolidation, backup/monitoring infra. #priority:high #source:report-5
      **Done**: Draft in ideas/blog-draft-2026-02-09-the-discipline-release.md
      #done:2026-02-08

- [x] Blog: "The 3:1 Ratio" -- formalize the consolidation ratio with
      evidence from git history, connect to /consolidate skill and
      convention drift patterns. #priority:medium #source:report-5
      **Done**: Draft in ideas/blog-draft-2026-02-10-the-3-1-ratio.md
      #done:2026-02-08

- [x] Blog: "Hooks: The Invisible Infrastructure" -- deep dive on hook
      ecosystem, bugs, $CLAUDE_PROJECT_DIR migration, repetition fatigue,
      cooldown mechanism. #priority:medium #source:report-5
      **Done**: Draft in ideas/blog-draft-2026-02-11-hooks-the-invisible-infrastructure.md
      #done:2026-02-08

- [x] Blog: "Context as Infrastructure" -- synthesis piece, broadest
      audience, elevates attention budget theory into thesis about
      persistent context. #priority:low #source:report-5
      **Done**: Draft in ideas/blog-draft-2026-02-12-context-as-infrastructure.md
      #done:2026-02-08

**Improvements** (from `ideas/REPORT-6-improvements.md`):
21 opportunities across 7 categories. Key items not already tracked elsewhere:

- [ ] Increase recall test coverage from 8.8% to 50%+. Core user-facing
      feature with near-zero safety net. Session format has already
      changed once. #priority:high #source:report-6

- [ ] CI coverage enforcement: extend `test-coverage` Makefile target
      beyond just `internal/context` to enforce project-wide and
      per-package minimums. #priority:medium #source:report-6

- [ ] Shell completion enrichment: add subcommand argument completions
      (e.g., `ctx add task|decision|learning|convention`). Cobra supports
      this natively. #priority:low #source:report-6

- [ ] MCP server integration: expose context as tools/resources via Model
      Context Protocol. Would enable deep integration with any
      MCP-compatible client. #priority:low #source:report-6

**User-Facing Documentation** (from `ideas/REPORT-7-documentation.md`):
Docs are feature-organized, not problem-organized. Key structural improvements:

- [ ] Create use-case page: "I Keep Re-Explaining My Codebase" -- the #1
      pain point driving adoption. Lead with the problem, show ctx
      solution, before/after comparison. #priority:high #source:report-7

- [ ] Create use-case page: "Is ctx Right for My Project?" -- decision
      framework for evaluation, good fit vs not right fit,
      try-in-5-minutes. #priority:high #source:report-7

- [ ] Expand Getting Started with guided first-session walkthrough showing
      full interaction end-to-end (what user types, what ctx outputs,
      what AI says back). #priority:medium #source:report-7

- [ ] Reorder site navigation: promote Prompting Guide and Context Files
      higher, demote Session Journal and Autonomous Loops. Current order
      front-loads advanced features. #priority:medium #source:report-7

- [ ] Cross-reference blog posts from documentation pages via "Further
      Reading" sections (6 specific link suggestions in
      report). #priority:low #source:report-7

- [ ] Create migration/adoption guide for existing projects: how ctx
      interacts with existing CLAUDE.md, .cursorrules, the --merge
      flag in workflow context. #priority:medium #source:report-7

- [ ] Create troubleshooting page: consolidate scattered troubleshooting
      into one page (ctx not in PATH, hooks not firing, context not
      loading). #priority:low #source:report-7

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

- [x] Check for the latest backup of .context (there is a script in ~/WORKSPACE)
  and if it is >2days old, remind the user to run the backup script.
  we can use a hook for this. Alternative could be detecting this,
  exiting early, and using a cron job; but the disadvantage is, if nothing has
  changed, the cron job will still make a redundant backup: If nobody is using
  the project, then there is no need to run the backup anyway.
  **Done**: UserPromptSubmit hook (check-backup-age.sh) checks marker file age
  and SMB mount status. Backup script moved to hack/backup-context.sh with env
  vars (CTX_BACKUP_SMB_URL). All hook paths migrated to $CLAUDE_PROJECT_DIR.
  `make backup` target added. #done:2026-02-05

- [x] Build minimal ctx monitor dashboard (see specs/monitor-architecture.md).
  - [x] MVP: a terminal script that runs in a separate window, finds the active
    session JSONL, estimates token usage, shows context health (green/yellow/red),
    and refreshes periodically. Later: pluggable auditors
    (drift, repetition detection), signal injection via hooks. #priority:medium #added:2026-02-04-223356
  **Rescoped**: Full dashboard deferred. Instead: check-context-size.sh hook
  with adaptive frequency + ctx-context-monitor skill for in-session awareness.
  context-watch.sh kept as-is for manual terminal monitoring. Repetition
  detection shelved (low ROI, user notices loops faster). #done:2026-02-05
  - [x] Deploy context-watch.sh to .context/tools/ via ctx init. Embedded in
    binary, deployed with 0755 permissions by createTools() in initialize.
    #done:2026-02-06

- [x] Search for a consolidation/code-hygiene skill that can enforce
  project-specific conventions. Today's codebase scan found concrete
  drift: 5 Is* predicate violations, magic strings in 7+ files, 80+ hardcoded
  file permissions. A skill that periodically checks for ctx's known drift
  patterns (godoc style, predicate naming, file organization, magic strings
  in config) would be the enforcement arm of CONVENTIONS.md. Look for skills
  focused on codebase consistency, convention enforcement, or technical debt
  scanning — not generic code review checklists. #priority:medium #added:2026-02-04-222224
  **Done**: Found generic skill (70% redundant per E/A/R), adapted to
  `/consolidate` with 9 project-specific checks. Added line width +
  duplication conventions. Blog post drafted. #completed:2026-02-05

**Features (priority order):**

- [ ] Fix block-non-path-ctx.sh hook: too aggressive matching blocks git -C path commands that don't invoke ctx #added:2026-02-07-211544

- [ ] T1.1: Topics system
      - Single `/topics/index.md` page
      - Popular topics (2+ sessions) get dedicated pages (`/topics/{topic}.md`)
      - Long-tail topics (1 session) listed inline with direct session links
      - All on one page for Ctrl+F discoverability
      #added:2026-02-03

- [ ] T1.2: Key files index
      - Reverse lookup: file → sessions that touched it
      - Uses `key_files` from frontmatter (not parsed from conversation)
      - `/files/index.md` or similar
      #added:2026-02-03

- [ ] T1.3: Session type pages
      - Dedicated page per type: `/types/debugging.md`, `/types/feature.md`, etc.
      - Groups sessions by type (feature, bugfix, refactor, exploration, debugging, documentation)
      #added:2026-02-03

- [ ] T1.4: Investigate inconsistent tool output collapsing
      - Some files have collapsed tool outputs (`<details>`), others don't
      - Example with uncollapsed long outputs: `2026-02-03-sunny-stirring-pillow-7900d2dc`
      - Collapsing threshold is >10 lines per normalize skill spec
      - Unclear where the existing collapsing happens — may be in export, normalize, or site pipeline
      - Investigate why some pass through uncollapsed, fix the gap
      #added:2026-02-06

**Deferred:**
- [-] Timeline narrative — dropped, duplicates `/ctx-blog` skill (see DECISIONS.md)
- Outcome filtering (uncertain value, revisit after seeing data)
- Stats dashboard (skipped - gamification, low ROI)
- Technology index (skipped - not useful for this project)

**Design note**: Topics and key-files both need content search across journal
entries (not just filename matching). Consider FTS5 or similar full-text
indexing — grep across 100+ journal files won't scale.

**Design status**: Understanding confirmed, ready for design approaches.

### Phase 2: Export Preservation `#priority:medium`

- [ ] T2.1: `ctx recall export --update` mode
      - Preserve YAML frontmatter and summary when re-exporting
      - Update only the raw conversation content
      - Solves: `--force` loses enrichments, no-force can't update
      #added:2026-02-03

### Maintenance

- [ ] Investigate ctx init overwriting user-generated content in .context/ files. Commit a9df9dd wiped 18 decisions from DECISIONS.md, replacing with empty template. Need guard to prevent reinit from destroying user data (decisions, learnings, tasks). Consider: skip existing files, merge strategy, or --force-only overwrite. #added:2026-02-06-182205
- [ ] Add ctx help command; use-case-oriented cheat sheet for lazy CLI users. Should cover: (1) core CLI commands grouped by workflow (getting started, tracking decisions, browsing history, AI context), (2) available slash-command skills with one-line descriptions, (3) common workflow recipes showing how commands and skills combine. One screen, no scrolling. Not a skill; a real CLI command. #added:2026-02-06-184257
- [ ] Add topic-based navigation to blog when post count reaches 15+ #priority:low #added:2026-02-07-015054
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

**GitHub Issues**:

- [ ] GH-6: Fix `ctx journal site --build` silent failure on macOS system
      Python 3.9. The stub package v0.0.2 installs but has no CLI binary.
      Update error message to suggest `pipx install zensical` and note
      Python >= 3.10 requirement. #priority:high #source:github-issue-6
      #added:2026-02-11

- [ ] GH-8: Replace `pip install zensical` → `pipx install zensical` across
      docs and Go error messages. 5 doc files + 4 Go source locations need
      updating. Keep Makefile's `.venv/bin/pip install` as-is (venv context).
      See issue for full file:line inventory. #priority:high #source:github-issue-8
      #added:2026-02-11

- [x] GH-7: Session/journal architecture simplification. Three overlapping
      storage layers (`.claude/projects/`, `.context/sessions/`, `.context/journal/`)
      with redundancy. `.context/sessions/` is a dead end -- nothing reads from it.
      **Done**: Eliminated `.context/sessions/` entirely. Deleted `internal/cli/session/`
      (15 files), removed auto-save hook, removed `--auto-save` from watch, removed
      pre-compact auto-save, removed `/ctx-save` skill, updated ~45 docs. Two stores
      remain: raw transcripts (`~/.claude/projects/`) and enriched journal
      (`.context/journal/`). See DECISIONS.md. #done:2026-02-11

## Blocked

## Reference

**Task Status Labels**:
- `[ ]` — pending
- `[x]` — completed
- `[-]` — skipped (with reason)
- `#in-progress` — currently being worked on (add inline, don't move task)
