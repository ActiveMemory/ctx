---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Breaking Migration
icon: lucide/arrow-right-left
---

![ctx](../../images/ctx-banner.png)

# Breaking Migration Guide

Template for upgrading across breaking CLI renames or behavior
changes. Use this as a starting point when writing migration
notes for a specific release, or hand it to your agent as
context for generating release-specific guidance.

**When to use**: When a release includes breaking changes
(command renames, removed flags, changed defaults) that require
user action.

**Companion**: [Upgrade guide](../upgrading.md) covers the
general upgrade flow. This runbook covers the breaking-change
specifics.

---

## Step 1: Identify What Changed

Ask your agent to diff the CLI surface between the old and new
version:

```
Compare the CLI command surface between the previous release tag
and HEAD. For each change, categorize as: renamed, removed,
new, or changed-behavior. Include old and new command signatures.
```

Or use the `/_ctx-command-audit` skill after the rename.

## Step 2: Regenerate Infrastructure

```bash
# Install the new binary
make build && sudo make install

# Regenerate CLAUDE.md and permissions
ctx init --reset --merge
```

`--merge` preserves your knowledge files (TASKS.md, DECISIONS.md,
etc.) while regenerating infrastructure (permissions, CLAUDE.md
managed sections).

## Step 3: Update the Plugin

```
/plugin -> select ctx -> Update now
```

Or, if using a local clone:

```bash
make plugin-reload
# restart Claude Code
```

## Step 4: Update Personal Scripts

Search your scripts and aliases for old command names:

```bash
# Example: find references to old command names
grep -r "ctx old-command" ~/scripts/ ~/.zshrc ~/.bashrc
```

Replace with the new names per the changelog.

## Step 5: Update Hook Configs

If you have custom hooks in `.claude/settings.local.json` that
reference ctx commands, update them:

```bash
jq '.hooks' .claude/settings.local.json | grep "ctx "
```

## Step 6: Verify

Activate the project first, otherwise `ctx status` and `ctx drift`
will fail with `Error: no context directory specified`:

```bash
eval "$(ctx activate)"
ctx status          # context files intact
ctx drift           # no broken references
make test           # if you're a contributor
```

See [Activating a Context Directory](../../recipes/activating-context.md).

---

## Writing Release-Specific Migration Notes

When preparing a release with breaking changes, create a section
in the release notes using this template:

```markdown
## Breaking Changes

### `old-command` renamed to `new-command`

**What changed**: `ctx old-command` is now `ctx new-command`.
The old name is removed (no deprecation alias).

**Action required**:
1. Run `ctx init --reset --merge` to update CLAUDE.md
2. Update any scripts referencing `ctx old-command`
3. Update hook configs if applicable

**Why**: [brief rationale for the rename]
```

Repeat for each breaking change. Users should be able to follow
the notes mechanically without needing to understand the
codebase.
