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

- [x] Create a Skills Reference page (docs/skills.md): standalone page listing 
  all ctx-shipped skills with name, one-liner description, and link to the recipe 
  or doc page that covers it in depth. Analogous to cli-reference.md but for 
  slash commands. Scope: ctx-bundled skills only
  (shipped via internal/tpl/claude/skills/); project-specific skills are out of 
  scope. Should include: skill name, description, when to use, key commands it 
  wraps, and cross-links to recipes where the skill appears. 
  #priority:medium #added:2026-02-14-125604

- [x] Blog: "Parallel Agents with Git Worktrees" — how to use worktrees
      with ctx to tackle large task backlogs. Narrative: 30 open tasks,
      grouping by file overlap, spawning parallel agents, merging results.
      Practical guide with the ctx project as the example. Connects to the
      agent-team-strategies report (REPORT-8). #priority:medium
      #added:2026-02-12 #done:2026-02-14

- [x] "Borrow from the future" merge skill: given two folders of the same
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
      CLI command. #priority:low #added:2026-02-12 #done:2026-02-14
      Done: Implemented as `/ctx-borrow` skill in `.claude/skills/ctx-borrow/SKILL.md`.
      Covers all 3 strategy tiers, conflict detection, dry-run, and selective apply.

**Documentation Drift** (from `ideas/REPORT-2-documentation-drift.md`):
Overall drift severity LOW. 14 existing doc.go files are accurate. Key gaps below.

**Maintainability** (from `ideas/REPORT-3-maintainability.md`):
Score 7.5/10. Strong structural discipline. Key refactoring opportunities below.

**Security** (from `ideas/REPORT-4-security.md`):
Overall risk LOW. No critical/high findings. 3 medium, 5 low.

- [x] M-1: Add path boundary validation on `--context-dir` / `CTX_DIR`.
      No guard prevents operations outside project root. In AI-agent
      context with auto-approved commands, a prompt-injected agent could
      write to sensitive locations. Add optional boundary check with
      `--allow-outside-cwd` escape hatch. #priority:medium #source:report-4
      Done: `ValidateBoundary()` + `CheckSymlinks()` in internal/validation/path.go,
      enforced globally in PersistentPreRun. Tests in path_test.go.

- [x] M-2: Add symlink detection before file read/write in `.context/`.
      No `Lstat()` or `EvalSymlinks()` calls anywhere. A malicious
      `.context/` with symlinks could cause reads/writes outside project
      boundary. #priority:medium #source:report-4
      Done: `CheckSymlinks()` in internal/validation/path.go uses
      `Lstat()` on dir + children. Called from context/loader.go:48.
      Tests in path_test.go.

- [x] M-3: Use secure temp file patterns for cooldown tombstones and
      counter files. Predictable `/tmp/ctx-*` paths with `$PPID` are
      vulnerable to symlink race. Use `os.CreateTemp()` or user-specific
      subdirectory. #priority:low #source:report-4
      Done: `secureTempDir()` in internal/cli/agent/cooldown.go uses
      $XDG_RUNTIME_DIR/ctx or os.TempDir()/ctx-<uid> (0700). Files 0600.

- [x] L-3: Fix prompt-coach.sh session tracking: uses `$$` (hook PID)
      instead of session ID, so counter state never persists across
      prompts. Rate limiting is broken. #priority:low #source:report-4
      Skipped: moot — prompt-coach hook removed entirely.

**Blog Themes** (from `ideas/REPORT-5-blog-themes.md`):
8 post proposals with narrative sequencing. See report for full details.

**Improvements** (from `ideas/REPORT-6-improvements.md`):
21 opportunities across 7 categories. Key items not already tracked elsewhere:

- [x] Align ctx recall list docs with CLI reality #priority:high #added:2026-02-15-191942

- [x] ~~Add drift check: verify .claude/hooks/*.sh files have execute permission~~ Moot: hooks are now Go subcommands (v0.6.0 plugin conversion) #priority:medium #added:2026-02-15-194829 #done:2026-02-16

- [ ] ctx init should distinguish partial .context/ (only logs) from fully initialized. Check for essential files (TASKS.md, CONSTITUTION.md) instead of just directory existence. #added:2026-02-16-110811

- [ ] Hooks should no-op when .context/ is not properly initialized (missing essential files). Currently hooks create .context/logs/ as side effect before ctx init, then ctx init thinks the dir is already initialized. #added:2026-02-16-110721

- [ ] Rename internal/tpl to internal/assets -- the package now holds both .context/ templates and the Claude Code plugin (skills, hooks, manifest). "tpl" is misleading. Mechanical refactor: rename dir, update all ~15 import sites, update embed.go package doc. Low priority, no behavior change. #added:2026-02-16-104745

- [ ] Align ctx recall list CLI output with docs: columnar table format with aligned headers (Slug, Project, Date, Duration, Turns, Tokens) #priority:high #added:2026-02-15-192053

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

- [x] Cross-reference blog posts from documentation pages via "Further
      Reading" sections (6 specific link suggestions in
      report). #priority:low #source:report-7 #done:2026-02-14
      Done: Added "Further Reading" to autonomous-loop, prompting-guide,
      integrations, context-files, and comparison pages.

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

**Documentation:**

- [ ] P3.19: Update Getting Started / Quick Start to mention scratchpad
      as part of `ctx init` output. #priority:low #added:2026-02-13

- [ ] P3.20: Update `ctx help` (when it exists) to include `pad` command.
      #priority:low #added:2026-02-13

### Phase 4: Obsidian Vault Export (`ctx journal obsidian`) `#priority:high`

**Context**: Export enriched journal entries as an Obsidian vault with wikilinks,
MOC pages, and graph-optimized cross-linking. Reuses existing journal scan/parse/index
infrastructure with an Obsidian-specific output layer.
Spec: `specs/journal-obsidian.md`

- [x] P4.0: Read `specs/journal-obsidian.md` before starting any P4 task.
      #added:2026-02-14-150803 #done:2026-02-14-152520

- [x] P4.1: Add Obsidian constants to `internal/config/`
      New constants: `ObsidianDirName`, `ObsidianDirEntries`, `ObsidianConfigDir`,
      `ObsidianAppConfig`, MOC filenames, `ObsidianMOCPrefix`. Follow existing
      `Journal*` constant naming pattern.
      Done: `internal/config/obsidian.go` — 3 const blocks covering dirs,
      config, MOC filenames, format templates, and README.
      #added:2026-02-14-150803 #done:2026-02-14-152520

- [x] P4.2: Implement wikilink conversion (`wikilink.go`)
      - `convertMarkdownLinks(content string) string` — regex-based markdown
        link → wikilink conversion. Skip external URLs. Strip `.md` extension
        and path prefixes from targets.
      - `formatWikilink(target, display string) string` — format `[[target|display]]`
      - Unit tests in `wikilink_test.go`
      Done: 11 test cases covering internal, external, multipart, mixed links.
      #added:2026-02-14-150803 #done:2026-02-14-152520

- [x] P4.3: Implement frontmatter transformation (`frontmatter.go`)
      - `transformFrontmatter(content string) string` — rename `topics` → `tags`,
        add `aliases` from title, add `source_file` field. Preserve all other
        fields. Parse and re-emit YAML frontmatter.
      - Unit tests for frontmatter transformation
      Done: `frontmatter.go` + `frontmatter_test.go` with 8 test cases.
      #added:2026-02-14-150803 #done:2026-02-14-152520

- [x] P4.4: Implement MOC page generation (`moc.go`)
      - `generateHomeMOC(entries, topics, keyFiles, sessionTypes)` — root hub
      - `generateObsidianTopicsMOC(topics)` — topics index with wikilinks
      - `generateObsidianTopicPage(topic)` — individual topic page
      - `generateObsidianFilesMOC(keyFiles)` — files index with wikilinks
      - `generateObsidianFilePage(kf)` — individual file page
      - `generateObsidianTypesMOC(types)` — types index with wikilinks
      - `generateObsidianTypePage(st)` — individual type page
      - All use wikilinks, grouped by month, same threshold logic as site
      Done: All 7 generators + `generateObsidianGroupedPage` helper.
      #added:2026-02-14-150803 #done:2026-02-14-152520

- [x] P4.5: Implement related sessions footer (`moc.go`)
      - `generateRelatedFooter(entry, topicIndex)` — append topic/type links
        and "see also" entries sharing topics. Creates bidirectional graph edges.
      Done: `generateRelatedFooter` + `collectRelated` with score-based
      prioritization. Lives in `moc.go` alongside MOC generators.
      #added:2026-02-14-150803 #done:2026-02-14-152520

- [x] P4.6: Implement vault orchestration (`vault.go`)
      - `runJournalObsidian(cmd, output)` → `buildObsidianVault(cmd, journalDir, output)`
        1. Scan entries (reuse `scanJournalEntries`)
        2. Create output dirs (`entries/`, `topics/`, `files/`, `types/`, `.obsidian/`)
        3. Write `.obsidian/app.json`
        4. Transform and write entries (normalize, convert links, transform
           frontmatter, add related footer)
        5. Build indices (reuse `buildTopicIndex` etc.)
        6. Generate and write MOC pages
        7. Generate and write Home.md
      - Normalize content but do NOT write back to source files
      Done: Full pipeline with extracted `buildObsidianVault` for testability.
      Filter helpers: `filterRegularEntries`, `filterEntriesWithTopics`,
      `filterEntriesWithKeyFiles`, `filterEntriesWithType`, `buildTopicLookup`.
      #added:2026-02-14-150803 #done:2026-02-14-152520

- [x] P4.7: Add Cobra command (`obsidian.go`)
      - `journalObsidianCmd()` — register under `ctx journal obsidian`
        with `--output` flag (default `.context/journal-obsidian/`)
      - Wire into `journal.go` command tree
      Done: `obsidian.go` + updated `journal.go` to add command.
      #added:2026-02-14-150803 #done:2026-02-14-152520

- [x] P4.8: Integration test
      - Test fixtures with sample enriched entries
      - Run full pipeline, verify: vault structure, wikilink format in entries,
        MOC pages contain wikilinks, `.obsidian/app.json` exists, frontmatter
        has `tags` not `topics`
      Done: `vault_test.go` with 12 test functions — unit tests for all
      generators, filter helpers, related footer, and a full integration test
      with 3 sample entries verifying vault structure end-to-end.
      #added:2026-02-14-150803 #done:2026-02-14-152520

### Phase 1: Journal Site Improvements `#priority:high`

**Context**: Enriched journal files now have YAML frontmatter (topics, type, outcome,
technologies, key_files). The site generator should leverage this metadata for
better discovery and navigation.

**Features (priority order):**

- [x] T1.1: Topics system
      - Single `/topics/index.md` page
      - Popular topics (2+ sessions) get dedicated pages (`/topics/{topic}.md`)
      - Long-tail topics (1 session) listed inline with direct session links
      - All on one page for Ctrl+F discoverability
      #added:2026-02-03
      Done: `index.go` — `buildTopicIndex`, `generateTopicsIndex`, `generateTopicPage`.
      Wired in `run.go`. Tests in `journal_test.go`.

- [x] T1.2: Key files index
      - Reverse lookup: file → sessions that touched it
      - Uses `key_files` from frontmatter (not parsed from conversation)
      - `/files/index.md` or similar
      #added:2026-02-03
      Done: `index.go` — `buildKeyFileIndex`, `generateKeyFilesIndex`, `generateKeyFilePage`.
      Wired in `run.go`.

- [x] T1.3: Session type pages
      - Dedicated page per type: `/types/debugging.md`, `/types/feature.md`, etc.
      - Groups sessions by type (feature, bugfix, refactor, exploration, debugging, documentation)
      #added:2026-02-03
      Done: `index.go` — `buildTypeIndex`, `generateTypesIndex`, `generateTypePage`.
      Wired in `run.go`.

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

- [ ] Investigate ctx init overwriting user-generated content in .context/ 
      files. Commit a9df9dd wiped 18 decisions from DECISIONS.md, replacing with 
      empty template. Need guard to prevent reinit from destroying user data 
     (decisions, learnings, tasks). Consider: skip existing files, merge strategy, 
      or --force-only overwrite. #added:2026-02-06-182205
- [ ] Add ctx help command; use-case-oriented cheat sheet for lazy CLI users. 
      Should cover: (1) core CLI commands grouped by workflow (getting started, tracking decisions, browsing history, AI context), (2) available slash-command skills with one-line descriptions, (3) common workflow recipes showing how commands and skills combine. One screen, no scrolling. Not a skill; a real CLI command. #added:2026-02-06-184257
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
