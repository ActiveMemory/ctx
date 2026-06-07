# Archived Learnings (consolidated 2026-06-07)

Originals replaced by consolidated entries in LEARNINGS.md.

## Group: ctx-dream design principles (consolidated)

## [2026-06-07-090909] A statistical merit score ranks attention; evidence decides eligibility

**Context**: Lifting the Hermes 'Dreaming' weighted scoring rubric (relevance/frequency/diversity/recency/consolidation/richness) into ctx-dream's open 'merit signal' question. Hermes uses the score as a promotion_threshold (score >= 0.6 -> auto-promote).

**Lesson**: A keyword/frequency/recency score measures ATTENTION (what to surface first), not TRUTH (whether the claim still holds). A frequent-but-stale topic scores high yet may be already-implemented or obsolete. Using such a score as an autonomous promote threshold is the trap; grounding-against-code/specs is what establishes eligibility. The score decides ordering; evidence decides eligibility.

**Application**: Adopt the merit/scoring rubric in ctx-dream only as a ranking signal feeding ruthless self-rejection (surface top-N to the human), never as a promotion threshold. Pair any statistical ranking with an evidence/grounding gate that decides what is eligible to surface at all.

---

## [2026-06-07-090904] Dream consolidation only proposes; it never autonomously writes canonical memory

**Context**: Cross-checking the Hermes 'Dreaming' sibling against specs/ctx-dream.md: Hermes' Deep-Sleep phase auto-promotes scored entries straight into MEMORY.md with no human in the loop.

**Lesson**: The load-bearing invariant of ctx-dream (Option B) is that consolidation emits PROPOSALS only; a human gate sits between the dream pass and any write to the five canonical files / MEMORY.md. The naive 'dream -> artifacts' direct arrow is exactly the failure mode of 'Useful Memories Become Faulty When Continuously Updated by LLMs' (2605.12978): naive autonomous consolidation can push memory utility BELOW the no-memory baseline. Independent designs (Hermes, OpenClaw) keep re-deriving the same sleep-phase architecture but omit the gate — corroboration of the shape, cautionary on the autonomy. Complements the existing 'a single agreeable LLM is not an adversarial gate' learning: that one says WHO the gate is; this one says a gate must EXIST at all.

**Application**: Any consolidation/dream/memory-rewrite feature in ctx must route through a human accept/reject step before touching canonical artifacts. When evaluating an external memory-consolidation design, the first check is: does it autonomously write canonical, or only propose? Autonomous-write is a reject.

---

## [2026-06-06-162156] A single agreeable LLM is not an adversarial gate

**Context**: Designing ctx-dream's promotion gate; reviewed the 28-source corpus at ideas/ctx-dreams/research/.

**Lesson**: Asked to critique a proposal, a reasoning model silently repairs the missing justification and approves it (ReportLogic finding). Robust gating needs multi-critic consensus + swap-consistency, or a human. Backing: Auto-Dreamer (2605.20616) is nearly this architecture; 'Useful Memories Become Faulty When Continuously Updated by LLMs' (2605.12978) is the threat model that naive continuous consolidation rots.

**Application**: For any 'agent reviews/approves agent output' design in ctx, never rely on a single LLM as the gate; use a human or independent multi-critic consensus.

---

## [2026-06-06-162156] Same proposals, two consumers, two interfaces

**Context**: A terse accept/reject CLI felt wrong for the ctx-dream serendipity review.

**Lesson**: A terse, action-coded accept/reject worklist is an agent's review interface; human serendipity needs substance-rich, semantically-generated summaries (no file-hunting). Same underlying data, different presentation per consumer.

**Application**: When an agent and a human consume the same proposals, render two views: dispositional/terse for the agent, substance-forward for the human.

---

## [2026-06-06-162156] Split agent/human work by comparative advantage: taste is the human's axis

**Context**: Deciding who does what in ctx-dream (agent vs human).

**Lesson**: The agent is the reliable gardener for mechanical/verifiable hygiene (never bored, never skips the 47th file); the human owns taste/serendipity, the axis humans still beat agents on. That is WHY the human is the gate, not merely a safety nicety.

**Application**: For curation/review features, give the agent the verifiable mechanical work and reserve the human for judgment/taste; design the human's surface for pleasure (substance to wander), not a queue to drain.

---

## [2026-06-06-162156] Don't-leak is a third safety axis: privacy class propagates from source to derived artifact

**Context**: ctx-dream derives summaries/proposals from ideas/, which is gitignored ("best kept hidden").

**Lesson**: A summary, backup, or ledger-line of a hidden file is itself hidden — derivation inherits the source's privacy class. This is a distinct safety axis alongside don't-corrupt and don't-obey-injected-instructions.

**Application**: Any agent process reading a gitignored source must keep every byproduct in gitignored locations; enforce structurally with `git check-ignore` on each write target (refuse tracked paths), never via a prompt. A deliberate human `promote` is the only sanctioned boundary crossing.

---

## Group: internal/audit & compliance gates for new code (consolidated)

## [2026-05-30-114436] New exported types must live in types.go or TestTypeFileConvention fails

**Context**: Defined Payload and Provenance structs alongside the Load/OverlayFlags funcs in a new payload.go; make test failed in internal/audit on TestTypeFileConvention with '2 NEW type definitions outside types.go'.

**Lesson**: The audit permits type definitions outside types.go only when the file is a 'pure type impl file' (only type defs + their methods, no standalone funcs) or the package is on the exempt list. A file that mixes struct definitions with standalone functions is a violation.

**Application**: When adding a new package that has both types and functions, put the type definitions in a dedicated types.go from the start; methods (with receivers) may live beside the behavior. Run 'go test ./internal/audit/ -run TestTypeFileConvention' to check.

---

## [2026-05-30-212102] tpl package is magic-string-audit-exempt but its call sites are not

**Context**: Migrating tpl_*.go format-string consts to text/template handles; a Render("name",...) sketch and map[string]any{"Key":...} render data would both trip audit/magic_strings_test.go (TestNoMagicStrings).

**Lesson**: internal/assets/tpl is in the magic-strings audit exemptStringPackages, so template-path literals are sanctioned there; but render data passed from non-exempt caller packages must be a typed struct (e.g. tpl.ObsidianData{...}), never a map[string]any with literal keys, which trips the audit at the call site.

**Application**: When adding a template, define a typed data struct in tpl/types.go and pass it at the call site; never pass map literals from caller packages.

---

## [2026-05-24-092924] Audit gates that bite when introducing new packages and helpers

**Context**: While landing the pad-undo Phase 1 work, the project audit suite (internal/audit) caught two violations on the new history.go file that aren't surfaced by golangci-lint or build errors: TestNoMixedVisibility and TestNoMagicStrings.

**Lesson**: TestNoMixedVisibility flags ANY unexported func in a file that also contains exported funcs — even with full Parameters/Returns doc sections. The fix is to split unexported helpers into a sibling file like <name>_internal.go in the same package. TestNoMagicStrings flags warn-format string literals passed to logWarn.Warn — they must live as named constants in internal/config/warn/, not inline. TestDocCommentStructure additionally requires Parameters: and Returns: sections on every helper regardless of visibility. The fuller catalog (from landing the audit-channel feature, a whole new CLI command + hook): TestNoMagicValues flags bare integers like `24` (use a named const, e.g. HoursPerDay). TestNoCmdPrintOutsideWrite forbids cmd.Println outside internal/write/ — route all output through a write/<area> function. TestNoNakedErrors forbids errors.New outside internal/err/ — even sentinel `var Err... = errors.New(...)` must live in the err package and be re-exported if a core package needs `errors.Is` against it. TestTypeFileConvention wants struct type definitions in a types.go file, not scattered in logic files. TestCmdDirPurity forbids unexported helper funcs in cmd/ dirs — they belong in a core/ package (so a hook's render helpers go to internal/cli/system/core/<area>/, not the cmd/<hook>/ dir). TestNoLiteralMdExtension forbids literal ".md" — use file.ExtMarkdown. TestDocGoSubcommandDrift requires the PARENT package's doc.go to list every new subcommand (both the cli-area doc.go and, for hooks, internal/cli/system/doc.go). TestDescKeyYAMLLinkage requires every DescKey constant to have a matching yaml entry. TestNoLiteralWhitespace forbids "\r\n"/"\n" literals — use token.NewlineCRLF / token.NewlineLF. And the hook-message registry has a hardcoded count test (TestRegistryCount) that must be bumped when you add a registry.yaml entry. staticcheck QF1012 also fights the audit here: it wants fmt.Fprintf(&b, ...) but TestNoUncheckedFmtWrite forbids discarding Fprintf's return — resolve by building the string with fmt.Sprintf first, then b.WriteString(s).

**Application**: When creating a new core/store-shaped file with both exported API and unexported helpers, split immediately into <name>.go (exported) + <name>_internal.go (unexported) — don't wait for the audit failure. When using logWarn.Warn for a new warning class, add the format constant to internal/config/warn/warn.go FIRST, then reference cfgWarn.<Name> at the call site. All new helpers (exported or not) get full godoc Parameters/Returns blocks. For a whole new CLI command, budget for the full gate set up front: types.go for structs, internal/err/<area>/ for ALL errors (including sentinels), internal/write/<area>/ for ALL output, a core/ package for any non-trivial helpers used by a cmd/ or hook dir, every format string and magic number as a named constant, every DescKey paired with a yaml entry, and the parent doc.go subcommand list updated. Run `go test ./internal/audit/ ./internal/compliance/` early and often — these gates are not surfaced by `go build` or `golangci-lint`.

---

## [2026-05-17-055500] Pre-emptive constants are dead exports; ship constants only when their caller lands

**Context**: During Phase KB Stage 3, I added the full set of expected constants
to `internal/config/kb/kb.go`: closeout-mode names, schema filenames, life-stage
tokens, pass-mode tokens, the LifeStageThreshold integer. Many of these had no
caller yet because their consumers (doctor advisories, the `ctx kb site build`
zensical wiring, doctor advisory checks) were Phase 7 work. The
`dead_exports_test.go` audit flagged 28 of them. Same for
`cli/kb/core/path/SchemasDir` and `KBArtifactFile`, plus `regex.SlugWithSlash`.

**Lesson**: ctx's dead-export audit is symbol-graph-strict: any exported const /
var / func without an internal reader fails the gate. You cannot scaffold
constants ahead of their callers, even if you know the caller is one phase away.
The constants must land in the same commit (or a strict precursor commit) as the
code that reads them.

**Application**: When defining configuration constants for a new feature, write
the caller first or in the same change. If a constant truly needs to ship ahead
of its caller (rare), park it in a TASKS.md line, not a config file. The audit
treats "future use" as dead.

---

## [2026-05-17-060000] naked_errors audit rejects fmt.Errorf wrapping outside internal/err/<area>/

**Context**: When fixing Phase KB audit failures, I initially assumed
`fmt.Errorf("desc: %w", err)` wrapping at the call site satisfies the
naked_errors audit. It does not. `internal/audit/naked_errors_test.go` flags
every `fmt.Errorf` and `errors.New` call outside `internal/err/**`. The ctx
convention requires error constructors to live in domain-scoped
`internal/err/<area>/` packages and pull their format strings from either
`internal/config/<area>/` Go-side constants OR `desc.Text(text.DescKey...)` YAML
keys.

**Lesson**: For Phase KB this meant building 14 new err packages (`closeout`,
`handover`, `gitmeta`, `kbevidence`, `kbsourcecoverage`, plus 7 kb-table
packages, `kbcli`, `initkb`) plus matching `internal/config/<area>/` packages
with `ErrMsg<Name>` and `Format<Name>` constants. The pattern: `var ErrX =
errors.New(cfgArea.ErrMsgX)` for sentinels; `func X(args, cause) error { return
fmt.Errorf(cfgArea.FormatX, args, cause) }` for wrapping constructors. Callers
do `errors.Is(err, errArea.ErrX)` for sentinel matching.

**Application**: Estimating the cost of "add a new feature" in ctx must include
the err-package + config-package wiring. Each new error surface is ~3 files per
area (config/<area>/messages.go, err/<area>/<area>.go, the calling code). The
Phase RG `MissingGitError` typed struct was the wrong shape for ctx; it became
`errGitmeta.ErrMissingGitTree` (sentinel) +
`errGitmeta.MissingGitTreeForCmd(cmdName, projectRoot)` (wrapping constructor).

---

## [2026-04-03-180000] Dead code detection (consolidated)

**Consolidated from**: 3 entries (2026-03-15 to 2026-03-30)

- Dead packages can build and test green while being completely unreachable —
  detection requires checking bootstrap registration, not just build success
  (e.g. internal/cli/recall/ existed with tests but was never wired into the
  command tree)
- Files created by `ctx init` that no agent, hook, or skill ever reads are dead
  on arrival — verify there is at least one consumer before adding to init
  scaffolding
- When touching legacy compat code, first ask whether the legacy path has real
  users — if not, delete it entirely rather than improving it (MigrateKeyFile
  had 5 callers and test coverage but zero users)

---

## Group: Error handling: sentinels, unwrapping, and silent discards (consolidated)

## [2026-06-02-051330] os.IsNotExist doesn't unwrap — detect file absence with os.Stat + errors.Is

**Context**: Hardening notify (P0.8.5), `LoadWebhook` needed to tell "encrypted file genuinely absent" (silent: not configured) from "present but broken" (surface it). `os.IsNotExist(loadErr)` on `crypto.LoadKey`'s error is always false: `LoadKey` wraps the os error via `errCrypto.ReadKey` → `fmt.Errorf(desc.Text(...), cause)`, and `os.IsNotExist` does not unwrap. The subtle part is `errors.Is(loadErr, os.ErrNotExist)`: it is **registry-dependent**. `errCrypto.ReadKey`'s format string comes from the externalized text registry (`'read key: %w'`); `fmt.Errorf` honors `%w` at runtime regardless of where the string came from, so in production (registry loaded) `errors.Is` correctly unwraps to `fs.ErrNotExist`. But in a unit-test binary that never initializes the text registry (verified: a probe in `internal/notify`), `desc.Text` returns a string with **no** `%w`, the cause is never wrapped, the error prints `%!(EXTRA *fs.PathError=...)`, and `errors.Is` also returns false. So the same call can behave differently in prod vs. a bare test binary.

**Lesson**: `os.IsNotExist` is the legacy, non-unwrapping check — false on any `fmt.Errorf("…%w…", …)` error; always prefer `errors.Is`. But `errors.Is(err, os.ErrNotExist)` only holds if the wrap actually carries `%w` at runtime, and a wrap whose format string is fetched from a text/i18n registry only carries `%w` when that registry is initialized. `go vet`'s wrap check sees only literal format strings, so a registry-templated wrap is vet-invisible and its wrapping is environment-dependent.

**Application**: To detect file absence reliably, stat the path directly: `os.Stat` returns an unwrapped `*fs.PathError`, so `errors.Is(statErr, os.ErrNotExist)` is dependable in every context. Branch on the stat (absent → not-configured; present → proceed to read/decrypt and surface any error). Reserve `errors.Is(…, os.ErrNotExist)` on a *returned library* error for chains you have confirmed wrap with `%w` independent of registry state.

---

## [2026-06-01-195111] An error-discard catalogue is an inventory, not a verdict — verify each site by reading before fixing

**Context**: Phase EH audited ~184 silent error-discard sites under internal/. The catalogue was built by grep + pattern/name classification (e.g. 'x, _ := SomethingMarshal' => B-marshal). When fixing, several name-inferred verdicts were wrong.

**Lesson**: Classifying a discard by the callee's name or a regex is a guess, not a fact. The discarded value's actual type and the call's role decide the category. Concrete false positives this pass: MergePublished returns (string, bool) — the discard is a 'markers missing' bool, not an error; LoadState returns a State value (not a pointer), so a 'nil-deref' was impossible; io/security's atomic writer already checked the meaningful close and only discarded error-path cleanup closes. All three would have been 'fixed' (churn or breakage) on name-inference alone.

**Application**: Treat any auto-generated audit/catalogue as a worklist of candidates, not findings. Before editing a flagged site, read the callee signature (is the discarded value even an error?) and the enclosing control flow (is it an already-failed path, a best-effort callback, or a data path?). Only then assign return-error vs logWarn vs annotate. This mirrors the Constitution's Context Integrity Invariants: never act on assumed content.

---

## [2026-05-17-180000] Sentinel errors use typed zero-data structs with lazy `desc.Text()` — never Go string consts

**Context**: In a prior Phase KB session I invented an intermediate
`ErrMsg* = "english string"` constant layer in
`internal/config/<pkg>/<pkg>.go`, then in `internal/err/<pkg>/<pkg>.go`
wrote `var ErrX = errors.New(cfgPkg.ErrMsgX)` — backed by a doc comment
claiming `desc.Text` could not be used because `var` initializers run
before `lookup.Init()` populates the embedded YAML table. The framing
was wrong, and the shape contradicted the convention already established
in the codebase. The pre-existing pattern lives in
`internal/err/context/context.go` (commit `e524dd98`): typed error
structs whose `Error()` method calls `assets.TextDesc(...)` /
`desc.Text(...)` lazily, at call time — not at package init.

**Lesson**: The canonical sentinel shape in this repo is a typed,
zero-data struct (for unparameterised sentinels) or a typed struct with
fields (for parameterised errors). The `Error()` method resolves text
via `desc.Text(text.DescKey...)` so the user-facing string lives in
`internal/assets/commands/text/errors.yaml`, keyed by a `DescKey<...>`
constant in `internal/config/embed/text/err_<pkg>.go`. The init-ordering
concern is genuine for `var ErrX = errors.New(desc.Text(...))` — but the
fix is to defer the `desc.Text` call into a method, not to materialise
the English at package init. Identity is preserved because empty-struct
values are comparable and `errors.Is` finds them through `fmt.Errorf("%w", …)`
wrappers.

**Application**: When you need an `errors.Is` target, write:

```go
type missingFooErr struct{}
func (missingFooErr) Error() string {
    return desc.Text(text.DescKeyErrPkgMissingFoo)
}
var ErrMissingFoo error = missingFooErr{}
```

For parameterised errors, follow `internal/err/context/context.go`'s
`NotFoundError` shape: exported struct type with fields, pointer
receiver on `Error()`, `errors.As` at the call site. Never define an
`ErrMsg*` string constant in `internal/config/<pkg>/`; never write
`var ErrX = errors.New("english")`. If you see those, sweep them: text
to YAML, sentinel to typed struct, doc comment justifying the const layer
deleted along with the const.

---

## [2026-04-08-074612] fmt.Fprintf to strings.Builder silently discards errors

**Context**: golangci-lint errcheck allows fmt.Fprintf to strings.Builder
because Write never fails, but project convention says zero silent discard

**Lesson**: Linter coverage gaps exist where language guarantees mask
conventions. AST tests fill the gap

**Application**: Created TestNoUncheckedFmtWrite to enforce fmt.Fprintf error
handling. Use if _, err := fmt.Fprintf(...) with log.Warn on the error path

---

## [2026-04-25-014704] filepath.Join('', rel) returns rel as CWD-relative, not error

**Context**: Recurring orphan jsonl-path-<sessionID> appeared at project root.
Older state.Dir() returned ('', nil) when CTX_DIR was undeclared, so
filepath.Join('', 'jsonl-path-XXX') = 'jsonl-path-XXX', writing relative to CWD.

**Lesson**: Functions returning a path-string must never return ('', nil).
Sentinel errors force callers to gate, closing the silent CWD-relative write.

**Application**: Audit any (string, error) path-returner that historically had a
('', nil) shortcut. Closed for state.Dir and rc.ContextDir; check remaining
resolvers.

---

## [2026-03-06-050125] Package-local err.go files invite broken windows from future agents

**Context**: Found err.go files in 5 CLI packages with heavily duplicated error
constructors (errFileWrite, errMkdir, errZensicalNotFound repeated across
packages)

**Lesson**: Centralizing errors in internal/err eliminates duplication and
prevents agents from continuing the pattern of adding local err.go files when
they see one exists

**Application**: New error constructors go to internal/err/errors.go. No err.go
files in CLI packages.

---

## Group: git CLI wrapping quirks (consolidated)

## [2026-05-22-220100] Group git flag constants by subcommand, not by "loose flags" — cross-group flags enable wrong-subcommand bugs

**Context**: `internal/config/git/git.go` had a constant group commented "Rev-parse flags" that contained `FlagShowCurrent`, but `--show-current` is a `git branch` flag — rev-parse doesn't recognize it. The misclassification meant `internal/gitmeta/branch.go` confidently wrote `Run(cfgGit.RevParse, cfgGit.FlagShowCurrent, ...)` and the call site looked internally consistent at review time: the constants it imported all came from the "Rev-parse flags" group. The bug (literal `branch: --show-current` in handover frontmatter) shipped because the constants file said the flag belonged where it didn't. Fixed in commit 5670f5b2 by splitting `FlagShowCurrent` into a new "Branch subcommand flags" group.

**Lesson**: When flag constants are grouped only by "what command surface they appear on" (e.g. "loose CLI flags") rather than by the subcommand they're actually valid for, future call sites can mix-and-match constants that the comment says are compatible but git rejects. The group comment functions as informal type information; let it tell the truth.

**Application**: In `internal/config/git/git.go` and any similar config package wrapping a CLI's flag surface, group constants by the subcommand whose argv they're valid in (`// Branch subcommand flags`, `// Rev-parse flags`, `// Log subcommand flags`). Flags that genuinely span subcommands (`-C`, `--`) go under a separate "Cross-subcommand flags" group with the spanning explicitly called out. When adding a new flag constant, the first question is "which `git X` subcommand accepts this?" — the answer dictates the group.

---

## [2026-05-22-220000] `git rev-parse` echoes unknown long-flag args back as literal stdout with exit 0 — the error guard never trips

**Context**: `internal/gitmeta.resolveBranchOrDetached` was invoking `git rev-parse --show-current` and returning the result if `runErr == nil`. The function has a defensive fallback (`return BranchDetached` on error), but the error path never fired because rev-parse exits 0 even when handed an unknown long-flag — it just echoes the literal arg back as its only line of output. Result: the resolver returned the string `"--show-current"` verbatim and shipped it into handover frontmatter. Confirmed on git 2.50.0: `$ git rev-parse --show-current` → `--show-current` (exit 0); compare `$ git rev-parse --not-a-real-flag` → same echo-back behavior.

**Lesson**: A non-zero exit guard around a git invocation does NOT catch wrong-subcommand-with-wrong-flag bugs against rev-parse. rev-parse treats unknown args as candidate revision/object names, fails to resolve them, and falls back to echoing them as literal output rather than erroring. Other subcommands (`git branch --bogus`) error loudly with exit ≠ 0; rev-parse specifically is the one that swallows silently. The defensive `if err != nil { return fallback }` pattern is necessary but not sufficient when wrapping rev-parse.

**Application**: When wrapping `git rev-parse`, validate the output shape (e.g. length, prefix, hex-ness for SHAs, no `--` prefix for branch names) before returning, not just the exit code. The `TestResolveHead_RealRepoReturnsBranchName` regression test that landed with the fix asserts both `ref.Branch == "trunk"` AND `!strings.Contains(ref.Branch, "--")` — the second assertion is the one that would catch a future regression where someone reintroduces a different wrong-flag invocation.

---

## [2026-03-24-000959] git describe --tags follows ancestry, not global tag list

**Context**: Release notes skill diffed against v0.3.0 instead of v0.6.0 because
the release branch diverged before v0.6.0 was tagged

**Lesson**: git describe --tags --abbrev=0 follows reachability from HEAD; use
git tag --sort=-v:refname | head -1 for the latest tag globally

**Application**: Any script or skill that needs the latest release should use
sorted tag list, not describe

---

## [2026-04-26-152842] Trailing word boundary in regex matches commit-tree as git commit

**Context**: First post-commit filter regex \bgit\s+commit\b in the OpenCode
plugin would have triggered on git commit-tree because \b matches between t and
-

**Lesson**: A trailing word boundary doesn't exclude hyphenated continuations
— \b matches every word/non-word transition. Use (?!-) negative lookahead to
specifically reject hyphen-suffixed siblings.

**Application**: For any porcelain with hyphenated cousins (commit-tree,
commit-graph, for-each-ref), append (?!-) to the boundary.

---

## Group: TypeScript/integration test surfaces & exclusion rot (consolidated)

## [2026-05-22-161720] Cross-language coverage gap: TS-typed integrations are a fourth surface beyond Go

**Context**: specs/cwd-anchored-context.md removed the CTX_DIR env channel. Three Go test suites caught orphan refs after deletion: audit/TestNoDeadExports (dead consts), audit/TestFlagYAMLMatchesConstants + TestExamplesYAMLLinkage + TestDescKeyYAMLLinkage (orphan YAML keys), compliance/TestDocGoSubcommandDrift (stale doc.go prose). Jumbo commit fc7db228 landed with all four green. But internal/assets/integrations/opencode/plugin/index.ts is a SEPARATE FOURTH surface — TypeScript, not Go — that local 'make lint' and 'go test ./...' never exercise. CI's tsc --noEmit (driven by tools/typecheck/opencode/) surfaced TS2339 on 'output.cwd does not exist on @opencode-ai/plugin shell.env output type'. Fix landed in 40d024a3 but cost a CI round-trip.

**Lesson**: When removing or renaming an env channel, feature flag, or any cross-language contract, the cleanup checklist is FOUR surfaces, not three: (1) Go code (build + lint + test), (2) audit/compliance tests (orphan consts, YAML keys, doc.go drift), (3) asset templates (CLAUDE.md, AGENT_PLAYBOOK, hooks.json, INSTRUCTIONS.md), (4) TypeScript-typed integrations — opencode plugin and the vscode extension. The TS surface is invisible to Go's test suite by design; the typecheck only runs in CI unless invoked explicitly from tools/typecheck/opencode/ or editors/vscode/.

**Application**: Before committing any change that touches internal/assets/integrations/opencode/plugin/ or editors/vscode/, run 'cd tools/typecheck/opencode && npx tsc --noEmit' (and the vscode equivalent). Longer-term: add a 'make typecheck' target wrapping both tsc invocations and include it in the pre-commit checklist alongside 'make lint' and 'go test ./...'. Add it to docs/operations/runbooks/release-checklist.md as a release gate too.

---

## [2026-05-11-202124] tsc cross-tree include resolves node_modules from source file, not tsconfig

**Context**: Set up tsc --noEmit gate for the embedded OpenCode plugin. tsconfig
lived in tools/typecheck/opencode/; include pointed at
internal/assets/integrations/opencode/plugin/index.ts via relative path. First
run failed with 'Cannot find module @opencode-ai/plugin' even though
node_modules was correctly populated in tools/typecheck/opencode/.

**Lesson**: When tsconfig.json sits in dir A but its 'include' points at .ts
files in dir B, tsc resolves node_modules by walking up from each source file's
location (dir B), NOT from the tsconfig's location (dir A). With
moduleResolution: bundler the behavior is the same. The 'node_modules' that
ships in dir A is invisible to a source file in a distant dir B.

**Application**: For any cross-tree tsc setup (typecheck gate for embedded
source elsewhere in the repo, monorepo-style references, etc.), add explicit
baseUrl + paths to the tsconfig. Example: baseUrl: '.', paths: {
'@opencode-ai/plugin': ['./node_modules/@opencode-ai/plugin/dist/index.d.ts'],
'@opencode-ai/plugin/*': ['./node_modules/@opencode-ai/plugin/dist/*'] }. Add
typeRoots ['./node_modules/@types', './node_modules'] for good measure. The cost
is some manual path mapping; the benefit is that node_modules can live wherever
the tooling does, not next to the source.

---

## [2026-05-22-230000] vitest's mocked `execFile` fires callbacks synchronously; real Node defers to `process.nextTick` — closure-capture patterns can TDZ-trap under the mock

**Context**: While scaffolding eslint for `editors/vscode/` (commit 198803de), the `prefer-const` rule flagged `let disposable: T | undefined;` in `runCtx()`. The `disposable` is referenced inside the `execFile` callback (`disposable?.dispose()`) but assigned only after `execFile` returns (the cancellation listener needs `child` to kill, and `child` only exists once `execFile` is called). My refactor: declare `const disposable` after `child = execFile(...)`, and let the inline callback close over `disposable` — relying on Node's `execFile` guarantee that callbacks fire on `process.nextTick` at the earliest (never synchronously, even on immediate-failure paths). This is safe in production. But under vitest, `cp.execFile` is replaced by `vi.mock("child_process")` whose mock callback **fires synchronously** at the point execFile returns. That synchronous invocation reads `disposable` from inside the callback before the `const disposable = ...` line has executed → `ReferenceError: Cannot access 'disposable' before initialization`. Reverted to `let` with an `// eslint-disable-next-line prefer-const` comment.

**Lesson**: vitest's mock factory (`vi.mock("child_process")`) does not preserve Node's async-deferral guarantees. Even APIs that are guaranteed to be asynchronous in production can fire synchronously in the test surface, because the mock is just `vi.fn()` returning a synchronous invocation of whatever the test wires up. This means a closure pattern that's *provably* safe by Node's contract can still TDZ-trap, because the TDZ check happens at runtime regardless of which environment fired the callback. The trap is invisible under typecheck (TypeScript can't reason about callback firing order) and invisible under static analysis (eslint flagged the const opportunity but couldn't see the temporal dependency).

**Application**: When eslint or any analyzer suggests tightening a `let` to `const` in code that captures the variable through an async callback, verify under the *test* runner, not just real-Node semantics. A safe heuristic: if the variable is referenced lexically *before* its declaration (via a closure that fires later), the safe form is `let` with an `eslint-disable-next-line` comment that names the test-mock constraint. Splitting the declaration earlier and assigning later is the lowest-friction pattern that's robust to mock-side synchronicity quirks. The general rule generalizes beyond execFile: any mocked-async API (`fs.readFile`, `dns.lookup`, `http.request`, etc.) can collapse to sync under `vi.mock()`.

---

## [2026-05-22-223000] Double-excluded tests rot compounding — re-enable cost = sum of all drift since last green, not just the original bug

**Context**: `editors/vscode/src/extension.test.ts` was excluded from CI's TypeScript typecheck via `tsconfig.ci.json`'s `**/*.test.ts` glob AND was never run under `npm test` in any CI job. The task to re-enable it (TASKS.md line 228) named two breakages — handler rename (`handleComplete`/`handleTasks` → `handleTask`) and a `fakeToken` listener signature mismatch. Both fixed quickly. But the moment vitest actually executed for the first time in months, 18 additional argv assertions failed: every handler in `extension.ts` had grown an `args.push("--no-color")` call between when the tests were written and now, and not one of those assertions had been updated. `expect.anything()` and `expect.any(Function)` happily passed the typecheck because they admit any shape — the typecheck would not have caught these even if the carve-out had been removed. Only execution did. Commit cf2a109c.

**Lesson**: A test suite excluded from BOTH typecheck and execution rots compounding, not linearly. Every unrelated change in the production code lands without resistance, and the cost of re-enabling is the sum of *all* drift since the suite was last green — not just the bug whose mention triggered the re-enable. The two exclusion layers (typecheck-side `exclude:` and CI-job-side missing-step) each provide false comfort that the other one might be catching something. Together they catch nothing.

**Application**: When adding a tooling exclude of any kind (`tsconfig` exclude glob, `go test ./... -short` skipping a directory, vitest `testPathIgnorePatterns`, `pytest --ignore`), file an immediate follow-up TASKS.md item whose acceptance criterion is *removal* of the exclude with a deadline or trigger. Treat the exclude as borrowed-time, not a stable state. When re-enabling, expect drift-debt: budget for fixing 5–20× more than the named scope and don't ship a partial fix that re-disables on first failure. In code review, an exclude addition without a paired follow-up should be a comment.

---

## Group: Editorial KB pipeline: design epistemology (consolidated)

## [2026-05-10-001859] An ongoing user's concrete workaround tax is the strongest validation evidence

**Context**: When extracting the editorial pipeline, the user pointed at
`your-project` as a project where they were already running the editorial pattern
manually, at concrete cost: CLAUDE.md disabling half of ctx code-dev skills
(/ctx-commit, /ctx-implement, /ctx-spec, /ctx-architecture, /ctx-brainstorm,
/ctx-wrap-up), 10-CONSTITUTION.md at repo root colliding with
.context/CONSTITUTION.md, hand-typed 8-item closeouts, hand-managed 20-INBOX.md,
dedicated reference/vcf/external-grounding.md for ground-mode. The workaround
was visible and the pain was specific.

**Lesson**: An ongoing user paying concrete workaround tax is the strongest
validation evidence; it beats hypothetical user research, beats N=2 design
discussion, beats 'this seems useful.' The shape of the workaround maps directly
to the gap the feature should fill. Validation is essentially complete before
any code is written; the new feature mechanizes what already works manually.

**Application**: When deciding whether to ship a feature, prefer 'a real user is
paying real workaround cost right now' over 'this seems valuable.' Use the
workaround details (which files they created, which conventions they bent, which
skills they disabled) as the inverse-spec of what to build. Ship the feature
shape that exactly matches what they hand-rolled, and use their project as the
regression test corpus (Phase KB-2 ports `your-project` as the validation step).

---

## [2026-05-10-001859] Lift renames alongside features when borrowing from battle-tested external designs

**Context**: When extracting the editorial pipeline from the sibling project,
noticed they named their editorial constitution 10-INGEST_RULES.md (not
10-CONSTITUTION.md), and explicitly recorded a 'domain-decisions.md is named to
disambiguate from .tool/DECISIONS.md (naming-by-rename rule)' note in their
schemas. They had hit and resolved naming conflicts that `your-project` was actively
re-fighting (with 10-CONSTITUTION.md at repo root colliding with
.context/CONSTITUTION.md).

**Lesson**: When lifting from a battle-tested external design, lift the renames
and disambiguation moves alongside the features. Intentional renames encode
resolved conflicts; treating them as cosmetic preferences re-litigates the
underlying fight in your codebase. The aesthetic difference between two names
often hides hard-won architectural learning.

**Application**: ctx editorial pipeline uses KB-RULES.md (not CONSTITUTION.md)
and domain-decisions.md (not DECISIONS.md) explicitly because the sibling did.
For any future external-design lift, scan the source for renames as signal of
resolved-conflict knowledge, and copy them with the rationale (in DECISIONS.md)
so future maintainers don't 'simplify' the names back into the conflict zone.

---

## [2026-05-10-001859] KB epistemology: in a KB you do not decide, you increase confidence

**Context**: Considered whether KB editorial decisions need a parallel
/ctx-kb-decide skill mirroring /ctx-decision-add. Got stuck on three resolutions
(skill surface doubles, mode-aware router, manual discipline) until the user
reframed: do you really decide in a KB, or do you just learn and improve
confidence? A claim with confidence greater than 0.9 is decided by contract;
lower confidence requires more evidence.

**Lesson**: In a knowledge base, the correct ontology has no 'decide' moment;
there are only evidence-capture events with confidence bands. Even
natural-language assertions like 'we are spinning off X, anchor on this' are
semantically evidence-capture (a high-confidence claim arriving), not
decision-capture (a choice between alternatives). The pipeline-only-writer model
is not rigid; it is the ontologically correct surface for evidence-tracked
knowledge.

**Application**: When a feature seems to require a parallel skill mirroring an
existing canonical capture skill, check whether the underlying domain has the
same ontology. If the new domain operates by 'increase confidence' rather than
'pick a choice,' the parallel skill is the wrong shape and the pipeline approach
is right. Useful general check: is this 'I made a call between alternatives' or
'I learned something about the world'? Different ontologies call for different
surfaces.

---

## [2026-05-10-001859] P2: A KB of KBs is a KB

**Context**: User raised 'KB of KBs' as a wished-for federation feature for
multi-team consolidation (research-master KB pulling several team KBs together).
Initial framing treated this as a v2 feature that might require v1 schema
decisions like KB-prefixed IDs (research-master/EV-019) or federation roots.
User reframed: 'kb is knowledge; knowledge is source; source is ingestable;
that's also what makes kb of kbs composable; because kb of kbs is a kb.'

**Lesson**: Recursive composability eliminates whole feature classes. When a
'thing-of-things' feature comes up, ask whether the standard pipeline applied to
its own output covers the case before designing a new mechanism. Federation as
'pipeline pointed at another instance of its own input shape' is dramatically
simpler than federation as a separate subsystem.

**Application**: Federation does not need v1 schema lockout: source-map kind: kb
plus the standard ingest pipeline covers it. Same insight applies to
taxonomy-was-wrong recovery (start fresh KB; ingest old as source; discard
irrelevant parts at extraction time) and multi-team consolidation (each team
owns a KB; master ingests them). Watch for this pattern in future ctx feature
design; the 'thing-of-things is a thing' shortcut may collapse the design
problem entirely.

---

## [2026-05-10-001859] P1: The LLM is the migration tool

**Context**: Designing schemas for the editorial pipeline raised the question of
whether to commit to specific aesthetic choices (EV-### IDs, four named modes,
four-band confidence) or hedge with abstract types that could absorb future
change. The unwind-cost analysis during /ctx-plan showed every category of
being-wrong is essentially cheap because the LLM absorbs the migration:
wholesale ID renumbering (LLM cleanup), taxonomy reshuffles
(start-fresh-and-ingest-old), schema-band remapping (mathematical and
scriptable), path renames (single sweep).

**Lesson**: When designing AI-assisted persistent storage, expensive migrations
are absorbed by LLM cleanup passes. Commit to the readable, opinionated,
aesthetic schema in v1 instead of hedging with abstract types. Be wrong cheaply:
the alternative (hedging upfront) ships a generic shape that nobody loves, and
migrations were never as expensive as we feared.

**Application**: For any future ctx feature where the schema-vs-flexibility
question arises, default to the specific shape; trust LLM cleanup as the
migration story. Surface dirty state via doctor advisories so the agent has a
work surface to operate on. Applies broadly: editorial KB schemas, closeout
shapes, future feature surfaces. Pair with the discipline of doctor flagging
duplicates / divergences so the LLM has clear cases to resolve.

---

## Group: Documentation, template & asset drift (consolidated)

## [2026-03-30-075941] Architecture diagrams drift silently during feature additions

**Context**: During the journal-recall merge, architecture-dia-build.md listed
23 CLI packages but 31 existed. 8 packages added over months without updating
the diagram.

**Lesson**: Exhaustive lists and counts in architecture docs go stale every time
a package is added. The drift is invisible because nobody re-counts.

**Application**: After adding a new CLI package, grep architecture diagrams for
package counts and directory listings. Consider adding a drift-check comment
that validates the count programmatically.

---

## [2026-03-25-173338] Template improvements don't propagate to existing projects

**Context**: 5 of 8 context files in the ctx project itself had stale/missing
comment headers — templates evolved but non-destructive init never re-synced
them

**Lesson**: Any template change is invisible to existing users until they run
ctx init --force

**Application**: Added drift detection (checkTemplateHeaders) to ctx drift.
Consider surfacing this during ctx status too.

---

## [2026-04-01-074419] Copilot CLI skills need a sync mechanism to prevent drift from ctx skills

**Context**: 5 Copilot CLI skills were condensed versions of ctx skills,
independently maintained with no drift detection

**Lesson**: Any time the same content exists in two locations without a sync
mechanism, it will drift silently

**Application**: make sync-copilot-skills added to build deps, make
check-copilot-skills added to audit target

---

## [2026-03-13-151952] sync-why mechanism existed but was not wired to build

**Context**: assets/why/ had drifted from docs/ — the sync targets existed in
the Makefile but build did not depend on sync-why

**Lesson**: Freshness checks that are not in the critical path will be
forgotten. Wire them as build prerequisites, not optional audit steps

**Application**: Any derived or copied asset should be a prerequisite of build,
not just audit

---

## [2026-03-25-234039] Machine-generated CLAUDE.md content consumes per-turn budget without proportional value

**Context**: GitNexus injected 121 lines (61% of CLAUDE.md) with auto-generated
skill pointers like 'Work in the Watch area (39 symbols)' — generic index data
loaded on every conversation turn

**Lesson**: CLAUDE.md is prime real estate — every token competes with
project-specific instructions. Auto-generated content belongs in on-demand
skills, not in always-loaded files

**Application**: Audit CLAUDE.md periodically for content that could be
delivered via skills instead. Prefer a one-line pointer over inline content for
companion tools

---

## [2026-02-26-100000] Documentation drift and auditing (consolidated)

**Consolidated from**: 6 entries (2026-01-29 to 2026-02-24)

- CLI reference docs can outpace implementation: ctx remind had no CLI, ctx
  recall sync had no Cobra wiring, key file naming diverged between docs and
  code. Always verify with `ctx <cmd> --help` before releasing docs.
- Structural doc sections (project layouts, command tables, skill counts) drift
  silently. Add `<!-- drift-check: <shell command> -->` markers above any
  section that mirrors codebase structure.
- Agent sweeps for style violations are unreliable (8 found vs 48+ actual).
  Always follow agent results with targeted grep and manual classification.
- ARCHITECTURE.md missed 4 core packages and 4 CLI commands. The /ctx-drift
  skill catches stale paths but not missing entries — run /ctx-architecture
  after adding new packages or commands.
- Documentation audits must compare against known-good examples and
  pattern-match for the COMPLETE standard, not just presence of any comment.
- Dead link checking belongs in /consolidate's check list (check 12), not as a
  standalone concern. When a new audit concern emerges, check if it fits an
  existing audit skill first.

---

## Group: User-facing text & magic-string discipline (consolidated)

## [2026-04-04-025813] Format-verb strings are localizable text, not exempt from magic string checks

**Context**: Strings like '%d entries checked' were passing TestNoMagicStrings
because the format-verb exemption was too broad

**Lesson**: Any string containing English words alongside format directives is
user-facing text that belongs in YAML assets

**Application**: Removed format-verb, URL-scheme, HTML-entity, and err/
exemptions from TestNoMagicStrings

---

## [2026-03-14-180903] Stderr error messages are user-facing text that belongs in assets

**Context**: Added fmt.Fprintf(os.Stderr) error reporting to event log,
initially with inline strings

**Lesson**: Any string that reaches the user, including stderr warnings, routes
through assets.TextDesc() for i18n readiness

**Application**: When adding stderr output, create text.yaml entries and asset
keys first

---

## [2026-03-31-224247] Magic string cleanup compounds: each pass reveals the next layer

**Context**: What started as fix 4 fmt.Fprintf(os.Stderr) calls expanded to
over-tokenized format strings, magic hex perms, unstandardized TOML parsing
tokens, missing docstrings on new constants — each fix exposed adjacent
violations

**Lesson**: Mechanical cleanup is fractal. The first sweep finds the obvious
violations, but fixing them puts adjacent code under scrutiny. Budget for 2-3x
the initial estimate

**Application**: When scoping cleanup tasks, do not commit to done in one pass.
Commit after each layer and let the user decide when to stop

---

## [2026-03-14-131202] Hardcoded _alt suffixes create implicit language favoritism

**Context**: Session parser had session_prefix_alt hardcoding Turkish as a
special case alongside English default

**Lesson**: Naming a constant _alt and hardcoding one non-English language as a
built-in default discriminates by giving that language special status. The
pattern doesn't scale (alt_2? alt_3?) and signals that adding languages requires
code changes.

**Application**: When a feature needs multi-value support, use configurable
lists from the start — not hardcoded pairs with _alt suffixes. Default to a
single canonical value; all extensions are user-configured equally.

---

## Group: Constant placement & helper smells (consolidated)

## [2026-03-12-133007] Constants belong in their domain package not in god objects

**Context**: file.go held agent scoring constants, budget percentages, cooldown
durations — none related to file config

**Lesson**: When a constant is only used by one domain (e.g. agent scoring), it
should live in that domain's config package

**Application**: Check callers before placing constants; if all callers are in
one domain, the constant belongs there

---

## [2026-03-07-221151] Always search for existing constants before adding new ones

**Context**: Added ExtJsonl constant to config/file.go but ExtJSONL already
existed with the same value, causing a duplicate

**Lesson**: Grep for the value (e.g. '.jsonl') across config/ before creating a
new constant — naming variations (camelCase vs ALLCAPS) make duplicates easy
to miss

**Application**: Before adding any new constant to internal/config, search by
value not just by name

---

## [2026-03-07-221148] SafeReadFile requires split base+filename paths

**Context**: During system/core cleanup, persistence.go passed a full path to
validation.SafeReadFile which expects (baseDir, filename) separately

**Lesson**: Use filepath.Dir(path) and filepath.Base(path) to split full paths
when adapting os.ReadFile calls to SafeReadFile

**Application**: When converting os.ReadFile to SafeReadFile, always check
whether the existing code has a full path or separate components

---

## [2026-03-12-133008] Project-root files vs context files are distinct categories

**Context**: Tried moving ImplementationPlan constant to config/ctx assuming it
was a context file. (Note: IMPLEMENTATION_PLAN.md was removed in 2026-03-25 as a
dead file — no agent consumer.)

**Lesson**: Files created by ctx init in the project root (Makefile) are
scaffolding, not context files loaded via ReadOrder. They belong in config/file,
not config/ctx

**Application**: Before moving a file constant, check whether it is in ReadOrder
(context) or created by init (project-root)

---

## [2026-03-16-022650] One-liner method wrappers hide dependencies without adding value

**Context**: checkBoundary() and loadContext() were methods on Handler that just
called validation.ValidateBoundary and context.Load with h.ContextDir

**Lesson**: If a method only passes a struct field to a stdlib function, inline
it — the wrapper obscures the real dependency

**Application**: Before extracting a helper method, check if it just forwards a
field to another function. If so, call the function directly.

---

## [2026-03-23-003353] Higher-order callbacks in param structs are a code smell

**Context**: MergeParams.UpdateFn and DeployParams.ListErr/ReadErr were function
pointers where all callers passed thin wrappers varying only by a text key

**Lesson**: If all callers pass thin wrappers around the same pattern
(fmt.Errorf with different keys), the callback is just data in disguise

**Application**: When a struct field is a function pointer, check if all callers
vary only by a string key — if so, replace the callback with the key and let
the consumer do the dispatch

---

## Group: Convention enforcement: mechanical gates over prose (consolidated)

## [2026-03-16-104146] Convention enforcement needs mechanical verification, not behavioral repetition

**Context**: Godoc Parameters/Returns sections were missed repeatedly across
sessions despite memory entries and feedback

**Lesson**: System-level brevity instructions outcompete context-injected
conventions. Memory shifts probability (~40% to ~70%) but doesn't create
invariants. The competing pressures are architectural, not a recall problem.

**Application**: Invest in linter rules or PreToolUse gates for
mechanically-checkable conventions. Reserve behavioral nudges for judgment calls
that can't be linted. See ideas/spec-convention-enforcement.md for the
three-tier strategy.

---

## [2026-03-16-114227] Docstring tasks require reading CONVENTIONS.md Documentation section first

**Context**: Agent was asked to review docstrings in server.go but skipped
convention loading, missed incomplete Parameter/Returns sections, and needed
three hints to recall the known issue

**Lesson**: Any task involving docstrings, comments, or documentation formatting
is a convention-sensitive task — read CONVENTIONS.md (Documentation section)
and LEARNINGS.md (for known gaps) before reviewing or writing

**Application**: On any docstring/comment task: (1) load CONVENTIONS.md
Documentation section, (2) check LEARNINGS.md for related entries, (3) audit all
functions in scope against the convention template, not just the ones in the
diff

---

## [2026-03-31-182054] Force-loaded behavioral prose gets ignored — action-gating hooks don't

**Context**: AGENT_PLAYBOOK was force-injected at ~14k tokens every session.
Agent routinely skipped its Context Readback directive when the user's first
message was a concrete task. Meanwhile, hooks that gate actions (qa-reminder,
specs-nudge, block-dangerous-commands) were consistently followed because they
fire at the moment of violation.

**Lesson**: Prose instructions compete with the user's immediate request and
lose. Hooks that intercept actions at execution time are enforceable. More
injected content means less attention per token — slim injection to only what
must be internalized before any action.

**Application**: When adding agent directives, prefer action-gating hooks over
injected prose. If it must be injected, keep it small and directive-only.
Reserve force-injection for hard rules (CONSTITUTION) and distilled actionable
checklists (gate file).

---

## [2026-04-08-074604] AST audit tests must cover unexported functions too

**Context**: TestDocCommentStructure only checked exported functions, so
agent-written helpers in format.go had no godoc enforcement

**Lesson**: Convention enforcement tests must default to scanning all documented
functions. Use explicit opt-outs (test files) not opt-ins (exported only)

**Application**: When adding AST audit tests, scan all functions. We fixed
TestDocCommentStructure to drop the IsExported gate and fixed 84 violations

---

## [2026-04-14-010134] AST stutter test only checks FuncDecl, not GenDecl

**Context**: tpl.TplEntryMarkdown stuttered for a long time because
TestNoStutteryFunctions in internal/audit walks *ast.FuncDecl only; the constant
slipped through.

**Lesson**: The audit suite has a real coverage gap for *ast.GenDecl (consts,
vars, types). Stuttery type/const names will not be caught until the audit is
extended to walk those node kinds.

**Application**: When a stuttery identifier is reported by a human, check both
the offending file and whether the audit can catch it; if not, file an
audit-extension task.

---

## [2026-04-04-025805] Agents add allowlist entries to make tests pass — guard every exemption

**Context**: Found that every exemption map/allowlist in audit tests is a
tempting shortcut for agents

**Lesson**: Added DO NOT widen guard comments to all 10 exemption data
structures across 7 test files

**Application**: Every new audit test with an exemption must include the guard
comment. Review PRs for drive-by allowlist additions.

---

## Group: Go toolchain, gofmt & build-tag pitfalls (consolidated)

## [2026-03-30-003734] Python-generated doc.go files need gofmt — formatter strips bare // padding lines

**Context**: Batch-generated doc.go files used blank // lines for padding, which
gofmt removes as unnecessary whitespace

**Lesson**: Programmatic Go file generation must produce substantive content
lines, not blank comment padding — gofmt enforces this

**Application**: Always run gofmt after any scripted Go file generation

---

## [2026-03-16-022642] Agents reliably introduce gofmt issues during bulk renames

**Context**: Subagents renamed consequences->consequence across 75+ files but
left formatting errors in 12 Go files

**Lesson**: Always run gofmt -l after agent-driven refactors before trusting the
build

**Application**: Add gofmt -w pass as a standard step after any agent-driven
bulk edit

---

## [2026-05-10-181418] Go compile/tool version mismatch comes from the cached toolchain, not the system Go

**Context**: Hit 'compile: version "go1.26.1" does not match go tool version
"go1.26.2"' on every go build / go test / make lint, even with my changes
stashed out. System Go was 1.26.2 (healthy); go.mod pinned 1.26.1, so Go's
auto-toolchain feature had downloaded 1.26.1 to
~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.1.darwin-arm64/. That cached
toolchain was internally inconsistent: its compile binary and stdlib export data
disagreed on version.

**Lesson**: When the compile-vs-tool version error appears, the bug is the
cached toolchain dir, not the installed Go. Reinstalling Go (brew, installer,
etc.) does NOT touch the cached download, so the error persists after reinstall.
Three real fixes: (1) rm -rf
~/go/pkg/mod/golang.org/toolchain@v0.0.1-go<X>.<platform>/ to force a clean
re-download (~30s); (2) bump go.mod to match the system Go so the cached one is
bypassed; (3) GOTOOLCHAIN=go<system version> to override the pin per-invocation.
go clean -cache and GOTOOLCHAIN=local do not help.

**Application**: First diagnostic on this error: check `go env GOROOT`. If it
points to `~/go/pkg/mod/golang.org/toolchain@...` the cached toolchain is in
play. Then either delete the cached dir (most surgical) or bump go.mod (one-line
diff, but lands in a commit). Do not waste time reinstalling Go.

---

## [2026-04-26-152850] make test exit code unreliable due to -cover covdata tooling issue

**Context**: make test exited 1 even with all 123 packages passing on this Go
install; root cause is missing covdata tool when -cover is enabled

**Lesson**: Don't trust make test exit code alone when verifying changes. The
-cover flag in the test target can fail with 'no such tool covdata' even when
every package passes.

**Application**: When make test fails, fall back to 'go test ./...' (no -cover)
and tally ^ok / ^FAIL counts to distinguish real failures from tooling issues.

---

## [2026-04-01-233248] go/packages respects build tags — darwin-only violations invisible on Linux

**Context**: TestNoExecOutsideExecPkg could not detect violations in _darwin.go
files when running on Linux

**Lesson**: AST checks using go/packages only see files matching the current
GOOS. Cross-platform violations need either multi-GOOS CI or a go/parser
fallback

**Application**: When writing audit checks for code with build tags, fix the
violations regardless (code correctness) but note that test coverage is
platform-dependent

---

## Group: Stale-task triage & verify-before-acting (consolidated)

## [2026-05-23-003000] Closing a stale TASKS.md item often means writing the test, not the code — verify before assuming the work is undone

**Context**: TASKS.md line 375 ("Improve hub failover client: distinguish auth errors from connection errors") had been open since 2026-04-08. On triage, `internal/hub/failover.go:61-63` already called `authErr(callErr)` and returned immediately on Unauthenticated/PermissionDenied; `internal/hub/err_check.go:22-30` `authErr()` checked exactly those two codes. The behavior was implemented in the original failover feature commit (8bcb6208) without the task being closed. But the test suite never asserted the invariant — three existing failover tests covered happy path, skip-bad-peer, and all-bad-peers, none of them exercised "auth fails → walk stops". A future refactor could have silently deleted the auth-fast-fail branch and all three would still pass. Commit 22cffc27 added `TestFailoverClient_FailsFastOnAuthError` and closed the task.

**Lesson**: Stale TASKS.md items frequently describe work that's *already done in code* but *not asserted in tests*. The task stays open not because nothing happened but because nothing pinned the behavior down so the task author could mark it complete. Reading a task description and assuming the code surface is missing is a misdiagnosis. The right pattern: `git log` / `git blame` / grep the symbols the task names; if the implementation exists, the task's value shifts from "build the thing" to "lock the thing down with a test that would catch its regression". Closes the task AND defends the behavior.

**Application**: When triaging TASKS.md, especially items older than a few weeks, run a "what's the implementation status?" sweep before scoping work. For each candidate: grep the function/file/behavior the task names; if it exists, check the test file for an assertion that exercises the named invariant (not just adjacent invariants). If the assertion is missing, the task closes by writing the regression test — frequently a single test function. This pattern applies to behavior-named tasks ("X should fail fast on Y", "Z should reject malformed W") much more than to feature-named tasks ("add the X command"). For ctx specifically, hub/connect/replication-adjacent tasks accreted this way during the original implementation push; the failover-auth task was one example, others (file locking on connect sync, fanout broadcast entry loss) are still on TASKS.md and may warrant the same triage.

---

## [2026-03-01-133014] Task descriptions can be stale in reverse — implementation done but task not marked complete

**Context**: ctx recall sync task said 'command is not registered in Cobra' but
the code was fully wired and all tests passed. The task description was stale.

**Lesson**: Tasks can become stale in the opposite direction from docs:
implementation gets completed but the task is not updated. Always verify with
ctx <cmd> --help before assuming work remains.

**Application**: Before starting implementation on a 'code exists but not wired'
task, run the command first to check if it already works.

---

## [2026-03-15-040642] Grep for callers must cover entire working tree before deleting functions

**Context**: Deleted 7 err/prompt functions as dead code, but callers existed in
unstaged refactoring files — caused build failures

**Lesson**: When the working tree has unstaged changes from a prior session,
grep hits only committed+staged code; must grep the full tree or build-test
before declaring functions dead

**Application**: Always run make build after deleting functions, even if grep
shows zero callers

---

## [2026-05-23-100000] Spec-trailer improvisation is heuristic drift — when no spec genuinely fits, the failure mode is reaching for the most-recent one

**Context**: Two commits on the `fix/journal-schema-drift` branch (a schema fix at `b84bc8e0` and a gitignore chore at `292e12ae`) both cited `ideas/spec-companion-intelligence.md` as their `Spec:` trailer. Neither commit had anything to do with companion intelligence (peer-MCP RAG integration). The agent had reached for that spec because it was the most recently mentioned spec in working memory from the previous commit's reasoning — not because it covered the work. The user caught the mismatch on review: "The spec you tagged has NOTHING TO DO with the commit." Audit of the session's trailers showed 2 genuinely wrong and ~4 stretches in 16 commits — a sustained drift pattern, not a one-off slip.

**Lesson**: When the CONSTITUTION mandates a `Spec:` trailer on every commit AND a particular commit has no on-topic spec available, the agent's path-of-least-resistance heuristic converges on "cite the most recent spec from context" because the local cost (scaffold a new spec) is higher than the local benefit (gate passes). The convergence satisfies the syntactic check (trailer present) but defeats the rule's semantic intent (truthful traceability). This is "heuristic drift" in the gradient-descent sense: the optimizer found a path that minimizes friction but not the loss function the rule was meant to enforce. The drift is silent — the trailer looks fine in `git log` unless a reader opens the cited spec and discovers the mismatch.

The deeper insight from this incident: session-scoped commitments ("I'll be more careful next time") do not survive across agent sessions. A fresh Claude Code session loads the project's persistent context (CONSTITUTION, AGENT_PLAYBOOK, LEARNINGS, files) but has no memory of any earlier session's self-imposed discipline. The structural fix must therefore live in persistent context, not in agent intention.

**Application**: When the closest candidate spec is the same as the previous commit's spec AND the work is qualitatively different, treat that as a red flag and stop. The Spec Verification Step in `AGENT_PLAYBOOK.md` (added 2026-05-23 in commit landing this learning) is the procedure: name the spec, articulate the overlap in one non-hand-waving sentence, and if you can't, choose one of three correct responses — scaffold a fresh spec, bundle the change into the next functional commit, or cite `specs/meta/chores.md` if the diff fits an explicitly listed chore category. Improvisation is no longer an option because the playbook closes that door. The CONSTITUTION's spec-trailer rule (`CONSTITUTION.md` Process Invariants) now also names the chore escape hatch and the verification gate explicitly. Both changes serve the same goal: remove the conditions under which improvisation can happen in the first place. See `specs/spec-trailer-discipline.md` for the design rationale.

---

## Group: Refactor mechanics: subagents, cascades & golden fixtures (consolidated)

## [2026-05-21-140230] Sentinel-removal refactors cascade through test surface

**Context**: Spec specs/cwd-anchored-context.md decomposed the work into 5 discrete steps; in practice steps 1 and 2 had to merge. Removing ErrDirNotDeclared from rc.ContextDir cascaded through ~10 errors.Is consumers and ~30 test fixtures that used t.Setenv(env.CtxDir, ...).

**Lesson**: Spec-level decomposition that treats 'swap resolver' and 'remove init guard' as separable does not survive contact when the second step references the soon-to-be-deleted sentinel from the first. Both have to compile against the new sentinel set in the same commit.

**Application**: When a future spec proposes step boundaries that hinge on a sentinel rename or removal, plan the merged commit up front rather than discover the cascade mid-implementation. The compile-surface analysis belongs at spec time, not implementation time.

---

## [2026-05-17-061000] Subagent parallelism shines for mechanical refactor with a worked-example reference

**Context**: Phase KB audit cleanup spanned 428 violations across 21 categories
in ~50 files. Doing it serially in the orchestrator would have burned the
session. Three subagents in parallel (one for 16 markdown templates, one for 10
schemas, one for 6 SKILL.md files) landed 32 files with zero integration churn.
A fourth subagent (9 kb writer packages) and a fifth (CLI cmd tree) followed the
same shape and cleared the bulk of audit failures while the orchestrator handled
handover + gitmeta + closeout itself.

**Lesson**: Subagents work well when (a) the work is well-bounded, (b) a
canonical worked example exists in the prompt or on disk, (c) the agent is told
to fix-or-fail-with-a-blocker rather than surface deferral options. The first
subagent I dispatched stopped at honest-scope reporting; the followups plowed
because the prompt explicitly invoked the Constitution's no-deferral rule and
pointed at a worked example.

**Application**: For mechanical refactor work at scale: do one worked example in
the orchestrator, then dispatch a subagent for the rest with the example as a
reference path in the prompt. Tell the subagent to either complete the work or
surface a specific blocker with a concrete next step, not options for the user
to choose between.

---

## [2026-05-30-212109] Capture golden fixtures from the live legacy code path before deleting it

**Context**: Behavior-preserving refactors of LoopScript composition and the recall <details>/<table> assembly had fragile whitespace where hand-transcribing the expected output risked silent drift from the original bytes.

**Lesson**: A throwaway test that runs the current (pre-refactor) code and writes its output to testdata/*.golden gives a regression baseline derived from real behavior, not a re-transcription; delete the throwaway, then have the committed test assert the new code is byte-identical to the fixtures.

**Application**: Use for any behavior-preserving refactor of formatting/rendering code: capture goldens from the legacy path before removing it, then assert byte-equality after.

---

## [2026-03-23-003544] Splitting core/ into subpackages reveals hidden structure

**Context**: init core/ was a flat bag of domain objects — splitting into
backup/, claude/, entry/, merge/, plan/, plugin/, project/, prompt/, tpl/,
validate/ exposed duplicated logic, misplaced types, and function-pointer
smuggling that were invisible in the flat layout

**Lesson**: Flat core/ packages hide coupling — circular dependency resolution
during splits naturally groups related items, increases cohesion, and surfaces
objects that don't belong

**Application**: When a core/ package grows, split it into subpackages even if
it creates temporary circular deps — resolving those deps is the design work
that reveals the right structure

---

## [2026-04-03-180000] Subagent scope creep and cleanup (consolidated)

**Consolidated from**: 4 entries (2026-03-06 to 2026-03-23)

- Subagents reliably rename functions, restructure files, change import aliases,
  and modify function signatures beyond their stated scope — even narrowly
  scoped tasks like fixing em-dashes in comments
- Subagents create new files during refactors but consistently fail to delete
  the originals — always audit for stale files, duplicate definitions, and
  orphaned imports afterward
- After any agent-driven refactor: run `git diff --stat` and `git diff
  --name-only HEAD`, revert anything outside the intended scope, and check for
  stale package declarations before building

---

## [2026-04-03-180000] Cross-cutting change ripple (consolidated)

**Consolidated from**: 4 entries (2026-02-19 to 2026-03-01)

- Path changes (e.g. key file location) ripple across 15+ doc files and 2 skills
  — grep broadly (not just code) and budget for 15+ file touches
- Removing embedded asset directories requires synchronized cleanup across 5+
  layers: embed directive, accessor functions, callers, tests, config constants,
  build targets, documentation — work outward from the embed
- Absorbing shell scripts into Go commands creates a discoverability gap —
  update contributing.md, common-workflows.md, and CLI index as part of the
  absorption checklist
- A feature without docs is invisible to users: always check feature page,
  cli-reference.md, relevant recipes, and zensical.toml nav after implementing a
  new CLI subcommand

---

## Group: Linting, gosec & I/O chokepoints (consolidated)

## [2026-03-01-222739] Gosec G306 flags test file WriteFile with 0644 permissions

**Context**: New tests used os.WriteFile(..., 0o644) for temp context files;
lint flagged all three occurrences

**Lesson**: Gosec enforces 0600 max on WriteFile even in test code. Use 0o600
for test temp files

**Application**: Default to 0o600 for os.WriteFile in tests; only use wider
permissions when testing permission behavior specifically

---

## [2026-04-03-180000] Lint suppression and gosec patterns (consolidated)

**Consolidated from**: 4 entries (2026-03-04 to 2026-03-19)

- Rename constants to avoid gosec G101 false positives (Tokens->Usage,
  Passed->OK) instead of adding nolint/nosec/path exclusions — exclusions
  break on file reorganization
- `nolint:goconst` for trivial values normalizes magic strings — use config
  constants instead of suppressing the linter
- `nolint:errcheck` in tests teaches agents to spread the pattern to production
  code — use `t.Fatal(err)` for setup, `defer func() { _ = f.Close() }()` for
  cleanup
- golangci-lint v2 ignores inline nolint directives for some linters — use
  config-level `exclusions.rules` for gosec patterns, fix the code instead of
  suppressing errcheck

---

## [2026-02-22-120002] Linting and static analysis (consolidated)

**Consolidated from**: 7 entries (2026-01-25 to 2026-02-20)

- Full pre-commit gate: (1) `CGO_ENABLED=0 go build ./cmd/ctx`, (2)
  `golangci-lint run`, (3) `CGO_ENABLED=0 go test` — all three, every time
- Own the codebase: fix pre-existing lint issues even if you didn't introduce
  them
- gosec G301/G306: use 0o750 for dirs, 0o600 for files everywhere including
  tests
- gosec G304 (file inclusion): safe to suppress with `//nolint:gosec` in test
  files using `t.TempDir()` paths
- golangci-lint errcheck: use `cmd.Printf`/`cmd.Println` in Cobra commands
  instead of `fmt.Fprintf`
- `defer os.Chdir(x)` fails errcheck; use `defer func() { _ = os.Chdir(x) }()`
- golangci-lint Go version mismatch in CI: use `install-mode: goinstall` to
  build linter from source

---

## [2026-04-01-233250] Raw I/O migration unlocks downstream checks for free

**Context**: TestNoRawPermissions had zero violations because the raw I/O
migration moved all octal literals into internal/io/ which already used
config/fs constants

**Lesson**: Chokepoint migrations have cascading benefits — centralizing one
concern (file I/O) automatically resolves other drift (raw permissions)

**Application**: Prioritize chokepoint migrations (io, exec, write, err) before
smaller checks that depend on them

---

## Group: Hook mechanics, output channels & compliance (consolidated)

## [2026-04-06-204226] Agents ignore system-reminder content without explicit relay instructions

**Context**: Provenance line (Session: abc | Branch: main @ hash) was emitted by
hook but agents in other projects silently ignored it. The line appeared in the
system-reminder but the agent treated it as internal metadata.

**Lesson**: Claude Code surfaces hook stdout as system-reminder tags. Agents
only relay content that has explicit display instructions. IMPORTANT: means pay
attention internally. Display this line verbatim means show to user. Without the
instruction, even correct output is invisible to the user.

**Application**: Any hook output intended for the user must include an explicit
relay instruction like Display this line verbatim at the start of your response.
Do not rely on IMPORTANT: alone — it signals internal priority, not
user-facing output.

---

## [2026-02-22-120000] Hook behavior and patterns (consolidated)

**Consolidated from**: 8 entries (2026-01-25 to 2026-02-17)

- Hook scripts receive JSON via stdin (not env vars); parse with
  `HOOK_INPUT=$(cat)` then jq
- Hook key names are case-sensitive: `PreToolUse` and `SessionEnd` (not
  `PreToolUseHooks`)
- Use `$CLAUDE_PROJECT_DIR` in hook paths, never hardcode absolute paths
- Hook regex can overfit: `ctx` as binary vs directory name differ; anchor
  patterns to command-start positions with `(^|;|&&|\|\|)\s*`
- grep patterns match inside quoted arguments — test with `ctx add learning
  "...blocked words..."` to verify no false positives
- Hook scripts can silently lose execute permission; verify with `ls -la
  .claude/hooks/*.sh` after edits
- Two-tier output is sufficient: unprefixed (agent context, may or may not
  relay) and `IMPORTANT: Relay VERBATIM` (guaranteed relay); don't add new
  severity prefixes
- Repeated injection causes agent repetition fatigue; use `--session $PPID
  --cooldown 10m` and pair with a readback instruction

---

## [2026-02-22-120001] UserPromptSubmit hook output channels (consolidated)

**Consolidated from**: 2 entries (2026-02-12)

- UserPromptSubmit hook stdout is prepended as AI context (not shown to user);
  stderr with exit 0 is swallowed entirely
- User-visible output requires `{"systemMessage": "..."}` JSON on stdout
  (warning banner) or exit 2 (blocks prompt)
- There is no non-blocking user-visible output channel for this hook type
- Design hooks for their actual audience: AI-facing = plain stdout, user-facing
  = systemMessage JSON

---

## [2026-02-26-100009] Hook compliance and output routing (consolidated)

**Consolidated from**: 3 entries (2026-02-22 to 2026-02-25)

- Plain-text hook output is silently ignored by the agent. Claude Code parses
  hook stdout starting with `{` as JSON directives; plain text is disposable.
  All hooks should return JSON via `printHookContext()`.
- Hook compliance degrades on narrow mid-session tasks (~15-25% partial skip
  rate). Root cause: CLAUDE.md's "may or may not be relevant" system reminder
  competes with hook authority. Fix: CLAUDE.md explicitly elevates hook
  authority. The mandatory checkpoint relay block is the compliance canary.
- No reliable agent-side before-session-end event exists. SessionEnd fires after
  the agent is gone. Mid-session nudges and explicit /ctx-wrap-up are the only
  reliable persistence mechanisms.

---

## [2026-02-27-002830] Context injection and compliance strategy (consolidated)

**Consolidated from**: 3 entries (2026-02-26)

- Verbal summaries with linked diagram files cut ARCHITECTURE.md from ~12K to
  ~3.8K tokens. Extract diagrams to linked files outside FileReadOrder; keep
  prose summaries inline. The 4-chars-per-token estimator is accurate —
  optimize content, not the estimator.
- Soft instructions have a ~75-85% compliance ceiling because "don't apply
  judgment" is itself evaluated by judgment. When 100% compliance is required,
  don't instruct — inject via `additionalContext`. Reserve soft instructions
  for ~80% acceptable compliance.
- Once ~7K tokens are auto-injected (fait accompli), the agent's rationalization
  inverts from "skip to save effort" to "marginal cost is trivial." Front-load
  highest-value content as injection, then use sunk cost to motivate on-demand
  reads for the remainder.

---

## Group: State, tombstones, logs & filesystem hygiene (consolidated)

## [2026-02-22-120006] Permission and settings drift (consolidated)

**Consolidated from**: 4 entries (2026-02-15)

- Permission drift is distinct from code drift — settings.local.json is
  gitignored, no review catches stale entries
- `Skill()` permissions don't support name prefix globs — list each skill
  individually
- Wildcard trusted binaries (`Bash(ctx:*)`, `Bash(make:*)`), but keep git
  commands granular (never `Bash(git:*)`)
- settings.local.json accumulates session debris; run periodic hygiene via
  `/sanitize-permissions` and `/ctx-drift`

---

## [2026-02-22-120008] Gitignore and filesystem hygiene (consolidated)

**Consolidated from**: 3 entries (2026-02-11 to 2026-02-15)

- Gitignored directories are invisible to `git status`; stale artifacts persist
  indefinitely — periodically `ls` gitignored working directories
- Add editor artifacts (*.swp, *.swo, *~) to .gitignore alongside IDE
  directories from day one
- Gitignore entries for sensitive paths are security controls, not documentation
  — never remove during cleanup sweeps

---

## [2026-03-05-205422] State directory accumulates silently without auto-prune

**Context**: Found 234 files in .context/state/ from weeks of sessions with no
cleanup mechanism

**Lesson**: Session tombstones are write-only. Without auto-prune, the state
directory grows unbounded. Added autoPrune(7) to context-load-gate so cleanup
happens once per session at startup.

**Application**: Auto-prune is now wired into session start via
context-load-gate. Manual prune still available via ctx system prune for
aggressive cleanup.

---

## [2026-03-05-205419] Global tombstones suppress hooks across all sessions

**Context**: Memory drift nudge used memory-drift-nudged with no session ID in
filename

**Lesson**: Any tombstone file intended to be session-scoped must include the
session ID in its filename, otherwise it suppresses across all concurrent and
future sessions. Use the UUID pattern so prune can clean them up.

**Application**: Audit all tombstone files for session-scoping; fixed
memory-drift, but backup-reminded, ceremony-reminded, check-knowledge,
journal-reminded, version-checked, ctx-wrapped-up still have this bug

---

## [2026-03-01-092611] Hook logs had no rotation; event log already did

**Context**: Investigated .context/logs/ and .context/state/ file management

**Lesson**: eventlog already rotates at 1MB with one previous generation.
logMessage() in state.go was pure append-only with no size check.

**Application**: When adding new log sinks, follow the established rotation
pattern (size-based, single previous generation)

---

## [2026-03-06-141506] Stale directory inodes cause invisible files over SSH

**Context**: Files created by Claude Code hooks were visible inside the VM but
not from the SSH terminal

**Lesson**: If a directory is recreated (e.g. by auto-prune), an SSH shell
holding the old directory inode will not see new files — ls returns no such
file even though cat with the full path works from other shells

**Application**: After ctx system prune or any state directory recreation, SSH
sessions need cd-dot or re-login to pick up the new inode

---

## Group: Host-pressure alerting: use derivatives, not levels (consolidated)

## [2026-05-28-201500] Swap occupancy is not memory pressure — use the kernel's derivative

**Context**: ctx's `check-resource` UserPromptSubmit hook alerted DANGER at swap-used ≥ 75% / memory-used ≥ 90%, generating false "wrap up the session" warnings at session start after hibernation. On macOS, swap doesn't recede when pressure ends — it's a sticky high-water mark, so static occupancy carries zero current information about whether the system is actually struggling.

**Lesson**: macOS and Windows swap proactively, and swap occupancy is STICKY — it doesn't recede when pressure ends. After hibernation, swap can be >75% full with zero current pressure. Any alert keyed on `SwapUsed/SwapTotal ≥ X%` will false-positive at session start. The signal isn't the *level*, it's the *derivative* — pages actively being pushed out, or the kernel's own pressure metric.

**Application**: For host-pressure detection, key on OS-native pressure signals (macOS `kern.memorystatus_vm_pressure_level` 1/2/4 → OK/Warning/Danger; Linux PSI `/proc/pressure/memory` `some.avg10` and `full.avg10`). These are kernel-computed derivatives — no snapshot state needed and they collapse to zero when the pressure ends. If native is unavailable, fall back to swap-out RATE (snapshot delta) gated on low available memory; never to occupancy alone. (Decision recorded same date; Windows exploratory task filed under Phase CLI-FIX.)

---

## [2026-04-13-153618] Load average measures a queue, not CPU utilization

**Context**: The 'Load Xx CPU count' resource alert fired at 1.74x while htop
showed per-core utilization well under 50% and idle cores. Load average counts
runnable + uninterruptible-sleep processes, smoothed over 1/5/15 minutes.

**Lesson**: Load average and CPU% measure different things. High load with low
CPU% typically means many short-lived processes or I/O-bound work (e.g., go test
spawning hundreds of parallel test binaries). The 1-minute average is too
reactive for dev machines that periodically run test suites — 5-minute smooths
transient spikes without hiding sustained pressure.

**Application**: For alerting thresholds based on system load, prefer 5-minute
over 1-minute averages. 1-minute is useful for interactive debugging; 5-minute
is better for automated alerts that should not fire on normal build/test
activity.

---

## Group: Go test isolation & patterns (consolidated)

## [2026-03-01-161459] Test HOME isolation is required for user-level path functions

**Context**: After adding ~/.ctx/.ctx.key as global key location, test suites
wrote real files to the developer home directory

**Lesson**: Any code that uses os.UserHomeDir() needs t.Setenv(HOME, tmpDir) in
tests — especially test helpers called by many tests (like setupEncrypted and
helper)

**Application**: When adding features that write to user-level paths (~/.ctx/,
~/.config/), always add HOME isolation to test setup functions first

---

## [2026-04-25-014704] Parallel go test ./... packages can race on ~/.claude/settings.json

**Context**: make test runs packages in parallel processes. Fourteen test files
invoked initialize.Cmd().Execute(), which read-modify-writes
~/.claude/settings.json without HOME isolation.

**Lesson**: Under load the races materialized as flaky 'FAIL coverage: [no
statements]' in cli/watch/core. Run alone the package passed; under parallel
make test it failed intermittently.

**Application**: testctx.Declare now sets HOME alongside CTX_DIR. Centralized
fix; future tests automatically isolate user-home writes.
## [2026-02-26-100005] Go testing patterns (consolidated)

**Consolidated from**: 7 entries (2026-01-19 to 2026-02-26)

- Compiler-driven refactoring misses test files: `go build ./...` catches
  production callsite breaks but not test files. Always run `go test ./...`
  after signature changes.
- All runCmd() returns must be consumed in tests: even setup calls need `_, _ =
  runCmd(...)` to satisfy errcheck.
- Set `color.NoColor = true` in a package-level init function to disable ANSI
  codes for CLI test string assertions.
- Recall CLI tests isolate via HOME env var: `t.Setenv("HOME", tmpDir)` with
  `.claude/projects/` structure gives full isolation from real session data.
- `formatDuration` accepts an interface with a Minutes method, not time.Duration
  directly. Use a stubDuration struct for testing.
- CI tests need `CTX_SKIP_PATH_CHECK=1` env var because init checks if ctx is in
  PATH.
- CGO must be disabled for ARM64 Linux (`CGO_ENABLED=0`) — CGO causes
  cross-compilation issues with `-m64` flag.

---

## [2026-03-01-222738] Converting PersistentPreRun to PersistentPreRunE changes exit behavior

**Context**: Boundary violation test used subprocess pattern because original
code called os.Exit(1)

**Lesson**: With PersistentPreRunE, errors propagate through Cobra Execute()
return — no os.Exit call. Subprocess-based tests that expected exit codes need
converting to direct error assertions

**Application**: When converting PreRun to PreRunE in Cobra commands, audit all
tests that relied on os.Exit behavior

---

## Archived: stale / superseded

## [2026-03-31-112534] Legacy key directory cleanup was specified but not automated

**Context**: ~/.local/ctx/keys/ accumulated 584 orphan keys from test runs
before the v0.8.0 migration to ~/.ctx/.ctx.key

**Lesson**: Migration specs that call for manual cleanup of old paths should
include an automated step — either in the migration code itself or as a
post-release cleanup task. Tests that write to global paths must isolate HOME.

**Application**: When writing migration specs, always include automated cleanup
of the old path. When writing tests that touch user-level directories, verify
HOME is isolated via t.Setenv.

---

## [2026-03-02-123613] Existing Projects is ambiguous framing for migration notes

**Context**: A doc admonition said Existing Projects: if you have an older key
at X, it auto-migrates. Every project is existing once installed — the framing
does not tell you how far behind you need to be.

**Lesson**: Version-anchored framing (Key Folder Change v0.7.0+) is clearer than
relative framing (Existing Projects, Legacy). State the version boundary and the
concrete action.

**Application**: When writing migration notes, anchor to a version number and
give copy-pasteable commands, not vague auto-handled assurances.

---

## [2026-03-05-042157] Claude Code has two separate memory systems behind feature flags

**Context**: Filesystem and behavioral analysis of Claude Code v2.1.69

**Lesson**: Claude Code has two separate memory systems behind feature flags.
Auto memory writes MEMORY.md to disk (user-visible, toggleable via settings).
Session memory is a separate background extraction pipeline with compaction and
team sync (push/pull model). The two systems serve different purposes and are
independently feature-flagged.

**Application**: ctx memory bridge targets auto memory (MEMORY.md on disk).
Session memory is API-side and not directly accessible. Full findings in
ideas/claude-code-project-directory-structure.md.

---

## [2026-03-06-184820] Claude Code supports PreCompact and SessionStart hooks that ctx does not use

**Context**: context-mode proves both hooks work in production across 5
platforms

**Lesson**: ctx's hook architecture only uses UserPromptSubmit, PreToolUse, and
PostToolUse — two lifecycle events are untapped

**Application**: PreCompact snapshot plus SessionStart re-injection would
eliminate post-compaction disorientation without any new persistence layer since
ctx agent already generates the content

---

## [2026-03-02-005217] Claude Code JSONL model ID does not distinguish 200k from 1M context

**Context**: Heartbeat hook was reporting 16% usage at 162k tokens because it
assumed claude-opus-4-6 always has 1M context window

**Lesson**: The JSONL model field is identical for both variants (both report
claude-opus-4-6). The 1M context requires a beta header, not a different model
ID. The user's model selection is stored in ~/.claude/settings.json with a [1m]
suffix when 1M is active.

**Application**: Auto-detect context window from ~/.claude/settings.json model
field containing [1m]. Default to 200k for all Claude models. The .ctxrc
context_window setting is a no-op for Claude Code users.

---

