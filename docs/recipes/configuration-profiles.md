---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Configuration Profiles
icon: lucide/sliders-horizontal
---

![ctx](../images/ctx-banner.png)

# Configuration Profiles

Switch between **dev** and **base** runtime configurations without
editing `.ctxrc` by hand. Useful when you want verbose logging and
webhook notifications during development, then clean defaults for
normal sessions.

**Uses**: `ctx config switch`, `ctx config status`, `/ctx-config`

---

## How it works

The `ctx` repo ships two source profiles committed to git:

| File           | Profile | Description                               |
|----------------|---------|-------------------------------------------|
| `.ctxrc.base`  | base    | All defaults, notifications off           |
| `.ctxrc.dev`   | dev     | Verbose logging, webhook notifications on |

The working copy (`.ctxrc`) is **gitignored**. Switching profiles
copies the source file over `.ctxrc`, so your runtime configuration
is always a clean snapshot of one of the two sources.

---

## Switching profiles

```bash
# Switch to dev (verbose logging, notifications)
ctx config switch dev

# Switch to base (defaults)
ctx config switch base

# Toggle to the opposite profile
ctx config switch

# "prod" is an alias for "base"
ctx config switch prod
```

The detection heuristic checks for an uncommented `notify:` line
in `.ctxrc`: present means dev, absent means base.

---

## Checking the active profile

```bash
ctx config status
```

Output examples:

```
active: dev (verbose logging enabled)
active: base (defaults)
active: none (.ctxrc does not exist)
```

---

## Typical workflow

1. **Start of a debugging session**: switch to dev for verbose
   logging and webhook notifications so you can trace hook
   activity and get push alerts.

   ```bash
   ctx config switch dev
   ```

2. **Work through the issue**: hooks log verbosely, webhooks fire
   on key events (commits, ceremony nudges, drift warnings).

3. **Done debugging**: switch back to base to silence the noise.

   ```bash
   ctx config switch base
   ```

---

## Customizing profiles

Edit the source files directly:

- **`.ctxrc.dev`** -- add any `.ctxrc` keys you want active during
  development (e.g., `log_level: debug`, `notify.events`,
  `notify.webhook_url`).
- **`.ctxrc.base`** -- keep this minimal. It represents your
  "production" defaults.

After editing a source file, re-run `ctx config switch <profile>`
to apply the changes to the working copy.

!!! tip "Commit your profiles"
    Both `.ctxrc.base` and `.ctxrc.dev` should be committed to git
    so team members share the same profile definitions. The working
    copy `.ctxrc` stays gitignored.

---

## Using the skill

In a Claude Code session, say any of:

- "switch to dev mode"
- "switch to base"
- "what profile am I on?"
- "toggle verbose logging"

The `/ctx-config` skill handles the rest.

**See also**: [`ctx config` reference](../cli/config.md),
[Configuration](../home/configuration.md)
