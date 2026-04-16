---
#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0

title: New Contributor
icon: lucide/user-plus
---

![ctx](../../images/ctx-banner.png)

# New Contributor Onboarding

Step-by-step onboarding sequence for new contributors. Consolidates
setup instructions currently scattered across the README,
[contributing guide](../../home/contributing.md), and setup docs.

**When to use**: First-time contributor setup, or when verifying
your development environment after a major upgrade.

---

## Step 1: Clone the Repository

```bash
git clone https://github.com/ActiveMemory/ctx.git
cd ctx
```

Or fork first on GitHub, then clone your fork.

## Step 2: Initialize Context

```bash
ctx init
```

This creates the `.context/` directory with knowledge files and
the `.claude/` directory with agent configuration. If `ctx` is
not yet installed, proceed to Step 3 first, then come back.

## Step 3: Build and Install

```bash
make build
sudo make install
```

Verify:

```bash
ctx --version
```

## Step 4: Install the Plugin (Claude Code Users)

If you use Claude Code, install the plugin from your local clone
so skills and hooks reflect your working tree:

1. Launch `claude`
2. Type `/plugin` and press Enter
3. Select **Marketplaces** -> **Add Marketplace**
4. Enter the absolute path to your clone (e.g., `~/WORKSPACE/ctx`)
5. Back in `/plugin`, select **Install** and choose `ctx`

Verify:

```bash
claude /plugin list   # should show ctx
```

See [Contributing: Install the Plugin](../../home/contributing.md#3-install-the-plugin-from-your-local-clone)
for details on cache clearing.

## Step 5: Switch to Dev Profile

```bash
ctx config switch dev
```

This enables verbose logging and notify events (useful during
development).

## Step 6: Verify Hooks

Start a Claude Code session and check that hooks fire:

```bash
claude
```

You should see ctx session hooks (ceremonies reminder, context
loading) on session start. If not, check that the plugin is
installed correctly (Step 4).

## Step 7: Run Your First Session

In Claude Code:

```
/ctx-status
```

This should show context file health, active tasks, and recent
decisions. If it works, your setup is complete.

## Step 8: Verify Context Persistence

End the session and start a new one:

```
/ctx-remember
```

The agent should recall what happened in the previous session.
This confirms that context persistence is working end-to-end.

## Step 9: Run Tests

```bash
make test     # unit tests
make audit    # full check: fmt + vet + lint + drift + docs + test
```

All tests should pass with a clean clone.

---

## Quick Reference

| Task | Command |
|------|---------|
| Build | `make build` |
| Install | `sudo make install` |
| Test | `make test` |
| Full audit | `make audit` |
| Rebuild docs site | `make site` |
| Serve docs locally | `make site-serve` |
| Clear plugin cache | `make plugin-reload` |
| Switch config profile | `ctx config switch dev` |

## Next Steps

- Read the [contributing guide](../../home/contributing.md)
  for project layout, code style, and PR process
- Check [TASKS.md](https://github.com/ActiveMemory/ctx/blob/main/.context/TASKS.md)
  for open work items
- Ask `/ctx-next` for suggested work
