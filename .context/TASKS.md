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
Analysis of 69 sessions found 8 recurring workflow patterns. 
7 have automation gaps worth addressing.

- [x] Create a Skills Reference page (docs/skills.md): standalone page listing all ctx-shipped skills with name, one-liner description, and link to the recipe or doc page that covers it in depth. Analogous to cli-reference.md but for slash commands. Scope: ctx-bundled skills only (shipped via internal/tpl/claude/skills/); project-specific skills are out of scope. Should include: skill name, description, when to use, key commands it wraps, and cross-links to recipes where the skill appears. #priority:medium #added:2026-02-14-125604

- [ ] Blog: "Parallel Agents with Git Worktrees" — how to use worktrees
      with ctx to tackle large task backlogs. Narrative: 30 open tasks,
      grouping by file overlap, spawning parallel agents, merging results.
      Practical guide with the ctx project as the example. Connects to the
      agent-team-strategies report (REPORT-8). #priority:medium
      #added:2026-02-12

- [ ] "Borrow from the future" merge skill: given two folders of the same
      project (a "present" and a "future" copy), extract and apply the delta.
      Use case: parallel worktree or separate checkout where commits can't
      be pushed/pulled normally. Strategy ladder:
      (1) Both are git repos → compare commit histories, cherry-pick or
          generate patch series (`git format-patch` / `git am`)
      (2) Only one is a git repo → export the non-git side to a temp git
          repo, diff against the other, produce a patch
      (3) Neither is git → recursive diff, present changes for selective
          application
      Should handle: file additions, deletions, renames, binary files.
      Should warn: conflicting changes in both copies. Could be a skill
      (`/ctx-borrow`) or a standalone `ctx borrow --from ../ctx-future`
      CLI command. #priority:low #added:2026-02-12

**Documentation Drift** (from `ideas/REPORT-2-documentation-drift.md`):
Overall drift severity LOW. 14 existing doc.go files are accurate. Key gaps below.

**Maintainability** (from `ideas/REPORT-3-maintainability.md`):
Score 7.5/10. Strong structural discipline. Key refactoring opportunities below.

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

- [-] L-3: Fix prompt-coach.sh session tracking: uses `$$` (hook PID)
      instead of session ID, so counter state never persists across
      prompts. Rate limiting is broken. #priority:low #source:report-4
      Skipped: moot — prompt-coach hook removed entirely.

**Blog Themes** (from `ideas/REPORT-5-blog-themes.md`):
8 post proposals with narrative sequencing. See report for full details.

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

### Phase 3: Encrypted Scratchpad (`ctx pad`) `#priority:high`

**Context**: Secure one-liner scratchpad, encrypted at rest, synced via git.
AES-256-GCM with symmetric key (gitignored). Plaintext fallback via config.
Spec: `specs/scratchpad.md`

**Core infrastructure:**




**CLI commands:**









**Conflict resolution:**


**Skill and permissions:**



**Tests:**




**Documentation:**




- [ ] P3.19: Update Getting Started / Quick Start to mention scratchpad
      as part of `ctx init` output. #priority:low #added:2026-02-13

- [ ] P3.20: Update `ctx help` (when it exists) to include `pad` command.
      #priority:low #added:2026-02-13

### Phase 1: Journal Site Improvements `#priority:high`

**Context**: Enriched journal files now have YAML frontmatter (topics, type, outcome,
technologies, key_files). The site generator should leverage this metadata for
better discovery and navigation.




**Features (priority order):**


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

- [ ] Recipes section needs human review. For example, certain workflows can
  be autonomously done by asking AI "can you record our learnings?" but
  from the documenation it's not clear. Spend as much time as necessary
  on every single recipe.

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


## Blocked

## Reference

**Task Status Labels**:
- `[ ]` — pending
- `[x]` — completed
- `[-]` — skipped (with reason)
- `#in-progress` — currently being worked on (add inline, don't move task)
