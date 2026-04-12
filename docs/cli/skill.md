---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Skill
icon: lucide/sparkles
---

![ctx](../images/ctx-banner.png)

## `ctx skill`

Manage reusable instruction bundles that can be installed into
`.context/skills/`.

A skill is a directory containing a `SKILL.md` file with YAML
frontmatter (`name`, `description`) and a Markdown instruction
body. Skills are loaded by the agent context packet when
`--skill <name>` is passed to `ctx agent`.

```bash
ctx skill <subcommand>
```

### `ctx skill install`

Install a skill from a source directory.

```bash
ctx skill install <source>
```

**Arguments**:

- `source`: Path to a directory containing `SKILL.md`

**Examples**:

```bash
ctx skill install ./my-skills/code-review
# Installed code-review → .context/skills/code-review
```

### `ctx skill list`

List all installed skills.

**Examples**:

```bash
ctx skill list
```

### `ctx skill remove`

Remove an installed skill.

**Arguments**:

- `name`: Skill name to remove

**Examples**:

```bash
ctx skill remove code-review
```

**See also**: [Building Project Skills recipe](../recipes/building-skills.md).
