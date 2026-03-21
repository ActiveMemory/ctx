# Decisions

<!-- INDEX:START -->
| Date | Decision |
|------|--------|
| 2026-03-21 | Output functions belong in write/, logic and types in core/ |
| 2026-03-20 | Shared formatting utilities belong in internal/format |
| 2026-03-20 | Go-YAML linkage check added to lint-drift as check 5 |
| 2026-03-18 | Eager Init() for static embedded data instead of per-accessor sync.Once |
| 2026-03-18 | Singular command names for all CLI entities |
| 2026-03-17 | Pre-compute-then-print for write package output blocks |
| 2026-03-16 | Server methods only handle dispatch and I/O, not struct construction |
| 2026-03-16 | Explicit Init over package-level init() for resource lookup |
| 2026-03-16 | Resource name constants in config/mcp/resource, mapping in server/resource |
| 2026-03-16 | Rename --consequences flag to --consequence for singular consistency |
| 2026-03-15 | Pure-logic CompactContext with no I/O — callers own file writes and reporting |
| 2026-03-15 | TextDescKey exhaustive test verifies all 879 constants resolve to non-empty YAML values |
| 2026-03-15 | Split text.yaml into 6 domain files loaded via loadYAMLDir |
| 2026-03-14 | Error package taxonomy: 22 domain files replace monolithic errors.go |
| 2026-03-14 | Session prefixes are parser vocabulary, not i18n text |
| 2026-03-14 | System path deny-list as safety net, not security boundary |
| 2026-03-14 | Config-driven freshness check with per-file review URLs |
| 2026-03-13 | Delete ctx-context-monitor skill — hook output is self-sufficient |
| 2026-03-13 | build target depends on sync-why to prevent embedded doc drift |
| 2026-03-13 | Templates and user-facing text live in assets, structural constants stay in config |
| 2026-03-12 | Recommend companion RAGs as peer MCP servers not bridge through ctx |
| 2026-03-12 | Split commands.yaml into 4 domain files |
| 2026-03-12 | Rename ctx-map skill to ctx-architecture |
| 2026-03-07 | Use composite directory path constants for multi-segment paths |
| 2026-03-06 | Drop fatih/color dependency — Unicode symbols are sufficient for terminal output, color was redundant |
| 2026-03-06 | Externalize all command descriptions to embedded YAML for i18n readiness — commands.yaml holds Short/Long for 105 commands plus flag descriptions, loaded via assets.CommandDesc() and assets.FlagDesc() |
| 2026-03-06 | cmd/root + core taxonomy for all CLI packages — single-command packages use cmd/root/{cmd.go,run.go}, multi-subcommand packages use cmd/<sub>/{cmd.go,run.go}, shared helpers in core/ |
| 2026-03-06 | Shared entry types and API live in internal/entry, not in CLI packages — domain types that multiple packages consume (mcp, watch, memory) belong in a domain package, not a CLI subpackage |
| 2026-03-06 | PR #27 (MCP server) meets v0.1 spec requirements — merge-ready pending 3 compliance fixes |
| 2026-03-06 | Skills stay CLI-based; MCP Prompts are the protocol equivalent |
| 2026-03-06 | Peer MCP model for external tool integration |
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

## [2026-03-21-084020] Output functions belong in write/, logic and types in core/

**Status**: Accepted

**Context**: PrintFeedReport was initially placed in cli/site/core/ but it calls cmd.Println — that's output formatting, not business logic

**Decision**: Output functions belong in write/, logic and types in core/

**Rationale**: The project taxonomy separates concerns: core/ owns domain logic, types, and helpers; write/ owns CLI output formatting that takes *cobra.Command for Println. Mixing them blurs the boundary and makes testing harder.

**Consequence**: All functions that call cmd.Print/Println/Printf belong in the write/ package tree. core/ never imports cobra for output purposes.

---

## [2026-03-20-232506] Shared formatting utilities belong in internal/format

**Status**: Accepted

**Context**: Pluralize, Duration, DurationAgo, and TruncateFirstLine were duplicated across memory/core, change/core, and other CLI packages

**Decision**: Shared formatting utilities belong in internal/format

**Rationale**: internal/format already existed with TimeAgo and Number formatters. Centralizing prevents duplication and matches the convention that domain-agnostic utilities live in shared packages, not CLI subpackages

**Consequence**: CLI packages import internal/format instead of defining local helpers. Local copies deleted.

---

## [2026-03-20-160103] Go-YAML linkage check added to lint-drift as check 5

**Status**: Accepted

**Context**: Prior refactoring sessions left broken and orphan linkages between Go DescKey constants and YAML entries that caused silent runtime failures

**Decision**: Go-YAML linkage check added to lint-drift as check 5

**Rationale**: Shell-based grep+comm approach fits the existing lint-drift pattern, runs at CI time, and is simpler than programmatic Go AST parsing

**Consequence**: CI-time check catches orphans in both directions plus cross-namespace duplicates, preventing recurrence

---

## [2026-03-18-193631] Eager Init() for static embedded data instead of per-accessor sync.Once

**Status**: Accepted

**Context**: 4 sync.Once guards + 4 exported maps + 4 Load functions + a wrapper package for YAML that never mutates.

**Decision**: Eager Init() for static embedded data instead of per-accessor sync.Once

**Rationale**: Data is static and required at startup. sync.Once per accessor is cargo cult. One Init() in main.go is sufficient. Tests call Init() in TestMain.

**Consequence**: Maps unexported, accessors are plain lookups, permissions and stopwords also loaded eagerly. Zero sync.Once remains in the lookup pipeline.

---

## [2026-03-18-193623] Singular command names for all CLI entities

**Status**: Accepted

**Context**: ctx add used learning (singular) but ctx learnings was plural. Inconsistency across 6 commands.

**Decision**: Singular command names for all CLI entities

**Rationale**: Less headache for i18n; one rule (singular = entity); developers think in OOP. Use field values come from DescKey constants for single-source-of-truth renaming.

**Consequence**: All commands singular: task, decision, learning, change, permission, dep. YAML keys, desc constants, directory names, and 50+ files updated.

---

## [2026-03-17-105627] Pre-compute-then-print for write package output blocks

**Status**: Accepted

**Context**: Audit of internal/write/ found 337 Println calls across 160 functions. Asked whether text/template or single format strings would clean up multi-Println functions like InfoLoopGenerated.

**Decision**: Pre-compute-then-print for write package output blocks

**Rationale**: text/template trades compile-time safety for runtime errors and only 38 of 160 functions benefit from consolidation. fmt.Sprintf with pre-computed conditional args handles all cases without new dependencies. Loop-based functions stay imperative.

**Consequence**: Functions with 4+ Printlns pre-compute conditionals into strings, then emit one cmd.Println with a multiline block template. Per-line Tpl* constants replaced with TplXxxBlock. Trivial (1-3 line) and loop-based functions excluded.

---

## [2026-03-16-122033] Server methods only handle dispatch and I/O, not struct construction

**Status**: Accepted

**Context**: MCP server had ok/error/writeError as methods plus prompt builders that didn't use Server state — they just constructed response structs

**Decision**: Server methods only handle dispatch and I/O, not struct construction

**Rationale**: Methods that don't access receiver state hide their true dependencies and inflate the Server interface. Free functions make the dependency graph explicit and are independently testable.

**Consequence**: New response helpers go in server/out, prompt builders in server/prompt. Server methods are limited to dispatch (handlePromptsGet) and I/O (writeJSON, emitNotification). Same principle applies to future tool/resource builders.

---

## [2026-03-16-104143] Explicit Init over package-level init() for resource lookup

**Status**: Accepted

**Context**: server/resource package used init() to silently build the URI lookup map

**Decision**: Explicit Init over package-level init() for resource lookup

**Rationale**: Implicit init hides startup dependencies, makes ordering unclear, and is harder to test. Explicit Init called from NewServer makes the dependency visible.

**Consequence**: res.Init() called explicitly from NewServer before ToList(); no package-level side effects

---

## [2026-03-16-104142] Resource name constants in config/mcp/resource, mapping in server/resource

**Status**: Accepted

**Context**: MCP resource handler had string literals scattered through handle_resource.go and rebuilt the resource list on every call

**Decision**: Resource name constants in config/mcp/resource, mapping in server/resource

**Rationale**: Constants follow the same pattern as config/mcp/tool. Mapping stays in server/resource because it bridges config constants with assets text (too many cross-cutting deps for a config package). Resource list and URI lookup are pre-built once at server init.

**Consequence**: URI-to-file lookup is O(1) via pre-built map; resource list built once in NewServer, not per request; no string literals in handler code

---

## [2026-03-16-022635] Rename --consequences flag to --consequence for singular consistency

**Status**: Accepted

**Context**: All other CLI flags (context, rationale, lesson, application) are singular nouns. consequences was the only plural.

**Decision**: Rename --consequences flag to --consequence for singular consistency

**Rationale**: Singular form matches the pattern. Consistency wins over natural language preference.

**Consequence**: 75+ files updated. Breaking change for --consequences users.

---

## [2026-03-15-230640] Pure-logic CompactContext with no I/O — callers own file writes and reporting

**Status**: Accepted

**Context**: MCP server and CLI compact command both implemented task compaction independently, with the MCP handler using a local WriteContextFile wrapper

**Decision**: Pure-logic CompactContext with no I/O — callers own file writes and reporting

**Rationale**: Separating pure logic from I/O lets both MCP (JSON-RPC responses) and CLI (cobra cmd.Println) callers control output and file writes. Eliminates duplication and the unnecessary mcp/server/fs package

**Consequence**: tidy.CompactContext returns a CompactResult struct; callers iterate FileUpdates and write them. Archive logic stays in callers since MCP and CLI have different archive policies

---

## [2026-03-15-101336] TextDescKey exhaustive test verifies all 879 constants resolve to non-empty YAML values

**Status**: Accepted

**Context**: PR #42 merged with ~80 new MCP text keys but no test coverage for key-to-YAML mapping

**Decision**: TextDescKey exhaustive test verifies all 879 constants resolve to non-empty YAML values

**Rationale**: A single table-driven test parsing embed.go source catches typos and missing YAML entries at test time — no manual key list to maintain

**Consequence**: New TextDescKey constants are automatically covered; orphaned keys fail CI

---

## [2026-03-15-040638] Split text.yaml into 6 domain files loaded via loadYAMLDir

**Status**: Accepted

**Context**: text.yaml grew to 1812 lines covering write, errors, mcp, doctor, hooks, and ui domains

**Decision**: Split text.yaml into 6 domain files loaded via loadYAMLDir

**Rationale**: Matches existing split pattern (commands.yaml, flags.yaml, examples.yaml); loadYAMLDir merges all files in commands/text/ transparently so TextDesc() API stays unchanged

**Consequence**: New domain files must go into commands/text/; loadYAMLDir reads all .yaml in the directory at init time

---

## [2026-03-14-180905] Error package taxonomy: 22 domain files replace monolithic errors.go

**Status**: Accepted

**Context**: internal/err/errors.go was 1995 lines with 188 functions in one file

**Decision**: Error package taxonomy: 22 domain files replace monolithic errors.go

**Rationale**: Convention requires files named by responsibility, not junk drawers; domain grouping makes it possible to find error constructors by domain

**Consequence**: 22 files (backup, config, crypto, date, fs, git, hook, init, journal, memory, notify, pad, parser, prompt, recall, reminder, session, site, skill, state, task, validation); errors.go deleted

---

## [2026-03-14-131152] Session prefixes are parser vocabulary, not i18n text

**Status**: Accepted

**Context**: Markdown session parser had hardcoded Session:/Oturum: pair in text.yaml as session_prefix/session_prefix_alt — didn't scale beyond two languages

**Decision**: Session prefixes are parser vocabulary, not i18n text

**Rationale**: Session header prefixes are recognition patterns for parsing, not user-facing interface strings. Separating content recognition from interface language lets users parse multilingual session files without code changes. Single-language default (Session:) avoids implicit favoritism.

**Consequence**: Prefixes moved to .ctxrc session_prefixes list. text.yaml entries and embed.go constants removed. Parser reads from rc.SessionPrefixes() with fallback to config/parser.DefaultSessionPrefixes. Users extend via .ctxrc.

---

## [2026-03-14-110748] System path deny-list as safety net, not security boundary

**Status**: Accepted

**Context**: Replacing nolint:gosec directives with centralized I/O wrappers in internal/io

**Decision**: System path deny-list as safety net, not security boundary

**Rationale**: ctx paths are internally constructed from config constants. The deny-list catches agent hallucinations (writing to /etc), not adversarial input. Public security docs would imply a threat model that does not exist.

**Consequence**: internal/io/doc.go documents limitations honestly for contributors. No user-facing security docs. The deny-list is a modicum of protection, not a promise.

---

## [2026-03-14-093748] Config-driven freshness check with per-file review URLs

**Status**: Accepted

**Context**: Building a hook to warn when technology-dependent constants go stale. Initially hardcoded the file list and Anthropic docs URL in the binary, but this only worked inside the ctx repo and assumed all projects care about Anthropic docs.

**Decision**: Config-driven freshness check with per-file review URLs

**Rationale**: Making the file list and review URLs configurable via .ctxrc freshness_files means any project can opt in. Per-file review_url avoids special-casing by project name — ctx sets Anthropic docs, other projects set their own vendor links or omit it entirely.

**Consequence**: The hook is a no-op by default (opt-in). ctx's own .ctxrc carries the tracked files. All nudge text goes through assets/text.yaml for localization. No project detection logic needed.

---

## [2026-03-13-223111] Delete ctx-context-monitor skill — hook output is self-sufficient

**Status**: Accepted

**Context**: The skill documented how to relay context window warnings, but the hook message already includes IMPORTANT: Relay this context window warning to the user VERBATIM which agents follow without the skill.

**Decision**: Delete ctx-context-monitor skill — hook output is self-sufficient

**Rationale**: No mechanism exists for hooks to trigger skills. The skill was never loaded during sessions. Adding enforcement elsewhere would either be too far back in context (playbook) or dilute the already-crisp hook message.

**Consequence**: One fewer skill to maintain. No behavioral change — agents continue relaying warnings as before.

---

## [2026-03-13-151955] build target depends on sync-why to prevent embedded doc drift

**Status**: Accepted

**Context**: assets/why/ files had silently drifted from their docs/ sources

**Decision**: build target depends on sync-why to prevent embedded doc drift

**Rationale**: Derived assets that are not in the build dependency chain will drift — the only reliable enforcement is making the build fail without sync

**Consequence**: Every make build now copies docs into assets before compiling

---

## [2026-03-13-151954] Templates and user-facing text live in assets, structural constants stay in config

**Status**: Accepted

**Context**: Ongoing refactoring session moving Tpl* constants out of config/

**Decision**: Templates and user-facing text live in assets, structural constants stay in config

**Rationale**: config/ is for structural constants (paths, limits, regexes); assets/ is for templates, labels, and text that would need i18n. Clean separation of concerns

**Consequence**: All tpl_entry.go, tpl_journal.go, tpl_loop.go, tpl_recall.go moved to assets/

---

## [2026-03-12-133007] Recommend companion RAGs as peer MCP servers not bridge through ctx

**Status**: Accepted

**Context**: Explored whether ctx should proxy RAG queries or integrate a RAG directly

**Decision**: Recommend companion RAGs as peer MCP servers not bridge through ctx

**Rationale**: MCP is the composition layer — agents already compose multiple servers. ctx is context, RAGs are intelligence. No bridging, no plugin system, no schema abstraction

**Consequence**: Spec created at ideas/spec-companion-intelligence.md; future work is documentation and UX only

---

## [2026-03-12-133007] Split commands.yaml into 4 domain files

**Status**: Accepted

**Context**: Single 2373-line YAML mixed commands, flags, text, and examples with inconsistent quoting

**Decision**: Split commands.yaml into 4 domain files

**Rationale**: Context is for humans — localization files should be human-readable block scalars. Separate files eliminate the underscore prefix namespace hack

**Consequence**: 4 files (commands.yaml, flags.yaml, text.yaml, examples.yaml) with dedicated loaders in embed.go

---

## [2026-03-12-133007] Rename ctx-map skill to ctx-architecture

**Status**: Accepted

**Context**: The name 'map' didn't convey the iterative, architectural nature of the ritual

**Decision**: Rename ctx-map skill to ctx-architecture

**Rationale**: 'architecture' better describes surveying and evolving project structure across sessions

**Consequence**: All cross-references updated across skills, docs, .context files, and settings

---

## [2026-03-07-221155] Use composite directory path constants for multi-segment paths

**Status**: Accepted

**Context**: Needed a constant for hooks/messages path used in message.go and message_cmd.go

**Decision**: Use composite directory path constants for multi-segment paths

**Rationale**: Matches existing pattern of DirClaudeHooks = '.claude/hooks' — keeps filepath.Join calls cleaner and avoids scattering path segments

**Consequence**: New multi-segment directory paths should be single constants (e.g. DirHooksMessages, DirMemoryArchive) rather than joined from individual segment constants

---

## [2026-03-06-200306] Drop fatih/color dependency — Unicode symbols are sufficient for terminal output, color was redundant

**Status**: Accepted

**Context**: fatih/color was used in 32 files for green checkmarks, yellow warnings, cyan headings, dim text

**Decision**: Drop fatih/color dependency — Unicode symbols are sufficient for terminal output, color was redundant

**Rationale**: Every colored output already had a semantic symbol (✓, ⚠, ○) that conveyed the same meaning; color added visual noise in non-terminal contexts (logs, pipes)

**Consequence**: Removed --no-color flag (only existed for color.NoColor); one fewer external dependency; FlagNoColor retained in config for CLI compatibility

---

## [2026-03-06-200257] Externalize all command descriptions to embedded YAML for i18n readiness — commands.yaml holds Short/Long for 105 commands plus flag descriptions, loaded via assets.CommandDesc() and assets.FlagDesc()

**Status**: Accepted

**Context**: Command descriptions were inline strings scattered across 105 cobra.Command definitions

**Decision**: Externalize all command descriptions to embedded YAML for i18n readiness — commands.yaml holds Short/Long for 105 commands plus flag descriptions, loaded via assets.CommandDesc() and assets.FlagDesc()

**Rationale**: Centralizing user-facing text in a single translatable file prepares for i18n without runtime cost (embedded at compile time)

**Consequence**: System's 30 hidden hook subcommands excluded (not user-facing); flag descriptions use _flags.scope.name convention

---

## [2026-03-06-200247] cmd/root + core taxonomy for all CLI packages — single-command packages use cmd/root/{cmd.go,run.go}, multi-subcommand packages use cmd/<sub>/{cmd.go,run.go}, shared helpers in core/

**Status**: Accepted

**Context**: 35 CLI packages had inconsistent flat structures mixing Cmd(), run logic, helpers, and types in the same directory

**Decision**: cmd/root + core taxonomy for all CLI packages — single-command packages use cmd/root/{cmd.go,run.go}, multi-subcommand packages use cmd/<sub>/{cmd.go,run.go}, shared helpers in core/

**Rationale**: Taxonomical symmetry: every package has the same predictable shape, making navigation instant and agent-friendly

**Consequence**: cmd/ contains only cmd.go + run.go; helpers go to core/; 474 files changed in initial restructuring

---

## [2026-03-06-200227] Shared entry types and API live in internal/entry, not in CLI packages — domain types that multiple packages consume (mcp, watch, memory) belong in a domain package, not a CLI subpackage

**Status**: Accepted

**Context**: External consumers were importing cli/add for EntryParams/ValidateEntry/WriteEntry, creating a leaky abstraction

**Decision**: Shared entry types and API live in internal/entry, not in CLI packages — domain types that multiple packages consume (mcp, watch, memory) belong in a domain package, not a CLI subpackage

**Rationale**: Domain types in CLI packages force consumers to depend on CLI internals; internal/entry provides a clean boundary

**Consequence**: entry aliases Params from add/core to avoid import cycle (entry imports add/core for insert logic); future work may move insert logic to entry to eliminate the cycle

---

## [2026-03-06-141507] PR #27 (MCP server) meets v0.1 spec requirements — merge-ready pending 3 compliance fixes

**Status**: Accepted

**Context**: Reviewed PR against specs/mcp-server.md; all 7 action items addressed, CI fails on 3 mechanical compliance issues

**Decision**: PR #27 (MCP server) meets v0.1 spec requirements — merge-ready pending 3 compliance fixes

**Rationale**: All spec requirements met; CI failures are trivial and low-risk; keeping PR open risks merge conflicts during active refactoring

**Consequence**: Merge and fix compliance issues in follow-up commit on main

---

## [2026-03-06-184816] Skills stay CLI-based; MCP Prompts are the protocol equivalent

**Status**: Accepted

**Context**: Question arose whether skills should switch from ctx CLI (Bash) to MCP tool calls once the MCP server ships

**Decision**: Skills stay CLI-based; MCP Prompts are the protocol equivalent

**Rationale**: CLI is always available (PATH prerequisite); MCP requires optional configuration. Hooks will always be CLI (shell commands). Two access patterns in the same tool is gratuitous complexity.

**Consequence**: Skills call CLI. MCP Prompts call MCP Tools. Hooks call CLI. Clean layer separation; no replacement, only parallel access paths.

---

## [2026-03-06-184812] Peer MCP model for external tool integration

**Status**: Accepted

**Context**: Evaluated three integration models (orchestrator, peer, hub) for how ctx relates to GitNexus and context-mode

**Decision**: Peer MCP model for external tool integration

**Rationale**: Peer model (side-by-side MCP servers, each queried independently by the agent) respects ctx's markdown-on-filesystem invariant and avoids coupling. ctx provides behavioral scaffolding; external tools provide their specialties.

**Consequence**: ctx MCP Prompts can reference external tools by convention without tight coupling. No plugin registry needed.

---

## [2026-03-06-050132] Create internal/parse for shared text-to-typed-value conversions

**Status**: Accepted

**Context**: parseDate with 2006-01-02 duplicated in 5+ files; needed a home that is not internal/utils or internal/strings (collides with stdlib)

**Decision**: Create internal/parse for shared text-to-typed-value conversions

**Rationale**: internal/parse scopes to convert text to typed values without becoming a junk drawer. Name invites sibling functions (duration, identifier parsing) naturally.

**Consequence**: parse.Date() is the first function; config.DateFormat holds the layout constant. Other time.Parse callers can migrate incrementally.

---

## [2026-03-06-050131] Centralize errors in internal/err, not per-package err.go files

**Status**: Accepted

**Context**: Duplicate error constructors across 5+ CLI packages; agents copying the pattern when they see a local err.go

**Decision**: Centralize errors in internal/err, not per-package err.go files

**Rationale**: Single location makes duplicates visible, enables future sentinel errors, and prevents broken-window accumulation

**Consequence**: All CLI err.go files migrated and deleted. New errors go to internal/err/errors.go exclusively.

---

## [2026-03-05-205424] Gitignore .context/memory/ for this project

**Status**: Accepted

**Context**: Memory mirror contains copies of MEMORY.md which holds strategic analysis and session notes

**Decision**: Gitignore .context/memory/ for this project

**Rationale**: Strategic content should not be in git history. Docs updated to say 'often git-tracked' for the general recommendation — this project is the exception.

**Consequence**: Mirror and archives are local-only for this project. Other projects can still track them. Sync and drift detection work the same way regardless.

---

## [2026-03-05-042154] Memory bridge design: three-phase architecture with hook nudge + on-demand

**Status**: Accepted

**Context**: Brainstormed how to bridge Claude Code MEMORY.md with ctx structured context files

**Decision**: Memory bridge design: three-phase architecture with hook nudge + on-demand

**Rationale**: Hook nudge + on-demand gives user choice and freedom. Wrap-up is the publish trigger, never commit (footgun). Heuristic classification for v1, no LLM. Marker-based merge for bidirectional conflict. Mirror is git-tracked + timestamped archives. Foundation spec delivers sync/status/diff/hook; import and publish are future phases.

**Consequence**: Foundation spec in specs/memory-bridge.md, import/publish specs deferred to ideas/. Tasked out as S-0.1.1 through S-0.1.10 in ideas/TASKS.md.

---

## [2026-03-05-023937] Revised strategic analysis: blog-first execution order, bidirectional sync as top-level section

**Status**: Accepted

**Context**: Editorial review of ideas/claude-memory-strategic-analysis.md surfaced six structural weaknesses in competitive positioning

**Decision**: Revised strategic analysis: blog-first execution order, bidirectional sync as top-level section

**Rationale**: 200-line cap is fragile differentiator (demoted); org-scoped memory is the real threat (elevated to HIGH); model agnosticism is premature (parked with trigger condition); bidirectional sync is the most underweighted insight (promoted); narrative shapes categories before implementation does (blog first)

**Consequence**: Execution order is now S-3 (blog) -> S-0 -> S-1 -> S-2. Strategic doc restructured from 9 to 10 sections. Blog post shipped as first deliverable.

---

## [2026-03-04-105238] Interface-based GraphBuilder for multi-ecosystem ctx deps

**Status**: Accepted

**Context**: P-1.3 questioned whether non-Go dependency support would introduce bloat and whether a semantic approach was better

**Decision**: Interface-based GraphBuilder for multi-ecosystem ctx deps

**Rationale**: The output pipeline (map[string][]string to Mermaid/table/JSON) was already language-agnostic. Each ecosystem builder is ~40 lines — this is finishing what was started, not bloat. Static manifest parsing (no external tools for Node/Python) keeps dependencies minimal.

**Consequence**: ctx deps now auto-detects Go, Node.js, Python, Rust. --type flag overrides detection. ctx-architecture skill works across ecosystems without changes.

---

## [2026-03-02-165038] Billing threshold piggybacks on check-context-size, not heartbeat

**Status**: Accepted

**Context**: User wanted a configurable token-count nudge for billing awareness (Claude Pro 1M context, extra cost after 200k). Heartbeat produces zero stdout and can't relay to user.

**Decision**: Billing threshold piggybacks on check-context-size, not heartbeat

**Rationale**: check-context-size already reads tokens, has VERBATIM relay working, and runs every prompt. Adding a third independent trigger there is minimal code and follows established patterns.

**Consequence**: New .ctxrc field billing_token_warn (default 0 = disabled). One-shot per session via billing-warned-{sessionID} state file. Template-overridable via check-context-size/billing.txt.

---

## [2026-03-02-123611] Replace auto-migration with stderr warning for legacy keys

**Status**: Accepted

**Context**: Auto-migration code existed for promoting keys from ~/.local/ctx/keys/ and .context/.ctx.key to ~/.ctx/.ctx.key. Userbase is small and this is alpha — no need to bloat the codebase.

**Decision**: Replace auto-migration with stderr warning for legacy keys

**Rationale**: Warn-only is simpler, avoids silent file operations, and puts the user in control. Migration instructions in docs are sufficient for the small userbase.

**Consequence**: MigrateKeyFile() now only warns on stderr. promoteToGlobal() helper deleted. Tests verify keys are not moved.

---

## [2026-03-02-005213] Consolidate all session state to .context/state/

**Status**: Accepted

**Context**: Session-scoped state (cooldown tombstones, pause markers, daily throttle markers) was split between /tmp (via secureTempDir()) and .context/state/ for project-scoped state

**Decision**: Consolidate all session state to .context/state/

**Rationale**: Single location simplifies mental model, eliminates duplicated secureTempDir() in two packages, removes the cleanup-tmp SessionEnd hook entirely. .context/state/ is already gitignored and project-scoped.

**Consequence**: All 18 callers updated. Tests switch from XDG_RUNTIME_DIR mocking to CTX_DIR + rc.Reset(). Hook lifecycle drops from 4 events to 3 (SessionEnd removed).

---

## [2026-03-01-222733] PersistentPreRunE init guard with three-level exemption

**Status**: Accepted

**Context**: ctx commands handled missing .context/ inconsistently — some caught errors, some got confusing file-not-found messages, some produced empty output

**Decision**: PersistentPreRunE init guard with three-level exemption

**Rationale**: Single PersistentPreRunE on root command gives one clear error. Three-level exemption (hidden commands, annotated commands, grouping commands) covers all edge cases without per-command boilerplate

**Consequence**: Boundary violation now returns an error instead of os.Exit(1), making it testable. The subprocess-based boundary test was simplified to a direct error assertion

---

## [2026-03-01-161457] Global encryption key at ~/.ctx/.ctx.key

**Status**: Superseded by [2026-03-02] global key simplification

**Context**: Key stored next to ciphertext (.context/.ctx.key) was a security antipattern and broke in worktrees. The slug-based per-project key system at ~/.local/ctx/keys/ was over-engineered for the common case (one user, one machine, one key).

**Decision**: Single global key at ~/.ctx/.ctx.key. Project-local override via .ctxrc key_path or .context/.ctx.key.

**Rationale**: One key per machine covers 99% of users. Per-project slug filenames and three-tier resolution added complexity without clear benefit. ~/.ctx/ is the natural home (matches ~/.claude/ convention). Tilde expansion in .ctxrc key_path fixes a standalone bug.

**Consequence**: Auto-migration promotes legacy keys (project-local, ~/.local/ctx/keys/) to ~/.ctx/.ctx.key. Deleted KeyDir(), ProjectKeySlug(), ProjectKeyPath(). ResolveKeyPath simplified to two params. 15+ doc files updated.

---

## [2026-03-01-112544] Heartbeat token telemetry: conditional fields, not always-present

**Status**: Accepted

**Context**: Adding tokens, context_window, usage_pct to heartbeat payloads. First prompt of a session has no JSONL usage data yet.

**Decision**: Heartbeat token telemetry: conditional fields, not always-present

**Rationale**: Token fields are only included in the template ref when tokens > 0. This avoids misleading pct=0% on the first heartbeat and keeps payloads clean for receivers that filter on field presence.

**Consequence**: Webhook consumers must handle heartbeats both with and without token fields. The message string also varies (with/without tokens=N pct=N% suffix).

---

## [2026-03-01-092613] Hook log rotation: size-based with one previous generation, matching eventlog pattern

**Status**: Accepted

**Context**: .context/logs/ files grow unbounded (~200KB after one month); needed a cap

**Decision**: Hook log rotation: size-based with one previous generation, matching eventlog pattern

**Rationale**: Architectural symmetry with eventlog, O(1) size check vs O(n) line counting, diagnostic logs don't need deep history (webhooks cover serious setups)

**Consequence**: Each log file caps at ~2MB (current + .1). config.LogMaxBytes = 1MB, same as EventLogMaxBytes

---

## [2026-03-01-090124] Promote 6 private skills to bundled plugin skills; keep 7 project-local

**Status**: Accepted

**Context**: Reviewed all 13 _ctx-* private skills to determine which are universally useful for any ctx user vs specific to the ctx codebase or personal infra.

**Decision**: Promote 6 private skills to bundled plugin skills; keep 7 project-local

**Rationale**: Promote if the skill benefits any ctx-powered project without project-specific hardcoding. Keep private if it references this repo's Go internals, personal infra, or language-specific tooling. Promote list: _ctx-spec (generic scaffolding), _ctx-brainstorm (design facilitation), _ctx-verify (claim verification), _ctx-skill-creator (skill authoring), _ctx-check-links (doc link audit), _ctx-sanitize-permissions (Claude Code permissions audit). Keep list: _ctx-audit (Go/ctx checks), _ctx-qa (Go Makefile), _ctx-backup (SMB infra), _ctx-release/_ctx-release-notes (ctx release workflow), _ctx-update-docs (ctx package mapping), _ctx-absorb (borderline, revisit later).

**Consequence**: Six skills move from .claude/skills/ to internal/assets/claude/skills/ and become available to all ctx users via ctx init. Cross-references between skills need updating (e.g., /_ctx-brainstorm becomes /ctx-brainstorm). The seven remaining private skills stay project-local.

---

## [2026-02-27-230718] Context window detection: JSONL-first fallback order

**Status**: Accepted

**Context**: check-context-size defaults to 200k but user runs 1M-context model, causing false 110% warnings. JSONL contains the model name which maps to actual window size.

**Decision**: Context window detection: JSONL-first fallback order

**Rationale**: effective_window = detect_from_jsonl(model) ?? ctxrc.context_window ?? 200_000. JSONL is ground truth (reflects actual model in use); ctxrc is fallback for first-hook-of-session or unknown models; 200k is safe last resort. Having ctxrc override JSONL would artificially restrict the check when a user forgets to update their config after switching models.

**Consequence**: Most users get correct window automatically. ctxrc context_window becomes a fallback, not an override. Task exists for implementation.

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

**Consequence**: Deleted `internal/cli/session/` (15 files), removed auto-save hook, removed `--auto-save` from watch, removed pre-compact auto-save from compact, removed `/ctx-save` skill, updated ~45 documentation files. Four earlier decisions superseded (SessionEnd hook, Auto-Save Before Compact, Session Filename Format, Two-Tier Persistence Model). Users who want session history use `ctx recall list/export` instead.

---

*Module-specific, already-shipped, and historical decisions:
[decisions-reference.md](decisions-reference.md)*
