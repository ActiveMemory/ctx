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
| 2026-05-22 | vitest's mocked `execFile` fires callbacks synchronously; real Node defers to `process.nextTick` — closure-capture patterns can TDZ-trap under the mock |
| 2026-05-22 | Double-excluded tests rot compounding — re-enable cost = sum of all drift since last green, not just the original bug |
| 2026-05-22 | Group git flag constants by subcommand, not by "loose flags" — cross-group flags enable wrong-subcommand bugs |
| 2026-05-22 | `git rev-parse` echoes unknown long-flag args back as literal stdout with exit 0 — the error guard never trips |
| 2026-05-22 | Cross-language coverage gap: TS-typed integrations are a fourth surface beyond Go |
| 2026-05-21 | Sentinel-removal refactors cascade through test surface |
| 2026-05-20 | macOS /var symlink trips path-equality; use EvalSymlinks with parent-resolution fallback |
| 2026-05-20 | Handover filenames are archaeology; parse by generated-at, not filename |
| 2026-05-20 | /ctx-plan is named after its input, not its output |
| 2026-05-17 | Creator confusion is the strongest doc-quality signal — louder than any user signal |
| 2026-05-17 | Sentinel errors use typed zero-data structs with lazy `desc.Text()` — never Go string consts |
| 2026-05-17 | `_helpers.go` / `_utils.go` filenames are project anti-pattern; use domain nouns |
| 2026-05-17 | Subagent parallelism shines for mechanical refactor with a worked-example reference |
| 2026-05-17 | naked_errors audit rejects fmt.Errorf wrapping outside internal/err/<area>/ |
| 2026-05-17 | Pre-emptive constants are dead exports; ship constants only when their caller lands |
| 2026-05-11 | Naive Markdown line-sweep corrupts multi-line code spans and YAML lists |
| 2026-05-11 | tsc cross-tree include resolves node_modules from source file, not tsconfig |
| 2026-05-10 | Go compile/tool version mismatch comes from the cached toolchain, not the system Go |
| 2026-05-10 | An ongoing user's concrete workaround tax is the strongest validation evidence |
| 2026-05-10 | Lift renames alongside features when borrowing from battle-tested external designs |
| 2026-05-10 | KB epistemology: in a KB you do not decide, you increase confidence |
| 2026-05-10 | P2: A KB of KBs is a KB |
| 2026-05-10 | P1: The LLM is the migration tool |
| 2026-05-08 | Cursor imports Claude Code hooks and sets CLAUDE_PROJECT_DIR per workspace |
| 2026-04-14 | Constitution forbids context window as a deferral excuse |
| 2026-04-14 | docs/cli/system.md and embed/cmd/system.go diverged on bootstrap promotion intent |
| 2026-04-14 | Raft-lite trade-off is the load-bearing choice in internal/hub |
| 2026-04-14 | AST stutter test only checks FuncDecl, not GenDecl |
| 2026-04-14 | Brand-name handling in title-case engines must cover possessives |
| 2026-04-13 | GPG signing from non-TTY contexts requires pinentry-mac (or equivalent) |
| 2026-04-13 | Load average measures a queue, not CPU utilization |
| 2026-04-13 | rc.ContextDir() is the single source of truth — fix the resolver, not callers |
| 2026-04-09 | Pad index shifting is a real UX bug in batch operations |
| 2026-04-08 | fmt.Fprintf to strings.Builder silently discards errors |
| 2026-04-08 | AST audit tests must cover unexported functions too |
| 2026-04-06 | Agents ignore system-reminder content without explicit relay instructions |
| 2026-04-04 | Format-verb strings are localizable text, not exempt from magic string checks |
| 2026-04-04 | Agents add allowlist entries to make tests pass — guard every exemption |
| 2026-04-03 | Subagent scope creep and cleanup (consolidated) |
| 2026-04-03 | Bulk rename and replace_all hazards (consolidated) |
| 2026-04-03 | Import cycles and package splits (consolidated) |
| 2026-04-03 | Lint suppression and gosec patterns (consolidated) |
| 2026-04-03 | Skill lifecycle and promotion (consolidated) |
| 2026-04-03 | Cross-cutting change ripple (consolidated) |
| 2026-04-03 | Dead code detection (consolidated) |
| 2026-04-03 | desc.Text() is the single highest-connectivity symbol in the codebase |
| 2026-04-01 | Raw I/O migration unlocks downstream checks for free |
| 2026-04-01 | go/packages respects build tags — darwin-only violations invisible on Linux |
| 2026-04-01 | Copilot CLI skills need a sync mechanism to prevent drift from ctx skills |
| 2026-04-01 | Contributor PRs based on older code reintroduce removed features |
| 2026-03-31 | Magic string cleanup compounds: each pass reveals the next layer |
| 2026-03-31 | Force-loaded behavioral prose gets ignored — action-gating hooks don't |
| 2026-03-31 | Legacy key directory cleanup was specified but not automated |
| 2026-03-31 | Convention audits must check cmd/ purity, not just types and docstrings |
| 2026-03-31 | JSON Schema default fields cause linter errors with some validators |
| 2026-03-30 | Architecture diagrams drift silently during feature additions |
| 2026-03-30 | Python-generated doc.go files need gofmt — formatter strips bare // padding lines |
| 2026-03-30 | lint-docstrings.sh greedy sed hid all return-type violations |
| 2026-03-25 | Machine-generated CLAUDE.md content consumes per-turn budget without proportional value |
| 2026-03-25 | Template improvements don't propagate to existing projects |
| 2026-03-24 | lint-drift false positives from conflating constant namespaces |
| 2026-03-24 | git describe --tags follows ancestry, not global tag list |
| 2026-03-23 | Typography detection script needs exclusion lists for intentional uses |
| 2026-03-23 | Splitting core/ into subpackages reveals hidden structure |
| 2026-03-23 | Higher-order callbacks in param structs are a code smell |
| 2026-03-20 | Commit messages containing script paths trigger PreToolUse hooks |
| 2026-03-18 | Lazy sync.Once per-accessor is a code smell for static embedded data |
| 2026-03-17 | Write package output census: 69 trivial/simple, 38 consolidation candidates, 18 complex |
| 2026-03-16 | Docstring tasks require reading CONVENTIONS.md Documentation section first |
| 2026-03-16 | Convention enforcement needs mechanical verification, not behavioral repetition |
| 2026-03-16 | One-liner method wrappers hide dependencies without adding value |
| 2026-03-16 | Agents reliably introduce gofmt issues during bulk renames |
| 2026-03-15 | Contributor PRs need post-merge follow-up commits for convention alignment |
| 2026-03-15 | Grep for callers must cover entire working tree before deleting functions |
| 2026-03-14 | Stderr error messages are user-facing text that belongs in assets |
| 2026-03-14 | Hardcoded _alt suffixes create implicit language favoritism |
| 2026-03-13 | sync-why mechanism existed but was not wired to build |
| 2026-03-12 | Project-root files vs context files are distinct categories |
| 2026-03-12 | Constants belong in their domain package not in god objects |
| 2026-03-07 | Always search for existing constants before adding new ones |
| 2026-03-07 | SafeReadFile requires split base+filename paths |
| 2026-03-06 | Stale directory inodes cause invisible files over SSH |
| 2026-03-06 | Stats sort uses string comparison on RFC3339 timestamps with mixed timezones |
| 2026-03-06 | Claude Code supports PreCompact and SessionStart hooks that ctx does not use |
| 2026-03-06 | Package-local err.go files invite broken windows from future agents |
| 2026-03-05 | State directory accumulates silently without auto-prune |
| 2026-03-05 | Global tombstones suppress hooks across all sessions |
| 2026-03-05 | Claude Code has two separate memory systems behind feature flags |
| 2026-03-05 | Blog post editorial feedback is higher-leverage than drafting |
| 2026-03-04 | CONSTITUTION hook compliance is non-negotiable — don't work around it |
| 2026-03-02 | Hook message registry test enforces exhaustive coverage of embedded templates |
| 2026-03-02 | Existing Projects is ambiguous framing for migration notes |
| 2026-03-02 | Claude Code JSONL model ID does not distinguish 200k from 1M context |
| 2026-03-01 | Gosec G306 flags test file WriteFile with 0644 permissions |
| 2026-03-01 | Converting PersistentPreRun to PersistentPreRunE changes exit behavior |
| 2026-03-01 | Test HOME isolation is required for user-level path functions |
| 2026-03-01 | Task descriptions can be stale in reverse — implementation done but task not marked complete |
| 2026-03-01 | Model-to-window mapping requires ordered prefix matching |
| 2026-03-01 | TASKS.md template checkbox syntax inside HTML comments is parsed by RegExTaskMultiline |
| 2026-03-01 | Hook logs had no rotation; event log already did |
| 2026-02-28 | ctx pad import, ctx pad export, and ctx system resources make three hack scripts redundant |
| 2026-02-28 | Getting-started docs assumed Claude Code as the only agent |
| 2026-02-28 | Plugin reload script must rebuild cache, not just delete it |
| 2026-02-27 | site/ directory must be committed with docs changes |
| 2026-02-27 | Doctor token_budget vs context_window confusion |
| 2026-02-27 | Drift detector false positives on illustrative code examples |
| 2026-02-27 | Context injection and compliance strategy (consolidated) |
| 2026-02-26 | Webhook silence after ctxrc profile swap is the most common notify debugging red herring |
| 2026-02-26 | Documentation drift and auditing (consolidated) |
| 2026-02-26 | Agent context loading and task routing (consolidated) |
| 2026-02-26 | Go testing patterns (consolidated) |
| 2026-02-26 | PATH and binary handling (consolidated) |
| 2026-02-26 | Task management and exit criteria (consolidated) |
| 2026-02-26 | Agent behavioral patterns (consolidated) |
| 2026-02-26 | Hook compliance and output routing (consolidated) |
| 2026-02-26 | ctx add and decision recording (consolidated) |
| 2026-02-24 | CLI tools don't benefit from in-memory caching of context files |
| 2026-02-22 | Hook behavior and patterns (consolidated) |
| 2026-02-22 | UserPromptSubmit hook output channels (consolidated) |
| 2026-02-22 | Linting and static analysis (consolidated) |
| 2026-02-22 | Permission and settings drift (consolidated) |
| 2026-02-22 | Gitignore and filesystem hygiene (consolidated) |
| 2026-01-28 | IDE is already the UI |
| 2026-04-29 | BunShell ctx.$ calls echo stdout to OpenCode's process unless .quiet() is set — leaks visible noise |
| 2026-04-29 | OpenCode plugin compaction interop is breadcrumb-mediated: own your context preservation explicitly |
| 2026-04-29 | @opencode-ai/plugin event hook is a single dispatcher, not an object of named handlers |
| 2026-04-29 | OpenCode plugin hooks like shell.env take (input, output) and mutate; returned objects are ignored |
| 2026-04-29 | OpenCode shell.env injects env only into agent's shell tool, not into plugin's own ctx.$ calls |
| 2026-04-26 | OpenCode auto-loads only flat .ts files under .opencode/plugins/; subdirectories are ignored |
| 2026-04-26 | OpenCode opencode.json MCP shape: command is Array<string>, no separate args field |
| 2026-04-26 | make test exit code unreliable due to -cover covdata tooling issue |
| 2026-04-26 | Trailing word boundary in regex matches commit-tree as git commit |
| 2026-04-26 | ctx system help can list project-local hooks not in the Go binary |
| 2026-04-25 | Confident code comments can pull an LLM away from first-principles knowledge |
| 2026-04-25 | filepath.Join('', rel) returns rel as CWD-relative, not error |
| 2026-04-25 | Parallel go test ./... packages can race on ~/.claude/settings.json |
<!-- INDEX:END -->

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

## [2026-05-22-161720] Cross-language coverage gap: TS-typed integrations are a fourth surface beyond Go

**Context**: specs/cwd-anchored-context.md removed the CTX_DIR env channel. Three Go test suites caught orphan refs after deletion: audit/TestNoDeadExports (dead consts), audit/TestFlagYAMLMatchesConstants + TestExamplesYAMLLinkage + TestDescKeyYAMLLinkage (orphan YAML keys), compliance/TestDocGoSubcommandDrift (stale doc.go prose). Jumbo commit fc7db228 landed with all four green. But internal/assets/integrations/opencode/plugin/index.ts is a SEPARATE FOURTH surface — TypeScript, not Go — that local 'make lint' and 'go test ./...' never exercise. CI's tsc --noEmit (driven by tools/typecheck/opencode/) surfaced TS2339 on 'output.cwd does not exist on @opencode-ai/plugin shell.env output type'. Fix landed in 40d024a3 but cost a CI round-trip.

**Lesson**: When removing or renaming an env channel, feature flag, or any cross-language contract, the cleanup checklist is FOUR surfaces, not three: (1) Go code (build + lint + test), (2) audit/compliance tests (orphan consts, YAML keys, doc.go drift), (3) asset templates (CLAUDE.md, AGENT_PLAYBOOK, hooks.json, INSTRUCTIONS.md), (4) TypeScript-typed integrations — opencode plugin and the vscode extension. The TS surface is invisible to Go's test suite by design; the typecheck only runs in CI unless invoked explicitly from tools/typecheck/opencode/ or editors/vscode/.

**Application**: Before committing any change that touches internal/assets/integrations/opencode/plugin/ or editors/vscode/, run 'cd tools/typecheck/opencode && npx tsc --noEmit' (and the vscode equivalent). Longer-term: add a 'make typecheck' target wrapping both tsc invocations and include it in the pre-commit checklist alongside 'make lint' and 'go test ./...'. Add it to docs/operations/runbooks/release-checklist.md as a release gate too.

---

## [2026-05-21-140230] Sentinel-removal refactors cascade through test surface

**Context**: Spec specs/cwd-anchored-context.md decomposed the work into 5 discrete steps; in practice steps 1 and 2 had to merge. Removing ErrDirNotDeclared from rc.ContextDir cascaded through ~10 errors.Is consumers and ~30 test fixtures that used t.Setenv(env.CtxDir, ...).

**Lesson**: Spec-level decomposition that treats 'swap resolver' and 'remove init guard' as separable does not survive contact when the second step references the soon-to-be-deleted sentinel from the first. Both have to compile against the new sentinel set in the same commit.

**Application**: When a future spec proposes step boundaries that hinge on a sentinel rename or removal, plan the merged commit up front rather than discover the cascade mid-implementation. The compile-surface analysis belongs at spec time, not implementation time.

---

## [2026-05-20-214839] macOS /var symlink trips path-equality; use EvalSymlinks with parent-resolution fallback

**Context**: TestRunInit_EnvCwdMatch_Succeeds in internal/cli/initialize/init_test.go failed on first run despite a deliberate setup where the env path and cwd candidate matched. Diagnosis: t.TempDir() returns paths like /var/folders/..., os.Getwd() after t.Chdir() returns the canonical /private/var/folders/... (because macOS's /var is a symlink to /private/var). filepath.Clean preserves the symlink form; equality fails.

**Lesson**: filepath.Clean alone is insufficient for path equality on macOS (and other systems with symlinked top-level dirs). filepath.EvalSymlinks resolves the symlinks but fails when the target path does not yet exist — common case for /Users/volkan/Desktop/WORKSPACE/ctx/.context BEFORE ctx init runs. The right pattern is a layered fallback: try EvalSymlinks(full), then EvalSymlinks(parent) + rejoin basename, then filepath.Clean as last resort.

**Application**: Encapsulated as internal/cli/initialize/core/envmatch/{envmatch.go,internal.go}. The Same(a, b) public function calls resolve() on each side; resolve() tries EvalSymlinks on the full path, falls back to EvalSymlinks on the parent (rejoining the basename), and falls through to filepath.Clean if both fail. Reusable for any future env-vs-cwd-style equality check. The package is per-feature (core/envmatch/) per the cmd/core/ purity rule enforced by internal/compliance/TestCmdDirPurity.

---

## [2026-05-20-214830] Handover filenames are archaeology; parse by generated-at, not filename

**Context**: User observed three coexisting handover filename shapes: .context/HANDOVER-2026-04-22.md (pre-skill root file), .context/handovers/YYYY-MM-DD-HHMMSS-slug.md (skill-era pre-CLI), .context/handovers/<RFC3339Compact>-slug.md (current CLI). User asked whether this was a regression or a skill-interpretation problem.

**Lesson**: Neither. The .context/HANDOVER-* root file predates the handovers/ directory contract entirely (the body even said 'delete this file after reading'). The YYYY-MM-DD-HHMMSS shape was an earlier skill iteration writing free-form before ctx handover write existed (commit 60543e46, 2026-05-17, introduced the CLI as sole writer per the anti-pattern note in /ctx-handover SKILL.md). The current parser at internal/write/handover/parse.go:75-107 keys on the 'generated-at' YAML frontmatter, not the filename — so legacy shapes still sort correctly via LatestHandoverCursor. Only files without frontmatter (the root April file) are invisible.

**Application**: When unifying filename shapes across history, use git mv to preserve rename detection. Derive the canonical timestamp from the file's own generated-at frontmatter rather than from the filename — that's the source of truth the parser uses anyway. If a handover predates frontmatter entirely (rare, pre-skill era), it's safe to delete because the parser never read it.

---

## [2026-05-20-214821] /ctx-plan is named after its input, not its output

**Context**: Agent (and apparently other agents in prior sessions per user observation) repeatedly inverted the canonical chain, treating /ctx-spec as the entry point and /ctx-plan as a post-spec step. The skill description starts 'stress-test a plan' (implying user brings a plan IN) while line 44 of the body says 'the deliverable is a debated brief, not a task list' (the OUTPUT is a brief, not a plan).

**Lesson**: Skill names that reference their INPUT bias the agent toward the wrong canonical position. The /ctx-plan skill takes a plan and produces a brief; the natural mental model when scanning the name is 'plan = output', which makes the agent place it AFTER spec instead of before. Also: /ctx-spec's 'When to Use' section listed /ctx-brainstorm as a predecessor but never /ctx-plan, so an agent skimming the top of the skill never learned the full chain.

**Application**: Made the canonical chain explicit at the top of both /ctx-plan and /ctx-spec skills (Canonical Chain block with the brainstorm → plan → spec → implement diagram) and in AGENT_PLAYBOOK_GATE Planning Work section. /ctx-spec When-to-Use now lists /ctx-plan as a predecessor; When-NOT-to-Use says 'when the bet is contested but not yet stress-tested, use /ctx-plan first'. /ctx-plan description now ends with '; produces a debated brief at .context/briefs/<TS>-<slug>.md that /ctx-spec --brief consumes'.

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

## [2026-04-08-074612] fmt.Fprintf to strings.Builder silently discards errors

**Context**: golangci-lint errcheck allows fmt.Fprintf to strings.Builder
because Write never fails, but project convention says zero silent discard

**Lesson**: Linter coverage gaps exist where language guarantees mask
conventions. AST tests fill the gap

**Application**: Created TestNoUncheckedFmtWrite to enforce fmt.Fprintf error
handling. Use if _, err := fmt.Fprintf(...) with log.Warn on the error path

---

## [2026-04-08-074604] AST audit tests must cover unexported functions too

**Context**: TestDocCommentStructure only checked exported functions, so
agent-written helpers in format.go had no godoc enforcement

**Lesson**: Convention enforcement tests must default to scanning all documented
functions. Use explicit opt-outs (test files) not opt-ins (exported only)

**Application**: When adding AST audit tests, scan all functions. We fixed
TestDocCommentStructure to drop the IsExported gate and fixed 84 violations

---

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

## [2026-04-04-025813] Format-verb strings are localizable text, not exempt from magic string checks

**Context**: Strings like '%d entries checked' were passing TestNoMagicStrings
because the format-verb exemption was too broad

**Lesson**: Any string containing English words alongside format directives is
user-facing text that belongs in YAML assets

**Application**: Removed format-verb, URL-scheme, HTML-entity, and err/
exemptions from TestNoMagicStrings

---

## [2026-04-04-025805] Agents add allowlist entries to make tests pass — guard every exemption

**Context**: Found that every exemption map/allowlist in audit tests is a
tempting shortcut for agents

**Lesson**: Added DO NOT widen guard comments to all 10 exemption data
structures across 7 test files

**Application**: Every new audit test with an exemption must include the guard
comment. Review PRs for drive-by allowlist additions.

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

## [2026-04-01-233250] Raw I/O migration unlocks downstream checks for free

**Context**: TestNoRawPermissions had zero violations because the raw I/O
migration moved all octal literals into internal/io/ which already used
config/fs constants

**Lesson**: Chokepoint migrations have cascading benefits — centralizing one
concern (file I/O) automatically resolves other drift (raw permissions)

**Application**: Prioritize chokepoint migrations (io, exec, write, err) before
smaller checks that depend on them

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

## [2026-04-01-074419] Copilot CLI skills need a sync mechanism to prevent drift from ctx skills

**Context**: 5 Copilot CLI skills were condensed versions of ctx skills,
independently maintained with no drift detection

**Lesson**: Any time the same content exists in two locations without a sync
mechanism, it will drift silently

**Application**: make sync-copilot-skills added to build deps, make
check-copilot-skills added to audit target

---

## [2026-04-01-074418] Contributor PRs based on older code reintroduce removed features

**Context**: PR #45 brought back prompt templates, PROMPT.md, and
IMPLEMENTATION_PLAN.md that were explicitly removed in March

**Lesson**: When resolving contributor merge conflicts, check decisions history
for intentional removals — do not assume the PR content is additive

**Application**: Cross-reference DECISIONS.md before accepting PR content that
adds files or features

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

## [2026-03-30-003734] Python-generated doc.go files need gofmt — formatter strips bare // padding lines

**Context**: Batch-generated doc.go files used blank // lines for padding, which
gofmt removes as unnecessary whitespace

**Lesson**: Programmatic Go file generation must produce substantive content
lines, not blank comment padding — gofmt enforces this

**Application**: Always run gofmt after any scripted Go file generation

---

## [2026-03-30-003707] lint-docstrings.sh greedy sed hid all return-type violations

**Context**: sed 's/.*) //' consumed return type parens, leaving { — functions
with return types were invisible to the script for months

**Lesson**: Greedy regex in shell scripts can silently suppress entire
categories of lint violations — test with edge cases, not just happy paths

**Application**: When writing sed-based lint checks, test with multi-paren
signatures (func Foo() (string, error))

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

## [2026-03-25-173338] Template improvements don't propagate to existing projects

**Context**: 5 of 8 context files in the ctx project itself had stale/missing
comment headers — templates evolved but non-destructive init never re-synced
them

**Lesson**: Any template change is invisible to existing users until they run
ctx init --force

**Application**: Added drift detection (checkTemplateHeaders) to ctx drift.
Consider surfacing this during ctx status too.

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

## [2026-03-24-000959] git describe --tags follows ancestry, not global tag list

**Context**: Release notes skill diffed against v0.3.0 instead of v0.6.0 because
the release branch diverged before v0.6.0 was tagged

**Lesson**: git describe --tags --abbrev=0 follows reachability from HEAD; use
git tag --sort=-v:refname | head -1 for the latest tag globally

**Application**: Any script or skill that needs the latest release should use
sorted tag list, not describe

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

## [2026-03-23-003353] Higher-order callbacks in param structs are a code smell

**Context**: MergeParams.UpdateFn and DeployParams.ListErr/ReadErr were function
pointers where all callers passed thin wrappers varying only by a text key

**Lesson**: If all callers pass thin wrappers around the same pattern
(fmt.Errorf with different keys), the callback is just data in disguise

**Application**: When a struct field is a function pointer, check if all callers
vary only by a string key — if so, replace the callback with the key and let
the consumer do the dispatch

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

## [2026-03-16-022650] One-liner method wrappers hide dependencies without adding value

**Context**: checkBoundary() and loadContext() were methods on Handler that just
called validation.ValidateBoundary and context.Load with h.ContextDir

**Lesson**: If a method only passes a struct field to a stdlib function, inline
it — the wrapper obscures the real dependency

**Application**: Before extracting a helper method, check if it just forwards a
field to another function. If so, call the function directly.

---

## [2026-03-16-022642] Agents reliably introduce gofmt issues during bulk renames

**Context**: Subagents renamed consequences->consequence across 75+ files but
left formatting errors in 12 Go files

**Lesson**: Always run gofmt -l after agent-driven refactors before trusting the
build

**Application**: Add gofmt -w pass as a standard step after any agent-driven
bulk edit

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

## [2026-03-15-040642] Grep for callers must cover entire working tree before deleting functions

**Context**: Deleted 7 err/prompt functions as dead code, but callers existed in
unstaged refactoring files — caused build failures

**Lesson**: When the working tree has unstaged changes from a prior session,
grep hits only committed+staged code; must grep the full tree or build-test
before declaring functions dead

**Application**: Always run make build after deleting functions, even if grep
shows zero callers

---

## [2026-03-14-180903] Stderr error messages are user-facing text that belongs in assets

**Context**: Added fmt.Fprintf(os.Stderr) error reporting to event log,
initially with inline strings

**Lesson**: Any string that reaches the user, including stderr warnings, routes
through assets.TextDesc() for i18n readiness

**Application**: When adding stderr output, create text.yaml entries and asset
keys first

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

## [2026-03-13-151952] sync-why mechanism existed but was not wired to build

**Context**: assets/why/ had drifted from docs/ — the sync targets existed in
the Makefile but build did not depend on sync-why

**Lesson**: Freshness checks that are not in the critical path will be
forgotten. Wire them as build prerequisites, not optional audit steps

**Application**: Any derived or copied asset should be a prerequisite of build,
not just audit

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

## [2026-03-06-141506] Stale directory inodes cause invisible files over SSH

**Context**: Files created by Claude Code hooks were visible inside the VM but
not from the SSH terminal

**Lesson**: If a directory is recreated (e.g. by auto-prune), an SSH shell
holding the old directory inode will not see new files — ls returns no such
file even though cat with the full path works from other shells

**Application**: After ctx system prune or any state directory recreation, SSH
sessions need cd-dot or re-login to pick up the new inode

---

## [2026-03-06-141504] Stats sort uses string comparison on RFC3339 timestamps with mixed timezones

**Context**: ctx system stats showed only old sessions, hiding the current one

**Lesson**: RFC3339 string comparison breaks when entries mix UTC (Z) and offset
(-08:00) formats — 13:00-08:00 sorts before 18:00Z lexicographically despite
being later in absolute time

**Application**: Always parse to time.Time before comparing RFC3339 timestamps;
never rely on lexicographic sort

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

## [2026-03-01-222739] Gosec G306 flags test file WriteFile with 0644 permissions

**Context**: New tests used os.WriteFile(..., 0o644) for temp context files;
lint flagged all three occurrences

**Lesson**: Gosec enforces 0600 max on WriteFile even in test code. Use 0o600
for test temp files

**Application**: Default to 0o600 for os.WriteFile in tests; only use wider
permissions when testing permission behavior specifically

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

## [2026-03-01-161459] Test HOME isolation is required for user-level path functions

**Context**: After adding ~/.ctx/.ctx.key as global key location, test suites
wrote real files to the developer home directory

**Lesson**: Any code that uses os.UserHomeDir() needs t.Setenv(HOME, tmpDir) in
tests — especially test helpers called by many tests (like setupEncrypted and
helper)

**Application**: When adding features that write to user-level paths (~/.ctx/,
~/.config/), always add HOME isolation to test setup functions first

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

## [2026-03-01-092611] Hook logs had no rotation; event log already did

**Context**: Investigated .context/logs/ and .context/state/ file management

**Lesson**: eventlog already rotates at 1MB with one previous generation.
logMessage() in state.go was pure append-only with no size check.

**Application**: When adding new log sinks, follow the established rotation
pattern (size-based, single previous generation)

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

## [2026-04-26-152850] make test exit code unreliable due to -cover covdata tooling issue

**Context**: make test exited 1 even with all 123 packages passing on this Go
install; root cause is missing covdata tool when -cover is enabled

**Lesson**: Don't trust make test exit code alone when verifying changes. The
-cover flag in the test target can fail with 'no such tool covdata' even when
every package passes.

**Application**: When make test fails, fall back to 'go test ./...' (no -cover)
and tally ^ok / ^FAIL counts to distinguish real failures from tooling issues.

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

## [2026-04-25-014704] Parallel go test ./... packages can race on ~/.claude/settings.json

**Context**: make test runs packages in parallel processes. Fourteen test files
invoked initialize.Cmd().Execute(), which read-modify-writes
~/.claude/settings.json without HOME isolation.

**Lesson**: Under load the races materialized as flaky 'FAIL coverage: [no
statements]' in cli/watch/core. Run alone the package passed; under parallel
make test it failed intermittently.

**Application**: testctx.Declare now sets HOME alongside CTX_DIR. Centralized
fix; future tests automatically isolate user-home writes.
