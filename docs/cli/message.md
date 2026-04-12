---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Message
icon: lucide/message-square-cog
---

![ctx](../images/ctx-banner.png)

### `ctx message`

Manage hook message templates.

Hook messages control the text hooks emit. The hook logic (when to
fire, counting, state tracking) is universal; the messages are
opinions that can be customized per-project.

```bash
ctx message <subcommand>
```

### `ctx message list`

Show all hook messages with category and override status.

```bash
ctx message list [--json]
```

**Flags**:

| Flag     | Description             |
|----------|-------------------------|
| `--json` | Output in JSON format   |

**Example**:

```bash
ctx message list
ctx message list --json | jq '.[] | select(.override)'
```

### `ctx message show`

Print the effective message template for a hook/variant pair. Shows
the user override if present, otherwise the embedded default.

```bash
ctx message show <hook> <variant>
```

**Example**:

```bash
ctx message show qa-reminder gate
ctx message show check-context-size checkpoint
```

### `ctx message edit`

Copy the embedded default template for `<hook> <variant>` to
`.context/hooks/messages/<hook>/<variant>.txt` so you can edit it
directly. The override takes effect the next time the hook fires.

```bash
ctx message edit <hook> <variant>
```

If an override already exists, the command fails and directs you to
edit it in place or reset it first.

**Example**:

```bash
ctx message edit qa-reminder gate
# Edit .context/hooks/messages/qa-reminder/gate.txt in your editor
```

### `ctx message reset`

Delete a user override and revert to the embedded default. Silent
no-op if no override exists.

```bash
ctx message reset <hook> <variant>
```

**Example**:

```bash
ctx message reset qa-reminder gate
```

See [Customizing hook messages](../recipes/customizing-hook-messages.md)
for the full workflow.
