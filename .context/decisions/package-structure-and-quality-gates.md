# package-structure-and-quality-gates

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

