# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|--------|
| 2026-03-06 | Create internal/parse for shared text-to-typed-value conversions |
| 2026-03-06 | Centralize errors in internal/err, not per-package err.go files |
| 2026-03-05 | Gitignore .context/memory/ for this project |
| 2026-03-05 | Memory bridge design: three-phase architecture with hook nudge + on-demand |
| 2026-03-05 | Revised strategic analysis: blog-first execution order, bidirectional sync as top-level section |
| 2026-03-04 | Interface-based GraphBuilder for multi-ecosystem ctx deps |
| 2026-03-02 | Billing threshold piggybacks on check-context-size, not heartbeat |
| 2026-03-02 | Replace auto-migration with stderr warning for legacy keys |
| 2026-03-02 | Consolidate all session state to .context/state/ |
| 2026-03-01 | PersistentPreRunE init guard with three-level exemption |
| 2026-03-01 | Global encryption key at ~/.ctx/.ctx.key |
| 2026-03-01 | Heartbeat token telemetry: conditional fields, not always-present |
| 2026-03-01 | Hook log rotation: size-based with one previous generation, matching eventlog pattern |
| 2026-03-01 | Promote 6 private skills to bundled plugin skills; keep 7 project-local |
| 2026-02-27 | Context window detection: JSONL-first fallback order |
| 2026-02-27 | Context injection architecture v2 (consolidated) |
| 2026-02-26 | .context/state/ directory for project-scoped runtime state |
| 2026-02-26 | Hook and notification design (consolidated) |
| 2026-02-26 | ctx init and CLAUDE.md handling (consolidated) |
| 2026-02-26 | Task and knowledge management (consolidated) |
| 2026-02-26 | Agent autonomy and separation of concerns (consolidated) |
| 2026-02-26 | Security and permissions (consolidated) |
| 2026-02-27 | Webhook and notification design (consolidated) |
<!-- INDEX:END -->

## [2026-03-06-050132] Create internal/parse for shared text-to-typed-value conversions

**Status**: Accepted

**Context**: parseDate with 2006-01-02 duplicated in 5+ files; needed a home that is not internal/utils or internal/strings (collides with stdlib)

**Decision**: Create internal/parse for shared text-to-typed-value conversions

**Rationale**: internal/parse scopes to convert text to typed values without becoming a junk drawer. Name invites sibling functions (duration, identifier parsing) naturally.

**Consequences**: parse.Date() is the first function; config.DateFormat holds the layout constant. Other time.Parse callers can migrate incrementally.

---

## [2026-03-06-050131] Centralize errors in internal/err, not per-package err.go files

**Status**: Accepted

**Context**: Duplicate error constructors across 5+ CLI packages; agents copying the pattern when they see a local err.go

**Decision**: Centralize errors in internal/err, not per-package err.go files

**Rationale**: Single location makes duplicates visible, enables future sentinel errors, and prevents broken-window accumulation

**Consequences**: All CLI err.go files migrated and deleted. New errors go to internal/err/errors.go exclusively.

---

## [2026-03-05-205424] Gitignore .context/memory/ for this project

**Status**: Accepted

**Context**: Memory mirror contains copies of MEMORY.md which holds strategic analysis and session notes

**Decision**: Gitignore .context/memory/ for this project

**Rationale**: Strategic content should not be in git history. Docs updated to say 'often git-tracked' for the general recommendation — this project is the exception.

**Consequences**: Mirror and archives are local-only for this project. Other projects can still track them. Sync and drift detection work the same way regardless.

---

## [2026-03-05-042154] Memory bridge design: three-phase architecture with hook nudge + on-demand

**Status**: Accepted

**Context**: Brainstormed how to bridge Claude Code MEMORY.md with ctx structured context files

**Decision**: Memory bridge design: three-phase architecture with hook nudge + on-demand

**Rationale**: Hook nudge + on-demand gives user choice and freedom. Wrap-up is the publish trigger, never commit (footgun). Heuristic classification for v1, no LLM. Marker-based merge for bidirectional conflict. Mirror is git-tracked + timestamped archives. Foundation spec delivers sync/status/diff/hook; import and publish are future phases.

**Consequences**: Foundation spec in specs/memory-bridge.md, import/publish specs deferred to ideas/. Tasked out as S-0.1.1 through S-0.1.10 in ideas/TASKS.md.

---

## [2026-03-05-023937] Revised strategic analysis: blog-first execution order, bidirectional sync as top-level section

**Status**: Accepted

**Context**: Editorial review of ideas/claude-memory-strategic-analysis.md surfaced six structural weaknesses in competitive positioning

**Decision**: Revised strategic analysis: blog-first execution order, bidirectional sync as top-level section

**Rationale**: 200-line cap is fragile differentiator (demoted); org-scoped memory is the real threat (elevated to HIGH); model agnosticism is premature (parked with trigger condition); bidirectional sync is the most underweighted insight (promoted); narrative shapes categories before implementation does (blog first)

**Consequences**: Execution order is now S-3 (blog) -> S-0 -> S-1 -> S-2. Strategic doc restructured from 9 to 10 sections. Blog post shipped as first deliverable.

---

## [2026-03-04-105238] Interface-based GraphBuilder for multi-ecosystem ctx deps

**Status**: Accepted

**Context**: P-1.3 questioned whether non-Go dependency support would introduce bloat and whether a semantic approach was better

**Decision**: Interface-based GraphBuilder for multi-ecosystem ctx deps

**Rationale**: The output pipeline (map[string][]string to Mermaid/table/JSON) was already language-agnostic. Each ecosystem builder is ~40 lines — this is finishing what was started, not bloat. Static manifest parsing (no external tools for Node/Python) keeps dependencies minimal.

**Consequences**: ctx deps now auto-detects Go, Node.js, Python, Rust. --type flag overrides detection. ctx-map skill works across ecosystems without changes.

---

## [2026-03-02-165038] Billing threshold piggybacks on check-context-size, not heartbeat

**Status**: Accepted

**Context**: User wanted a configurable token-count nudge for billing awareness (Claude Pro 1M context, extra cost after 200k). Heartbeat produces zero stdout and can't relay to user.

**Decision**: Billing threshold piggybacks on check-context-size, not heartbeat

**Rationale**: check-context-size already reads tokens, has VERBATIM relay working, and runs every prompt. Adding a third independent trigger there is minimal code and follows established patterns.

**Consequences**: New .ctxrc field billing_token_warn (default 0 = disabled). One-shot per session via billing-warned-{sessionID} state file. Template-overridable via check-context-size/billing.txt.

---

## [2026-03-02-123611] Replace auto-migration with stderr warning for legacy keys

**Status**: Accepted

**Context**: Auto-migration code existed for promoting keys from ~/.local/ctx/keys/ and .context/.ctx.key to ~/.ctx/.ctx.key. Userbase is small and this is alpha — no need to bloat the codebase.

**Decision**: Replace auto-migration with stderr warning for legacy keys

**Rationale**: Warn-only is simpler, avoids silent file operations, and puts the user in control. Migration instructions in docs are sufficient for the small userbase.

**Consequences**: MigrateKeyFile() now only warns on stderr. promoteToGlobal() helper deleted. Tests verify keys are not moved.

---

## [2026-03-02-005213] Consolidate all session state to .context/state/

**Status**: Accepted

**Context**: Session-scoped state (cooldown tombstones, pause markers, daily throttle markers) was split between /tmp (via secureTempDir()) and .context/state/ for project-scoped state

**Decision**: Consolidate all session state to .context/state/

**Rationale**: Single location simplifies mental model, eliminates duplicated secureTempDir() in two packages, removes the cleanup-tmp SessionEnd hook entirely. .context/state/ is already gitignored and project-scoped.

**Consequences**: All 18 callers updated. Tests switch from XDG_RUNTIME_DIR mocking to CTX_DIR + rc.Reset(). Hook lifecycle drops from 4 events to 3 (SessionEnd removed).

---

## [2026-03-01-222733] PersistentPreRunE init guard with three-level exemption

**Status**: Accepted

**Context**: ctx commands handled missing .context/ inconsistently — some caught errors, some got confusing file-not-found messages, some produced empty output

**Decision**: PersistentPreRunE init guard with three-level exemption

**Rationale**: Single PersistentPreRunE on root command gives one clear error. Three-level exemption (hidden commands, annotated commands, grouping commands) covers all edge cases without per-command boilerplate

**Consequences**: Boundary violation now returns an error instead of os.Exit(1), making it testable. The subprocess-based boundary test was simplified to a direct error assertion

---

## [2026-03-01-161457] Global encryption key at ~/.ctx/.ctx.key

**Status**: Superseded by [2026-03-02] global key simplification

**Context**: Key stored next to ciphertext (.context/.ctx.key) was a security antipattern and broke in worktrees. The slug-based per-project key system at ~/.local/ctx/keys/ was over-engineered for the common case (one user, one machine, one key).

**Decision**: Single global key at ~/.ctx/.ctx.key. Project-local override via .ctxrc key_path or .context/.ctx.key.

**Rationale**: One key per machine covers 99% of users. Per-project slug filenames and three-tier resolution added complexity without clear benefit. ~/.ctx/ is the natural home (matches ~/.claude/ convention). Tilde expansion in .ctxrc key_path fixes a standalone bug.

**Consequences**: Auto-migration promotes legacy keys (project-local, ~/.local/ctx/keys/) to ~/.ctx/.ctx.key. Deleted KeyDir(), ProjectKeySlug(), ProjectKeyPath(). ResolveKeyPath simplified to two params. 15+ doc files updated.

---

## [2026-03-01-112544] Heartbeat token telemetry: conditional fields, not always-present

**Status**: Accepted

**Context**: Adding tokens, context_window, usage_pct to heartbeat payloads. First prompt of a session has no JSONL usage data yet.

**Decision**: Heartbeat token telemetry: conditional fields, not always-present

**Rationale**: Token fields are only included in the template ref when tokens > 0. This avoids misleading pct=0% on the first heartbeat and keeps payloads clean for receivers that filter on field presence.

**Consequences**: Webhook consumers must handle heartbeats both with and without token fields. The message string also varies (with/without tokens=N pct=N% suffix).

---

## [2026-03-01-092613] Hook log rotation: size-based with one previous generation, matching eventlog pattern

**Status**: Accepted

**Context**: .context/logs/ files grow unbounded (~200KB after one month); needed a cap

**Decision**: Hook log rotation: size-based with one previous generation, matching eventlog pattern

**Rationale**: Architectural symmetry with eventlog, O(1) size check vs O(n) line counting, diagnostic logs don't need deep history (webhooks cover serious setups)

**Consequences**: Each log file caps at ~2MB (current + .1). config.LogMaxBytes = 1MB, same as EventLogMaxBytes

---

## [2026-03-01-090124] Promote 6 private skills to bundled plugin skills; keep 7 project-local

**Status**: Accepted

**Context**: Reviewed all 13 _ctx-* private skills to determine which are universally useful for any ctx user vs specific to the ctx codebase or personal infra.

**Decision**: Promote 6 private skills to bundled plugin skills; keep 7 project-local

**Rationale**: Promote if the skill benefits any ctx-powered project without project-specific hardcoding. Keep private if it references this repo's Go internals, personal infra, or language-specific tooling. Promote list: _ctx-spec (generic scaffolding), _ctx-brainstorm (design facilitation), _ctx-verify (claim verification), _ctx-skill-creator (skill authoring), _ctx-check-links (doc link audit), _ctx-sanitize-permissions (Claude Code permissions audit). Keep list: _ctx-audit (Go/ctx checks), _ctx-qa (Go Makefile), _ctx-backup (SMB infra), _ctx-release/_ctx-release-notes (ctx release workflow), _ctx-update-docs (ctx package mapping), _ctx-absorb (borderline, revisit later).

**Consequences**: Six skills move from .claude/skills/ to internal/assets/claude/skills/ and become available to all ctx users via ctx init. Cross-references between skills need updating (e.g., /_ctx-brainstorm becomes /ctx-brainstorm). The seven remaining private skills stay project-local.

---

## [2026-02-27-230718] Context window detection: JSONL-first fallback order

**Status**: Accepted

**Context**: check-context-size defaults to 200k but user runs 1M-context model, causing false 110% warnings. JSONL contains the model name which maps to actual window size.

**Decision**: Context window detection: JSONL-first fallback order

**Rationale**: effective_window = detect_from_jsonl(model) ?? ctxrc.context_window ?? 200_000. JSONL is ground truth (reflects actual model in use); ctxrc is fallback for first-hook-of-session or unknown models; 200k is safe last resort. Having ctxrc override JSONL would artificially restrict the check when a user forgets to update their config after switching models.

**Consequences**: Most users get correct window automatically. ctxrc context_window becomes a fallback, not an override. Task exists for implementation.

---

## [2026-02-27-002830] Context injection architecture v2 (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-02-26)

- **Diagram extraction**: ARCHITECTURE.md contained ~600 lines of ASCII/Mermaid diagrams (~12K tokens). Extracted to 5 architecture-dia-*.md files outside FileReadOrder. Agents get verbal summaries at session start; diagrams available on demand. Total injection dropped 53% (20K→9.5K tokens).
- **Auto-injection replaces directives**: Soft instructions have ~75-85% compliance ceiling because "don't apply judgment" is itself evaluated by judgment. The v2 context-load-gate injects content directly via `additionalContext` — agents never choose whether to comply. Injection strategy: CONSTITUTION, CONVENTIONS, ARCHITECTURE, AGENT_PLAYBOOK verbatim; DECISIONS, LEARNINGS index-only; TASKS mention-only. Total ~7,700 tokens. See: `specs/context-load-gate-v2.md`.
- **Imperative framing**: Advisory framing allowed agents to assess relevance and skip files. Imperative framing with unconditional compliance checkpoint removes the escape hatch. Verbatim relay is fallback safety net, not primary instruction.

---

## [2026-02-26-200001] .context/state/ directory for project-scoped runtime state

**Status**: Accepted

New gitignored directory under `context_dir` resolution for ephemeral project-scoped state. Follows `.context/logs/` precedent — added to `config.GitignoreEntries` and root `.gitignore`.

First use: injection oversize flag written by context-load-gate when injected tokens exceed the configurable `injection_token_warn` threshold (`.ctxrc`, default 15000). The check-context-size VERBATIM hook reads the flag and nudges the user to run `/ctx-consolidate`.

See: `specs/injection-oversize-nudge.md`.

---

## [2026-02-26-100001] Hook and notification design (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-02-12 to 2026-02-24)

- Tone down proactive content suggestion claims in docs rather than add more hooks. Already have 9 UserPromptSubmit hooks; adding another risks fatigue. Conversational prompting already works.
- Hook commands must use structured JSON output (hookSpecificOutput.additionalContext) instead of plain text, because Claude Code treats plain text as ignorable ambient context.
- Drop prompt-coach hook entirely: zero useful tips fired, output channel invisible to user, orphan temp file accumulation. The prompting guide already covers best practices.
- De-emphasize /ctx-journal-normalize from the default journal pipeline. The normalize skill is expensive and nondeterministic; programmatic normalization handles most cases. Skill remains available for targeted per-file use.

---

## [2026-02-26-100002] ctx init and CLAUDE.md handling (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-01-20)

- `ctx init` handles CLAUDE.md intelligently: creates if missing, backs up and offers merge if existing, uses marker comment for idempotency. The `--merge` flag enables non-interactive append.
- `ctx init` always generates `.claude/hooks/` alongside `.context/` with no flag needed. Other AI tools ignore `.claude/`; Claude Code users get seamless zero-config experience.
- Core tool stays generic and tool-agnostic, with optional Claude Code enhancements via `.claude/hooks/`. Other AI tools can be supported similarly (`ctx hook cursor`, etc.).

---

## [2026-02-26-100004] Task and knowledge management (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-27 to 2026-02-18)

- Tasks must include explicit deliverables, not just implementation steps. Parent tasks define WHAT the user gets; subtasks decompose HOW to build it. Without explicit deliverables, AI optimizes for checking boxes.
- Use reverse-chronological order (newest first) for DECISIONS.md and LEARNINGS.md. Ensures most recent items are read first regardless of token budget.
- Add quick reference index to DECISIONS.md: compact table at top allows scanning; agents can grep for full timestamp to jump to entry. Auto-updated on `ctx add decision`.
- Knowledge scaling via archive path for decisions and learnings: follow the task archive pattern, move old entries to `.context/archive/`, extend `ctx compact --archive` to cover all three file types.

---

## [2026-02-26-100005] Agent autonomy and separation of concerns (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-01-21 to 2026-01-28)

- Removed AGENTS.md from project root. Consolidated on CLAUDE.md (auto-loaded) + .context/AGENT_PLAYBOOK.md as the canonical agent instruction path. Projects using ctx should not create AGENTS.md.
- Separate orchestrator directive from agent tasks: `.context/TASKS.md` is the agent's mind (tasks the agent owns); `IMPLEMENTATION_PLAN.md` is the orchestrator's thin directive layer ("check your tasks"). Prevents task list drift.
- No custom UI -- IDE is the interface. UI is a liability; IDEs already excel at file browsing, search, markdown editing, and git integration. Focus CLI efforts on good markdown output.

---

## [2026-02-26-100006] Security and permissions (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-21 to 2026-02-24)

- Keep CONSTITUTION.md minimal: only truly inviolable rules (security, correctness, process invariants). Style preferences go in CONVENTIONS.md. Overly strict constitution gets ignored.
- Centralize constants with semantic prefixes in `internal/config/config.go`: `Dir*` for directories, `File*` for paths, `Filename*` for names, `UpdateType*` for entry types. Single source of truth, compile-time typo checks.
- Hooks use `ctx` from PATH, not hardcoded absolute paths. Standard Unix practice; portable across machines/users. `ctx init` checks PATH availability before proceeding.
- Drop absolute-path-to-ctx regex from block-dangerous-commands shell script. The block-non-path-ctx Go subcommand already covers this with better patterns; duplicating creates two sources of truth.

---

## [2026-02-27-002831] Webhook and notification design (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-02-22 to 2026-02-26)

- **Session attribution**: All webhook payloads must include session_id. Reading it from stdin costs nothing and enables multi-agent diagnostics. All run functions take stdin parameter; tests use createTempStdin.
- **Opt-in events**: Notify events are opt-in, not opt-out. EventAllowed returns false for nil/empty event lists. The correct default for notifications is silence. `ctx notify test` bypasses the filter as a special case.
- **Shared encryption key**: Webhook URLs encrypted with the shared .ctx.key (AES-256-GCM), not a dedicated key. One key, one gitignore entry, one rotation cycle. Notify is a peer of scratchpad — both store user secrets encrypted at rest.

---

## [2026-02-11] Remove .context/sessions/ storage layer and ctx session command

**Status**: Accepted

**Context**: The session/recall/journal system had three overlapping storage layers: `~/.claude/projects/` (raw JSONL transcripts, owned by Claude Code), `.context/sessions/` (JSONL copies + context snapshots), and `.context/journal/` (enriched markdown from `ctx recall export`). The recall pipeline reads directly from `~/.claude/projects/`, making `.context/sessions/` a dead-end write sink that nothing reads from. The auto-save hook copied transcripts to a directory nobody consumed. The `ctx session save` command created context snapshots that git already provides through version history. This was ~15 Go source files, a shell hook, ~20 config constants, and 30+ doc references supporting infrastructure with no consumers.

**Decision**: Remove `.context/sessions/` entirely. Two stores remain: raw transcripts (global, tool-owned in `~/.claude/projects/`) and enriched journal (project-local in `.context/journal/`).

**Rationale**: Dead-end write sinks waste code surface, maintenance effort, and user attention. The recall pipeline already proved that reading directly from `~/.claude/projects/` is sufficient. Context snapshots are redundant with git history. Removing the middle layer simplifies the architecture from three stores to two, eliminates an entire CLI command tree (`ctx session`), and removes a shell hook that fired on every session end.

**Consequences**: Deleted `internal/cli/session/` (15 files), removed auto-save hook, removed `--auto-save` from watch, removed pre-compact auto-save from compact, removed `/ctx-save` skill, updated ~45 documentation files. Four earlier decisions superseded (SessionEnd hook, Auto-Save Before Compact, Session Filename Format, Two-Tier Persistence Model). Users who want session history use `ctx recall list/export` instead.

---

*Module-specific, already-shipped, and historical decisions:
[decisions-reference.md](decisions-reference.md)*
