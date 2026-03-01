# Archived Learnings (consolidated 2026-02-26)

Originals replaced by consolidated entries in LEARNINGS.md.

## Group: Documentation drift and auditing

## [2026-02-24-204548] CLI reference docs can outpace implementation — always verify against Cobra registration

**Context**: Found 3 commands fully documented in cli-reference.md (ctx remind, ctx recall sync, key file naming) that don't match the binary. Documentation was written speculatively before Cobra subcommands were registered.

**Lesson**: ctx remind has no CLI at all, ctx recall sync has Go code but no Cobra wiring, and key file naming diverged between docs (.ctx.key) and code (.scratchpad.key). Docs can describe commands that are unreachable.

**Application**: Before releasing docs for new commands, verify with ctx <cmd> --help that the command is actually reachable. Add a drift check to the QA gate.

---

## [2026-02-24-171233] Drift-check comments prevent documentation staleness

**Context**: Contributing.md project layout was stale: missing 2 internal packages (notify, sysinfo), 2 top-level dirs (assets, examples), wrong skill count (27 vs 29)

**Lesson**: Structural documentation sections (project layouts, command tables, skill counts) drift silently after code changes. HTML comment markers like drift-check give agents a verification command to run

**Application**: Add drift-check markers above any doc section that mirrors codebase structure. Format: <!-- drift-check: <shell command> -->

---

## [2026-02-24-171231] Documentation style audits require multiple targeted passes

**Context**: Initial agent sweep for filename backticks found only 8 violations; manual re-check found 48+. Same pattern repeated for parenthetical emphasis and quoted terms.

**Lesson**: Automated/agent searches for style violations are unreliable for prose rules with many exception categories (code blocks, table cells, admonitions). Multiple verification passes with direct grep patterns catch what agents miss.

**Application**: When auditing prose style rules, always follow agent results with a targeted grep and manual classification of the results

---

## [2026-02-24-015941] ARCHITECTURE.md had significant drift — 4 core packages and 4 CLI commands missing

**Context**: During the first /ctx-map run, analysis revealed crypto, sysinfo, notify, journal/state were missing from the core packages table, and notify, pad, permissions, system CLI commands were absent. The doc claimed 19 commands (actually 22) and 16 skill templates (actually 28).

**Lesson**: ARCHITECTURE.md drifts silently when new packages are added without updating the doc. The existing /ctx-drift skill catches stale paths but not missing packages — it cannot detect what is absent from a table.

**Application**: Run /ctx-map after adding new packages or CLI commands. The tracking file staleness detection catches modules with new commits, but new modules need the first-run survey to be discovered.

---

## [2026-01-29-164322] Documentation audits require verification against actual standards

**Context**: Agent claimed 'no Go docstring issues found' but manual inspection revealed many functions missing Parameters/Returns sections. The agent only checked if comments existed, not if they followed the standard format.

**Lesson**: When auditing documentation, compare against a known-good example first. Pattern-match for the COMPLETE standard (e.g., '// Parameters:' AND '// Returns:' sections), not just presence of any comment.

**Application**: Before declaring 'no issues', manually verify at least 5 random samples match the documented standard. Use grep patterns that detect missing sections, not just missing comments.

---

## [2026-02-15-231022] Dead link checking is consolidation check 12, not a standalone concern

**Context**: User identified dead links in rendered site as a problem. Initial instinct was a standalone task or drift extension.

**Lesson**: Doc link rot is code-level drift — same category as magic strings or stale architecture diagrams. It belongs in /consolidate's check list, with a standalone /check-links skill that consolidate invokes.

**Application**: When a new audit concern emerges, check if it fits an existing audit skill before creating an isolated one. Consolidate is the natural home for anything that drifts silently between sessions.

---

## Group: ctx init and CLAUDE.md behavior

## [2026-02-14-164029] ctx init is non-destructive toward tool-specific configs

**Context**: Verified by reading run.go — no code paths touch .cursorrules, .aider.conf.yml, or copilot instructions

**Lesson**: ctx init only creates .context/, CLAUDE.md, .claude/, PROMPT.md, and IMPLEMENTATION_PLAN.md. It has zero awareness of other tools' config files.

**Application**: State this definitively in docs rather than hedging — it's confirmed by the code

---

## [2026-02-14-164013] merge insertion is position-aware, not append

**Context**: Reading fs.go findInsertionPoint() to document --merge behavior

**Lesson**: The --merge flag finds the first H1 heading, skips trailing blank lines, and inserts the ctx block there. If no H1 is found, it inserts at the top. Content is never appended to the end.

**Application**: Document the insertion position clearly — users care about where their content ends up in the merged file

---

## [2026-02-14-164011] ctx init CLAUDE.md handling is a 3-state machine

**Context**: Reading claude.go to write the migration guide

**Lesson**: ctx init checks for: no file (create), file without ctx markers (merge/prompt), file with markers (skip or force-replace). The markers <!-- ctx:context --> / <!-- ctx:end --> are the pivot.

**Application**: When documenting merge behavior, describe all three states explicitly rather than just the happy path

---

## [2026-01-20-080000] Always Backup Before Modifying User Files

**Context**: `ctx init` needs to create/modify CLAUDE.md, but user may have existing customizations.

**Lesson**: When modifying user files (especially config files like CLAUDE.md):
1. **Always backup first** — `file.bak` before any modification
2. **Check for existing content** — use marker comments for idempotency
3. **Offer merge, don't overwrite** — respect user's customizations
4. **Provide escape hatch** — `--merge` flag for automation, manual merge for control

**Application**: Any `ctx` command that modifies user files should follow this pattern.

---

## Group: Agent context loading and task routing

## [2026-01-23-180000] ctx agent vs Manual File Reading Trade-offs

**Context**: User asked "Do you remember?" and agent used parallel file reads
instead of `ctx agent`. Compared outputs to understand the delta.

**Lesson**: `ctx agent` is optimized for task execution:
- Filters to pending tasks only
- Surfaces constitution rules inline
- Provides prioritized read order
- Token-budget aware

Manual file reading is better for exploratory/memory questions:
- Session history access
- Timestamps ("modified 8 min ago")
- Completed task context
- Parallel reads for speed

**Application**: No need to mandate one approach. Agents naturally pick appropriately:
- "Do you remember?" → parallel file reads (need history)
- "What should I work on?" → `ctx agent` (need tasks)

---

## [2026-01-23-140000] Infer Intent on "Do You Remember?" Questions

**Context**: User asked "Do you remember?" at session start. Agent asked for
clarification instead of proactively checking context files.

**Lesson**: In a ctx-enabled project, "do you remember?" has an obvious
meaning: check the `.context/` files and report what you know from previous
sessions. Don't ask for clarification - just do it.

**Application**: When user asks memory-related questions ("do you remember?",
"what were we working on?", "where did we leave off?"), immediately:
1. Read `.context/TASKS.md`, `DECISIONS.md`, `LEARNINGS.md`
2. Run `ctx recall list --limit 5` for recent session history
3. Summarize what you find

Don't ask "would you like me to check the context files?" - that's the
obvious intent.

---

## [2026-01-20-180000] .context/ Is NOT a Claude Code Primitive

**Context**: User asked if Claude Code natively understands `.context/`.

**Lesson**: Claude Code only natively reads:
- `CLAUDE.md` (auto-loaded at session start)
- `.claude/settings.json` (hooks and permissions)

The `.context/` directory is a ctx convention. Claude won't know about it unless:
1. A hook runs `ctx agent` to inject context
2. CLAUDE.md explicitly instructs reading `.context/`

**Application**: Always create CLAUDE.md as the bootstrap entry point.

---

## [2026-01-21-160000] Orchestrator vs Agent Tasks Must Be Separate

**Context**: Ralph Loop checked `IMPLEMENTATION_PLAN.md`, found all tasks
done, exited — ignoring `.context/TASKS.md`.

**Lesson**: Separate concerns:
- **`IMPLEMENTATION_PLAN.md`** = Orchestrator directive ("check your tasks")
- **`.context/TASKS.md`** = Agent's mind (actual task list)

The orchestrator shouldn't maintain a parallel ledger. It just says
"check your mind."

**Application**: For new projects, `IMPLEMENTATION_PLAN.md` has ONE task:
"Check `.context/TASKS.md`"

---

## [2026-01-25-170000] AGENTS.md Is Not Auto-Loaded

**Context**: Had both AGENTS.md and CLAUDE.md in project root, causing confusion.

**Lesson**: Only CLAUDE.md is read automatically by Claude Code. Projects
using ctx should rely on the CLAUDE.md → AGENT_PLAYBOOK.md chain, not AGENTS.md.

---

## Group: Blog and content publishing

## [2026-02-18-071508] Blog posts are living documents

**Context**: Session spent enriching two blog posts with cross-links, update admonitions, citations, and contextual notes. Every post had 3-6 places where a link or admonition improved reader experience.

**Lesson**: Blog posts benefit from periodic enrichment passes: cross-linking to newer content, adding update admonitions for superseded features, citing sources, and adding contextual admonitions that connect ideas across posts.

**Application**: Schedule blog enrichment as part of consolidation sessions. When a new feature supersedes something described in a blog post, add an update admonition immediately rather than waiting.

---

## [2026-02-17] Blog publishing from ideas/ requires a consistent checklist

**Context**: Published 4 blog posts from ideas/ drafts in one session. Each required the same steps: date update, path fixes, cross-links, Arc section, blog index, See also in companions. Missing any step left broken links or orphaned posts.

**Lesson**: Blog publishing is a repeatable workflow with 7 steps: (1) update date and frontmatter, (2) fix relative paths from ideas/ to docs/blog/, (3) add cross-links to/from companion posts, (4) add "The Arc" section connecting to the series narrative, (5) update blog index, (6) add "See also" in related posts, (7) verify all link targets exist.

**Application**: Follow this checklist for every ideas/ → docs/blog/ promotion. Consider making it a recipe in hack/runbooks/ if the pattern continues.

---

## [2026-01-28-072838] Changelogs vs Blogs serve different audiences

**Context**: Synthesizing session history into documentation

**Lesson**: Changelogs document WHAT; blogs explain WHY. Same information,
different engagement. Changelogs are for machines (audits, dependency trackers).
Blogs are for humans (narrative, context, lessons).

**Application**: When synthesizing session history, output both: changelog for
completeness, blog for readability.

---

## [2026-02-15-040313] Cross-repo links to published docs should use ctx.ist

**Context**: hack/runbooks/persistent-irc.md linked to docs/ via relative paths, getting-started.md linked to MANIFESTO.md via GitHub — both bypass ctx.ist rendering (admonitions, nav, search)

**Lesson**: When content is published on ctx.ist, always link to the site URL, not the GitHub blob or a relative file path. GitHub won't render zensical admonitions and readers lose navigation context.

**Application**: When adding See Also or cross-references in runbooks or docs, use https://ctx.ist/... URLs for anything the site publishes. Reserve GitHub links for repo-only content (issues, releases, security tab, source files not on the site).

---

## Group: Worktrees and parallel agents

## [2026-02-21-200039] Parallel agents work cleanly on disjoint file sets

**Context**: The .contextrc → .ctxrc rename touched ~40 files across Go source, docs, specs, and context files

**Lesson**: Splitting work across 4 parallel agents by file-type (Go source, new files, docs, specs/context) with zero file overlap completed the entire rename in one pass with no conflicts

**Application**: For broad renames or refactors, partition by file domain and run agents in parallel rather than sequentially

---

## [2026-02-24-030252] Worktree agents lack key-dependent features by design

**Context**: Investigated whether hooks/notify/pad would break in git worktrees when .ctx.key is gitignored

**Lesson**: ctx pad fails gracefully (no key), ctx notify silently no-ops, journal enrichment writes to the worktree and is orphaned on teardown. All path resolution is cwd-relative with no git-root awareness. The only real gap is notify — pad and journal are naturally avoided by workflow.

**Application**: Document worktree limitations in skill docs, recipes, and reference pages. File a task to enable notify in worktrees (the one feature that would genuinely help autonomous opaque agents).

---

## [2026-02-21-195820] Multi-agent work requires commit auditing before batching

**Context**: Found 122 uncommitted files from 3+ agent sessions spanning different features (knowledge health, recall safety, system ergonomics)

**Lesson**: When multiple agents work in parallel, their changes intermingle in the working tree with no commit boundaries. An audit pass is needed to group changes logically before committing.

**Application**: After multi-agent sessions, run git diff --stat and group by logical feature before committing. Don't commit everything in one giant batch.

---

## [2026-02-17-183937] rsync between worktrees can clobber permissions and gitignored files

**Context**: Used rsync -av to borrow upstream changes; it overwrote .claude/hooks/*.sh with non-executable copies and clobbered gitignored settings.local.json

**Lesson**: rsync -av preserves source permissions, not destination. Gitignored files have no git safety net. Use --no-perms or --chmod=+x for scripts, and --exclude gitignored paths explicitly.

**Application**: When borrowing between worktrees: 1) exclude gitignored paths (.claude/settings.local.json, ideas/, .context/logs/) 2) restore +x on hook scripts after sync 3) consider the /absorb skill which handles these edge cases

---

## [2026-02-12] Git worktrees for parallel agent development

**Context**: Explored using git worktrees for parallel agent development across a large task backlog

**Lesson**: Git worktrees enable parallel Claude Code agent sessions without file conflicts. Create worktrees OUTSIDE the project as sibling directories (`git worktree add ../ctx-docs -b work/docs`). Each worktree gets its own branch, staging area, and working files but shares the same `.git` object database. Group tasks by blast radius (files touched) to minimize merge conflicts. 3-4 parallel worktrees is the practical limit before merge complexity outweighs productivity gains.

**Application**: When tackling many independent tasks: (1) group by file overlap, (2) create worktrees as siblings with `git worktree add ../name -b work/name`, (3) launch separate claude sessions in each, (4) merge back to main as tracks complete, (5) cleanup with `git worktree remove`. Don't run `ctx init` in worktrees — `.context/` is already tracked.

---

## Group: Go testing patterns

## [2026-02-26-000104] Compiler-driven refactoring misses test files

**Context**: Changed notify.Send signature to force callsite updates. go build caught all production code but test files broke separately

**Lesson**: go build ./... catches production callsite breaks but not test callsite breaks. go test ./... compiles tests too

**Application**: After changing a Go function signature, run go test ./... (not just go build) to catch test file callsites. Or grep for the function name in *_test.go as a belt-and-suspenders check

---

## [2026-02-24-015827] All runCmd() returns must be consumed in tests

**Context**: golangci-lint errcheck flagged every runCmd(...) call in remind_test.go that didn't capture the return

**Lesson**: Even setup calls in tests that run commands as preconditions need '_, _ = runCmd(...)' to satisfy errcheck

**Application**: When writing test helpers that call cobra commands as setup, always capture both returns

---

## [2026-02-14-163552] color.NoColor in init for CLI test files

**Context**: Recall CLI tests had ANSI escape codes in output making string assertions unreliable

**Lesson**: Setting color.NoColor = true in a package-level init function disables ANSI codes for all tests in the package

**Application**: Add init with color.NoColor = true in test files for CLI packages that use fatih/color. Cleaner than per-test setup.

---

## [2026-02-14-163551] Recall CLI tests isolate via HOME env var

**Context**: Needed integration tests for recall list/show/export without touching real session data

**Lesson**: parser.FindSessions reads os.UserHomeDir which uses HOME env var. Setting t.Setenv HOME to tmpDir with .claude/projects/ structure gives full isolation.

**Application**: For recall integration tests: t.Setenv HOME to tmpDir, create .claude/projects/dir/ with JSONL fixtures. See internal/cli/recall/run_test.go.

---

## [2026-02-14-163550] formatDuration accepts interface not time.Duration

**Context**: Writing unit tests for formatDuration in recall/fmt.go

**Lesson**: formatDuration takes interface with Minutes method, not time.Duration directly. A stub type is needed for testing.

**Application**: Use a stubDuration struct with a mins field and Minutes method when testing formatDuration. See internal/cli/recall/fmt_test.go.

---

## [2026-01-25-180000] CI Tests Need CTX_SKIP_PATH_CHECK

**Context**: CI tests were failing because ctx binary isn't installed on CI runners.

**Lesson**: Tests that call `ctx init` will fail without `CTX_SKIP_PATH_CHECK=1`
env var, because init checks if ctx is in PATH.

---

## [2026-01-19-120000] CGO Must Be Disabled for ARM64 Linux

**Context**: `go test` failed with `gcc: error: unrecognized command-line option '-m64'`

**Lesson**: On ARM64 Linux, CGO causes cross-compilation issues. Always use `CGO_ENABLED=0`.

**Application**:
```bash
CGO_ENABLED=0 go build -o ctx ./cmd/ctx
CGO_ENABLED=0 go test ./...
```

---

## Group: PATH and binary handling

## [2026-01-23-120000] Always Use ctx from PATH

**Context**: Agent used `./dist/ctx-linux-arm64` and `go run ./cmd/ctx`
instead of just `ctx`, even though the binary was installed to PATH.

**Lesson**: When working on a ctx-enabled project, always use `ctx` directly:
```bash
ctx status        # correct
ctx agent         # correct
./dist/ctx        # avoid hardcoded paths
go run ./cmd/ctx  # avoid unless developing ctx itself
```

**Application**: Check `which ctx` if unsure. The binary is installed during
setup (`sudo make install` or `sudo cp ./ctx /usr/local/bin/`).

---

## [2026-01-21-120000] Hooks Should Use PATH, Not Hardcoded Paths

**Context**: Original hooks used hardcoded absolute paths like
`/home/user/project/dist/ctx-linux-arm64`. This caused issues when dogfooding
or sharing configs.

**Lesson**: Hooks should assume `ctx` is in the user's PATH:
- More portable across machines/users
- Standard Unix practice
- `ctx init` now checks if `ctx` is in PATH before proceeding
- Hooks use `ctx agent` instead of `/full/path/to/ctx-linux-arm64 agent`

**Application**:
1. Users must install ctx to PATH: `sudo make install` or `sudo cp ./ctx /usr/local/bin/`
2. `ctx init` will fail with clear instructions if ctx is not in PATH
3. Tests can skip this check with `CTX_SKIP_PATH_CHECK=1`

**Supersedes**: Previous learning "Binary Path Must Be Absolute" (2026-01-20)

---

## [2026-02-17] Agent must never place binaries — nudge the user to install

**Context**: Agent removed ~/go/bin/ctx and discussed copying to /usr/local/bin. Both actions bypass the proper installation path (make install with elevated privileges) which the agent cannot run.

**Lesson**: The agent must never place binaries in any bin directory — not via cp, mv, go install, or any other mechanism. When a rebuild is needed, the agent builds with `make build` and asks the user to run the privileged install step themselves.

**Application**: When ctx binary is stale or missing: (1) run `make build`, (2) ask the user to install it (requires privileges), (3) wait for confirmation before continuing. Hooks in block-dangerous-commands.sh now block cp/mv to bin dirs and `go install` as a command.

---

## Group: Task management and exit criteria

## [2026-02-13-133314] Specs get lost without cross-references from TASKS.md

**Context**: Designed encrypted scratchpad feature, wrote spec in specs/scratchpad.md, tasked it out in TASKS.md. Realized a new session picking up the tasks might never find the spec.

**Lesson**: Agents read TASKS.md early but may never discover specs/ on their own. Single-layer instructions get skipped under pressure; redundancy across layers is the only reliable mitigation for probabilistic instruction-following.

**Application**: Three-layer defense for every spec: (1) playbook instruction for the general pattern, (2) spec reference in the Phase header, (3) bold breadcrumb in the first task of the phase. Added 'Planning Non-Trivial Work' section to AGENT_PLAYBOOK.md to codify this.

---

## [2026-01-28-040915] Subtasks complete does not mean parent task complete

**Context**: AI marked parent task done after finishing subtasks but missing
actual deliverable

**Lesson**: Subtask completion is implementation progress, not delivery.
The parent task defines what the user gets.

**Application**: Parent tasks should have explicit deliverables; don't close
until deliverable is verified.

---

## [2026-01-21-180000] Exit Criteria Must Include Verification

**Context**: Dogfooding experiment had another Claude rebuild `ctx` from specs.
All tasks were marked complete, Ralph Loop exited successfully. But the built
binary didn't work — commands just printed help text instead of executing.

**Lesson**: "All tasks checked off" ≠ "Implementation works." This applies to
US too, not just the dogfooding clone. Our own verification is based on manual
testing, not automated proof. Blind spots exist in both projects.

Exit criteria must include:
- **Integration tests**: Binary executes commands correctly (not just unit tests)
- **Coverage targets**: Quantifiable proof that code paths are tested
- **Smoke tests**: Basic "does it run" verification in CI

**Application**:
1. Add integration test suite that invokes the actual binary
2. Set coverage targets (e.g., 70% for core packages)
3. Add verification tasks to TASKS.md — we have the same blind spot
4. Being proud of our achievement doesn't prove its validity

---

## [2026-02-17] Reports graduate to ideas/done/ only after all items are tracked or resolved

**Context**: Moving REPORT-6 and REPORT-7 to ideas/done/. Each had a mix of completed, skipped, and untracked items. Moving before tracking would lose the untracked items.

**Lesson**: Before graduating a report: (1) cross-reference every item against TASKS.md and the codebase, (2) add trackers for undone items, (3) create specs for items that need design, (4) put remaining low-priority items in a future-considerations document, (5) update TASKS.md path references, (6) then move.

**Application**: Always do the full cross-reference before moving reports to done/. The report is the source of truth until every item has a home elsewhere.

---

## Group: Agent behavioral patterns

## [2026-02-22-212342] Interaction pattern capture risks softening agent rigor

**Context**: Considered a skill to analyze session history and extract user interaction patterns into .context/

**Lesson**: Automated pattern capture from sessions risks training the agent to please rather than push back. Scoping to 'process only' is insufficient — the agent doing the analysis has an incentive to learn the wrong lessons. Existing mechanisms (learnings, hooks, constitution) already capture process preferences explicitly and with human review.

**Application**: Do not build implicit user-modeling from session history. Rely on explicit, human-reviewed context (learnings, conventions, hooks) for behavioral shaping. If interaction patterns are ever revisited, require a CONSTITUTION clause preventing compliance-oriented learning.

---

## [2026-02-11-124635] Chain-of-thought prompting improves agent reasoning accuracy

**Context**: Research shows accuracy on reasoning tasks jumps from 17.7% to 78.7% by adding think step-by-step to prompts. Applied this across agent guidelines.

**Lesson**: Explicit think step-by-step instructions in agent prompts dramatically improve reasoning accuracy at negligible token cost. This applies to skill files, playbooks, and autonomous loop prompts — anywhere the agent makes decisions before acting.

**Application**: Added Reason Before Acting section to AGENT_PLAYBOOK.md and reasoning nudges to 7 skills (ctx-implement, brainstorm, ctx-reflect, ctx-loop, qa, verify, consolidate). For autonomous loops, include reasoning instructions in PROMPT.md.

---

## [2026-01-30-120009] Say 'project conventions' not 'idiomatic X'

**Context**: When asking Claude to follow documentation style, saying
'idiomatic Go' triggered training priors (stdlib conventions) instead of
project-specific standards.

**Lesson**: Use 'follow project conventions' or 'check AGENT_PLAYBOOK' rather
than 'idiomatic [language]' to ensure Claude looks at project files first.

**Application**: In prompts requesting style alignment, reference project
files explicitly rather than language-wide conventions.

---

## [2026-01-25-140000] Autonomous Mode Creates Technical Debt

**Context**: Compared commits from autonomous "YOLO mode" (auto-accept,
agent-driven) vs human-guided refactoring sessions.

**Lesson**: YOLO mode is effective for feature velocity but accumulates technical debt:

| YOLO Pattern                           | Human-Guided Fix                      |
|----------------------------------------|---------------------------------------|
| `"TASKS.md"` scattered in 10 files     | `config.FilenameTask` constant        |
| `dir + "/" + file`                     | `filepath.Join(dir, file)`            |
| `{"task": "TASKS.md"}`                 | `{UpdateTypeTask: FilenameTask}`      |
| Monolithic `cli_test.go` (1500+ lines) | Colocated `package/package_test.go`   |
| `package initcmd` in `init/` folder    | `package initialize` in `initialize/` |

**Application**:
1. Schedule periodic consolidation sessions (not just feature sprints)
2. When same literal appears 3+ times, extract to constant
3. Constants should reference constants (self-referential maps)
4. Tests belong next to implementations, not in monoliths

---

## [2026-02-19-215204] Trust the binary output over source code analysis

**Context**: Wrongly concluded ctx decisions archive was missing from the installed binary based on a single CLI test that showed parent help instead of subcommand help. The user's own terminal showed it working fine.

**Lesson**: A single ambiguous CLI output is not proof of absence. Re-run the exact command before claiming something is missing. When the user contradicts your finding, they are probably right.

**Application**: When checking if a subcommand exists, run the subcommand directly (e.g., ctx decisions archive --help) and if results are ambiguous, retry before drawing conclusions.

---

## Group: Hook compliance and output routing

## [2026-02-22-194441] Plain-text hook output is silently ignored by the agent

**Context**: The qa-reminder PreToolUse:Edit hook was firing successfully (confirmed in debug log) but producing no visible effect on agent behavior

**Lesson**: Claude Code parses hook stdout that starts with { as JSON directives; plain text is injected as context but the agent treats it as disposable. Structured JSON with hookSpecificOutput.additionalContext is reliably processed.

**Application**: All non-blocking PreToolUse/PostToolUse hooks should return JSON via printHookContext(), not cmd.Println() plain text

---

## [2026-02-25-180000] Hook compliance degrades on narrow mid-session tasks — relay block is the compliance canary

**Context**: Prior agent sessions skipped context file reads despite PreToolUse hook explicitly requiring them. Analysis with user revealed failure modes.

**Lesson**:
- Prior agents skipped context files when given narrow tasks ("fix the off-by-one in budget.go")
- Root cause: CLAUDE.md's "may or may not be relevant" system reminder competed with hook's "Do not assess relevance. Do not skip."
- Fix: CLAUDE.md now explicitly elevates hook authority over relevance filtering
- Risk: Mid-session narrow tasks still have ~15-25% partial skip rate
- Mitigation: The mandatory checkpoint relay block ensures visibility even on partial skip — if you don't see the block, the files weren't read
- Constitution now includes: context loading is step one of every session, not a detour

**Application**: If you're reading this, the system is working — you loaded context files and encountered your prior self's failure analysis. Don't skip the reads. The 30-second cost is always cheaper than a decision made without context.

---

## [2026-02-22-194443] No reliable agent-side before-session-end event exists

**Context**: Investigated whether SessionEnd hooks could trigger context persistence automatically

**Lesson**: SessionEnd fires after the agent is gone — it can only run fire-and-forget shell commands (cleanup), not LLM reasoning. CTRL+C terminates immediately. The playbook already encodes the correct strategy: persist as you go.

**Application**: Mid-session nudges (check-persistence) and explicit /ctx-wrap-up are the only reliable persistence mechanisms; do not design hooks that assume the agent gets a final turn

---

## Group: ctx add and decision recording

## [2026-02-14-164053] ctx add learning/decision requires structured flags, not just a string

**Context**: Repeatedly suggested bare ctx add learning '...' in session endings despite this being wrong

**Lesson**: Learnings require --context, --lesson, --application. Decisions require --context, --rationale, --consequences. A bare string only sets the title — the command will fail without the required flags.

**Application**: Never suggest ctx add learning 'text' as a one-liner. Always show the full flag form. The CLAUDE.md template and session-end prompts should model the correct syntax.

---

## [2026-01-28-191951] Required flags now enforced for learnings

**Context**: Implemented ctx add learning flags to match decision's ADR
(Architectural Decision Record) pattern

**Lesson**: Structured entries with Context/Lesson/Application are more useful
than one-liners

**Application**: Always use ctx add learning with all three flags; agents
guided via AGENT_PLAYBOOK.md

---

## [2026-01-27-180000] Always Complete Decision Record Sections

**Context**: Decisions added via `ctx add decision` were left with placeholder
text like "[Add context here]".

**Lesson**: When recording decisions, always fill in Context
(what prompted this), Rationale (why this choice over alternatives), and
Consequences (what changes as a result). Placeholder text is a code smell -
decisions without rationale lose their value over time.

**Application**: After using `ctx add decision`, immediately edit the file to
complete all sections. Future: use `--context`, `--rationale`, `--consequences`
flags when available.

---

## [2026-01-27-160000] Slash Commands Require Matching Permissions

**Context**: Claude Code slash commands using `!` bash syntax require matching
permissions in settings.local.json.

**Lesson**: When adding new /ctx-* commands, ensure ctx init pre-seeds the
required `Bash(ctx <subcommand>:*)` permissions. Use additive merging for user
config - never remove existing permissions.

---

## Group: Plugin and marketplace architecture

## [2026-02-16-164547] Duplicate skills appear with namespace prefix in Claude Code

**Context**: Had both .claude/skills/ctx-status and the marketplace plugin providing the same skill

**Lesson**: When a repo-local .claude/skills/ directory and a marketplace plugin both define the same skill name, Claude Code lists both: the local version unprefixed and the plugin version with a ctx: namespace prefix (e.g. ctx-status and ctx:ctx-status)

**Application**: To avoid confusing duplicates, ensure distributed skills live only in the plugin source (internal/assets/claude/skills/) and not also in .claude/skills/. Dev-only skills that aren't in the plugin won't collide.

---

## [2026-02-16-164518] Local marketplace plugin enables live skill editing

**Context**: Setting up the contributor workflow for ctx development

**Lesson**: Claude Code marketplace plugins source from the repo root where `.claude-plugin/marketplace.json` lives (e.g. ~/WORKSPACE/ctx). The marketplace.json points to the actual plugin in `internal/assets/claude`. Edits to skills and hooks under that path take effect on the next Claude Code load — no reinstall needed

**Application**: The contributor docs instruct devs to add their local clone as a marketplace source rather than using the GitHub URL. This gives them live feedback on skill changes without a rebuild cycle.

---

## [2026-02-16-164547] Security docs are most vulnerable to stale paths after architecture migrations

**Context**: Migrated from per-project .claude/hooks/ and .claude/skills/ to plugin model; found 5 security docs still referencing the old paths

**Lesson**: When moving infrastructure from per-project files to a plugin/external model, audit security docs first — stale paths in security guidance give users a false sense of protection (e.g. 'make .claude/hooks/ immutable' for a directory that no longer exists)

**Application**: After any file-layout migration, grep security and agent-security docs for old paths before anything else

---
