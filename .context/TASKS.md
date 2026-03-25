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

### Phase -3: DevEx

[ ] Plugin enablement gap: Ref: `ideas/plugin-enablement-gap.md`. 
Local-installed plugins get registered in `installed_plugins.json` but not 
auto-added to `enabledPlugins`, so slash commands are invisible in non-ctx 
projects.

- [ ] Add cobra Example fields to CLI commands via examples.yaml #added:2026-03-20-163413

- [ ] Evaluate Gemini Search MCP server as peer MCP for grounded web queries — try gemini-grounding, document in multi-tool-setup recipe if useful. See ideas/gemini-search-mcp.md #added:2026-03-20-141022

- [ ] Create ctx-docstrings skill: audit and fix docstrings against CONVENTIONS.md Documentation section. Skill loads CONVENTIONS.md, scans functions in scope for missing/incomplete docstring sections (Parameters, Returns), reports violations, and optionally fixes them. Language-agnostic design with Go as first implementation. Deterministic enforcement via linter is tracked separately in ideas/spec-convention-enforcement.md #added:2026-03-16-114445

### Phase -2: Task completion nudge:

- [ ] Design UserPromptSubmit hook that runs `make audit` at session start and 
  surfaces failures as a consolidation-debt warning before the agent acts on stale assumptions. Project-level hook (not bundled in ctx), configurable via .ctxrc or settings.json. Related: consolidation nudge hook spec. #added:2026-03-23-223500

- [ ] Bug: check-version hook missing throttle touch on plugin version read error (run.go:70). When claude.PluginVersion() fails, the hook returns without touching the daily throttle marker, causing repeated checks on days when plugin.json is missing or corrupted. Fix: add internalIo.TouchFile(markerFile) before the early return. See docs/recipes/hook-sequence-diagrams.md check-version diagram which documents the expected behavior. #added:2026-03-23-162802

- [ ] Design UserPromptSubmit hook that runs go build and surfaces compilation errors before the agent acts on stale assumptions #added:2026-03-23-120136

- [ ]: Architecture mapping skill refactoring:
  - [ ] Update ctx-architecture skill based on the following findings; remove 
    gitnexus from the template and the actual skill; have a separate follow-up enrichment
    skill (see the next task where it also has a spec)
        - [2026-03-25-021557] Code intelligence tools trade depth for breadth in 
        architecture analysis
          - **Context**: Compared three sessions analyzing a large codebase 
          (~34k symbols): Session 1 (broken MCP) produced 5,866 lines of 
          DETAILED_DESIGN with per-controller data flows, scale math, 
          startup sequences. Session 2 (full MCP + same skill) produced 1,124 lines 
          (5.2x less). Session 3 (enrichment) added verified graph data but couldn't 
          recover the intimate code knowledge.
          - **Lesson**: When graph query tools are available, agents satisfice 
          instead of maximize. They get structural answers without reading code, 
          missing operational details (defaults, timeouts, scale math, edge cases) 
          that only emerge from line-by-line reading. The tool answers the question 
          asked but prevents discovery of answers to questions never asked.
          - **Application**: Architecture analysis skills should NOT offer MCP tools:
          force code reading first. Use a separate enrichment skill to verify and 
          extend with tools afterward. Constraint is the feature.

- [ ] Architecture Mapping (Enrichment):
  **Context**: Skill that incrementally builds and maintains ARCHITECTURE.md
  and DETAILED_DESIGN.md. Coverage tracked in map-tracking.json.
  Spec: `specs/ctx-architecture.md`
  - [ ] Create ctx-architecture-enrich skill: takes existing /ctx-architecture 
  principal-mode artifacts as baseline, runs comprehensive enrichment pass via 
  GitNexus MCP (blast radius verification, registration site discovery, 
  execution flow tracing, domain clustering comparison, shallow module 
  deep-dive). Spec: `ideas/spec-architecture-enrich.md`. Reference 
  implementation: kubernetes-service enrichment pass 2026-03-25. 
  #added:2026-03-25-120000

- [ ]: ctx-architecture-failure-analysis
      **Context**: Adversarial analysis skill that identifies where a codebase will
      silently betray you. Requires `ctx-architecture` artifacts as input (ARCHITECTURE.md,
      DETAILED_DESIGN*.md, map-tracking.json). Does its own targeted deep reads focusing
      on mutation points, shared mutable state, error swallowing, concurrency, implicit
      ordering, missing enforcement, and scaling cliffs. Uses available tooling (GitNexus,
      Gemini Search) to cross-reference patterns.

      Produces `DANGER-ZONES.md` — a ranked inventory of silent failure points with:
      location, failure mode, blast radius, detection gap, and suggested fix. Two tiers:
      "most likely to cause production incidents" and "less likely but equally dangerous."

      Distinct from a security threat model (which would be `ctx-threat-model` — a
      separate skill for auth bypass, injection, privilege escalation, supply chain).
      This skill focuses on correctness: race conditions, ordering assumptions, cache
      staleness, fan-out amplification, non-atomic ownership, inverted logic,
      force-delete orphans, global state mutation.

      - [ ] Design SKILL.md for ctx-architecture-failure-analysis: inputs 
      (architecture artifacts), analysis phases, output format (DANGER-ZONES.md), 
      quality checklist #added:2026-03-25-060000
      - [ ] Define the adversarial analysis framework: categories of silent 
      failure (concurrency, ordering, cache, amplification, ownership, error 
      swallowing, global state) with heuristics for each #added:2026-03-25-060000
      - [ ] Implement skill with GitNexus integration: use impact analysis for 
      blast radius estimation, use context for shared-state detection #added:2026-03-25-060000
      - [ ] Add Gemini Search integration: cross-reference discovered patterns 
      against known failure modes in similar systems. #added:2026-03-25-060000

- [ ] dependency sanity check: on session start; I should be able to know
  gitnexus mcp is up; gemini mcp is up. maybe a status report during ctx-remember.

- [ ] ctx-architecture-extend
  **Context**: Companion to `ctx-architecture` and `ctx-failure-analysis`, completing
  a trilogy: how does it work → where will it break → where does it grow.
      Reads architecture artifacts → identifies registration patterns (interfaces, factory
      functions, plugin systems, ordered slices, scheme registrations) → traces recent
      additions via git log to confirm which extension points are actually used → produces
      `EXTENSION-POINTS.md` ranked by frequency, with exact file locations, function
      signatures, and the typical feature pattern (e.g., "most features require a variable
      + a mutator + a machine-agent task").

      Valuable for onboarding ("I need to add feature X, where do I start?") and
      architecture review ("are we adding features in the right places?").

      - [ ] Design SKILL.md for ctx-extension-map: inputs (architecture artifacts + 
      git log), analysis phases, output format (EXTENSION-POINTS.md), 
      quality checklist #added:2026-03-25-062000
      - [ ] Define extension point detection heuristics: interface registrations, 
      factory patterns, ordered slices, scheme init blocks, //go:embed directories, 
      feature flag structs with tags #added:2026-03-25-062000
      - [ ] Add git log frequency analysis: trace recent commits to confirm 
     which extension points are actively used vs. dormant #added:2026-03-25-062000
      - [ ] Integrate with GitNexus: use cluster/process data to identify 
      registration call sites and their callers #added:2026-03-25-062000

[ ] drift check should notify if claude permissions have insecure stuff in it.

[ ] task: sync workspace to ARI_INBOX

### Phase -1: Hack Script Absorption

Absorb remaining `hack/` scripts into Go subcommands. Eliminates shell
dependencies, improves portability, and makes the skill layer call `ctx`
directly instead of `make` targets.

### Phase 0.9: Suppress Nudges After Wrap-Up

Spec: `specs/suppress-nudges-after-wrap-up.md`. Read the spec before starting
any P0.9 task.

**Phase 3 — Skill integration:**

- [-] P0.9.2: Split cli-reference.md — moved to Future
  #added:2026-02-24-204208

- [-] P0.9.3: Investigate proactive content suggestions — moved to Future
  #added:2026-02-24-185754

### Phase 0.8: RSS/Atom Feed Generation (`ctx site feed`)

Spec: `specs/rss-feed.md`. Read the spec before starting any P0.8 task.

### Phase 0.4: Hook Message Templates

Spec: `specs/future-complete/hook-message-templates.md`. Read the spec before
starting any P0.4 task.

**Phase 2 — Discoverability + documentation:**

Spec: `specs/future-complete/hook-message-customization.md`.

- [ ] Migrate hook message templates from .txt files to YAML localization #added:2026-03-20-163801

### Phase 0.4.9: Injection Oversize Nudge

Spec: `specs/injection-oversize-nudge.md`. Read the spec before starting
any P0.4.9 task.

### Phase 0.4.10: Context Window Token Usage

Spec: `specs/context-window-usage.md`. Read the spec before starting any
P0.4.10 task.

### Phase 0.5 Cleanup

* Human: internal/recall/parser requires a serious refactoring; for example
  the parser object and its private and public methods need to go to its own
  package and other helper functions need to go to a different adjacent package.
* Human: internal/notify/notify.go requires refactoring (all functions bagged in
  one file; types need to go to types.go per convention etc etc)
* Human: split err package into sub packages.

- [ ] Add Use* constants for all system subcommands #added:2026-03-21-092550

- [ ] Refactor site/cmd/feed: extract helpers and types to core/, make Run public #added:2026-03-21-074859

- [ ] Add Use* constants for all cobra subcommand Use strings #added:2026-03-20-184639

- [ ] Systematic audit: extract all magic flag name strings across CLI commands into config/flag constants #added:2026-03-20-175155

- [ ] Move generic string helpers from cli/add/core/strings.go to internal/format #added:2026-03-20-175046

- [ ] Add missing flag name constants (priority, section, file) and priority level constants (high, medium, low) to config/flag #added:2026-03-20-170842

### Phase 0: Ideas

**User-Facing Documentation** (from `ideas/done/REPORT-7-documentation.md`):
Docs are feature-organized, not problem-organized. Key structural improvements:

**Agent Team Strategies** (from `ideas/REPORT-8-agent-teams.md`):
8 team compositions proposed. Reference material, not tasks. Key takeaways:


- [ ] Scan all config/**/* constants and catalog which ones should be ctxrc entries for user configurability #priority:medium #added:2026-03-22-095552

- [ ] Update user-facing documentation for changed CLI flag shorthands #added:2026-03-21-102755

- [ ] Add Unicode-aware slugification for non-ASCII content #added:2026-03-21-070953

- [ ] Make TitleSlugMaxLen configurable via .ctxrc #added:2026-03-21-070944

- [ ] Spec and implement CRLF-to-LF newline normalization for journal and context files #added:2026-03-20-224845

- [ ] Test ctx on Windows — validate build, init, agent, drift, journal pipeline #added:2026-03-20-224835

- [ ] Evaluate Windows support for sysinfo.Collect and path handling #added:2026-03-20-194930

- [ ] Make doctor thresholds configurable via .ctxrc #added:2026-03-20-194923

- [ ] Evaluate cross-platform path handling in change/core/scan.go — git always uses "/" but UniqueTopDirs should consider filepath.ToSlash for Windows robustness #added:2026-03-20-182103

- [ ] Replace English-only Pluralize helper in change/core/detect.go with i18n-safe approach #added:2026-03-20-180502

- [ ] Replace ASCII-only alnum check in agent/core/score.go with unicode.IsLetter/IsDigit #added:2026-03-20-175943

### Phase S-0: Memory Bridge Groundwork

Prerequisites that unblocked the memory bridge phases.


### Phase MB: Memory Bridge Foundation (`ctx memory`)

Spec: `specs/memory-bridge.md`. Read the spec before starting any MB task.

Bridge Claude Code's auto memory (MEMORY.md) into `.context/` with discovery,
mirroring, and drift detection. Foundation for future import/publish phases.

### Phase MI: Memory Import Pipeline (`ctx memory import`)

Spec: `specs/memory-import.md`. Read the spec before starting any MI task.

Import entries from Claude Code's MEMORY.md into structured `.context/` files
using heuristic classification. Builds on Phase MB foundation (discover, mirror, state).

- [-] MI.future: `--interactive` mode for agent-assisted classification — skipped: `--dry-run` covers review; agents can use `ctx add` directly for overrides; interactive CLI prompts don't compose with agent workflows

### Phase S-3: Blog Post — "Agent Memory is Infrastructure"

Spec: `specs/blog-agent-memory-infrastructure.md`.


### Phase MP: Memory Publish (`ctx memory publish`)

Spec: `specs/memory-publish.md`. Read the spec before starting any MP task.

Push curated context from `.context/` into Claude Code's MEMORY.md so the agent
sees structured project context on session start without needing hooks.

### Phase 9: Context Consolidation Skill `#priority:medium`

**Context**: `/ctx-consolidate` skill that groups overlapping entries by keyword
similarity and merges them with user approval. Originals archived, not deleted.
Spec: `specs/context-consolidation.md`
Ref: https://github.com/ActiveMemory/ctx/issues/19 (Phase 3)

- [ ] Implement consolidation nudge hook: count sessions since last consolidation, nudge after 6. Spec: `specs/consolidation-nudge-hook.md` #added:2026-03-23-223000

- [ ] Auto-record consolidation baseline commit: `/ctx-consolidate` and `ctx system mark-consolidation` should stamp HEAD hash + date into `.context/state/consolidation.json` only on first invocation (write-once until reset). Subsequent consolidation sessions preserve the original baseline. The baseline resets only when the consolidation nudge counter resets (i.e., when a new feature cycle begins). This way multi-pass consolidation keeps the true starting point. Related: `specs/consolidation-nudge-hook.md` #added:2026-03-23-224000

### Phase EM: Extension Map Skill (`/ctx-extension-map`)

question: is this done; or needs planning?

### Phase WC: Write Consolidation

Baseline commit: `4ec5999` (Auto-prune state directory on session start).
Goal: consolidate user-facing messages into `internal/write/` as the central
output package. All CLI commands should route printed output through this package.

- [ ] Migrate moc.go hardcoded strings to YAML or Go templates #added:2026-03-20-214922

- [ ] Design terminal-aware truncation for CLI output #added:2026-03-20-184509

### Phase SP: Configurable Session Prefixes

Spec: `specs/session-prefixes.md`. Read the spec before starting any SP task.

Replace hardcoded `session_prefix` / `session_prefix_alt` pair with a
user-extensible `session_prefixes` list in `.ctxrc`. Parser vocabulary
is not i18n text — it belongs in runtime config.

### Phase EH: Error Handling Audit

Systematic audit of silently discarded errors across the codebase.
Many call sites use `_ =` or `_, _ =` to discard errors without
any feedback. Some are legitimate (best-effort cleanup), most are
lazy escapes that hide failures.

- [ ] Add drift check: MCP prompt coverage vs bundled skills — programmatic check comparing config/mcp/prompt constants against assets.ListSkills() to detect skills without MCP prompt equivalents. Pair with the tool coverage drift check. @CoderMungan #priority:medium #added:2026-03-15-120519

- [ ] MCP v0.3: expand MCP prompts to cover more skills — current 5 prompts (session-start, add-decision, add-learning, reflect, checkpoint) are a subset of ~30 bundled skills. Evaluate which skills benefit from protocol-native MCP prompt equivalents. Decision 2026-03-06 established 'Skills stay CLI-based; MCP Prompts are the protocol equivalent.' @CoderMungan #priority:medium #added:2026-03-15-120519

- [ ] Add drift check: MCP tool coverage vs CLI commands — programmatic check that compares registered MCP tool names (config/mcp/tool) against ctx CLI subcommands to detect newly added CLI commands without MCP equivalents. Could be a drift detector check or a compliance test. @CoderMungan #priority:medium #added:2026-03-15-120116

- [ ] MCP v0.3: expose additional CLI commands as MCP tools — candidates: ctx_load (full context packet), ctx_agent (token-budgeted packet), ctx_reindex (rebuild indices), ctx_sync (reconcile docs/code), ctx_doctor (health check). Evaluate which provide value over the protocol vs requiring terminal interaction. @CoderMungan #priority:medium #added:2026-03-15-120025

- [ ] Make MCP defaults configurable via .ctxrc — add mcp_recall_limit, mcp_truncate_len, mcp_truncate_content_len, mcp_min_word_len, mcp_min_word_overlap fields to .ctxrc schema; expose via rc.MCP*() with fallback to config/mcp/cfg defaults; update tools.go to read from rc instead of cfg constants. @CoderMungan #priority:medium #added:2026-03-15-114700

- [ ] MCP tools.go cleanup pass: magic strings, duplicated fragments, nested templates. Lines: 461:481 + 186:196 duplicated code; 335 magic number; 382:385 nested TextDescs → single template; 390+851 magic time literal; 443+499+800 magic words; 557+892+902 magic numbers; 590+638 nested TextDesc templating; 820 prefixed %s; 854 suffix %s #priority:high #added:2026-03-15-110429

- [ ] EH.0: Create central warning sink — `internal/log/warn.go` with
      `var Sink io.Writer = os.Stderr` and `func Warn(format string, args ...any)`.
      All stderr warnings (`fmt.Fprintf(os.Stderr, ...)`) route through this
      function. The `fmt.Fprintf` return error is handled once, centrally.
      The sink is swappable (tests use `io.Discard`, future: syslog, file).
      EH.2–EH.4 should use `log.Warn()` instead of raw `fmt.Fprintf`.
      DoD: `grep -rn 'fmt.Fprintf(os.Stderr' internal/` returns zero hits
      #priority:high #added:2026-03-15

- [ ] EH.1: Catalogue all silent error discards — recursive walk of `internal/`
      for patterns: `_ = `, `_, _ = `, `//nolint:errcheck`, bare `return` after
      error-producing calls. Group by category:
      (a) file close in defer — often legitimate but should log on failure
      (b) file write/read — data loss risk, must surface
      (c) os.Remove/Rename — state corruption risk
      (d) fmt.Fprint to stderr — truly best-effort, acceptable
      Commands: `grep -rn '_ =' internal/`, `grep -rn 'nolint:errcheck' internal/`
      Output: spreadsheet in `.context/` with file, line, expression, category,
      and recommended action (log-stderr, return-error, acceptable-as-is).
      DoD: every `_ =` in the codebase is categorised and has a recommended action
      #priority:high #added:2026-03-14

- [ ] EH.2: Address category (b) — file write/read discards. These risk silent
      data loss. Fix: return the error, or at minimum emit to stderr with
      `fmt.Fprintf(os.Stderr, "ctx: ...: %v\n", err)` following the pattern
      established in `internal/log/event.go`.
      DoD: no write/read error is silently discarded
      #priority:high #added:2026-03-14

- [ ] EH.3: Address category (a) — file close in defer. Most are `defer func()
      { _ = f.Close() }()`. For read-only files, close errors are rare but
      should still surface. For write/append files, close can fail if the
      final flush fails — these are data loss. Fix: `if err := f.Close();
      err != nil { fmt.Fprintf(os.Stderr, "ctx: close %s: %v\n", path, err) }`.
      DoD: all defer-close sites log failures to stderr
      #priority:medium #added:2026-03-14

- [ ] EH.4: Address category (c) — os.Remove/Rename discards. These are state
      operations (rotation, pruning, temp file cleanup). Silent failure leaves
      stale state. Fix: stderr warning at minimum; for rotation/rename, consider
      returning the error.
      DoD: no Remove/Rename error is silently discarded
      #priority:medium #added:2026-03-14

- [ ] EH.5: Validate — `grep -rn '_ =' internal/` returns only category (d)
      entries (fmt.Fprint to stderr) and entries explicitly annotated as
      acceptable. Run `make lint && make test` to confirm no regressions.
      DoD: grep output is clean or fully annotated; CI green
      #priority:high #added:2026-03-14

- [ ] Add AST-based lint test to detect exported functions with no external callers #added:2026-03-21-070357

- [ ] Audit exported functions used only within their own package and make them private #added:2026-03-21-070346

- [ ] Audit and remove side-effect output from error-returning functions #added:2026-03-20-212212

### Phase ET: Error Package Taxonomy (`internal/err/`)

`errors.go` is 1995 lines with 188 functions in a single file. Split into
domain-grouped files. No API changes — same package, same function signatures,
just file reorganization.

Taxonomy (from prefix analysis):

| File             | Prefixes / Domain                                     | ~Count |
|------------------|-------------------------------------------------------|--------|
| `memory.go`      | Memory*, Discover*                                    | 17     |
| `parser.go`      | Parser*                                               | 7      |
| `crypto.go`      | Crypto*, Encrypt*, Decrypt*, GenerateKey, SaveKey, LoadKey, NoKeyAt | 14     |
| `task.go`        | Task*, NoTaskSpecified, NoTaskMatch, NoCompletedTasks | 8      |
| `journal.go`     | LoadJournalState*, SaveJournalState*, ReadJournalDir, NoJournalDir, NoJournalEntries, ScanJournal, UnknownStage, StageNotSet | 10 |
| `session.go`     | Session*, FindSessions, NoSessionsFound, All*, Ambiguous* | 8 |
| `pad.go`         | Edit*, Blob*, ReadScratchpad, OutFlagRequiresBlob, NoConflict*, Resolve* | 10 |
| `recall.go`      | Reindex*, Stats*, EventLog*                           | 6      |
| `fs.go`          | Read*, Write*, Open*, Stat*, File*, Mkdir*, CreateDir, DirNotFound, NotDirectory, Boundary* | 30 |
| `backup.go`      | Backup*, CreateBackup*, CreateArchive*                | 6      |
| `prompt.go`      | Prompt*, NoPromptTemplate, ListTemplates, ReadTemplate, NoTemplate | 7 |
| `hook.go`        | Embedded*, Override*, UnknownHook, UnknownVariant, MarkerNotFound | 6 |
| `skill.go`       | Skill*                                                | 2      |
| `config.go`      | UnknownProfile, ReadProfile, UnknownFormat, UnknownProjectType, InvalidTool, UnsupportedTool, NotInitialized, ContextNotInitialized, ContextDirNotFound, FlagRequires* | 12 |
| `errors.go`      | Remaining general-purpose: WorkingDirectory, CtxNotInPath, ReadInput, InvalidDate*, Reminder*, Drift*, Git*, Webhook*, etc. | ~25 |




- [ ] Add freshness_files to .ctxrc defaults seeded by ctx init — currently the freshness config is only in the gitignored .ctxrc, so new clones don't get it. Consider a .ctxrc.defaults pattern or seeding via ctx init template. #priority:medium #added:2026-03-14-105143

- [ ] SEC.1: Security-sensitive file change hook — PostToolUse on Edit/Write matching security-critical paths (.claude/settings.local.json, .claude/settings.json, CLAUDE.md, .claude/CLAUDE.md, .context/CONSTITUTION.md). Three actions: (1) nudge user in-session, (2) relay to webhook for out-of-band alerting (autonomous loops), (3) append to dedicated security log (.context/state/security-events.jsonl) for forensics. Separate from general event log. Spec needed. #priority:high #added:2026-03-13

- [ ] O.5: Session timeline view — add --sessions flag to ctx system events. Per-session breakdown of eval/fired counts with hook list. See ideas/spec-hook-observability.md Phase 5 #added:2026-03-12-145401

- [ ] O.4: Doctor hook health check — surface hook activity in ctx doctor output (active/evaluated-never-fired/never-evaluated). See ideas/spec-hook-observability.md Phase 4 #added:2026-03-12-145401

- [ ] O.3: Skip reason logging — add eventlog.Skip() with standard reason constants (paused, throttled, condition-not-met). Instrument 19 hook early-exit paths. See ideas/spec-hook-observability.md Phase 3 #added:2026-03-12-145401

- [ ] O.2: Event summary view — add --summary flag to ctx system events. Aggregates eval/fired counts per hook, shows last-eval/last-fired timestamps, lists never-evaluated hooks. See ideas/spec-hook-observability.md Phase 2 #added:2026-03-12-145401

- [ ] O.1: Hook eval logging — wrap hook cobra commands to log 'eval' events on every invocation. Refactor Run() signatures from os.Stdin to io.Reader (peek+replay pattern). Adds eventlog.Eval(), EventTypeEval constant. See ideas/spec-hook-observability.md Phase 1 #added:2026-03-12-145401

- [ ] Companion intelligence recommendation: implement spec from ideas/spec-companion-intelligence.md — ctx doctor companion detection, ctx init recommendation tip, ctx agent awareness in packets #added:2026-03-12-133008

- [ ] Add configurable assets layer: allow users to plug their own YAML files for localization (language selection, custom text overrides). Currently all user-facing text is hardcoded in commands.yaml; need a mechanism to load user-provided YAML that overlays or replaces built-in text. This enables i18n without forking. #priority:low #added:2026-03-07-233756

- [ ] Cleanup internal/cli/system/core/persistence.go: move 10 (base for ParseInt) to config constant #priority:low #added:2026-03-07-220825

- [ ] Cleanup internal/cli/system/core/session_tokens.go: move SessionStats from state.go to types.go #priority:low #added:2026-03-07-220825

- [ ] Cleanup internal/cli/system/core/wrapup.go: line 18 constant should go to config; make WrappedUpExpiry configurable via ctxrc #priority:low #added:2026-03-07-220825

- [ ] Cleanup internal/cli/system/core/version.go: line 81 newline should come from config #priority:low #added:2026-03-07-220819

- [ ] Add taxonomy to internal/cli/system/core/ — currently an unstructured bag of files; group by domain (backup, hooks, session, knowledge, etc.) #priority:medium #added:2026-03-07-220819

- [ ] Cleanup internal/cli/system/core/version_drift.go: line 53 string formatting should use assets #priority:medium #added:2026-03-07-220819

- [ ] Cleanup internal/cli/system/core/state.go: magic permissions (0o750), magic strings ('Context: ' prefix, etc.) #priority:medium #added:2026-03-07-220819

- [ ] Cleanup internal/cli/system/core/smb.go: errors should come from internal/err; lines 101, 116, 111 need assets text #priority:medium #added:2026-03-07-220819

- [ ] Make AutoPruneStaleDays configurable via ctxrc. Currently hardcoded to 7 days in config.AutoPruneStaleDays; add a ctxrc key (e.g., auto_prune_days) and fallback to the default. #priority:low #added:2026-03-07-220512

- [ ] Refactor check_backup_age/run.go: move consts (lines 23-24) to config, magic directories (line 59) to config, symbolic constants for strings (line 72), messages to assets (lines 79, 90-91), extract non-Run functions to system/core, fix docstrings #priority:medium #added:2026-03-07-180020

- [ ] Add ctxrc support for recall.list.limit to make the default --limit for recall list configurable. Currently hardcoded as config.DefaultRecallListLimit (20). #priority:low #added:2026-03-07-164342

- [ ] Extract journal/core into a standalone journal parser package — functionally isolated enough for its own package rather than remaining as core/ #added:2026-03-07-093815

- [ ] Move PluginInstalled/PluginEnabledGlobally/PluginEnabledLocally from initialize to internal/claude — these are Claude Code plugin detection functions, not init-specific #added:2026-03-07-091656

- [ ] Move guide/cmd/root/run.go text to assets, listCommands to separate file + internal/write #added:2026-03-07-090322

- [ ] Move drift/core/sanitize.go strings to assets #added:2026-03-07-090322

- [ ] Move drift/core/out.go output functions to internal/write per convention #added:2026-03-07-090322

- [ ] Move drift/core/fix.go fmt.Sprintf strings to assets — user-facing output text for i18n #added:2026-03-07-090322

- [ ] Move drift/cmd/root/run.go cmd.Print* output strings to internal/write per convention #added:2026-03-07-084152

- [ ] Extract doctor/core/checks.go strings — 105 inline Name/Category/Message values to assets (i18n) and config (Name/Category constants) #added:2026-03-07-083428

- [ ] Split deps/core builders into per-ecosystem packages — go.go, node.go, python.go, rust.go are specific enough for their own packages under deps/core/ or deps/builders/ #added:2026-03-07-082827

- [ ] Audit git graceful degradation — verify all exec.Command(git) call sites degrade gracefully when git is absent, per project guide recommendation #added:2026-03-07-081625

- [ ] Fix 19 doc.go quality issues: system (13 missing subcmds), agent (phantom refs), load/loop (header typo), claude (stale migration note), 13 minimal descriptions (pause, resume, task, notify, decision, learnings, remind, context, eventlog, index, rc, recall/parser, task/core) #added:2026-03-07-075741

- [ ] Move cmd.Print* output strings in compact/cmd/root/run.go to internal/write per convention #added:2026-03-07-074737

- [ ] Extract changes format.go rendering templates to assets — headings, labels, and format strings are user-facing text for i18n #added:2026-03-07-074719

- [ ] Lift HumanAgo and Pluralize to a common package — reusable time formatting, used by changes and potentially status/recall #added:2026-03-07-074649

- [ ] Extract isAlnum predicate for localization — currently ASCII-only in agent keyword extraction (score.go:141) #added:2026-03-07-073900

- [ ] Make stopwords configurable via .ctxrc — currently embedded in assets, domain users need custom terms #added:2026-03-07-073900

- [ ] Make recency scoring thresholds and relevance match cap configurable via .ctxrc — currently hardcoded in config (7/30/90 days, cap 3) #added:2026-03-07-073900

- [ ] Make DefaultAgentCooldown configurable via .ctxrc — currently hardcoded at 10 minutes in config #added:2026-03-07-073106

- [ ] Make TaskBudgetPct and ConventionBudgetPct configurable via .ctxrc — currently hardcoded at 0.40 and 0.20 in config #added:2026-03-07-072714

- [ ] Localization inventory: audit config constants, write package templates, and assets YAML for i18n mapping — low priority, most users are English-first developers #added:2026-03-06-192419

- [ ] Consider indexing tasks and conventions in TASKS.md and CONVENTIONS.md (currently only decisions and learnings have index tables) #added:2026-03-06-190225


- [ ] Remove FlagNoColor and fatih/color dependency — replace with stdlib terminal coloring or plain output #added:2026-03-06-182831

- [ ] Validate .ctxrc against ctxrc.schema.json at load time — schema is embedded but never enforced, doctor does field-level checks without using it #added:2026-03-06-174851

- [ ] Fix 3 CI compliance issues from PR #27 after merge: missing copyright header on internal/mcp/server_test.go, missing doc.go for internal/cli/mcp/, literal newlines in internal/mcp/resources.go and tools.go #added:2026-03-06-141508

- [ ] Add PostToolUse session event capture. Append lightweight event records (tool name, files touched, timestamp) to .context/state/session-events.jsonl on significant PostToolUse events (file edits, git operations, errors). Not SQLite — just JSONL append. This feeds the PreCompact snapshot hook with richer input so it can report what the agent was actively working on, not just static file state. #added:2026-03-06-185126

- [ ] Add next-step hints to ctx agent and ctx status output. Append actionable suggestions based on context health (e.g. stale tasks, high completion ratio, drift findings). Pattern learned from GitNexus self-guiding agent workflows. #added:2026-03-06-184829

- [ ] Implement PreCompact and SessionStart hooks for session continuity across compaction. Wire ctx agent --budget 4000 to both events: PreCompact outputs context packet before compaction so compactor preserves key info; SessionStart re-injects context packet so fresh/post-compact sessions start oriented. Two thin ctx system subcommands, two entries in hooks.json. See ideas/gitnexus-contextmode-analysis.md for design rationale. #added:2026-03-06-184825

- [ ] Audit fatih/color removal across ~35 files — removed from recall/run.go, recall/lock.go, write/validate.go; ~30 files remain. Separate consolidation pass. #added:2026-03-06-050140

- [ ] Audit remaining 2006-01-02 usages across codebase — 5+ files still use the literal instead of config.DateFormat. Incremental migration. #added:2026-03-06-050140

- [ ] WC.2: Audit CLI packages for direct fmt.Print/Println usage — candidates for migration #added:2026-03-06

### Phase WC2: Write Output Block Consolidation

Spec: `specs/write-output-consolidation.md`. Read the spec before starting any WC2 task.

Consolidate multi-line imperative `cmd.Println` sequences in `internal/write/`
into pre-computed single-print block patterns. Separates conditional logic from
I/O and replaces 4-8 individual Tpl\* constants per function with one block template.

- [ ] WC2.1: Tier 1 — Consolidate multi-line functions with no conditionals: `InfoInitNextSteps`, `InfoObsidianGenerated`, `InfoJournalSiteGenerated`, `InfoDepsNoProject`, `ArchiveDryRun`, `ImportScanHeader`. Add `TplXxxBlock` YAML entries, wire through embed.go + config.go, remove replaced individual constants. #added:2026-03-17
- [ ] WC2.2: Tier 2a — Consolidate conditional functions in info.go: `InfoLoopGenerated` (pre-compute iterLine). Prove the pre-computation pattern on the function that motivated this spec. #added:2026-03-17
- [ ] WC2.3: Tier 2b — Consolidate conditional functions in sync/recall/notify: `SyncResult`, `CtxSyncHeader`, `CtxSyncAction`, `SessionMetadata`, `TestResult`, `SyncDryRun`, `PruneSummary`. Each needs 1-3 pre-computed strings before the single print call. #added:2026-03-17
- [ ] WC2.4: Constant cleanup — verify all replaced individual `TplXxx*` config vars, `TextDescKey*` constants, and YAML entries are removed. Run `make lint` and `go test ./internal/write/...` to confirm no regressions. #added:2026-03-17
- [ ] WC2.5: Update CONVENTIONS.md — add a "Write Package Output" subsection documenting the pre-compute-then-print pattern for future functions with 4+ Printlns and conditionals. #added:2026-03-17

## MCP-related

### Phase MCP-SAN: MCP Server Input Sanitization

[ ] Assignee: @CoderMungan -- https://github.com/ActiveMemory/ctx/issues/49

### Phase MCP-COV: MCP Test Coverage

[ ] Assignee: @CoderMungan -- https://github.com/ActiveMemory/ctx/issues/50

## Later

### Phase PR: State Pruning (`ctx system prune`)

Clean stale per-session state files from `.context/state/`. Files with UUID
session ID suffixes accumulate ~6-8 per session with no cleanup. Strategy:
age-based — prune files older than N days (default 7).

- [ ] Regenerate site/ for state-maintenance recipe (docs/recipes/state-maintenance.md added but site not rebuilt) #added:2026-03-05-205425

- [ ] Audit remaining global tombstones for session-scoping: backup-reminded, ceremony-reminded, check-knowledge, journal-reminded, version-checked, ctx-wrapped-up all have the same cross-session suppression bug as memory-drift-nudged #added:2026-03-05-205425

- [ ] F.2: ctx recall import — import Claude Code session JSONLs from local or remote (~/.claude/projects/) into local ~/.claude/projects/. Pure Go: local copy with os.CopyFS-style walk, remote via os/exec ssh+scp (no rsync dependency). --source flag accepts local path or user@host. --dry-run shows what would be copied. Skips existing files (content-addressed by UUID filenames). Enables journal export from sessions that ran on other machines. #added:2026-03-05-141912

- [ ] P0.5: Blog: "Building a Claude Code Marketplace Plugin" — narrative from session
      history, journals, and git diff of feat/plugin-conversion branch.
      Covers: motivation (shell hooks to Go subcommands), plugin directory
      layout, marketplace.json, eliminating make plugin, bugs found during
      dogfooding (hooks creating partial .context/), and the fix. Use
      /ctx-blog-changelog with branch diff as source material. #added:2026-02-16-111948
- [ ] P9.2: Test manually on this project's LEARNINGS.md (20+ entries).
      #priority:medium #added:2026-02-19
- [ ] P0.8.1: Install golangci-lint on the integration server #for-human
      #priority:medium #added:2026-02-23 #added:2026-02-23-170213
- [ ] PM.3: Review hook diagnostic logs after a long session. Check
      `.context/logs/check-persistence.log` and
       `.context/logs/check-context-size.log` to verify hooks fire correctly.
       Tune nudge frequency if needed. #priority:medium #added:2026-02-09
- [ ] PM.4: Run `/consolidate` to address codebase drift. Considerable drift has
      accumulated (predicate naming, magic strings, hardcoded permissions,
      godoc style). #priority:medium #added:2026-02-06
- [ ] Improve test coverage for core packages at 0% #added:2026-03-20-164324

- [ ] PM.7: Aider/Cursor parser implementations: the recall architecture was
      designed for extensibility (tool-agnostic Session type with
      tool-specific parsers). Adding basic Aider and Cursor parsers would
      validate the parser interface, broaden the user base, and fulfill
      the "works with any AI tool" promise. Aider format is simpler than
      Claude Code's. #priority:medium #source:report-6 #added:2026-02-17

## Future

- [ ] P0.8.5: Enable webhook notifications in worktrees. Currently `ctx notify`
      silently fails because `.context.key` is gitignored and absent in
      worktrees. For autonomous runs with opaque worktree agents, notifications
      are the one feature that would genuinely be useful. Possible approaches:
      resolve the key via `git rev-parse --git-common-dir` to find the main
      checkout, or copy the key into worktrees at creation time (ctx-worktree
      skill). #priority:medium #added:2026-02-22
- [ ] P0.9.2: Split cli-reference.md (1633 lines) into command group pages:
  cli-overview, cli-init-status, cli-context, cli-recall, cli-tools, cli-system —
  each page covers a natural command group with its subcommands and flags
  #added:2026-02-24-204208
- [ ] P0.9.3: Investigate proactive content suggestions: docs/recipes/publishing.md claims
  agents suggest blog posts and journal rebuilds at natural moments, but no hook
  or playbook mechanism exists to trigger this — either wire it up (e.g.
  post-task-completion nudge) or tone down the docs to match reality
  #added:2026-02-24-185754
- [ ] PG.1: Add agent/tool compatibility matrix to prompting guide — document which
      patterns degrade gracefully when agents lack file access, CLI tools, or
      ctx integration. Treat as a "works best with / degrades to" table.
      #priority:medium #added:2026-02-25
- [ ] PG.2: Add versioning/stability note to prompting guide — "these principles are
      stable; examples evolve" + doc date in frontmatter. Needed once the guide
      becomes canonical and people start quoting it. #priority:low #added:2026-02-25
- [ ] P0.1: Brainstorm: Standardize drift-check comment format and integrate with
  `/ctx-drift` — formalize ad-hoc `<!-- drift-check: ... -->` markers, teach
  drift skill to parse/execute them, publish pattern in docs/recipes. Benefits
  tooling/CLI but AI handles ad-hoc fine for now. #priority:medium #added:2026-02-28
- [ ] F.1: MCP server integration: expose context as tools/resources via Model
  Context Protocol. Would enable deep integration with any
  MCP-compatible client. #priority:low #source:report-6
- [ ] Q.1: Docstring cross-reference audit — compliance test that flags docstrings
  mentioning domains that don't match their callers. Start with `write/**`,
  extend to all `internal/`. Spec: `specs/docstring-cross-reference-audit.md`
  #priority:medium #added:2026-03-17

- [ ] Migrate Sprintf-based templates (tpl_*.go) to Go text/template or embedded template files — ObsidianReadme, LoopScript, and other multi-line format strings that can't move to YAML #added:2026-03-18-163629

- [ ] Split internal/assets/embed_test.go — tests that call read/ packages must move to their respective read/ package to avoid import cycles #added:2026-03-18-192914

- [ ] Improve recall/core format tests — replace hardcoded string assertions (e.g. Contains Tokens) with semantic checks that verify structure and values, not label text #added:2026-03-19-194645

### Phase BT: Build Tooling — `cmd/ctxctl`

Replace shell-based build scripts (Makefile shell expansions, `hack/build-all.sh`,
`hack/release.sh`, `hack/tag.sh`, `sync-*`/`check-*` targets) with a first-class
Go binary at `cmd/ctxctl`. Shares internal packages with `ctx` (version, assets,
embed FS). Installable: `go install github.com/ActiveMemory/ctx/cmd/ctxctl@latest`.
Eliminates `jq` build dependency. Testable, cross-platform.

- [ ] Bug: release script versions.md table insertion fails silently. The sed pattern on line 133 uses `$` anchor but the actual Markdown table header has column padding spaces before the trailing `|`. The row is never inserted. Fix: relax the header match pattern or switch to a simpler approach (e.g., insert after the separator line directly). Also verify the "latest stable" sed handles trailing `).\n` correctly. #priority:high #added:2026-03-23-221500

- [ ] Replace hack/lint-drift.sh with AST-based Go tests in internal/audit/. Spec: `specs/ast-audit-tests.md` #added:2026-03-23-210000


Dividing line: `ctx` is the user/agent tool, `ctxctl` is the maintainer/contributor
tool. If a developer clones the repo and needs to build, test, release, or validate
— that's `ctxctl`. If a user is working in a project and needs context — that's `ctx`.

Strong fits beyond build/release:
- `ctxctl plugin package` — package .claude-plugin for marketplace publishing
- `ctxctl plugin validate` — validate plugin.json, hooks.json, skill structure
- `ctxctl doctor` — contributor pre-flight (Go version, tools, GPG, hooks);
  absorbs `hack/gpg-fix.sh` and `hack/gpg-test.sh`
- `ctxctl changelog` — deterministic release notes from git log

Reasonable fits if project grows:
- `ctxctl test smoke` — replaces the shell pipeline in `make smoke`
- `ctxctl site build/serve` — wraps zensical + feed generation
- `ctxctl mcp register` — replaces `hack/gemini-search.sh` and future MCP registrations

Not a fit (keep in `ctx`):
- Anything user-facing in a project context (status, agent, drift, recall)
- Anything Claude Code hooks call — hooks must call `ctx`, not `ctxctl`

- [ ] Design `ctxctl` CLI surface: `ctxctl sync`, `ctxctl build`, `ctxctl release`, `ctxctl check`, `ctxctl tag` #added:2026-03-25-050000
- [ ] Implement `ctxctl sync` — stamps VERSION into plugin.json + syncs why docs; replaces `sync-version`, `sync-why` #added:2026-03-25-050000
- [ ] Implement `ctxctl check` — drift checks: version sync, why docs, lint-drift, lint-docs; replaces `check-*` targets #added:2026-03-25-050000
- [ ] Implement `ctxctl build` — cross-platform builds with version stamping; replaces `build-all.sh` #added:2026-03-25-050000
- [ ] Implement `ctxctl release` — full release flow (sync, build, tag, checksums); replaces `release.sh` + `tag.sh` #added:2026-03-25-050000
- [ ] Simplify Makefile to thin wrappers: `make build` → `go run ./cmd/ctxctl build` #added:2026-03-25-050000
- [ ] Remove `jq` build dependency once ctxctl handles JSON natively #added:2026-03-25-050000

- [ ] Implement MCP warm-up in /ctx-remember session ceremony — when a graph/RAG tool is configured in .ctxrc, run one orientation query at session start to build procedural familiarity. Spec: `ideas/spec-mcp-warm-up-ceremony.md` #added:2026-03-25-120000

- [ ] Update ctx doctor to check for graph tool availability — detect if a graph/RAG MCP is configured in .ctxrc, verify connection status, recommend installation if missing #added:2026-03-25-120000

- [ ] Explore pluggable graph tool interface — replace hardcoded GitNexus references in skill text with configurable .ctxrc graph_tool key. Skills use template placeholder instead of literal tool names. Define minimum interface contract (query, context, impact). Spec: `ideas/spec-mcp-warm-up-ceremony.md` #added:2026-03-25-120000
