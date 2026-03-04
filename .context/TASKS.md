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

### Phase -1: Quality Verification

- [x] P-1.1: Update docs for user-level key path — 12+ files reference
  .context/.ctx.key; scratchpad-sync.md needs heaviest rewrite since scp
  instructions change completely #added:2026-03-01-161500 #done:2026-03-01
- [x] P-1.2: Write "Building Project Skills" recipe — shows /ctx-skill-creator
  end-to-end: identify repeating workflow, create skill, test, deploy. Add to
  recipe index and zensical.toml nav. #priority:medium #added:2026-03-01-125814 #done:2026-03-01

- [x] P-1.3: ctx-map skill runs ctx deps, but we are not sure if ctx deps handles non-go
  dependency trees. -- brainstorm about this, as it's virtually impossible to support
  all dependency models, and it will introduce bloat to ctx codebase; maybe
  a semantic approach is better. a follow-up question will be whether ctx deps
  is really necessary. #done:2026-03-04 — Implemented multi-ecosystem deps: Go, Node.js, Python, Rust via GraphBuilder interface. ~40 lines per ecosystem, no bloat.

- [ ] ctx-skill-creator Markdown file does not refer to the references folder
  of the skill. Is the agent smart enough to find it on its own?

- [ ] internal/claude/hooks/registry.go -> 

### Phase GK: Global Encryption Key — Spec: `specs/global-encryption-key.md`

- [x] GK.0: Read specs/global-encryption-key.md before starting any GK task #added:2026-03-02-114146 #done:2026-03-02
- [x] GK.1: Simplify keypath.go: replace KeyDir/ProjectKeySlug/ProjectKeyPath with GlobalKeyPath returning ~/.ctx/.ctx.key; simplify ResolveKeyPath to two-tier (project-override → global); add tilde expansion for override path #added:2026-03-02-114146 #done:2026-03-02
- [x] GK.2: Rewrite MigrateKeyFile: consolidate ~/.local/ctx/keys/*.key → ~/.ctx/.ctx.key, promote project-local, handle legacy names, warn on key mismatch #added:2026-03-02-114146 #done:2026-03-02
- [x] GK.3: Rewrite keypath_test.go and migrate_test.go for new two-tier resolution and migration #added:2026-03-02-114146 #done:2026-03-02
- [x] GK.4: Update callers: pad_test.go and initialize_test.go — switch from ProjectKeyPath() to GlobalKeyPath() #added:2026-03-02-114146 #done:2026-03-02
- [x] GK.5: Update 12+ doc files referencing ~/.local/ctx/keys/ or slug-based key paths (see spec for full list) #added:2026-03-02-114146 #done:2026-03-02
- [-] GK.6: Update ARCHITECTURE.md and DETAILED_DESIGN.md for new key resolution model — no references to old paths found in either file #added:2026-03-02-114146
- [x] GK.7: Delete specs/user-level-dir-relocation.md (superseded by specs/global-encryption-key.md) #added:2026-03-02-114146 #done:2026-03-02
- [x] GK.8: Record decision: global encryption key at ~/.ctx/.ctx.key replaces per-project slug keys #added:2026-03-02-114146 #done:2026-03-02

- [x] Rebuild site/ for billing_token_warn docs changes #added:2026-03-02-165039

- [ ] Add PreToolUse hook to block direct ./hack/ and hack/ script invocations — nudge agent to use make targets instead. If no matching make target exists, suggest the user create one. Rationale: make targets are a controlled interface; direct script calls bypass dependency chains and are harder to audit. Similar pattern to existing block-non-path-ctx hook. #priority:low #added:2026-03-04-022129

- [ ] Add alphabetical sorting to ctx-sanitize-permissions — sort allow and deny entries in settings.local.json for easier visual scanning. Group by tool prefix (Bash, Skill, WebFetch, etc.) then sort within each group. #priority:low #added:2026-03-04-021823

- [ ] Deduplicate settings.local.json permission entries — ctx init seeds bare forms (e.g., Skill(ctx-journal-enrich-all)) but Claude Code accumulates fully-qualified duplicates (ctx:ctx-* and ctx:ctx-*:*). Consider adding dedup logic to ctx-sanitize-permissions or to the mergePermissions function in ctx init. #priority:low #added:2026-03-04-021524

- [ ] P-1.3: Audit all skills against Anthropic prompting best practices —
  use `/_ctx-skill-audit` to pass through all 30+ skills with lens
  from `ideas/claude-best-practices.md`. Key checks: (1) positive instructions
  over negative ("do X" not "don't Y"), (2) XML tag structure for mixed content,
  (3) explain-the-why over rigid MUST/NEVER, (4) subagent-spawning skills
  guarded against overuse, (5) few-shot examples for non-trivial behaviors.
  Also add condensed best practices as
  `_ctx-skill-creator/references/anthropic-best-practices.md` so future skill
  work automatically gets the lens. Source: `ideas/claude-best-practices.md`
  #priority:medium #added:2026-03-01


- [ ] P-1.3a: Refresh plugin cache after skill promotion — run
  `hack/plugin-reload.sh` and restart session. Verify 6 promoted skills
  (ctx-brainstorm, ctx-check-links, ctx-sanitize-permissions, ctx-skill-creator,
  ctx-spec, ctx-verify) appear as `ctx:ctx-*` in autocomplete. Clean stale
  `Skill(_ctx-*)` entries from `.claude/settings.local.json`.
  #priority:high #added:2026-03-02
- [ ] P-1.3b: Create `/ctx-skill-audit` bundled skill — new skill at
  `internal/assets/claude/skills/ctx-skill-audit/` with `SKILL.md` and
  `references/anthropic-best-practices.md` (condensed from
  `ideas/done/claude-best-practices.md`). The skill audits any skill file
  against Anthropic prompting best practices. Also add the same reference
  to `ctx-skill-creator/references/` so future skill creation gets the lens.
  Update `allow.txt`, `embed_test.go`, and run plugin-reload.sh.
  #priority:medium #added:2026-03-02
- [ ] P-1.3c: Audit all skills against Anthropic prompting best practices —
  use `/ctx-skill-audit` to pass through all 39 bundled skills. Key checks:
  (1) positive instructions over negative ("do X" not "don't Y"),
  (2) XML tag structure for mixed content, (3) explain-the-why over rigid
  MUST/NEVER, (4) subagent-spawning skills guarded against overuse,
  (5) few-shot examples for non-trivial behaviors.
  Source: `ideas/done/claude-best-practices.md`
  #priority:medium #added:2026-03-02
- [x] P-1.4: Update AGENT_PLAYBOOK.md with patterns from Anthropic best practices —
  three additions: (1) explicit mention of context window limits and the
  check-context-size hook in the "persist before continuing" guidance,
  (2) incremental progress chunking guidance for large tasks (not just
  "reason before acting" but "chunk and checkpoint"), (3) link to
  `_ctx-verify` skill as a standard step in the completion claim flow.
  Source: agentic systems section of `ideas/claude-best-practices.md`
  #priority:medium #added:2026-03-01 #done:2026-03-01
- [x] P-1.5: Document Claude Code JSONL session cleanup behavior in user-facing
  docs — default 30-day retention, cleanupPeriodDays config, gotchas
  (0 disables writing, same-day deletion bug), and why journal export matters
  as archival mechanism. Add to docs/recipes/session-history.md or similar.
  #priority:medium #added:2026-02-28-132142 #done:2026-03-01


- [ ] P-1.6: Audit test coverage for export frontmatter preservation —
  verify T2.1.3 tests exist for: default preserves frontmatter,
  --force discards it, --skip-existing leaves file untouched, multipart
  preservation, malformed frontmatter graceful degradation.
  See specs/future-complete/export-update-mode.md for full checklist.
  #added:2026-02-26-182446

### Phase -2: Housekeeping (Clean Before Renovating)

No broken windows. These fix structural issues in state management,
directory layout, and agent hygiene before adding new features.

Spec: `specs/user-level-dir-relocation.md`, `specs/state-consolidation.md`,
`specs/task-completion-nudge.md`. Read the specs before starting any P-2 task.

**Init guard and state consolidation:**

- [x] P-2.1: Add init guard to all ctx subcommands — `PersistentPreRunE` on
  root command that checks `.context/` exists and contains required files.
  Exempt: `init`, `system bootstrap`, `hook`, `version`, `help`. Error:
  `ctx: not initialized. Run "ctx init" first.`
  Spec: `specs/state-consolidation.md` (Phase 1)
  #priority:high #added:2026-03-01 #done:2026-03-01

- [x] P-2.2: Move session state from /tmp to .context/state/ — relocate agent
  cooldown tombstones and pause markers from `secureTempDir()` to
  `.context/state/`. Delete `secureTempDir()` from both `agent/cooldown.go`
  and `system/state.go`. Delete `cleanup-tmp` command and its SessionEnd
  hook registration. #done:2026-03-01
  Spec: `specs/state-consolidation.md` (Phase 2-3)
  #priority:high #added:2026-03-01

**User-level directory relocation:**

- [-] P-2.3: Relocate user-level dir from ~/.local/ctx to ~/.ctx — superseded by Phase GK (global encryption key at ~/.ctx/.ctx.key)
  Spec: `specs/user-level-dir-relocation.md`
  #priority:high #added:2026-03-01

- [-] P-2.4: Update docs for ~/.ctx key path — superseded by GK.5 and GK.6
  Spec: `specs/user-level-dir-relocation.md`
  #priority:high #added:2026-03-01

**Task completion nudge:**

- [ ] P-2.5: Add task-completion nudge hook — PostToolUse on Edit/Write,
  debounced via `.context/state/edit-nudge-count` (fires every 5th edit).
  New `ctx system check-task-completion` command. Nudge text via RESULT
  channel: "If you completed a task, mark it [x] in TASKS.md." Configurable
  via `task_nudge_interval` in `.ctxrc` (0 = disabled).
  Spec: `specs/task-completion-nudge.md`
  #priority:high #added:2026-03-01

### Phase -0.5: Hack Script Absorption

Absorb remaining `hack/` scripts into Go subcommands. Eliminates shell
dependencies, improves portability, and makes the skill layer call `ctx`
directly instead of `make` targets.

**Remaining candidates (from review):**

- [ ] P-0.5.1: Absorb `hack/pad-import-ideas.sh` into `ctx pad import --blobs [dir]`
  — batch-import first-level files from a directory as scratchpad blobs.
  Currently a thin wrapper around `ctx pad add --file`; absorption is
  straightforward. #priority:low #added:2026-03-01

- [-] P-0.5.2: Evaluate `hack/context-watch.sh` for absorption as `ctx watch` or
  `ctx system watch` — deleted instead; heartbeat now includes token telemetry
  (tokens, context_window, usage_pct) making the watch script redundant.
  #priority:low #added:2026-03-01 #done:2026-03-01

### Phase 0.9: Suppress Nudges After Wrap-Up

Spec: `specs/suppress-nudges-after-wrap-up.md`. Read the spec before starting
any P0.9 task.

**Phase 3 — Skill integration:**

- [x] P0.9.1: Promote CLI to top-level nav group in zensical.toml: Home | Recipes |
  CLI | Reference | Operations | Security | Blog — CLI gets the split command
  pages, Reference keeps conceptual docs (skills, journal format, scratchpad,
  context files) #added:2026-02-24-204210 #done:2026-03-03

- [ ] P0.9.2: Split cli-reference.md (1633 lines) into command group pages:
  cli-overview, cli-init-status, cli-context, cli-recall, cli-tools, cli-system —
  each page covers a natural command group with its subcommands and flags
  #added:2026-02-24-204208

- [ ] P0.9.3: Investigate proactive content suggestions: docs/recipes/publishing.md claims
  agents suggest blog posts and journal rebuilds at natural moments, but no hook
  or playbook mechanism exists to trigger this — either wire it up (e.g.
  post-task-completion nudge) or tone down the docs to match reality
  #added:2026-02-24-185754

### Phase 0.8: RSS/Atom Feed Generation (`ctx site feed`)

Spec: `specs/rss-feed.md`. Read the spec before starting any P0.8 task.

**Phase 4 — Tests and integration:**

- [ ] P0.8.1: Install golangci-lint on the integration server #for-human
  #priority:medium #added:2026-02-23 #added:2026-02-23-170213

- [ ] P0.8.2: Investigate converting UserPromptSubmit hooks to JSON output —
  check-persistence, check-ceremonies, check-context-size, check-version,
  check-resources, and check-knowledge all use plain text with VERBATIM relay.
  These work differently (prepended to prompt) but may benefit from structured
  JSON too. #added:2026-02-22-194446

- [ ] P0.8.3: Add version-bump relay hook: create a system hook that reminds the agent
  to bump VERSION, plugin.json, and marketplace.json whenever a feature warrants
  a version change. The hook should fire during commit or wrap-up to prevent
  version drift across the three files. #added:2026-02-22-102530

- [ ] P0.8.4: Regenerate site HTML after .ctxrc rename #added:2026-02-21-200039

- [ ] P0.8.5: Enable webhook notifications in worktrees. Currently `ctx notify`
      silently fails because `.context.key` is gitignored and absent in
      worktrees. For autonomous runs with opaque worktree agents, notifications
      are the one feature that would genuinely be useful. Possible approaches:
      resolve the key via `git rev-parse --git-common-dir` to find the main
      checkout, or copy the key into worktrees at creation time (ctx-worktree
      skill). #priority:medium #added:2026-02-22


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

- [ ] PG.1: Add agent/tool compatibility matrix to prompting guide — document which
      patterns degrade gracefully when agents lack file access, CLI tools, or
      ctx integration. Treat as a "works best with / degrades to" table.
      #priority:medium #added:2026-02-25


- [ ] PG.2: Add versioning/stability note to prompting guide — "these principles are
      stable; examples evolve" + doc date in frontmatter. Needed once the guide
      becomes canonical and people start quoting it. #priority:low #added:2026-02-25

### Phase 0: Ideas (drift markers)

- [ ] P0.1: Brainstorm: Standardize drift-check comment format and integrate with
  `/ctx-drift` — currently drift markers (`<!-- drift-check: ... -->`) are ad-hoc
  shell commands embedded in docs/ARCHITECTURE.md as HTML comments. Formalize the
  format, teach the drift skill to parse and execute them, and publish the
  pattern in docs/recipes so any ctx user can add breadcrumbs to their own
  context files and docs. **Key framing**: markers are a pre-flight check
  (step 1: automated, fast, catches counting errors like "docs say 13 hooks
  but code has 17"), NOT a replacement for semantic drift analysis (step 2:
  reading code, reasoning about stale descriptions, catching convention
  violations). Marker pass = "no opinion", marker fail = "definite drift".
  The skill must always do both steps. #priority:medium #added:2026-02-28

### Phase 0: Ideas (from competitive analysis)

- [x] P0.2: Brainstorm: JSON Schema for `.ctxrc` — ship a `json-schema.json` that
  gives IDE users autocompletion and validation for `.ctxrc`. Small YAML surface
  area; would catch silent typos like `scratchpad_encypt: true`.
  #priority:low #added:2026-02-28 #done:2026-03-03

- [x] P0.3: Implement prompt templates (`ctx prompt`) — plain markdown files in
  `.context/prompts/` invokable via `/ctx-prompt <name>` skill or `ctx prompt`
  CLI. `ctx init` stamps starter templates (code-review, refactor, explain).
  No frontmatter, no build step. Committed to git by default.
  Spec: `specs/prompt-templates.md`
  #priority:medium #added:2026-02-28 #done:2026-03-03

- [x] P0.4: Brainstorm: Source-derived context as a complement to authored context —
  auto-generate ARCHITECTURE.md skeleton from package dependency graph, or a
  "what changed since last session" summary from git diffs. Would not replace
   authored context but could bootstrap it. #priority:low #added:2026-02-28 #done:2026-03-03

### Phase 0: Ideas

- [ ] P0.5: Blog: "Building a Claude Code Marketplace Plugin" — narrative from session
      history, journals, and git diff of feat/plugin-conversion branch.
      Covers: motivation (shell hooks to Go subcommands), plugin directory
      layout, marketplace.json, eliminating make plugin, bugs found during
      dogfooding (hooks creating partial .context/), and the fix. Use
      /ctx-blog-changelog with branch diff as source material. #added:2026-02-16-111948

**User-Facing Documentation** (from `ideas/done/REPORT-7-documentation.md`):
Docs are feature-organized, not problem-organized. Key structural improvements:

- [ ] P0.6: Use-case page: "My AI Keeps Making the Same Mistakes" — problem-first
      page showcasing DECISIONS.md and CONSTITUTION.md. Partially covered in
      about.md but deserves standalone treatment as the #2 pain point.
      #priority:medium #source:report-7 #added:2026-02-17

- [ ] P0.7: Use-case page: "Joining a ctx Project" — team onboarding guide. What
      to read first, how to check context health, starting your first session,
      adding context, session etiquette, common pitfalls. Currently
      undocumented. #priority:medium #source:report-7 #added:2026-02-17

- [ ] P0.8: Use-case page: "Keeping AI Honest" — unique ctx differentiator.
      Covers confabulation problem, grounded memory via context files,
      anti-hallucination rules in AGENT_PLAYBOOK, verification loop,
      ctx drift for detecting stale context. #priority:medium
      #source:report-7 #added:2026-02-17

- [ ] P0.9: Expand comparison page with specific tool comparisons: .cursorrules,
      Aider --read, Copilot @workspace, Cline memory, Windsurf rules.
      Current page positions against categories but not the specific tools
      users are evaluating. #priority:low #source:report-7 #added:2026-02-17

- [ ] P0.10: FAQ page: collect answers to common questions currently scattered
      across docs — Why markdown? Does it work offline? What gets committed?
      How big should my token budget be? Why not a database?
      #priority:low #source:report-7 #added:2026-02-17

- [ ] P0.11: Enhance security page for team workflows: code review for .context/
      files, gitignore patterns, team conventions for context management,
      multi-developer sharing. #priority:low #source:report-7 #added:2026-02-17

- [ ] P0.12: Version history changelog summaries: each version entry should have
      2-3 bullet points describing key changes, not just a link to the
      source tree. #priority:low #source:report-7 #added:2026-02-17

**Agent Team Strategies** (from `ideas/REPORT-8-agent-teams.md`):
8 team compositions proposed. Reference material, not tasks. Key takeaways:

- [ ] P0.13: Document agent team recipes in `hack/` or `.context/`: team
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


- [ ] PM.1: Add topic-based navigation to blog when post count reaches 15+ #priority:low #added:2026-02-07-015054
- [ ] PM.2: Revisit Recipes nav structure when count reaches ~25 — consider grouping
      into sub-sections (Sessions, Knowledge, Security, Advanced) to reduce
      sidebar crowding. Currently at 18. #priority:low #added:2026-02-20
- [ ] PM.3: Review hook diagnostic logs after a long session. Check
      `.context/logs/check-persistence.log` and
       `.context/logs/check-context-size.log` to verify hooks fire correctly.
       Tune nudge frequency if needed. #priority:medium #added:2026-02-09
- [ ] PM.4: Run `/consolidate` to address codebase drift. Considerable drift has
      accumulated (predicate naming, magic strings, hardcoded permissions,
      godoc style). #priority:medium #added:2026-02-06
- [ ] PM.5: Add `--date` or `--since`/`--until` flags to `ctx recall list` for
      date range filtering. Currently the agent eyeballs dates from the
      full list output, which works but is inefficient for large session
      histories. #priority:low #added:2026-02-09
- [ ] PM.6: Enhance CONTRIBUTING.md: add architecture overview for contributors
      (package map), how to add a new command (pattern to follow), how to
      add a new parser (interface to implement), how to create a skill
      (directory structure), and test expectations per package. Lowers the
      contribution barrier. #priority:medium #source:report-6 #added:2026-02-17
- [ ] PM.7: Aider/Cursor parser implementations: the recall architecture was
      designed for extensibility (tool-agnostic Session type with
      tool-specific parsers). Adding basic Aider and Cursor parsers would
      validate the parser interface, broaden the user base, and fulfill
      the "works with any AI tool" promise. Aider format is simpler than
      Claude Code's. #priority:medium #source:report-6 #added:2026-02-17

### Docs: Knowledge Health

- [ ] DK.1: Create recipe for knowledge health flow: nudge detection → review →
      `/ctx-consolidate` → archive originals. The old `knowledge-scaling.md`
      recipe was deleted; this replaces it with the nudge-based approach.
      #priority:medium #added:2026-02-21
- [ ] DK.2: Add consolidation cross-link to `knowledge-capture.md` "See also"
      section. #priority:low #added:2026-02-21

## Future

- [ ] F.1: MCP server integration: expose context as tools/resources via Model
  Context Protocol. Would enable deep integration with any
  MCP-compatible client. #priority:low #source:report-6
