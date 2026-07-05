# Spec: ctx statusline

## Problem

A session's cost and context pressure are invisible while it
runs. `ctx usage` reports token stats after the fact; the
context-size hook nudges at thresholds; nothing shows the
operator, continuously, what this session costs and how full
the window is. Claude Code provides exactly this data to a
configured statusline command as stdin JSON on every render.

## Design

`ctx system statusline`: a hidden system subcommand wired into
Claude Code via settings.json:

```json
"statusLine": {
  "type": "command",
  "command": "ctx system statusline"
}
```

No `refreshInterval`: Claude Code already re-renders after each
assistant message (event-driven), which is exactly when model,
context usage, and cost change. A timer would only re-spawn the
process to print the same line.

Reads the statusLine stdin JSON only. No transcript parsing,
no file walking: this runs on every render and must stay cheap.

Rendered line, all segments from stdin:

```
<user>@<host> <dir> | <model> | ctx: <N>% | $<C.CC>
```

- model from `model.display_name`
- context percentage from `context_window.used_percentage`
- cost from `cost.total_cost_usd`, printed as a plain figure

Behavior contract:

- Output sanitized to bounded printable ASCII; a malformed or
  hostile stdin payload cannot corrupt the terminal.
- Always exits zero. A broken statusline must never surface as
  a Claude Code error loop.
- Missing fields degrade to omitted segments, never to "$0.00"
  (printing a wrong number is worse than printing none).

Deploy hygiene:

- `ctx init` merges the statusLine block into
  `.claude/settings.local.json`; any pre-existing user statusline
  is backed up to `.context/state/previous-statusline.json` and
  restored when the feature is disabled. A statusLine that is not
  ctx's is never deleted on disable: it is not ours to remove.
- Disable semantics are two-stage because init refuses to re-run
  on populated projects (overwrite safety): the renderer honors
  `statusline.enabled: false` immediately by printing an empty
  line, and the settings-entry restore/removal happens whenever
  the init merge next runs.
- The settings merge is raw-map surgery: only the keys ctx owns
  are rewritten, so unmodeled settings (env, model, a user's own
  keys) survive byte-for-byte. This also fixed a pre-existing
  typed round-trip in the permissions merge that silently dropped
  unknown keys.
- `.ctxrc` block:

```yaml
# statusline:
#   enabled: true       # deploy/keep the statusline
#   show_cost: true     # omit the $ segment when false
```

`show_cost: false` exists for screen-sharing, streaming, and
recorded demos, where a running dollar figure is unwanted
on-screen: not for hiding cost from the operator.

## Decisions

**Inform, do not gate.** A rejected design escalates the line
to a spend alarm with a model-switch suggestion ("switch to a
cheaper model") whenever the live model matches a gated-family
list. ctx deliberately does not do this, on three grounds:

1. Audience: that posture fits enterprise budget enforcement
   across thousands of seats. ctx users own their bills; the
   correct posture is information, not enforcement.
2. A family-substring rule contains no task context. The
   expensive model is frequently the *cheaper* choice per
   outcome; a statusline cannot see outcomes, so it must not
   prescribe. Prescription without judgment is worse than
   silence.
3. The alarm cannot reach the case it claims to protect
   against: statuslines are seen only by present operators,
   who already see the dollar figure. Unattended overspend
   (loops, cron sessions) is addressed by notifications in
   those subsystems, not by shouting at an empty chair.

**No model suggestions anywhere in the line.** Choosing a
model is the operator's call; the line supplies the numbers
that inform it.

**Family/window extensibility without a rebuild** (a config
map teaching new model families their context-window size) is
deferred until the stdin payload proves insufficient; today
the percentage arrives pre-computed.

## Non-Goals

- Cost gating, spend alarms, or model-switch nudges
- Per-turn cost attribution (that is `ctx usage`'s job)
- Budget tracking across sessions
- Custom statusline layout DSL

## Acceptance

- [ ] `echo <sample stdin JSON> | ctx system statusline`
      renders the line; malformed JSON renders a degraded line
      and exits 0
- [ ] Cost segment absent when `cost.total_cost_usd` is
      missing and when `show_cost: false`
- [ ] `ctx setup` backs up an existing statusline before
      merging; disable restores it
- [ ] No file or transcript reads on the render path
- [ ] Output contains only bounded printable ASCII regardless
      of input
