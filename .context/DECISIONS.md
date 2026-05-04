# Decisions

<!-- INDEX:START -->
| Date | Decision |
|----|--------|
| 2026-04-16 | Deprecate and remove ctx backup |
| 2026-04-14 | doc.go quality floor: behavior-grounded, ~25-100 body lines, related-packages section required |
| 2026-04-14 | Bootstrap stays under ctx system bootstrap (reverted experimental top-level promotion) |
| 2026-04-14 | Title Case style for docs is AP-leaning with explicit ambiguity carve-outs |
| 2026-04-13 | Walk boundary uses git as a hint, not a requirement |
| 2026-04-11 | Journal stays local; LEARNINGS.md is the shareable layer |
| 2026-04-11 | `Entry.Author` is server-authoritative, not client-authoritative |
| 2026-04-09 | Architecture skill pipeline is a triad not a quartet |
| 2026-04-08 | Remove #done tag convention, simplify task archival |
| 2026-04-06 | Use hook relay for session provenance instead of JSONL parsing or env vars |
| 2026-04-04 | TestNoMagicStrings and TestNoMagicValues no longer exempt const/var definitions outside config/ |
| 2026-04-04 | String-typed enums belong in config/, not domain packages |
| 2026-04-03 | Output functions belong in write/ (consolidated) |
| 2026-04-03 | YAML text externalization pipeline (consolidated) |
| 2026-04-03 | Package taxonomy and code placement (consolidated) |
| 2026-04-03 | Eager init over lazy loading (consolidated) |
| 2026-04-03 | Pure logic separation of concerns (consolidated) |
| 2026-04-03 | config/ explosion is correct — fix is documentation, not restructuring |
| 2026-04-01 | IRC to Discord as primary community channel |
| 2026-04-01 | AST audit tests live in internal/audit/, one file per check |
| 2026-04-01 | Split assets/hooks/ into assets/integrations/ + assets/hooks/messages/ |
| 2026-04-01 | Rename ctx hook → ctx setup to disambiguate from the hook system |
| 2026-03-31 | Split log into log/event and log/warn to break import cycles |
| 2026-03-31 | Context-load-gate injects only CONSTITUTION and AGENT_PLAYBOOK_GATE, not full ReadOrder |
| 2026-03-31 | Spec signal words and nudge threshold are user-configurable via .ctxrc |
| 2026-03-30 | Flags-not-subcommands for journal source: list and show are view modes on a noun, not independent entities |
| 2026-03-30 | Journal consumed recall — recall CLI package deleted |
| 2026-03-30 | Classify rules are user-configurable via .ctxrc |
| 2026-03-25 | Architecture analysis and enrichment are separate skills — constraint is the feature |
| 2026-03-25 | Companion tools documented as optional MCP enhancements with runtime check |
| 2026-03-25 | Prompt templates removed — skills are the single agent instruction mechanism |
| 2026-03-24 | Write-once baseline with explicit end-consolidation for consolidation lifecycle |
| 2026-03-23 | Pre/pre HTML tags promoted to shared constants in config/marker |
| 2026-03-22 | Output functions belong in write/, never in core/ or cmd/ |
| 2026-03-20 | Shared formatting utilities belong in internal/format |
| 2026-03-20 | Go-YAML linkage check added to lint-drift as check 5 |
| 2026-03-18 | Singular command names for all CLI entities |
| 2026-03-17 | Pre-compute-then-print for write package output blocks |
| 2026-03-16 | Resource name constants in config/mcp/resource, mapping in server/resource |
| 2026-03-16 | Rename --consequences flag to --consequence for singular consistency |
| 2026-03-14 | Error package taxonomy: 22 domain files replace monolithic errors.go |
| 2026-03-14 | Session prefixes are parser vocabulary, not i18n text |
| 2026-03-14 | System path deny-list as safety net, not security boundary |
| 2026-03-14 | Config-driven freshness check with per-file review URLs |
| 2026-03-13 | Delete ctx-context-monitor skill — hook output is self-sufficient |
| 2026-03-13 | build target depends on sync-why to prevent embedded doc drift |
| 2026-03-12 | Recommend companion RAGs as peer MCP servers not bridge through ctx |
| 2026-03-12 | Rename ctx-map skill to ctx-architecture |
| 2026-03-07 | Use composite directory path constants for multi-segment paths |
| 2026-03-06 | Drop fatih/color dependency — Unicode symbols are sufficient for terminal output, color was redundant |
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
| 2026-04-25 | Use t.Setenv for subprocess env in tests, not append(os.Environ(), ...) |
| 2026-04-25 | Tighten state.Dir / rc.ContextDir to (string, error) with sentinel errors |
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

## [2026-04-13-153617] Walk boundary uses git as a hint, not a requirement

**Status**: Accepted

**Context**: ctx init failed when a non-ctx-initialized repo lived inside a
ctx-initialized parent workspace. walkForContextDir walked up and found the
parent's .context, then the boundary check rejected it. We considered
project-marker heuristics (go.mod, package.json) and making git mandatory.

**Decision**: Walk boundary uses git as a hint, not a requirement

**Rationale**: Project markers are unreliable (e.g. package.json for customer
shipments, Haskell projects have no common marker). Making git mandatory breaks
ctx's 'git recommended but not required' stance. Git-as-hint resolves the bug
without new dependencies: walk finds candidate, validate against git root,
discard if outside; fall back to CWD when no git is found.

**Consequence**: walkForContextDir now consults findGitRoot to anchor ancestor
.context candidates. Monorepos, submodules, and nested workspaces resolve
correctly. No-git projects still work via CWD fallback.

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

## [2026-04-04-025755] TestNoMagicStrings and TestNoMagicValues no longer exempt const/var definitions outside config/

**Status**: Accepted

**Context**: The isConstDef/isVarDef blanket exemption masked 156+ string and 7
numeric constants in the wrong package

**Decision**: TestNoMagicStrings and TestNoMagicValues no longer exempt
const/var definitions outside config/

**Rationale**: Const definitions outside config/ are magic values in the wrong
place — naming them does not fix the structural problem

**Consequence**: All new code with string/numeric constants outside config/
fails these tests immediately

---

## [2026-04-04-025746] String-typed enums belong in config/, not domain packages

**Status**: Accepted

**Context**: Debated whether type IssueType string with const values belongs in
domain or config. The string value is the same regardless of type annotation.

**Decision**: String-typed enums belong in config/, not domain packages

**Rationale**: Types without behavior belong in config. Promote to entity/ only
when methods/interfaces appear.

**Consequence**: All type Foo string + const blocks outside config/ are now
caught by TestNoMagicStrings.

---

## [2026-04-03-180000] Output functions belong in write/ (consolidated)

**Status**: Accepted

**Consolidated from**: 2 entries (2026-03-21 to 2026-03-22)

**Decision**: Output functions belong in write/, logic and types in core/,
orchestration in cmd/

**Rationale**: The write/ taxonomy is flat by domain — each CLI feature gets
its own write/ package. core/ owns domain logic and types. cmd/ owns Cobra
orchestration. Functions that call cmd.Print/Println/Printf belong in write/.
core/ never imports cobra for output purposes.

**Consequence**: All new CLI output must go through a write/ package. No
cmd.Print* calls in internal/cli/ outside of internal/write/.

---

## [2026-04-03-180000] YAML text externalization pipeline (consolidated)

**Status**: Accepted

**Consolidated from**: 5 entries (2026-03-06 to 2026-04-03)

**Decision**: All user-facing text externalized to embedded YAML domain files,
justified by agent legibility and drift prevention — not i18n

**Rationale**: The real justification is agent legibility (named DescKey
constants as traversable graphs) and drift prevention (TestDescKeyYAMLLinkage
catches orphans mechanically). i18n is a free downstream consequence. The
exhaustive test verifies all constants resolve to non-empty YAML values — new
keys are automatically covered.

**Consequence**: commands.yaml split into 4 domain files (commands, flags, text,
examples) loaded via dedicated loaders. text.yaml split into 6 domain files
loaded via loadYAMLDir. The 3-file ceremony (DescKey + YAML + write/err
function) is the cost of agent-legible, drift-proof output.

---

## [2026-04-03-180000] Package taxonomy and code placement (consolidated)

**Status**: Accepted

**Consolidated from**: 3 entries (2026-03-06 to 2026-03-13)

**Decision**: Three-zone taxonomy: cmd/ for Cobra wiring (cmd.go + run.go),
core/ for logic and types, assets/ for templates and user-facing text. config/
for structural constants only.

**Rationale**: Taxonomical symmetry makes navigation instant and agent-friendly.
Domain types that multiple packages consume belong in domain packages
(internal/entry), not CLI subpackages. Templates and user-facing text live in
assets/ for i18n readiness; structural constants (paths, limits, regexes) stay
in config/.

**Consequence**: Every CLI package has the same predictable shape. Shared entry
types live in internal/entry. Template files (tpl_*.go) moved from config/ to
assets/. 474 files changed in initial restructuring.

---

## [2026-04-03-180000] Eager init over lazy loading (consolidated)

**Status**: Accepted

**Consolidated from**: 2 entries (2026-03-16 to 2026-03-18)

**Decision**: Explicit Init() called eagerly at startup for static embedded data
and resource lookups, instead of per-accessor sync.Once or package-level init()

**Rationale**: Static embedded data is required at startup — sync.Once per
accessor is cargo cult. Package-level init() hides startup dependencies and
makes ordering unclear. Explicit Init() called from main.go / NewServer makes
the dependency visible and testable.

**Consequence**: Maps unexported, accessors are plain lookups. Tests call Init()
in TestMain. res.Init() called from NewServer before ToList(). No package-level
side effects, zero sync.Once in the lookup pipeline.

---

## [2026-04-03-180000] Pure logic separation of concerns (consolidated)

**Status**: Accepted

**Consolidated from**: 3 entries (2026-03-15 to 2026-03-23)

**Decision**: Pure-logic functions return data structs; callers own I/O, file
writes, and reporting. Function pointers in param structs replaced with text
keys.

**Rationale**: Pure logic with no I/O lets both MCP (JSON-RPC) and CLI (cobra)
callers control output independently. Methods that don't access receiver state
hide their true dependencies — make them free functions. If all callers of a
callback vary only by a string key, the callback is data in disguise.

**Consequence**: CompactContext returns CompactResult; callers iterate
FileUpdates. Server response helpers in server/out, prompt builders in
server/prompt. All cross-cutting param structs in entity are
function-pointer-free.

---

## [2026-04-03-133244] config/ explosion is correct — fix is documentation, not restructuring

**Status**: Accepted

**Context**: Architecture analysis flagged 60+ config sub-packages as a
bottleneck. Evaluation showed the alternative (8-10 domain packages) trades
granular imports for fat dependency units. Current structure gives zero internal
dependencies, surgical dependency tracking, and minimal recompile scope.

**Decision**: config/ explosion is correct — fix is documentation, not
restructuring

**Rationale**: Go's compilation unit is the package. Granular packages mean
precise dependency tracking. The developer experience cost (IDE noise, package
discovery) is real but solvable with a README decision tree, not restructuring.
Restructuring would be massive mechanical churn for cosmetic benefit.

**Consequence**: config/README.md written with organizational guide and decision
tree. No restructuring planned. embed/text/ file count will shrink naturally
when tpl/ migrates to text/template.

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

## [2026-04-01-074417] Split assets/hooks/ into assets/integrations/ + assets/hooks/messages/

**Status**: Accepted

**Context**: The directory mixed Copilot integration templates with hook message
templates

**Decision**: Split assets/hooks/ into assets/integrations/ +
assets/hooks/messages/

**Rationale**: Integration assets (Copilot instructions, AGENTS.md, CLI
scripts/skills) are not hooks. Hook messages ARE the hook system templates.

**Consequence**: integrations/ for tool integration assets, hooks/messages/ for
hook system templates. Embed directives and all config constants updated.

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

## [2026-03-31-182003] Context-load-gate injects only CONSTITUTION and AGENT_PLAYBOOK_GATE, not full ReadOrder

**Status**: Accepted

**Context**: Force-loading ~14k tokens of context files (8 files) every session
diluted attention without proportional value. CLAUDE.md already instructs agents
to read full context files on-demand. Behavioral prose in force-loaded content
was routinely skipped.

**Decision**: Context-load-gate injects only CONSTITUTION and
AGENT_PLAYBOOK_GATE, not full ReadOrder

**Rationale**: Hard rules (CONSTITUTION) must be present before any action.
Distilled directives (gate file) provide actionable session-start guidance in
~2k tokens. Full playbook, conventions, architecture, decisions, learnings are
pulled on-demand when task context requires them.

**Consequence**: New AGENT_PLAYBOOK_GATE.md file must stay in sync with
AGENT_PLAYBOOK.md. HTML comment cross-reference added to playbook header for
contributor discoverability.

---

## [2026-03-31-005113] Spec signal words and nudge threshold are user-configurable via .ctxrc

**Status**: Accepted

**Context**: Initially hardcoded signal words and 150-char threshold in run.go.
User pointed out these are localizable vocabulary, following the
session_prefixes / classify_rules pattern

**Decision**: Spec signal words and nudge threshold are user-configurable via
.ctxrc

**Rationale**: Signal words are language-dependent and project-dependent — a
Spanish-speaking user or a non-Go project would have different signal terms

**Consequence**: Added spec_signal_words and spec_nudge_min_len to CtxRC struct,
rc accessors with defaults in config/entry, JSON schema updated

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

## [2026-03-30-003745] Classify rules are user-configurable via .ctxrc

**Status**: Accepted

**Context**: Memory entry classification used hardcoded keyword rules that could
not be customized

**Decision**: Classify rules are user-configurable via .ctxrc

**Rationale**: Users may work in domains where the default keywords do not match
(non-English, specialized terminology). Same pattern as session_prefixes.

**Consequence**: classify_rules in .ctxrc overrides defaults; schema updated;
rc.ClassifyRules() accessor with fallback to config/memory.DefaultClassifyRules

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

## [2026-03-25-173337] Companion tools documented as optional MCP enhancements with runtime check

**Status**: Accepted

**Context**: Gemini Search and GitNexus improve skills but no docs mentioned
them and no code checked their availability

**Decision**: Companion tools documented as optional MCP enhancements with
runtime check

**Rationale**: Users should know what tools enhance their workflow without being
forced to install them. Suppressible via .ctxrc for users who don't want them.

**Consequence**: /ctx-remember smoke-tests MCPs at session start.
companion_check: false suppresses.

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

## [2026-03-23-165612] Pre/pre HTML tags promoted to shared constants in config/marker

**Status**: Accepted

**Context**: Two packages (normalize and format) used hardcoded pre strings
independently

**Decision**: Pre/pre HTML tags promoted to shared constants in config/marker

**Rationale**: Cross-package magic strings belong in config constants per
CONVENTIONS.md

**Consequence**: marker.TagPre and marker.TagPreClose are the canonical
references; package-local constants deleted

---

## [2026-03-22-084316] Output functions belong in write/, never in core/ or cmd/

**Status**: Accepted

**Context**: System write migration revealed that cmd.Print* calls scattered
across core/ and cmd/ packages prevented localization and violated separation of
concerns

**Decision**: Output functions belong in write/, never in core/ or cmd/

**Rationale**: The write/ taxonomy is flat by domain — each CLI feature gets
its own write/ package. core/ owns logic and types, cmd/ owns orchestration,
write/ owns all output.

**Consequence**: All new CLI output must go through a write/ package. No
cmd.Print* calls in internal/cli/ outside of internal/write/.

---

## [2026-03-20-232506] Shared formatting utilities belong in internal/format

**Status**: Accepted

**Context**: Pluralize, Duration, DurationAgo, and TruncateFirstLine were
duplicated across memory/core, change/core, and other CLI packages

**Decision**: Shared formatting utilities belong in internal/format

**Rationale**: internal/format already existed with TimeAgo and Number
formatters. Centralizing prevents duplication and matches the convention that
domain-agnostic utilities live in shared packages, not CLI subpackages

**Consequence**: CLI packages import internal/format instead of defining local
helpers. Local copies deleted.

---

## [2026-03-20-160103] Go-YAML linkage check added to lint-drift as check 5

**Status**: Accepted

**Context**: Prior refactoring sessions left broken and orphan linkages between
Go DescKey constants and YAML entries that caused silent runtime failures

**Decision**: Go-YAML linkage check added to lint-drift as check 5

**Rationale**: Shell-based grep+comm approach fits the existing lint-drift
pattern, runs at CI time, and is simpler than programmatic Go AST parsing

**Consequence**: CI-time check catches orphans in both directions plus
cross-namespace duplicates, preventing recurrence

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

## [2026-03-17-105627] Pre-compute-then-print for write package output blocks

**Status**: Accepted

**Context**: Audit of internal/write/ found 337 Println calls across 160
functions. Asked whether text/template or single format strings would clean up
multi-Println functions like InfoLoopGenerated.

**Decision**: Pre-compute-then-print for write package output blocks

**Rationale**: text/template trades compile-time safety for runtime errors and
only 38 of 160 functions benefit from consolidation. fmt.Sprintf with
pre-computed conditional args handles all cases without new dependencies.
Loop-based functions stay imperative.

**Consequence**: Functions with 4+ Printlns pre-compute conditionals into
strings, then emit one cmd.Println with a multiline block template. Per-line
Tpl* constants replaced with TplXxxBlock. Trivial (1-3 line) and loop-based
functions excluded.

---

## [2026-03-16-104142] Resource name constants in config/mcp/resource, mapping in server/resource

**Status**: Accepted

**Context**: MCP resource handler had string literals scattered through
handle_resource.go and rebuilt the resource list on every call

**Decision**: Resource name constants in config/mcp/resource, mapping in
server/resource

**Rationale**: Constants follow the same pattern as config/mcp/tool. Mapping
stays in server/resource because it bridges config constants with assets text
(too many cross-cutting deps for a config package). Resource list and URI lookup
are pre-built once at server init.

**Consequence**: URI-to-file lookup is O(1) via pre-built map; resource list
built once in NewServer, not per request; no string literals in handler code

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

## [2026-03-14-180905] Error package taxonomy: 22 domain files replace monolithic errors.go

**Status**: Accepted

**Context**: internal/err/errors.go was 1995 lines with 188 functions in one
file

**Decision**: Error package taxonomy: 22 domain files replace monolithic
errors.go

**Rationale**: Convention requires files named by responsibility, not junk
drawers; domain grouping makes it possible to find error constructors by domain

**Consequence**: 22 files (backup, config, crypto, date, fs, git, hook, init,
journal, memory, notify, pad, parser, prompt, recall, reminder, session, site,
skill, state, task, validation); errors.go deleted

---

## [2026-03-14-131152] Session prefixes are parser vocabulary, not i18n text

**Status**: Accepted

**Context**: Markdown session parser had hardcoded Session:/Oturum: pair in
text.yaml as session_prefix/session_prefix_alt — didn't scale beyond two
languages

**Decision**: Session prefixes are parser vocabulary, not i18n text

**Rationale**: Session header prefixes are recognition patterns for parsing, not
user-facing interface strings. Separating content recognition from interface
language lets users parse multilingual session files without code changes.
Single-language default (Session:) avoids implicit favoritism.

**Consequence**: Prefixes moved to .ctxrc session_prefixes list. text.yaml
entries and embed.go constants removed. Parser reads from rc.SessionPrefixes()
with fallback to config/parser.DefaultSessionPrefixes. Users extend via .ctxrc.

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

## [2026-03-13-151955] build target depends on sync-why to prevent embedded doc drift

**Status**: Accepted

**Context**: assets/why/ files had silently drifted from their docs/ sources

**Decision**: build target depends on sync-why to prevent embedded doc drift

**Rationale**: Derived assets that are not in the build dependency chain will
drift — the only reliable enforcement is making the build fail without sync

**Consequence**: Every make build now copies docs into assets before compiling

---

## [2026-03-12-133007] Recommend companion RAGs as peer MCP servers not bridge through ctx

**Status**: Accepted

**Context**: Explored whether ctx should proxy RAG queries or integrate a RAG
directly

**Decision**: Recommend companion RAGs as peer MCP servers not bridge through
ctx

**Rationale**: MCP is the composition layer — agents already compose multiple
servers. ctx is context, RAGs are intelligence. No bridging, no plugin system,
no schema abstraction

**Consequence**: Spec created at ideas/spec-companion-intelligence.md; future
work is documentation and UX only

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

## [2026-03-07-221155] Use composite directory path constants for multi-segment paths

**Status**: Accepted

**Context**: Needed a constant for hooks/messages path used in message.go and
message_cmd.go

**Decision**: Use composite directory path constants for multi-segment paths

**Rationale**: Matches existing pattern of DirClaudeHooks = '.claude/hooks' —
keeps filepath.Join calls cleaner and avoids scattering path segments

**Consequence**: New multi-segment directory paths should be single constants
(e.g. DirHooksMessages, DirMemoryArchive) rather than joined from individual
segment constants

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

## [2026-03-06-141507] PR #27 (MCP server) meets v0.1 spec requirements — merge-ready pending 3 compliance fixes

**Status**: Accepted

**Context**: Reviewed PR against specs/mcp-server.md; all 7 action items
addressed, CI fails on 3 mechanical compliance issues

**Decision**: PR #27 (MCP server) meets v0.1 spec requirements — merge-ready
pending 3 compliance fixes

**Rationale**: All spec requirements met; CI failures are trivial and low-risk;
keeping PR open risks merge conflicts during active refactoring

**Consequence**: Merge and fix compliance issues in follow-up commit on main

---

## [2026-03-06-184816] Skills stay CLI-based; MCP Prompts are the protocol equivalent

**Status**: Accepted

**Context**: Question arose whether skills should switch from ctx CLI (Bash) to
MCP tool calls once the MCP server ships

**Decision**: Skills stay CLI-based; MCP Prompts are the protocol equivalent

**Rationale**: CLI is always available (PATH prerequisite); MCP requires
optional configuration. Hooks will always be CLI (shell commands). Two access
patterns in the same tool is gratuitous complexity.

**Consequence**: Skills call CLI. MCP Prompts call MCP Tools. Hooks call CLI.
Clean layer separation; no replacement, only parallel access paths.

---

## [2026-03-06-184812] Peer MCP model for external tool integration

**Status**: Accepted

**Context**: Evaluated three integration models (orchestrator, peer, hub) for
how ctx relates to GitNexus and context-mode

**Decision**: Peer MCP model for external tool integration

**Rationale**: Peer model (side-by-side MCP servers, each queried independently
by the agent) respects ctx's markdown-on-filesystem invariant and avoids
coupling. ctx provides behavioral scaffolding; external tools provide their
specialties.

**Consequence**: ctx MCP Prompts can reference external tools by convention
without tight coupling. No plugin registry needed.

---

## [2026-03-06-050132] Create internal/parse for shared text-to-typed-value conversions

**Status**: Accepted

**Context**: parseDate with 2006-01-02 duplicated in 5+ files; needed a home
that is not internal/utils or internal/strings (collides with stdlib)

**Decision**: Create internal/parse for shared text-to-typed-value conversions

**Rationale**: internal/parse scopes to convert text to typed values without
becoming a junk drawer. Name invites sibling functions (duration, identifier
parsing) naturally.

**Consequence**: parse.Date() is the first function; config.DateFormat holds the
layout constant. Other time.Parse callers can migrate incrementally.

---

## [2026-03-06-050131] Centralize errors in internal/err, not per-package err.go files

**Status**: Accepted

**Context**: Duplicate error constructors across 5+ CLI packages; agents copying
the pattern when they see a local err.go

**Decision**: Centralize errors in internal/err, not per-package err.go files

**Rationale**: Single location makes duplicates visible, enables future sentinel
errors, and prevents broken-window accumulation

**Consequence**: All CLI err.go files migrated and deleted. New errors go to
internal/err/errors.go exclusively.

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

## [2026-03-05-042154] Memory bridge design: three-phase architecture with hook nudge + on-demand

**Status**: Accepted

**Context**: Brainstormed how to bridge Claude Code MEMORY.md with ctx
structured context files

**Decision**: Memory bridge design: three-phase architecture with hook nudge +
on-demand

**Rationale**: Hook nudge + on-demand gives user choice and freedom. Wrap-up is
the publish trigger, never commit (footgun). Heuristic classification for v1, no
LLM. Marker-based merge for bidirectional conflict. Mirror is git-tracked +
timestamped archives. Foundation spec delivers sync/status/diff/hook; import and
publish are future phases.

**Consequence**: Foundation spec in specs/memory-bridge.md, import/publish specs
deferred to ideas/. Tasked out as S-0.1.1 through S-0.1.10 in ideas/TASKS.md.

---

## [2026-03-05-023937] Revised strategic analysis: blog-first execution order, bidirectional sync as top-level section

**Status**: Accepted

**Context**: Editorial review of ideas/claude-memory-strategic-analysis.md
surfaced six structural weaknesses in competitive positioning

**Decision**: Revised strategic analysis: blog-first execution order,
bidirectional sync as top-level section

**Rationale**: 200-line cap is fragile differentiator (demoted); org-scoped
memory is the real threat (elevated to HIGH); model agnosticism is premature
(parked with trigger condition); bidirectional sync is the most underweighted
insight (promoted); narrative shapes categories before implementation does (blog
first)

**Consequence**: Execution order is now S-3 (blog) -> S-0 -> S-1 -> S-2.
Strategic doc restructured from 9 to 10 sections. Blog post shipped as first
deliverable.

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

## [2026-03-02-123611] Replace auto-migration with stderr warning for legacy keys

**Status**: Accepted

**Context**: Auto-migration code existed for promoting keys from
~/.local/ctx/keys/ and .context/.ctx.key to ~/.ctx/.ctx.key. Userbase is small
and this is alpha — no need to bloat the codebase.

**Decision**: Replace auto-migration with stderr warning for legacy keys

**Rationale**: Warn-only is simpler, avoids silent file operations, and puts the
user in control. Migration instructions in docs are sufficient for the small
userbase.

**Consequence**: MigrateKeyFile() now only warns on stderr. promoteToGlobal()
helper deleted. Tests verify keys are not moved.

---

## [2026-03-02-005213] Consolidate all session state to .context/state/

**Status**: Accepted

**Context**: Session-scoped state (cooldown tombstones, pause markers, daily
throttle markers) was split between /tmp (via secureTempDir()) and
.context/state/ for project-scoped state

**Decision**: Consolidate all session state to .context/state/

**Rationale**: Single location simplifies mental model, eliminates duplicated
secureTempDir() in two packages, removes the cleanup-tmp SessionEnd hook
entirely. .context/state/ is already gitignored and project-scoped.

**Consequence**: All 18 callers updated. Tests switch from XDG_RUNTIME_DIR
mocking to CTX_DIR + rc.Reset(). Hook lifecycle drops from 4 events to 3
(SessionEnd removed).

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

## [2026-03-01-161457] Global encryption key at ~/.ctx/.ctx.key

**Status**: Superseded by [2026-03-02] global key simplification

**Context**: Key stored next to ciphertext (.context/.ctx.key) was a security
antipattern and broke in worktrees. The slug-based per-project key system at
~/.local/ctx/keys/ was over-engineered for the common case (one user, one
machine, one key).

**Decision**: Single global key at ~/.ctx/.ctx.key. Project-local override via
.ctxrc key_path or .context/.ctx.key.

**Rationale**: One key per machine covers 99% of users. Per-project slug
filenames and three-tier resolution added complexity without clear benefit.
~/.ctx/ is the natural home (matches ~/.claude/ convention). Tilde expansion in
.ctxrc key_path fixes a standalone bug.

**Consequence**: Auto-migration promotes legacy keys (project-local,
~/.local/ctx/keys/) to ~/.ctx/.ctx.key. Deleted KeyDir(), ProjectKeySlug(),
ProjectKeyPath(). ResolveKeyPath simplified to two params. 15+ doc files
updated.

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

## [2026-03-01-092613] Hook log rotation: size-based with one previous generation, matching eventlog pattern

**Status**: Accepted

**Context**: .context/logs/ files grow unbounded (~200KB after one month);
needed a cap

**Decision**: Hook log rotation: size-based with one previous generation,
matching eventlog pattern

**Rationale**: Architectural symmetry with eventlog, O(1) size check vs O(n)
line counting, diagnostic logs don't need deep history (webhooks cover serious
setups)

**Consequence**: Each log file caps at ~2MB (current + .1). config.LogMaxBytes =
1MB, same as EventLogMaxBytes

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

## [2026-02-27-002830] Context injection architecture v2 (consolidated)

**Status**: Accepted

**Consolidated from**: 3 decisions (2026-02-26)

- **Diagram extraction**: ARCHITECTURE.md contained ~600 lines of ASCII/Mermaid
  diagrams (~12K tokens). Extracted to 5 architecture-dia-*.md files outside
  FileReadOrder. Agents get verbal summaries at session start; diagrams
  available on demand. Total injection dropped 53% (20K→9.5K tokens).
- **Auto-injection replaces directives**: Soft instructions have ~75-85%
  compliance ceiling because "don't apply judgment" is itself evaluated by
  judgment. The v2 context-load-gate injects content directly via
  `additionalContext` — agents never choose whether to comply. Injection
  strategy: CONSTITUTION, CONVENTIONS, ARCHITECTURE, AGENT_PLAYBOOK verbatim;
  DECISIONS, LEARNINGS index-only; TASKS mention-only. Total ~7,700 tokens. See:
  `specs/context-load-gate-v2.md`.
- **Imperative framing**: Advisory framing allowed agents to assess relevance
  and skip files. Imperative framing with unconditional compliance checkpoint
  removes the escape hatch. Verbatim relay is fallback safety net, not primary
  instruction.

---

## [2026-02-26-200001] .context/state/ directory for project-scoped runtime state

**Status**: Accepted

New gitignored directory under `context_dir` resolution for ephemeral
project-scoped state. Follows `.context/logs/` precedent — added to
`config.GitignoreEntries` and root `.gitignore`.

First use: injection oversize flag written by context-load-gate when injected
tokens exceed the configurable `injection_token_warn` threshold (`.ctxrc`,
default 15000). The check-context-size VERBATIM hook reads the flag and nudges
the user to run `/ctx-consolidate`.

See: `specs/injection-oversize-nudge.md`.

---

## [2026-02-26-100001] Hook and notification design (consolidated)

**Status**: Accepted

**Consolidated from**: 4 decisions (2026-02-12 to 2026-02-24)

- Tone down proactive content suggestion claims in docs rather than add more
  hooks. Already have 9 UserPromptSubmit hooks; adding another risks fatigue.
  Conversational prompting already works.
- Hook commands must use structured JSON output
  (hookSpecificOutput.additionalContext) instead of plain text, because Claude
  Code treats plain text as ignorable ambient context.
- Drop prompt-coach hook entirely: zero useful tips fired, output channel
  invisible to user, orphan temp file accumulation. The prompting guide already
  covers best practices.
- De-emphasize /ctx-journal-normalize from the default journal pipeline. The
  normalize skill is expensive and nondeterministic; programmatic normalization
  handles most cases. Skill remains available for targeted per-file use.

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

**Context**: TestBinaryIntegration spawns subprocesses; the prior helper did append(os.Environ(), CTX_DIR=...) to override the developer-shell value. Wrong abstraction.

**Decision**: Use t.Setenv for subprocess env in tests, not append(os.Environ(), ...)

**Rationale**: t.Setenv mutates the live process env, exec.Cmd with nil Env inherits it, and cleanup is automatic at test end. One line replaces the helper.

**Consequence**: Helper deleted, six call sites simplified, no env-dedup logic to maintain. Pattern reusable for other subprocess tests.

---

## [2026-04-25-014704] Tighten state.Dir / rc.ContextDir to (string, error) with sentinel errors

**Status**: Accepted

**Context**: Old single-return form returned ('', nil) when CTX_DIR was undeclared. Callers that filtered only on err != nil joined empty stateDir with relative names and wrote state files into CWD instead of .context/state/.

**Decision**: Tighten state.Dir / rc.ContextDir to (string, error) with sentinel errors

**Rationale**: Returning a sentinel ErrDirNotDeclared makes the empty-path case unrepresentable in a 'looks fine' branch. Forces every caller through the same explicit gate.

**Consequence**: All callers needed migration; tests had to declare CTX_DIR explicitly. In return, the filepath.Join('', rel) trap is closed by construction.
