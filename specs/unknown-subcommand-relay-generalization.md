# Spec: generalize the unknown-subcommand relay to `ctx hook`

**Status:** accepted (impl 2026-05-30)
**Supersedes the deferral in:** `specs/system-unknown-subcommand-relay.md`
(follow-up #1, "scoped to ctx system only")

## Problem

The unknown-subcommand relay added for `ctx system`
(`specs/system-unknown-subcommand-relay.md`) makes a missing verb fail
loud: it emits a verbatim NudgeBox, fires a best-effort event-log +
webhook relay, suppresses cobra's help dump, and exits non-zero. It was
deliberately scoped to `ctx system` because that group is the one wired
into `hooks.json`, where a help-dump-at-exit-0 is read by a
UserPromptSubmit hook as success and injected every prompt.

But the relay's value is not only the every-prompt amplification — it is
also making **drift loud instead of silent**. `ctx hook` is agent- and
script-consumed (the ctx-doctor/ctx-pause/ctx-resume skills run
`ctx hook event|message|pause|resume`; loop scripts bake in
`ctx hook notify --event loop`). Every caller invokes a *known* verb and
reads its output. If such a verb drifts out of the binary, today's
behavior is: cobra prints the `ctx hook` group help and **exits 0**. The
agent "happily ignores" it (or misparses it), the human is not notified,
and nothing is relayed. For `ctx hook notify` that is especially bad:
the loop believes it notified the human when it did not.

`ctx hook` is *user-facing* (not `Hidden`, unlike `ctx system`), so the
fix must preserve the friendly bare-`ctx hook` help while making an
**unknown** verb loud.

## Design

Lift the handler out of `internal/cli/system/core/unknown` into a neutral,
parameterized package `internal/cli/unknown`, so a group opts in by
setting its `RunE`:

```go
c.RunE = unknown.HandlerFor(unknown.SystemConfig)   // ctx system
c.RunE = unknown.HandlerFor(unknown.HookConfig)     // ctx hook
```

- `unknown.Config` carries the per-group text keys (relay prefix, box
  title, body, relay message) and the relay ref (`HookName`, `Variant`).
- `unknown.SystemConfig` reproduces the existing `ctx system` behavior
  byte-for-byte; `unknown.HookConfig` uses `ctx hook`-specific copy
  (CLI-drift framing, not hooks.json/version-skew framing) and the new
  `hook.Hook` relay label.
- Behavior is unchanged: bare group → `cmd.Help()` + exit 0; unknown verb
  → box + best-effort session-gated relay + `SilenceUsage` + non-zero exit.
- The `relay` package seam (`= nudge.Relay`) is preserved for tests.

### PreRunE exemption (the subtle part)

The root `PersistentPreRunE` (internal/bootstrap/cmd.go) early-returns for
"grouping commands without a Run/RunE." `ctx system` instead relies on its
`Hidden` early-return. `ctx hook` is visible and currently exempt **only**
via the no-RunE rule — so adding a `RunE` would newly subject bare
`ctx hook` and `ctx hook <verb>` to the context/git preconditions, a
regression (bare `ctx hook` help must work outside a project).

Fix: annotate the `ctx hook` group with `AnnotationSkipInit`. The
annotation is evaluated against the *target* command, so it exempts only
the group-level invocation (bare or unknown-verb); valid subcommands
(`event`, `pause`, …) are evaluated on their own and keep their existing
preconditions.

## Scope

- In: `ctx hook` opt-in; shared parameterized package; the move of the
  handler + its tests out of `system/core/unknown`.
- Out: a build-time guard that scans skills/loops for `ctx hook <verb>`
  references against the cobra tree (the system guard covers only
  hooks.json-wired verbs; `ctx hook` is not hooks.json-wired). Noted as a
  possible future guard, not built here.

## Why not fold into `parent.Cmd`

`parent.Cmd` is imported by many top-level groups; wiring the relay deps
(message/nudge/session) into it would widen every group's dependency
surface and force a behavior choice on groups that do not want it. A
dedicated opt-in handler keeps the choice explicit and per-group, which
is what the deferred task asked for.
