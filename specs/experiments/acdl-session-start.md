# ACDL Experiment — Session-Start Context Assembly

**Status:** experiment / findings record (not a feature spec)
**Date:** 2026-05-27
**Author:** session 96765858 (branch `feat/pad-undo-snapshot`)
**Ref:** ACDL — Agentic Context Description Language,
<https://nogaplab.github.io/acdl-website/tutorial.html>

## Purpose

Test the thesis from the ctx-vs-ACDL discussion: does expressing
ctx's *session-start context flow* in ACDL surface doc/behavior
deltas that prose (CLAUDE.md / AGENT_PLAYBOOK.md) currently hides?

ACDL is a formal *description language* for how an LLM agent's
context evolves across turns: role markers (`S:/U:/A:/T:`),
namespaces (`env.` external inputs, `sys.` agent state, `resp.`
model output), `@T` turn-indexing, and `ForEach`/`If` control
flow. It describes; it does not store or execute. The bet was that
naming each injected fragment and its source would force latent
mismatches into the open.

**Result: thesis confirmed.** Writing the spec surfaced six deltas
and one confirmed, root-caused production bug (see §Root Cause).

## TL;DR

- ctx has **no `SessionStart` hook**. "On Session Start" in
  CLAUDE.md is unenforced prose; the real flow is a per-prompt
  (`UserPromptSubmit`) + per-tool (`PreToolUse`) push pipeline.
  The actor inverts: prose is *model-pull*, reality is *harness-push*.
- The session was **running plugin 0.8.1**, whose installed
  `hooks.json` **differs** from the repo working-tree asset. The
  first draft specced the wrong file. Spec the *installed* config.
- The every-turn ~52-line `system` help dump is a **version skew**, not
  a phantom: the installed plugin (0.8.1, *pre*-cwd-anchored
  `hooks.json`) wires `check-anchor-drift` first, but the on-PATH binary
  (0.8.1, *post*-cwd-anchored) already deleted that command. `ctx system
  <unknown>` prints help and **exits 0**, so the harness labels it
  `"hook success"` and injects it. `check-anchor-drift` was a real,
  specced, tested feature, deliberately retired — see §Root Cause.
- `check-audit` / `check-freshness` / `check-skill-discovery` are
  **`Hidden` commands by design** — their absence from
  `ctx system --help` is *not* drift (correction to first analysis).
- `ctx agent` budget disagrees across three sources: **4000 / unset /
  8000**.

---

## What "session start" means in ctx

There is no `SessionStart` event in either hooks config. "Session
start" is just `@T = 1` of a pipeline that fires on **every** user
prompt and **every** tool call:

- `UserPromptSubmit` → a stack of `ctx system check-*` hooks, each
  self-gating on `sys.*` state (most go silent after turn 1).
- `PreToolUse` → `context-load-gate` + `ctx agent` (best-effort,
  stderr-silenced), materialising only when the model first calls a
  tool — i.e. *after* it has already begun responding.

---

## Spec A — what the prose implies

`claude/CLAUDE.md` → "On Session Start": run bootstrap → read
playbook → run `ctx agent`. Modeled, that is a model-driven turn-1
sequence:

```acdl
// CLAUDE.md "On Session Start": model-driven, three ordered steps.
SessionStart_Intended[@T=1]: {
    S: PROJECT_CLAUDE_MD
    U: env.user_input[@1]
    A: {
        T: sys.run("ctx system bootstrap")     // step 1 — "CRITICAL, not optional"
        T: sys.read("AGENT_PLAYBOOK.md")        // step 2
        T: sys.run("ctx agent")                 // step 3 — content summary
        resp.reply[@1]
    }
}
```

## Spec B — what actually assembles (AS RUNNING, plugin 0.8.1)

Source: `~/.claude/plugins/cache/activememory-ctx/ctx/0.8.1/hooks/hooks.json`.
This is the config that ran *this* session.

```acdl
// ctx session-start, as actually running under plugin 0.8.1.
// @T indexes user-prompt turns. No SessionStart hook exists;
// CLAUDE.md's startup ritual is honour-system prose, unenforced.
SessionStart[@T]: {

    // Static system layer, present every turn.
    S: {
        HARNESS_SYSTEM_PROMPT
        PROJECT_CLAUDE_MD                      // claude/CLAUDE.md — prose ritual lives here
    }

    // UserPromptSubmit stack — fires EVERY @T, in declared order.
    S: {
        // ⚠ VERSION SKEW (0.8.1 plugin hooks.json:42, FIRST in stack):
        // check-anchor-drift was a real, retired feature. This pre-cwd-anchored
        // plugin config still wires it, but the post-cwd-anchored binary deleted
        // it → cobra prints `system` Long help, exits 0 → injected every turn,
        // labelled "hook success". ~52 lines of pollution. See §Root Cause.
        env.ctx_system_help_dump[@T]           // check-anchor-drift  ⚠ SKEW

        env.session_banner[@T]                 // check-context-size — "Display verbatim … Context: N% free"

        If sys.ceremony_unadopted  { env.ceremony_nudge }      // check-ceremony
        If sys.unpersisted         { env.persistence_nudge }   // check-persistence
        If sys.journal_backlog     { env.journal_reminder }    // check-journal        [fired @1, silent @2+]
        If sys.pending_reminders   { env.reminder_relay }      // check-reminder
        If sys.version_stale       { env.version_nudge }       // check-version
        If sys.resource_danger     { env.resource_warning }    // check-resource (DANGER only)
        If sys.knowledge_oversize  { env.knowledge_nudge }     // check-knowledge      [fired @1, silent @2+]
        If sys.map_stale           { env.map_nudge }           // check-map-staleness
        If sys.memory_drift        { env.memory_drift_nudge }  // check-memory-drift   [fired @1, silent @2+]
        If sys.freshness_stale     { env.freshness_nudge }     // check-freshness  (Hidden cmd)
        If sys.skill_undiscovered  { env.skill_discovery }     // check-skill-discovery (Hidden cmd)
        sys.heartbeat[@T]                                       // heartbeat (no stdout)
    }

    U: env.user_input[@T]

    // PreToolUse — does NOT exist at prompt submit. Materialises on the
    // model's first tool call, after it has started responding to @T.
    If env.tool_invoked[@T] {
        S: {
            env.context_load_directive         // context-load-gate (matcher .*)
            env.agent_packet                   // `ctx agent --budget 8000` (stderr-silenced, best-effort)
        }
        If env.tool_is_bash[@T]     { S: { env.path_guard; env.qa_reminder } }  // block-non-path-ctx; qa-reminder
        If env.tool_is_planmode[@T] { S: env.specs_nudge }                      // specs-nudge
    }

    A: resp.reply[@T]
}
```

## Spec B′ — working-tree asset divergence (post-0.8.1)

Source: `internal/assets/claude/hooks/hooks.json` (uncommitted working
tree). This is what 0.8.2 *will* install — already different from what
ran this session:

- **Drops** `check-anchor-drift` — retired by the cwd-anchored migration
  (`fc7db228`), since the drift it detected cannot occur once `CTX_DIR`
  is gone.
- **Adds** `check-audit` after `check-reminder` (out-of-band audit relay).
- **Invocation style** changed: `cd "$CLAUDE_PROJECT_DIR" && ctx system …`
  (asset, cwd-anchored) vs `CTX_DIR="…/.context" ctx system …` (0.8.1,
  pre-migration). This `cd` migration *is* `fc7db228` Step 3.

The installed-vs-asset gap is itself a finding (delta #7).

---

## The delta (amended, evidence-backed)

Certainty: **[confirmed]** = read from source/config; **[inferred]** =
reasoned, not directly verified.

1. **Actor inversion — no `SessionStart` hook. [confirmed]**
   CLAUDE.md tells the *model* to run bootstrap → playbook → `ctx
   agent`. Neither hooks config has a `SessionStart` event; none of
   the three steps is enforced. Reality is harness-push (~13–14
   `UserPromptSubmit` hooks + `PreToolUse`). ACDL exposes this because
   the prose steps have no hook to bind to — they only exist as
   model-driven `A:` actions that don't match the push reality.

2. **`ctx agent` timing + triple-valued budget. [confirmed]**
   Prose implies an explicit startup step; actually a `PreToolUse .*`
   hook (`ctx agent --budget 8000 2>/dev/null || true`) that lands on
   the *first tool call*, after the model is already replying. Budget
   disagrees three ways:
   - `WORKSPACE/CLAUDE.md` (outer): `--budget 4000`
   - `ctx/CLAUDE.md` (inner): bare `ctx agent`
   - both hooks.json: `--budget 8000`

3. **Command surface — corrected. [confirmed]**
   First analysis called `check-audit` / `check-freshness` /
   `check-skill-discovery` "wired but undocumented drift." **Wrong.**
   They are registered **`Hidden` cobra subcommands**
   (`internal/cli/system/cmd/check{audit,freshness,skilldiscovery}/doc.go`
   each say "hidden hook"), so their omission from `ctx system --help`
   is by design. The `check-anchor-drift` line in the hand-written help
   prose (`internal/assets/commands/commands.yaml:1165`) is **not** a
   never-real phantom — it is a *leftover* from a retired feature.
   `check-anchor-drift` was specced (`single-source-context-anchor.md`
   §F), implemented and tested (`internal/cli/system/cmd/checkanchordrift/`,
   incl. a 153-line `run_test.go`), then deliberately deleted in
   `fc7db228` (cwd-anchored migration). That commit pruned the command
   *key* from `commands.yaml` but missed this *prose* mention in the
   `system` `long:` block. See §Root Cause.

4. **Undocumented verbatim banner. [confirmed]**
   The mandatory `Session: … | Branch: … | Context: N% free` line is
   emitted by `check-context-size` but documented in no CLAUDE.md.
   Minor, but ACDL forced it a name (`env.session_banner`) and a
   source.

5. **Every-turn help dump — root-caused to a version skew. [confirmed]**
   Promoted from "candidate." See §Root Cause. Not a phantom: a real
   retired feature (`check-anchor-drift`) whose pre-cwd-anchored plugin
   `hooks.json` outlived the binary that deleted it. The prose layer had
   no way to reveal this — only git history + the spec did.

6. **"Session start" is a misnomer (@T gating). [confirmed]**
   These are per-prompt, self-gating hooks, not one-time startup
   context. Observed: memory-drift / knowledge / journal fired at @1,
   went silent @2+. ACDL's `@T` + `If` captures the conditional
   every-turn shape prose flattens.

7. **Installed ≠ repo (version skew). [confirmed]**
   The session ran plugin/binary **0.8.1**, but the spec was first
   drafted against the **working-tree asset** (post-0.8.1, already
   divergent). The retired `check-anchor-drift` hook polluting this
   session wasn't even in the asset. Lesson: when specifying agent
   context flows, spec the *installed* config, not the repo template.

---

## Root Cause — the every-turn `ctx system` help dump

Two stacked defects:

### Bug #1 — Version skew across the cwd-anchored migration

**Not a phantom — a retired feature, half-installed.** `check-anchor-drift`
was a real diagnostic: specced in `single-source-context-anchor.md` §F,
implemented and tested in `internal/cli/system/cmd/checkanchordrift/`
(`cmd.go`/`doc.go`/`run.go` + a 153-line `run_test.go`), backed by
`internal/cli/system/core/anchor/` (symlink-aware `anchor.Equal`). Its
job: warn when the user's shell-level `CTX_DIR` (from `ctx activate`)
diverged from the hook-injected `$CLAUDE_PROJECT_DIR/.context`.

It was **deliberately deleted** in `fc7db228` (2026-05-22, *"feat(cwd-anchored):
drop CTX_DIR + ctx activate/deactivate"*, implementing
`specs/cwd-anchored-context.md`). Commit message, Step 3, verbatim:
*"check-anchor-drift hook entry removed (the drift it detected cannot
occur under cwd-anchored)."* Under cwd-anchoring there is no `CTX_DIR`
channel, so the two values it compared no longer exist; the failure
mode it guarded is now a resolver-level hard refusal
(`errCtx.NoContextHere($PWD)`) — strictly stronger.

**Why it still fires here:** `fc7db228` is in **no 0.8.x tag** (only
`v0.8.0` exists); the migration is unreleased. This machine runs a
*post*-migration binary (0.8.1, no `check-anchor-drift`) against a
*pre*-migration installed plugin (`0.8.1/hooks/hooks.json:42`, still
injecting `CTX_DIR=` and wiring `check-anchor-drift` first). The hook
calls a command the binary correctly dropped → help dump. The working
tree is already consistent (asset `hooks.json` is `cd`-based, no
`check-anchor-drift`); the fix is to **release/reinstall so plugin and
binary come from the same post-`fc7db228` commit** — nothing is removed.
A prose leftover remains at `commands.yaml:1165` (the dumped blob);
left in place pending an explicit decision.

### Bug #2 — Silent-help-on-unknown (latent; still present)

`internal/cli/parent/parent.go:25` (`parent.Cmd`) builds the `system`
group with **no `Run`/`RunE` and no `Args` validator**:

```go
c := &cobra.Command{Use: use, Short: short, Long: long, Example: …}
c.AddCommand(subs...)
```

Cobra's default `legacyArgs` raises `unknown command` **only for the
root command** (`!cmd.HasParent()`). For a non-root group like `ctx
system`, an unknown subcommand falls through to `Help()` and returns
`nil` → **exit 0**. The harness sees exit 0, labels the hook
`"success"`, and injects the full ~52-line Long help into context.

Empirically confirmed on the 0.8.1 binary:

```
$ CTX_DIR=.../.context ctx system check-anchor-drift ; echo EXIT=$?
Hook plumbing namespace. Hosts Claude Code hook logic …   (52 lines)
EXIT=0
```

**This is the insidious one:** it converts any one-line hook-name typo
into invisible, every-turn context pollution disguised as success, and
it applies to **every** `parent.Cmd` group (`system`, `hook`, …) — so
future renamed/removed hook commands will silently regress the same way.

### Proposed fix (bug #2)

Give grouping commands an unknown-subcommand guard so they exit
non-zero instead of printing help. Minimal:

```go
c := &cobra.Command{Use: use, Short: short, Long: long, Example: …,
    Args: cobra.NoArgs,   // unknown subcommand → "unknown command", exit 1
}
```

`cobra.NoArgs` does not affect valid subcommands (cobra descends into a
matched subcommand before applying the parent's `Args`). It only fires
when the parent itself receives leftover args — i.e. an unknown
subcommand — turning silent pollution into a loud, debuggable hook
failure. Needs its own spec + branch; not done in this experiment.

---

## Verdict

The ACDL exercise paid off: ~40 lines of notation surfaced two
doc-vs-reality mismatches (#1, #2), one corrected surface claim (#3),
an undocumented invariant (#4), a confirmed root-caused bug (#5), a
conceptual misframing (#6), and an installed-vs-repo skew (#7) — none
visible from the prose layer.

The counter-point from the design discussion holds: this is a second
artifact to keep in sync. Its value is conditional on it being either
(a) checkable against the *installed* `hooks.json`, or (b) treated as
the audit target for `_ctx-surface-audit` / `_ctx-alignment-audit` —
not maintained by hand.

## Recommended follow-ups

Nothing here removes the `check-anchor-drift` feature or its history —
it was already correctly retired by `fc7db228`. These finish the
delivery and harden the surface.

1. **Resolve the version skew (live fix).** Cut/install a release where
   the plugin's `hooks.json` and the on-PATH binary come from the same
   post-`fc7db228` commit, so the pre-cwd-anchored plugin config stops
   calling a command the binary deleted. *Stops the every-turn dump.*
2. **Harden `parent.Cmd` against unknown subcommands (bug #2).**
   `Args: cobra.NoArgs` (or equivalent) so `ctx system <bogus>` exits
   non-zero instead of printing help + exit 0; add a regression test.
   Independent of the skew — prevents *any* future typo'd hook name
   from silently polluting context. Own spec + branch.
3. **Decide on the `commands.yaml:1165` prose leftover.** `fc7db228`
   pruned the command key but left the textual mention in the `system`
   `long:` block. Keep (historical) or prune — explicit decision, not a
   reflex deletion.
4. **Reconcile the `ctx agent` budget** across the two CLAUDE.md files
   and the hook (4000 / unset / 8000 — pick one or document why).
5. Consider whether the hidden `check-*` commands belong in the
   hand-written help prose for discoverability, or stay hidden
   consistently.
