# Tasks

<!--
STRUCTURE RULES (see CONSTITUTION.md):
- Tasks stay in their Phase section permanently — never move them
- Use inline labels: #in-progress, #blocked, #priority:high
- Mark completed: [x], skipped: [-] (with reason)
- Never delete tasks, never remove Phase headers
-->

### Phase -1: Quality Verification

- [-] Session pattern analysis skill — rejected. Automated pattern capture from sessions risks training the agent to please rather than push back. Existing mechanisms (learnings, hooks, constitution) already capture process preferences explicitly. See LEARNINGS.md. #added:2026-02-22-212143

- [ ] Review remaining hack/ scripts for ctx absorption candidates (pad-import-ideas.sh could become ctx pad import --blobs; context-watch.sh monitors raw JSONL token usage which no ctx command covers yet) #added:2026-02-28-184800

- [ ] Document Claude Code JSONL session cleanup behavior in user-facing docs — default 30-day retention, cleanupPeriodDays config, gotchas (0 disables writing, same-day deletion bug), and why journal export matters as archival mechanism. Add to docs/recipes/session-history.md or similar. #priority:medium #added:2026-02-28-132142

- [ ] Add system resource health check to ctx doctor — call sysinfo.Collect() and report memory/swap/disk/load status as a new 'Resources' category. Use the same threshold logic from check-resources (WARNING at 80%/50%/85%/0.8x, DANGER at 90%/75%/95%/1.5x). Graceful degradation: if sysinfo returns Supported:false for a metric, skip it. Add tests with constructed Snapshot values. #added:2026-02-27-230202

- [ ] Auto-detect context window size from session JSONL model field — the JSONL contains the model name (e.g. "claude-opus-4-5-20251101") which can be mapped to the actual window size (200k for standard, 1M for 1M-context models). Currently defaults to 200k via DefaultContextWindow, causing check-context-size to report '110% full' when a 1M-context model is in use with ~220k tokens. **Resolution**: three-tier fallback: `effective_window = detect_from_jsonl(model) ?? ctxrc.context_window ?? 200_000`. JSONL is ground truth (reflects actual model in use); .ctxrc is fallback for first-hook-of-session (no JSONL yet) or unknown models; 200k is safe last resort. **Approach**: (1) parse model field from JSONL in readSessionTokenUsage, (2) maintain a model-to-window lookup (opus/sonnet standard=200k, 1M suffix=1000000), (3) JSONL detection wins when available, .ctxrc fills in when JSONL can't determine window. (4) improve the warning message to show 'X tokens out of Y' so users notice which model tier they're on. **Keep context_window in .ctxrc** for: first-hook-of-session (no JSONL yet), unknown model IDs not in mapping. Workaround until implemented: set context_window: 1000000 in .ctxrc manually. #added:2026-02-27-222206

- [ ] Audit test coverage for export frontmatter preservation — verify T2.1.3 tests exist for: default preserves frontmatter, --force discards it, --skip-existing leaves file untouched, multipart preservation, malformed frontmatter graceful degradation. See specs/future-complete/export-update-mode.md for full checklist. #added:2026-02-26-182446


- [-] Suppress context checkpoint nudges after wrap-up — replaced by Phase 0.9 breakdown below #added:2026-02-24-205402

### Phase 0.9: Suppress Nudges After Wrap-Up

Spec: `specs/suppress-nudges-after-wrap-up.md`. Read the spec before starting any P0.9 task.

**Phase 1 — Plumbing command:**



**Phase 2 — Hook suppression:**


**Phase 3 — Skill integration:**



- [ ] Promote CLI to top-level nav group in zensical.toml: Home | Recipes | CLI | Reference | Operations | Security | Blog — CLI gets the split command pages, Reference keeps conceptual docs (skills, journal format, scratchpad, context files) #added:2026-02-24-204210

- [ ] Split cli-reference.md (1633 lines) into command group pages: cli-overview, cli-init-status, cli-context, cli-recall, cli-tools, cli-system — each page covers a natural command group with its subcommands and flags #added:2026-02-24-204208

- [x] Fix key file naming inconsistency — docs say .context/.context.key, binary says .context/.scratchpad.key. Reconcile naming across code and docs (related to the key relocation task) #added:2026-02-24-201813 #done:2026-02-28

- [ ] Implement ctx recall sync subcommand — propagates locked: true from frontmatter to .state.json and vice versa. Go code exists in internal/cli/recall/sync.go with tests but the command is not registered in Cobra. Docs at cli-reference.md lines 795-816 describe the expected interface #added:2026-02-24-201812

- [ ] Implement ctx remind CLI command — add, list, dismiss subcommands for managing reminders. The check-reminders hook already reads reminders.json but there is no CLI to create or dismiss them. Docs at cli-reference.md lines 1334-1410 describe the expected interface #added:2026-02-24-201810

- [ ] Investigate proactive content suggestions: docs/recipes/publishing.md claims agents suggest blog posts and journal rebuilds at natural moments, but no hook or playbook mechanism exists to trigger this — either wire it up (e.g. post-task-completion nudge) or tone down the docs to match reality #added:2026-02-24-185754

- [ ] Fix enrichment to honor locked state: (1) Add locked: true frontmatter check to /ctx-journal-enrich and /ctx-journal-enrich-all skills — refuse to enrich and tell the user (2) Update docs to clarify that lock protects against both export and enrichment #added:2026-02-24-183246

- [x] Rename .context.key to .ctx.key as part of the key relocation — shorter name aligned with CLI binary name, update all code and doc references from .context.key to .ctx.key #added:2026-02-24-181448 #done:2026-02-28

- [ ] Make encryption key path configurable in .ctxrc (e.g. notify.key_path or crypto.key_path) with default falling back to ~/.local/ctx/keys/<project-hash>.key #added:2026-02-24-172643

- [x] Scan docs for .context/.context.key references and update to reflect new user-level key path — check webhook-notifications.md, scratchpad.md, configuration.md, and any other docs mentioning the key location #added:2026-02-24-172642 #done:2026-02-28

- [ ] Move encryption key to user-level path (~/.local/ctx/keys/<project-hash>.key) instead of .context/.context.key — decouples key from project, removes git-centric assumption, prevents key-next-to-ciphertext antipattern #added:2026-02-24-172517

- [ ] Commit the docs audit changes: nav indexing, ctx brand, parenthetical emphasis, project layout, filename backticks, quoted-term emphasis, drift markers, missing skill entries #added:2026-02-24-171234

- [-] Implement RSS/Atom feed generation for ctx.ist blog — replaced by Phase 0.8 breakdown below #added:2026-02-24-025015

### Phase 0.8: RSS/Atom Feed Generation (`ctx site feed`)

Spec: `specs/rss-feed.md`. Read the spec before starting any P0.8 task.

**Phase 1 — Types and scaffolding:**



**Phase 2 — Blog scanner and feed generator:**



**Phase 3 — CLI wiring:**


**Phase 4 — Tests and integration:**




- [ ] Install golangci-lint on the integration server #for-human #priority:medium #added:2026-02-23 #added:2026-02-23-170213


- [ ] Investigate converting UserPromptSubmit hooks to JSON output — check-persistence, check-ceremonies, check-context-size, check-version, check-resources, and check-knowledge all use plain text with VERBATIM relay. These work differently (prepended to prompt) but may benefit from structured JSON too. #added:2026-02-22-194446

- [ ] Add version-bump relay hook: create a system hook that reminds the agent to bump VERSION, plugin.json, and marketplace.json whenever a feature warrants a version change. The hook should fire during commit or wrap-up to prevent version drift across the three files. #added:2026-02-22-102530


- [ ] Regenerate site HTML after .ctxrc rename #added:2026-02-21-200039



- [ ] Enable webhook notifications in worktrees. Currently `ctx notify`
      silently fails because `.context.key` is gitignored and absent in
      worktrees. For autonomous runs with opaque worktree agents, notifications
      are the one feature that would genuinely be useful. Possible approaches:
      resolve the key via `git rev-parse --git-common-dir` to find the main
      checkout, or copy the key into worktrees at creation time (ctx-worktree
      skill). #priority:medium #added:2026-02-22

- [x] AI: verify and archive completed tasks in TASK.md; the file has gotten
      crowded. Verify each task individually before archiving. #done:2026-02-28

### Phase 0.4: Hook Message Templates

Spec: `specs/future-complete/hook-message-templates.md`. Read the spec before starting any P0.4 task.

**Phase 1 — Core + defaults (no behavioral change):**


**Phase 2 — Discoverability + documentation:**

Spec: `specs/future-complete/hook-message-customization.md`.


### Phase 0.4.9: Injection Oversize Nudge

Spec: `specs/injection-oversize-nudge.md`. Read the spec before starting any P0.4.9 task.


### Phase 0.4.10: Context Window Token Usage

Spec: `specs/context-window-usage.md`. Read the spec before starting any P0.4.10 task.


### Phase 0.5: Spec Scaffolding Skill

- [ ] Create `/ctx-spec` skill — scaffolds a new spec from `specs/spec-template.md`,
      prompts for feature name, creates `specs/{name}.md`, and walks through sections
      with the user (especially edge cases, error handling, validation). Complements
      `/_ctx-brainstorm` (dialogue) by producing the written artifact (document).
      Template: `specs/spec-template.md` #priority:medium #added:2026-02-25

### Prompting Guide — Canonical Reference

- [ ] Add agent/tool compatibility matrix to prompting guide — document which
      patterns degrade gracefully when agents lack file access, CLI tools, or
      ctx integration. Treat as a "works best with / degrades to" table.
      #priority:medium #added:2026-02-25


- [ ] Add versioning/stability note to prompting guide — "these principles are
      stable; examples evolve" + doc date in frontmatter. Needed once the guide
      becomes canonical and people start quoting it. #priority:low #added:2026-02-25

### Phase 0: Ideas (drift markers)

- [ ] Brainstorm: Standardize drift-check comment format and integrate with `/ctx-drift` — currently drift markers (`<!-- drift-check: ... -->`) are ad-hoc shell commands embedded in docs/ARCHITECTURE.md as HTML comments. Formalize the format, teach the drift skill to parse and execute them, and publish the pattern in docs/recipes so any ctx user can add breadcrumbs to their own context files and docs. **Key framing**: markers are a pre-flight check (step 1: automated, fast, catches counting errors like "docs say 13 hooks but code has 17"), NOT a replacement for semantic drift analysis (step 2: reading code, reasoning about stale descriptions, catching convention violations). Marker pass = "no opinion", marker fail = "definite drift". The skill must always do both steps. #priority:medium #added:2026-02-28

### Phase 0: Ideas (from competitive analysis)

- [ ] Brainstorm: JSON Schema for `.ctxrc` — ship a `json-schema.json` that gives IDE users autocompletion and validation for `.ctxrc`. Small YAML surface area; would catch silent typos like `scratchpad_encypt: true`. #priority:low #added:2026-02-28

- [ ] Brainstorm: Lightweight prompt snippets — reusable prompt templates lighter than full skills. Our skills are heavier (full SKILL.md). A "prompt snippet" concept could fill the gap between a skill and a raw instruction. #priority:low #added:2026-02-28

- [ ] Brainstorm: Source-derived context as a complement to authored context — auto-generate ARCHITECTURE.md skeleton from package dependency graph, or a "what changed since last session" summary from git diffs. Would not replace authored context but could bootstrap it. #priority:low #added:2026-02-28

### Phase 0: Ideas

- [ ] Blog: "Building a Claude Code Marketplace Plugin" — narrative from session 
      history, journals, and git diff of feat/plugin-conversion branch. 
      Covers: motivation (shell hooks to Go subcommands), plugin directory 
      layout, marketplace.json, eliminating make plugin, bugs found during 
      dogfooding (hooks creating partial .context/), and the fix. Use 
      /ctx-blog-changelog with branch diff as source material. #added:2026-02-16-111948

**User-Facing Documentation** (from `ideas/done/REPORT-7-documentation.md`):
Docs are feature-organized, not problem-organized. Key structural improvements:


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


### Maintenance

- [x] Human: Ensure the new journal creation /ctx-journal-normalize and
  /ctx-journal-enrich-all works.

- [x] Recipes section needs human review. For example, certain workflows can
  be autonomously done by asking AI "can you record our learnings?" but
  from the documenation it's not clear. Spend as much time as necessary
  on every single recipe.

- [-] Investigate ctx init overwriting user-generated content in .context/
      files. Commit a9df9dd wiped 18 decisions from DECISIONS.md, replacing with
      empty template. — Already fixed: runInit skips existing files unless --force
      is passed (per-file stat check at run.go:90), and prompts for confirmation
      when essential files exist (run.go:51-64). #added:2026-02-06-182205
- [x] Add ctx help command; use-case-oriented cheat sheet for lazy CLI users.
      Should cover: (1) core CLI commands grouped by workflow (getting started,
      tracking decisions, browsing history, AI context), (2) available
      slash-command skills with one-line descriptions, (3) common workflow
      recipes showing how commands and skills combine. One screen,
      no scrolling. Not a skill; a real CLI command. #added:2026-02-06-184257 #done:2026-02-28
- [ ] Add topic-based navigation to blog when post count reaches 15+ #priority:low #added:2026-02-07-015054
- [ ] Revisit Recipes nav structure when count reaches ~25 — consider grouping 
      into sub-sections (Sessions, Knowledge, Security, Advanced) to reduce 
      sidebar crowding. Currently at 18. #priority:low #added:2026-02-20
- [ ] Review hook diagnostic logs after a long session. Check 
      `.context/logs/check-persistence.log` and 
       `.context/logs/check-context-size.log` to verify hooks fire correctly. 
       Tune nudge frequency if needed. #priority:medium #added:2026-02-09
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
