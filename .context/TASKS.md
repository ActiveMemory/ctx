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

- [x] Explore: Replace prompt counter heuristic with actual JSONL token counts in check-context-size hook — parse session JSONL mid-session via recall/parser to sum token fields instead of incrementing a counter. See ideas/done/toke-count.md for original discussion. #added:2026-02-26-180313 #done:2026-02-27

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

### Phase 0.4.10: Context Window Token Usage

Spec: `specs/context-window-usage.md`. Read the spec before starting any P0.4.10 task.

- [x] P0.4.10.1: Add `ContextWindow` field to `CtxRC`, `DefaultContextWindow = 200000`,
      `ContextWindow()` accessor in rc package. #added:2026-02-27 #done:2026-02-27
- [x] P0.4.10.2: Create `internal/cli/system/session_tokens.go` — JSONL path finder
      with caching, tail reader, usage parser, token formatting. #added:2026-02-27 #done:2026-02-27
- [x] P0.4.10.3: Create `check-context-size/window.txt` template + registry entry.
      #added:2026-02-27 #done:2026-02-27
- [x] P0.4.10.4: Modify `check_context_size.go` — token reading, >80% independent
      trigger, token line in checkpoint box, window warning template.
      #added:2026-02-27 #done:2026-02-27
- [x] P0.4.10.5: Write tests — `session_tokens_test.go` (13 tests), additions to
      `check_context_size_test.go` (4 tests). #added:2026-02-27 #done:2026-02-27
- [x] P0.4.10.6: Update docs — configuration.md (`context_window` key),
      customizing-hook-messages.md (16 customizable, window variant + vars).
      #added:2026-02-27 #done:2026-02-27
- [x] P0.4.10.7: Verify: `make build`, `make test`, `make lint`.
      #added:2026-02-27 #done:2026-02-27

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

### Phase 0.6: Event Log and Doctor

Spec: `specs/event-log.md`. Read the spec before starting any P0.6 task.

**Phase 1 — Event log infrastructure:**

- [ ] P0.6.1: Add `EventLog bool` field to `CtxRC` in `internal/rc/types.go`,
      `DefaultEventLog = false` in `default.go`, wire into `Default()` and add
      `EventLog()` accessor in `rc.go`.
      DOD: `make build` passes. `rc.EventLog()` returns false by default, true
      when `.ctxrc` has `event_log: true`. #added:2026-02-27

- [ ] P0.6.2: Add constants in `internal/config/dir.go`: `FileEventLog = "events.jsonl"`,
      `FileEventLogPrev = "events.1.jsonl"`, `EventLogMaxBytes = 1 << 20`.
      Add `state/events.jsonl` and `state/events.1.jsonl` to `GitignoreEntries`.
      DOD: Constants compile. `ctx init` adds event log paths to `.gitignore`.
      #added:2026-02-27

- [ ] P0.6.3: Create `internal/eventlog/eventlog.go` with `Append()` and rotation.
      `Append()` is a noop when `rc.EventLog()` is false. Builds `notify.Payload`,
      marshals to JSON, appends line to `.context/state/events.jsonl`. Creates
      `state/` dir if missing. Rotates when file exceeds `EventLogMaxBytes`.
      DOD: `Append()` writes valid JSONL. Noop when disabled. Rotation works
      (current → `.1`, old `.1` removed). State dir auto-created. #added:2026-02-27

- [ ] P0.6.4: Create `internal/eventlog/eventlog_test.go` — 11 tests covering:
      disabled noop, basic append + readback, state dir auto-creation, rotation
      trigger, rotation overwrite of `.1` file, query with no file, filter by
      hook, filter by session, `--last N`, include rotated, corrupt line skip.
      DOD: All 11 tests pass. `go test ./internal/eventlog/` green. #added:2026-02-27

- [ ] P0.6.5: Add `eventlog.Append()` calls to all system hooks that currently
      call `notify.Send()`. Same arguments, parallel call. Files: check_ceremonies,
      check_persistence, check_context_size, check_journal, check_reminders,
      check_knowledge, check_map_staleness, check_version, check_resources,
      context_load_gate, post_commit, qa_reminder, specs_nudge.
      DOD: Every `notify.Send()` in `internal/cli/system/` has a matching
      `eventlog.Append()`. `make build` passes. No behavioral change to hook
      output. #added:2026-02-27

**Phase 2 — `ctx system events` command:**

- [ ] P0.6.6: Create `internal/cli/system/events.go` with `eventsCmd()` and
      `runEvents()`. Flags: `--hook`, `--session`, `--event`, `--last` (default 50),
      `--json`, `--all`. Human format: `timestamp  event  hook  message` columns.
      JSON format: raw JSONL passthrough. Register in `system.go`.
      DOD: `ctx system events` outputs last 50 events in human-readable format.
      All filter flags work (intersection). `--json` outputs raw JSONL. `--all`
      includes rotated file. "No events logged." when file missing. #added:2026-02-27

- [ ] P0.6.7: Create `internal/cli/system/events_test.go` — 4 tests: default
      human output, `--json` raw output, no log file message, combined filter
      intersection.
      DOD: All tests pass. `go test ./internal/cli/system/` green. #added:2026-02-27

**Phase 3 — `ctx doctor` command:**

- [ ] P0.6.8: Create `internal/cli/doctor/doctor.go` with `Cmd()`, `runDoctor()`,
      and individual check functions. Checks: context initialized, required files
      present, drift detected (via `drift.Detect()`), hook config valid, event
      logging status, webhook configured, pending reminders count, task completion
      ratio, context token size, last event timestamp. `--json` flag for
      machine-readable output. Register in `bootstrap.go`.
      DOD: `ctx doctor` outputs structured health report with categories
      (Structure, Quality, Hooks, State), status indicators (ok/warning/error),
      and summary line. `--json` outputs valid JSON matching `Report` struct.
      All 10 checks implemented. #added:2026-02-27

- [ ] P0.6.9: Create `internal/cli/doctor/doctor_test.go` — 6 tests: healthy
      project (all pass), no `.context/` (structure error), drift warnings
      surfaced, event log off (info note not error), `--json` valid output,
      high task completion ratio (warning).
      DOD: All tests pass. `go test ./internal/cli/doctor/` green. #added:2026-02-27

**Phase 4 — `/ctx-doctor` skill:**

- [ ] P0.6.10: Create `internal/assets/claude/skills/ctx-doctor/SKILL.md`.
      Trigger phrases: "diagnose", "troubleshoot", "doctor", "health check",
      "why didn't my hook fire?", "hooks seem broken", "context seems stale".
      Diagnostic playbook: (1) run `ctx doctor --json` for structural baseline,
      (2) run `ctx system events --json --last 100` if event logging enabled,
      (3) correlate findings across sources, (4) present structured findings
      with evidence, (5) suggest actionable next steps without auto-fixing.
      Graceful degradation: works without event log, notes reduced capability.
      DOD: Skill embedded in binary (`make build`). Frontmatter valid. When to
      Use / When NOT to Use sections present. Playbook covers all 6 data sources
      from spec. Degradation path documented. #added:2026-02-27

- [ ] P0.6.11: Verify full suite: `make build`, `make test`, `make lint`.
      DOD: Zero build errors. All tests pass (eventlog, system, doctor packages).
      Zero new lint issues. `ctx doctor` runs successfully on this project.
      `ctx system events` runs (shows "No events logged." or real data if
      `event_log: true`). #added:2026-02-27

**Phase 5 — Documentation:**

- [ ] P0.6.12: Update CLI docs — add `ctx system events` section to
      `docs/cli/system.md` (flags table, examples, human/JSON output format).
      Create `docs/cli/doctor.md` for `ctx doctor` (command syntax, checks
      table, output examples, when-to-use guidance vs `ctx status` vs
      `/ctx-doctor`). Add `ctx doctor` row to `docs/cli/index.md` commands table.
      DOD: All three doc files updated. Command syntax matches implementation.
      Examples are copy-pasteable. Cross-links work. #added:2026-02-27

- [ ] P0.6.13: Update configuration docs — add `event_log` to `.ctxrc`
      reference table in `docs/home/configuration.md` (or equivalent `.ctxrc`
      section in `docs/cli/index.md`). Type: bool, default: false, description
      matches spec.
      DOD: `event_log` documented in the `.ctxrc` reference table. #added:2026-02-27

- [ ] P0.6.14: Add `/ctx-doctor` entry to `docs/reference/skills.md` — name,
      description, trigger phrases.
      DOD: Skill listed with description and trigger phrases matching SKILL.md.
      #added:2026-02-27

- [ ] P0.6.15: Update existing recipes — add event logging mentions to
      `docs/recipes/system-hooks-audit.md` (local alternative to Sheets),
      `docs/recipes/context-health.md` (`ctx doctor` as superset of drift),
      `docs/recipes/webhook-notifications.md` (local complement to webhooks).
      DOD: Each recipe has a paragraph or section mentioning the new feature
      with cross-link to the troubleshooting recipe. No broken links.
      #added:2026-02-27

- [ ] P0.6.16: Create `docs/recipes/troubleshooting.md` recipe — The Problem,
      TL;DR, Commands and Skills table, workflow sections (quick check with
      `ctx doctor`, deep dive with `/ctx-doctor`, raw event inspection),
      Common Problems section (hook not firing, too many nudges, stale context,
      agent not following instructions), prerequisites, See Also links.
      DOD: Recipe follows existing recipe structure (title, icon, banner, TL;DR,
      commands table, workflow steps, tips, see also). Common Problems section
      has 4 subsections with concrete diagnostic steps. #added:2026-02-27

- [ ] P0.6.17: Update `docs/recipes/index.md` — add Troubleshooting entry
      under Maintenance section. Update `zensical.toml` — add nav entries for
      `docs/cli/doctor.md` and `docs/recipes/troubleshooting.md`.
      DOD: Recipe index lists troubleshooting with description and uses list.
      `zensical.toml` has both new nav entries. Site builds without errors.
      #added:2026-02-27

### Phase 0.7: Session Pause

Spec: `specs/session-pause.md`. Read the spec before starting any P0.7 task.

**Phase 1 — Core infrastructure:**

- [x] P0.7.1: Add pause helpers to `internal/cli/system/state.go` —
      `pauseMarkerPath()`, `paused()` (returns turn count, increments counter),
      `pausedMessage()` (graduated reminder string).
      DOD: Helpers compile. `paused()` returns 0 when no marker file exists.
      `paused()` creates/increments counter when marker exists. `pausedMessage()`
      returns `"ctx:paused"` for turns 1–5, longer string for 6+.
      Unit tests in `state_test.go` cover: no marker → 0, marker exists →
      increment, message for turns 1/5/6/100. #added:2026-02-27

- [x] P0.7.2: Create `internal/cli/system/pause.go` — `ctx system pause`
      plumbing command. Reads session ID from stdin (same `readInput()` pattern
      as other hooks) or `--session-id` flag. Creates pause marker file with
      counter initialized to 0.
      DOD: `ctx system pause` creates `secureTempDir()/ctx-paused-{sessionID}`.
      Double-pause resets counter to 0. `make build` passes. #added:2026-02-27

- [x] P0.7.3: Create `internal/cli/system/resume.go` — `ctx system resume`
      plumbing command. Removes pause marker file. Silent no-op if not paused.
      DOD: `ctx system resume` removes marker file. No error when file doesn't
      exist. `make build` passes. #added:2026-02-27

- [x] P0.7.4: Register `pauseCmd()` and `resumeCmd()` in `system.go`.
      DOD: Commands appear in `ctx system --help` as hidden plumbing commands.
      `make build` passes. #added:2026-02-27

- [x] P0.7.5: Create `internal/cli/system/pause_test.go` — tests for pause,
      resume, and counter behavior.
      DOD: Tests cover: pause creates marker, resume removes marker, resume
      when not paused is no-op, double-pause resets counter, `paused()` increments
      on each call, `pausedMessage()` output for turns 1/5/6/100. All tests pass.
      #added:2026-02-27

**Phase 2 — Hook integration:**

- [x] P0.7.6: Add pause check to `check_context_size.go` — early return when
      paused, emit graduated reminder (this hook is the designated single emitter).
      DOD: When paused, hook emits `pausedMessage()` instead of checkpoint/warning.
      Counter increments. Normal behavior when not paused. Test added to
      `check_context_size_test.go`. #added:2026-02-27

- [x] P0.7.7: Add silent pause check to all other pausable hooks —
      `check_ceremonies.go`, `check_persistence.go`, `check_journal.go`,
      `check_reminders.go`, `check_version.go`, `check_resources.go`,
      `check_knowledge.go`, `check_map_staleness.go`, `context_load_gate.go`,
      `qa_reminder.go`, `post_commit.go`, `specs_nudge.go`.
      DOD: Each hook calls `paused()` after reading input. If paused, returns
      nil (no output). Security hooks (`block_non_path_ctx.go`,
      `block_dangerous_commands.go`) and `cleanup_tmp.go` are NOT modified.
      `make build` passes. #added:2026-02-27

- [x] P0.7.8: Verify pause integration — `make build`, `make test`, `make lint`.
      DOD: Zero build errors. All existing + new tests pass. No new lint issues.
      Manual verification: run `ctx system pause` with a session ID, then run
      `ctx system check-context-size` — should emit `ctx:paused`. Run
      `ctx system block-non-path-ctx` — should still fire normally. #added:2026-02-27

**Phase 3 — Top-level commands:**

- [x] P0.7.9: Create `internal/cli/pause/pause.go` — top-level `ctx pause`
      command. Reads session ID from stdin or `--session-id` flag. Delegates to
      `ctx system pause` logic (shared function, not shell-out).
      DOD: `ctx pause` creates pause marker. `ctx pause --help` shows usage.
      `make build` passes. #added:2026-02-27

- [x] P0.7.10: Create `internal/cli/resume/resume.go` — top-level `ctx resume`
      command. Same pattern as pause.
      DOD: `ctx resume` removes pause marker. `ctx resume --help` shows usage.
      `make build` passes. #added:2026-02-27

- [x] P0.7.11: Register in `internal/bootstrap/bootstrap.go`.
      DOD: `ctx --help` shows `pause` and `resume` in command list.
      `make build` passes. #added:2026-02-27

**Phase 4 — Skills:**

- [x] P0.7.12: Create `/ctx-pause` skill template at
      `internal/assets/claude/skills/ctx-pause/SKILL.md`. Trigger phrases:
      "pause ctx", "pause context", "stop the nudges", "quiet mode".
      Runs `ctx pause`, confirms with short message.
      DOD: Skill has valid frontmatter. When to Use / When NOT to Use sections.
      `make build` passes (skill embedded). #added:2026-02-27

- [x] P0.7.13: Create `/ctx-resume` skill template at
      `internal/assets/claude/skills/ctx-resume/SKILL.md`. Trigger phrases:
      "resume ctx", "resume context", "turn nudges back on", "unpause".
      Runs `ctx resume`, confirms with short message.
      DOD: Skill has valid frontmatter. When to Use / When NOT to Use sections.
      `make build` passes (skill embedded). #added:2026-02-27

**Phase 5 — Documentation:**

- [x] P0.7.14: Add `ctx pause` / `ctx resume` to CLI docs (`docs/cli/tools.md`
      or appropriate CLI page). Include command syntax, flags table, examples.
      DOD: Both commands documented with synopsis, flags, and 2–3 examples.
      Cross-link to recipe. #added:2026-02-27

- [x] P0.7.15: Add `/ctx-pause` and `/ctx-resume` entries to
      `docs/reference/skills.md` — name, description, trigger phrases.
      DOD: Both skills listed with description and trigger phrases matching
      SKILL.md. #added:2026-02-27

- [x] P0.7.16: Create recipe `docs/recipes/session-pause.md` — "Pausing
      Context Hooks". Covers: the problem (nudge overhead for quick tasks),
      TL;DR, commands table, workflow (pause → work → resume), what gets
      paused vs what doesn't, graduated reminder behavior, tips (resume
      before wrap-up, use for quick investigations).
      DOD: Recipe follows existing structure (title, icon, banner, TL;DR,
      commands table, workflow steps, tips, see also). Explains security
      hooks exemption. Mentions ~8k initial load tradeoff. #added:2026-02-27

- [x] P0.7.17: Update `docs/recipes/index.md` — add Session Pause entry.
      Update `zensical.toml` — add nav entry for the recipe. Update
      `docs/home/configuration.md` if any `.ctxrc` keys are added (none
      expected for v1).
      DOD: Recipe index lists session pause with description. `zensical.toml`
      has nav entry. Site builds without errors. #added:2026-02-27

- [x] P0.7.18: Update `docs/recipes/session-ceremonies.md` and
      `docs/recipes/session-lifecycle.md` — add a note about `/ctx-pause`
      as an escape hatch when ceremonies aren't needed.
      DOD: Each recipe has a short paragraph or admonition mentioning pause
      with cross-link to the pause recipe. #added:2026-02-27

- [x] P0.7.19: Verify full suite: `make build`, `make test`, `make lint`.
      DOD: Zero build errors. All tests pass. No new lint issues. Site builds.
      #added:2026-02-27

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
