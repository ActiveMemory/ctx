# `ctx system` Unknown-Subcommand Verbatim Relay

## Problem

`ctx system <unknown>` prints the `system` group's ~51-line Long
help and exits **0**. Cobra's default `legacyArgs` raises an
"unknown command" error only for the *root* command
(`!cmd.HasParent()`); for a non-root grouping command like
`ctx system`, leftover args fall through and — because the group
has no `Run`/`RunE` — cobra calls `Help()` and returns `nil`.

In a Claude Code `UserPromptSubmit` hook, exit 0 is read by the
harness as "hook success", so the entire help blob is injected
into the agent's context **on every prompt**. This is exactly the
failure mode that the version-skew bug rode in on: a stale
`hooks.json` wired `ctx system check-anchor-drift` after the
binary deleted that command, and instead of a visible error the
agent silently ate 51 lines of help each turn (see
`specs/experiments/acdl-session-start.md` §Root Cause, Bug #2).

A non-zero exit code **alone** does not fix this: in a
`UserPromptSubmit` hook a non-zero exit with help suppressed is
swallowed by the harness — the user/agent sees nothing. The signal
has to travel on the channel the harness actually surfaces: hook
**stdout**. So the fix must *emit a message* (not just fail) when
`ctx system` is handed a verb it does not recognise.

The companion build-time guard already shipped
(`specs/hooks-wiring-guard.md`,
`TestShippedHooksResolveToRegisteredCommands`): it stops *our own*
package from shipping a hooks.json that names a missing verb. This
spec covers the **runtime** half — making a skew that reaches a
user's machine (old plugin, new binary, or vice versa) fail
*loud and legible* instead of silent.

## Approach

Give the `ctx system` command — and only that command — a `RunE`
that intercepts the unknown-subcommand case:

- **Bare `ctx system`** (no args): print help, exit 0. Unchanged
  behavior.
- **`ctx system <verb>` where `<verb>` is a real subcommand**:
  cobra descends into it and runs *its* `RunE`; `system`'s own
  `RunE` is never reached. Unchanged behavior.
- **`ctx system <unknown>`**: cobra finds no matching subcommand,
  so the leftover args reach `system`'s `RunE`. The handler emits a
  verbatim-relay box to stdout naming the unknown verb and hinting
  at the likely cause (plugin/binary version skew), sets
  `SilenceUsage`, and returns a non-nil error → exit non-zero. When
  a real session ID is present on stdin (i.e. the command was run
  by a hook), it *also* records a relay event — event-log append
  then webhook — via the same `nudge.Relay` path the `check-*`
  hooks use, so the skew is captured for out-of-band/autonomous
  alerting, not just the live session.

The behavior attaches to `system.Cmd()` after it is built by the
shared `parent.Cmd(...)`. The shared `parent.Cmd` is **not**
changed: every other group (`ctx hook`, …) keeps cobra's current
default. Scope is `ctx system` only, because that is the only group
wired into `hooks.json` and therefore the only one whose skew
pollutes a session.

Why `RunE` and not `Args: cobra.NoArgs`: `NoArgs` would surface
cobra's *generic* "unknown command" error on stderr and still emit
nothing to stdout — the harness would swallow it just like a bare
non-zero exit. We need a controlled message on stdout, which a
`RunE` gives us.

`ctx system` is registered `Hidden` (group.go `hiddenCmds`), so
`RootCmd`'s `PersistentPreRunE` early-returns for it
(`if cmd.Hidden { return nil }`) — adding a `RunE` does **not**
subject `ctx system` to the context-dir / init / git
preconditions. Verify this holds during implementation.

## Behavior

### Happy Path

1. A stale `hooks.json` (or a typo) runs
   `ctx system check-anchor-drift`.
2. Cobra finds no `check-anchor-drift` subcommand; the arg falls
   through to `system`'s `RunE` as `args = ["check-anchor-drift"]`.
3. The handler renders a verbatim-relay box to **stdout**:
   - names the unknown subcommand,
   - states the likely cause (a hook referencing a command this
     binary no longer ships — version skew between the installed
     plugin's `hooks.json` and the on-PATH `ctx` binary),
   - gives the recovery (align plugin and binary to the same
     release).
4. Best-effort: the handler reads the session ID from stdin
   (`session.ReadID`, which is TTY-safe and timeout-guarded). If a
   real session ID is present, it calls `nudge.Relay` to append a
   relay event to the local event log and fire the relay webhook
   (log-first: webhook only on a successful append).
5. The handler sets `cmd.SilenceUsage = true` and returns a
   non-nil error.
6. The harness surfaces the stdout box to the user/agent; the
   non-zero exit makes the failure debuggable instead of disguised
   as success; the event-log/webhook record gives out-of-band
   visibility. No 51-line help dump.

### Edge Cases

| Case | Expected behavior |
|------|-------------------|
| Bare `ctx system` (no args) | Print help, exit 0. Unchanged. |
| Valid hidden subcommand (`ctx system bootstrap`) | Routes to the subcommand; `system`'s `RunE` not invoked. Unchanged. |
| Valid subcommand + bad flag (`ctx system heartbeat --nope`) | Cobra's flag parsing on the *subcommand* errors as today; `system`'s `RunE` not involved. |
| Multiple leftover args (`ctx system foo bar`) | Treated as unknown; relay names the first token (`foo`). Exit non-zero. |
| Unknown verb that is a prefix of a real one (`ctx system check`) | No exact match → unknown; relay names `check`. (Cobra prefix-matching is not enabled.) |
| `ctx system --help` | Cobra's help flag short-circuits before `RunE`; prints help, exit 0. Unchanged. |
| Manual `ctx system typo` at a terminal (stdin is a TTY) | `session.ReadID` detects the char device and returns `IDUnknown` immediately — **no hang**. Stdout box + exit non-zero; relay leg (event log + webhook) is skipped (no session context to attribute it to). |
| Hook stdin present but blank/garbage session JSON | `ReadID` returns `IDUnknown` within the 2s timeout; relay leg skipped; stdout box + exit non-zero still happen. |
| Hook runs the dead command under a non-`UserPromptSubmit` event | Stdout may not be injected for every event, but the non-zero exit + stderr error + event-log/webhook relay still beat a silent exit-0 help dump; no regression. |

### Validation Rules

The "unknown" determination is cobra's, not ours: if `system`'s
`RunE` runs with `len(args) > 0`, cobra has already failed to match
a subcommand. The handler treats any non-empty `args` as an unknown
subcommand and `args[0]` as its name. No additional parsing.

### Error Handling

| Error condition | User-facing message (stdout verbatim box) | Recovery |
|-----------------|-------------------------------------------|----------|
| `ctx system <unknown>` | Verbatim-relay box: `ctx system: unknown subcommand "<verb>". A Claude Code hook is likely calling a ctx command this binary no longer ships (version skew between the installed plugin and the on-PATH ctx binary). Align the plugin and binary to the same release.` | Update plugin/binary to matching versions; or fix the hook command. |

The returned error (stderr, via `main.go`'s `writeErr`; cobra's
own printing stays silenced by `RootCmd`'s `SilenceErrors`) is a
terse companion to the stdout box — it carries the unknown verb so
a human reading logs sees the cause without the help dump.

## Interface

### CLI

No new command or flag. Changes the behavior of the existing
(hidden) `ctx system` grouping command when invoked with an
unrecognised subcommand.

```
ctx system <unknown>   # was: help + exit 0 ; now: verbatim relay + exit 1
ctx system             # unchanged: help, exit 0
ctx system bootstrap   # unchanged: runs bootstrap
```

## Implementation

### Files to Create/Modify

| File | Change |
|------|--------|
| `internal/cli/system/system.go` | After `parent.Cmd(...)`, set `c.RunE` to the unknown-subcommand handler. Keep grouping construction via `parent.Cmd`. |
| `internal/cli/system/core/.../unknown.go` (new, e.g. `core/unknown/`) | Handler: render the verbatim box, set `SilenceUsage`, return the error. Keeps `system.go` thin. |
| `internal/err/<system>` | New `UnknownSubcommand(verb string) error` constructor. Locate the existing system/CLI error file rather than creating a package. |
| `internal/assets/commands/text/*.yaml` + `internal/config/embed/text` | New `DescKey` for the relay box title and body text (no hardcoded user-facing strings, per convention). |
| `internal/cli/system/*_test.go` | Tests per Testing section. |

### Key Functions

```go
// in internal/cli/system/system.go (sketch)
c := parent.Cmd(cmd.DescKeySystem, cmd.UseSystem, /* subs... */)
c.RunE = unknown.Handler // name TBD
return c

// in the new handler (sketch)
func Handler(cmd *cobra.Command, args []string) error {
    if len(args) == 0 {
        return cmd.Help() // bare `ctx system`: unchanged
    }
    verb := args[0]
    box := message.NudgeBox(relayPrefix, title, bodyFor(verb))
    fmt.Fprintln(cmd.OutOrStdout(), box)

    // Best-effort relay leg: only when a hook supplied a session.
    // ReadID is TTY-safe and timeout-guarded, so a manual typo at a
    // terminal returns IDUnknown without blocking.
    if sid := session.ReadID(os.Stdin); sid != cfgSession.IDUnknown {
        ref := entity.NewTemplateRef(hook.System, variantUnknownSub, nil)
        if relayErr := nudge.Relay(relayMsgFor(verb), sid, ref); relayErr != nil {
            logWarn.Warn(warn.RelayUnknownSubcommand, relayErr) // log, don't mask
        }
    }

    cmd.SilenceUsage = true            // do NOT re-dump help on error
    return errSystem.UnknownSubcommand(verb)
}
```

Note: a relay-leg failure is logged via the existing warn path and
does **not** change the returned error — the user's actual problem
is the unknown subcommand, and that is what the exit reflects. The
hook name/variant for the `TemplateRef` (`hook.System` /
`variantUnknownSub`) may need new constants; reuse existing ones if
a suitable pair already exists.

### Helpers to Reuse

- `message.NudgeBox(relayPrefix, title, content)` — the same
  verbatim-box framing the `check-*` hooks use
  (`internal/cli/system/core/message`).
- `session.ReadID(os.Stdin)` — TTY-safe, timeout-guarded session-ID
  read from hook stdin; returns `cfgSession.IDUnknown` when absent
  (`internal/cli/system/core/session`).
- `nudge.Relay(msg, sessionID, ref)` — log-first event-log append +
  relay webhook (`internal/cli/system/core/nudge`).
- `entity.NewTemplateRef(hook, variant, vars)` — relay event ref for
  filtering/aggregation.
- `desc.Text(...)` — load box title/body from assets, not literals.
- `parent.Cmd` — still builds the grouping command; only the `RunE`
  is added on top.

## Configuration

None. No `.ctxrc` keys, env vars, or settings.

## Testing

- **Unit (handler):** call the handler with `args=["bogus"]`;
  assert (a) stdout contains the verbatim box, the token `bogus`,
  and the skew hint; (b) it returns a non-nil error;
  (c) `cmd.SilenceUsage` is true after the call.
- **Unit (bare):** handler with `args=[]` prints help and returns
  nil.
- **Cobra-level:** build `system.Cmd()`, execute with
  `["definitely-not-a-cmd"]`, capture output; assert non-nil error
  and the box on stdout, and assert **no** Long-help body is
  printed. Execute with `[]` (bare) and assert help + nil error.
- **Routing:** assert a known subcommand name resolves to its own
  command (e.g. `system.Cmd().Find([]string{"bootstrap"})` returns
  `bootstrap`), proving the handler is bypassed for valid verbs.
- **Relay leg — fired:** with a stdin carrying a real session ID,
  assert `nudge.Relay` records an event (event-log append observed
  in a temp state dir) for the unknown verb. Assert a relay failure
  is logged but the handler still returns the `UnknownSubcommand`
  error (relay failure does not mask it).
- **Relay leg — skipped:** with a TTY/blank stdin (`ReadID` →
  `IDUnknown`), assert no relay event is recorded and the call does
  not block; stdout box + error still happen.
- **Scope:** assert `parent.Cmd` is unchanged — another group built
  from `parent.Cmd` (e.g. `ctx hook`) still exits 0 on an unknown
  subcommand. This documents the intentional scoping and will flag
  if someone later moves the fix into the shared helper.

## Non-Goals

- **Not** changing the shared `parent.Cmd`. Other groups keep
  cobra's default unknown-subcommand behavior.
- **Not** fixing `ctx hook` or any other group. They are not wired
  into `hooks.json`, so their skew does not pollute a session. If
  that changes, revisit.
- **Not** the build-time wiring guard — already shipped
  (`specs/hooks-wiring-guard.md`). This is its runtime complement.
- **Not** removing the leftover `check-anchor-drift` prose mention
  in `commands.yaml` (acdl follow-up #3 — a separate explicit
  decision).
- **Not** generalizing the unknown-subcommand handler into the
  shared `parent.Cmd`. `ctx hook` (and any future group) has the
  same latent exit-0-on-unknown behavior, but it is not wired into
  `hooks.json`, so its skew does not pollute a session. Captured as
  a follow-up task, not built here.

## Resolved Decisions

Both decisions below were settled with the user at spec time
(session 0066d49b):

1. **Relay leg is in scope.** The handler fires the event-log +
   webhook relay (`nudge.Relay`), gated on a real session ID read
   best-effort from stdin. When no session is present (manual TTY
   typo, blank JSON), it emits the stdout box and exits non-zero
   without a relay record.
2. **Scoped to `ctx system` only.** `parent.Cmd` is untouched;
   other groups keep cobra's default. Generalization is recorded as
   a follow-up task (see Non-Goals).
