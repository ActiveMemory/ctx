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

### `ctx hook message`

Manage hook message templates.

Hook messages control the text hooks emit. The hook logic (when to
fire, counting, state tracking) is universal; the messages are
opinions that can be customized per-project.

```bash
ctx hook message <subcommand>
```

### `ctx hook message list`

Show all hook messages with category and override status.

```bash
ctx hook message list [--json]
```

**Flags**:

| Flag     | Description             |
|----------|-------------------------|
| `--json` | Output in JSON format   |

**Example**:

```bash
ctx hook message list
ctx hook message list --json | jq '.[] | select(.override)'
```

### `ctx hook message show`

Print the effective message template for a hook/variant pair. Shows
the user override if present, otherwise the embedded default.

```bash
ctx hook message show <hook> <variant>
```

**Example**:

```bash
ctx hook message show qa-reminder gate
ctx hook message show check-context-size checkpoint
```

### `ctx hook message edit`

Copy the embedded default template for `<hook> <variant>` to
`.context/hooks/messages/<hook>/<variant>.txt` so you can edit it
directly. The override takes effect the next time the hook fires.

```bash
ctx hook message edit <hook> <variant>
```

If an override already exists, the command fails and directs you to
edit it in place or reset it first.

**Example**:

```bash
ctx hook message edit qa-reminder gate
# Edit .context/hooks/messages/qa-reminder/gate.txt in your editor
```

### `ctx hook message reset`

Delete a user override and revert to the embedded default. Silent
no-op if no override exists.

```bash
ctx hook message reset <hook> <variant>
```

**Example**:

```bash
ctx hook message reset qa-reminder gate
```

See [Customizing hook messages](../recipes/customizing-hook-messages.md)
for the full workflow.
