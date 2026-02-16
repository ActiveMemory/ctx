---
title: "Permission Snapshots"
icon: lucide/camera
---

![ctx](../images/ctx-banner.png)

## The Problem

Claude Code's `.claude/settings.local.json` accumulates one-off permissions
every time you click "Allow." After busy sessions the file is full of
session-specific entries that expand the agent's surface area beyond intent.

Since `settings.local.json` is `.gitignore`d, there is no PR review or CI
check. The file drifts independently on every machine, and there is no
built-in way to reset to a known-good state.

## The Solution

Save a curated `settings.local.json` as a **golden image**, then restore
from it to drop session-accumulated permissions. The golden file
(`.claude/settings.golden.json`) is committed to version control and shared
with the team.

## Commands and Skills Used

| Command/Skill                | Role in this workflow                          |
|------------------------------|------------------------------------------------|
| `ctx permissions snapshot`   | Save settings.local.json as golden image       |
| `ctx permissions restore`    | Reset settings.local.json from golden image    |
| `/sanitize-permissions`      | Audit for dangerous patterns before snapshotting |

## Step by Step

### 1. Curate your permissions

Start with a clean `settings.local.json`. Optionally run `/sanitize-permissions`
to remove dangerous patterns first.

Review the file manually. Every entry should be there because **you** decided
it belongs, not because you clicked "Allow" once during debugging.

See the [Permission Hygiene](claude-code-permissions.md) recipe for
recommended defaults.

### 2. Take a snapshot

```bash
ctx permissions snapshot
# Saved golden image: .claude/settings.golden.json
```

This creates a byte-for-byte copy. No re-encoding, no indent changes.

### 3. Commit the golden file

```bash
git add .claude/settings.golden.json
git commit -m "Add permission golden image"
```

The golden file is **not** gitignored (unlike `settings.local.json`). This
is intentional: it becomes a team-shared baseline.

### 4. Auto-restore at session start

Add this instruction to your `CLAUDE.md`:

```markdown
## On Session Start

Run `ctx permissions restore` to reset permissions to the golden image.
```

The agent will restore the golden image at the start of every session,
automatically dropping any permissions accumulated during previous sessions.

### 5. Update when intentional changes are made

When you add a new permanent permission (not a one-off debugging entry):

```bash
# Edit settings.local.json with the new permission
# Then update the golden image:
ctx permissions snapshot
git add .claude/settings.golden.json
git commit -m "Update permission golden image: add cargo test"
```

## Conversational Approach

You don't need to remember exact commands. These natural-language prompts
work with agents trained on the ctx playbook:

| What you say                              | What happens                                |
|-------------------------------------------|---------------------------------------------|
| "Save my current permissions as baseline" | Agent runs `ctx permissions snapshot`        |
| "Reset permissions to the golden image"   | Agent runs `ctx permissions restore`         |
| "Clean up my permissions"                 | Agent runs `/sanitize-permissions` then snapshot |
| "What permissions did I accumulate?"      | Agent diffs local vs golden                  |

## Next Up

**[Keeping Context in a Separate Repo →](external-context.md)**: Share
`.context/` across projects or keep it in a dedicated repository.

## See Also

* [Permission Hygiene](claude-code-permissions.md): recommended defaults and
  maintenance workflow
* [CLI Reference: ctx permissions](../cli-reference.md#ctx-permissions):
  full command documentation
