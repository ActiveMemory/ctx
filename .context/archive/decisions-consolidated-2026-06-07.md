# Archived Decisions (consolidated 2026-06-07)

Originals replaced by consolidated entries in DECISIONS.md.

## Group: Output belongs in write/ — taxonomy and emission style

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

## Group: Package taxonomy and shared-code placement

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
## [2026-05-17-181500] `entity.Sentinel` lives in `internal/entity/` because the cross-package-types audit treats `entity/` as the canonical home for shared types

**Status**: Accepted

**Context**: While converting the prior session's
`ErrMsg`-string-sentinel anti-pattern to typed-string sentinels
with lazy `desc.Text` resolution, the natural home for the
`Sentinel` type was a small shared helper used by every
`internal/err/<area>/` package. The first draft placed it at
`internal/err/sentinel/`, but `TestCrossPackageTypes` (which has
zero grandfathered violations and forbids weakening or
allowlist-bumping) flagged the cross-package usage with the hint
"consider entity/".

**Alternatives Considered**:
- Per-package sentinel type duplicated across 9 err packages.
  Pros: no cross-package type. Cons: 18 boilerplate declarations
  (type + Error method × 9) with doc comments; convention drift
  risk as the duplicated shape can diverge.
- Keep `internal/err/sentinel/` and add it to `typeExemptPackages`
  in the audit. Pros: semantic home matches the type's role
  (behavioral mixin for errors). Cons: the audit explicitly
  forbids exemption-list growth as the mechanism for new code;
  the test header says "If a test fails after your change, fix
  the code under test."
- Move `Sentinel` to `internal/entity/`. Pros: passes the audit
  without weakening; one shared declaration; consistent with
  every other cross-cutting type. Cons: `Sentinel` is a
  behavioral helper, not a domain data shape — semantically
  stretches `entity/`'s usual contents.

**Decision**: Place `Sentinel` in `internal/entity/sentinel.go`.

**Rationale**: The audit's rule is the project's hardline: every
cross-package type goes in `entity/`. The semantic stretch is
real but small, and writing exceptions to the audit is more
expensive long-term than absorbing a one-type semantic blur in
a package whose contract is already "things used cross-package."
Per-package duplication was rejected because the convention is
load-bearing — the next session that touches an err package
needs one obvious shape to copy, not a choice between 9 nearly
identical copies.

**Consequence**: `entity/` now houses a typed-string error
helper alongside its data shapes. Future readers landing in
`entity/` will find one file (`sentinel.go`) that doesn't
match the package's "data" theme; the doc comment on `Sentinel`
explains why. If `entity/` grows more behavioral helpers, the
package contract should be revisited; for now the precedent is
contained to this single type.

**Related**: LEARNINGS.md `[2026-05-17-180000] Sentinel errors
use typed zero-data structs with lazy desc.Text()` records the
shape itself.

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

## Group: Error handling: centralized in internal/err, domain-file taxonomy

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
## [2026-03-06-050131] Centralize errors in internal/err, not per-package err.go files

**Status**: Accepted

**Context**: Duplicate error constructors across 5+ CLI packages; agents copying
the pattern when they see a local err.go

**Decision**: Centralize errors in internal/err, not per-package err.go files

**Rationale**: Single location makes duplicates visible, enables future sentinel
errors, and prevents broken-window accumulation

**Consequence**: All CLI err.go files migrated and deleted. New errors go to
internal/err/errors.go exclusively.

## Group: config/ as constants home and the magic-value audit

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

## Group: YAML text externalization, init, and drift guards

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
## [2026-03-13-151955] build target depends on sync-why to prevent embedded doc drift

**Status**: Accepted

**Context**: assets/why/ files had silently drifted from their docs/ sources

**Decision**: build target depends on sync-why to prevent embedded doc drift

**Rationale**: Derived assets that are not in the build dependency chain will
drift — the only reliable enforcement is making the build fail without sync

**Consequence**: Every make build now copies docs into assets before compiling

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

## Group: CWD-anchored context model

## [2026-05-20-214812] Anchor ctx to CWD; drop activate, drop env-var resolver, drop all walks (proposed)

**Status**: Accepted

**Context**: Even after strict-CWD activate landed, eval $(ctx activate) remains an opaque per-shell ceremony. Two-channel resolution (env CTX_DIR + cwd) is the residual complexity; activate/deactivate exist only because of the env channel; the env channel exists to avoid the walk. With .context/ mandated as .git/'s sibling (CONSTITUTION require-git), if cwd must contain .context/ then both .context/ AND .git/ are in cwd — and every resolver across rc, gitmeta, and the activate commands collapses to os.Stat.

**Decision**: Anchor ctx to CWD; drop activate, drop env-var resolver, drop all walks (proposed)

**Rationale**: User counter to the agent's walk-to-.git/ proposal: the walk infrastructure (rc.ScanCandidates, gitmeta upward walk) is precisely what we want to delete; keeping ANY walk forces us to maintain two implementations. Mental model anchor matches zensical (zensical.toml), helm (Chart.yaml), terraform (.tf), Claude Code ($CLAUDE_PROJECT_DIR). Subdir convenience tax is a fixed per-shell cost (cd $(git rev-parse --show-toplevel)) for the user who knows their project root; agents pay no tax (cd is mechanical for them).

**Consequence**: Spec written at specs/cwd-anchored-context.md (314L); supersedes specs/activate-strict-cwd.md entirely and large sections of specs/single-source-context-anchor.md. Implementation queued as TASKS.md item at #priority:medium #added:2026-05-20 — multi-step (rc + gitmeta resolver simplification → init guard removal → hook cd migration → activate/deactivate deletion → docs sweep), estimated ~600-1000 LOC net deletion. Four open questions to resolve before code: CTX_DIR transition policy, deprecation shim, editor-integration grep, implementation order.

---

## [2026-05-20-214801] ctx activate is strict-CWD; drop upward walk

**Status**: Accepted

**Context**: Bug TASKS:58 — fresh git init under a workspace with its own .context/ silently bound the parent context because activate walked up past the git boundary. Previous design (specs/single-source-context-anchor.md) preserved walk-up under 'interactive discovery' on the rationale that workspace-shared .context/ next to per-project ones was a legitimate layout.

**Decision**: ctx activate is strict-CWD; drop upward walk

**Rationale**: ctx activate is a state-setting command (exports CTX_DIR); state commands follow git's read-vs-state pattern (read walks freely, state refuses to cross repo boundaries). The workspace-shared use case is preserved by user action (cd to workspace before activating), not by inferred walk. The 'also visible upward' stderr advisory was invisible to eval-bindable invocations anyway.

**Consequence**: scan() in internal/cli/activate/core/resolve/internal.go collapsed from 49 LOC walking via rc.ScanCandidates to a single os.Stat; resolve.Selected() signature went (string, []string, error) → (string, error); writeActivate.AlsoVisible and FormatAlsoVisibleAdvisory deleted; errActivate.NoCandidates renamed to NoLocalContext(cwd) and now names PWD verbatim. Spec: specs/activate-strict-cwd.md.

---

## [2026-05-21-140236] Spec steps 1+2 merged into a single commit (cwd-anchored-context)

**Status**: Accepted

**Context**: Yesterday's spec (specs/cwd-anchored-context.md) decomposed the cwd-anchored refactor into 5 sequential steps, each intended to land as a separate commit. Step 1 (resolver swap, rc.ContextDir → cwd-anchored os.Stat) cannot compile without Step 2 (init guard removal, deletion of internal/cli/initialize/core/envmatch/) because envmatch references the soon-to-be-deleted ErrDirNotDeclared sentinel.

**Decision**: Spec steps 1+2 merged into a single commit (cwd-anchored-context)

**Rationale**: Cleanest commit boundaries beat strict spec adherence when the spec's boundaries are mechanically infeasible. Steps 1 and 2 were merged into one atomic commit; remaining steps 3 (hook cd migration), 4 (activate/deactivate deletion), 5 (docs sweep) stay as discrete commits per the spec.

**Consequence**: Spec stays authoritative for what; commit-slicing diverges for practical reasons. Future cwd-anchored work follows a 4-commit (merged) decomposition, not the spec's 5. Spec text remains as-written; the divergence is documented here, not in the spec.

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

## [2026-05-21-203052] Substrate vs. artifact placement: .context/ vs. project root

**Status**: Accepted

**Context**: Question surfaced while scaffolding specs/ctx-ai-backend.md and specs/ctx-ai-extraction-and-recall.md. User observed that specs/ is the only folder (aside from GETTING_STARTED.md) ctx-managed but outside .context/, and asked whether the placement was philosophically correct. Initial 'state vs. artifact' framing was challenged with 'by that token, isn't kb a project artifact?' — exposing that the binary cut was too coarse.

**Decision**: Substrate vs. artifact placement: .context/ vs. project root

**Rationale**: Distinguish cognitive substrate (lives under .context/) from project artifact (lives at root) by the *consumption/mutation path*, not by who manages the files. Substrate is read AND written through ctx-mediated paths (ctx agent, ctx decision add, /ctx-kb-ingest, /ctx-handover, ceremonies); artifacts are read AND edited directly by humans (specs/, CLAUDE.md, GETTING_STARTED.md, docs/). Three coupling tests sharpen the line: (a) queried via ctx-mediated paths, (b) tightly coupled to ctx pipeline machinery, (c) authored under ctx skill discipline. The kb passes all three (kb closeouts fold into handovers, /ctx-kb-ingest enforces pass-mode and citations, /ctx-kb-ask is the primary read path) so it stays under .context/. Specs pass none (referenced by commits, never loaded by ctx agent, no pipeline coupling) so they live at root. Rejected alternatives: (1) move specs/ under .context/specs/ for boundary cleanliness — fails because specs are project artifacts written for humans/reviewers/community devs and hiding them under a dotfile breaks navigability; (2) move kb/ to project root because it has artifact-like properties — fails because kb machinery (closeouts, source-coverage ledger, evidence-index schema) cannot be lifted out of .context/ without splitting things that live together; (3) keep the original 'state vs. artifact' framing — too binary, kb pushback proved a third axis was needed.

**Consequence**: Codified as a CONVENTIONS.md entry under 'File Organization'. Placement test for new ctx-related files or folders: is this consumed/mutated through ctx-mediated paths (substrate, .context/) or read/edited directly by humans (artifact, root)? Visibility complaint about .context/ being a dotfile is acknowledged but acceptable — humans navigate substrate via ctx commands and generated views (ctx site kb build, ctx serve), not via file browsers. Trade-off: the rule's correctness depends on the ctx-mediated paths actually existing for substrate files; if substrate is added but no skill/command consumes it, the placement test misclassifies. See also: CONVENTIONS.md 'File Organization' section.

---

## Group: Encryption key resolution and migration

## [2026-06-02-051330] Remove the implicit project-local .ctx.key resolution tier

**Status**: Accepted

**Context**: Picking up TASKS.md P0.8.5 ("notify fails in worktrees"), we
found `crypto.ResolveKeyPath` still auto-detects a project-local
`<contextDir>/.ctx.key` (a stat-gated tier) and prefers it over the global
`~/.ctx/.ctx.key`. That file is gitignored, so it is absent in a fresh
worktree checkout: resolution silently falls back to the (different) global
key and webhook/pad decryption fails. The v0.8.0 global-encryption-key spec
already collapsed per-project keys into one global key — calling project-local
keys "a security antipattern [key next to ciphertext]" that "broke in
worktrees" — but left the implicit auto-detect tier in place. Empirically
(built binary + isolated repo + fake webhook sink): the default global key
works in worktrees; only a project-local key reproduces the failure, and the
fire path swallows it silently.

**Alternatives Considered**:
- Approach A — worktree-aware key fallback via `git rev-parse --git-common-dir`
  to resolve the main checkout's key from inside a worktree: keeps project-local
  keys working / but adds git-awareness to key resolution, contradicts the
  CWD-anchored model, larger blast radius, and props up a deprecated mechanism.
- Approach B — copy the key into the worktree at creation (ctx-worktree skill):
  no resolver change / but agent-driven and unenforceable, widens skill
  permissions, and is redundant under a global key.
- Keep the tier, only fix the silent failure: smallest change / but leaves the
  documented security antipattern and the worktree divergence in place.

**Decision**: Remove the implicit project-local `.context/.ctx.key`
auto-detection tier from `ResolveKeyPath`. Resolution becomes: (1) explicit
`.ctxrc key_path` override, (2) global `~/.ctx/.ctx.key`, (3) project-local
path only as a degenerate fallback when the home dir is unavailable. Genuine
per-project isolation stays available via the explicit `key_path` override.
Paired with surfacing the silent fire-path failure so any stranded-key decrypt
failure is visible, not silent.

**Rationale**: The project-local key is the only thing that makes a worktree
behave differently from N side-by-side terminals in the same directory;
removing it makes them indistinguishable (the desired model) and deletes a
security antipattern the project already named. It is net deletion, consistent
with the global-key and cwd-anchored simplifications. The explicit override
covers the rare real isolation need without the ciphertext-adjacent footgun.

**Consequence**: Projects on the default global key are unaffected. Projects
with a project-local `.context/.ctx.key` resolve to the global key; their
existing local-key-encrypted `.notify.enc` / pad will fail to decrypt — now
visibly (a warning on the fire path, a surfaced error on pad/test paths)
instead of silently. Documented remedy: back up the local key, then re-key to
global or set an explicit `key_path`. No auto-migration (none exists in-tree).

**Related**: Spec specs/notify-resolution-hardening.md | Supersedes the
project-local auto-detection portion of
specs/released/v0.8.0/global-encryption-key.md | Relates to
specs/cwd-anchored-context.md and [2026-03-01] Global encryption key

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
## Group: ctxctl maintainer binary and out-of-band audit channel

## [2026-05-28-201000] ctxctl PATH-installed alongside ctx for clean roots and one binary across worktrees

**Status**: Accepted

**Context**: Initial ctxctl design wired the hook to `./ctxctl` at repo root, forcing a per-worktree build, dirtying the root, and contradicting the project's PATH-only convention (`block-non-path-ctx` enforces it for ctx).

**Decision**: ctxctl PATH-installed alongside ctx for clean roots and one binary across worktrees

**Rationale**: Mirror ctx's install pattern: build to `dist/`, install to `/usr/local/bin/ctxctl`. One binary serves all worktrees and repo copies; the local hook calls `ctxctl` from PATH so no repo-root binary is needed. Defensive `/ctxctl` + `tools/ctxctl/ctxctl` gitignores stay so stray binaries can never be committed.

**Consequence**: New Makefile targets `install-ctxctl` and `reinstall-ctxctl` mirror `install`/`reinstall`. Hook in `.claude/settings.local.json`: `cd "$CLAUDE_PROJECT_DIR" && ctxctl audit-relay`. Sets the convention for future maintainer-only binaries (`tools/<name>/` separate module, `dist/` build, PATH install). `specs/ctxctl-bootstrap.md` Interface section updated to match.

---

## [2026-05-27-161302] ctxctl is a separate Go module at tools/ctxctl (own go.mod), not cmd/ctxctl in the same module

**Status**: Accepted

**Context**: Migrating the maintainer-only audit channel out of the ctx binary (specs/ctxctl-bootstrap.md). The prior decision (handover 2026-05-26) chose same-module cmd/ctxctl, on the belief that a separate go.mod could not import ctx's internal/ packages and would force relocating/duplicating ~25 files.

**Decision**: ctxctl is a separate Go module at tools/ctxctl (own go.mod), not cmd/ctxctl in the same module

**Rationale**: That blocker was empirically disproved this session: a nested module whose path is lexically under github.com/ActiveMemory/ctx CAN import the parent module's internal/ packages (verified by build test; a non-nested 'outsider' module path is rejected). Given that, a hard module boundary beats an in-module import-graph test for the asymmetric requirement that actually matters: ctx must never break because of ctxctl. ctx's go.mod will not require tools/ctxctl, so ctx literally cannot import ctxctl; the one-directional ctxctl->ctx coupling is acceptable because ctxctl is disposable maintainer tooling ('nobody whines if ctxctl breaks; everyone suffers if ctxctl leaks into ctx'). Full self-containment (duplicating the ~20 shared internal foundations: rc, desc, config, nudge, io...) was rejected as a DRY catastrophe and a worse broken window than the one being fixed.

**Consequence**: New module tools/ctxctl (module path github.com/ActiveMemory/ctx/tools/ctxctl) reuses ctx's internal/ foundations in place; audit-channel-specific logic relocates to internal/ctxctl/; ctxctl owns its relay/CLI text as plain English Go constants under tools/ctxctl (no YAML localization, no desc/i18n engine for its own output -- no French ctxctl); a repo-root go.work (committed) wires the workspace; an import-graph guard test asserts cmd/ctx never imports internal/ctxctl. Supersedes the same-module cmd/ctxctl decision. specs/ctxctl-bootstrap.md is rewritten to match.

---

## [2026-05-24-123908] ctxctl lives at cmd/ctxctl in the same Go module, not a separate go.mod

**Status**: Accepted

**Context**: Deciding where the planned ctxctl maintainer binary lives and how to house the audit channel (which should not ship in the ctx user binary). User initially proposed tools/ctx/ctxctl with its own go.mod for dependency isolation; the Phase BT saga (TASKS.md) specified cmd/ctxctl in the same module. The audit channel is already ~25 files under internal/ (internal/cli/audit, internal/config/audit, internal/err/audit, internal/write/audit, internal/cli/system/core/audit).

**Decision**: ctxctl lives at cmd/ctxctl in the same Go module, not a separate go.mod

**Rationale**: Go compiles a package into a binary only if that binary's main transitively imports it. So audit packages under internal/ imported ONLY by cmd/ctxctl/main are excluded from the ctx binary — binary-level isolation without a module split, and zero relocation of the existing internal/ audit files. A separate go.mod cannot cleanly import the parent module's internal/ (Go module + internal/ visibility friction), forcing relocation or duplication. The only real win of a separate go.mod is dependency isolation — keeping heavy build/release deps out of ctx's module graph — which the audit channel does not need (only yaml, already a ctx dep). Defer the module-split question until a future ctxctl subcommand actually pulls in heavy isolated deps.

**Consequence**: ctxctl reuses internal/ packages verbatim. An import-graph guard test must enforce that cmd/ctx never transitively imports internal/cli/audit (so the channel stays out of the shipped binary). Refined companion rule: shipped product hooks call ctx; repo-local dev hooks (ctx's own gitignored .claude/settings.local.json) may call ctxctl. If a future ctxctl subcommand needs heavy isolated deps, revisit the module split then — not now.

---

## [2026-05-24-112626] Discipline enforcement belongs on the verbatim-relay channel, run out-of-band

**Status**: Accepted

**Context**: pad-undo Phase 1 shipped a user-facing command (ctx pad undo) without matching SKILL.md/recipe updates. The agent had read CONVENTIONS.md at session start AND knew the Constitution forbids 'I can create a follow-up task', yet still labeled the docs work 'Phase 2'. The user asked: how do we prevent this for future agents, not just this session? In-band advisory prose demonstrably does not survive mid-task tunnel vision.

**Decision**: Discipline enforcement belongs on the verbatim-relay channel, run out-of-band

**Rationale**: Verbatim relay is the ONE discipline channel in this codebase that empirically survives tunnel vision: the bordered reminder boxes (ctx remind, journal/knowledge notices) get echoed by agents every turn without filtering because the relay bypasses agent judgment. So move discipline checks onto that proven channel rather than inventing a new mechanism. Run the auditor OUT OF BAND (separate Claude Code session) for two reasons: (1) fresh-context judgment — the implementer cannot grade its own homework; (2) cost — a per-commit in-band AI gate burns API tokens on every commit, whereas a manually-triggered separate session bills against the user's interactive plan and lets them choose when to spend cycles. Programmatic test gates (internal/audit, internal/compliance) stay for mechanical checks but cannot make judgment calls like 'which recipe should mention this flag'.

**Consequence**: New generic channel: out-of-band-skill writes .context/audit/<kind>.md, ctx system check-audit hook relays unread reports verbatim, ctx audit list/show/dismiss manages lifecycle. Dismissal is digest-bound so fresh findings re-surface. The channel is kind-agnostic — the hook relays any report file, so sibling skills (/ctx-spec-trailer-audit, /ctx-capture-audit) plug in with zero hook changes. Trade-off: no automated trigger in Phase 1 (no cron/post-commit) — relies on user discipline to actually run the auditor; a user who never runs it gets no nags. Naming collision with the existing internal/audit/ AST-tests package is tolerated (different layers, no compile conflict) but flagged in the spec.

---

## Group: KB editorial pipeline (Phase KB) design

## [2026-05-16-000000] Phase KB lifts the current upstream editorial-pipeline shape, superseding the 4-phase predecessor in the brief

**Status**: Accepted

**Context**: The Phase KB spec at `specs/kb-editorial-pipeline.md` was
originally lifted from the upstream editorial pipeline in May 2026, at which
point that pipeline encoded a 4-phase model (triage / extract / reconcile /
surface). The upstream design has since evolved past that shape into a pass-mode
contract (`topic-page` / `triage` / `evidence-only`) with up-front declaration,
a 4-invariant completion circuit breaker, a source-coverage state-machine
ledger, a topic-adjacency pre-flight, a cold-reader orientation rubric,
folder-shaped topics from day one, and an explicit CLI-as-scaffold-authority
rule. The comparison note at `ideas/upstream-pipeline-comparison.md` enumerated
the deltas. The fork was whether to implement the spec as written (older shape;
faster to type; weaker as a feature) or to revise the spec to absorb the
upstream design's current shape before any code is written.

**Decision**: Phase KB lifts the current upstream editorial-pipeline shape.
`specs/kb-editorial-pipeline.md` was rewritten in place on 2026-05-16 to encode
pass-mode contract, completion circuit breaker, source-coverage state-machine
ledger, topic-adjacency pre-flight, cold-reader rubric, folder-shaped topics
from day one, CLI-as-scaffold-authority, and explicit failure-analysis section.
The original 4-phase model is superseded; the brief's two organizing principles
(LLM as migration tool; KB-of-KBs is a KB) carry forward.

**Rationale**: The upstream pipeline's evolution after the brief was drafted
reflects real pain: false-finish drift, ledger-vs-reality divergence, adjacency
invisibility, mode-muddying under operator pressure. Lifting the older shape
would mean re-fighting those wounds. The user's lift-the-whole-shape posture
(feedback memory `feedback_no_defer_unfamiliar_scope`) extends here: lift the
patterns the upstream author chose, not just the structure visible at the moment
of first contact. Concretely: folder-shaped topics from day one avoid a v1.1
migration (the upstream reference's live kb has 12 sub-topic folders under
`topics/claude-code/` alone; that depth arrives fast); the pass-mode contract
makes promise=result visible per pass instead of buried in a closeout the
operator might not read; the state-machine ledger replaces the spec's flat
`source-map.md` so "what is incomplete?" has a canonical answer; the circuit
breaker turns CONSTITUTION's "Completion Over Motion" from prose into a
mechanical gate.

**Consequence**: Phase KB tasks in `.context/TASKS.md` (line 1832 onward) now
reference the revised spec; concrete additions cover the new shape (path
constants under `internal/cli/kb/core/`, new helpers for passmode /
circuitbreaker / ledger / adjacency / coldreader / lifestage, new doctor
advisories for ledger drift + pass-mode mismatch + illegal state transitions,
generalized closeout naming `<TS>-<mode>-closeout.md`). The `internal/store/`
shape from the original spec is replaced with `internal/write/` per existing ctx
convention (writers live in `internal/write/<area>/`). Folder-shaped topics from
day one means `.context/kb/topics/<slug>/index.md` is the canonical surface, not
flat `<slug>.md`; `ctx kb topic new` is the sole scaffold writer.
Failure-analysis section is now part of the spec, with three concrete loss modes
(pass-mode bypass, ledger drift, adjacency trivialization) each carrying v1
mitigations. Spec: `specs/kb-editorial-pipeline.md`. Source:
`ideas/upstream-pipeline-comparison.md`.

---

## [2026-05-10-001857] Editorial constitution at .context/ingest/KB-RULES.md, not CONSTITUTION.md

**Status**: Accepted

**Context**: `your-project` hand-rolled an editorial pipeline at the repo root with
10-CONSTITUTION.md, colliding with .context/CONSTITUTION.md. CLAUDE.md spent
paragraphs explaining the layer split (workflow infra at repo root vs ctx layer
at .context/ vs domain content at docs/). The naming collision is the core
friction.

**Decision**: Editorial constitution at .context/ingest/KB-RULES.md, not
CONSTITUTION.md

**Rationale**: Sibling project hit and named-their-way-out-of this exact
conflict (their file is 10-INGEST_RULES.md, with an explicit naming-by-rename
rule recorded in their domain-decisions.md schema header: 'KB-side filename is
domain-decisions.md to disambiguate from the root file'). Lift the rename, not
just the feature; learn from their resolved wound rather than re-fight the
conflict.

**Consequence**: Pipeline templates use KB-RULES.md throughout
(specs/kb-editorial-pipeline.md and brief reflect this); ctx CONSTITUTION.md
retains its singular meaning as the project-level invariants file; no
layer-bleed documentation needed in CLAUDE.md to cover an avoided collision;
same naming discipline carries through to domain-decisions.md (kept separate
from DECISIONS.md by the same logic).

---

## [2026-05-10-001856] Phase KB ships handover plus editorial paired, not split

**Status**: Accepted

**Context**: Trade-off considered: handover and editorial pipeline are
technically separable. Handover alone gives narrative thread between sessions.
Editorial alone piles up closeouts that 'do you remember?' reads via the
postdated-unfolded-closeout path. Either could ship without the other; question
was whether to split into two ships for smaller risk per release.

**Decision**: Phase KB ships handover plus editorial paired, not split

**Rationale**: The closeout/fold mechanism is the integration point between the
two features. Shipping paired guarantees the fold gets real-world stress on day
one rather than being added retroactively when the second feature lands.
Better-together over smaller-ship; integration coherence over delivery cadence;
the user's lift-the-whole-shape posture extends to shipping coherence.

**Consequence**: Phase KB is bigger than either feature alone; KB-2 sub-phase
covers `your-project` port as the integration regression suite; ideas/001 handover
work folds into Phase KB rather than shipping as its own phase; the polish-PR
(Phase SK) and git-mandate (Phase RG) Phase 0 prerequisites land first to keep
Phase KB clean.

---

## [2026-05-10-001856] KB ontology is pipeline-only-writer; no /ctx-kb-decide parallel skill

**Status**: Accepted

**Context**: Designing the KB editorial layer raised the question of whether KB
editorial decisions need a parallel /ctx-kb-decide skill mirroring
/ctx-decision-add. Three resolutions tested: alpha) skill surface doubles (every
capture skill gets a kb sibling); beta) capture skills become mode-aware
routers; gamma) capture skills stay single-purpose with user discipline.

**Decision**: KB ontology is pipeline-only-writer; no /ctx-kb-decide parallel
skill

**Rationale**: All three rejected after a deeper reframe surfaced by the user:
in a KB you don't decide, you increase confidence. A claim with confidence
greater than 0.9 is fact-by-contract; lower confidence needs more evidence. Even
natural-language assertions ('we are spinning off X, anchor on this') are
semantically evidence-capture events, not decision-capture events. The sibling
pipeline-only-writer model is not rigid; it is the ontologically correct surface
for evidence-tracked knowledge.

**Consequence**: KB skill surface stays small: 4 mode skills
(ingest/ask/site-review/ground) plus 1 lightweight ctx kb note for
capture-without-pipeline; existing /ctx-decision-add etc. unchanged in
authority; users who want to record a KB editorial framing instead drop a
finding into the inbox or hand-edit the markdown directly. No router question on
every capture; no parallel skill maintenance burden.

---

## [2026-05-10-001856] Mandate git as architectural precondition

**Status**: Accepted

**Context**: ctx today silently degrades without git via commit:none sentinels
in provenance flags; doctor effectively says 'git required for this to work
properly' without enforcing. Sibling project mandates git architecturally and
says so explicitly. User confirmed N approximately 0 ctx projects in practice
run without git. Editorial pipeline lift inherits the git-required assumption
(closeout sha:/branch:, evidence-index SHA-pinned in-repo citations, handover
Provenance from git HEAD).

**Decision**: Mandate git as architectural precondition

**Rationale**: Persistent-memory promise is dishonest without an undo layer: LLM
agents are not trustworthy stewards of files; git reflog is the recovery path.
Eliminates dead-code branches across every git-touching path. Trust boundary:
refuse-on-no-git rather than auto-git-init (ctx never modifies user filesystem
outside .context/). User: we should have done this on day zero.

**Consequence**: Breaking change in next minor release; specs/require-git.md
written; commit:none sentinel becomes unreachable across gitmeta and doctor
advisories; CONSTITUTION.md amendment + DECISIONS.md entry will land during
Phase RG implementation; release notes carry one-command migration ('run git
init in any pre-existing git-less ctx project before upgrading').

---

## [2026-05-10-001820] Lift sibling editorial pipeline shape into ctx as v1, paired with handover

**Status**: Accepted

**Context**: Sibling clean-room project (analyzed undercover; not named to avoid
carryover) ships a battle-tested editorial pipeline (4 modes, 9 KB artifacts,
closeout/fold mechanism, browseable site rendering). `your-project` has been
hand-rolling the same shape for weeks at workaround cost: CLAUDE.md disables
half of ctx code-dev skills, 10-CONSTITUTION.md at repo root collides with
.context/CONSTITUTION.md, hand-typed 8-item closeouts, hand-managed 20-INBOX.md.
Considered lift-intact vs hedge-and-defer.

**Decision**: Lift sibling editorial pipeline shape into ctx as v1, paired with
handover

**Rationale**: The sibling design is field-tested under production use;
`your-project` is a live validation corpus already paying the workaround tax (N=1
lived validation beats hypothetical user research). Initial defer-on-uncertainty
instinct corrected by user pushback to lift the whole shape with a non-colliding
rename (KB-RULES.md, not CONSTITUTION.md). Two organizing principles (P1: LLM is
the migration tool; P2: a KB of KBs is a KB) make lift-the-whole-shape rational
rather than reckless.

**Consequence**: specs/kb-editorial-pipeline.md written; three TASKS.md phases
added (SK polish, RG require-git, KB editorial+handover); KB has its own write
authority separate from canonical files; closeout/fold mechanism integrates
editorial work with session continuity via handover; ideas/003 brief produced as
design source.

---

## Group: Companion-tool integration: peer-MCP, no gateway

## [2026-05-23-030000] Skill body text uses capability-first language with canonical tools as examples; install-guide docs name canonical implementations; `allowed-tools` frontmatter stays MCP-specific

**Status**: Accepted

**Context**: The 2026-05-23 "MCP gateway not worth the coupling cost" decision rejected pluggable abstraction over companion tools at the code/protocol layer (no gateway, no plugin registry). But that decision left an open question: skill body text was still hard-coding specific tool names (GitNexus, Gemini Search), and so were several `docs/` pages. The hard-coding is *its own* form of vouching — just static prescription instead of dynamic dispatch. A user with Firecrawl / sourcegraph-cody / vLLM read the skill and saw instructions naming tools they don't have; the agent couldn't self-route because the skill text told it to use specific MCP server names.

Three rule choices were considered for the body-text layer:

1. Pluggable abstraction with `.ctxrc`-declared capability mapping — rejected by the prior decision (it IS the interface-contract ownership cost we ruled out).
2. Per-tool skill variants (`ctx-architecture-enrich-gitnexus`, `…-sourcegraph`, …) — explodes the skill count without removing the prescription, just sliced thinner.
3. **Capability-first body text with canonical tools as examples** — chosen.

A parallel question existed for `docs/`: an install guide LEGITIMATELY names tools (its job is "tell me what to install"). Genericizing install commands would harm newcomers. The right split: operational/descriptive docs use the same capability-first phrasing as skills; install-guide docs name canonical implementations explicitly, with a one-liner noting equivalents work.

The `allowed-tools` frontmatter is a separate concern. Genericizing to `mcp__*` would grant skills access to EVERY connected MCP — a permission expansion, not a cosmetic change. Operators with different toolchains edit `allowed-tools` in their local skill copy or fork. A separate spec can revisit if needed.

**Decision**: Three layered rules.

1. **Skill body text** uses capability-first language ("a code-intelligence MCP", "a web-search-with-citations MCP") with the canonical implementation listed as an example ("canonical: GitNexus; equivalents include sourcegraph-cody"). Operational example calls (e.g. `mcp__gitnexus__impact({…})`) stay as canonical-impl illustrations.
2. **Install-guide docs** (`docs/home/getting-started.md`, `docs/recipes/multi-tool-setup.md`) name canonical implementations directly and provide concrete setup commands. A preamble notes that equivalents work for non-canonical toolchains.
3. **`allowed-tools` frontmatter** stays MCP-specific. Skills ship with `mcp__gitnexus__*`, `mcp__gemini-search__*` in the allowlist. Operators using different MCP servers edit the allowlist in their local skill copies.

**Rationale**: Three reinforcing properties:

- **Manifesto-aligned.** ctx no longer prescribes specific tools in skill bodies. Agents self-route based on what's connected.
- **No new abstraction layer.** Pure text rewrite. Zero code change, zero interface contract, zero coupling.
- **Discoverability preserved.** Canonical tools stay first-listed in every section so newcomers immediately learn what to install if they're starting from zero.

Alternatives explicitly rejected: code-level pluggability (2026-05-23 MCP-gateway decision); per-tool skill variants (maintenance explosion without solving the smell); "remove all tool names" (loses discoverability for new users who do want a recommendation).

**Consequence**:

- Eight skill files updated (commit f554f758): ctx-refactor, ctx-explain, ctx-code-review, ctx-remember (claude + copilot-cli), ctx-architecture, ctx-architecture-enrich, ctx-architecture-failure-analysis. Prescriptive references to specific tools rewritten as capability-first with canonical examples.
- Six docs updated alongside (this commit): architecture-exploration runbook, architecture-deep-dive recipe, skills.md reference, cli/index.md schema, getting-started.md install guide, multi-tool-setup.md recipe.
- `specs/skill-audit-companion-tool-neutrality.md` documents the per-file rewrites and the install-guide-vs-operational split for future contributors.
- New skill authors follow this rule: describe the capability, name the canonical implementation as an example, leave `allowed-tools` MCP-specific.
- If a real second-viable graph-tool ecosystem emerges and operators consistently ask for pluggable `allowed-tools`, the prior MCP-gateway decision can be revisited; the present decision doesn't preclude that future evolution.

See also: `specs/skill-audit-companion-tool-neutrality.md`, `specs/ctx-remember-silent-companion-fallback.md` (the install-nag fix that preceded this audit), the 2026-05-23 "MCP gateway not worth the coupling cost" decision above.

---

## [2026-05-23-020000] MCP gateway not worth the coupling cost; companion tools stay peer-MCP and remain not-vouched-for-by-ctx

**Status**: Accepted

**Context**: Builds on the 2026-03-12 "Recommend companion RAGs as peer MCP servers not bridge through ctx" and the earlier 2026-03-06 "Peer MCP model for external tool integration" decisions. Those framed the choice as architectural (markdown-on-filesystem invariant, avoid plugin registries). The new framing, surfaced during the triage of architecture-pipeline tasks, names a stronger ownership-shaped reason: an MCP gateway through ctx would couple ctx to the lifecycle of every gatewayed tool. If ctx proxied GitNexus, users couldn't independently `pip install gitnexus` or uninstall it — ctx would become the install/uninstall surface, the upgrade path, the version-compatibility owner. That coupling is a tax we don't want to pay for a tool we don't ship.

**Decision**: MCP gateway not worth the coupling cost; companion tools stay peer-MCP and remain not-vouched-for-by-ctx.

**Rationale**: Three independent considerations converge:

1. **Composition is already MCP's job.** Agents already compose multiple MCP servers. Adding a gateway through ctx duplicates the composition layer without adding capability — the agent could just talk to GitNexus directly. The peer model preserves that property.
2. **Ownership coupling is bidirectional.** A gateway makes ctx vouch for the peer (install, uninstall, version compatibility, error surface translation). It also makes the peer's failures surface as ctx failures from the agent's perspective, blurring the diagnostic boundary. Both directions add support burden disproportionate to the value of "one extra abstraction layer".
3. **The skills already work without it.** `/ctx-architecture-enrich` and `/ctx-architecture-failure-analysis` reference GitNexus by name in their SKILL.md instructions. The agent invokes GitNexus directly via its own MCP client. No gateway involved, no abstraction needed — the skill names the tool it expects and the agent either has it configured or doesn't. Doctor-style checks (existing TASKS.md item at line 1346) handle the "is it there?" surface without proxying.

Alternatives considered and rejected: (1) Gateway through ctx — rejected for the ownership reasons above. (2) Pluggable graph-tool abstraction with multiple candidate implementations (the now-skipped TASKS.md item) — implies ctx vouches for the interface contract across implementations, same ownership trap. (3) Optional gateway as opt-in — added complexity without removing the coupling for users who opt in; cleaner to have no gateway at all.

**Consequence**: 

- **Pluggable graph tool interface task** (TASKS.md "Explore pluggable graph tool interface", `#added:2026-03-25-120000`) **skipped** as a direct consequence — pluggability without ownership is incoherent.
- **GitNexus stays named-by-convention** in skill text. SKILL.md instructions can reference `gitnexus.*` MCP tool names directly; agents either have the configuration or fail explicitly.
- **Architecture pipeline 4th step** (`ctx-architecture-next`, added today) is *itself* gateway-free: it consumes only the Markdown artifacts produced by the prior three steps, so the synthesis layer has no MCP dependency at all. That's the right shape for any future pipeline-completing skill: read what's on disk, write a new artifact.
- **Doctor / preflight checks** for companion-tool availability remain valid (TASKS.md line 1346, "Update `ctx doctor` to check for graph tool availability"). Checking that a peer exists is not the same as proxying through it.
- **The earlier 2026-03-12 peer-MCP decision is not superseded** — it's reinforced. This entry adds the ownership lens; the architectural reasoning from that entry still applies.

See also: `ideas/spec-companion-intelligence.md` (the original peer-MCP design), `ideas/gitnexus-contextmode-analysis.md`, the now-skipped pluggable-interface task in TASKS.md.

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

## Group: Localizable vocabulary and i18n primitives

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
## [2026-05-23-001500] Keep `i18n.Fold` strict; add `i18n.MatchKey` as the separate diacritic-insensitive primitive

**Status**: Accepted

**Context**: The placeholder localization task (line 287, specs/placeholder-i18n.md) introduced `internal/i18n.Fold` (commit 435d6670) as the project-mandated case-fold primitive. Field testing in the validator integration test surfaced an ergonomic problem: `Fold` preserves Unicode-defined linguistic distinctions (`İ` ≠ `i`, `ü` ≠ `u`), so a Turkish user with a Turkish keyboard typing `İPTAL` would not reject against an `iptal` entry in `.ctxrc` — they'd need to enumerate every diacritic variant of their vocabulary. Same problem for German `Straße`/`strasse`, French `café`/`cafe`, etc. The bilingual case (English keyboard plus Turkish prose) made the friction unavoidable for non-English users.

**Decision**: Keep `i18n.Fold` strict; add `i18n.MatchKey` as the separate diacritic-insensitive primitive.

**Rationale**: Two distinct primitives with explicit contracts beats one primitive that conflates them. `Fold` stays a strict Unicode case-fold (`cases.Fold` semantics, `İ` ≠ `i`) — required for callers that need linguistic-precision: identifier deduplication, parsing, security-relevant comparison. `MatchKey` is `Fold + NFKD + strip(U+0300..U+036F)` — collapses Latin/general diacritics (Turkish dotted-I, German umlaut, French accents, Vietnamese horn) so casual keyboard variation matches transparently. Alternatives considered: (1) tighten `Fold` itself to include the strip step — rejected as conflating two contracts; any future caller that wants Unicode-precise comparison would silently get the looser semantics, with no compile-time signal. (2) Provide one primitive with an options/flags arg — rejected as bloated API for two distinct use cases. (3) Document the friction and let users enumerate variants — rejected as user-hostile for non-English projects, which is exactly the population the localization spec was meant to serve. (4) Two primitives, picked at call site — CHOSEN. The `Picking the right primitive` section in `internal/i18n/doc.go` gives the rule: "if your matcher compares user input against a vocabulary list and the user might type with or without diacritics, use MatchKey; otherwise Fold."

**Consequence**: Two primitives to maintain (small — both are ~10 LoC over the upstream `cases` package). Call sites pick the right one explicitly. The placeholder validator uses MatchKey at all three sites (loader, .ctxrc merge, input lookup). Tests guard both halves: MatchKey collapses Turkish/German/French/Spanish/Catalan/Czech/Vietnamese as expected; preserves script-essential marks for Arabic/Indic/Hebrew/CJK; Fold stays strict. The compliance AST ban applies to both — no new direct `strings.ToLower` callers can enter the codebase without using one of these. See also: specs/i18n-fold-helper-and-ban.md, LEARNINGS.md `Unicode block separation makes diacritic-stripping surgical`.

---

## [2026-05-10-181404] Placeholder overrides use EXTEND not REPLACE semantics

**Status**: Accepted

**Context**: When localizing the placeholder set used by
validate.RejectPlaceholder, .ctxrc gains a placeholders: list. The existing
precedent (rc.SessionPrefixes) uses REPLACE semantics: any non-empty user list
completely replaces the shipped defaults. Placeholders need a different rule.

**Decision**: Placeholder overrides use EXTEND not REPLACE semantics

**Rationale**: The dominant case in this codebase is Tarzan Turkish —
bilingual EN+TR projects where users need both English (TBD, n/a, see chat) and
Turkish (iptal, yapılacak, görüşülecek) placeholders rejected
simultaneously. REPLACE would force users to re-list every English default just
to add one Turkish term, which they would skip and silently lose half the
validator's coverage. EXTEND appends user list onto the shipped defaults so
partial overrides do not regress baseline protection.

**Consequence**: rc.Placeholders() must combine defaults + user list with
case-folded de-duplication, diverging from the SessionPrefixes pattern. A future
maintainer reading both accessors side-by-side will notice the inconsistency;
the divergence is intentional and Spec: specs/placeholder-i18n.md captures why.
If REPLACE is later wanted, add an opt-in placeholders_replace: true toggle
rather than flipping the default.

---

## Group: Embedded assets and editor-integration harnesses

## [2026-05-11-211246] Embedded and separately-published harnesses use distinct CI and release pipelines

**Status**: Accepted

**Context**: ctx ships two kinds of artifact. Embedded harnesses (OpenCode
plugin, Copilot CLI scripts, Claude/OpenCode/Copilot CLI skills, git trace
hooks, etc.) live under internal/assets/, are //go:embed'd into the ctx Go
binary, and reach users via 'ctx setup' writing their bytes to disk.
Separately-published harnesses (currently just the VS Code extension under
editors/vscode/) build to their own artifact (.vsix), publish to a third-party
channel (VS Code Marketplace under publisher 'activememory'), version
independently, and reach users via that channel's update mechanism. Until this
session, the boundary was implicit: doc.go and embed_test.go talked only about
the embedded tree; release.yml only built the Go binary; nothing in CI exercised
the vscode extension at all. A reviewer's first read of
internal/assets/integrations/ was 'this is a dumping ground' precisely because
the contract was not documented.

**Decision**: Embedded and separately-published harnesses use distinct CI and
release pipelines

**Rationale**: Conflating the two would have one of two consequences: (a)
shoehorning vscode into //go:embed, which means baking a .vsix or its sources
into the Go binary and writing them out at setup time -- bloating the binary
with bytes most users never use, and forcing the Go release cadence onto
something with its own marketplace cadence; or (b) leaving the vscode harness
ungated 'because it's different' -- which is what we had, and which is how typos
ship. The right move is to acknowledge the two patterns are first-class peers,
give each a documented home (internal/assets/ vs. editors/<editor>/), and gate
each in CI with the toolchain appropriate to its release pipeline (Go
test/build/vet for embedded; npm ci + esbuild + tsc for vscode). Future
harnesses pick a pattern explicitly at placement time rather than drifting.

**Consequence**: internal/assets/README.md now carries the 'Embedded vs.
Separately-Published: At a Glance' table as the canonical reference.
.github/workflows/ci.yml gained a vscode-extension job that gates the
marketplace publish path. editors/vscode/README.md gained a 'Release' section
with checklist and explicit notes on which CI gates protect the manual vsce
publish. The two patterns are now first-class: a new harness must declare which
it follows before placing files. Open implications: (1) anyone proposing to lift
integrations/ out of internal/assets/ should re-read this decision -- the no-../
//go:embed constraint plus the pattern-asymmetry are the load-bearing reasons
against; (2) the embedded-only quality gaps tracked in TASKS.md (shellcheck,
PSScriptAnalyzer, skill frontmatter validity) and the separately-published
quality gaps (vscode test rot, lint, vsce package dry-run) live in distinct
gap-task clusters and should not be merged. Spec:
specs/internal-assets-readme.md.

---

## [2026-05-11-000000] Embedded foreign-language assets under internal/assets/ are intentional, not a smell

**Status**: Accepted

**Context**: A diagnostic conversation surfaced that
`internal/assets/integrations/` contains TypeScript
(`opencode/plugin/index.ts`), Bash and PowerShell scripts
(`copilot-cli/scripts/`), JSON, YAML, and Markdown — none of it Go source. The
first-glance read was "internal/ has become a dumping ground for non-Go tooling;
lift integrations/ out." Audit of `embed.go` proved otherwise: every file under
`integrations/` is captured by an explicit `//go:embed` directive and shipped
inside the ctx binary as raw bytes, then written to the user's filesystem at
`ctx setup` time. The smell was real (no contract document existed to explain
this) but the architectural diagnosis was wrong.

**Decision**: Embedded foreign-language assets stay under `internal/assets/`.
The `internal/` directory is honoring Go's import-privacy convention; the
contract is "everything in this tree is `//go:embed`'d into the binary as
bytes." A `README.md` at `internal/assets/README.md` documents the contract;
`internal/assets/doc.go` continues to serve the Go-doc audience.

**Rationale**: Three reasons against lifting:

1. **Hard Go constraint**: `//go:embed` directives cannot reference parents (no
`../`). Moving assets out of the embed.go directory tree forces moving (or
duplicating) the embed package itself, with import-path blast radius across
every consumer. The relocation cost is disproportionate to the readability win.
2. **Idiomatic Go**: `internal/` is about import privacy, not source language.
Projects like Kubernetes and Cobra ship embedded foreign-language payloads from
`internal/` without considering it a smell.
3. **The actual fix is cheaper**: the smell was a missing contract document, not
a misplaced directory. A README that names the rule ("everything here is
`//go:embed`'d; foreign-language files are intentional payload") resolves the
legibility problem at zero structural cost. Dev tooling *about* the embedded
payload (e.g. `tsconfig.json` for the TS plugin) is what does not belong inside
the embed tree — that goes in a sibling tooling directory.

**Consequence**: Future contributors who feel the same "internal/ is a dumping
ground" instinct will find a README documenting why the layout is correct. The
README also enumerates current quality gates (presence, format parse, schema
integrity) and the known gaps (TypeScript type-check, shellcheck,
PSScriptAnalyzer, skill frontmatter validation) — gaps now spawned as discrete
Phase 0 tasks. The line-30 `tsc --noEmit` task is redirected: its tooling files
must live in a sibling directory outside `internal/assets/` to honor the embed
contract.

**Related**: Spec: specs/internal-assets-readme.md

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

## [2026-05-22-161800] OpenCode plugin: agent shell tool not anchored to project root under cwd-anchored

**Status**: Accepted

**Context**: specs/cwd-anchored-context.md changed ctx's resolver from CTX_DIR env-var to $PWD/.context/. The opencode plugin (internal/assets/integrations/opencode/plugin/index.ts) previously injected CTX_DIR into the agent's shell tool via the shell.env hook so agent-issued 'ctx' commands resolved to the right project. Under cwd-anchored, ctx no longer reads CTX_DIR; the only way to make ctx resolve correctly is to ensure the shell tool's cwd is the project root. But @opencode-ai/plugin v1.4.x exposes only 'env' on the shell.env hook output type ({ env: Record<string, string>; }) — no 'cwd' field. The plugin cannot force the agent shell into the project root from inside the SDK contract.

**Decision**: OpenCode plugin: agent shell tool not anchored to project root under cwd-anchored

**Rationale**: Decision: drop the shell.env handler entirely and document that users must launch OpenCode from the project root. Plugin-internal subprocess calls (ctx.$.cwd(ctx.directory)) remain anchored, so the ceremony invocations (session.created, session.idle, tool.execute.after, experimental.session.compacting) still work. Only the agent-issued shell commands lack an anchoring channel. Alternatives considered: (1) keep the handler with a dummy env injection 'in case the SDK adds cwd' — rejected as dead code with no semantic load; (2) inject PWD/OLDPWD to influence the shell's cwd — rejected as brittle and outside the SDK type contract; (3) patch @opencode-ai/plugin upstream to expose cwd on shell.env — deferred (real upstream work, coordination required, degrades gracefully without it); (4) document the launch-from-root requirement and remove the handler — CHOSEN. The cwd-anchored error message ('ctx: no .context/ at <pwd>. Run `ctx init` here, or cd to a project that has one.') is itself clear and self-fixing, so the friction is bounded.

**Consequence**: Agent-issued 'ctx' commands fail with the clear cwd-anchored error when OpenCode is launched from outside the project root. User re-launches from the right directory. Plugin's own ceremony calls continue to work. Trade-off: minor user-facing friction in exchange for not building unsupported SDK behaviour into the plugin. Escalation path if this becomes recurring: alternative 3 (upstream SDK PR adding cwd to shell.env output type). See also: specs/cwd-anchored-context.md, LEARNINGS.md 'Cross-language coverage gap'.

---

## [2026-04-26-231517] OpenCode tool.execute.before omission is permanent; block-dangerous-commands will not become a ctx Go subcommand

**Status**: Accepted

**Context**: The 2026-04-26-152858 decision shipped the OpenCode plugin without
a tool.execute.before hook and noted "Re-add when block-dangerous-commands is
promoted to the ctx Go binary." Revisited: that promotion is no longer planned.
Keeping the open task on the books makes future sessions believe a re-add is
pending.

**Decision**: We will not promote block-dangerous-commands to a ctx system Go
subcommand. The OpenCode plugin's missing tool.execute.before hook is permanent,
not deferred.

**Rationale**: The Cobra exit-1 / `{ blocked: true }` interaction makes any shim
hostile to users without the Claude wrapper, and the safety-hook gap is
acceptable given OpenCode's positioning. Recording this avoids the tax of a
perpetually-pending follow-up that no one intends to land.

**Consequences**: TASKS.md item "Promote 'block-dangerous-commands' to a real
ctx system Go subcommand…" marked `[-]` skipped. The 2026-04-26-152858
rationale's "Re-add when…" clause is void; the underlying
ship-without-the-hook decision remains in force. Other (non-OpenCode) editor
integrations that want a dangerous-command safety net will need a different
mechanism.

**Related**: Amends [2026-04-26-152858] OpenCode plugin ships without
tool.execute.before hook (rationale's deferred re-add is now closed).

---
## [2026-04-26-152905] Editor-integration plugins must filter post-commit to actual git commit invocations

**Status**: Accepted

**Context**: Original PR #72 OpenCode plugin ran 'ctx system post-commit' after
every shell tool call, not only after real commits

**Decision**: Editor-integration plugins must filter post-commit to actual git
commit invocations

**Rationale**: post-commit is meaningful only after a real commit lands; firing
on every shell call is noise that trains users to ignore the resulting nudges

**Consequences**: Editor plugins always sniff the actual command string (regex
on the extracted command) before triggering capture nudges that target specific
commands. Same pattern applies to any future hook that targets a specific
porcelain command.

## [2026-04-26-152858] OpenCode plugin ships without tool.execute.before hook

**Status**: Accepted

**Context**: The natural fit (block-dangerous-commands) doesn't exist as a ctx
system Go subcommand; shimming to it would block every shell call on installs
without the Claude wrapper because Cobra's unknown-command exit 1 is read as {
blocked: true } by OpenCode

**Decision**: OpenCode plugin ships without tool.execute.before hook

**Rationale**: Better to ship a feature-narrower plugin than one that bricks the
editor for users without the wrapper. Re-add when block-dangerous-commands is
promoted to the ctx Go binary.

**Consequences**: OpenCode users get bootstrap, persistence, post-commit, and
task-completion nudges but no dangerous-command safety net.
specs/opencode-integration.md records the deliberate omission.

## Group: Context injection, hooks, and session-state architecture

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
## [2026-05-08-195040] Gate mkdir inside state.Dir() rather than per-caller

**Status**: Accepted

**Context**: Closing the cross-IDE Cursor leak required preventing state.Dir()
from materializing .context/state/ in uninitialized projects. Two viable
options: (A) gate inside state.Dir itself; (B) require every caller to check
Initialized() first.

**Decision**: Gate mkdir inside state.Dir() rather than per-caller

**Rationale**: Option (A) makes the invariant ('no .context/state/ in
uninitialized projects') structurally enforced. The leak's root cause was
exactly the (B)-style assumption — checkreminder.Run deliberately skipped the
gate to print provenance unconditionally, and that path silently produced the
leak via Preamble -> nudge.Paused -> PauseMarkerPath -> state.Dir. As long as
Dir() mkdirs unconditionally, every future caller is one missed gate away from
re-introducing the bug.

**Consequence**: state.Dir() now returns errCtx.ErrNotInitialized for uninit
projects. Hook callers' existing 'if dirErr != nil { return nil }' branches
absorb it silently; interactive callers (ctx add, task complete, prune) surface
a path-bearing message via cobra. cooldown.TombstonePath was refactored to
delegate to state.Dir so the gate also covers the PreToolUse 'ctx agent' path.
memory.SaveState/LoadState were left alone because they use 0755 (different leak
class) and are user-initiated, not auto-triggered.

---

## [2026-04-25-014704] Tighten state.Dir / rc.ContextDir to (string, error) with sentinel errors

**Status**: Accepted

**Context**: Old single-return form returned ('', nil) when CTX_DIR was
undeclared. Callers that filtered only on err != nil joined empty stateDir with
relative names and wrote state files into CWD instead of .context/state/.

**Decision**: Tighten state.Dir / rc.ContextDir to (string, error) with sentinel
errors

**Rationale**: Returning a sentinel ErrDirNotDeclared makes the empty-path case
unrepresentable in a 'looks fine' branch. Forces every caller through the same
explicit gate.

**Consequence**: All callers needed migration; tests had to declare CTX_DIR
explicitly. In return, the filepath.Join('', rel) trap is closed by
construction.
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

## Archived: stale / superseded

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
