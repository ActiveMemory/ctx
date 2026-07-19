# audit-lint-compliance

## [2026-07-18-084406] Adding a new internal/ package must clear a fixed audit gauntlet — satisfy it up front

**Context**: The internal/disclosure + internal/cli/disclosure packages tripped ~6 make audit rounds while building pd-m1/m2.

**Lesson**: All mechanically enforced: doc.go per package (compliance), docstring Parameters/Returns floor on ALL funcs incl. unexported (lint-docstrings), string literals must live in config/* (magic-strings), exported API and unexported helpers in SEPARATE files (mixed-visibility), <=80-char lines, DescKey<->YAML bijection (both directions), no dead exports, and TestDocGoSubcommandDrift mis-parses a group doc.go bullet like '- internal/cli/x/cmd/y:' as subcommand claims (bulletRe backtracks on the slash).

**Application**: For a new package, from the first file: one doc.go with a package comment (never a package comment in a regular file); split exported vs unexported functions across files; every literal as a config constant; full Parameters/Returns docstrings; and in a group doc.go use SPACED bullets ('- ctx x sub: ...') not slash-paths, or leave documented empty like kb/doc.go.

---

## [2026-07-15-141726] gosec G101 has two independent triggers: identifier-keyword match and value-entropy

**Context**: Removed the blanket G101 path-exclusion for internal/config/embed/text/; gosec then flagged 3 lines. Two were identifier matches (DescKeyErrHubGenerateToken, ...AdminTokenRequired — 'token' substring in the name), but the third flagged only sourcemap.write-row out of 8 identical *WriteRow consts.

**Lesson**: G101 matches on (a) the variable/const IDENTIFIER containing a credential keyword (token/pass/pwd/secret/cred/apiKey — case-insensitive substring), AND (b) a Shannon-entropy heuristic on the string VALUE. The single sourcemap hit was entropy, not name — which is why 7 sibling *WriteRow consts weren't flagged.

**Application**: Fix each trigger at the source, never with #nosec or a path exclusion: for a name-match, rename the identifier off the keyword (keep the value/i18n-key/env-var intact so linkage/behavior are unchanged); for an entropy hit, change the value (e.g. write-row->append-row). Confirm empirically: renaming the identifier while keeping the value clears a name-match, proving it wasn't value-based.

---

## [2026-07-04-160749] Compliance tests are the style guide for new commands

**Context**: Implementing ctx system statusline tripped eight architecture gates in sequence: magic strings/values must live in internal/config/, one visibility per file, cmd/ files restricted to cmd.go+run.go+doc.go with helpers in core/, types in types.go, the ctxRC schema mirror struct, the system doc.go subcommand registry, godoc Parameters/Returns structure, and 80-char lines

**Lesson**: House conventions are enforced by go test (internal/compliance/ and sibling suites), not documented prose; the test suite is the authoritative style guide

**Application**: Before scaffolding a new command, read internal/compliance/ and copy an existing cmd package's structure; run the full suite early rather than discovering gates by failure at commit time

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

