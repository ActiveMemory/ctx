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

- [ ] Update docs for user-level key path — 12+ files reference .context/.ctx.key; scratchpad-sync.md needs heaviest rewrite since scp instructions change completely #added:2026-03-01-161500

- [ ] Write "Building Project Skills" recipe — shows /ctx-skill-creator end-to-end: identify repeating workflow, create skill, test, deploy. Add to recipe index and zensical.toml nav. #priority:medium #added:2026-03-01-125814

- [x] Write "Design Before Coding" recipe — chains /ctx-brainstorm, /ctx-spec, /ctx-add-task, /ctx-implement into a workflow recipe at docs/recipes/design-before-coding.md. Add to recipe index and zensical.toml nav. #priority:medium #added:2026-03-01-125812 #done:2026-03-01

- [x] Update ctx init --force smoke test to verify no tools/ or sessions/ directories are created in .context/ — these were removed from init but no automated test enforces their absence #added:2026-03-01-112548 #done:2026-03-01

- [x] Elevate six private skills to bundled ctx:ctx-* plugin skills: _ctx-spec, _ctx-brainstorm, _ctx-verify, _ctx-skill-creator, _ctx-check-links, _ctx-sanitize-permissions. Move from .claude/skills/ to internal/assets/claude/skills/, update cross-references between skills, and test deployment via ctx init. #priority:medium #added:2026-03-01-075815 #done:2026-03-01

- [x] Review remaining hack/ scripts for ctx absorption candidates (pad-import-ideas.sh could become ctx pad import --blobs; context-watch.sh monitors raw JSONL token usage which no ctx command covers yet) #added:2026-02-28-184800 #done:2026-03-01

- [ ] Audit all skills against Anthropic prompting best practices — use `/_ctx-skill-audit` to pass through all 30+ skills with lens from `ideas/claude-best-practices.md`. Key checks: (1) positive instructions over negative ("do X" not "don't Y"), (2) XML tag structure for mixed content, (3) explain-the-why over rigid MUST/NEVER, (4) subagent-spawning skills guarded against overuse, (5) few-shot examples for non-trivial behaviors. Also add condensed best practices as `_ctx-skill-creator/references/anthropic-best-practices.md` so future skill work automatically gets the lens. Source: `ideas/claude-best-practices.md` #priority:medium #added:2026-03-01

- [ ] Update AGENT_PLAYBOOK.md with patterns from Anthropic best practices — three additions: (1) explicit mention of context window limits and the check-context-size hook in the "persist before continuing" guidance, (2) incremental progress chunking guidance for large tasks (not just "reason before acting" but "chunk and checkpoint"), (3) link to `_ctx-verify` skill as a standard step in the completion claim flow. Source: agentic systems section of `ideas/claude-best-practices.md` #priority:medium #added:2026-03-01

- [ ] Document Claude Code JSONL session cleanup behavior in user-facing docs — default 30-day retention, cleanupPeriodDays config, gotchas (0 disables writing, same-day deletion bug), and why journal export matters as archival mechanism. Add to docs/recipes/session-history.md or similar. #priority:medium #added:2026-02-28-132142

- [x] Add system resource health check to ctx doctor — call sysinfo.Collect() and report memory/swap/disk/load status as a new 'Resources' category. Use the same threshold logic from check-resources (WARNING at 80%/50%/85%/0.8x, DANGER at 90%/75%/95%/1.5x). Graceful degradation: if sysinfo returns Supported:false for a metric, skip it. Add tests with constructed Snapshot values. #added:2026-02-27-230202 #done:2026-03-01

- [x] BUG: Context checkpoint reports "200k" window size even when using 1M-context model (e.g. claude-opus-4-6[1m]). The check-context-size hook hardcodes DefaultContextWindow=200000. Observed during a deep debugging session where checkpoint said "~96k tokens (~48% of 200k)" but the actual model had a 1M context window. Immediate symptom: premature "wrap up soon" nudges. Related to the auto-detect task below. #priority:high #added:2026-02-28 #done:2026-03-01

- [x] BUG: Skill `ctx-journal-enrich-all` has no fallback heuristic for detecting enrichment. When `ctx system mark-journal --check` fails (e.g. due to the `-q` bootstrap bug below), agents fall back to grepping frontmatter — but the skill doesn't specify which fields distinguish export metadata (`title`, `date`, `time`, `model`, `tokens_in`) from enrichment metadata (`type`, `outcome`, `topics`, `technologies`, `summary`). Add explicit guidance: "Enriched entries have `type`, `outcome`, `topics`, `technologies`, and `summary` in frontmatter. Do NOT use `title` or `date` — these are set by export, not enrichment." #priority:high #added:2026-03-01 #done:2026-03-01

- [x] BUG: Skill `ctx-journal-enrich-all` references `ctx system bootstrap -q` but `-q` flag doesn't exist. The skill template assumes a quiet/machine-readable output mode that was never implemented. Either add `-q` flag to `bootstrap` (output only context_dir path) or update the skill to use `--json` and parse. #priority:high #added:2026-02-28 #done:2026-03-01

- [-] Consider ignoring unknown flags gracefully in CLI commands (especially `system bootstrap`) instead of hard-failing. Skipped: added `-q`/`--quiet` to bootstrap instead — strict flag parsing catches real errors; lenient parsing would mask them and produce wrong-format output silently. #priority:medium #added:2026-02-28 #done:2026-03-01

- [x] Auto-detect context window size from session JSONL model field — the JSONL contains the model name (e.g. "claude-opus-4-5-20251101") which can be mapped to the actual window size (200k for standard, 1M for 1M-context models). Currently defaults to 200k via DefaultContextWindow, causing check-context-size to report '110% full' when a 1M-context model is in use with ~220k tokens. **Resolution**: three-tier fallback: `effective_window = detect_from_jsonl(model) ?? ctxrc.context_window ?? 200_000`. JSONL is ground truth (reflects actual model in use); .ctxrc is fallback for first-hook-of-session (no JSONL yet) or unknown models; 200k is safe last resort. **Approach**: (1) parse model field from JSONL in readSessionTokenUsage, (2) maintain a model-to-window lookup (opus/sonnet standard=200k, 1M suffix=1000000), (3) JSONL detection wins when available, .ctxrc fills in when JSONL can't determine window. (4) improve the warning message to show 'X tokens out of Y' so users notice which model tier they're on. **Keep context_window in .ctxrc** for: first-hook-of-session (no JSONL yet), unknown model IDs not in mapping. Workaround until implemented: set context_window: 1000000 in .ctxrc manually. #added:2026-02-27-222206 #done:2026-03-01

- [ ] Audit test coverage for export frontmatter preservation — verify T2.1.3 tests exist for: default preserves frontmatter, --force discards it, --skip-existing leaves file untouched, multipart preservation, malformed frontmatter graceful degradation. See specs/future-complete/export-update-mode.md for full checklist. #added:2026-02-26-182446

### Phase -0.5: Hack Script Absorption

Absorb remaining `hack/` scripts into Go subcommands. Eliminates shell dependencies, improves portability, and makes the skill layer call `ctx` directly instead of `make` targets.

**Backup → `ctx system backup`:**

- [x] Implement `ctx system backup` subcommand — Go-native tar+gzip archival with `--scope project|global|all` flag. Project scope archives `.context/`, `.claude/`, `ideas/`; global scope archives `~/.claude/` (excluding `todos/`). Config via env vars `CTX_BACKUP_SMB_URL` and `CTX_BACKUP_SMB_SUBDIR` (same interface as current scripts). GVFS mount detection via `os.Stat`, `gio mount` fallback via `exec.Command`. Touch `~/.local/state/ctx-last-backup` on success. #priority:medium #added:2026-03-01 #done:2026-03-01

- [x] Update `_ctx-backup` skill to call `ctx system backup --scope` instead of `make backup*`. Remove Makefile target dependency from skill's `allowed-tools`. #priority:medium #added:2026-03-01 #done:2026-03-01

- [x] Delete `hack/backup-context.sh`, `hack/backup-global.sh`, and related Makefile targets (`backup`, `backup-global`, `backup-all`). #priority:medium #added:2026-03-01 #done:2026-03-01

**Config profiles → `ctx config switch`:**

- [x] Implement `ctx config switch` subcommand — `ctx config switch [dev|base|status]`. Reads `.ctxrc.base` and `.ctxrc.dev` profile files, copies the selected one to `.ctxrc`. Default action (no arg) swaps to the other profile. `status` reports which profile is active. #priority:medium #added:2026-03-01 #done:2026-03-01

- [x] Create `/ctx-config` skill wrapping `ctx config switch` with natural language triggers ("switch to dev mode", "what profile am I on?", "toggle verbose logging"). #priority:medium #added:2026-03-01 #done:2026-03-01

- [x] Delete `hack/ctxrc-swap.sh` and related Makefile targets (`rc-dev`, `rc-base`). #priority:medium #added:2026-03-01 #done:2026-03-01

**Cleanup:**

- [x] Delete `hack/start-dogfood.sh` — obsolete bootstrapper from early development. Its functionality is now covered by `ctx init`. #added:2026-03-01 #done:2026-03-01

**Remaining candidates (from review):**

- [ ] Absorb `hack/pad-import-ideas.sh` into `ctx pad import --blobs [dir]` — batch-import first-level files from a directory as scratchpad blobs. Currently a thin wrapper around `ctx pad add --file`; absorption is straightforward. #priority:low #added:2026-03-01

- [-] Evaluate `hack/context-watch.sh` for absorption as `ctx watch` or `ctx system watch` — deleted instead; heartbeat now includes token telemetry (tokens, context_window, usage_pct) making the watch script redundant. #priority:low #added:2026-03-01 #done:2026-03-01

### Phase 0.9: Suppress Nudges After Wrap-Up

Spec: `specs/suppress-nudges-after-wrap-up.md`. Read the spec before starting any P0.9 task.

**Phase 3 — Skill integration:**

- [ ] Promote CLI to top-level nav group in zensical.toml: Home | Recipes | CLI | Reference | Operations | Security | Blog — CLI gets the split command pages, Reference keeps conceptual docs (skills, journal format, scratchpad, context files) #added:2026-02-24-204210

- [ ] Split cli-reference.md (1633 lines) into command group pages: cli-overview, cli-init-status, cli-context, cli-recall, cli-tools, cli-system — each page covers a natural command group with its subcommands and flags #added:2026-02-24-204208

- [x] Fix key file naming inconsistency — docs say .context/.context.key, binary says .context/.scratchpad.key. Reconcile naming across code and docs (related to the key relocation task) #added:2026-02-24-201813 #done:2026-02-28

- [x] Implement ctx recall sync subcommand — propagates locked: true from frontmatter to .state.json and vice versa. Go code exists in internal/cli/recall/sync.go with tests but the command is not registered in Cobra. Docs at cli-reference.md lines 795-816 describe the expected interface #added:2026-02-24-201812 #done:2026-03-01-000000

- [x] Implement ctx remind CLI command — add, list, dismiss subcommands for managing reminders. The check-reminders hook already reads reminders.json but there is no CLI to create or dismiss them. Docs at cli-reference.md lines 1334-1410 describe the expected interface #added:2026-02-24-201810 #done:2026-03-01

- [ ] Investigate proactive content suggestions: docs/recipes/publishing.md claims agents suggest blog posts and journal rebuilds at natural moments, but no hook or playbook mechanism exists to trigger this — either wire it up (e.g. post-task-completion nudge) or tone down the docs to match reality #added:2026-02-24-185754

- [x] Fix enrichment to honor locked state: (1) Add locked: true frontmatter check to /ctx-journal-enrich and /ctx-journal-enrich-all skills — refuse to enrich and tell the user (2) Update docs to clarify that lock protects against both export and enrichment #added:2026-02-24-183246 #done:2026-03-01

- [x] Rename .context.key to .ctx.key as part of the key relocation — shorter name aligned with CLI binary name, update all code and doc references from .context.key to .ctx.key #added:2026-02-24-181448 #done:2026-02-28

- [x] Make encryption key path configurable in .ctxrc (e.g. notify.key_path or crypto.key_path) with default falling back to ~/.local/ctx/keys/<project-hash>.key #added:2026-02-24-172643 #done:2026-03-01

- [x] Scan docs for .context/.context.key references and update to reflect new user-level key path — check webhook-notifications.md, scratchpad.md, configuration.md, and any other docs mentioning the key location #added:2026-02-24-172642 #done:2026-02-28

- [x] Move encryption key to user-level path (~/.local/ctx/keys/<project-hash>.key) instead of .context/.context.key — decouples key from project, removes git-centric assumption, prevents key-next-to-ciphertext antipattern #added:2026-02-24-172517 #done:2026-03-01

- [ ] Commit the docs audit changes: nav indexing, ctx brand, parenthetical emphasis, project layout, filename backticks, quoted-term emphasis, drift markers, missing skill entries #added:2026-02-24-171234

- [-] Implement RSS/Atom feed generation for ctx.ist blog — replaced by Phase 0.8 breakdown below #added:2026-02-24-025015

### Phase 0.8: RSS/Atom Feed Generation (`ctx site feed`)

Spec: `specs/rss-feed.md`. Read the spec before starting any P0.8 task.

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

**Phase 2 — Discoverability + documentation:**

Spec: `specs/future-complete/hook-message-customization.md`.

### Phase 0.4.9: Injection Oversize Nudge

Spec: `specs/injection-oversize-nudge.md`. Read the spec before starting any P0.4.9 task.

### Phase 0.4.10: Context Window Token Usage

Spec: `specs/context-window-usage.md`. Read the spec before starting any P0.4.10 task.

### Phase 0.6: Plugin Enablement Gap

Ref: `ideas/plugin-enablement-gap.md`. Local-installed plugins get registered in `installed_plugins.json` but not auto-added to `enabledPlugins`, so slash commands are invisible in non-ctx projects.

- [x] P0.6.1: Update `ctx init` output to mention plugin enablement — after the `/plugin install` line, add a note that local installs require manual `enabledPlugins` entry in `~/.claude/settings.json` for cross-project use. Also update docs (getting-started, plugin install instructions). #priority:high #added:2026-03-01 #done:2026-03-01

- [x] P0.6.2: Auto-enable plugin globally on `ctx init` — write `{"enabledPlugins": {"ctx@activememory-ctx": true}}` to `~/.claude/settings.json` (merge, not overwrite) during init. On by default; add `--no-plugin-enable` flag to suppress. Reads `~/.claude/plugins/installed_plugins.json` first to confirm the plugin is actually installed before writing. #priority:high #added:2026-03-01 #done:2026-03-01

- [x] P0.6.3: Add plugin-not-enabled warning to `ctx system bootstrap` — check if `ctx@activememory-ctx` exists in `~/.claude/plugins/installed_plugins.json` but is missing from both `~/.claude/settings.json` and `.claude/settings.local.json` `enabledPlugins`. Emit a one-line warning with the fix command. #priority:medium #added:2026-03-01 #done:2026-03-01

- [x] P0.6.4: Add plugin enablement check to `ctx doctor` — new diagnostic category "Plugin" that detects installed-but-not-enabled state. Reports: plugin installed (yes/no), globally enabled (yes/no), locally enabled (yes/no), and suggests fix if installed but not enabled anywhere. #priority:medium #added:2026-03-01 #done:2026-03-01

### Phase 0.5: Spec Scaffolding Skill

- [x] Create `/ctx-spec` skill — scaffolds a new spec from `specs/spec-template.md`,
      prompts for feature name, creates `specs/{name}.md`, and walks through sections
      with the user (especially edge cases, error handling, validation). Complements
      `/_ctx-brainstorm` (dialogue) by producing the written artifact (document). #done:2026-03-01
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
- [x] `/ctx-journal-enrich-all` should handle export-if-needed: check for
      unexported sessions before enriching and export them automatically,
      so the user can say "process the journal" and the skill handles the
      full pipeline (export → normalize → enrich). #priority:medium #added:2026-02-09 #done:2026-03-01
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
- [x] Fix skills page (`docs/skills.md`): `/ctx-consolidate` entry says
      "runs `ctx reindex`" — should say `ctx learnings reindex` /
      `ctx decisions reindex`. #priority:low #added:2026-02-21
      #done:2026-03-01 — resolved by implementing `ctx reindex`; the reference is now correct
- [ ] Add consolidation cross-link to `knowledge-capture.md` "See also"
      section. #priority:low #added:2026-02-21
- [x] `ctx reindex` convenience command — runs `ctx decisions reindex` and
      `ctx learnings reindex` in one call. Both files grow at similar rates;
      users always want to reindex both. #priority:low #added:2026-02-21
      #done:2026-03-01

## Future

- [ ] MCP server integration: expose context as tools/resources via Model
  Context Protocol. Would enable deep integration with any
  MCP-compatible client. #priority:low #source:report-6
