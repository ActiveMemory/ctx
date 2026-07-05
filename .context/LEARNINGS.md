# Learnings

<!--
UPDATE WHEN:
- Discover a gotcha, bug, or unexpected behavior
- Debugging reveals non-obvious root cause
- External dependency has quirks worth documenting
- "I wish I knew this earlier" moments
- Production incidents reveal gaps

DO NOT UPDATE FOR:
- Well-documented behavior (link to docs instead)
- Temporary workarounds (use TASKS.md for follow-up)
- Opinions without evidence
-->

<!-- INDEX:START -->
| Date | Learning |
|----|--------|
| 2026-06-10 | Stock macOS bash 3.2 treats empty-array expansion as unbound under set -u |
| 2026-07-04 | Typed JSON round-trips silently drop user-owned keys |
| 2026-07-03 | New Claude model families silently fall to the 200k default in ModelContextWindow |
| 2026-06-07 | Pin an on-disk contract before splitting work across parallel agents |
| 2026-06-07 | site/ is tracked build output — rebuild and bundle it with doc commits |
| 2026-06-07 | ctx-dream is headless-first; invoking /ctx-dream interactively is debugging, not the UX |
| 2026-06-07 | ctx-dream design principles (consolidated) |
| 2026-06-07 | internal/audit & compliance gates for new code (consolidated) |
| 2026-06-07 | Error handling: sentinels, unwrapping, and silent discards (consolidated) |
| 2026-06-07 | git CLI wrapping quirks (consolidated) |
| 2026-06-07 | TypeScript/integration test surfaces & exclusion rot (consolidated) |
| 2026-06-07 | Editorial KB pipeline: design epistemology (consolidated) |
| 2026-06-07 | Documentation, template & asset drift (consolidated) |
| 2026-06-07 | User-facing text & magic-string discipline (consolidated) |
| 2026-06-07 | Constant placement & helper smells (consolidated) |
| 2026-06-07 | Convention enforcement: mechanical gates over prose (consolidated) |
| 2026-06-07 | Go toolchain, gofmt & build-tag pitfalls (consolidated) |
| 2026-06-07 | Stale-task triage & verify-before-acting (consolidated) |
| 2026-06-07 | Refactor mechanics: subagents, cascades & golden fixtures (consolidated) |
| 2026-06-07 | Linting, gosec & I/O chokepoints (consolidated) |
| 2026-06-07 | Hook mechanics, output channels & compliance (consolidated) |
| 2026-06-07 | State, tombstones, logs & filesystem hygiene (consolidated) |
| 2026-06-07 | Host-pressure alerting: use derivatives, not levels (consolidated) |
| 2026-06-07 | Go test isolation & patterns (consolidated) |
| 2026-06-01 | Guard managed blocks before regenerating; don't trust the span to be machine-owned |
| 2026-05-31 | macOS GUI apps inherit a minimal PATH; augment it to find a user-installed CLI |
| 2026-05-31 | Tauri 2 requires rustc >= 1.88; bump the toolchain before cargo check |
| 2026-05-28 | ctx kb: single topic-enumeration site; life-stage count is consumer-side |
| 2026-05-28 | A non-root Go module nested under the main module's path CAN import its internal/ packages |
| 2026-05-28 | cobra's legacyArgs lets unknown subcommands silently succeed on non-root groups |
| 2026-05-25 | Skill shipping location: _ctx- prefix is repo-internal, internal/assets/claude/skills/ctx-* is bundled and shipped |
| 2026-05-23 | Unicode block separation makes diacritic-stripping surgical — no per-script handling needed for Arabic/Indic/Hebrew/CJK |
| 2026-05-20 | macOS /var symlink trips path-equality; use EvalSymlinks with parent-resolution fallback |
| 2026-05-20 | Handover filenames are archaeology; parse by generated-at, not filename |
| 2026-05-20 | /ctx-plan is named after its input, not its output |
| 2026-05-17 | Creator confusion is the strongest doc-quality signal — louder than any user signal |
| 2026-05-17 | `_helpers.go` / `_utils.go` filenames are project anti-pattern; use domain nouns |
| 2026-05-11 | Naive Markdown line-sweep corrupts multi-line code spans and YAML lists |
| 2026-05-08 | Cursor imports Claude Code hooks and sets CLAUDE_PROJECT_DIR per workspace |
| 2026-04-14 | Constitution forbids context window as a deferral excuse |
| 2026-04-14 | docs/cli/system.md and embed/cmd/system.go diverged on bootstrap promotion intent |
| 2026-04-14 | Raft-lite trade-off is the load-bearing choice in internal/hub |
| 2026-04-14 | Brand-name handling in title-case engines must cover possessives |
| 2026-04-13 | GPG signing from non-TTY contexts requires pinentry-mac (or equivalent) |
| 2026-04-13 | rc.ContextDir() is the single source of truth — fix the resolver, not callers |
| 2026-04-09 | Pad index shifting is a real UX bug in batch operations |
| 2026-04-03 | Bulk rename and replace_all hazards (consolidated) |
| 2026-04-03 | Import cycles and package splits (consolidated) |
| 2026-04-03 | Skill lifecycle and promotion (consolidated) |
| 2026-04-03 | desc.Text() is the single highest-connectivity symbol in the codebase |
| 2026-04-01 | Contributor PRs based on older code reintroduce removed features |
| 2026-03-31 | Convention audits must check cmd/ purity, not just types and docstrings |
| 2026-03-31 | JSON Schema default fields cause linter errors with some validators |
| 2026-03-30 | lint-docstrings.sh greedy sed hid all return-type violations |
| 2026-03-24 | lint-drift false positives from conflating constant namespaces |
| 2026-03-23 | Typography detection script needs exclusion lists for intentional uses |
| 2026-03-20 | Commit messages containing script paths trigger PreToolUse hooks |
| 2026-03-18 | Lazy sync.Once per-accessor is a code smell for static embedded data |
| 2026-03-17 | Write package output census: 69 trivial/simple, 38 consolidation candidates, 18 complex |
| 2026-03-15 | Contributor PRs need post-merge follow-up commits for convention alignment |
| 2026-03-06 | Stats sort uses string comparison on RFC3339 timestamps with mixed timezones |
| 2026-03-05 | Blog post editorial feedback is higher-leverage than drafting |
| 2026-03-04 | CONSTITUTION hook compliance is non-negotiable — don't work around it |
| 2026-03-02 | Hook message registry test enforces exhaustive coverage of embedded templates |
| 2026-03-01 | Model-to-window mapping requires ordered prefix matching |
| 2026-03-01 | TASKS.md template checkbox syntax inside HTML comments is parsed by RegExTaskMultiline |
| 2026-02-28 | ctx pad import, ctx pad export, and ctx system resources make three hack scripts redundant |
| 2026-02-28 | Getting-started docs assumed Claude Code as the only agent |
| 2026-02-28 | Plugin reload script must rebuild cache, not just delete it |
| 2026-02-27 | site/ directory must be committed with docs changes |
| 2026-02-27 | Doctor token_budget vs context_window confusion |
| 2026-02-27 | Drift detector false positives on illustrative code examples |
| 2026-02-26 | Webhook silence after ctxrc profile swap is the most common notify debugging red herring |
| 2026-02-26 | Agent context loading and task routing (consolidated) |
| 2026-02-26 | PATH and binary handling (consolidated) |
| 2026-02-26 | Task management and exit criteria (consolidated) |
| 2026-02-26 | Agent behavioral patterns (consolidated) |
| 2026-02-26 | ctx add and decision recording (consolidated) |
| 2026-02-24 | CLI tools don't benefit from in-memory caching of context files |
| 2026-01-28 | IDE is already the UI |
| 2026-04-29 | BunShell ctx.$ calls echo stdout to OpenCode's process unless .quiet() is set — leaks visible noise |
| 2026-04-29 | OpenCode plugin compaction interop is breadcrumb-mediated: own your context preservation explicitly |
| 2026-04-29 | @opencode-ai/plugin event hook is a single dispatcher, not an object of named handlers |
| 2026-04-29 | OpenCode plugin hooks like shell.env take (input, output) and mutate; returned objects are ignored |
| 2026-04-29 | OpenCode shell.env injects env only into agent's shell tool, not into plugin's own ctx.$ calls |
| 2026-04-26 | OpenCode auto-loads only flat .ts files under .opencode/plugins/; subdirectories are ignored |
| 2026-04-26 | OpenCode opencode.json MCP shape: command is Array<string>, no separate args field |
| 2026-04-26 | ctx system help can list project-local hooks not in the Go binary |
| 2026-04-25 | Confident code comments can pull an LLM away from first-principles knowledge |
<!-- INDEX:END -->

---

## [2026-06-10-223128] Stock macOS bash 3.2 treats empty-array expansion as unbound under set -u

**Context**: make audit aborted at hack/lint-drift.sh line 39 on a stock Mac (bash 3.2.57) with 'exclude_args[@]: unbound variable' while gating the #93 TOCTOU fix; the script works fine on Linux bash 4+

**Lesson**: bash 3.2 (what every stock macOS ships, GPLv2 freeze) treats "${arr[@]}" on an empty array as an unbound-variable error under set -u; bash 4.4+ does not. Any 'set -u' script that expands a possibly-empty array breaks for every Mac contributor

**Application**: Guard with ${arr[@]+"${arr[@]}"} (parameter-expansion alternate form) wherever a possibly-empty array is expanded in hack/ scripts; test hack/ scripts with /bin/bash, not just homebrew bash

## [2026-07-04-152957] Typed JSON round-trips silently drop user-owned keys

**Context**: The init permissions merge read settings.local.json into a typed struct and re-marshaled it, so every key ctx does not model (env, a user's statusLine) vanished on re-init

**Lesson**: A typed read-modify-write of a file users also edit by hand is silent data loss; unknown fields do not survive encoding/json round-trips

**Application**: Mutate shared JSON via raw-map surgery (map[string]json.RawMessage), rewriting only the keys ctx owns; see internal/cli/initialize/core/merge/settings.go

---

## [2026-07-03-182238] New Claude model families silently fall to the 200k default in ModelContextWindow

**Context**: claude-fable-5 sessions warned '104 percent full' while the 1M window was 79 percent free; ModelContextWindow only recognized the [1m] suffix and the opus substring

**Lesson**: Unmapped model families inherit rc.DefaultContextWindow (200k), and every EffectiveContextWindow consumer (check-context-size hook, heartbeat, nudges, provenance) misreports against the wrong window

**Application**: At every model launch: add the family substring to internal/config/claude, extend TestModelContextWindow, and sanity-check against a live session's /context output before trusting hook percentages

---

## [2026-06-07-162840] Pin an on-disk contract before splitting work across parallel agents

**Context**: The ctx-dream skill (proposal writer) and the ctx dream CLI
(proposal reader) were built by parallel tracks; each independently invented the
proposals layout — one per-file proposals/<id>.json, the other a single
proposals.json array — requiring reconciliation at integration.

**Lesson**: When parallel agents share a serialized artifact (a file path +
schema), the exact on-disk shape must be fixed in BOTH prompts up front. 'Each
track invents it' guarantees a mismatch that only surfaces at integration.

**Application**: Before fanning out producer/consumer work to parallel agents,
write the precise file path + JSON schema into both agents' specs as a fixed
contract, and reconcile/verify it first thing at integration.

---

## [2026-06-07-162840] site/ is tracked build output — rebuild and bundle it with doc commits

**Context**: A docs change (docs/cli/dream.md, etc.) produced a surprise
189-file site/ drift mid-session; site/ is tracked (zensical build via 'make
site'), not gitignored.

**Lesson**: Any change under docs/ requires regenerating site/ with 'make site'
and committing the rebuilt site/ in the SAME commit (cf. f0f100a0, which bundled
its site rebuild). Otherwise the built output silently drifts and shows up as a
large untracked/modified set later.

**Application**: After editing anything under docs/, run 'make site' and 'git
add site/' and include it in the doc commit. Don't treat site/ as ephemeral —
it's versioned.

---

## [2026-06-07-142015] ctx-dream is headless-first; invoking /ctx-dream interactively is debugging, not the UX

**Context**: An end user ran /ctx-dream in their foreground terminal session and
watched the agent hand-execute the pass (grep/cat/hash/write JSON). Their words:
'I'm not dreaming but viewing a dream being debugged.'

**Lesson**: The dream is a sleep-time/headless product: cron 'claude -p' runs
the pass out-of-band, then the human is nagged and reviews via /ctx-serendipity.
The /ctx-dream SKILL is the executor's instruction set, not a user command —
driving it interactively makes the agent perform the executor's mechanical work
visibly, which is a debugging affordance, not the end-user experience. Exposing
/ctx-dream as a user slash-command invites exactly this confusion.

**Application**: End-user entry points are 'ctx dream' (on-demand pass that
prints a digest) and cron (scheduled); review is /ctx-serendipity. Do NOT have
users invoke /ctx-dream directly. Reconsider de-listing /ctx-dream from the
user-invocable skill set, ensure 'ctx dream' gives a clean run->digest
experience, and wire the 'serendipity round waiting' nag so the dream->review
loop closes without the user watching a pass.

---

## [2026-06-07-170001] ctx-dream design principles (consolidated)

**Consolidated from**: 6 entries (2026-06-06 to 2026-06-07)

- Merit/scoring rubric
  (relevance/frequency/recency/diversity/consolidation/richness, à la Hermes
  "Dreaming") measures ATTENTION (what to surface first), never TRUTH; use it
  only as a ranking signal feeding ruthless self-rejection, never as an
  autonomous promotion threshold — pair any statistical ranking with an
  evidence/grounding gate that decides eligibility.
- Load-bearing invariant (Option B): dream consolidation emits PROPOSALS only; a
  human accept/reject gate sits between the dream pass and any write to the five
  canonical files / MEMORY.md. Autonomous canonical writes are the documented
  rot failure mode (arXiv 2605.12978); independent designs (Hermes, OpenClaw,
  Auto-Dreamer) re-derive the sleep-phase shape but omit the gate. When
  evaluating any external memory-consolidation design, first check: does it
  autonomously write canonical, or only propose? Autonomous-write is a reject.
- A single LLM asked to critique a proposal silently repairs the missing
  justification and approves it (ReportLogic finding) — a single agreeable LLM
  is not an adversarial gate. Robust gating needs human or independent
  multi-critic consensus + swap-consistency. (This says a gate must EXIST; the
  proposes-only entry says one must sit before canonical writes; together they
  define WHO and WHETHER.)
- Same proposals, two consumers, two interfaces: render a terse/dispositional
  accept-reject worklist for the agent reviewer and a substance-rich,
  semantically-generated summary for the human (no file-hunting). Same data,
  presentation per consumer.
- Split agent/human work by comparative advantage: the agent is the reliable
  gardener for mechanical/verifiable hygiene (never skips the 47th file); the
  human owns taste/serendipity — which is WHY the human is the gate, not
  merely a safety nicety. Design the human's surface for pleasure (substance to
  wander), not a queue to drain.
- Don't-leak is a third safety axis alongside don't-corrupt and
  don't-obey-injected-instructions: a summary/backup/ledger-line of a gitignored
  source inherits its privacy class. Keep every byproduct in gitignored
  locations; enforce structurally with `git check-ignore` on each write target
  (refuse tracked paths), never via prompt. A deliberate human `promote` is the
  only sanctioned boundary crossing.

---

## [2026-06-07-170002] internal/audit & compliance gates for new code (consolidated)

**Consolidated from**: 6 entries (2026-03-15 to 2026-05-30)

- New exported types must live in types.go: TestTypeFileConvention permits types
  outside types.go only in pure-type files (defs+methods, no standalone funcs)
  or exempt packages; a file mixing structs with standalone funcs fails. Put
  type defs in a dedicated types.go from the start.
- internal/assets/tpl is on the magic-strings exempt list, so template-path
  literals are sanctioned THERE — but render data passed from non-exempt
  callers must be a typed struct (tpl.ObsidianData{...}), never map[string]any
  with literal keys, which trips the audit at the call site.
- Full gate catalog for a new package/CLI command (none surfaced by `go
  build`/`golangci-lint` — run `go test ./internal/audit/
  ./internal/compliance/`): TestNoMixedVisibility (split unexported helpers into
  <name>_internal.go), TestNoMagicStrings/Values (named consts in
  internal/config/warn/ for warn formats; named const for bare ints),
  TestDocCommentStructure (Parameters/Returns on every helper, exported or not),
  TestNoCmdPrintOutsideWrite (route output through internal/write/<area>/),
  TestNoNakedErrors, TestTypeFileConvention, TestCmdDirPurity (no helpers in
  cmd/ — use core/<area>/), TestNoLiteralMdExtension (file.ExtMarkdown),
  TestDocGoSubcommandDrift (parent doc.go lists every subcommand),
  TestDescKeyYAMLLinkage, TestNoLiteralWhitespace (token.NewlineCRLF/LF),
  TestRegistryCount (bump on registry.yaml additions). staticcheck QF1012 vs
  TestNoUncheckedFmtWrite: build with fmt.Sprintf then b.WriteString.
- naked_errors audit flags every fmt.Errorf/errors.New outside internal/err/**
  — call-site wrapping does NOT satisfy it. Error constructors live in
  domain-scoped internal/err/<area>/ pulling format strings from
  internal/config/<area>/ or desc.Text. Pattern: `var ErrX =
  errors.New(cfgArea.ErrMsgX)` (sentinel); `func X(args, cause) error { return
  fmt.Errorf(cfgArea.FormatX, …) }` (wrapper). Budget ~3 files/area for any
  new error surface.
- Pre-emptive constants are dead exports: TestNoDeadExports is
  symbol-graph-strict — any exported const/var/func without an internal reader
  fails. Land constants in the same commit (or strict precursor) as their
  caller; never scaffold config ahead of consumers. Genuine future-use goes in a
  TASKS.md line, not a config file.
- Dead-code detection: packages can build+test green while unreachable — check
  bootstrap registration, not build success (e.g. internal/cli/recall/ had
  tests, never wired). Files created by `ctx init` with no agent/hook/skill
  reader are dead on arrival. When touching legacy compat code, first ask if the
  legacy path has real users; if not, delete rather than improve (MigrateKeyFile
  had 5 callers, zero users).

---

## [2026-06-07-170003] Error handling: sentinels, unwrapping, and silent discards (consolidated)

**Consolidated from**: 6 entries (2026-03-06 to 2026-06-02)

- os.IsNotExist does NOT unwrap — it is false on any fmt.Errorf("…%w…")
  error; prefer errors.Is(err, os.ErrNotExist). But errors.Is only holds if the
  wrap carries %w at runtime, and a wrap whose format string comes from the
  text/i18n registry only carries %w when that registry is initialized (so it
  behaves differently in prod vs a bare test binary; go vet can't see it). To
  detect file absence reliably, stat directly: os.Stat returns an unwrapped
  *fs.PathError so errors.Is(statErr, os.ErrNotExist) is dependable everywhere.
- An error-discard catalogue (grep + name/regex classification) is an inventory
  of candidates, not findings. Name-inference produces false positives (a
  discarded bool mistaken for an error; a value type that can't nil-deref; an
  already-failed cleanup-close path). Read the callee signature and enclosing
  control flow before assigning return-error vs logWarn vs annotate.
- Canonical sentinel shape: a typed zero-data struct (or fielded struct for
  parameterised errors) whose Error() resolves text via
  desc.Text(text.DescKey…) lazily at call time — never `var ErrX =
  errors.New("english")` and never an ErrMsg* string-const layer. Empty-struct
  values are comparable and errors.Is finds them through %w wraps. Reference:
  internal/err/context/context.go.
- fmt.Fprintf to strings.Builder silently discards errors (Write never fails) so
  errcheck allows it, but project convention forbids any silent discard —
  TestNoUncheckedFmtWrite enforces `if _, err := fmt.Fprintf(...)`.
- A path-returning (string, error) function must never return ('', nil):
  filepath.Join('', rel) yields rel as a CWD-relative path, causing orphan
  writes at project root. Sentinel errors force callers to gate. Audit any
  path-returner with a historic ('', nil) shortcut (fixed: state.Dir,
  rc.ContextDir).
- Package-local err.go files in CLI packages invite agents to duplicate error
  constructors (errFileWrite, errMkdir repeated). Centralize in internal/err; no
  err.go files in CLI packages.

---

## [2026-06-07-170004] git CLI wrapping quirks (consolidated)

**Consolidated from**: 4 entries (2026-03-24 to 2026-05-22)

- `git rev-parse` exits 0 on an unknown long-flag and echoes the literal arg
  back as its only stdout line (treats it as a candidate revision name). A
  non-zero-exit guard never trips, so `--show-current` shipped verbatim into
  handover frontmatter. Validate the OUTPUT shape (length, no `--` prefix,
  hex-ness for SHAs) when wrapping rev-parse, not just the exit code.
  (`--show-current` is a `git branch` flag, not rev-parse.)
- Group git flag constants by the subcommand whose argv they're valid in (//
  Branch subcommand flags, // Rev-parse flags), not by "loose CLI flags" — the
  group comment is informal type info; mis-grouping enables wrong-subcommand
  bugs. Genuinely-spanning flags (-C, --) go under an explicit Cross-subcommand
  group.
- `git describe --tags --abbrev=0` follows reachability from HEAD, not the
  global tag list (diffed against v0.3.0 instead of v0.6.0 on a diverged release
  branch). For "latest release globally" use `git tag --sort=-v:refname | head
  -1`.
- A trailing regex word boundary \b does NOT exclude hyphenated continuations
  (\bgit commit\b matches `git commit-tree`). For porcelain with hyphenated
  cousins (commit-tree, commit-graph, for-each-ref) append a (?!-) negative
  lookahead.

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

## [2026-06-07-170006] Editorial KB pipeline: design epistemology (consolidated)

**Consolidated from**: 5 entries (all 2026-05-10)

- An ongoing user paying concrete workaround tax (disabled skills, hand-typed
  closeouts, colliding root constitution files) is the strongest validation
  evidence — beats user research, N=2 discussion, "seems useful." Use the
  workaround details as the inverse-spec; ship the shape they hand-rolled and
  use their project as the regression corpus.
- When lifting from a battle-tested external design, lift the renames and
  disambiguation moves alongside the features: intentional renames encode
  resolved conflicts (KB-RULES.md not CONSTITUTION.md; domain-decisions.md not
  DECISIONS.md). Treating them as cosmetic re-litigates the underlying fight.
- KB epistemology: a knowledge base has no "decide" moment — only
  evidence-capture events with confidence bands (>0.9 = decided by contract).
  Even NL assertions ("anchor on this") are evidence-capture, not
  decision-capture. So a parallel /ctx-kb-decide skill is the wrong shape; the
  pipeline-only-writer model is ontologically correct. General check: "I chose
  between alternatives" vs "I learned about the world."
- Recursive composability eliminates feature classes: a KB of KBs is a KB
  (source-map kind: kb + the standard ingest pipeline covers federation; no v1
  schema lockout). Ask whether the standard pipeline pointed at its own output
  covers a "thing-of-things" before designing a new mechanism.
- The LLM is the migration tool: every category of being-wrong about a schema
  (ID renumbering, taxonomy reshuffle, band remapping, path renames) is cheap
  because LLM cleanup absorbs the migration. Commit to the readable, opinionated
  v1 schema instead of hedging with abstract types; surface dirty state via
  doctor advisories so the agent has a work surface.

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

## [2026-06-07-170009] Constant placement & helper smells (consolidated)

**Consolidated from**: 6 entries (2026-03-07 to 2026-03-23)

- A constant used by only one domain (agent scoring, budget %, cooldowns)
  belongs in that domain's config package, not a god-object file.go. Check
  callers before placing.
- Before adding any constant to internal/config, grep by VALUE (".jsonl") not
  just name — camelCase vs ALLCAPS variants hide duplicates (ExtJsonl vs
  existing ExtJSONL).
- Project-root files created by `ctx init` (Makefile) are scaffolding
  (config/file), NOT context files loaded via ReadOrder (config/ctx). Check
  ReadOrder membership before moving a file constant.
- SafeReadFile / validation.SafeReadFile take (baseDir, filename) separately —
  split full paths with filepath.Dir + filepath.Base when adapting os.ReadFile
  calls.
- One-liner method wrappers that just forward a struct field to a stdlib/pkg
  function (checkBoundary → validation.ValidateBoundary with h.ContextDir)
  obscure the real dependency — inline them.
- A param-struct field that is a function pointer where all callers pass thin
  wrappers varying only by a text key (MergeParams.UpdateFn) is "data in
  disguise" — replace the callback with the key and let the consumer dispatch.

---

## [2026-06-07-170010] Convention enforcement: mechanical gates over prose (consolidated)

**Consolidated from**: 6 entries (2026-03-16 to 2026-04-14)

- System-level brevity instructions outcompete context-injected conventions;
  memory shifts probability (~40%→~70%) but doesn't create invariants. Invest
  in linter/PreToolUse gates for mechanically-checkable conventions; reserve
  behavioral nudges for judgment calls.
- Force-loaded behavioral prose (AGENT_PLAYBOOK at ~14k tokens) gets skipped
  when the user's first message is a concrete task; action-gating hooks
  (qa-reminder, specs-nudge) are followed because they fire at the moment of
  violation. More injected content = less attention per token. Prefer
  action-gating hooks; reserve force-injection for hard rules + distilled
  checklists.
- Any docstring/comment/documentation-formatting task is convention-sensitive:
  read CONVENTIONS.md (Documentation section) + LEARNINGS.md for known gaps
  FIRST, and audit all functions in scope against the template, not just diffed
  ones.
- AST audit tests must default to scanning ALL documented functions (use
  opt-outs not exported-only opt-ins) — TestDocCommentStructure missed
  unexported helpers (84 violations fixed). And the stutter test
  (TestNoStutteryFunctions) walks *ast.FuncDecl only, not GenDecl — stuttery
  const/var/type names slip through until the audit is extended.
- Every exemption map/allowlist in audit tests is a tempting agent shortcut: add
  DO-NOT-widen guard comments to every exemption data structure (10 across 7
  files) and review PRs for drive-by allowlist additions.

---

## [2026-06-07-170011] Go toolchain, gofmt & build-tag pitfalls (consolidated)

**Consolidated from**: 5 entries (2026-03-16 to 2026-05-10)

- gofmt strips bare `//` padding lines as unnecessary whitespace, so
  programmatic Go generation must produce substantive content lines; always run
  gofmt after any scripted Go-file generation.
- Agents reliably introduce gofmt issues during bulk renames (75+ files, 12
  broken); run `gofmt -l` (then `-w`) as a standard step after any agent-driven
  bulk edit before trusting the build.
- The "compile version X does not match go tool version Y" error comes from the
  CACHED toolchain (~/go/pkg/mod/golang.org/toolchain@…), not the system Go
  — reinstalling Go does nothing. Diagnose via `go env GOROOT`; fix by
  deleting the cached dir, bumping go.mod, or GOTOOLCHAIN=go<system>. `go clean
  -cache` and GOTOOLCHAIN=local don't help.
- `make test` exit code is unreliable: the -cover flag can fail with "no such
  tool covdata" even when every package passes. Fall back to `go test ./...` (no
  -cover) and tally ^ok/^FAIL.
- AST checks via go/packages only see files matching the current GOOS —
  darwin-only (_darwin.go) violations are invisible on Linux. Fix violations
  regardless; note coverage is platform-dependent (need multi-GOOS CI or a
  go/parser fallback).

---

## [2026-06-07-170012] Stale-task triage & verify-before-acting (consolidated)

**Consolidated from**: 4 entries (2026-03-01 to 2026-05-23)

- Stale TASKS.md items often describe work already done in code but not asserted
  in tests — the task stayed open because nothing pinned the behavior. Triage
  older items by grep/git-blame on the named symbols; if implemented, close by
  writing the regression test (often one function). Applies to behavior-named
  tasks more than feature-named ones.
- Tasks can be stale in reverse: implementation completed but task not marked
  done (recall sync was fully wired despite a "not registered" description). Run
  `ctx <cmd> --help` before assuming work remains.
- Grep for callers must cover the ENTIRE working tree before deleting functions
  — with unstaged changes from a prior session, grep hits only
  committed+staged code. Always `make build` after deleting functions even when
  grep shows zero callers.
- Spec-trailer improvisation is heuristic drift: when no on-topic spec exists,
  the path of least resistance cites the most-recent spec from context,
  satisfying the syntactic gate but defeating truthful traceability — and
  session-scoped "I'll be careful" commitments don't survive across sessions, so
  the fix must live in persistent context. Correct responses: scaffold a fresh
  spec, bundle into the next functional commit, or cite specs/meta/chores.md.
  (See specs/spec-trailer-discipline.md; AGENT_PLAYBOOK Spec Verification Step.)

---

## [2026-06-07-170013] Refactor mechanics: subagents, cascades & golden fixtures (consolidated)

**Consolidated from**: 6 entries (2026-02-19 to 2026-05-30)

- Behavior-preserving refactors of formatting/rendering code: capture golden
  fixtures from the LIVE legacy path before deleting it (throwaway test writes
  testdata/*.golden), then assert byte-equality after — avoids silent drift
  from hand-transcribing expected output.
- Removing a sentinel (ErrDirNotDeclared) cascades through ~10 errors.Is
  consumers and ~30 test fixtures; spec-level step boundaries that separate
  "swap resolver" from "remove guard" don't survive when the second references
  the soon-deleted sentinel. Plan the merged commit at spec time; do the
  compile-surface analysis then.
- Subagent parallelism shines for well-bounded mechanical refactor WITH a
  canonical worked example on disk and an explicit fix-or-fail-with-a-blocker
  instruction (invoke the no-deferral rule). Do one worked example in the
  orchestrator, then dispatch subagents pointing at it.
- Subagents reliably exceed scope (rename funcs, change signatures, restructure
  files even for em-dash fixes) and create new files without deleting originals.
  After any agent refactor: `git diff --stat`, `git diff --name-only HEAD`,
  revert out-of-scope changes, check for stale package decls/duplicate
  defs/orphaned imports, run gofmt + `go test ./...`.
- Splitting a flat core/ package into subpackages exposes duplicated logic,
  misplaced types, and function-pointer smuggling invisible in the flat layout;
  circular-dep resolution during the split IS the design work that reveals the
  right structure.
- Cross-cutting change ripple: path/asset/feature changes ripple across 15+ doc
  files + multiple layers (embed directive, accessors, callers, tests, config
  consts, build targets, docs). Grep broadly (not just code); a feature without
  docs (feature page, cli-reference, recipes, nav) is invisible.

---

## [2026-06-07-170014] Linting, gosec & I/O chokepoints (consolidated)

**Consolidated from**: 4 entries (2026-01-25 to 2026-04-03)

- Full pre-commit gate, every time: (1) CGO_ENABLED=0 go build ./cmd/ctx, (2)
  golangci-lint run, (3) CGO_ENABLED=0 go test. Own the codebase — fix
  pre-existing lint issues you didn't introduce.
- gosec permissions: 0o600 for files (incl. tests — G306 flags 0644 even in
  test code), 0o750 for dirs (G301); G304 file-inclusion is safe to
  //nolint:gosec in tests using t.TempDir(). Prefer renaming constants to avoid
  G101 false positives (Tokens→Usage, Passed→OK) over nolint/nosec/path
  exclusions, which break on file reorg.
- Suppression anti-patterns: nolint:goconst normalizes magic strings (use config
  consts); nolint:errcheck in tests teaches agents to spread the pattern to
  production (use t.Fatal for setup, `defer func(){ _ = f.Close() }()` for
  cleanup); golangci-lint v2 ignores inline nolint for some linters — use
  config-level exclusions.rules for gosec, fix the code for errcheck. Use
  cmd.Printf/Println in Cobra commands instead of fmt.Fprintf. `defer
  os.Chdir(x)` fails errcheck — wrap in `defer func(){ _ = os.Chdir(x) }()`.
  CI Go-version mismatch: install-mode goinstall.
- Chokepoint migrations have cascading benefits: centralizing file I/O into
  internal/io/ (already using config/fs consts) zeroed out TestNoRawPermissions
  for free. Prioritize chokepoint migrations (io, exec, write, err) before
  smaller dependent checks.

---

## [2026-06-07-170015] Hook mechanics, output channels & compliance (consolidated)

**Consolidated from**: 5 entries (2026-01-25 to 2026-04-06)

- Hook scripts receive JSON via stdin (HOOK_INPUT=$(cat) then jq), not env vars;
  key names are case-sensitive (PreToolUse, SessionEnd); use
  $CLAUDE_PROJECT_DIR, never hardcode paths; anchor regex to command-start
  `(^|;|&&|\|\|)\s*` ('ctx' binary vs dir); grep matches inside quoted args
  (test with blocked words); scripts silently lose execute permission (verify ls
  -la).
- Output routing: plain-text hook stdout is silently ignored — Claude Code
  parses stdout starting with `{` as JSON directives; return JSON via
  printHookContext(). For UserPromptSubmit specifically, stdout is prepended as
  AI context (not user-visible), stderr+exit0 is swallowed, user-visible output
  requires {"systemMessage":"…"} or exit 2 (blocks); there is NO non-blocking
  user-visible channel. Two-tier severity is sufficient: unprefixed (agent
  context, may relay) and "IMPORTANT: Relay VERBATIM" (guaranteed); don't add
  more prefixes.
- Agents only relay content with explicit display instructions: a
  system-reminder line with no "Display this line verbatim" is invisible to the
  user even when correct. IMPORTANT: signals internal priority, not user-facing
  output.
- Compliance: soft instructions have a ~75–85% ceiling because "don't apply
  judgment" is itself judgment; for 100% compliance inject via additionalContext
  rather than instruct. Hook compliance degrades on narrow mid-session tasks
  (~15–25% skip) because CLAUDE.md's "may or may not be relevant" competes
  with hook authority — fix by elevating hook authority explicitly; the
  mandatory checkpoint relay block is the compliance canary. No reliable
  agent-side before-session-end event exists (SessionEnd fires after the agent
  is gone) — mid-session nudges + explicit /ctx-wrap-up are the only reliable
  persistence. Repeated injection causes repetition fatigue — gate with
  --session $PPID --cooldown and pair with a readback instruction.
- Context-budget injection strategy: once ~7K tokens are auto-injected (fait
  accompli), the agent's rationalization inverts from "skip to save effort" to
  "marginal cost is trivial." Front-load highest-value content as injection,
  then leverage sunk cost for on-demand reads. Verbal summaries + linked diagram
  files cut ARCHITECTURE.md ~12K→3.8K (extract diagrams outside FileReadOrder;
  the 4-chars/token estimator is accurate — optimize content not the
  estimator).

---

## [2026-06-07-170016] State, tombstones, logs & filesystem hygiene (consolidated)

**Consolidated from**: 6 entries (2026-02-11 to 2026-03-06)

- Permission drift is distinct from code drift — settings.local.json is
  gitignored so no review catches stale entries; it accumulates session debris
  (run /sanitize-permissions + /ctx-drift). Skill() permissions don't support
  name-prefix globs (list each); wildcard trusted binaries (Bash(ctx:*),
  Bash(make:*)) but keep git granular (never Bash(git:*)).
- Gitignored directories are invisible to git status — stale artifacts persist
  indefinitely (periodically ls them). Add editor artifacts (*.swp,*.swo,*~) to
  .gitignore from day one. Gitignore entries for sensitive paths are security
  controls, not documentation — never remove during cleanup.
- The state directory accumulates write-only session tombstones and grows
  unbounded without auto-prune (234 files found); autoPrune(7) now runs once per
  session at startup via context-load-gate (manual `ctx system prune` still
  available).
- A session-scoped tombstone must include the session ID in its filename, else
  it suppresses hooks across ALL concurrent and future sessions (memory-drift
  fixed; backup-reminded, ceremony-reminded, check-knowledge, journal-reminded,
  version-checked, ctx-wrapped-up still carry this bug). Use the UUID pattern so
  prune can clean them.
- New log sinks must follow the established rotation pattern (size-based, single
  previous generation): eventlog rotated at 1MB but logMessage() in state.go was
  append-only with no size check.
- If a directory is recreated (auto-prune), an SSH shell holding the old inode
  won't see new files (ls returns "no such file" though cat with the full path
  works elsewhere); after `ctx system prune` or any state recreation, SSH
  sessions need cd-. or re-login.

---

## [2026-06-07-170017] Host-pressure alerting: use derivatives, not levels (consolidated)

**Consolidated from**: 2 entries (2026-04-13 to 2026-05-28)

- Swap occupancy is NOT memory pressure: macOS/Windows swap proactively and
  occupancy is a sticky high-water mark that doesn't recede when pressure ends,
  so any alert keyed on SwapUsed/SwapTotal ≥ X% false-positives at session
  start (e.g. after hibernation). Key on OS-native pressure derivatives instead:
  macOS kern.memorystatus_vm_pressure_level (1/2/4 → OK/Warning/Danger), Linux
  PSI /proc/pressure/memory some.avg10/full.avg10; fall back to swap-out RATE
  gated on low available memory, never occupancy.
- Load average measures a queue (runnable + uninterruptible-sleep), not CPU
  utilization — high load with low CPU% means many short-lived/I/O-bound
  processes (e.g. go test spawning hundreds of binaries). For automated alerts
  prefer the 5-minute average over the reactive 1-minute, which fires on normal
  build/test activity.

---

## [2026-06-07-170018] Go test isolation & patterns (consolidated)

**Consolidated from**: 4 entries (2026-01-19 to 2026-04-25)

- Any code using os.UserHomeDir() / user-level paths (~/.ctx/, ~/.config/) needs
  t.Setenv("HOME", tmpDir) in tests — especially shared setup helpers. Under
  parallel `make test`, fourteen test files invoking initialize.Cmd().Execute()
  raced on read-modify-write of ~/.claude/settings.json, surfacing as flaky
  "FAIL coverage: [no statements]"; testctx.Declare now sets HOME alongside
  CTX_DIR (centralized fix).
- Go testing patterns: `go build ./...` misses test-file callsite breaks —
  always `go test ./...` after signature changes. Consume all runCmd() returns
  (`_, _ = runCmd(...)`) for errcheck. Disable ANSI via color.NoColor=true in
  package init for string assertions. Recall tests isolate via t.Setenv("HOME",
  tmpDir) with .claude/projects/. formatDuration takes an interface with
  Minutes() (use a stubDuration). CI needs CTX_SKIP_PATH_CHECK=1 (init checks
  PATH). CGO_ENABLED=0 for ARM64 Linux.
- Converting PersistentPreRun → PersistentPreRunE changes exit behavior:
  errors propagate through Cobra Execute() return with no os.Exit.
  Subprocess-based tests expecting exit codes must convert to direct error
  assertions.

---
## [2026-06-01-174927] Guard managed blocks before regenerating; don't trust the span to be machine-owned

**Context**: ctx learning add silently deleted entry bodies that lived between
INDEX:START/END markers: index.Update replaced the whole marker span with a
regenerated table, and ParseHeaders scanning the full file made the result look
complete, hiding the loss.

**Lesson**: Code that 'replaces the managed block' (index regen, KB managed
blocks, moc.go) assumes the span between its markers is disposable and
machine-owned. That assumption breaks the moment user content drifts inside the
markers, and the regenerated output looks correct so the loss is invisible. The
fix is a precondition guard that refuses to mutate when regeneration would lose
data — not smarter parsing of the trapped content.

**Application**: Before any 'replace between markers' write, validate the span:
refuse on entry/content found where only generated output belongs, and on
malformed/duplicated/out-of-order markers. Fail loud and leave the file
byte-identical rather than regenerate. Run the guard at the read-before-mutate
choke point so nothing is written on refusal.

---

## [2026-05-31-094649] macOS GUI apps inherit a minimal PATH; augment it to find a user-installed CLI

**Context**: A bundled Tauri app launched via Finder/launchd gets a minimal PATH (/usr/bin:/bin:...), so /usr/local/bin/ctx is not found even though it resolves in a terminal-launched dev run.

**Lesson**: Do not rely on inherited PATH when spawning user-installed CLIs from a desktop GUI.

**Application**: ctx_adapter prepends /usr/local/bin:/opt/homebrew/bin to PATH on every std::process::Command invocation.

---

## [2026-05-31-094649] Tauri 2 requires rustc >= 1.88; bump the toolchain before cargo check

**Context**: cargo check for the ctx-desktop Tauri app failed: darling, serde_with, time and plist transitive deps require rustc 1.88, but the local toolchain was 1.87.

**Lesson**: Tauri 2's dependency tree tracks recent rustc releases; the pinned-stable assumption breaks builds.

**Application**: Run rustup update stable before building a Tauri 2 app; this project moved 1.87 -> 1.96.

---

## [2026-05-28-215214] ctx kb: single topic-enumeration site; life-stage count is consumer-side

**Context**: kb reindex blanked the CTX:KB:TOPICS block for grouped kbs
(things-wtf-dr regrouped 49 topics into folders); the task speculated a sibling
life-stage topic-count glob was also affected.

**Lesson**: reindex.ListTopics (internal/cli/kb/core/reindex/topic.go) is the
ONLY topic enumeration/count in ctx, and CTX:KB:TOPICS is the only managed
block. The life-stage concept in ctx is the ingest/closeout frontmatter field,
unrelated to topics. Any per-life-stage topic count lives in the consumer kb,
which ctx neither generates nor owns.

**Application**: Localize nested-topic fixes to ListTopics; treat
per-group/per-life-stage topic counts as consumer territory (same recurse +
exclude-group-landing pattern, fixed in their repo).

---

## [2026-05-28-201400] A non-root Go module nested under the main module's path CAN import its internal/ packages

**Context**: While designing the ctxctl module split, the initial spec (and a
lot of online consensus) claimed a separate `go.mod` cannot import the parent
module's `internal/` packages, which would have forced relocating or duplicating
~25 foundation packages (`rc`, `desc`, `nudge`, `config/*`, …). The "obvious"
reading made same-module the only viable option.

**Lesson**: Go's internal-import rule is **lexical on import paths, not
module-scoped**. A separate module whose path is
`github.com/<owner>/<main>/tools/<x>` CAN import
`github.com/<owner>/<main>/internal/...` — verified by an empirical build
experiment this session. An outsider path (`example.com/...`) is rejected with
`use of internal package … not allowed`. The rule fires on the import-path
prefix relative to the `internal/` directory's parent, not on module boundaries.

**Application**: For monorepo splits (maintainer-only tooling, isolated
experiments, ancillary CLIs), choose a module path nested under the main module
so the new module reuses the parent's foundations via the lexical-internal
allowance. Full self-containment of a maintainer module would be a DRY
catastrophe; the lexical allowance is the correct shape. Prove it with a
throwaway `go build` against a representative `internal/` import before
designing around the *wrong* constraint.

---

## [2026-05-28-201300] cobra's legacyArgs lets unknown subcommands silently succeed on non-root groups

**Context**: Every prompt of this session injected 52 lines of `ctx system` help
text into agent context, labeled "hook success." Investigation traced it to the
0.8.1 plugin's `hooks.json` wiring `ctx system check-anchor-drift` as the first
UserPromptSubmit hook — a command the 0.8.1 binary no longer has (the command
was deleted by the cwd-anchored migration in `fc7db228`, but the plugin's hook
config wasn't updated). The harness reported "hook success" because cobra exits
0 on the unknown subcommand.

**Lesson**: cobra's `legacyArgs` only raises "unknown command" for the **root**
command (`!cmd.HasParent()`); any non-root group (built with `parent.Cmd`)
treats an unknown subcommand as non-error: it falls through to `Help()` and
returns nil → exit 0. In a UserPromptSubmit hook this is **invisible** — the
harness logs "hook success" and injects the whole help text into agent context
every prompt. The 0.8.1 plugin's stale wiring of the retired
`check-anchor-drift` caused exactly this for the entire session.

**Application**: Non-root cobra groups must have an explicit unknown-subcommand
guard. Two routes: (a) `Args: cobra.NoArgs` so unknown subcommands error loud
(non-zero exit + "unknown command" stderr); (b) a `RunE` that emits a **verbatim
relay** — which is what actually reaches the user in a UserPromptSubmit hook
context where a non-zero exit alone is invisible. Tracked under Phase CLI-FIX as
the verbatim-relay guard on `ctx system`.

---

## [2026-05-25-221357] Skill shipping location: _ctx- prefix is repo-internal, internal/assets/claude/skills/ctx-* is bundled and shipped

**Context**: Created /ctx-surface-audit under internal/assets/claude/skills/
(the shipped path), but it audits ctx's own internal/ source layout — useless
in an end-user project that installs ctx. There is an established _ctx-* family
(_ctx-command-audit, _ctx-audit, _ctx-release, _ctx-qa, etc.) in .claude/skills/
for repo-only dev skills; the user caught the misplacement.

**Lesson**: A skill that references ctx's own source tree (internal/,
docs/recipes/, cmd/) or dev workflow is repo-internal and must live in
.claude/skills/_<name>/ (underscore prefix, committed to the repo but NOT
bundled). Only genuinely user-facing skills belong in
internal/assets/claude/skills/, which ctx init / ctx setup install into end-user
projects. The same ship-vs-repo-internal question applies one layer up:
user-facing CLI commands go in ctx, maintainer commands go in ctxctl; shipped
hooks live in internal/assets/claude/hooks/hooks.json and call ctx, repo-local
dev hooks live in the gitignored .claude/settings.local.json and may call
ctxctl.

**Application**: Before creating a skill, command, or hook, ask: does this serve
a user working in their project, or a ctx maintainer working in this repo?
Maintainer-facing → _-prefixed skill in .claude/skills/ + ctxctl command +
repo-local hook. User-facing → internal/assets/claude/skills/ + ctx command +
shipped hooks.json. Putting maintainer tooling in the shipped paths taxes every
end user (e.g. a UserPromptSubmit hook firing on every prompt for a feature they
never use).

---

## [2026-05-23-001000] Unicode block separation makes diacritic-stripping surgical — no per-script handling needed for Arabic/Indic/Hebrew/CJK

**Context**: While building `i18n.MatchKey` (commit 978582f5) for
diacritic-insensitive placeholder matching, the natural reflex was "this is
going to need per-script special cases — CJK doesn't have case, Arabic has
shadda/fatha that are meaning-changing, Bengali vowel signs are
script-essential, Hebrew niqqud distinguishes words." I sized the work assuming
we'd need a script-aware policy, possibly with a locale config or an opt-in flag
for "strip all combining marks" vs "strip only Latin-style decoration".
Empirical test across Turkish/German/French/Spanish/Catalan/Czech/Vietnamese
(should collapse) and Arabic/Bengali/Devanagari/Hindi/Hebrew/Chinese/Korean
(should preserve) showed the entire policy fits in one numeric range:
U+0300..U+036F.

**Lesson**: Unicode pre-separated combining marks by intent at the codepoint
level. The "Combining Diacritical Marks" block (U+0300–U+036F) holds
Latin/general decorative marks: acute, grave, diaeresis, tilde, cedilla, caron,
the Turkish combining dot, the Vietnamese horn, etc. Script-essential marks live
in separate blocks per script: Arabic in U+0610–U+06ED, Bengali in
U+0980–U+09FF, Devanagari in U+0900–U+097F, Hebrew niqqud in
U+0591–U+05C7, and so on. The block boundaries are not coincidental — they
encode the same distinction a reasonable design would want to make. So a narrow
byte-range strip is exactly the right primitive: it expresses "remove
decoration, keep structural marks" in one comparison, without needing to know
anything about the input's script.

**Application**: When designing comparison/normalization primitives for
international input, check the Unicode block boundaries before reaching for
per-script special cases or a config field. Often the standardization committee
already drew the line you want, and an arithmetic range check (`r >= 0x0300 && r
<= 0x036F`) does the work. Verify empirically across the scripts you care about
— but expect the answer to be cleaner than your initial sizing. The general
rule: when Unicode has put related characters in their own block, treat that
block as a meaningful unit of policy. (For ctx, this is now
`cfgI18n.CombiningMarksLatinStart`/`End` and the `MatchKey` implementation in
`internal/i18n/matchkey.go`.)

---

## [2026-05-20-214839] macOS /var symlink trips path-equality; use EvalSymlinks with parent-resolution fallback

**Context**: TestRunInit_EnvCwdMatch_Succeeds in
internal/cli/initialize/init_test.go failed on first run despite a deliberate
setup where the env path and cwd candidate matched. Diagnosis: t.TempDir()
returns paths like /var/folders/..., os.Getwd() after t.Chdir() returns the
canonical /private/var/folders/... (because macOS's /var is a symlink to
/private/var). filepath.Clean preserves the symlink form; equality fails.

**Lesson**: filepath.Clean alone is insufficient for path equality on macOS (and
other systems with symlinked top-level dirs). filepath.EvalSymlinks resolves the
symlinks but fails when the target path does not yet exist — common case for
/Users/volkan/Desktop/WORKSPACE/ctx/.context BEFORE ctx init runs. The right
pattern is a layered fallback: try EvalSymlinks(full), then EvalSymlinks(parent)
+ rejoin basename, then filepath.Clean as last resort.

**Application**: Encapsulated as
internal/cli/initialize/core/envmatch/{envmatch.go,internal.go}. The Same(a, b)
public function calls resolve() on each side; resolve() tries EvalSymlinks on
the full path, falls back to EvalSymlinks on the parent (rejoining the
basename), and falls through to filepath.Clean if both fail. Reusable for any
future env-vs-cwd-style equality check. The package is per-feature
(core/envmatch/) per the cmd/core/ purity rule enforced by
internal/compliance/TestCmdDirPurity.

---

## [2026-05-20-214830] Handover filenames are archaeology; parse by generated-at, not filename

**Context**: User observed three coexisting handover filename shapes:
.context/HANDOVER-2026-04-22.md (pre-skill root file),
.context/handovers/YYYY-MM-DD-HHMMSS-slug.md (skill-era pre-CLI),
.context/handovers/<RFC3339Compact>-slug.md (current CLI). User asked whether
this was a regression or a skill-interpretation problem.

**Lesson**: Neither. The .context/HANDOVER-* root file predates the handovers/
directory contract entirely (the body even said 'delete this file after
reading'). The YYYY-MM-DD-HHMMSS shape was an earlier skill iteration writing
free-form before ctx handover write existed (commit 60543e46, 2026-05-17,
introduced the CLI as sole writer per the anti-pattern note in /ctx-handover
SKILL.md). The current parser at internal/write/handover/parse.go:75-107 keys on
the 'generated-at' YAML frontmatter, not the filename — so legacy shapes still
sort correctly via LatestHandoverCursor. Only files without frontmatter (the
root April file) are invisible.

**Application**: When unifying filename shapes across history, use git mv to
preserve rename detection. Derive the canonical timestamp from the file's own
generated-at frontmatter rather than from the filename — that's the source of
truth the parser uses anyway. If a handover predates frontmatter entirely (rare,
pre-skill era), it's safe to delete because the parser never read it.

---

## [2026-05-20-214821] /ctx-plan is named after its input, not its output

**Context**: Agent (and apparently other agents in prior sessions per user
observation) repeatedly inverted the canonical chain, treating /ctx-spec as the
entry point and /ctx-plan as a post-spec step. The skill description starts
'stress-test a plan' (implying user brings a plan IN) while line 44 of the body
says 'the deliverable is a debated brief, not a task list' (the OUTPUT is a
brief, not a plan).

**Lesson**: Skill names that reference their INPUT bias the agent toward the
wrong canonical position. The /ctx-plan skill takes a plan and produces a brief;
the natural mental model when scanning the name is 'plan = output', which makes
the agent place it AFTER spec instead of before. Also: /ctx-spec's 'When to Use'
section listed /ctx-brainstorm as a predecessor but never /ctx-plan, so an agent
skimming the top of the skill never learned the full chain.

**Application**: Made the canonical chain explicit at the top of both /ctx-plan
and /ctx-spec skills (Canonical Chain block with the brainstorm → plan →
spec → implement diagram) and in AGENT_PLAYBOOK_GATE Planning Work section.
/ctx-spec When-to-Use now lists /ctx-plan as a predecessor; When-NOT-to-Use says
'when the bet is contested but not yet stress-tested, use /ctx-plan first'.
/ctx-plan description now ends with '; produces a debated brief at
.context/briefs/<TS>-<slug>.md that /ctx-spec --brief consumes'.

---

## [2026-05-17-200000] Creator confusion is the strongest doc-quality signal — louder than any user signal

**Context**: In this session the project author asked *"why
external sources only? I can ground on a repo, a MCP query, a
markdown I dropped into ./inbox — are they also considered
'external'. Or is there a nomenclature confusion here?"* — and
explicitly noted *"it is confusing to the very creator of this
pipeline. -- and that's not a good sign."* Investigation
confirmed the input contract accepted in-tree paths and MCP
resources all along, but the SKILL.md ledes, the CLI docs table
row, and the recipe all framed ground as "external" — which
the creator's own mental model couldn't reconcile with the
contract.

**Lesson**: A normal-user reading-confusion signal is "I don't
understand this." A creator reading-confusion signal is "this
contradicts what I built." The second is louder by an order
of magnitude — the creator has the full internal model and a
strong prior on what the system should say. If they trip over
the words, the words are wrong, full stop. Don't defend the
existing framing; don't explain what was meant. Rewrite to
match what the contract actually does. The creator was a
control instrument; if even that instrument deflected, the
docs are mis-anchoring everyone.

**Application**: When the project's own creator asks a
"do we even need X?" or "wait, isn't X actually doing Y?"
question, treat it as a doc-bug report, not an architecture
question. Investigate the literal contract (input/output
shapes, code-level reality) before debating semantics. If the
contract is correct but the docs misframe it, the action is
"rewrite the framing across every doc surface that touches
it" — skills, recipes, CLI tables, anything user-facing.
Concrete instance handled this session: dropped "external"
from ctx-kb-ground's prose, description, pass-mode value,
recipe Step 4, and CLI table row across all three skill trees.

---

## [2026-05-17-061500] `_helpers.go` / `_utils.go` filenames are project anti-pattern; use domain nouns

**Context**: During Phase KB / Phase RG audit cleanup, the first file split I
attempted to satisfy the mixed-visibility audit named the new file
`read_helpers.go`. The user vetoed on sight: "utils; helpers, etcs are ALL lazy
naming; I will veto them the moment I see them; find proper domain objects."

**Lesson**: ctx's per-package file layout follows domain nouns, not
visibility-suffixed catch-alls. The canonical reference shape is
`internal/journal/parser/` which splits 18 files by domain (envelope, markdown,
parse, validate, claude, copilot, ...). The mixed-visibility audit demands a
split, but the split target must be a real noun: `frontmatter.go` (YAML
parsing/validation), `markdown.go` (rendering), `filename.go` (filename
derivation), `provenance.go` (sha/branch resolution), `parse.go` (one-shot
parser), `cursor.go` (latest-pointer logic). Never `_helpers.go`.

**Application**: When splitting a file to satisfy `mixed_visibility_test`, name
the new file for what the helpers ARE about, not for what visibility they have.
If you can't name it cleanly, the split itself may be wrong and the funcs may
belong in a different package entirely.

---

## [2026-05-11-231025] Naive Markdown line-sweep corrupts multi-line code spans and YAML lists

**Context**: Performed a programmatic typographic sweep across docs/*.md to wrap
bare 'ctx' tokens in backticks (commit 61aab858). 81 source files, 236 lines
changed. First pass corrupted two indented JSON snippets in MkDocs admonitions
because the fence regex anchored to start-of-line and missed admonition-indented
fences. After fixing the fence regex, two more corruptions surfaced (multi-line
inline-code spans where the opening backtick is on line N and the closing on
line N+1: the line-at-a-time transformer treated each line independently,
leading to misjudged span boundaries on the second line). After post-sweep
validation, a YAML parse error on docs/blog/2026-02-03-the-attention-budget.md
surfaced one more breakage: a 'topics:' list-item like '- ctx primitives' got
wrapped to '- `ctx` primitives', which is invalid YAML (a value starting with
backtick is not a valid unquoted scalar). Total: 2 multi-line span corruptions +
1 YAML breakage, all detected only by post-sweep validation (make site + grep
audit), not by the dry-run.

**Lesson**: A naive line-at-a-time regex sweep across Markdown documents must
respect a wider 'skip' set than the obvious cases. The full safe-skip list is:
(1) triple-backtick fenced code blocks, BOTH root-level and indented inside
MkDocs admonitions or list items (fence regex must allow leading whitespace,
e.g. '^\\s*```'); (2) inline backtick spans on the same line; (3) multi-line
inline-code spans crossing line boundaries (line-at-a-time logic cannot detect
both ends, so either track fence-like 'odd-count' state across lines or treat
any unclosed-on-line backtick as 'protect rest of line'); (4) the ENTIRE YAML
frontmatter block (delimited by '---' at top and next '---'), not just specific
keys like title/description/icon, because list-item values under
topics/tags/keywords are also YAML and break on a leading backtick; (5) image
alt-text '![alt]' (alt-text does not render in monotype); (6) link-reference
definitions '[name]: url "title"'; (7) project copyright header comment blocks.
Dry-run output never catches YAML or multi-line span breakage; validation MUST
include a parser-level check (make site for YAML, post-grep for '``name`'
double-backtick patterns near the wrapped token).

**Application**: When designing any future programmatic sweep across docs/
(typography passes, internationalization, brand renames, em-dash replacement,
link-text rewrites): (1) implement the full skip set above, not a subset; (2)
for fence detection use '^\\s*```', not '^```'; (3) for the frontmatter, detect
the entire block between '---' delimiters, not specific keys; (4) for multi-line
inline-code, choose between cross-line backtick-pair tracking (complex but
correct) or the simpler 'unclosed backtick protects rest of line' heuristic
(corrupts ~1 per 100 files but recoverable manually); (5) ALWAYS validate
post-sweep with 'make site' (zensical surfaces YAML errors) and a grep for
'``\\w' double-backtick patterns near the wrapped token; (6) commit only after
both validations are clean. For one-shot sweeps the script can be ad-hoc, but
record the validation gate as part of the commit message so the next contributor
knows what to check.

---

## [2026-05-08-195031] Cursor imports Claude Code hooks and sets CLAUDE_PROJECT_DIR per workspace

**Context**: Investigating why .context/state/ appeared in non-ctx projects
opened in Cursor. Hypothesis was a Cursor extension or shell hook; turned out to
be Cursor's documented Claude-compatibility behavior
(https://cursor.com/docs/hooks): it loads ~/.claude hooks and injects
CLAUDE_PROJECT_DIR=workspace_root so they 'just work'. Globally-enabled Claude
plugins therefore fire in every Cursor workspace.

**Lesson**: When debugging cross-tool side effects, check whether the host tool
advertises compatibility with the implicated tool's config format. The leak
surface of any global Claude plugin is now 'every Cursor workspace + every
Claude Code project', not just 'every Claude Code project'.

**Application**: Hooks must be safe to fire in non-ctx projects: silent bail
when state.Initialized() is false, no filesystem side effects. The ctx code-side
fix lives in state.Dir's Initialized gate; the design rule is broader: assume
hooks may run anywhere, not just where the user invoked ctx init.

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

## [2026-04-14-010134] Raft-lite trade-off is the load-bearing choice in internal/hub

**Context**: Discovered while writing thorough doc.go for internal/hub. The
package embeds HashiCorp Raft for leader election only; data replication is
sequence-based gRPC sync over the append-only JSONL store.

**Lesson**: A leader crash window between accept and replicate can lose the most
recent write. Append-only storage plus idempotent clients make this acceptable;
full Raft log replication would not be needed and would not be simpler.

**Application**: Any future make hub stronger proposal must engage with this
trade-off explicitly. Do not abandon Raft-lite accidentally by introducing
log-replicated state; that would invalidate the simplicity argument.

---

## [2026-04-14-010105] Brand-name handling in title-case engines must cover possessives

**Context**: First pass of hack/title-case-headings.py produced 'Ctx's' from
'ctx's' because the brand check matched the bare token only.

**Lesson**: A brand allowlist needs to recognize <brand>, <brand>'s, <brand>s,
and short apostrophe-suffixed variants. Single-word matching misses contractions
and possessives.

**Application**: When adding a new always-lowercase brand to
hack/title-case-headings.py, extend the suffix-aware loop in title_case_word,
not just the BRAND_LOWER set.

---

## [2026-04-13-153618] GPG signing from non-TTY contexts requires pinentry-mac (or equivalent)

**Context**: git commit failed from Claude Code's shell with 'gpg: signing
failed: No such file or directory' — the default pinentry-curses cannot open a
TTY in agent-invoked shells. Manual commits from a real terminal worked fine.

**Lesson**: GPG's default curses pinentry requires an interactive TTY. In
non-TTY contexts (Claude Code, CI, scripts, cron), signing fails silently-ish.
The fix is to configure a GUI pinentry that uses the OS keychain: brew install
pinentry-mac; echo 'pinentry-program $(brew --prefix)/bin/pinentry-mac' >>
~/.gnupg/gpg-agent.conf; gpgconf --kill gpg-agent. Once the passphrase is saved
in Keychain, signing works from any context.

**Application**: If agents or CI need to sign commits, configure pinentry-mac
(macOS) or pinentry-gtk/pinentry-qt (Linux) with the OS keychain, not
pinentry-curses. This is a one-time setup per machine.

---

## [2026-04-13-153618] rc.ContextDir() is the single source of truth — fix the resolver, not callers

**Context**: When ctx init failed with a boundary error, my first instinct was
to have init bypass rc.ContextDir() and use filepath.Join(cwd, dir.Context)
directly. Volkan shut that down: rc.ContextDir() encodes invariants (team
shares, symlinks, network mounts, .ctxrc overrides) that individual commands
cannot reason about.

**Lesson**: Resolution chains with multiple fallbacks are contracts. If one
command bypasses the chain, it silently diverges from every other command's
notion of 'the context directory.' When a resolver produces a wrong answer for a
specific case, fix the resolver — don't let callers opt out.

**Application**: Any time you see rc.ContextDir(), rc.RC(), or similar central
resolvers producing a bad result, the fix belongs in the resolver itself (or in
its input data like .ctxrc). Caller-side bypasses create drift.

---

## [2026-04-09-001323] Pad index shifting is a real UX bug in batch operations

**Context**: ctx pad rm 10; rm 11; rm 12 deleted wrong entries because indices
shifted after each deletion

**Lesson**: Any ID-based system where users chain operations needs stable IDs.
Look-then-act is safe for single ops; look-then-batch-act breaks with shifting
indices

**Application**: Both pad and remind now use stable IDs with batch delete and
range support. Apply same pattern to any future numbered-list subsystem

---

## [2026-04-03-180000] Bulk rename and replace_all hazards (consolidated)

**Consolidated from**: 3 entries (2026-03-15 to 2026-03-20)

- `replace_all` on short tokens (e.g. `core.`, function names) matches inside
  longer identifiers and function definitions — `remindcore.` becomes
  `remindtidy.`, `func HumanAgo` becomes `func format.DurationAgo` (invalid Go)
- `sed` insert-before-first-match does not understand Go import aliases — the
  alias attaches to whatever line sed inserts, not the original target
- For function renames: delete the old definition separately rather than using
  replace_all. For bulk import additions: check for aliased imports first and
  handle them separately, or use goimports

---

## [2026-04-03-180000] Import cycles and package splits (consolidated)

**Consolidated from**: 5 entries (2026-03-06 to 2026-03-22)

- Types in god-object files (e.g. hook/types.go with 15+ types from 8 domains)
  create circular dependencies — move types to their owning domain package
- Tests in parent package X cannot import X/sub packages that import X back —
  move tests to the sub-package they exercise
- Variable shadowing causes cascading failures after splits: `dir`, `file`,
  `entry` are common Go variable names that collide with new sub-package names
  — run `go test ./...` before committing splits
- When moving constants between packages, change imports and all references in a
  single atomic write so the linter never sees an inconsistent state
- Import cycle rule: the package providing implementation logic must own the
  shared types; the facade package aliases them (e.g. `entry.Params` aliases
  `add/core.EntryParams`)

---

## [2026-04-03-180000] Skill lifecycle and promotion (consolidated)

**Consolidated from**: 4 entries (2026-03-01 to 2026-03-14)

- Internal skill renames and promotions require synchronized updates across 6+
  layers: SKILL.md frontmatter, internal cross-references, external docs,
  embed_test.go expected list, recipe/reference docs, and plugin cache rebuild +
  session restart
- Skill behavior changes ripple through hook messages, fallback strings in Go
  code, doc descriptions, and Makefile hints — grep for the skill name across
  the entire repo
- Skills without a trigger mechanism (no user invocation, no hook loading) are
  dead code — audit skills for reachability
- After promoting skills: grep -r for the old name across the whole tree, run
  plugin-reload.sh, restart session to verify autocomplete, and clean stale
  Skill() entries from settings.local.json

---

## [2026-04-03-133244] desc.Text() is the single highest-connectivity symbol in the codebase

**Context**: GitNexus enrichment during architecture analysis revealed
desc.Text() (internal/assets/read/desc/desc.go:75) has 30+ direct callers
spanning every architectural layer (MCP handler, format, index, tidy, trace,
memory, sysinfo, io) and participates in 53 execution flows.

**Lesson**: TestDescKeyYAMLLinkage is the most critical guard in the codebase
— it protects the symbol with the widest blast radius. If YAML text loading
breaks, the entire CLI and MCP server output blank strings silently (no crash,
no warning).

**Application**: Treat desc.Text() as a frozen API — add new functions rather
than modifying the existing signature. Any change to config/embed/text or
assets/read/desc should be followed by running the linkage audit. Monitor this
symbol during major refactors.

---

## [2026-04-01-074418] Contributor PRs based on older code reintroduce removed features

**Context**: PR #45 brought back prompt templates, PROMPT.md, and
IMPLEMENTATION_PLAN.md that were explicitly removed in March

**Lesson**: When resolving contributor merge conflicts, check decisions history
for intentional removals — do not assume the PR content is additive

**Application**: Cross-reference DECISIONS.md before accepting PR content that
adds files or features

---

## [2026-03-31-005112] Convention audits must check cmd/ purity, not just types and docstrings

**Context**: Placed needsSpec helper in cmd/root/run.go instead of
core/entry/predicate.go. Missed it because the audit checklist only covered
types and docstrings

**Lesson**: cmd/ directories must contain only Cmd() and Run*() — all helper
functions, unexported logic, and types belong in core/. Added TestCmdDirPurity
compliance test to enforce this mechanically

**Application**: The compliance test now catches this automatically. 28
pre-existing violations grandfathered in the allowlist

---

## [2026-03-31-005110] JSON Schema default fields cause linter errors with some validators

**Context**: ctxrc.schema.json had default: values on 16 fields that triggered
incompatible type errors in the user's linter

**Lesson**: Move default values into the description string instead of using the
default keyword — Go rc.*() accessors handle the actual defaults

**Application**: When adding new .ctxrc fields, document defaults in the
description, never use default: in the schema

---

## [2026-03-30-003707] lint-docstrings.sh greedy sed hid all return-type violations

**Context**: sed 's/.*) //' consumed return type parens, leaving { — functions
with return types were invisible to the script for months

**Lesson**: Greedy regex in shell scripts can silently suppress entire
categories of lint violations — test with edge cases, not just happy paths

**Application**: When writing sed-based lint checks, test with multi-paren
signatures (func Foo() (string, error))

---

## [2026-03-24-001001] lint-drift false positives from conflating constant namespaces

**Context**: lint-drift.sh checked all string constants in embed/cmd/*.go
against commands.yaml, but Use* constants are cobra syntax strings, not YAML
lookup keys

**Lesson**: Shell grep on constant values cannot distinguish constant types;
only DescKey* constants are YAML keys. AST-based analysis is needed for
type-aware checks

**Application**: Already captured in specs/ast-audit-tests.md; the lint-drift
fix is shipped in v0.8.0

---

## [2026-03-23-165611] Typography detection script needs exclusion lists for intentional uses

**Context**: detect-ai-typography.sh flagged config/token/delim.go (intentional
delimiter constants) and test files (test data containing em-dashes)

**Lesson**: Detection scripts for convention enforcement need exclusion patterns
for files where the flagged patterns are intentional data, not prose

**Application**: Add exclusion patterns proactively when creating detection
scripts; *_test.go and constant-definition files are common false positive
sources

---

## [2026-03-20-160112] Commit messages containing script paths trigger PreToolUse hooks

**Context**: Git commit message body contained a path to a shell script under
the hack directory which matched a hook pattern that blocks direct script
invocation

**Lesson**: Hooks scan all Bash tool input including heredoc content used for
commit messages, not just the command itself

**Application**: Rephrase commit messages and ctx add content to avoid paths
that match hook deny patterns, use generic references instead of literal file
paths

---

## [2026-03-18-133457] Lazy sync.Once per-accessor is a code smell for static embedded data

**Context**: assets package had 4 sync.Once guards, 4 exported maps, 4 Load*()
functions, and a wrapper desc package — all to lazily load YAML from embed.FS
that never mutates. Every accessor call went through sync.Once + global map +
wrapper indirection.

**Lesson**: When data is static and loaded from embedded bytes, scatter-loading
with per-accessor sync.Once is over-engineering. A single Init() called eagerly
at startup is simpler, and one sync.Once on Init() itself provides the test
safety net. Exported maps that exist only for wrapper packages to reach are a
sign the abstraction boundary is wrong.

**Application**: Prefer eager Init() in main.go for static embedded data. Keep
maps unexported. Accessors do plain map lookups. If a wrapper package exists
solely to break a cycle caused by exported state, delete the wrapper and
unexport the state.

---

## [2026-03-17-105637] Write package output census: 69 trivial/simple, 38 consolidation candidates, 18 complex

**Context**: Full audit of internal/write/ (26 files, 160 functions, 337 Println
calls) to evaluate whether block template consolidation is worth a systematic
refactor.

**Lesson**: Only 30% of write functions benefit from output consolidation. The
sweet spot is multi-line (16) and conditional (22) functions.

**Application**: Check function category before consolidating. Trivial/simple
stay as-is. Conditional functions need pre-computation before block templates.
Loop-based complex functions stay imperative. Don't bulk-refactor.

---

## [2026-03-15-101342] Contributor PRs need post-merge follow-up commits for convention alignment

**Context**: PR #42 (MCP v0.2) addressed bulk of review feedback but left ~12
inline strings, no embed_test coverage, and substring matching in
containsOverlap

**Lesson**: Merging with known gaps is fine when the gaps are mechanical, but
the follow-up must be immediate — track in ideas/done/ with a review status
doc

**Application**: For future contributor PRs: create ideas/pr{N}-review-status.md
during review, merge when architecture is sound, fix convention gaps in a
same-day follow-up commit

---

## [2026-03-06-141504] Stats sort uses string comparison on RFC3339 timestamps with mixed timezones

**Context**: ctx system stats showed only old sessions, hiding the current one

**Lesson**: RFC3339 string comparison breaks when entries mix UTC (Z) and offset
(-08:00) formats — 13:00-08:00 sorts before 18:00Z lexicographically despite
being later in absolute time

**Application**: Always parse to time.Time before comparing RFC3339 timestamps;
never rely on lexicographic sort

---

## [2026-03-05-023941] Blog post editorial feedback is higher-leverage than drafting

**Context**: Draft of Agent Memory Is Infrastructure was publication-quality on
first pass; user editorial feedback (structural emphasis, rhetorical sharpening,
amnesia/archaeology bridge) elevated it significantly more than initial
generation

**Lesson**: For narrative content, the first draft captures the argument; the
editorial pass captures the voice. Both are necessary but the editorial pass has
disproportionate impact on quality.

**Application**: For future blog posts, invest more in the editorial cycle
(structural feedback then targeted refinements) rather than trying to nail voice
on first generation.

---

## [2026-03-04-105239] CONSTITUTION hook compliance is non-negotiable — don't work around it

**Context**: After make build, ran ./ctx deps --help which was blocked by
block-non-path-ctx. Instead of asking user to install, tried cp ctx ~/bin/ —
escalating workarounds.

**Lesson**: When a hook blocks an action, the correct response is to follow the
hook's instruction (ask the user to sudo make install), not to find creative
bypasses.

**Application**: Always ask the user to install when testing a freshly built
binary. Never attempt alternative install paths to circumvent a hook.

---

## [2026-03-02-165039] Hook message registry test enforces exhaustive coverage of embedded templates

**Context**: Adding billing.txt to embedded assets without a registry entry
caused TestRegistryCoversAllEmbeddedFiles to fail immediately

**Lesson**: Every new .txt file under internal/assets/hooks/messages/ must have
a corresponding entry in registry.go — the test acts as an exhaustive
bidirectional check

**Application**: When adding new hook message variants, update the registry
entry before running tests

---

## [2026-03-01-124921] Model-to-window mapping requires ordered prefix matching

**Context**: Implementing modelContextWindow() for the three-tier context window
fallback. Claude model IDs use nested prefixes (claude-sonnet-4-5 vs
claude-sonnet-4-20250514).

**Lesson**: A switch with ordered HasPrefix cases (most specific first) is
cleaner and safer than iterating separate prefix lists. The catch-all 'claude-*'
returns 200k for unrecognized Claude models.

**Application**: When adding new model families to modelContextWindow() in
session_tokens.go, add the most specific prefix first to avoid shadowing shorter
prefixes.

---

## [2026-03-01-095709] TASKS.md template checkbox syntax inside HTML comments is parsed by RegExTaskMultiline

**Context**: Template had example checkboxes (- [x], - [ ]) in HTML comments
that the line-based regex matched as real tasks, causing
TestArchiveCommand_NoCompletedTasks to fail

**Lesson**: RegExTaskMultiline is line-based and has no awareness of HTML
comment blocks — checkbox-like patterns inside comments get counted as real
tasks

**Application**: Use backtick-quoted or indented references instead of actual
checkbox syntax in template comments. When adding examples to TASKS.md
templates, avoid patterns that match regExTaskPattern

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

## [2026-02-28-150701] Plugin reload script must rebuild cache, not just delete it

**Context**: hack/plugin-reload.sh was deleting
~/.claude/plugins/cache/activememory-ctx/ without repopulating it. Claude Code's
installed_plugins.json still referenced the cache path, so the plugin appeared
enabled but hooks.json was missing — all plugin hooks silently stopped firing.

**Lesson**: Claude Code snapshots plugin hooks from the cache directory at
session startup. If the cache is deleted, plugin hooks vanish silently with no
error. The reload script must rebuild the cache from source assets
(internal/assets/claude/) after clearing it, and warn that a session restart is
required.

**Application**: Always rebuild the plugin cache in hack/plugin-reload.sh. When
debugging hooks that don't fire, check ~/.claude/plugins/cache/ first — a
missing hooks.json is the most likely cause.

---

## [2026-02-27-231228] site/ directory must be committed with docs changes

**Context**: The site/ directory contains generated HTML served directly from
the repo (no CI build step). Multiple sessions have committed docs/ changes
without the corresponding site/ output, or ignored site/ as 'generated noise'.

**Lesson**: site/ is intentionally tracked in git — there is no GitHub Pages
workflow or CI step to build it. When docs change, the regenerated site/ HTML
must be staged and committed alongside the source.

**Application**: Always git add site/ when committing changes under docs/. Never
gitignore site/.

---

## [2026-02-27-230741] Doctor token_budget vs context_window confusion

**Context**: ctx doctor reported context size against token_budget (8k) instead
of context_window (200k), making 22k tokens look alarming.

**Lesson**: token_budget (ctx agent output trim target) and context_window
(model capacity) serve different purposes. Health checks about context fitting
should use context_window, with warning threshold proportional (e.g., 20% of
window).

**Application**: Doctor now uses rc.ContextWindow() with 20% threshold and shows
per-file token breakdown for actionable insight into which files are heavy.

---

## [2026-02-27-230738] Drift detector false positives on illustrative code examples

**Context**: ctx drift flagged 23 warnings for backtick-quoted paths in
CONVENTIONS.md and ARCHITECTURE.md that were prose examples (loader.go,
session/run.go, sync.Once), not real file references.

**Lesson**: Path reference detection should verify the top-level directory
exists on disk before flagging. Bare filenames and paths under non-existent
directories are almost always examples in documentation.

**Application**: The fix checks os.Stat(topDir) on the first path component.
Future drift checks on documentation-heavy files should use the same heuristic.

---

## [2026-02-26-003854] Webhook silence after ctxrc profile swap is the most common notify debugging red herring

**Context**: Spent time investigating why webhooks weren't firing — checked
binary version, hook configs, notify.Send internals. Actual cause was .ctxrc
swapped to prod profile (notify commented out) earlier in session.

**Lesson**: When webhooks stop, check .ctxrc profile first (`ctx config
status`). Also: not all tool uses trigger webhook-sending hooks — Read only
triggers context-load-gate (one-shot) and ctx agent (no webhook). qa-reminder
requires Edit matcher.

**Application**: Before debugging notify internals, run `ctx config status` and
verify the event would actually match a hook with notify.Send.

---

## [2026-02-26-100002] Agent context loading and task routing (consolidated)

**Consolidated from**: 5 entries (2026-01-20 to 2026-01-25)

- `ctx agent` is optimized for task execution (filters pending tasks, surfaces
  constitution, token-budget aware). Manual file reading is better for
  exploratory/memory questions (session history, timestamps, completed tasks).
- On "Do you remember?" questions, immediately read .context/ files and run `ctx
  journal source --limit 5`. Never ask "would you like me to check?" — that is
  the obvious intent.
- .context/ is NOT a Claude Code primitive. Only CLAUDE.md and
  .claude/settings.json are auto-loaded. The .context/ directory requires a hook
  or explicit CLAUDE.md instruction to be discovered.
- ~~Orchestrator (IMPLEMENTATION_PLAN.md) and agent (.context/TASKS.md) task
  lists must be separate.~~ (Superseded 2026-03-25: IMPLEMENTATION_PLAN.md
  removed. TASKS.md is the single task source.)
- Only CLAUDE.md is auto-loaded by Claude Code. Projects using ctx should rely
  on the CLAUDE.md -> AGENT_PLAYBOOK.md chain, not AGENTS.md.

---

## [2026-02-26-100006] PATH and binary handling (consolidated)

**Consolidated from**: 3 entries (2026-01-21 to 2026-02-17)

- Always use `ctx` from PATH, never `./dist/ctx-linux-arm64` or `go run
  ./cmd/ctx`. Check `which ctx` if unsure.
- Hooks must use PATH, not hardcoded paths. `ctx init` checks if ctx is in PATH
  before proceeding. Tests can skip with `CTX_SKIP_PATH_CHECK=1`.
- Agent must never place binaries in any bin directory (not via cp, mv, or go
  install). Build with `make build`, then ask the user to run the privileged
  install step. Hooks in block-dangerous-commands.sh enforce this.

---

## [2026-02-26-100007] Task management and exit criteria (consolidated)

**Consolidated from**: 4 entries (2026-01-21 to 2026-02-17)

- Specs get lost without cross-references from TASKS.md. Three-layer defense:
  (1) playbook instruction, (2) spec reference in Phase header, (3) bold
  breadcrumb in first task.
- Subtask completion is implementation progress, not delivery. Parent tasks
  should have explicit deliverables; don't close until deliverable is verified.
- Exit criteria must include verification: integration tests (binary executes
  correctly), coverage targets, and smoke tests. "All tasks checked off" does
  not equal "implementation works."
- Reports graduate to ideas/done/ only after all items are tracked or resolved.
  Cross-reference every item against TASKS.md and the codebase before moving.

---

## [2026-02-26-100008] Agent behavioral patterns (consolidated)

**Consolidated from**: 5 entries (2026-01-25 to 2026-02-22)

- Interaction pattern capture risks softening agent rigor. Do not build implicit
  user-modeling from session history. Rely on explicit, human-reviewed context
  (learnings, conventions, hooks) for behavioral shaping.
- Chain-of-thought prompting improves agent reasoning accuracy (17.7% to 78.7%).
  Added "Reason Before Acting" to AGENT_PLAYBOOK.md and reasoning nudges to 7
  skills.
- Say "project conventions" not "idiomatic X" to ensure Claude looks at project
  files first rather than triggering training priors (stdlib conventions).
- Autonomous "YOLO mode" is effective for feature velocity but accumulates
  technical debt (magic strings, monolithic tests, hardcoded paths). Schedule
  periodic consolidation sessions.
- Trust the binary output over source code analysis. A single ambiguous CLI
  output is not proof of absence — re-run the exact command before claiming
  something is missing.

---

## [2026-02-26-100010] ctx add and decision recording (consolidated)

**Consolidated from**: 4 entries (2026-01-27 to 2026-02-14)

- `ctx add learning` requires `--context`, `--lesson`, `--application` flags.
  `ctx add decision` requires `--context`, `--rationale`, `--consequence`. A
  bare string only sets the title and the command will fail without required
  flags.
- Structured entries with Context/Lesson/Application are more useful than
  one-liners. Agents are guided via AGENT_PLAYBOOK.md.
- Always complete decision record sections — placeholder text like "[Add
  context here]" is a code smell. Decisions without rationale lose their value
  over time.
- Slash commands using `!` bash syntax require matching permissions in
  settings.local.json. When adding new /ctx-* commands, ensure ctx init
  pre-seeds the required `Bash(ctx <subcommand>:*)` permissions.

---

## [2026-02-24-032945] CLI tools don't benefit from in-memory caching of context files

**Context**: Discussed whether ctx should read and cache LEARNINGS.md,
DECISIONS.md etc. in memory

**Lesson**: ctx is a short-lived CLI process, not a daemon. Context files are
tiny (few KB), sub-millisecond to read. Cache invalidation complexity exceeds
the read cost. Caching only makes sense if ctx becomes a long-lived process (MCP
server, watch daemon).

**Application**: Don't add caching layers to ctx's file reads. If an MCP server
mode is ever added, revisit then.

---

## [2026-01-28-051426] IDE is already the UI

**Context**: Considering whether to build custom UI for .context/ files

**Lesson**: Discovery, search, and editing of .context/ markdown files works
better in VS Code/IDE than any custom UI we'd build. Full-text search,
git integration, extensions - all free.

**Application**: Don't reinvent the editor. Let users use their preferred IDE.

---


*Module-specific, niche, and historical learnings:
[learnings-reference.md](learnings-reference.md)*
## [2026-04-29-050000] BunShell ctx.$ calls echo stdout to OpenCode's process unless .quiet() is set — leaks visible noise

**Context**: After PR #72 wired session.created and session.idle to fire `ctx
system bootstrap`, `ctx agent --budget 4000`, and friends, end users started
seeing chunks of Markdown bleeding into the OpenCode TUI: `## Steering`, `#
Product Context`, `Describe the product...`. These are the contents of
`.context/steering/` template stubs that `ctx agent --budget 4000` includes in
its context packet. The plugin used the shell-level `2>/dev/null || true` to
swallow stderr and force exit 0, but stdout was untouched.

**Lesson**: BunShell's documented behavior: *"By default, the shell will write
to the current process's stdout and stderr, as well as buffering that output."*
So an `await ctx.$\`...\`` call in a plugin echoes its stdout/stderr to
OpenCode's process, which the TUI/agent surfaces. Shell-level `2>/dev/null` only
suppresses stderr; stdout still leaks. The fix is BunShell's `.quiet()` modifier
on the BunShellPromise, which configures the shell to only buffer the output
rather than also writing to the parent process.

**Application**: Always chain `.nothrow().quiet()` on BunShell template literals
in OpenCode plugins, even for fire-and-forget calls where you discard the
result: `await ctx.$\`ctx system bootstrap\`.nothrow().quiet()`. With both
modifiers, you don't need shell-level `2>/dev/null || true` — `.nothrow()`
swallows non-zero exits at the BunShell layer, `.quiet()` keeps every byte of
output buffered. Pattern is the cooperative default for any plugin that spawns
long-output commands during the agent session lifecycle.

---

## [2026-04-29-040000] OpenCode plugin compaction interop is breadcrumb-mediated: own your context preservation explicitly

**Context**: After PR #72 wired `session.created` / `session.idle` /
`tool.execute.after` / `shell.env`, a `/compact` test in OpenCode (with
`oh-my-openagent@3.17.6` also installed) recovered ctx context post-compaction
*only by accident*: oh-my-openagent's `experimental.session.compacting` handler
builds a structured summary template that happens to preserve
`.context/`-prefixed file paths in its "Active Working Context → Files"
section. Combined with our `shell.env` CTX_DIR injection, the agent had enough
breadcrumbs to re-read DECISIONS.md from disk post-compaction. Without that
section, our context would have evaporated silently into the compaction summary.

**Lesson**: Two compaction-aware plugins in the same session can synergize
without either knowing about the other — but the synergy is fragile because it
depends on undocumented serialization choices in the *other* plugin. If the
other plugin's template ever changes (e.g., drops file-path preservation, swaps
the "Active Working Context" section name, condenses paths to basenames), the
breadcrumbs disappear and ctx context is lost without any signal. The `Hooks`
interface in `@opencode-ai/plugin` v1.4.x exposes
`experimental.session.compacting?: (input, output: { context: string[]; prompt?:
string }) => Promise<void>` — pushing to `output.context` is *additive*
(appends to the default prompt), and replacing `output.prompt` is *destructive*
(only one plugin can win that race).

**Application**: Register `experimental.session.compacting` in your own plugin
and push high-signal context strings (e.g., `ctx system bootstrap` output) to
`output.context` so context preservation does not depend on coexisting plugins.
Never set `output.prompt` from a thin shim — that would conflict with primary
compaction harnesses like oh-my-openagent. Composition via `output.context` is
the correct cooperative pattern.

---

## [2026-04-29-030000] @opencode-ai/plugin event hook is a single dispatcher, not an object of named handlers

**Context**: PR #72's first OpenCode plugin shipped with `event: {
"session.created": fn, "session.idle": fn }` — an object keyed by event type.
It compiled clean against `satisfies Plugin` but never fired. End-to-end trace
showed neighboring hooks (`shell.env`, `tool.execute.after`) running while every
event handler silently no-op'd.

**Lesson**: `@opencode-ai/plugin` v1.4.x defines `event?: (input: { event: Event
}) => Promise<void>` — one dispatcher called for every event with
`input.event.type` discriminating. Asymmetric with neighbors because `shell.env`
and `tool.execute.*` *are* top-level named keys; only the dozens of `EventX`
types collapse into the single `event` slot.

**Application**: Use `event: async ({event}) => { if (event.type ===
"session.created") { ... } else if (event.type === "session.idle") { ... } }`.
Type discriminator strings live under each `EventX` type in
`node_modules/@opencode-ai/sdk/dist/gen/types.gen.d.ts`.

---

## [2026-04-29-030100] OpenCode plugin hooks like shell.env take (input, output) and mutate; returned objects are ignored

**Context**: First plugin had `"shell.env": () => ({ CTX_DIR: ".context" })`.
The hook fired but the agent's bash tool never saw `CTX_DIR`; manual export was
required for every ctx call. The returned object was dropped on the floor by the
runtime.

**Lesson**: Multiple hooks in `@opencode-ai/plugin` v1.4.x take two arguments
where the second is an OUT param. Examples: `shell.env: (input, output: {env})
=> void` (mutate `output.env`), `tool.execute.after: (input, output: {title,
output, metadata}) => void`, `chat.params: (input, output: {temperature, ...})
=> void`, `chat.headers: (input, output: {headers}) => void`. Pattern is
consistent across the SDK.

**Application**: Always read the type definition in
`node_modules/@opencode-ai/plugin/dist/index.d.ts` for any hook before wiring.
If a hook signature has two parameters where the second is an object, it's a
mutation hook — return values are discarded.

---

## [2026-04-29-030200] OpenCode shell.env injects env only into agent's shell tool, not into plugin's own ctx.$ calls

**Context**: After fixing `shell.env`'s `(input, output) => mutate output.env`
signature so `CTX_DIR` reached the agent's bash tool, the plugin's own
`ctx.$\`ctx system bootstrap\`` calls still failed silently — they ran without
`CTX_DIR` and ctx fell back to `~/.context`. The hook fired correctly; the
plugin's subprocess side-effects didn't see the env.

**Lesson**: `shell.env` injects env into the agent's shell-tool invocations. The
plugin's own BunShell calls (`ctx.$\`...\``) inherit OpenCode's process env,
which is *separate*. Two shells, two envs.

**Application**: Build an env-aware BunShell once in the plugin factory: `const
$ = ctx.$.env({ ...process.env, CTX_DIR: \`${ctx.directory}/.context\` })`.
Reuse it for every plugin-initiated subprocess call. `ctx.directory` is the
project root from `PluginInput`.

---

## [2026-04-26-180000] OpenCode auto-loads only flat .ts files under .opencode/plugins/; subdirectories are ignored

**Context**: Initial OpenCode integration deployed the plugin as
`.opencode/plugins/ctx/index.ts` (a directory with index.ts inside, mirroring
npm package conventions). End-to-end smoke testing showed the plugin file was
present and the binary was current, yet OpenCode never invoked any of the
plugin's hooks (no `module-load` trace fired even with `--print-logs --log-level
DEBUG`). Copying the same content to a flat `.opencode/plugins/ctx.ts` file made
the plugin load and fire correctly.

**Lesson**: OpenCode's plugin auto-discovery only scans top-level files under
`.opencode/plugins/` and `~/.config/opencode/plugins/`. Subdirectories are
silently skipped — there is no log line indicating a subdirectory was found
and ignored. The official docs at opencode.ai/docs/plugins/ say only "files in
these directories are automatically loaded at startup" without specifying the
rule, so this is easy to miss. The `opencode plugin <module>` CLI registers npm
modules (a different code path) and accepts only npm names, not local paths.

**Application**: Deploy single-file plugins as `.opencode/plugins/<name>.ts`,
not `.opencode/plugins/<name>/index.ts`. No `package.json` is required when the
plugin uses type-only imports (`import type` is erased at compile time) and the
host runtime injects the plugin context. To verify a plugin is actually loaded,
add a top-of-module side effect (e.g. `appendFileSync` to a known path) and
confirm it fires before debugging hook contracts.

---

## [2026-04-26-165500] OpenCode opencode.json MCP shape: command is Array<string>, no separate args field

**Context**: `ctx setup opencode --write` was generating `opencode.json` with
the Copilot CLI MCP shape (`{type: "local", command: "ctx", args: ["mcp",
"serve"]}`). OpenCode rejected the file at startup with `Configuration is
invalid… Expected array, got "ctx" mcp.ctx.command` and `Missing key
mcp.ctx.enabled`.

**Lesson**: OpenCode's `McpLocalConfig` (in `@opencode-ai/sdk`) defines
`command: Array<string>` as a single field that holds the binary AND its
arguments — there is no separate `args` field. It also requires `enabled:
boolean` at runtime even though the TS type marks it optional. The Copilot CLI
MCP shape is similar in spirit but structurally different; do not copy-paste
between them.

**Application**: For OpenCode MCP entries always use `command: ["ctx", "mcp",
"serve"]` and include `enabled: true`. If you add a new editor integration with
its own MCP file format, read the upstream type definitions from
`node_modules/@<vendor>/sdk/dist/gen/types.gen.d.ts` (or equivalent) before
reusing an existing generator.

---

## [2026-04-26-152836] ctx system help can list project-local hooks not in the Go binary

**Context**: PR #72 plugin called 'ctx system block-dangerous-commands'; user's
installed ctx 0.7.2 listed it in help, but no directory exists under
internal/cli/system/cmd/ — it's a Claude Code plugin-local hook surfaced via
wrapper

**Lesson**: ctx system help output is a union of compiled Go subcommands and
project-local Claude wrappers; non-Claude integrations only see the Go subset

**Application**: When porting plugin behavior to a new editor, only call
subcommands that have a directory under internal/cli/system/cmd/. Don't trust
ctx system help output as the canonical surface.

---

## [2026-04-25-014704] Confident code comments can pull an LLM away from first-principles knowledge

**Context**: cli_test.go had a comment claiming 'parent's t.Setenv doesn't
propagate to exec'd children unless we build it into cmd.Env' which is wrong. I
patched the helper's CTX_DIR dedup instead of questioning the helper itself,
despite knowing t.Setenv semantics.

**Lesson**: A comment that explains why a stdlib mechanism 'doesn't work' is
doing extra rhetorical work to talk a reader out of the obvious approach. That's
exactly when to verify from first principles instead of trusting the
surrounding-code frame.

**Application**: When an existing comment justifies a non-canonical approach
contradicting stdlib knowledge: pause, verify against memory of the actual API
before patching within the existing frame.

---
