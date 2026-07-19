# skills-agents-and-tasks

## [2026-07-18-084406] A /ctx-task-out measurement gate caught a pre-existing bug before it could compound

**Context**: pd-m1's E4 'layout proof' epic was decomposed specifically to prove-or-kill the progressive-disclosure premise that 'ctx add needs zero change' BEFORE building the mover on it. The empty-staging case failed: ctx learning add destroyed the ## Themes section.

**Lesson**: The root cause was a shipped insert.AfterHeader tail-truncation bug (content[:i]+entry, dropping content[i:]), not a design flaw. The gate surfaced a real data-loss defect in ~4 tasks instead of after building on a false premise. Decomposing a milestone's load-bearing assumption into an early falsifiable task pays for itself.

**Application**: When running /ctx-task-out, make the riskiest load-bearing assumption its OWN early task with a falsifiable acceptance check; treat its failure as signal about reality (a real bug), not just about the plan. Log the gate firing in the plan's Amendments.

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

## [2026-07-06-a] `ctx journal site` rewrites source entry bodies in place

**Context**: The self-heal render-hash guard assumed only `execute.go`
(import) authors entry bodies. But `ctx journal site` (`make journal`)
soft-wraps/collapses the *source* `.md` in place (site/run.go), and did
NOT refresh the render hash. Result: after any site build, growth-aware
import saw a hash mismatch and false-flagged every normalized entry as
"edited outside ctx", stranding the growth.

**Lesson**: A "prove the file is unchanged since we wrote it" hash is
only valid if *every* code path that writes the body refreshes it. Grep
all writers of the artifact, not just the obvious one.

**Application**: When adding a content-integrity hash, enumerate every
writer (`grep -rn SafeWriteFile <dir>`) and make each refresh the hash —
or route all writes through one helper. Refresh unconditionally on every
non-failed pass so an idempotent transform still re-syncs.

---

## [2026-07-06-b] zensical bundles a fixed lucide icon snapshot

**Context**: `make site` failed with "template not found:
.icons/lucide/message-square-cog.svg". A doc frontmatter `icon:` field
referenced a lucide icon absent from zensical 0.0.47's bundled set (and
0.0.47 is the latest). Blocks the entire build, not just that page.

**Lesson**: zensical (like mkdocs-material) ships a frozen lucide
snapshot; a valid-on-lucide.dev icon may not exist in the installed
version. A single bad `icon:` fails the whole site build.

**Application**: Before a site rebuild, validate every doc `icon:
lucide/X` against the installed bundle
(`find <zensical>/templates/.icons/lucide -name X.svg`). Pick from the
bundled set, not lucide.dev.

---

## [2026-07-06-c] PATH `ctx` is a stale install; verify against the fresh binary

**Context**: After editing an asset (commands.yaml) and `make build`,
`ctx journal import --help` still showed the OLD text. Cause: `ctx` on
PATH is `/usr/local/bin/ctx` (a prior `make install`), while `make
build` produces `./ctx` in the repo. A project hook blocks running
`./ctx` directly.

**Lesson**: `make build` does not update the PATH binary. Verifying
behavior via PATH `ctx` after a source change tests the stale install,
not your change.

**Application**: To verify a fresh build, copy it out
(`cp ./ctx <scratch>/ctx-new`) and run that — it dodges the `./ctx`
hook and reflects your edits. Or `sudo make reinstall`. Never trust
PATH `ctx` help/behavior right after `make build`.

---

## [2026-07-06-d] `internal/assets` ← `internal/rc` blocks package-assets tests importing rc

**Context**: Moving the `.ctxrc` schema-vs-struct bijection test to
reflect over the real `rc.CtxRC` (drift-proofing it) failed: `rc`
imports `internal/assets/read/placeholders`, which imports top-level
`internal/assets`. So a `package assets` test importing `rc` forms a
cycle (assets → rc → placeholders → assets).

**Lesson**: `internal/assets` is a low-level leaf that `rc` transitively
depends on; tests in `package assets` cannot import `rc`.

**Application**: Home a struct-vs-schema guard in the package that owns
the struct (`package rc`), importing `internal/assets` to read the
embedded schema — that direction is acyclic (assets never imports rc).

---

