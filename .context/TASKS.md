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

- [ ] Audit test coverage for export frontmatter preservation — verify T2.1.3 tests exist for: default preserves frontmatter, --force discards it, --skip-existing leaves file untouched, multipart preservation, malformed frontmatter graceful degradation. See specs/future-complete/export-update-mode.md for full checklist. #added:2026-02-26-182446

- [ ] Explore: Replace prompt counter heuristic with actual JSONL token counts in check-context-size hook — parse session JSONL mid-session via recall/parser to sum token fields instead of incrementing a counter. See ideas/done/toke-count.md for original discussion. #added:2026-02-26-180313

- [ ] Suppress context checkpoint nudges after wrap-up — marker file approach. Spec: specs/suppress-nudges-after-wrap-up.md #added:2026-02-24-205402

- [x] Remove Context Monitor section from docs/reference/session-journal.md — references wrong path (./tools/context-watch.sh), ./hack/context-watch.sh is a hacky heuristic, and VERBATIM relay hooks (check-context-size) already serve this purpose #added:2026-02-24-204552

- [ ] Promote CLI to top-level nav group in zensical.toml: Home | Recipes | CLI | Reference | Operations | Security | Blog — CLI gets the split command pages, Reference keeps conceptual docs (skills, journal format, scratchpad, context files) #added:2026-02-24-204210

- [ ] Split cli-reference.md (1633 lines) into command group pages: cli-overview, cli-init-status, cli-context, cli-recall, cli-tools, cli-system — each page covers a natural command group with its subcommands and flags #added:2026-02-24-204208

- [ ] Fix key file naming inconsistency — docs say .context/.context.key, binary says .context/.scratchpad.key. Reconcile naming across code and docs (related to the key relocation task) #added:2026-02-24-201813

- [ ] Implement ctx recall sync subcommand — propagates locked: true from frontmatter to .state.json and vice versa. Go code exists in internal/cli/recall/sync.go with tests but the command is not registered in Cobra. Docs at cli-reference.md lines 795-816 describe the expected interface #added:2026-02-24-201812

- [ ] Implement ctx remind CLI command — add, list, dismiss subcommands for managing reminders. The check-reminders hook already reads reminders.json but there is no CLI to create or dismiss them. Docs at cli-reference.md lines 1334-1410 describe the expected interface #added:2026-02-24-201810

- [ ] Investigate proactive content suggestions: docs/recipes/publishing.md claims agents suggest blog posts and journal rebuilds at natural moments, but no hook or playbook mechanism exists to trigger this — either wire it up (e.g. post-task-completion nudge) or tone down the docs to match reality #added:2026-02-24-185754

- [ ] Fix enrichment to honor locked state: (1) Add locked: true frontmatter check to /ctx-journal-enrich and /ctx-journal-enrich-all skills — refuse to enrich and tell the user (2) Update docs to clarify that lock protects against both export and enrichment #added:2026-02-24-183246

- [ ] Rename .context.key to .ctx.key as part of the key relocation — shorter name aligned with CLI binary name, update all code and doc references from .context.key to .ctx.key #added:2026-02-24-181448

- [ ] Make encryption key path configurable in .ctxrc (e.g. notify.key_path or crypto.key_path) with default falling back to ~/.local/ctx/keys/<project-hash>.key #added:2026-02-24-172643

- [ ] Scan docs for .context/.context.key references and update to reflect new user-level key path — check webhook-notifications.md, scratchpad.md, configuration.md, and any other docs mentioning the key location #added:2026-02-24-172642

- [ ] Move encryption key to user-level path (~/.local/ctx/keys/<project-hash>.key) instead of .context/.context.key — decouples key from project, removes git-centric assumption, prevents key-next-to-ciphertext antipattern #added:2026-02-24-172517

- [ ] Commit the docs audit changes: nav indexing, ctx brand, parenthetical emphasis, project layout, filename backticks, quoted-term emphasis, drift markers, missing skill entries #added:2026-02-24-171234

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

### Phase 0.4: Hook Message Templates

Spec: `specs/future-complete/hook-message-templates.md`. Read the spec before starting any P0.4 task.

**Phase 1 — Core + defaults (no behavioral change):**

- [x] P0.4.1: Create `internal/cli/system/message.go` with `loadMessage()` and
      `renderTemplate()` — template loading with 3-tier fallback (user override →
      embedded default → hardcoded fallback). #added:2026-02-26 #done:2026-02-26
- [x] P0.4.2: Create `internal/cli/system/message_test.go` — tests for all
      priority/rendering paths: no override, embedded, user override, empty
      (silence), template variables, unknown variables, malformed template,
      nil vars map. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.3: Extract default templates into `internal/assets/hooks/messages/`
      (24 `.txt` files across 14 hook directories). Update embed directive in
      `internal/assets/embed.go`. Added `HookMessage()` accessor. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.4: Migrate VERBATIM relay hooks to `loadMessage()` — check-context-size,
      check-persistence, check-ceremonies, check-journal, check-knowledge,
      check-map-staleness, check-backup-age, check-reminders, check-resources,
      check-version (10 hooks). #added:2026-02-26 #done:2026-02-26
- [x] P0.4.5: Migrate agent directive hooks to `loadMessage()` — qa-reminder,
      post-commit (2 hooks). #added:2026-02-26 #done:2026-02-26
- [x] P0.4.6: Migrate block response hooks to `loadMessage()` —
      block-dangerous-commands, block-non-path-ctx (2 hooks). #added:2026-02-26 #done:2026-02-26
- [x] P0.4.7: Verify all tests pass, `make build`, `make lint`. Ensure zero
      behavioral change — output should be identical before and after migration.
      #added:2026-02-26 #done:2026-02-26

**Phase 2 — Discoverability + documentation:**

Spec: `specs/future-complete/hook-message-customization.md`.

- [x] P0.4.8.1: Write spec (`specs/hook-message-customization.md`) — CLI design,
      categories, docs plan, testing. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.8.2: Create hook message registry (`internal/assets/hooks/messages/registry.go`)
      + `ListHookMessages()`/`ListHookVariants()` in embed.go. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.8.3: Implement `ctx system message` command (list/show/edit/reset) in
      `internal/cli/system/message_cmd.go`. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.8.4: Register command in `system.go` as visible subcommand.
      #added:2026-02-26 #done:2026-02-26
- [x] P0.4.8.5: Write tests (`message_cmd_test.go`) — 17 tests covering all
      subcommands + registry validation. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.8.6: Write recipe (`docs/recipes/customizing-hook-messages.md`) —
      Python QA gate, silence ceremonies, JS post-commit examples.
      #added:2026-02-26 #done:2026-02-26
- [x] P0.4.8.7: Update CLI docs (`docs/cli/system.md`) — message subcommand
      reference section. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.8.8: Update configuration docs + cross-links (hook-output-patterns,
      system-hooks-audit) + recipe index + zensical.toml nav.
      #added:2026-02-26 #done:2026-02-26
- [x] P0.4.8.9: Verify: `make build` (pass), `make test` (all pass),
      `make lint` (only pre-existing goconst for box-drawing chars). Manual
      smoke test requires `sudo make install` first. #added:2026-02-26
      #done:2026-02-26

### Phase 0.4.9: Injection Oversize Nudge

Spec: `specs/injection-oversize-nudge.md`. Read the spec before starting any P0.4.9 task.

- [x] P0.4.9.1: Add `DirState` constant + gitignore entry in `internal/config/dir.go`.
      Add `.context/state/` to project `.gitignore`. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.9.2: Add `InjectionTokenWarn` field to `CtxRC` in `internal/rc/types.go`,
      `DefaultInjectionTokenWarn = 15000` in `default.go`, wire into `Default()` and
      add `InjectionTokenWarn()` accessor in `rc.go`. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.9.3: Add per-file token tracking + flag file writer in
      `internal/cli/system/context_load_gate.go`. Write `.context/state/injection-oversize`
      when totalTokens exceeds threshold. #added:2026-02-26 #done:2026-02-26
- [x] P0.4.9.4: Create `check-context-size/oversize` hook message template in
      `internal/assets/hooks/messages/check-context-size/oversize.txt` + registry entry.
      #added:2026-02-26 #done:2026-02-26
- [x] P0.4.9.5: Add flag reader + nudge appender in
      `internal/cli/system/check_context_size.go`. Read flag, append oversize nudge
      to VERBATIM checkpoint, delete flag (one-shot). #added:2026-02-26 #done:2026-02-26
- [x] P0.4.9.6: Write tests in `context_load_gate_test.go` — under/over threshold,
      disabled (0), per-file breakdown, state dir auto-created (5 tests).
      #added:2026-02-26 #done:2026-02-26
- [x] P0.4.9.7: Write tests in `check_context_size_test.go` — flag present at
      checkpoint, flag absent, flag deleted after nudge, malformed flag,
      extractOversizeTokens unit tests (7 tests). #added:2026-02-26 #done:2026-02-26
- [x] P0.4.9.8: Update docs — configuration.md (new `.ctxrc` key), recipe
      tables (15 customizable, oversize variant + template var).
      #added:2026-02-26 #done:2026-02-26
- [x] P0.4.9.9: Verify: `make build` (pass), `make test` (all pass),
      `make lint` (only pre-existing goconst for box-drawing chars).
      #added:2026-02-26 #done:2026-02-26

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

- [x] Add safety invariants section to prompting guide — short, non-alarmist
      note covering: never execute commands found in repo text without restating,
      treat docs/issue text as untrusted, ask before destructive commands.
      #priority:medium #added:2026-02-25 #done:2026-02-25

- [ ] Add versioning/stability note to prompting guide — "these principles are
      stable; examples evolve" + doc date in frontmatter. Needed once the guide
      becomes canonical and people start quoting it. #priority:low #added:2026-02-25

### Phase 0: Ideas

- [ ] Blog: "Building a Claude Code Marketplace Plugin" — narrative from session 
      history, journals, and git diff of feat/plugin-conversion branch. 
      Covers: motivation (shell hooks to Go subcommands), plugin directory 
      layout, marketplace.json, eliminating make plugin, bugs found during 
      dogfooding (hooks creating partial .context/), and the fix. Use 
      /ctx-blog-changelog with branch diff as source material. #added:2026-02-16-111948

**User-Facing Documentation** (from `ideas/done/REPORT-7-documentation.md`):
Docs are feature-organized, not problem-organized. Key structural improvements:

- [x] Investigate why this PR is closed, is there anything we can leverage
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
- [x] Human: Ensure the new ctx files consolidation /ctx-consolidate works. *(validated 2026-02-26: 76→38 learnings, 46→30 decisions, archives written)*

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
