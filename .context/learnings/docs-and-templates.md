# docs-and-templates

## [2026-07-16-120002] site/ is tracked build output — rebuild and commit it with docs changes (consolidated)

**Consolidated from**: 2 entries (2026-02-27, 2026-06-07)

- `site/` is intentionally git-tracked (zensical build via `make site`); there is NO GitHub Pages workflow or CI build step — nothing rebuilds it downstream.
- Any change under `docs/` requires regenerating `site/` (`make site`) and committing the rebuilt output in the SAME commit; otherwise it silently drifts and surfaces later as a large modified set (observed: a surprise 189-file `site/` drift mid-session).
- Always `git add site/` with doc commits; never gitignore `site/` or treat it as ephemeral noise.

---

## [2026-07-16-120003] Contributor PRs: check for removal-reintroduction and schedule convention follow-up (consolidated)

**Consolidated from**: 2 entries (2026-03-15, 2026-04-01)

- PRs based on older code reintroduce intentionally-removed features (PR #45 brought back prompt templates, PROMPT.md, IMPLEMENTATION_PLAN.md removed in March). Cross-reference DECISIONS.md before accepting additive PR content — do not assume it is purely additive.
- Merging with known MECHANICAL gaps is fine if the follow-up is immediate (PR #42 left ~12 inline strings, no embed_test coverage, substring matching in containsOverlap). Track in `ideas/pr{N}-review-status.md` during review, merge when the architecture is sound, and fix the convention gaps in a same-day follow-up commit.

---

## [2026-07-16-120004] Drift/detection scripts need existence checks and exclusion lists for intentional & illustrative uses (consolidated)

**Consolidated from**: 3 entries (2026-02-27, 2026-03-23, 2026-03-24)

- Path-reference drift detection must `os.Stat` the top-level directory of a path before flagging it — bare filenames and paths under non-existent dirs are almost always documentation examples (fixed 23 false positives from prose in CONVENTIONS.md/ARCHITECTURE.md like loader.go, session/run.go, sync.Once).
- Shell grep on constant VALUES cannot distinguish constant TYPES (cobra `Use*` syntax strings vs `DescKey*` YAML lookup keys); type-aware checks need AST analysis, not grep (captured in specs/ast-audit-tests.md; the lint-drift fix shipped in v0.8.0).
- Convention-enforcement detection scripts (e.g. typography) need proactive exclusion patterns for files where the flagged pattern is intentional data, not prose: `*_test.go` and constant-definition files (e.g. config/token/delim.go) are the common false-positive sources.

---

## [2026-06-07-170005] TypeScript/integration test surfaces & exclusion rot (consolidated)

**Consolidated from**: 4 entries (2026-05-11 to 2026-05-22)

- Removing/renaming any cross-language contract (env channel, feature flag) is a
  FOUR-surface cleanup, not three: (1) Go build+lint+test, (2) audit/compliance
  tests, (3) asset templates (CLAUDE.md, AGENT_PLAYBOOK, hooks.json), (4)
  TypeScript-typed integrations (opencode plugin, vscode extension). The TS
  surface is invisible to `go test ./...` by design; tsc --noEmit only runs in
  CI unless invoked from tools/typecheck/opencode/ or editors/vscode/. Want: a
  `make typecheck` target wrapping both, in pre-commit + release checklist.
- tsc resolves node_modules by walking up from each SOURCE file's location, not
  the tsconfig's location. For a cross-tree setup (tsconfig in dir A, include
  points at dir B), add explicit baseUrl + paths (+ typeRoots) to the tsconfig
  so node_modules can live with the tooling.
- vitest's vi.mock() does NOT preserve Node's async-deferral guarantees: a
  mocked execFile (or fs.readFile, dns.lookup, http.request) can fire its
  callback synchronously, TDZ-trapping a closure that's provably safe by Node's
  contract. When a linter suggests tightening let→const on a var captured
  through an async callback, verify under the test runner; the safe form is
  `let` + an eslint-disable naming the mock constraint.
- A test suite excluded from BOTH typecheck and execution rots compounding:
  re-enable cost = sum of ALL drift since last green (named 2 breakages, found
  18 more on first run), not just the named bug. expect.anything()/expect.any()
  pass typecheck so only execution catches the drift. When adding any tooling
  exclude (tsconfig glob, vitest ignore, pytest --ignore), file an immediate
  follow-up whose acceptance criterion is removal; budget 5–20× the named
  scope on re-enable.

---

## [2026-06-07-170007] Documentation, template & asset drift (consolidated)

**Consolidated from**: 6 entries (2026-02-24 to 2026-04-01)

- Exhaustive lists/counts in architecture docs (package lists, command tables,
  skill counts) drift silently because nobody re-counts (23 listed vs 31
  actual). Add `<!-- drift-check: <shell command> -->` markers; run
  /ctx-architecture after adding packages/commands (/ctx-drift catches stale
  paths but not missing entries).
- Template changes are invisible to existing projects until `ctx init --force`;
  non-destructive init never re-syncs. checkTemplateHeaders was added to `ctx
  drift`.
- Any content duplicated in two locations without a sync mechanism drifts
  silently (Copilot CLI skills as condensed ctx skills; assets/why/ vs docs/).
  Wire freshness checks as build PREREQUISITES, not optional audit steps (make
  sync-copilot-skills, make sync-why must be build deps).
- Machine-generated CLAUDE.md content (GitNexus injected 121 lines / 61%)
  consumes per-turn budget without proportional value. Auto-generated content
  belongs in on-demand skills; prefer a one-line pointer over inline content.
  Audit CLAUDE.md periodically.
- CLI reference docs outpace implementation (ctx remind had no CLI, recall sync
  no Cobra wiring) — verify with `ctx <cmd> --help` before releasing docs.
  Agent style-violation sweeps are unreliable (8 found vs 48+ actual); follow
  with targeted grep + manual classification. Documentation audits must compare
  against known-good examples for the COMPLETE standard, not mere presence. New
  audit concerns (e.g. dead links) belong in an existing audit skill's checklist
  before becoming standalone.

---

## [2026-06-07-170008] User-facing text & magic-string discipline (consolidated)

**Consolidated from**: 4 entries (2026-03-14 to 2026-04-04)

- Any string containing English words alongside format directives ("%d entries
  checked") is user-facing text belonging in YAML assets — the format-verb
  (and URL-scheme, HTML-entity, err/) exemptions were removed from
  TestNoMagicStrings.
- Any string reaching the user, including stderr warnings, routes through
  assets.TextDesc() for i18n readiness; create text.yaml entries and asset keys
  first.
- Magic-string cleanup is fractal: each fix puts adjacent code under scrutiny (4
  Fprintf calls → over-tokenized formats, magic hex perms, TOML tokens,
  missing docstrings). Budget 2–3× the initial estimate; commit per layer.
- Naming a constant _alt and hardcoding one non-English language as a built-in
  default is implicit language favoritism that doesn't scale (alt_2? alt_3?).
  Use configurable lists from the start; default to a single canonical value,
  all extensions user-configured equally.

---

## [2026-04-14-010134] Constitution forbids context window as a deferral excuse

**Context**: Mid-session, agent proposed pacing through doc.go rewrites with the
reasoning that context budget was tight.

**Lesson**: The CONSTITUTION explicitly lists 'We are running out of context
window' as a forbidden deferral phrase under No Excuse Generation. The rule is
real and applies to agent self-pacing, not just user-facing answers.

**Application**: When tempted to scope down because context is tight, re-read
the constitution. The right move is to do the work end-to-end, not to ask the
user which slice to skip.

---

## [2026-04-14-010134] docs/cli/system.md and embed/cmd/system.go diverged on bootstrap promotion intent

**Context**: Header comment in internal/config/embed/cmd/system.go claimed
bootstrap was promoted to top-level; the bootstrap.go registration never
actually promoted it. Two contradictory sources of truth coexisted silently.

**Lesson**: Header-comment claims about command-tree structure are unaudited;
they can drift from registrations without any test failing. Trust the code, not
the comment.

**Application**: When evaluating any package_name namespace cleanup type claim
about command structure, verify against the actual cobra registration in
internal/bootstrap/group.go before acting.

---

## [2026-02-28-184758] ctx pad import, ctx pad export, and ctx system resources make three hack scripts redundant

**Context**: Audited hack/ scripts against ctx CLI surface

**Lesson**: As ctx CLI grew, several hack scripts became wrappers around
built-in commands (pad-import.sh -> ctx pad import, pad-export-blobs.sh -> ctx
pad export, resource-watch.sh -> watch -n5 ctx system resources)

**Application**: Periodically audit hack/ for scripts that ctx has absorbed

---

## [2026-02-28-184647] Getting-started docs assumed Claude Code as the only agent

**Context**: The installation section opened with 'A full ctx installation has
two parts' — binary + Claude Code plugin — leaving non-Claude-Code users
without a clear path

**Lesson**: Installation docs should lead with the universal requirement (the
binary) and present agent-specific integration as conditional

**Application**: When writing docs for multi-tool projects, frame the common
denominator first, then branch by tool

---

