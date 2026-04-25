---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Security Design
icon: lucide/shield-half
---

![ctx](../images/ctx-banner.png)

How `ctx` thinks about security: trust boundaries, what the system
does and does not do for you, the engineering principle behind the
audit trail, and the permission hygiene workflow.

For vulnerability disclosure, see
[Reporting Vulnerabilities](reporting.md).

## Trust Model

`ctx` operates within a single trust boundary: **the local
filesystem**.

The person who authors `.context/` files is the same person who runs
the agent that reads them. There is no remote input, no shared state,
and no server component.

This means:

* **`ctx` does not sanitize context files for prompt injection.** This
  is a deliberate design choice, not an oversight. The files are
  authored by the developer who owns the machine: sanitizing their
  own instructions back to them would be counterproductive.
* **If you place adversarial instructions in your own `.context/`
  files, your agent will follow them.** This is expected behavior.
  You control the context; the agent trusts it.

!!! warning "Shared Repositories"
    In shared repositories, `.context/` files should be reviewed in
    code review (*the same way you would review CI/CD config or
    Makefiles*). A malicious contributor could add harmful
    instructions to `CONSTITUTION.md` or `TASKS.md`.

## What `ctx` Does for Security

`ctx` is designed with security in mind:

* **No secrets in context**: The constitution explicitly forbids
  storing secrets, tokens, API keys, or credentials in `.context/`
  files.
* **Local only**: `ctx` runs entirely locally with no external
  network calls.
* **No code execution**: `ctx` reads and writes Markdown files only;
  it does not execute arbitrary code.
* **Git-tracked**: Core context files are meant to be committed, so
  they should never contain sensitive data. Exception: `sessions/`
  and `journal/` contain raw conversation data and should be
  gitignored.

## Permission Hygiene

Claude Code evaluates permissions in deny → ask → allow order.
`ctx init` automatically populates `permissions.deny` with rules
that block dangerous operations before the allow list is ever
consulted.

**Default deny rules block:**

* `sudo`, `git push`, `rm -rf /`, `rm -rf ~`, `curl`, `wget`,
  `chmod 777`
* `Read` / `Edit` of `.env`, credentials, secrets, `.pem`, `.key`
  files

Even with deny rules in place, the allow list accumulates one-off
permissions over time. Periodically review for:

* **Destructive commands**: `git reset --hard`, `git clean -f`, etc.
* **Config injection vectors**: permissions that allow modifying
  files controlling agent behavior (`CLAUDE.md`,
  `settings.local.json`).
* **Broad wildcards**: overly permissive patterns that pre-approve
  more than intended.

For the full hygiene workflow, see the
[Claude Code Permission Hygiene](../recipes/claude-code-permissions.md)
recipe.

## State File Management

Hook state files (throttle markers, prompt counters, pause markers)
are stored in `.context/state/`, which is project-scoped and
gitignored. State files are automatically managed by the hooks that
create them; no manual cleanup is needed.

## Log-First Audit Trail

The event log (`.context/state/events.jsonl`) is the authoritative
record of what `ctx` hooks did during a session. Several
audit-adjacent features depend on that log being trustworthy, not
merely best-effort:

* `ctx event` / `ctx system view-events` replays session history
  from the log.
* Webhook notifications give operators a real-time signal that
  assumes every notification corresponds to a logged event.
* Drift, freshness, and map-staleness checks count events over
  time and surface regressions.

A log that silently drops entries while the rest of the system
claims success is worse than no log at all: operators see a green
TUI and a webhook notification and conclude "it happened," even
when the audit trail never landed. The codebase treats this as a
correctness problem, not a UX polish problem.

### The Rule

> Any code path that emits an observable side effect (webhook,
> stdout marker, throttle-file touch, state mutation) must append
> the corresponding event-log entry **first** and gate the side
> effect on the append succeeding. If the log write fails, the
> side effect must not fire.

In code, this shape:

```go
if appendErr := event.Append(channel, msg, sessionID, ref); appendErr != nil {
    return appendErr // do NOT send the webhook or touch the marker
}
if sendErr := notify.Send(channel, msg, sessionID, ref); sendErr != nil {
    return sendErr
}
// downstream side effects (marker touch, stdout, etc.)
```

The `nudge.Relay` helper in `internal/cli/system/core/nudge`
enforces this for the common "log + webhook" pair. Hook `Run`
functions that compose their own sequence (`session_event`,
`heartbeat`, several `check_*` hooks) follow the same ordering
explicitly.

### Known Gaps

* **Nudge webhooks have no log channel.** `nudge.EmitAndRelay`
  sends a "nudge" notification before the "relay" event is logged.
  The nudge leg is fire-and-forget because no event-log channel
  records nudges today. A future refactor may add one; until then
  this is the one documented exception.
* **`ctx agent --cooldown` and `ctx doctor` propagate rather than
  gate.** They surface real errors to the caller (usually Cobra)
  rather than deciding what to do with them locally. Editors that
  invoke these commands may display errors in an ugly way; the
  ugliness is the correct signal (something persisted is broken),
  not a defect to smooth over.
* **Verbose hook logs in `core/log.Message` stay best-effort.**
  That logger captures per-hook activity (how many prompts, which
  percent, etc.) for debugging; it is NOT the event audit trail.
  Its failures go to stderr via `log/warn.Warn` rather than
  propagating, because losing an operational log line is not a
  correctness problem.

### Background

The `error` returns on `event.Append`, `io.AppendBytes`,
`nudge.Relay`, and `cooldown.Active` / `cooldown.TouchTombstone`
were introduced as part of the resolver-tightening refactor.
Before that change, most hook paths called these helpers and
silently discarded their errors. The principle above was extracted
from the observation that every user-visible correctness problem
hit during the refactor traced back to some function saying "this
succeeded" when the underlying write never landed.

## Best Practices

1. **Review before committing**: Always review `.context/` files
   before committing.
2. **Use `.gitignore`**: If you must store sensitive notes locally,
   add them to `.gitignore`.
3. **Drift detection**: Run `ctx drift` to check for potential
   issues.
4. **Permission audit**: Review `.claude/settings.local.json` after
   busy sessions.
