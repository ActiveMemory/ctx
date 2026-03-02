---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Config
icon: lucide/settings-2
---

### `ctx config`

Manage runtime configuration profiles.

```bash
ctx config <subcommand>
```

The ctx repo ships two `.ctxrc` source profiles (`.ctxrc.base` and
`.ctxrc.dev`). The working copy (`.ctxrc`) is gitignored and switched
between them using subcommands below.

#### `ctx config switch`

Switch between `.ctxrc` configuration profiles.

```bash
ctx config switch [dev|base]
```

With no argument, toggles between dev and base. Accepts `prod` as an
alias for `base`.

| Argument | Description                                |
|----------|--------------------------------------------|
| `dev`    | Switch to dev profile (verbose logging)    |
| `base`   | Switch to base profile (all defaults)      |
| *(none)* | Toggle to the opposite profile             |

**Profiles**:

| Profile | Description                                 |
|---------|---------------------------------------------|
| `dev`   | Verbose logging, webhook notifications on   |
| `base`  | All defaults, notifications off             |

**Examples**:

```bash
ctx config switch dev     # Switch to dev profile
ctx config switch base    # Switch to base profile
ctx config switch         # Toggle (dev → base or base → dev)
ctx config switch prod    # Alias for "base"
```

The detection heuristic checks for an uncommented `notify:` line in
`.ctxrc`: present means dev, absent means base.

#### `ctx config status`

Show which `.ctxrc` profile is currently active.

```bash
ctx config status
```

**Output examples**:

```
active: dev (verbose logging enabled)
active: base (defaults)
active: none (.ctxrc does not exist)
```

**See also**: [Configuration](../home/configuration.md),
[Contributing: Configuration Profiles](../home/contributing.md#configuration-profiles)
