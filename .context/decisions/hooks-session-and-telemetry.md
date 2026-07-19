# hooks-session-and-telemetry

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

