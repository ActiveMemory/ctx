---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: Plugin Release
icon: lucide/puzzle
---

![ctx](../../images/ctx-banner.png)

# Plugin Release

Plugin-specific release procedure. The general
[release checklist](release-checklist.md) covers the full `ctx`
release; this runbook covers the plugin-specific steps that are
not part of that flow.

**When to use**: When releasing plugin changes (new skills, hook
updates, permission changes) independently of a `ctx` binary
release, or as a sub-procedure within the full release.

---

## What Ships in the Plugin

The plugin lives at `internal/assets/claude/` and includes:

| Component | Path | What it does |
|-----------|------|-------------|
| Skills | `internal/assets/claude/skills/` | User-facing `/ctx-*` slash commands |
| Hooks | `internal/assets/claude/hooks/` | Pre/post tool-use hooks |
| Plugin manifest | `internal/assets/claude/.claude-plugin/plugin.json` | Declares skills, hooks, version |
| Marketplace | `.claude-plugin/marketplace.json` | Points Claude Code to the plugin |

## Step 1: Update hooks.json (If Hooks Changed)

If you added, removed, or modified hooks:

```bash
# Verify hook definitions match implementations
make audit
```

Check that `plugin.json` lists all hooks correctly. Missing
hooks silently fail to fire.

## Step 2: Bump Version

Update the version in three places:

- `internal/assets/claude/.claude-plugin/plugin.json`
- `.claude-plugin/marketplace.json` (two fields)
- `editors/vscode/package.json` + `package-lock.json`
  (if VS Code extension is affected)

!!! tip "The Release Script Does This"
    If you're running `make release`, the script bumps these
    automatically from `VERSION`. Only bump manually if you're
    releasing the plugin independently.

## Step 3: Test Against a Fresh Install

```bash
# Clear cached plugin
make plugin-reload

# Restart Claude Code, then:
claude /plugin list    # verify version
```

Test the critical paths:

- [ ] `/ctx-status` works
- [ ] Session hooks fire (ceremonies, context loading)
- [ ] At least one user-facing skill works end-to-end
- [ ] Pre-tool-use hooks block when they should

## Step 4: Test Against a Clean Project

Create a temporary project to verify the plugin works outside
the ctx repo:

```bash
mkdir /tmp/test-ctx-plugin && cd /tmp/test-ctx-plugin
git init
ctx init
claude   # start a session, verify hooks fire
```

## Step 5: Verify Skill Count

The plugin manifest declares all user-invocable skills. Verify
the count matches:

```bash
# Count skills in plugin.json
jq '.skills | length' internal/assets/claude/.claude-plugin/plugin.json

# Count skill directories
ls -d internal/assets/claude/skills/ctx-*/ | wc -l
```

These numbers should match (some skills are not user-invocable
and won't appear in both counts).

## Step 6: Commit and Tag

If releasing independently of a binary release:

```bash
git add internal/assets/claude/ .claude-plugin/
git commit -m "chore: release plugin v0.X.Y"
git tag plugin-v0.X.Y
git push origin main --tags
```

If part of a full release, the
[release checklist](release-checklist.md) handles this.

---

## Troubleshooting

### Skills Don't Appear After Update

Claude Code caches plugin files aggressively:

```bash
make plugin-reload    # clears cache
# restart Claude Code
```

### Hooks Don't Fire

Check that the hook is registered in `plugin.json` and that
the command it calls exists:

```bash
jq '.hooks' internal/assets/claude/.claude-plugin/plugin.json
```

### Version Mismatch

If `claude /plugin list` shows an old version after updating:

```bash
make plugin-reload
# restart Claude Code
claude /plugin list   # should show new version
```
