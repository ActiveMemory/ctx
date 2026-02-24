# Tasks

<!--
STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently — never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers
-->

### Phase -1: Quality Verification

- [x] AI: ctx-borrow project skill is confusing as `ctx-` prefix implies a
      ctx skill; needs rename. Renamed to /absorb. #done:2026-02-21
- [-] Session pattern analysis skill — rejected. Automated pattern capture from sessions risks training the agent to please rather than push back. Existing mechanisms (learnings, hooks, constitution) already capture process preferences explicitly. See LEARNINGS.md. #added:2026-02-22-212143

- [ ] Implement RSS/Atom feed generation for ctx.ist blog (see specs/rss-feed.md) #added:2026-02-24-025015

- [ ] Install golangci-lint on the integration server #for-human #priority:medium #added:2026-02-23 #added:2026-02-23-170213

- [x] Convert shell hook scripts to `ctx system` subcommands #done:2026-02-24
      Spec: `specs/shell-hooks-to-go.md`. Subtasks:
      - [x] `block-dangerous-commands.go` + tests
      - [x] `check-backup-age.go` + tests
      - [x] Wire into system.go + doc.go
      - [x] Update settings.local.json
      - [x] Delete .claude/hooks/ shell scripts

- [ ] Investigate converting UserPromptSubmit hooks to JSON output — check-persistence, check-ceremonies, check-context-size, check-version, check-resources, and check-knowledge all use plain text with VERBATIM relay. These work differently (prepended to prompt) but may benefit from structured JSON too. #added:2026-02-22-194446

- [ ] Add version-bump relay hook: create a system hook that reminds the agent to bump VERSION, plugin.json, and marketplace.json whenever a feature warrants a version change. The hook should fire during commit or wrap-up to prevent version drift across the three files. #added:2026-02-22-102530

- [x] Rename .scratchpad.key to .context.key #priority:medium #added:2026-02-22-101118

- [ ] Regenerate site HTML after .ctxrc rename #added:2026-02-21-200039

- [x] Fix mark-journal --check to handle locked stage #added:2026-02-21-191851

- [x] `ctx recall sync` — frontmatter-to-state lock sync #done:2026-02-22
      Spec: `specs/recall-sync.md`. Subtasks:
      - [x] Core command (`sync.go`)
      - [x] Wire into recall.go + help text
      - [x] Tests (`sync_test.go`)
      - [x] Docs: cli-reference.md, session-journal.md, session-archaeology.md

- [ ] Enable webhook notifications in worktrees. Currently `ctx notify`
      silently fails because `.context.key` is gitignored and absent in
      worktrees. For autonomous runs with opaque worktree agents, notifications
      are the one feature that would genuinely be useful. Possible approaches:
      resolve the key via `git rev-parse --git-common-dir` to find the main
      checkout, or copy the key into worktrees at creation time (ctx-worktree
      skill). #priority:medium #added:2026-02-22

- [ ] AI: verify and archive completed tasks in TASK.md; the file has gotten
      crowded. Verify each task individually before archiving.

### Phase 0: Ideas

- [ ] Blog: "Building a Claude Code Marketplace Plugin" — narrative from session 
      history, journals, and git diff of feat/plugin-conversion branch. 
      Covers: motivation (shell hooks to Go subcommands), plugin directory 
      layout, marketplace.json, eliminating make plugin, bugs found during 
      dogfooding (hooks creating partial .context/), and the fix. Use 
      /ctx-blog-changelog with branch diff as source material. #added:2026-02-16-111948

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

### Phase 9: Context Consolidation Skill `#priority:medium`

**Context**: `/ctx-consolidate` skill that groups overlapping entries by keyword
similarity and merges them with user approval. Originals archived, not deleted.
Spec: `specs/context-consolidation.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 3)

- [ ] P9.2: Test manually on this project's LEARNINGS.md (20+ entries).
      #priority:medium #added:2026-02-19

### Phase 10: Architecture Mapping Skill (`/ctx-map`)

**Context**: Skill that incrementally builds and maintains ARCHITECTURE.md
and DETAILED_DESIGN.md. Coverage tracked in map-tracking.json.
Spec: `specs/ctx-map.md`

- [x] P10.1: Write spec `specs/ctx-map.md`
      DOD: Covers overview, behavior (first/subsequent/opt-out/nudge),
      tracking.json schema, confidence rubric, staleness detection,
      document constraints, file manifest, non-goals #priority:high
      #done:2026-02-23
- [x] P10.2: Create skill `internal/assets/claude/skills/ctx-map/SKILL.md`
      DOD: Standard template (frontmatter, When to Use, When NOT to Use,
      Execution phases, Quality Checklist). Covers first-run, subsequent-run,
      opt-out, nudge. References confidence rubric. #priority:high
      #done:2026-02-23
- [x] P10.3: Register skill in `internal/config/file.go`
      DOD: `FileDetailedDesign`, `FileMapTracking` constants added.
      `Skill(ctx-map)` in DefaultClaudePermissions. `make build` passes.
      #priority:high #done:2026-02-23
- [x] P10.4: Verify build and tests
      DOD: `make build` and `make test` pass. Skill is embedded (verify
      via `ctx init --force` in temp dir). #priority:high #done:2026-02-23
- [x] P10.5: Run first mapping session on ctx codebase
      DOD: DETAILED_DESIGN.md created with per-module sections.
      map-tracking.json created with coverage data. ARCHITECTURE.md
      reviewed and updated if needed. #priority:medium #done:2026-02-23

### Maintenance

- [ ] Human: Ensure the new journal creation /ctx-journal-normalize and
  /ctx-journal-enrich-all works.
- [ ] Human: Ensure the new ctx files consolidation /ctx-consolidate works.

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

### Docs: Knowledge Health

- [ ] Create recipe for knowledge health flow: nudge detection → review →
      `/ctx-consolidate` → archive originals. The old `knowledge-scaling.md`
      recipe was deleted; this replaces it with the nudge-based approach.
      #priority:medium #added:2026-02-21
- [ ] Fix skills page (`docs/skills.md`): `/ctx-consolidate` entry says
      "runs `ctx reindex`" — should say `ctx learnings reindex` /
      `ctx decisions reindex`. #priority:low #added:2026-02-21
- [ ] Add consolidation cross-link to `knowledge-capture.md` "See also"
      section. #priority:low #added:2026-02-21

- [ ] `ctx reindex` convenience command — runs `ctx decisions reindex` and
      `ctx learnings reindex` in one call. Both files grow at similar rates;
      users always want to reindex both. #priority:low #added:2026-02-21

## Future

- [ ] MCP server integration: expose context as tools/resources via Model
  Context Protocol. Would enable deep integration with any
  MCP-compatible client. #priority:low #source:report-6

## Reference

**Task Status Labels**:
- `[ ]` — pending
- `[x]` — completed
- `[-]` — skipped (with reason)
- `#in-progress` — currently being worked on (add inline, don't move task)
