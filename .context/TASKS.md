# Tasks

<!--
STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently — never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers
-->

### Phase -1: Quality Verification

- [ ] Human: there are still minor render issues in journal site.
      Collaborate with AI to sort them out.

- [ ] Update enrich-all skill: remove normalize-first prerequisite #priority:medium #added:2026-02-20-062259

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

### Phase 6: Session Ceremonies `#priority:medium`

**Context**: Sessions have two bookend rituals — init (`/ctx-remember`) and
wrap-up (`/ctx-wrap-up`). Unlike other ctx skills that encourage conversational
use, these should be explicitly invoked as slash commands for precision and
completeness. `/ctx-remember` already exists; `/ctx-wrap-up` is new.
Spec: `specs/session-wrap-up.md`

- [ ] P6.1: Create `/ctx-wrap-up` skill file at
      `.claude/skills/ctx-wrap-up/SKILL.md`. Skill should: gather signal
      (git diff, git log, conversation themes), propose candidates grouped
      by type, persist approved candidates via `ctx add`, optionally offer
      `/ctx-commit`. #priority:medium #added:2026-02-18

- [ ] P6.2: Update `check-persistence` hook to suggest `/ctx-wrap-up`
      instead of generic persistence advice.
      #priority:low #added:2026-02-18

- [ ] P6.3: Create `docs/recipes/session-ceremonies.md` recipe — covers the
      two bookend skills as explicit rituals. Why explicit invocation (not
      conversational) for these two. Session start: `/ctx-remember`. Session
      end: `/ctx-wrap-up`. Quick reference card. When to skip.
      #priority:medium #added:2026-02-18

- [ ] P6.4: Update docs to reflect ceremony pattern — add "Session
      Ceremonies" grouping in docs/skills.md, cross-link from
      recipes/session-lifecycle.md, add callout in docs/prompting-guide.md
      distinguishing ceremony skills (explicit) from workflow skills
      (conversational), mention `/ctx-remember` in docs/first-session.md
      as recommended session start.
      #priority:low #added:2026-02-18

- [ ] P6.5: `ctx system check-ceremonies` hook — scans recent journal
      entries for `/ctx-remember` and `/ctx-wrap-up` usage. If missing from
      last 3 sessions, emits a VERBATIM relay nudge explaining the benefit.
      Journal-first detection (cheap string scan); if journals are stale or
      missing, nudge user to export journals instead of falling back to
      JSONL (avoids expensive large-file scanning that eats context budget).
      Daily throttle, self-silencing when habits form.
      Spec: `specs/ceremony-nudge.md`
      #priority:medium #added:2026-02-18

- [ ] P6.6: Register `check-ceremonies` in hooks.json under
      UserPromptSubmit. Add to `system.go` command tree. Update doc.go
      package comment.
      #priority:medium #added:2026-02-18

- [ ] P6.7: Tests for check-ceremonies — unit tests for journal scanning,
      throttling, output variants (both missing, one missing, neither).
      Integration test with sample journal directory.
      #priority:medium #added:2026-02-18

### Phase 7: Smart Retrieval `#priority:high`

**Context**: `ctx agent --budget` is cosmetic — the budget value is displayed
but never used for content selection. LEARNINGS.md is entirely excluded.
Decisions are title-only. No relevance filtering. This is the highest-impact
improvement for agent context quality. Spec: `specs/smart-retrieval.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 1)

### Phase 8: Drift Nudges `#priority:high`

**Context**: Context files grow without feedback. A project with 47 learnings
gets the same `ctx drift` output as one with 5. Entry count warnings nudge
users to consolidate or archive. Spec: `specs/drift-nudges.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 2)

### Phase 9: Context Consolidation Skill `#priority:medium`

**Context**: `/ctx-consolidate` skill that groups overlapping entries by keyword
similarity and merges them with user approval. Originals archived, not deleted.
Spec: `specs/context-consolidation.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 3)

- [ ] P9.2: Test manually on this project's LEARNINGS.md (20+ entries).
      #priority:medium #added:2026-02-19

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
