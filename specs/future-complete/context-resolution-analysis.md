# Context Directory Resolution: Problem Statement & Observations

**Status:** analysis, not a proposal
**Scope:** how `ctx` resolves `.context/` and `.ctxrc` across CLI, hook, and
sub-agent invocations — and why the current algorithm produces silent
failures in real setups.

---

## TL;DR

`ctx` resolves `.context/` by walking up from CWD with a git-root boundary
check. This is a heuristic that silently picks the wrong project when:

- Projects are legitimately nested (e.g. `WORKSPACE` meta-project containing
  `WORKSPACE/ctx`).
- A sub-agent runs with a CWD that happens to be inside a nested project,
  but the *session's intent* is the outer one.
- Another agent (or anything with write access) left a stray `.context/` dir
  along the walk path.
- Git submodules create a `.git` that isn't the project's real root.

`.ctxrc` has a separate, *weaker* resolution: CWD only, no walk-up. This
creates an asymmetry where a project-root `.ctxrc` becomes invisible to a
hook firing from a subdirectory — silently dropping webhook routes,
diagnostics config, and other behavior-altering directives.

A full fix isn't achievable at the `ctx` process level alone: `ctx` cannot
reliably distinguish who invoked it (parent agent, sub-agent, human),
because the harness (Claude Code) does not propagate identifying env vars
into sub-agent environments (Anthropic issue
[#26429](https://github.com/anthropics/claude-code/issues/26429)).

The realistic ceiling is **observability over enforcement**: make every
write's resolution decision loud and attributable, so fragmentation is
detectable post-hoc instead of invisible at write-time.

---

## Current Behavior

### `.context/` Resolution (`internal/rc/walk.go`)

Priority order, highest first:

1. **CLI override** (`--context-dir`, stored via `OverrideContextDir`,
   `rc.go:74-85`). Absolutized, no walk.
2. **Absolute path from `.ctxrc` or `CTX_DIR` env** (`walk.go:37-39`).
   Returned as-is, no walk.
3. **Walk-up from CWD** looking for a directory whose basename matches the
   configured name (default `.context`) (`walk.go:49-60`).
4. **Git-root or CWD fallback** if the walk finds nothing
   (`walk.go:64-69`).

Key behaviors:

- **Innermost wins.** First matching directory up the tree is chosen.
- **Git-root boundary check.** If the candidate is found *above* CWD, it
  must be within the git root (`walk.go:84-88`), otherwise we distrust the
  ancestor and use the git root as anchor (`walk.go:91-93`).
- **No `.git` anywhere → fall back to `cwd/.context`** (`walk.go:69`).
  `ctx init` creates it; hooks would error out.

### `.ctxrc` Resolution (`internal/rc/load.go`)

Resolution is literal `os.ReadFile(".ctxrc")` against process CWD
(`load.go:30` → `SafeReadUserFile` → `os.ReadFile`). Missing file is
silently ignored; defaults apply. **There is no walk-up for `.ctxrc`.**

This asymmetry is the first structural bug: `.context/` searches up the
tree, `.ctxrc` does not. A hook firing from a subdirectory resolves to the
project-root `.context/` but never sees the project-root `.ctxrc`.

---

## Scenarios Where This Breaks

### Scenario A: Hook from a Subdirectory

Hook fires during an edit to `/repo/src/components/foo.tsx`.

- `.context/` walk-up correctly finds `/repo/.context`.
- `.ctxrc` read from CWD (`/repo/src/components/`) — miss.
- Any directives in `/repo/.ctxrc` (webhook routing, event filters,
  billing warnings, session prefixes, …) are silently ignored.
- Observability hole with no error signal.

### Scenario B: Nested Legitimate Projects

```
~/Desktop/WORKSPACE/          .git + .context    (meta-project)
~/Desktop/WORKSPACE/ctx/      .git + .context    (ctx project)
```

User opens a session at `WORKSPACE` to run arch-exploration across all
sub-projects. Agent `cd`s into `WORKSPACE/ctx/` to investigate.

- Any direct `ctx` call from `WORKSPACE/ctx/` resolves innermost-wins →
  writes go to `WORKSPACE/ctx/.context/`.
- The session's *intent* was `WORKSPACE/.context/` (arch-exploration
  narrative).
- Result: tasks, decisions, learnings, journal entries, and webhook
  events fragment across two projects.
- The filesystem cannot distinguish "legitimate sub-project running
  independently" from "sub-agent working on behalf of outer session."

### Scenario C: Rogue / Stray `.context/` Directories

```
a/.context                       (real project A — different concern)
a/b/.git + a/b/.context          (real project B)
a/b/c/.context                   (left behind by an agent earlier)
a/b/c/d/.git                     (submodule of B)
a/b/c/d/e/.context               (left behind by another agent)
a/b/c/d/e/f/                     (CWD)
```

Walk-up from `f`:

- Finds `e/.context` first. `findGitRoot(f)` returns `d` (first `.git`
  up). Candidate `e/.context` is under `d` → passes boundary check →
  **wins.** Wrong.
- If `e/.context` didn't exist: walk continues to `c/.context` — also
  wrong.
- If `c/.context` didn't exist either: `b/.context` — correct.
- If `b/.context` didn't exist: `a/.context` — wrong again (different
  project).

The git-root boundary check helps only when rogue dirs live outside any
`.git` subtree. Submodules and nested `.git` entries defeat it.

### Scenario D: Sub-Agent Fragmentation

Parent session at `WORKSPACE`. Parent spawns a sub-agent to debug
`WORKSPACE/ctx`.

- Sub-agent inherits PWD (possibly `WORKSPACE/ctx` if the parent `cd`'d).
- Sub-agent's direct `ctx` calls from its Bash tool: no
  `CLAUDE_PROJECT_DIR`, no session id in env. Walk-up picks
  `WORKSPACE/ctx/.context/`.
- Hooks firing on the sub-agent's tool calls: `CLAUDE_PROJECT_DIR` *is*
  set (hook contract), points at the parent session's project —
  `WORKSPACE`.
- Within a single sub-agent, direct calls and hook-triggered calls
  **disagree** about which project they belong to. Sub-sessional
  fragmentation.

### Scenario E: Intentional Shared `.context/`

User wants three projects to share one external `.context/` via
`.ctxrc` indirection (absolute `context_dir` + `allow_outside_cwd: true`).

- Works if `ctx` is invoked from the project root.
- Breaks from any subdirectory: the project-root `.ctxrc` isn't found
  (no walk-up), so the override never takes effect. Falls back to
  default walk → resolves somewhere else entirely.

---

## Empirical Findings (Verified This Session)

1. **`CLAUDE_PROJECT_DIR` is not set in plain Claude Code sessions.**
   `env | grep CLAUDE` in both parent agent and sub-agent returned only
   `CLAUDECODE=1`, `CLAUDE_CODE_ENTRYPOINT=cli`,
   `CLAUDE_CODE_EXECPATH=…`. Documented to be set in *hook* payloads
   only.

2. **`CLAUDE_PROJECT_DIR` is not propagated into sub-agent
   environments.** Per
   [anthropics/claude-code#26429](https://github.com/anthropics/claude-code/issues/26429),
   this is a known harness limitation. The mechanism exists (other
   `CLAUDE_*` vars propagate) but this one doesn't.

3. **Session id arrives via stdin JSON, not env.**
   `internal/entity/hook.go:15-18`: `HookInput.SessionID` is
   unmarshaled from the hook's stdin payload. Outside hooks, there is no
   in-process way to recover it.

4. **`ctx` does not currently reference `CLAUDE_PROJECT_DIR`.** Ripgrep
   across `internal/` returns zero matches. The var is unused by ctx
   today — any design leveraging it is net-new integration.

5. **Sub-agent inherits parent env and PWD by default.** Confirmed via
   sub-agent probe (`CLAUDE_CODE_ENTRYPOINT` survived; `PWD` matched).

6. **`cd` does not persist between Bash calls within a sub-agent.** Each
   Bash tool call starts a fresh shell from the sub-agent's invocation
   PWD. So CWD-based reasoning inside a sub-agent is unstable across
   calls, not just across agents.

---

## Why In-ctx Enforcement Can't Work

`ctx` as a process cannot reliably distinguish among:

- Parent agent invoking `ctx` via Bash tool
- Sub-agent invoking `ctx` via Bash tool
- Human running `ctx` in an adjacent terminal

All three share `CLAUDECODE=1` (if invoked from inside Claude Code),
inherited PATH, similar PPID shapes, comparable CWD semantics. No
kernel-visible signal differentiates them. Every candidate heuristic
either false-positives on the human or false-negatives on some sub-agent
configuration.

Enforcement primitives that *would* work require harness cooperation:

- Claude Code tagging sub-agent invocations with a distinguishing env var.
- Claude Code supporting per-sub-agent command filters that block `ctx`.

Neither exists today.

Consequence: any rule like "sub-agents don't use ctx" is a **norm, not an
invariant**. Useful, worth documenting, but not enforceable by ctx.

---

## Design Ideas Considered

Brief, with why each helps or fails. None are solo solutions.

### `.context/PROJECT` Marker File

`ctx init` writes a marker (self-reference or pointer) inside `.context/`.
Walk-up only adopts marked directories; unmarked `.context/` crashes.

- **Helps:** Rejects stray `.context/` dirs without markers. Catches the
  accidental rogue case. Makes "what is a project root" an explicit
  declaration, not an inference.
- **Fails:** A rogue agent can also create the marker. Self-consistent
  rogue is indistinguishable from legitimate nested sub-project.
- **Verdict:** Necessary but not sufficient. Raises the floor.

### `PROJECT` Contains Enclosing Git-Remote Fingerprint

Marker stores the git remote URL (or a hash of it) of the enclosing repo
at init time. Walk-up rejects markers whose fingerprint doesn't match the
enclosing repo's current remote.

- **Helps:** Detects cross-project adoption (marker from repo X found
  inside repo Y's checkout).
- **Fails:** Doesn't help within the same repo (sub-agent inside the
  correct repo but wrong sub-project). A sophisticated rogue can forge
  the fingerprint.
- **Verdict:** Raises the bar from "drive-by accident" to "active
  spoofing." Useful but not a ceiling.

### `CLAUDE_PROJECT_DIR` as Session Anchor

Use harness-provided env var as the authoritative project root.

- **Helps:** In *hook* contexts, this is reliable and solves most of the
  fragmentation. Hook-invoked `ctx` would always resolve correctly.
- **Fails:** Not available in sub-agent environments (issue #26429). Not
  available in non-hook CLI invocations. Not available for non-Claude
  tools (Cursor, Cline, Kiro, Codex).
- **Verdict:** Use where available, but cannot be the primary mechanism.

### Session Cache at `~/.ctx/sessions/<id>.json`

First hook in a session writes resolved project to a session-keyed file.
Subsequent `ctx` calls read it to stay consistent.

- **Helps:** Decouples resolution from CWD/env heuristics. Hooks and
  direct calls that know the session id agree.
- **Fails:** Sub-agent direct calls don't know the session id (no
  `CLAUDE_SESSION_ID` in env, not in stdin for non-hook calls).
  Fallbacks to PPID-chain walking get hacky fast.
- **Verdict:** Works for the part we can solve; doesn't reach sub-agent
  direct calls.

### Sub-Agents Don't Use ctx at All

Architectural rule: sub-agents report results to parent, parent persists.
No direct `ctx` calls from sub-agents.

- **Helps:** Eliminates sub-agent fragmentation entirely. Matches the
  architectural intent of sub-agents (task-scoped minions, not sessions).
  Parent has richer context to write entries anyway.
- **Fails:** Not enforceable at the ctx layer (see previous section). A
  norm, not an invariant.
- **Verdict:** Correct intent. Must be paired with observability for
  violation detection.

### Observability: Audit-Log Every Write

On any write path (task/decision/learning add, journal emit, webhook
fire), append a sidecar entry to `~/.ctx/audit.jsonl` with: resolved
context_dir, session_id (if available), PWD, PPID chain, argv, env
signature, timestamp.

A reconciliation command (`ctx audit` or session-end hook) flags writes
whose resolved context_dir diverges from the session's dominant project.

- **Helps:** Fragmentation becomes visible and attributable.
  Post-hoc migration of outlier writes becomes possible.
- **Fails:** Does not prevent wrong writes. Webhook fires still happen
  against the wrong URL — no way to recall them.
- **Verdict:** The only honest answer given in-ctx enforcement is
  impossible. Necessary floor for any design that accepts the norm.

---

## Where This Lands

Given that enforcement is not available at the ctx layer, the realistic
design shape looks like:

1. **Fix the `.ctxrc` asymmetry.** Walk up for `.ctxrc` the same way we
   walk for `.context/`. Innermost wins. Closes Scenario A entirely.
   Low cost, high value. No harness dependency.

2. **Introduce `.context/PROJECT` marker with git-remote fingerprint.**
   Unmarked `.context/` during walk-up → crash with paths listed and
   repair suggestion. Stray agent-made dirs surface immediately.
   Raises the floor on Scenarios C and E.

3. **Use `CLAUDE_PROJECT_DIR` as the anchor *in hook contexts only*.**
   Hooks are where the var is actually present. When set, skip walk and
   use it directly. Fixes sub-agent-triggered hooks cleanly — which is
   the biggest source of silent webhook loss.

4. **Adopt "sub-agents don't use ctx" as a documented norm.** Agent
   definitions and playbook guidance. Cannot be enforced; can be
   encouraged.

5. **Add audit logging for all write paths.** Sidecar JSONL capturing
   full resolution context. `ctx audit` command for post-hoc
   reconciliation. Makes fragmentation visible even when the norm is
   violated.

6. **File upstream: request `CLAUDE_PROJECT_DIR` propagation into
   sub-agent env.** A one-line patch in Claude Code would collapse most
   of (3) and (4) into a clean anchor-based solution. Until it lands,
   audit-log reconciliation is the best ceiling.

7. **For non-Claude harnesses**: per-tool session probe in the `Tool`
   config; shell-PID ancestor anchor as universal fallback; explicit
   `CTX_SESSION_ID` env as break-glass.

---

## Open Questions

1. **Is observability-only acceptable?** The original pain was *silent*
   loss of webhook events. Detection + recovery may be enough — or it
   may not, if some events are unrecoverable (external systems already
   acted on the wrong webhook). If unrecoverable, we need a stronger
   pre-write guard even at the cost of UX.

2. **Should `strict_single_root: true` be an option?** Opt-in
   configuration that makes any nested `.context/` a crash. Useful for
   users who know they have exactly one project per tree. Trivial to
   implement if (2) above lands.

3. **What breaks if `.ctxrc` walk-up is added?** Probably nothing — but
   worth checking if any test relies on the CWD-only behavior. The
   change is small; the risk is in what assumed the old behavior
   elsewhere.

4. **How do we handle `PROJECT` marker migration?** Existing `.context/`
   directories in the wild have no marker. Options: (a) treat unmarked
   as legacy-valid until a grace period ends; (b) `ctx init --migrate`
   writes markers into all detected roots. (a) is gentler; (b) is
   cleaner.

5. **What's the session-end trigger for reconciliation?** Claude Code
   has a SessionEnd hook. Non-Claude tools may not. Do we reconcile on
   every write (expensive)? On every N writes? On explicit `ctx audit`?
   Probably the latter for MVP.

6. **Does the parent agent reliably know when to persist on behalf of a
   sub-agent?** The "sub-agents don't use ctx" norm assumes parents
   internalize sub-agent results and persist them. If the parent is also
   an agent following a playbook, it can be taught. If it's the human
   driving, they need to be taught. Either way, this is a habit-change,
   not a tool feature.

---

## Appendix: Things to Verify Before Committing to This

- [ ] Does Claude Code's `SessionStart` hook payload include
  `CLAUDE_PROJECT_DIR`? (Expected yes, but unverified here.)
- [ ] Does the `Agent` tool accept a tool-allowlist parameter that
  could be used to block Bash-based `ctx` invocation per-sub-agent?
  (Worth asking upstream.)
- [ ] Does `transcript_path` or any other hook-payload field vary
  between parent-triggered and sub-agent-triggered tool events?
  (If yes, it's another discrimination signal for audit attribution.)
- [ ] For the non-Claude tools (Cursor, Cline, Kiro, Codex), what
  session-identity mechanisms exist? Need per-tool discovery before
  the per-tool probe design can be fleshed out.
- [ ] Test the `.ctxrc` walk-up fix behind a flag first; some setups
  may rely on the CWD-only behavior implicitly.

---

## References

- `internal/rc/walk.go` — walk-up implementation
- `internal/rc/rc.go:74-85` — resolution entry point (`ContextDir()`)
- `internal/rc/load.go:26-49` — `.ctxrc` loading (CWD-only read)
- `internal/entity/hook.go:15-18` — hook input contract (session id
  via stdin)
- [anthropics/claude-code#26429](https://github.com/anthropics/claude-code/issues/26429)
  — `CLAUDE_PROJECT_DIR` not propagated to sub-agents
