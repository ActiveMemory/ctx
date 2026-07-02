# Decisions

<!-- INDEX:START -->
| Date | Decision |
|----|--------|
| 2026-06-19 | `ctx ai` is a separate namespace with proposed-patch output only |
| 2026-06-07 | ctx-dream executor is a documented contract, not a hardcoded cron/claude assumption |
| 2026-06-07 | Output belongs in write/ — taxonomy and emission style (consolidated) |
| 2026-06-07 | Package taxonomy and shared-code placement (consolidated) |
| 2026-06-07 | Error handling: centralized in internal/err, domain-file taxonomy (consolidated) |
| 2026-06-07 | config/ as constants home and the magic-value audit (consolidated) |
| 2026-06-07 | YAML text externalization, init, and drift guards (consolidated) |
| 2026-06-07 | CWD-anchored context model (consolidated) |
| 2026-06-07 | Encryption key resolution and migration (consolidated) |
| 2026-06-07 | ctxctl maintainer binary and out-of-band audit channel (consolidated) |
| 2026-06-07 | KB editorial pipeline (Phase KB) design (consolidated) |
| 2026-06-07 | Companion-tool integration: peer-MCP, no gateway (consolidated) |
| 2026-06-07 | Localizable vocabulary and i18n primitives (consolidated) |
| 2026-06-07 | Embedded assets and editor-integration harnesses (consolidated) |
| 2026-06-07 | Context injection, hooks, and session-state architecture (consolidated) |
| 2026-06-06 | ctx-dream: standalone proposing memory consolidator (Option B), human-gated via serendipity |
| 2026-05-31 | Journal screen renders ctx journal source verbatim instead of parsing it |
| 2026-05-31 | ctx Desktop shells out via std::process::Command, not tauri-plugin-shell |
| 2026-05-30 | Name the add JSON-ingest flag --json-file, not --json |
| 2026-05-28 | Memory pressure detection uses OS-native signals (macOS pressure level + Linux PSI), not occupancy |
| 2026-05-24 | Pad snapshot-on-mutate at the store.WriteEntries choke point |
| 2026-05-20 | Gitignore .context/handovers/; track only .gitkeep |
| 2026-04-16 | Deprecate and remove ctx backup |
| 2026-04-14 | doc.go quality floor: behavior-grounded, ~25-100 body lines, related-packages section required |
| 2026-04-14 | Bootstrap stays under ctx system bootstrap (reverted experimental top-level promotion) |
| 2026-04-14 | Title Case style for docs is AP-leaning with explicit ambiguity carve-outs |
| 2026-04-11 | Journal stays local; LEARNINGS.md is the shareable layer |
| 2026-04-11 | `Entry.Author` is server-authoritative, not client-authoritative |
| 2026-04-09 | Architecture skill pipeline is a triad not a quartet |
| 2026-04-08 | Remove #done tag convention, simplify task archival |
| 2026-04-06 | Use hook relay for session provenance instead of JSONL parsing or env vars |
| 2026-04-01 | IRC to Discord as primary community channel |
| 2026-04-01 | AST audit tests live in internal/audit/, one file per check |
| 2026-04-01 | Rename ctx hook → ctx setup to disambiguate from the hook system |
| 2026-03-31 | Split log into log/event and log/warn to break import cycles |
| 2026-03-30 | Flags-not-subcommands for journal source: list and show are view modes on a noun, not independent entities |
| 2026-03-30 | Journal consumed recall — recall CLI package deleted |
| 2026-03-25 | Architecture analysis and enrichment are separate skills — constraint is the feature |
| 2026-03-25 | Prompt templates removed — skills are the single agent instruction mechanism |
| 2026-03-24 | Write-once baseline with explicit end-consolidation for consolidation lifecycle |
| 2026-03-18 | Singular command names for all CLI entities |
| 2026-03-16 | Rename --consequences flag to --consequence for singular consistency |
| 2026-03-14 | System path deny-list as safety net, not security boundary |
| 2026-03-14 | Config-driven freshness check with per-file review URLs |
| 2026-03-13 | Delete ctx-context-monitor skill — hook output is self-sufficient |
| 2026-03-12 | Rename ctx-map skill to ctx-architecture |
| 2026-03-06 | Drop fatih/color dependency — Unicode symbols are sufficient for terminal output, color was redundant |
| 2026-03-05 | Gitignore .context/memory/ for this project |
| 2026-03-04 | Interface-based GraphBuilder for multi-ecosystem ctx deps |
| 2026-03-02 | Billing threshold piggybacks on check-context-size, not heartbeat |
| 2026-03-01 | PersistentPreRunE init guard with three-level exemption |
| 2026-03-01 | Heartbeat token telemetry: conditional fields, not always-present |
| 2026-03-01 | Promote 6 private skills to bundled plugin skills; keep 7 project-local |
| 2026-02-27 | Context window detection: JSONL-first fallback order |
| 2026-02-26 | ctx init and CLAUDE.md handling (consolidated) |
| 2026-02-26 | Task and knowledge management (consolidated) |
| 2026-02-26 | Agent autonomy and separation of concerns (consolidated) |
| 2026-02-26 | Security and permissions (consolidated) |
| 2026-02-27 | Webhook and notification design (consolidated) |
| 2026-04-25 | Use t.Setenv for subprocess env in tests, not append(os.Environ(), ...) |
<!-- INDEX:END -->

<!-- DECISION FORMATS

## Quick Format (Y-Statement)

For lightweight decisions, a single statement suffices:

> "In the context of [situation], facing [constraint], we decided for [choice]
> and against [alternatives], to achieve [benefit], accepting that [trade-off]."

## Full Format

For significant decisions:

## [YYYY-MM-DD] Decision Title

**Status**: Accepted | Superseded | Deprecated

**Context**: What situation prompted this decision? What constraints exist?

**Alternatives Considered**:
- Option A: [Pros] / [Cons]
- Option B: [Pros] / [Cons]

**Decision**: What was decided?

**Rationale**: Why this choice over the alternatives?

**Consequence**: What are the implications? (Include both positive and negative)

**Related**: See also [other decision] | Supersedes [old decision]

## When to Record a Decision

✓ Trade-offs between alternatives
✓ Non-obvious design choices
✓ Choices that affect architecture
✓ "Why" that needs preservation

✗ Minor implementation details
✗ Routine maintenance
✗ Configuration changes
✗ No real alternatives existed

-->

## [2026-06-19-064500] `ctx ai` is a separate namespace with proposed-patch output only

**Status**: Accepted

**Context**: GitHub issue #92 (`ctx ai backend`) introduces an optional,
local-first AI backend layer. The governing spec left several Block A choices
open: whether AI use should live under a new `ctx ai <verb>` namespace or as
flags on existing commands, where the validation consumer should write
proposals, whether Block A needs a companion skill, and how to reconcile the
spec's TOML examples with the repo's existing YAML `.ctxrc` parser.

**Decision**: Block A uses a new top-level `ctx ai <verb>` namespace. Backend
configuration uses the existing YAML `.ctxrc` shape under `backends:` rather
than TOML-style `[backends]` tables. `ctx setup --backend <name>` is a distinct
setup mode that can run without the current `<tool>` positional argument. The
provisional proposal queue is `.context/proposals/ai/`, and the Block A
validation consumer is `ctx ai propose <input> --emit ...`, a generic
validation-only proposer that writes proposed-patch JSON artifacts and never
mutates `.context/*.md`. No new companion skill ships in Block A; command
assets, setup docs, the `ctx ai` CLI reference, the vLLM recipe, and the agent
playbook note cover the user-facing surface.

**Rationale**: A separate `ctx ai` namespace keeps optional AI behavior out of
deterministic commands (`ctx status`, `ctx agent`, ceremonies, and hooks), makes
fail-closed behavior explicit, and avoids smuggling backend availability into
existing deterministic verbs through flags like `--use-ai`. YAML config matches
the actual `internal/rc` implementation and avoids landing a spec that cannot
be implemented without replacing the config parser. `.context/proposals/ai/`
keeps AI-produced suggestions inside the cognitive substrate while preserving
the human gate before canonical memory changes. A generic `propose` verb proves
backend dispatch and artifact writing without pretending to settle the final
Block B taxonomy; later `ctx ai compact` or `ctx ai ingest` commands remain
available once the extraction-and-recall spec is promoted.

**Consequence**: Implementations must add a deterministic-boundary guard so
agent/status/ceremony paths cannot import `internal/backend` or invoke `ctx ai`.
Multiple configured backends require `--backend <name>` or `backends.default`;
there is no implicit selection and no deterministic fallback. Proposed-patch
artifacts need a minimal schema with backend, model, input reference, emit
kinds, proposed rows, source spans or citations when available, and status
metadata. Documentation and command assets are part of the same deliverable
because `ctx ai` is a user-facing surface. Spec: `specs/ctx-ai-backend.md`.

---

## [2026-06-07-112203] ctx-dream executor is a documented contract, not a hardcoded cron/claude assumption

**Status**: Accepted

**Context**: Settling ctx-dream v1 open questions. The executor runs the
out-of-band dream pass (read ideas/, classify+ground, write proposals). Question
was cron 'claude -p' vs a raw Anthropic-API scheduled loop.

**Decision**: ctx-dream executor is a documented contract, not a hardcoded
cron/claude assumption

**Rationale**: cron 'claude -p' is the reference executor (reuses Claude Code
auth, tool-calling, and PreToolUse hooks so the three guards are structural for
free; matches the existing skill draft and the cheap-validation goal). But we
must NOT assume it is the only executor: other harnesses (different AI CLI, raw
API loop, CI runner) must be able to run the same dream. So ctx owns an
executor-agnostic Go core (dreams/ layout, state record, ledger, proposal
schema, the three guards as callable logic) and the executor is a documented
contract: run one bounded pass, enforce the three guards STRUCTURALLY (Claude
Code via PreToolUse hooks; API loop via in-loop tool executor), fail loud, write
proposals-only into dreams/. Dream is opt-in, not enabled by default.

**Consequence**: Guards live as reusable Go logic in internal/dream/, not only
as a hook script. Two user-facing docs are required: a Claude Code enablement
guide and an executor-contract reference for other harnesses. The serendipity
review skill is split into its own spec (specs/ctx-serendipity.md). v1 ships the
cron/claude-p reference path but the data contract + guards stay
executor-portable.

---

## [2026-06-07-180001] Output belongs in write/ — taxonomy and emission style (consolidated)

**Consolidated from**: 3 entries (2026-03-17 to 2026-04-03)

- Output functions belong in write/ (flat by domain, one package per CLI
  feature); core/ owns logic and types, cmd/ owns Cobra orchestration. No
  cmd.Print* calls in internal/cli/ outside internal/write/ — enables
  localization and clean separation.
- Within write/, use pre-compute-then-print: functions with 4+ Printlns
  pre-compute conditional strings then emit one multiline block (TplXxxBlock),
  rejecting text/template (runtime errors, only 38/160 functions benefit);
  trivial and loop-based functions stay imperative.

---

## [2026-06-07-180002] Package taxonomy and shared-code placement (consolidated)

**Consolidated from**: 6 entries (2026-03-06 to 2026-05-17)

- Three-zone taxonomy: cmd/ for Cobra wiring, core/ for logic and types, assets/
  for templates and user-facing text; config/ for structural constants only.
  Symmetry makes navigation agent-friendly; shared domain types live in domain
  packages (internal/entry), not CLI subpackages.
- Pure-logic functions return data structs; callers own I/O, file writes, and
  reporting — lets MCP and CLI callers control output independently.
  Receiver-stateless methods become free functions; callbacks that vary only by
  a string key become text-key data.
- Shared formatting utilities (Pluralize, Duration, TruncateFirstLine, etc.)
  live in internal/format, not duplicated across CLI subpackages.
- internal/parse is the home for shared text-to-typed-value conversions
  (parse.Date first), scoped to avoid becoming a junk drawer.
- Every cross-package type goes in internal/entity/ — the cross-package-types
  audit (zero grandfathered violations) is the hardline; entity.Sentinel lives
  there even though it is a behavioral helper, over per-package duplication
  across 9 err packages.
- Multi-segment directory paths are single composite constants
  (DirHooksMessages, DirMemoryArchive), not joined from segment constants.

---

## [2026-06-07-180003] Error handling: centralized in internal/err, domain-file taxonomy (consolidated)

**Consolidated from**: 2 entries (2026-03-06 to 2026-03-14)

- Errors centralize in internal/err, not per-package err.go files — single
  location makes duplicates visible, enables sentinel errors, prevents
  broken-window accumulation; all CLI err.go files migrated and deleted.
- The monolithic 1995-line errors.go (188 functions) was split into 22 domain
  files (backup, config, crypto, …, validation) named by responsibility, so
  error constructors are findable by domain.

---

## [2026-06-07-180004] config/ as constants home and the magic-value audit (consolidated)

**Consolidated from**: 4 entries (2026-03-23 to 2026-04-04)

- String-typed enums (type Foo string + const blocks) belong in config/, not
  domain packages — types without behavior live in config; promote to entity/
  only when methods/interfaces appear.
- TestNoMagicStrings/TestNoMagicValues dropped the const/var exemption outside
  config/ (it masked 156+ string and 7 numeric constants in the wrong place);
  naming a constant in the wrong package does not fix the structural problem.
- The 60+ config/ sub-package "explosion" is correct, not a bottleneck: Go's
  compile unit is the package, so granular packages give precise dependency
  tracking and minimal recompile; the DX cost is fixed by a README decision
  tree, not restructuring.
- Cross-package magic strings (e.g. <pre> HTML tags used by normalize and
  format) promote to shared config constants (config/marker TagPre/TagPreClose);
  package-local copies deleted.

---

## [2026-06-07-180005] YAML text externalization, init, and drift guards (consolidated)

**Consolidated from**: 5 entries (2026-03-13 to 2026-04-03)

- All user-facing text externalizes to embedded YAML domain files
  (commands/flags/text/examples split via dedicated loaders), justified by agent
  legibility (named DescKey constants as traversable graphs) and drift
  prevention, not i18n; the 3-file ceremony (DescKey + YAML + write/err fn) is
  the accepted cost.
- Static embedded data and resource lookups use an explicit Init() called
  eagerly at startup, not per-accessor sync.Once or package-level init() —
  makes the startup dependency visible and testable; maps unexported, accessors
  are plain lookups.
- A Go↔YAML linkage check (lint-drift check 5, shell grep+comm) catches
  orphaned/broken DescKey↔YAML links and cross-namespace duplicates at CI
  time, preventing silent runtime failures.
- The build target depends on sync-why so derived assets/why/ files cannot drift
  from their docs/ sources — build fails without sync.
- MCP resource name constants live in config/mcp/resource (parallel to
  config/mcp/tool); the resource→file mapping stays in server/resource (too
  many cross-cutting deps for a config package), pre-built once at server init
  for O(1) lookup.

---

## [2026-06-07-180006] CWD-anchored context model (consolidated)

**Consolidated from**: 5 entries (2026-04-13 to 2026-05-21)

- Walk boundary uses git as a hint, not a requirement: walkForContextDir
  consults findGitRoot to anchor ancestor .context candidates and falls back to
  CWD when no git is found — fixes nested-repo binding without making git
  mandatory or relying on unreliable project markers.
- ctx activate is strict-CWD (drop upward walk): state-setting commands follow
  git's read-vs-state pattern (read walks, state refuses to cross repo
  boundaries); workspace-shared layouts are preserved by user action (cd first),
  not inferred walk.
- Anchor ctx to CWD entirely: drop activate/deactivate, the env-var (CTX_DIR)
  resolver, and all walks. With .context/ mandated as .git/'s sibling, every
  resolver collapses to os.Stat; keeping any walk would force maintaining two
  implementations. Mental model matches helm/terraform/Claude Code; ~600-1000
  LOC net deletion (specs/cwd-anchored-context.md).
- Spec steps 1+2 (resolver swap + init-guard removal) merged into one commit
  because step 1 cannot compile without step 2; cleanest commit boundaries beat
  strict spec adherence — remaining steps stay discrete (4-commit
  decomposition, not the spec's 5).
- Substrate vs. artifact placement: cognitive substrate (read AND written via
  ctx-mediated paths) lives under .context/; project artifacts (read/edited
  directly by humans, e.g. specs/, CLAUDE.md, docs/) live at root. kb passes all
  three coupling tests (mediated queries, pipeline coupling, skill discipline)
  so it stays under .context/.

---

## [2026-06-07-180007] Encryption key resolution and migration (consolidated)

**Consolidated from**: 3 entries (2026-03-01 to 2026-06-02)

- Single global key at ~/.ctx/.ctx.key (matches ~/.claude/ convention); one key
  per machine covers ~99% of users. Replaced the over-engineered slug-based
  per-project key system; project-local key-next-to-ciphertext was a security
  antipattern that broke in worktrees. [Original 2026-03-01 entry was marked
  Superseded by the 2026-03-02 simplification.]
- Legacy-key auto-migration replaced with a stderr warning only: warn-only is
  simpler, avoids silent file operations, and keeps the (small, alpha) userbase
  in control; docs carry migration instructions.
- Removed the implicit project-local .context/.ctx.key auto-detection tier from
  ResolveKeyPath: resolution is now (1) explicit .ctxrc key_path, (2) global
  ~/.ctx/.ctx.key, (3) project-local only as a degenerate fallback when home is
  unavailable. The local tier was the only thing making worktrees differ from
  side-by-side terminals; its removal is net deletion, and the previously-silent
  fire-path decrypt failure is now surfaced.

---

## [2026-06-07-180008] ctxctl maintainer binary and out-of-band audit channel (consolidated)

**Consolidated from**: 4 entries (2026-05-24 to 2026-05-28)

- Discipline enforcement belongs on the verbatim-relay channel, run out-of-band:
  relay is the one discipline channel that survives tunnel vision; run the
  auditor in a separate Claude Code session for fresh-context judgment and cost
  control. New generic channel: a skill writes .context/audit/<kind>.md, a
  check-audit hook relays unread reports verbatim, ctx audit list/show/dismiss
  manages lifecycle (digest-bound dismissal).
- [Superseded] ctxctl first placed at cmd/ctxctl in the same Go module:
  binary-level isolation via transitive-import exclusion, zero relocation of
  existing internal/audit files, on the belief a separate go.mod couldn't import
  the parent's internal/.
- That belief was empirically disproved: a nested module lexically under the
  parent path CAN import internal/. So ctxctl became a separate Go module at
  tools/ctxctl (own go.mod) — a hard module boundary guarantees ctx can never
  import ctxctl (the asymmetric requirement that matters); one-directional
  ctxctl→ctx coupling is acceptable for disposable maintainer tooling. A
  go.work wires the workspace; a guard test asserts cmd/ctx never imports
  internal/ctxctl.
- ctxctl is PATH-installed alongside ctx (build to dist/, install to
  /usr/local/bin/ctxctl) for clean repo roots and one binary across all
  worktrees, mirroring ctx's install pattern; the local hook calls ctxctl from
  PATH.

---

## [2026-06-07-180009] KB editorial pipeline (Phase KB) design (consolidated)

**Consolidated from**: 6 entries (2026-05-10 to 2026-05-16)

- Lift the sibling clean-room project's battle-tested editorial pipeline into
  ctx as v1, paired with handover: it is field-tested under production use and
  your-project is already paying the workaround tax (N=1 lived validation); lift
  the whole shape with a non-colliding rename, not hedge-and-defer.
- Mandate git as an architectural precondition: persistent-memory is dishonest
  without an undo layer (git reflog); refuse-on-no-git rather than auto-git-init
  (ctx never modifies the filesystem outside .context/); eliminates commit:none
  dead-code branches. Breaking change in next minor.
- KB ontology is pipeline-only-writer; no /ctx-kb-decide skill: in a KB you
  don't decide, you increase confidence — even NL assertions are
  evidence-capture events, not decision-capture. KB surface stays small (4 mode
  skills + ctx kb note); canonical capture skills unchanged.
- Phase KB ships handover + editorial paired, not split: the closeout/fold
  mechanism is the integration point; shipping paired stresses the fold on day
  one rather than retrofitting it.
- Editorial constitution lives at .context/ingest/KB-RULES.md, not
  CONSTITUTION.md: lifts the sibling project's resolved naming-collision (their
  10-INGEST_RULES.md rename) so ctx CONSTITUTION.md keeps its singular meaning;
  same discipline carries to domain-decisions.md vs DECISIONS.md.
- Phase KB lifts the *current* upstream pipeline shape (pass-mode contract,
  completion circuit breaker, source-coverage state-machine ledger,
  topic-adjacency pre-flight, cold-reader rubric, folder-shaped topics from day
  one, CLI-as-scaffold-authority), superseding the brief's 4-phase model —
  lifting the older shape would re-fight wounds the upstream author already
  healed.

---

## [2026-06-07-180010] Companion-tool integration: peer-MCP, no gateway (consolidated)

**Consolidated from**: 6 entries (2026-03-06 to 2026-05-23)

- Peer MCP model for external tools (GitNexus, context-mode): side-by-side
  servers each queried independently by the agent, chosen over orchestrator/hub
  models to respect ctx's markdown-on-filesystem invariant and avoid
  coupling/plugin registries.
- Skills stay CLI-based; MCP Prompts are the protocol equivalent: CLI is always
  available (PATH prereq), MCP is optional config, hooks are always CLI — two
  access patterns in one tool is gratuitous complexity.
- Recommend companion RAGs as peer MCP servers, not bridged through ctx: MCP is
  the composition layer; ctx is context, RAGs are intelligence — no bridging,
  plugin system, or schema abstraction.
- Companion tools documented as optional MCP enhancements with a runtime check
  (/ctx-remember smoke-tests MCPs at session start; companion_check:false
  suppresses) so users learn what enhances their workflow without being forced
  to install.
- MCP gateway not worth the coupling cost: a gateway would make ctx own
  install/uninstall/version/error-surface for tools it doesn't ship
  (bidirectional ownership coupling); composition is already MCP's job and the
  skills already work peer-to-peer. The pluggable-graph-tool task was skipped as
  a direct consequence (pluggability without ownership is incoherent).
- Skill body text uses capability-first language with canonical tools as
  examples; install-guide docs name canonical implementations directly
  (newcomers need a recommendation); allowed-tools frontmatter stays
  MCP-specific (genericizing to mcp__* is a permission expansion). Pure text
  rewrite, no new abstraction layer.

---

## [2026-06-07-180011] Localizable vocabulary and i18n primitives (consolidated)

**Consolidated from**: 5 entries (2026-03-14 to 2026-05-23)

- Session prefixes are parser vocabulary, not i18n text: header-recognition
  patterns move to .ctxrc session_prefixes (default Session:), separating
  content recognition from interface language so users parse multilingual
  session files without code changes.
- Classify rules are user-configurable via .ctxrc (classify_rules overrides
  config/memory defaults) — same pattern as session_prefixes, for
  non-English/specialized domains.
- Spec signal words and the nudge threshold (spec_signal_words,
  spec_nudge_min_len) are .ctxrc-configurable — signal words are language- and
  project-dependent.
- Keep i18n.Fold strict (Unicode case-fold, İ≠i, for identifier
  dedup/parsing/security comparison); add i18n.MatchKey (Fold + NFKD + strip
  combining marks) as a separate diacritic-insensitive primitive for matching
  user input against vocabulary lists. Two explicit-contract primitives beat one
  conflated primitive or an options flag.
- Placeholder overrides use EXTEND, not REPLACE, semantics (diverging from
  SessionPrefixes' REPLACE): the dominant bilingual EN+TR case needs both
  default and added placeholders rejected simultaneously; REPLACE would silently
  lose baseline coverage. Opt-in placeholders_replace:true reserved if REPLACE
  is later wanted.

---

## [2026-06-07-180012] Embedded assets and editor-integration harnesses (consolidated)

**Consolidated from**: 7 entries (2026-04-01 to 2026-05-22)

- Embedded foreign-language assets (TS/Bash/PowerShell/YAML) under
  internal/assets/ are intentional, not a smell: every file is //go:embed'd into
  the ctx binary and written at ctx setup; internal/ is about import privacy,
  not source language. The fix for the legibility gap was a contract README, not
  relocation (//go:embed can't reference ../).
- assets/hooks/ split into assets/integrations/ (tool-integration assets:
  Copilot instructions, AGENTS.md, CLI scripts/skills) + assets/hooks/messages/
  (hook-system templates) — integration assets are not hooks.
- Embedded harnesses (//go:embed'd, shipped via ctx setup) and
  separately-published harnesses (e.g. VS Code extension → marketplace, own
  cadence) are first-class peers with distinct CI/release pipelines; a new
  harness declares which pattern it follows before placing files.
- OpenCode plugin ships without a tool.execute.before hook: the natural fit
  (block-dangerous-commands) isn't a ctx Go subcommand and shimming would brick
  the editor (Cobra exit-1 read as {blocked:true}) on installs without the
  Claude wrapper. This omission is permanent — block-dangerous-commands will
  not be promoted to a ctx Go subcommand; the perpetually-pending re-add task is
  closed.
- Under cwd-anchored, the OpenCode plugin's agent shell tool can't be anchored
  to project root (the @opencode-ai/plugin SDK exposes only env, not cwd on
  shell.env); drop the shell.env handler and document launch-from-root.
  Plugin-internal ceremony calls stay anchored; the cwd-anchored error message
  is self-fixing.
- Editor-integration plugins must filter post-commit to actual git commit
  invocations (regex on the extracted command), not fire on every shell call —
  firing on noise trains users to ignore nudges.

---

## [2026-06-07-180013] Context injection, hooks, and session-state architecture (consolidated)

**Consolidated from**: 8 entries (2026-02-26 to 2026-05-08)

- Context injection v2: extract ~600 lines of diagrams out of FileReadOrder (53%
  token drop); auto-inject content via additionalContext (soft directives hit a
  ~75-85% compliance ceiling); imperative framing with an unconditional
  compliance checkpoint, verbatim relay as fallback. Inject
  CONSTITUTION/CONVENTIONS/ARCHITECTURE/PLAYBOOK verbatim, DECISIONS/LEARNINGS
  index-only, TASKS mention-only (~7,700 tokens).
- Context-load-gate injects only CONSTITUTION + AGENT_PLAYBOOK_GATE (~2k
  tokens), not the full ReadOrder: hard rules must be present pre-action;
  everything else is pulled on-demand. AGENT_PLAYBOOK_GATE.md must stay in sync
  with AGENT_PLAYBOOK.md.
- .context/state/ is the gitignored, project-scoped home for ephemeral runtime
  state (following the .context/logs/ precedent); all session state (cooldown
  tombstones, pause/throttle markers) consolidated there from /tmp, dropping the
  cleanup-tmp SessionEnd hook (4 hook events → 3).
- Gate mkdir inside state.Dir() rather than per-caller so "no .context/state/ in
  uninitialized projects" is structurally enforced; state.Dir() returns
  ErrNotInitialized (hook callers absorb silently, interactive callers surface a
  path-bearing message).
- Tighten state.Dir / rc.ContextDir to (string, error) with sentinel
  ErrDirNotDeclared: makes the empty-path case unrepresentable in a "looks fine"
  branch, closing the filepath.Join("", rel) trap that wrote state into CWD.
- Hook/notification design: prefer toning down docs claims over adding hooks
  (fatigue from 9 UserPromptSubmit hooks); hook output must be structured JSON
  (additionalContext), not plain text; dropped prompt-coach hook (zero useful
  tips, invisible channel); de-emphasized /ctx-journal-normalize (expensive,
  nondeterministic).
- Hook log rotation is size-based with one previous generation (current + .1,
  ~2MB cap), matching the eventlog pattern — O(1) size check, diagnostic logs
  don't need deep history.

---

## [2026-06-06-133805] ctx-dream: standalone proposing memory consolidator (Option B), human-gated via serendipity

**Status**: Accepted

**Context**: We explored whether ctx should grow a scheduled, background 'dream'
(a sleep-time memory process) and how it should relate to canonical memory. Felt
pain: the author's ideas/ folder is too overwhelming to triage, and canonical
files bloat over time (109 decisions, 151 learnings, 154 unimported sessions).
The risk to avoid: a background LLM job autonomously rewriting authoritative
memory and silently corrupting it (the research shows continuous LLM
consolidation is lossy and non-monotonic). Full debate:
.context/briefs/20260606T203414Z-ctx-dream-disciplined-consolidator.md

**Decision**: ctx-dream: standalone proposing memory consolidator (Option B),
human-gated via serendipity

**Rationale**: Chose a NEW, standalone, PROPOSING consolidator (Option B): it
writes only to its own sidecar + proposals queue + ledger + per-dream archive,
never autonomously to the five canonical files; a human 'serendipity' review
session is the sole bridge (accept/reject/amend) into canonical. One skill, two
modes: discipline (default; grounded, structured, provenanced proposals) and
creative/exploration (a safe relaxation: resurface + chance, reader-only).
Principle: decouple the cognition, reuse the plumbing (own the consolidation
logic; reuse import/enrich/kb-ingest via the enriched-journal data contract).
Standalone so mechanics evolve independently and changes to existing curation
skills can't break it, and for creative freedom (don't assume existing verbs
suffice). Discipline-first because it is the hard load-bearing substrate and
creative is a strict, safer relaxation of it. Grounded in
ideas/ctx-dreams/research: Auto-Dreamer (2605.20616) for the architecture,
'Useful Memories Become Faulty When Continuously Updated by LLMs' (2605.12978)
for the threat model, and the deep-research eval cluster for the finding that a
single agreeable LLM is not an adversarial gate (it silently repairs the missing
justification), which is why the gate must be human. Rejected: Option A (dream
owns a parallel canonical store, which does not fix bloat and creates two
divergent substrates); autonomous mutation / auto-approve (violates 'each memory
entry needs dedicated human attention'); pure-garden-only (under-serves
engineering's need for grounding and actionability); coupling to existing
skills' internals; garden-first build order.

**Consequence**: Positive: nothing autonomous touches canonical, so the system
is reversible by construction; the dream's mechanics can evolve freely; v1
(disciplined ideas/ triage, validated via a ctx-remind-nagged ~15-minute review
round) is low-stakes and validates the mechanism and author engagement cheaply.
Negative / trade-off: no human serendipity session = no consolidation, so the
dream's entire value is gated behind human review cadence, and the author
historically under-runs curation; mitigated only by ctx-remind nags + targeting
felt pain (ideas/) + a pleasure-not-chore framing. Validation of the full
product thesis (disciplined consolidation of canonical memory for engineering
teams) is deferred to a later test on a project where bloat actually bites. Spec
work proceeds via /ctx-spec --brief on the brief above; key mechanics remain
open (executor, proposal schema, ledger schema, .context/ layout).

---

## [2026-05-31-094649] Journal screen renders ctx journal source verbatim instead of parsing it

**Status**: Accepted

**Context**: The Journal timeline needs session data, but ctx journal source emits a whitespace-aligned table with no --json mode, and columns (usage/turns) drop out on short sessions.

**Decision**: Journal screen renders ctx journal source verbatim instead of parsing it

**Rationale**: A column parser misaligns when fields are absent; rendering the text verbatim in a monospace panel is honest and robust for P0.

**Consequence**: Journal is a styled text view, not structured cards; a proper timeline awaits a journal source --json mode upstream (same pattern as the artifact list --json commands).

---

## [2026-05-31-094649] ctx Desktop shells out via std::process::Command, not tauri-plugin-shell

**Status**: Accepted

**Context**: The GUI must run the ctx binary for every read and write. Tauri 2 offers tauri-plugin-shell with capability-scoped command allowlists.

**Decision**: ctx Desktop shells out via std::process::Command, not tauri-plugin-shell

**Rationale**: Running ctx inside our own #[tauri::command] via std::process::Command avoids all shell-plugin permission/capability wiring and keeps the adapter a single Rust module. ctx resolves its context from $PWD/.context, so each call sets current_dir to the selected project root.

**Consequence**: No shell capability in capabilities/default.json; the adapter owns PATH augmentation and git provenance synthesis; a CLI/output change is a one-file fix in ctx_adapter.rs.

---

## [2026-05-30-114429] Name the add JSON-ingest flag --json-file, not --json

**Status**: Accepted

**Context**: The CLI-FIX spec specified the literal flag --json <file>, but
--json is already a bool output-format flag across the CLI (ctx
status/drift/doctor/bootstrap --json all mean 'emit machine-readable output').

**Decision**: Name the add JSON-ingest flag --json-file, not --json

**Rationale**: Overloading --json as a string input-path on the add commands
would break that cross-command convention and confuse muscle memory. --json-file
is unambiguous, parallels the existing --file/-f source flag, and leaves -j
free. Pushed back on the spec's literal wording rather than satisfice.

**Consequence**: The add commands intentionally diverge from the spec's literal
--json; the spec was updated to reflect --json-file. Any future JSON-input flag
elsewhere should follow the --json-file naming, reserving --json for bool
output.

---

## [2026-05-28-200500] Memory pressure detection uses OS-native signals (macOS pressure level + Linux PSI), not occupancy

**Status**: Accepted

**Context**: `check-resource` alerted DANGER at swap-used ≥ 75% / memory-used
≥ 90% — pure occupancy. macOS swap is sticky (never recedes);
post-hibernation swap stays >75% with idle RAM, producing false "wrap up the
session" DANGER at session start. Memory occupancy on macOS includes reclaimable
cache — also a poor pressure proxy.

**Decision**: Memory pressure detection uses OS-native signals (macOS pressure
level + Linux PSI), not occupancy

**Rationale**: Occupancy is a level; pressure is a derivative. Only the kernel's
derivative reflects current struggle. macOS: `sysctl
kern.memorystatus_vm_pressure_level` (1/2/4 → OK/Warning/Danger). Linux:
`/proc/pressure/memory` (PSI) `some.avg10 ≥ 10.0` → warn, `full.avg10 ≥
10.0` → danger. Windows: filed as an exploratory task; unsupported for now
("other" platform falls through to `PressureSupported=false`, no alert).

**Consequence**: `MemInfo` gains `Pressure` + `PressureSupported`;
`threshold.go` drops both occupancy `byteCheck`s and emits a single pressure
alert. Doctor swap row removed (no longer a health signal); occupancy fields
retained for `ctx stats` display. PSI 10.0 defaults named in `config/stats` —
retunable in one place. `make lint` 0 issues, `make test` ok on the change.

---

## [2026-05-24-092912] Pad snapshot-on-mutate at the store.WriteEntries choke point

**Status**: Accepted

**Context**: Adding a safety net for accidental `ctx pad rm` (and any other
destructive pad mutation) required choosing where to insert the snapshot logic:
per-subcommand (in each cmd/<op>/run.go), or at the persistence choke point
(store.WriteEntriesWithIDs).

**Decision**: Pad snapshot-on-mutate at the store.WriteEntries choke point

**Rationale**: store.WriteEntriesWithIDs is invoked by every mutating pad
subcommand (add/edit/mv/rm/merge/normalize/resolve/tag and undo itself);
instrumenting it once gives universal coverage with one site of truth.
Per-subcommand instrumentation would need maintenance every time a new pad
mutation lands and is easy to forget. The snapshot itself is a byte-for-byte
copy of the existing pad blob (no re-encryption), so plaintext and encrypted
modes use identical logic; the existing ciphertext IS the snapshot.

**Consequence**: All future pad mutations get the safety net automatically
without per-command wiring. The op label for the snapshot filename is derived
from cmd.Name() at the call site, so the cmd parameter that already flowed in
for diagnostic output now carries semantic weight too. New constraint: any
future code path that bypasses WriteEntriesWithIDs to mutate the pad will
silently bypass the safety net — a guardrail test could enforce this if/when
that risk materializes.

---

## [2026-05-20-214753] Gitignore .context/handovers/; track only .gitkeep

**Status**: Accepted

**Context**: Per-session, operator-specific artifacts that grow without bound
and can leak host/internal identifiers (ari, asgard, broadcom-class) into public
mirrors when the project's .context/ is committed.

**Decision**: Gitignore .context/handovers/; track only .gitkeep

**Rationale**: Aligns with the existing per-personal-state gitignore family
(journal, memory, state, logs, reminders.json, scratchpad.enc); the directory's
.gitkeep keeps the read-side missing-dir gate passing on fresh clones; the rest
of the closeout-fold pipeline already lives in .context/archive/closeouts/ which
IS tracked.

**Consequence**: ctx init template (internal/config/file/ignore.go) added
.context/handovers/* and !.context/handovers/.gitkeep; existing tracked
handovers untracked via git rm --cached but kept on disk; the 'handover is the
sole authoritative recall artifact' phrasing in KB-RULES.md still holds — it's
local-machine authoritative.

---

## [2026-04-16-011520] Deprecate and remove ctx backup

**Status**: Accepted

**Context**: ctx backup is environment-specific (SMB/GVFS), fires nag hooks for
unconfigured users, and solves a problem that belongs to the OS layer. ctx hub
already handles cross-machine knowledge persistence.

**Decision**: Deprecate and remove ctx backup

**Rationale**: Hub handles persistence, backup is env-specific, wrong layer for
ctx to own. No external users depend on it. Broadcom mirror issue and GVFS
Linux-only dependency add maintenance burden.

**Consequence**: Need backup-strategy runbook before removal. Maintainer must
set up replacement cron job. About 60 files to remove across CLI, config, hooks,
docs, skills. Spec: specs/deprecate-ctx-backup.md

---

## [2026-04-14-010205] doc.go quality floor: behavior-grounded, ~25-100 body lines, related-packages section required

**Status**: Accepted

**Context**: About 140 doc.go files were rewritten this session. User flagged
the original 5-line Key exports + See source files + Part of subsystem pattern
as lazy minimum effort.

**Decision**: doc.go quality floor: behavior-grounded, ~25-100 body lines,
related-packages section required

**Rationale**: Behavior-grounded rewrites (read source first, then write) are
the only acceptable form for any non-trivial package. The lazy template
communicates nothing a future reader cannot grep for; it satisfies tooling
without adding signal.

**Consequence**: Every non-trivial package's doc.go now leads with the package's
actual purpose, names key behaviors, calls out non-obvious design choices
(Raft-lite, two-step indirection, idempotency contracts), and lists related
packages with paths. New packages should follow the same shape.

---

## [2026-04-14-010205] Bootstrap stays under ctx system bootstrap (reverted experimental top-level promotion)

**Status**: Accepted

**Context**: Mid-session promoted ctx bootstrap to top-level to make a stale
CLAUDE.md instruction work. User reverted it and reaffirmed the original design.

**Decision**: Bootstrap stays under ctx system bootstrap (reverted experimental
top-level promotion)

**Rationale**: The ctx system namespace is for agent and hook plumbing the user
does not type by hand. Bootstrap is invoked by AI agents at session start;
surfacing it at top-level pollutes ctx --help for humans without benefit.

**Consequence**: internal/bootstrap/group.go reverted;
internal/config/embed/cmd/system.go header now correctly states bootstrap is
intentionally not promoted. The CLAUDE.md template across the repo (and the
workspace copy) updated to reference ctx system bootstrap as canonical.

---

## [2026-04-14-010205] Title Case style for docs is AP-leaning with explicit ambiguity carve-outs

**Status**: Accepted

**Context**: Needed a deterministic Title Case engine for headings and
admonition titles across docs/. User precedent (Working with AI lowercase with)
ruled out strict Chicago.

**Decision**: Title Case style for docs is AP-leaning with explicit ambiguity
carve-outs

**Rationale**: AP lowercase prepositions regardless of length matches
user-approved titles. But strict AP would lowercase ambiguous prep/conj/adv
words like before, after, since, until, past, near, down, up, off, hurting
common cases. Carve-outs leave them at default-cap and let the engine reach a
sensible result for ~95 percent of headings without manual review.

**Consequence**: hack/title-case-headings.py ships an AP-leaning with ambiguity
carve-outs PREPOSITIONS set. Future style changes must touch that set explicitly
with reasoning. New brand or acronym additions go through the same audited
pattern.

---

## [2026-04-11-200000] Journal stays local; LEARNINGS.md is the shareable layer

**Status**: Accepted

**Context**: With the hub now carrying shared project context between machines
and eventually between teammates, the question came up whether enriched
journal entries should ride along — either the raw `.context/journal/` files
or an "export enriched entries as shareable learning items" pipeline layered
on top of `/ctx-journal-enrich`. The journal is already gitignored per the
2026-03-05 `.context/memory/` decision and for the same reason: it's a
first-person log of raw prompts, half-formed thoughts, dead ends, personal
names, and things the user talks through with themselves. It sits in the
same trust tier as shell history or a private notebook.

The trade-off is real: shared journals would make it trivial for teammates
(or future-me on another machine) to see the full reasoning trail behind a
decision. But "full reasoning trail" is precisely the thing that makes a
journal journal and not a changelog — it includes the parts the author
hasn't decided to stand behind yet, plus incidental private content.

**Decision**: The journal is **Tier-0 personal** and never leaves the
originating machine. No hub sync, no export-by-default, no
enriched-entries-as-shareable-items pipeline. The enrichment pipeline
(`/ctx-journal-enrich`) stays as-is: journal → human-in-the-loop review →
explicit promotion to LEARNINGS.md / DECISIONS.md / CONVENTIONS.md via the
existing `/ctx-learning-add`, `/ctx-decision-add`, `/ctx-convention-add`
commands. Those distilled artifacts are **Tier-1 shareable** and are what
the hub syncs when a team opts into shared context.

The promotion boundary is therefore the enrichment step, not a new export
pipeline. The user is the gate.

**Rationale**: Any "shareable enriched journal entry" pipeline would have to
re-implement the trust boundary that `/ctx-learning-add` already enforces:
the human decides what's worth sharing, strips incidental private content,
and rewrites it as a standalone artifact. A second pipeline that tries to
do this automatically would either (a) leak private content by accident, or
(b) require the same human review and thus collapse back into
`/ctx-learning-add`. The principled answer is that there is no second
pipeline — LEARNINGS.md *is* the shareable form of the journal.

This also preserves the psychological safety of the journal: the author
can write freely because they know nothing they write is one sync away
from a teammate's screen. Lose that property and the journal stops being a
journal and starts being a changelog draft.

**Consequence**:

- Journal files stay gitignored and stay out of `ctx hub` sync paths. Any
  future code that walks context files for replication must exclude
  `.context/journal/` explicitly and be covered by a test.
- `/ctx-journal-enrich` remains the promotion boundary. Its output targets
  are LEARNINGS.md / DECISIONS.md / CONVENTIONS.md, never a separate
  "shareable journal" bucket.
- Hub docs (`docs/home/hub.md`, `docs/recipes/hub-personal.md`,
  `docs/recipes/hub-team.md`, `docs/security/hub.md`) should state the
  Tier-0 / Tier-1 split explicitly so users building team workflows don't
  assume "shared context" means "shared everything."
- The sync code path in `internal/hub/sync_helper.go` and any future
  replication of context files must enforce this exclusion at the
  code level — a gitignore entry is a user-convenience signal, not a
  hub-trust boundary.
- A potential future "personal multi-machine journal sync" (same human,
  different laptops) is explicitly **out of scope** of this decision. If
  it ever ships, it rides a different transport (encrypted-at-rest,
  single-user, not the team hub) and needs its own decision record.

**Alternatives considered**:

- **Sync raw journal files via hub**: rejected. Inverts the gitignore
  decision, leaks private content by construction, destroys the
  journal's "safe to write freely" property.
- **Auto-export enriched entries as a new shareable artifact type**:
  rejected. Duplicates `/ctx-learning-add` without the human gate, or
  collapses back into it. No real difference from the status quo except
  the opportunity for accidental leakage.
- **Opt-in per-entry "publish to hub" flag in the journal**: rejected as
  premature. If the user wants an entry on the hub, the existing flow is
  one command away — write it as a learning or decision. A second path
  adds surface area without adding capability.

**Related**: Reinforces the 2026-03-05 `.context/memory/` gitignore
decision (same trust-tier reasoning for a different private artifact).

## [2026-04-11-180000] `Entry.Author` is server-authoritative, not client-authoritative

**Status**: Accepted

**Context**: The `Entry.Author` field on hub entries is copied verbatim from
the client's publish request (`handler.go:82`). It's optional, freeform, and
unauthenticated — a client with a valid token for project `alpha` can publish
entries claiming `Author: "bob@acme.com"` regardless of who actually
authenticated. This is the same spoofing pattern as `Origin` (audit finding
H-04) and was flagged as audit finding H-22 with three options: keep, drop,
override, or promote. The decision was never formally closed.

The premise that resolved it: **identity is eventually part of the token**.
Under the sysadmin-registry MVP, the server already knows `{user_id, project}`
from the authenticated token. Under the PKI stretch, the signed claim carries
identity cryptographically. In both models, the client has nothing to say about
authorship that the server doesn't already know with higher confidence.

**Decision**: `Entry.Author` is **server-authoritative**. The server stamps it
from the authenticated identity source on every publish. The client's
`pe.Author` input is ignored (or rejected — implementation choice, not
semantic difference). The field stays in the wire format but its semantics
change from "whatever the client said" to "whatever the server's auth layer
resolved."

Stamping source by phase:

- **Today (pre-registry)**: `Author = ClientInfo.ProjectName`, same source as
  the `Origin` server-enforcement fix (H-04). Lossy but consistent.
- **Registry MVP**: `Author = users.json` row's `user_id` (e.g.,
  `alice@acme.com`). Precise per-human attribution.
- **PKI stretch**: `Author = signed claim's sub field`. Cryptographic identity.

**Rationale**: Dropping the field is wrong because the registry MVP will
already give us a per-user identity to stamp — removing Author just to re-add
it later is churn. "Override" and "promote" are cosmetically different forms
of the same decision (server fills from auth context); "promote" is what
happens naturally once the registry MVP types the field as `UserID`.
Client-sourced Author is indefensible because it replicates the Origin
spoofing vector in a second field.

**Consequence**:

- The Author field stays on the wire and in `Entry{}`.
- Client-side code that populates `pe.Author` from local config becomes a
  no-op. Audit `ctx connect publish` and `ctx add --share` for any such
  code paths before the server-enforcement fix lands.
- `handler.go publish()` fills Author from the authenticated context (the
  same `ClientInfo` that H-04 pulls for Origin). Single unified
  auth-to-handler pipe.
- `docs/security/hub.md` "Compromised client token" section gets rewritten:
  attribution becomes **wrong** on compromise (attacker's token maps to
  attacker's identity), not **forgeable** (attacker cannot stamp someone
  else's name).
- The sysadmin-registry spec (`specs/hub-identity-registry.md`, tasked)
  MUST include a `user_id` field per row — it's the stamping source.
- Three open tasks collapse into one: H-22 resolves to "implement
  server-authoritative Author" instead of "decide Author fate." TASKS.md
  updated.

**Alternatives considered**:

- **Keep client-authoritative**: rejected. Same spoofing vector as Origin;
  trivially defeats any downstream attribution check.
- **Drop the field**: rejected. The registry MVP will need per-human
  attribution anyway. Dropping today is churn that gets undone
  immediately.
- **Override at client-side before publish**: rejected. Puts the security
  boundary on the wrong side of the trust zone. Must be server-side.

**Follow-up — client-advisory metadata**: the client still has useful
information to share that isn't an identity claim: a human-friendly
display name, the machine that made the publish, the tool version, a
CI system label, a team/role handle. This lives on a **new sibling
field `Meta`** (a `ClientMetadata` sub-struct), not on `Author`. The
separation of types is what protects the security property: `Author`
is reserved for server-authoritative identity, `Meta` is
client-advisory and explicitly labeled as such in any rendered
surface. `Meta` fields are size-capped individually (256 bytes) and
in aggregate (2 KB), validated for plain-string content (no
newlines, no control characters), and never claimed as attribution
in any API response. The renderer MUST label `Meta`-sourced values
with prose like "client label" or "client-reported" so readers
cannot mistake them for authoritative identity. See TASKS.md for
the implementation task.

---

## [2026-04-09-001332] Architecture skill pipeline is a triad not a quartet

**Status**: Accepted

**Context**: Had a proposed ctx-architecture-extend for extension point mapping,
making four skills

**Decision**: Architecture skill pipeline is a triad not a quartet

**Rationale**: Extension points already covered per-module in DETAILED_DESIGN
and by registration site discovery in enrich. Fourth skill fragments pipeline
without distinct value

**Consequence**: Pipeline is map enrich hunt. Three skills three questions: how
does it work, how well does it connect, where will it break

---

## [2026-04-08-013731] Remove #done tag convention, simplify task archival

**Status**: Accepted

**Context**: Tasks had #done:YYYY-MM-DD timestamps that agents added
inconsistently and nobody read. compact --archive filtered by age using these
timestamps.

**Decision**: Remove #done tag convention, simplify task archival

**Rationale**: [x] checkbox is semantically sufficient. git blame provides the
completion timestamp. Removing #done eliminates redundant ceremony and
simplifies compact --archive to archive all completed tasks regardless of age.

**Consequence**: compact --archive no longer filters by archive_after_days for
tasks. The .ctxrc field is inert but retained for backwards compatibility.
Historical #done tags in archives are preserved.

---

## [2026-04-06-204212] Use hook relay for session provenance instead of JSONL parsing or env vars

**Status**: Accepted

**Context**: Needed to give agents awareness of their session ID, branch, and
commit hash for task/decision/learning provenance. Considered three approaches:
(1) parsing most-recent JSONL at runtime, (2) CTX_SESSION_ID env var, (3) hook
relay via UserPromptSubmit.

**Decision**: Use hook relay for session provenance instead of JSONL parsing or
env vars

**Rationale**: JSONL parsing breaks with parallel sessions (wrong file picked).
Env vars aren't exported by Claude Code. Hook relay is zero-state: the hook
receives session_id from Claude Code on every prompt, emits it, agent absorbs
through repetition. No counters, no cleanup, no resume edge cases.

**Consequence**: Provenance depends on the hook being registered (enabledPlugins
in settings.local.json). Projects without plugin registration get no provenance.
Filed as separate bug.

---

## [2026-04-01-233247] IRC to Discord as primary community channel

**Status**: Accepted

**Context**: Discord server exists at https://ctx.ist/discord; IRC/libera.chat
references were stale

**Decision**: IRC to Discord as primary community channel

**Rationale**: Discord is faster for async community support; IRC was historical

**Consequence**: Updated zensical.toml, README, community docs, journal
template. Added community footer to ctx help and ctx init output via YAML assets
pipeline

---

## [2026-04-01-233246] AST audit tests live in internal/audit/, one file per check

**Status**: Accepted

**Context**: Needed a home for AST-based codebase invariant tests separate from
the existing compliance_test.go monolith

**Decision**: AST audit tests live in internal/audit/, one file per check

**Rationale**: One test per file prevents the 1200+ line monster pattern. Shared
helpers in helpers_test.go with sync.Once caching. Package is all _test.go
except doc.go — produces no binary, not importable

**Consequence**: New checks are added as individual *_test.go files; the pattern
(loadPackages, walk AST, collect violations, t.Error) is established and
repeatable

---

## [2026-04-01-074416] Rename ctx hook → ctx setup to disambiguate from the hook system

**Status**: Accepted

**Context**: PR #45 contributor assumed hook meant the setup command, causing
naming collisions with the PreToolUse/PostToolUse hook system

**Decision**: Rename ctx hook → ctx setup to disambiguate from the hook system

**Rationale**: hook has a specific meaning in ctx; setup accurately describes
generating AI tool integration configs

**Consequence**: CLI breaking change. All docs, specs, TypeScript extension, and
YAML assets updated. Released specs left as historical.

---

## [2026-03-31-224245] Split log into log/event and log/warn to break import cycles

**Status**: Accepted

**Context**: io and notify could not import log.Warn because log imported both
of them for event logging, creating circular dependencies

**Decision**: Split log into log/event and log/warn to break import cycles

**Rationale**: Separating concerns (stderr sink vs JSONL event log) into
subpackages eliminated the cycle. Warn sink is foundation-level with only config
imports, event logging is higher-level

**Consequence**: All stderr warnings now route through logWarn.Warn(). New code
importing log/warn has no cycle risk. Event types moved to internal/entity

---

## [2026-03-30-075927] Flags-not-subcommands for journal source: list and show are view modes on a noun, not independent entities

**Status**: Accepted

**Context**: During the journal-recall merge, recall had separate list and show
subcommands. Merging them into journal created a design choice: source list +
source show (three levels) vs source --show (two levels).

**Decision**: Flags-not-subcommands for journal source: list and show are view
modes on a noun, not independent entities

**Rationale**: Keeps CLI nesting to two levels max. Default behavior (bare
source) lists sessions; --show switches to inspect mode. When two operations
differ only in how they view the same data, make them flags on one command.

**Consequence**: journal source dispatches via --show flag rather than
positional subcommand. Future view-mode toggles should follow this pattern.

---

## [2026-03-30-003756] Journal consumed recall — recall CLI package deleted

**Status**: Accepted

**Context**: ctx recall was never registered in bootstrap; ctx journal had all
the same subcommands

**Decision**: Journal consumed recall — recall CLI package deleted

**Rationale**: One dead command group creates confusion in docs and skills.
Journal is the canonical command group.

**Consequence**: internal/cli/recall/ deleted, 19 doc files updated,
docs/cli/recall.md renamed to journal.md, zensical.toml updated. MCP tool
ctx_recall rename tasked separately (API contract)

---

## [2026-03-25-233646] Architecture analysis and enrichment are separate skills — constraint is the feature

**Status**: Accepted

**Context**: Observed that agents take shortcuts when code intelligence tools
are available during architecture analysis. A 5.2x depth reduction was measured
(5866 vs 1124 lines) when GitNexus was available during reading. Mentioning
unavailable tools by name in a skill plants the idea for the agent to use them.

**Decision**: Architecture analysis and enrichment are separate skills —
constraint is the feature

**Rationale**: Discovery requires forced reading without shortcuts. Validation
and quantification are a separate pass. Two-pass compiler analogy: semantic
parsing (human-style reading) then static analysis (graph enrichment). Never
mention tools you want the agent to avoid — absence is the only reliable
constraint.

**Consequence**: ctx-architecture deliberately excludes code intelligence tools
from allowed-tools and never mentions them. ctx-architecture-enrich is a
separate skill that runs after, using the deep artifacts as baseline. Gemini is
allowed in both for upstream/external lookups only.

---


## [2026-03-25-173336] Prompt templates removed — skills are the single agent instruction mechanism

**Status**: Accepted

**Context**: Prompt templates (.context/prompts/) overlapped with skills but had
no discoverability — even the project creator didn't know they existed

**Decision**: Prompt templates removed — skills are the single agent
instruction mechanism

**Rationale**: Adding metadata to prompts to fix discoverability would recreate
the skill system. One concept is better than two.

**Consequence**: code-review, explain, refactor promoted to proper skills. ctx
prompt CLI removed. loop.md retained as ctx loop config file at
.context/loop.md.

---

## [2026-03-24-001001] Write-once baseline with explicit end-consolidation for consolidation lifecycle

**Status**: Accepted

**Context**: Designing the consolidation nudge hook; multi-pass consolidation
spans dozens of sessions and you cannot programmatically distinguish feature
from consolidation sessions

**Decision**: Write-once baseline with explicit end-consolidation for
consolidation lifecycle

**Rationale**: First ctx-consolidate stamps baseline (write-once), user runs
end-consolidation when done. Failure mode is silence (no stale nudges), not
wrong behavior

**Consequence**: Requires mark-consolidation, end-consolidation, and
snooze-consolidation plumbing commands. Spec: specs/consolidation-nudge-hook.md

---


## [2026-03-18-193623] Singular command names for all CLI entities

**Status**: Accepted

**Context**: ctx add used learning (singular) but ctx learnings was plural.
Inconsistency across 6 commands.

**Decision**: Singular command names for all CLI entities

**Rationale**: Less headache for i18n; one rule (singular = entity); developers
think in OOP. Use field values come from DescKey constants for
single-source-of-truth renaming.

**Consequence**: All commands singular: task, decision, learning, change,
permission, dep. YAML keys, desc constants, directory names, and 50+ files
updated.

---

## [2026-03-16-022635] Rename --consequences flag to --consequence for singular consistency

**Status**: Accepted

**Context**: All other CLI flags (context, rationale, lesson, application) are
singular nouns. consequences was the only plural.

**Decision**: Rename --consequences flag to --consequence for singular
consistency

**Rationale**: Singular form matches the pattern. Consistency wins over natural
language preference.

**Consequence**: 75+ files updated. Breaking change for --consequences users.

---



## [2026-03-14-110748] System path deny-list as safety net, not security boundary

**Status**: Accepted

**Context**: Replacing nolint:gosec directives with centralized I/O wrappers in
internal/io

**Decision**: System path deny-list as safety net, not security boundary

**Rationale**: ctx paths are internally constructed from config constants. The
deny-list catches agent hallucinations (writing to /etc), not adversarial input.
Public security docs would imply a threat model that does not exist.

**Consequence**: internal/io/doc.go documents limitations honestly for
contributors. No user-facing security docs. The deny-list is a modicum of
protection, not a promise.

---

## [2026-03-14-093748] Config-driven freshness check with per-file review URLs

**Status**: Accepted

**Context**: Building a hook to warn when technology-dependent constants go
stale. Initially hardcoded the file list and Anthropic docs URL in the binary,
but this only worked inside the ctx repo and assumed all projects care about
Anthropic docs.

**Decision**: Config-driven freshness check with per-file review URLs

**Rationale**: Making the file list and review URLs configurable via .ctxrc
freshness_files means any project can opt in. Per-file review_url avoids
special-casing by project name — ctx sets Anthropic docs, other projects set
their own vendor links or omit it entirely.

**Consequence**: The hook is a no-op by default (opt-in). ctx's own .ctxrc
carries the tracked files. All nudge text goes through assets/text.yaml for
localization. No project detection logic needed.

---

## [2026-03-13-223111] Delete ctx-context-monitor skill — hook output is self-sufficient

**Status**: Accepted

**Context**: The skill documented how to relay context window warnings, but the
hook message already includes IMPORTANT: Relay this context window warning to
the user VERBATIM which agents follow without the skill.

**Decision**: Delete ctx-context-monitor skill — hook output is
self-sufficient

**Rationale**: No mechanism exists for hooks to trigger skills. The skill was
never loaded during sessions. Adding enforcement elsewhere would either be too
far back in context (playbook) or dilute the already-crisp hook message.

**Consequence**: One fewer skill to maintain. No behavioral change — agents
continue relaying warnings as before.

---



## [2026-03-12-133007] Rename ctx-map skill to ctx-architecture

**Status**: Accepted

**Context**: The name 'map' didn't convey the iterative, architectural nature of
the ritual

**Decision**: Rename ctx-map skill to ctx-architecture

**Rationale**: 'architecture' better describes surveying and evolving project
structure across sessions

**Consequence**: All cross-references updated across skills, docs, .context
files, and settings

---

---

## [2026-03-06-200306] Drop fatih/color dependency — Unicode symbols are sufficient for terminal output, color was redundant

**Status**: Accepted

**Context**: fatih/color was used in 32 files for green checkmarks, yellow
warnings, cyan headings, dim text

**Decision**: Drop fatih/color dependency — Unicode symbols are sufficient for
terminal output, color was redundant

**Rationale**: Every colored output already had a semantic symbol (✓, ⚠,
○) that conveyed the same meaning; color added visual noise in non-terminal
contexts (logs, pipes)

**Consequence**: Removed --no-color flag (only existed for color.NoColor); one
fewer external dependency; FlagNoColor retained in config for CLI compatibility

---




---

## [2026-03-05-205424] Gitignore .context/memory/ for this project

**Status**: Accepted

**Context**: Memory mirror contains copies of MEMORY.md which holds strategic
analysis and session notes

**Decision**: Gitignore .context/memory/ for this project

**Rationale**: Strategic content should not be in git history. Docs updated to
say 'often git-tracked' for the general recommendation — this project is the
exception.

**Consequence**: Mirror and archives are local-only for this project. Other
projects can still track them. Sync and drift detection work the same way
regardless.

---



## [2026-03-04-105238] Interface-based GraphBuilder for multi-ecosystem ctx deps

**Status**: Accepted

**Context**: P-1.3 questioned whether non-Go dependency support would introduce
bloat and whether a semantic approach was better

**Decision**: Interface-based GraphBuilder for multi-ecosystem ctx deps

**Rationale**: The output pipeline (map[string][]string to Mermaid/table/JSON)
was already language-agnostic. Each ecosystem builder is ~40 lines — this is
finishing what was started, not bloat. Static manifest parsing (no external
tools for Node/Python) keeps dependencies minimal.

**Consequence**: ctx deps now auto-detects Go, Node.js, Python, Rust. --type
flag overrides detection. ctx-architecture skill works across ecosystems without
changes.

---

## [2026-03-02-165038] Billing threshold piggybacks on check-context-size, not heartbeat

**Status**: Accepted

**Context**: User wanted a configurable token-count nudge for billing awareness
(Claude Pro 1M context, extra cost after 200k). Heartbeat produces zero stdout
and can't relay to user.

**Decision**: Billing threshold piggybacks on check-context-size, not heartbeat

**Rationale**: check-context-size already reads tokens, has VERBATIM relay
working, and runs every prompt. Adding a third independent trigger there is
minimal code and follows established patterns.

**Consequence**: New .ctxrc field billing_token_warn (default 0 = disabled).
One-shot per session via billing-warned-{sessionID} state file.
Template-overridable via check-context-size/billing.txt.

---



## [2026-03-01-222733] PersistentPreRunE init guard with three-level exemption

**Status**: Accepted

**Context**: ctx commands handled missing .context/ inconsistently — some
caught errors, some got confusing file-not-found messages, some produced empty
output

**Decision**: PersistentPreRunE init guard with three-level exemption

**Rationale**: Single PersistentPreRunE on root command gives one clear error.
Three-level exemption (hidden commands, annotated commands, grouping commands)
covers all edge cases without per-command boilerplate

**Consequence**: Boundary violation now returns an error instead of os.Exit(1),
making it testable. The subprocess-based boundary test was simplified to a
direct error assertion

---

---

## [2026-03-01-112544] Heartbeat token telemetry: conditional fields, not always-present

**Status**: Accepted

**Context**: Adding tokens, context_window, usage_pct to heartbeat payloads.
First prompt of a session has no JSONL usage data yet.

**Decision**: Heartbeat token telemetry: conditional fields, not always-present

**Rationale**: Token fields are only included in the template ref when tokens >
0. This avoids misleading pct=0% on the first heartbeat and keeps payloads clean
for receivers that filter on field presence.

**Consequence**: Webhook consumers must handle heartbeats both with and without
token fields. The message string also varies (with/without tokens=N pct=N%
suffix).

---

---

## [2026-03-01-090124] Promote 6 private skills to bundled plugin skills; keep 7 project-local

**Status**: Accepted

**Context**: Reviewed all 13 _ctx-* private skills to determine which are
universally useful for any ctx user vs specific to the ctx codebase or personal
infra.

**Decision**: Promote 6 private skills to bundled plugin skills; keep 7
project-local

**Rationale**: Promote if the skill benefits any ctx-powered project without
project-specific hardcoding. Keep private if it references this repo's Go
internals, personal infra, or language-specific tooling. Promote list: _ctx-spec
(generic scaffolding), _ctx-brainstorm (design facilitation), _ctx-verify (claim
verification), _ctx-skill-create (skill authoring), _ctx-link-check (doc link
audit), _ctx-permission-sanitize (Claude Code permissions audit). Keep list:
_ctx-audit (Go/ctx checks), _ctx-qa (Go Makefile), _ctx-backup (SMB infra),
_ctx-release/_ctx-release-notes (ctx release workflow), _ctx-update-docs (ctx
package mapping), _ctx-absorb (borderline, revisit later).

**Consequence**: Six skills move from .claude/skills/ to
internal/assets/claude/skills/ and become available to all ctx users via ctx
init. Cross-references between skills need updating (e.g., /_ctx-brainstorm
becomes /ctx-brainstorm). The seven remaining private skills stay project-local.

---

## [2026-02-27-230718] Context window detection: JSONL-first fallback order

**Status**: Accepted

**Context**: check-context-size defaults to 200k but user runs 1M-context model,
causing false 110% warnings. JSONL contains the model name which maps to actual
window size.

**Decision**: Context window detection: JSONL-first fallback order

**Rationale**: effective_window = detect_from_jsonl(model) ??
ctxrc.context_window ?? 200_000. JSONL is ground truth (reflects actual model in
use); ctxrc is fallback for first-hook-of-session or unknown models; 200k is
safe last resort. Having ctxrc override JSONL would artificially restrict the
check when a user forgets to update their config after switching models.

**Consequence**: Most users get correct window automatically. ctxrc
context_window becomes a fallback, not an override. Task exists for
implementation.

---



---

## [2026-02-26-100002] ctx init and CLAUDE.md handling (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-01-20)

- `ctx init` handles CLAUDE.md intelligently: creates if missing, backs up and
  offers merge if existing, uses marker comment for idempotency. The `--merge`
  flag enables non-interactive append.
- `ctx init` always generates `.claude/hooks/` alongside `.context/` with no
  flag needed. Other AI tools ignore `.claude/`; Claude Code users get seamless
  zero-config experience.
- Core tool stays generic and tool-agnostic, with optional Claude Code
  enhancements via `.claude/hooks/`. Other AI tools can be supported similarly
  (`ctx hook cursor`, etc.).

---

## [2026-02-26-100004] Task and knowledge management (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-27 to 2026-02-18)

- Tasks must include explicit deliverables, not just implementation steps.
  Parent tasks define WHAT the user gets; subtasks decompose HOW to build it.
  Without explicit deliverables, AI optimizes for checking boxes.
- Use reverse-chronological order (newest first) for DECISIONS.md and
  LEARNINGS.md. Ensures most recent items are read first regardless of token
  budget.
- Add quick reference index to DECISIONS.md: compact table at top allows
  scanning; agents can grep for full timestamp to jump to entry. Auto-updated on
  `ctx add decision`.
- Knowledge scaling via archive path for decisions and learnings: follow the
  task archive pattern, move old entries to `.context/archive/`, extend `ctx
  compact --archive` to cover all three file types.

---

## [2026-02-26-100005] Agent autonomy and separation of concerns (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-01-21 to 2026-01-28)

- Removed AGENTS.md from project root. Consolidated on CLAUDE.md (auto-loaded) +
  .context/AGENT_PLAYBOOK.md as the canonical agent instruction path. Projects
  using ctx should not create AGENTS.md.
- ~~Separate orchestrator directive from agent tasks~~ (superseded 2026-03-25:
  IMPLEMENTATION_PLAN.md removed — TASKS.md is the single source of truth for
  work items, AGENT_PLAYBOOK.md covers agent behavior).
- No custom UI -- IDE is the interface. UI is a liability; IDEs already excel at
  file browsing, search, markdown editing, and git integration. Focus CLI
  efforts on good markdown output.

---

## [2026-02-26-100006] Security and permissions (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-01-21 to 2026-02-24)

- Keep CONSTITUTION.md minimal: only truly inviolable rules (security,
  correctness, process invariants). Style preferences go in CONVENTIONS.md.
  Overly strict constitution gets ignored.
- Centralize constants with semantic prefixes in `internal/config/config.go`:
  `Dir*` for directories, `File*` for paths, `Filename*` for names,
  `UpdateType*` for entry types. Single source of truth, compile-time typo
  checks.
- Hooks use `ctx` from PATH, not hardcoded absolute paths. Standard Unix
  practice; portable across machines/users. `ctx init` checks PATH availability
  before proceeding.
- Drop absolute-path-to-ctx regex from block-dangerous-commands shell script.
  The block-non-path-ctx Go subcommand already covers this with better patterns;
  duplicating creates two sources of truth.

---

## [2026-02-27-002831] Webhook and notification design (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-02-22 to 2026-02-26)

- **Session attribution**: All webhook payloads must include session_id. Reading
  it from stdin costs nothing and enables multi-agent diagnostics. All run
  functions take stdin parameter; tests use createTempStdin.
- **Opt-in events**: Notify events are opt-in, not opt-out. EventAllowed returns
  false for nil/empty event lists. The correct default for notifications is
  silence. `ctx notify test` bypasses the filter as a special case.
- **Shared encryption key**: Webhook URLs encrypted with the shared .ctx.key
  (AES-256-GCM), not a dedicated key. One key, one gitignore entry, one rotation
  cycle. Notify is a peer of scratchpad — both store user secrets encrypted at
  rest.

---

## [2026-02-11] Remove .context/sessions/ storage layer and ctx session command

**Status**: Accepted

**Context**: The session/recall/journal system had three overlapping storage
layers: `~/.claude/projects/` (raw JSONL transcripts, owned by Claude Code),
`.context/sessions/` (JSONL copies + context snapshots), and `.context/journal/`
(enriched markdown from `ctx recall import`). The recall pipeline reads directly
from `~/.claude/projects/`, making `.context/sessions/` a dead-end write sink
that nothing reads from. The auto-save hook copied transcripts to a directory
nobody consumed. The `ctx session save` command created context snapshots that
git already provides through version history. This was ~15 Go source files, a
shell hook, ~20 config constants, and 30+ doc references supporting
infrastructure with no consumers.

**Decision**: Remove `.context/sessions/` entirely. Two stores remain: raw
transcripts (global, tool-owned in `~/.claude/projects/`) and enriched journal
(project-local in `.context/journal/`).

**Rationale**: Dead-end write sinks waste code surface, maintenance effort, and
user attention. The recall pipeline already proved that reading directly from
`~/.claude/projects/` is sufficient. Context snapshots are redundant with git
history. Removing the middle layer simplifies the architecture from three stores
to two, eliminates an entire CLI command tree (`ctx session`), and removes a
shell hook that fired on every session end.

**Consequence**: Deleted `internal/cli/session/` (15 files), removed auto-save
hook, removed `--auto-save` from watch, removed pre-compact auto-save from
compact, removed `/ctx-save` skill, updated ~45 documentation files. Four
earlier decisions superseded (SessionEnd hook, Auto-Save Before Compact, Session
Filename Format, Two-Tier Persistence Model). Users who want session history use
`ctx journal source`/`ctx journal import` instead.

---


*Module-specific, already-shipped, and historical decisions:
[decisions-reference.md](decisions-reference.md)*

---

## [2026-04-25-014704] Use t.Setenv for subprocess env in tests, not append(os.Environ(), ...)

**Status**: Accepted

**Context**: TestBinaryIntegration spawns subprocesses; the prior helper did
append(os.Environ(), CTX_DIR=...) to override the developer-shell value. Wrong
abstraction.

**Decision**: Use t.Setenv for subprocess env in tests, not append(os.Environ(),
...)

**Rationale**: t.Setenv mutates the live process env, exec.Cmd with nil Env
inherits it, and cleanup is automatic at test end. One line replaces the helper.

**Consequence**: Helper deleted, six call sites simplified, no env-dedup logic
to maintain. Pattern reusable for other subprocess tests.

---
